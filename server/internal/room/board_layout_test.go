package room

import (
	"encoding/json"
	"math/rand"
	"strings"
	"testing"

	"server/internal/locale"
)

func TestCellTypeSpecs_MatchDispatchTable(t *testing.T) {
	if len(cellTypeSpecs) != len(EffectDispatchTable) {
		t.Fatalf("cell type catalog has %d types, dispatch table has %d", len(cellTypeSpecs), len(EffectDispatchTable))
	}
	for _, spec := range cellTypeSpecs {
		if _, ok := EffectDispatchTable[spec.Type]; !ok {
			t.Errorf("cell type %q has no effect handler", spec.Type)
		}
	}
}

func TestValidateBoard_DefaultBoardIsValid(t *testing.T) {
	cells, err := ValidateBoard(DefaultBoard())
	if err != nil {
		t.Fatalf("default board rejected: %v", err)
	}
	if cells[0].Params["amount"] != 100 || cells[0].Params["pass_bonus"] != 50 {
		t.Errorf("deploy params = %v, want amount 100 and pass_bonus 50", cells[0].Params)
	}
	if _, ok := cells[1].Params["size"]; ok {
		t.Errorf("legacy size param was kept: %v", cells[1].Params)
	}
}

func TestValidateBoard_OrdersCellsFromJSON(t *testing.T) {
	board := DefaultBoard()
	board[3], board[4] = board[4], board[3]
	raw, _ := json.Marshal(board)
	var decoded []BoardCell
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	cells, err := ValidateBoard(decoded)
	if err != nil {
		t.Fatalf("shuffled board rejected: %v", err)
	}
	for i, cell := range cells {
		if cell.Index != i {
			t.Fatalf("cell at position %d has index %d", i, cell.Index)
		}
	}
	if cells[1].Params["amount"] != 10 {
		t.Errorf("JSON amount not normalized to int: %#v", cells[1].Params["amount"])
	}
}

func TestValidateBoard_Rejections(t *testing.T) {
	cases := []struct {
		name   string
		mutate func([]BoardCell) []BoardCell
		want   string
	}{
		{"too few cells", func(b []BoardCell) []BoardCell { return b[:31] }, "exactly 32 cells"},
		{"duplicate index", func(b []BoardCell) []BoardCell { b[5].Index = 4; return b }, "more than once"},
		{"index out of range", func(b []BoardCell) []BoardCell { b[5].Index = 32; return b }, "out of range"},
		{"unknown type", func(b []BoardCell) []BoardCell { b[5].Type = "skip_turn"; return b }, "unknown cell type"},
		{"start cell not deploy", func(b []BoardCell) []BoardCell { b[0].Type = "mystery"; b[0].Params = nil; return b }, "start cell must be deploy"},
		{"second deploy", func(b []BoardCell) []BoardCell { b[5].Type = "deploy"; b[5].Params = nil; return b }, "only be placed on the start cell"},
		{"blank name", func(b []BoardCell) []BoardCell { b[5].Name = locale.New("  ", "x", "y"); return b }, "name is required"},
		{"unknown param", func(b []BoardCell) []BoardCell { b[1].Params = map[string]any{"amount": 10, "size": "S"}; return b }, "unknown parameter"},
		{"fractional param", func(b []BoardCell) []BoardCell { b[1].Params = map[string]any{"amount": 10.5}; return b }, "must be an integer"},
		{"param out of range", func(b []BoardCell) []BoardCell { b[1].Params = map[string]any{"amount": 0}; return b }, "between 1 and 500"},
		{"teleport to self", func(b []BoardCell) []BoardCell { b[12].Params = map[string]any{"target_position": 12}; return b }, "its own cell"},
		{"teleport off board", func(b []BoardCell) []BoardCell { b[12].Params = map[string]any{"target_position": 40}; return b }, "between 0 and 31"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateBoard(tc.mutate(DefaultBoard()))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want containing %q", err, tc.want)
			}
		})
	}
}

func TestNormalizeCellParams_DefaultsAndClamping(t *testing.T) {
	got := NormalizeCellParams("deadline", map[string]any{"gain": float64(900), "stale": 1})
	if got["gain"] != maxCellXP || got["loss"] != 20 || len(got) != 2 {
		t.Errorf("deadline params = %v", got)
	}
	if got := NormalizeCellParams("xp_loss", map[string]any{"size": "M"}); got["amount"] != 25 {
		t.Errorf("legacy size not folded into amount: %v", got)
	}
	if got := NormalizeCellParams("mystery", map[string]any{"amount": 5}); len(got) != 0 {
		t.Errorf("mystery should take no params: %v", got)
	}
}

func TestConfigurableEffects_UseCellParams(t *testing.T) {
	r := NewRoom("test-room-config", &MockBroadcaster{})
	p := NewPlayer("c1", "Tester")

	deploy := BoardCell{Type: "deploy", Params: map[string]any{"amount": 70}}
	if res := r.ApplyCellEffect(p, deploy); res.XPDelta != 70 {
		t.Errorf("deploy landing bonus = %d, want 70", res.XPDelta)
	}
	coffee := BoardCell{Type: "coffee_break", Params: map[string]any{"amount": float64(12)}}
	if res := r.ApplyCellEffect(p, coffee); res.XPDelta != 12 {
		t.Errorf("coffee break = %d, want 12", res.XPDelta)
	}
	if res := r.ApplyCellEffect(p, BoardCell{Type: "coffee_break"}); res.XPDelta != 5 {
		t.Errorf("coffee break default = %d, want 5", res.XPDelta)
	}

	deadline := BoardCell{Type: "deadline", Params: map[string]any{"gain": 80, "loss": 30}}
	seenGain, seenLoss := false, false
	r.SetRNG(rand.New(rand.NewSource(1)))
	for i := 0; i < 20; i++ {
		p.XP = 100
		switch res := r.ApplyCellEffect(p, deadline); res.XPDelta {
		case 80:
			seenGain = true
		case -30:
			seenLoss = true
		default:
			t.Fatalf("deadline delta = %d, want +80 or -30", res.XPDelta)
		}
	}
	if !seenGain || !seenLoss {
		t.Errorf("deadline outcomes gain=%v loss=%v, want both", seenGain, seenLoss)
	}
}

func TestSetBoard_PassBonusFromStartCell(t *testing.T) {
	board := DefaultBoard()
	board[0].Params = map[string]any{"amount": 100, "pass_bonus": 75}
	for i := 1; i < BoardSize; i++ {
		board[i].Type = "skip_next"
		board[i].Params = map[string]any{}
	}
	r := NewRoom("test-room-pass", &MockBroadcaster{})
	r.SetBoard(BoardLayout{ID: "b1", Name: "Custom", Cells: board})
	if r.BoardName != "Custom" || r.Board()[0].Params["pass_bonus"] != 75 {
		t.Fatalf("board not applied: name=%q start=%v", r.BoardName, r.Board()[0].Params)
	}

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AdminStartGame()
	r.mu.Lock()
	player := r.playerMap["c1"]
	player.Position = BoardSize - 1
	player.ChosenDifficulty = "easy"
	rolls := r.rollPlayerLocked(player)
	r.mu.Unlock()

	if len(rolls) != 1 || !rolls[0].PassedGO || rolls[0].LapBonus != 75 {
		t.Fatalf("roll = %+v, want passed GO with lap bonus 75", rolls)
	}
}

func TestBeginnerBoard_IsValidAndGentlerThanDefault(t *testing.T) {
	cells, err := ValidateBoard(BeginnerBoard())
	if err != nil {
		t.Fatalf("beginner board rejected: %v", err)
	}
	for _, cell := range cells {
		if cell.Name.Ru == "" || cell.Name.Kz == "" || cell.Explanation.En == "" || cell.Explanation.Ru == "" || cell.Explanation.Kz == "" {
			t.Errorf("cell %d is missing translated text: %+v", cell.Index, cell)
		}
	}

	type tally struct{ negative, positive, teleport int }
	count := func(board []BoardCell) tally {
		var c tally
		for _, cell := range board {
			switch cell.Type {
			case "xp_loss", "skip_next", "code_freeze", "deadline":
				c.negative++
			case "xp_gain", "special_challenge", "double_xp", "free_pass", "coffee_break":
				c.positive++
			case "teleport":
				c.teleport++
			}
		}
		return c
	}
	beginner, standard := count(cells), count(DefaultBoard())
	if beginner.negative >= standard.negative || beginner.positive <= standard.positive || beginner.teleport >= standard.teleport {
		t.Fatalf("beginner %+v is not gentler than default %+v", beginner, standard)
	}
}
