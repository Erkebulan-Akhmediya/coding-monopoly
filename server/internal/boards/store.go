// Package boards persists reusable board layouts that admins pick when
// creating a room.
package boards

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"server/internal/room"
)

var (
	ErrNotFound      = errors.New("board not found")
	ErrNameTaken     = errors.New("a board with this name already exists")
	ErrDeleteDefault = errors.New("the default board cannot be deleted")
)

// DefaultName is the name given to the built-in layout when it is stored.
const DefaultName = "Default"

// DB is the subset of pgx shared by *pgxpool.Pool, *pgx.Conn, and pgx.Tx.
type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Board is a stored board layout. Cells are only loaded for single-board reads.
type Board struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	IsDefault bool             `json:"is_default"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Cells     []room.BoardCell `json:"cells,omitempty"`
}

// List returns all boards without cells, the default board first.
func List(ctx context.Context, db DB) ([]Board, error) {
	rows, err := db.Query(ctx, `SELECT id, name, is_default, created_at, updated_at FROM boards ORDER BY is_default DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	boards := make([]Board, 0)
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.Name, &b.IsDefault, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

// Load returns one board with its cells ordered by index and their params
// normalized to the keys the engine reads.
func Load(ctx context.Context, db DB, id string) (Board, error) {
	var b Board
	if !IsUUID(id) {
		return b, ErrNotFound
	}
	err := db.QueryRow(ctx, `SELECT id, name, is_default, created_at, updated_at FROM boards WHERE id = $1`, id).
		Scan(&b.ID, &b.Name, &b.IsDefault, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return b, ErrNotFound
	}
	if err != nil {
		return b, err
	}
	rows, err := db.Query(ctx, `SELECT cell_index, name_en, name_ru, name_kz, explanation_en, explanation_ru, explanation_kz, type::text, params
		FROM board_cells WHERE board_id = $1 ORDER BY cell_index`, id)
	if err != nil {
		return b, err
	}
	defer rows.Close()
	for rows.Next() {
		var cell room.BoardCell
		var params []byte
		if err := rows.Scan(&cell.Index, &cell.Name.En, &cell.Name.Ru, &cell.Name.Kz, &cell.Explanation.En, &cell.Explanation.Ru, &cell.Explanation.Kz, &cell.Type, &params); err != nil {
			return b, err
		}
		raw := map[string]any{}
		if err := json.Unmarshal(params, &raw); err != nil {
			return b, err
		}
		cell.Params = room.NormalizeCellParams(cell.Type, raw)
		b.Cells = append(b.Cells, cell)
	}
	return b, rows.Err()
}

// DefaultID returns the id of the default board, or ErrNotFound if none is set.
func DefaultID(ctx context.Context, db DB) (string, error) {
	var id string
	err := db.QueryRow(ctx, `SELECT id FROM boards WHERE is_default`).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return id, err
}

// Create stores a new non-default board. cells must already be validated.
func Create(ctx context.Context, db DB, name string, cells []room.BoardCell) (string, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var id string
	if err := tx.QueryRow(ctx, `INSERT INTO boards (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return "", mapNameConflict(err)
	}
	if err := insertCells(ctx, tx, id, cells); err != nil {
		return "", err
	}
	return id, tx.Commit(ctx)
}

// Update renames a board and replaces all of its cells. cells must already be validated.
func Update(ctx context.Context, db DB, id, name string, cells []room.BoardCell) error {
	if !IsUUID(id) {
		return ErrNotFound
	}
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `UPDATE boards SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, name, id)
	if err != nil {
		return mapNameConflict(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, `DELETE FROM board_cells WHERE board_id = $1`, id); err != nil {
		return err
	}
	if err := insertCells(ctx, tx, id, cells); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Delete removes a non-default board. Rooms already created with it keep
// their own copy of the cells.
func Delete(ctx context.Context, db DB, id string) error {
	if !IsUUID(id) {
		return ErrNotFound
	}
	var isDefault bool
	err := db.QueryRow(ctx, `SELECT is_default FROM boards WHERE id = $1`, id).Scan(&isDefault)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if isDefault {
		return ErrDeleteDefault
	}
	result, err := db.Exec(ctx, `DELETE FROM boards WHERE id = $1 AND NOT is_default`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SetDefault makes id the only default board.
func SetDefault(ctx context.Context, db DB, id string) error {
	if !IsUUID(id) {
		return ErrNotFound
	}
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE boards SET is_default = false WHERE is_default AND id <> $1`, id); err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `UPDATE boards SET is_default = true WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return tx.Commit(ctx)
}

// EnsureDefault makes sure a default board exists: it stores the built-in
// layout as the default when there are no boards, or promotes the oldest
// board when none is marked default.
func EnsureDefault(ctx context.Context, db DB) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var hasDefault, hasAny bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM boards WHERE is_default), EXISTS(SELECT 1 FROM boards)`).Scan(&hasDefault, &hasAny); err != nil {
		return err
	}
	switch {
	case hasDefault:
		return nil
	case hasAny:
		if _, err := tx.Exec(ctx, `UPDATE boards SET is_default = true WHERE id = (SELECT id FROM boards ORDER BY created_at, id LIMIT 1)`); err != nil {
			return err
		}
	default:
		var id string
		if err := tx.QueryRow(ctx, `INSERT INTO boards (name, is_default) VALUES ($1, true) RETURNING id`, DefaultName).Scan(&id); err != nil {
			return err
		}
		if err := insertCells(ctx, tx, id, room.DefaultBoard()); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func insertCells(ctx context.Context, tx pgx.Tx, boardID string, cells []room.BoardCell) error {
	for _, cell := range cells {
		params, err := json.Marshal(cell.Params)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO board_cells (board_id, cell_index, name_en, name_ru, name_kz, explanation_en, explanation_ru, explanation_kz, type, params)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			boardID, cell.Index, cell.Name.En, cell.Name.Ru, cell.Name.Kz, cell.Explanation.En, cell.Explanation.Ru, cell.Explanation.Kz, cell.Type, params); err != nil {
			return err
		}
	}
	return nil
}

func mapNameConflict(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrNameTaken
	}
	return err
}

// IsUUID reports whether s is a canonical hyphenated UUID, so malformed ids
// are reported as not found instead of as database errors.
func IsUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
				return false
			}
		}
	}
	return true
}
