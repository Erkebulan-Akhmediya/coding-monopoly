package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"server/internal/boards"
	"server/internal/room"
)

const maxBoardNameLength = 128

type boardInput struct {
	Name  string           `json:"name"`
	Cells []room.BoardCell `json:"cells"`
}

// validate trims the name and replaces Cells with the validated, ordered,
// normalized cells that should be stored.
func (input *boardInput) validate() error {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return errors.New("board name is required")
	}
	if utf8.RuneCountInString(input.Name) > maxBoardNameLength {
		return fmt.Errorf("board name exceeds maximum length of %d characters", maxBoardNameLength)
	}
	cells, err := room.ValidateBoard(input.Cells)
	if err != nil {
		return err
	}
	input.Cells = cells
	return nil
}

func (h *Handler) listCellTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cell_types": room.CellTypeSpecs(), "board_size": room.BoardSize})
}

func (h *Handler) boards(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/boards"), "/")
	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.listBoards(w, r)
		case http.MethodPost:
			h.createBoard(w, r)
		default:
			methodNotAllowed(w)
		}
		return
	}
	parts := strings.Split(path, "/")
	if len(parts) == 2 && parts[1] == "default" {
		if r.Method == http.MethodPost {
			h.setDefaultBoard(w, r, parts[0])
			return
		}
		methodNotAllowed(w)
		return
	}
	if len(parts) != 1 {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getBoard(w, r, parts[0])
	case http.MethodPut:
		h.updateBoard(w, r, parts[0])
	case http.MethodDelete:
		h.deleteBoard(w, r, parts[0])
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) listBoards(w http.ResponseWriter, r *http.Request) {
	list, err := boards.List(r.Context(), h.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list boards")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"boards": list})
}

func (h *Handler) getBoard(w http.ResponseWriter, r *http.Request, id string) {
	board, err := boards.Load(r.Context(), h.db, id)
	if err != nil {
		writeBoardError(w, err, "read board")
		return
	}
	writeJSON(w, http.StatusOK, board)
}

func (h *Handler) createBoard(w http.ResponseWriter, r *http.Request) {
	var input boardInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := input.validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := boards.Create(r.Context(), h.db, input.Name, input.Cells)
	if err != nil {
		writeBoardError(w, err, "create board")
		return
	}
	board, err := boards.Load(r.Context(), h.db, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read created board")
		return
	}
	writeJSON(w, http.StatusCreated, board)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request, id string) {
	var input boardInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := input.validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := boards.Update(r.Context(), h.db, id, input.Name, input.Cells); err != nil {
		writeBoardError(w, err, "update board")
		return
	}
	board, err := boards.Load(r.Context(), h.db, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read updated board")
		return
	}
	writeJSON(w, http.StatusOK, board)
}

func (h *Handler) deleteBoard(w http.ResponseWriter, r *http.Request, id string) {
	if err := boards.Delete(r.Context(), h.db, id); err != nil {
		writeBoardError(w, err, "delete board")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) setDefaultBoard(w http.ResponseWriter, r *http.Request, id string) {
	if err := boards.SetDefault(r.Context(), h.db, id); err != nil {
		writeBoardError(w, err, "set default board")
		return
	}
	list, err := boards.List(r.Context(), h.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list boards")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"boards": list})
}

// roomBoard returns the layout a new room is played on: boardID, or the
// default board when boardID is empty. With no default board stored, the
// empty layout makes the room use the built-in board.
func (h *Handler) roomBoard(ctx context.Context, boardID string) (room.BoardLayout, error) {
	if boardID == "" {
		id, err := boards.DefaultID(ctx, h.db)
		if errors.Is(err, boards.ErrNotFound) {
			return room.BoardLayout{}, nil
		}
		if err != nil {
			return room.BoardLayout{}, err
		}
		boardID = id
	}
	board, err := boards.Load(ctx, h.db, boardID)
	if err != nil {
		return room.BoardLayout{}, err
	}
	cells, err := room.ValidateBoard(board.Cells)
	if err != nil {
		return room.BoardLayout{}, fmt.Errorf("board %q is invalid: %w", board.Name, err)
	}
	return room.BoardLayout{ID: board.ID, Name: board.Name, Cells: cells}, nil
}

func writeBoardError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, boards.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, boards.ErrNameTaken), errors.Is(err, boards.ErrDeleteDefault):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, fallback)
	}
}
