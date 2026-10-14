package room

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"server/internal/locale"
)

const (
	// BoardSize is the fixed number of perimeter cells on every board.
	BoardSize = 32
	// StartCellIndex is the Deploy (GO) cell; passing it pays the lap bonus.
	StartCellIndex = 0

	maxCellXP             = 500
	maxCellNameLength     = 255
	maxCellExplanationLen = 2000
)

// BoardLayout is a named board configuration a room is created with.
type BoardLayout struct {
	ID    string
	Name  string
	Cells []BoardCell
}

// CellParamSpec describes one integer parameter a cell type reads from Params.
type CellParamSpec struct {
	Key     string `json:"key"`
	Default int    `json:"default"`
	Min     int    `json:"min"`
	Max     int    `json:"max"`
}

// CellTypeSpec describes a cell type that can be placed on a board.
type CellTypeSpec struct {
	Type      string          `json:"type"`
	StartOnly bool            `json:"start_only"`
	Params    []CellParamSpec `json:"params"`
}

// cellTypeSpecs is the catalog of placeable cell types, in display order.
// Every entry must have a handler in EffectDispatchTable.
var cellTypeSpecs = []CellTypeSpec{
	{Type: "deploy", StartOnly: true, Params: []CellParamSpec{
		{Key: "amount", Default: 100, Min: 0, Max: maxCellXP},
		{Key: "pass_bonus", Default: 50, Min: 0, Max: maxCellXP},
	}},
	{Type: "xp_gain", Params: []CellParamSpec{{Key: "amount", Default: 10, Min: 1, Max: maxCellXP}}},
	{Type: "xp_loss", Params: []CellParamSpec{{Key: "amount", Default: 10, Min: 1, Max: maxCellXP}}},
	{Type: "special_challenge", Params: []CellParamSpec{{Key: "bonus", Default: 30, Min: 1, Max: maxCellXP}}},
	{Type: "double_xp", Params: []CellParamSpec{}},
	{Type: "free_pass", Params: []CellParamSpec{}},
	{Type: "skip_next", Params: []CellParamSpec{}},
	{Type: "teleport", Params: []CellParamSpec{{Key: "target_position", Default: 0, Min: 0, Max: BoardSize - 1}}},
	{Type: "mystery", Params: []CellParamSpec{}},
	{Type: "code_freeze", Params: []CellParamSpec{}},
	{Type: "coffee_break", Params: []CellParamSpec{{Key: "amount", Default: 5, Min: 0, Max: maxCellXP}}},
	{Type: "deadline", Params: []CellParamSpec{
		{Key: "gain", Default: 50, Min: 0, Max: maxCellXP},
		{Key: "loss", Default: 20, Min: 0, Max: maxCellXP},
	}},
}

// CellTypeSpecs returns the catalog of placeable cell types and their parameters.
func CellTypeSpecs() []CellTypeSpec {
	return cellTypeSpecs
}

func findCellTypeSpec(cellType string) (CellTypeSpec, bool) {
	for _, spec := range cellTypeSpecs {
		if spec.Type == cellType {
			return spec, true
		}
	}
	return CellTypeSpec{}, false
}

// intParam reads an integer parameter; JSON numbers decode as float64.
func intParam(params map[string]any, key string) (int, bool) {
	switch v := params[key].(type) {
	case int:
		return v, true
	case float64:
		if v == math.Trunc(v) && math.Abs(v) <= math.MaxInt32 {
			return int(v), true
		}
	}
	return 0, false
}

// cellParam returns the integer parameter key of cell, falling back to the
// cell type's default when it is absent.
func cellParam(cell BoardCell, key string) int {
	if value, ok := intParam(cell.Params, key); ok {
		return value
	}
	spec, _ := findCellTypeSpec(cell.Type)
	for _, param := range spec.Params {
		if param.Key == key {
			return param.Default
		}
	}
	return 0
}

// NormalizeCellParams reduces params to the keys cellType reads, filling in
// defaults, folding legacy keys (xp "size"), and clamping to the allowed range.
func NormalizeCellParams(cellType string, params map[string]any) map[string]any {
	out := map[string]any{}
	spec, ok := findCellTypeSpec(cellType)
	if !ok {
		return out
	}
	for _, param := range spec.Params {
		value, ok := intParam(params, param.Key)
		if !ok {
			value = param.Default
			if param.Key == "amount" {
				value = parseAmount(params, param.Default)
			}
		}
		out[param.Key] = min(max(value, param.Min), param.Max)
	}
	return out
}

// ValidateBoard checks that cells form a playable board and returns them
// ordered by index with text trimmed and params normalized. Every index
// 0..BoardSize-1 must appear exactly once, the start cell must be Deploy (and
// only it), English names are required, and params must be integers within
// their type's range.
func ValidateBoard(cells []BoardCell) ([]BoardCell, error) {
	if len(cells) != BoardSize {
		return nil, fmt.Errorf("board must have exactly %d cells", BoardSize)
	}
	ordered := make([]BoardCell, BoardSize)
	seen := make([]bool, BoardSize)
	for _, cell := range cells {
		if cell.Index < 0 || cell.Index >= BoardSize {
			return nil, fmt.Errorf("cell index %d is out of range", cell.Index)
		}
		if seen[cell.Index] {
			return nil, fmt.Errorf("cell %d appears more than once", cell.Index)
		}
		seen[cell.Index] = true
		cell.Name = trimText(cell.Name)
		cell.Explanation = trimText(cell.Explanation)
		if err := validateCell(cell); err != nil {
			return nil, fmt.Errorf("cell %d: %w", cell.Index, err)
		}
		cell.Params = NormalizeCellParams(cell.Type, cell.Params)
		ordered[cell.Index] = cell
	}
	return ordered, nil
}

func validateCell(cell BoardCell) error {
	spec, ok := findCellTypeSpec(cell.Type)
	if !ok {
		return fmt.Errorf("unknown cell type %q", cell.Type)
	}
	if cell.Index == StartCellIndex && !spec.StartOnly {
		return errors.New("the start cell must be deploy")
	}
	if spec.StartOnly && cell.Index != StartCellIndex {
		return fmt.Errorf("%s can only be placed on the start cell", cell.Type)
	}
	if cell.Name.En == "" {
		return errors.New("name is required")
	}
	for _, name := range []string{cell.Name.En, cell.Name.Ru, cell.Name.Kz} {
		if utf8.RuneCountInString(name) > maxCellNameLength {
			return fmt.Errorf("name exceeds %d characters", maxCellNameLength)
		}
	}
	for _, explanation := range []string{cell.Explanation.En, cell.Explanation.Ru, cell.Explanation.Kz} {
		if utf8.RuneCountInString(explanation) > maxCellExplanationLen {
			return fmt.Errorf("explanation exceeds %d characters", maxCellExplanationLen)
		}
	}
	for key := range cell.Params {
		known := false
		for _, param := range spec.Params {
			known = known || param.Key == key
		}
		if !known {
			return fmt.Errorf("unknown parameter %q for %s", key, cell.Type)
		}
	}
	for _, param := range spec.Params {
		if _, present := cell.Params[param.Key]; !present {
			continue
		}
		value, ok := intParam(cell.Params, param.Key)
		if !ok {
			return fmt.Errorf("parameter %q must be an integer", param.Key)
		}
		if value < param.Min || value > param.Max {
			return fmt.Errorf("parameter %q must be between %d and %d", param.Key, param.Min, param.Max)
		}
	}
	if cell.Type == "teleport" && cellParam(cell, "target_position") == cell.Index {
		return errors.New("teleport cannot target its own cell")
	}
	return nil
}

func trimText(text locale.Text) locale.Text {
	return locale.New(strings.TrimSpace(text.En), strings.TrimSpace(text.Ru), strings.TrimSpace(text.Kz))
}
