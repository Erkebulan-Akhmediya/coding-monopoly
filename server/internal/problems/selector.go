// Package problems contains database reads used by future turn assignment.
package problems

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// SelectPublishedForAssignment returns a random published problem only. Draft
// content is intentionally unavailable to game assignment code.
func SelectPublishedForAssignment(ctx context.Context, db interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, difficulty string) (id, problemType, prompt string, err error) {
	err = db.QueryRow(ctx, `SELECT id, type, prompt FROM problems
		WHERE difficulty = $1 AND is_published = true
		ORDER BY random() LIMIT 1`, difficulty).Scan(&id, &problemType, &prompt)
	return
}
