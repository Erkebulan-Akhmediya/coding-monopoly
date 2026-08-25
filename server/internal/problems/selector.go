// Package problems contains database reads used by future turn assignment.
package problems

import (
	"context"

	"github.com/jackc/pgx/v5"

	"server/internal/locale"
)

type Option struct {
	ID      string
	Text    locale.Text
	Correct bool
}

// SelectPublishedForAssignment returns a random published problem only. Draft
// content is intentionally unavailable to game assignment code.
func SelectPublishedForAssignment(ctx context.Context, db interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, difficulty string) (id, problemType string, prompt locale.Text, err error) {
	err = db.QueryRow(ctx, `SELECT id, type, prompt_en, prompt_ru, prompt_kz FROM problems
		WHERE difficulty = $1 AND is_published = true
		ORDER BY random() LIMIT 1`, difficulty).Scan(&id, &problemType, &prompt.En, &prompt.Ru, &prompt.Kz)
	return
}

// SelectPublishedQuestion loads the complete server-side question needed for
// grading. Correct values are deliberately kept in this package's return
// value and are only copied into the active room state.
func SelectPublishedQuestion(ctx context.Context, db interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, difficulty string) (id, problemType string, prompt locale.Text, options []Option, acceptedAnswers []locale.Text, err error) {
	err = db.QueryRow(ctx, `SELECT id, type, prompt_en, prompt_ru, prompt_kz FROM problems
		WHERE difficulty = $1 AND is_published = true
		ORDER BY random() LIMIT 1`, difficulty).Scan(&id, &problemType, &prompt.En, &prompt.Ru, &prompt.Kz)
	if err != nil {
		return
	}

	if problemType == "mcq" {
		rows, queryErr := db.Query(ctx, `SELECT id, text_en, text_ru, text_kz, is_correct FROM problem_options WHERE problem_id = $1 ORDER BY id`, id)
		if queryErr != nil {
			err = queryErr
			return
		}
		defer rows.Close()
		for rows.Next() {
			var option Option
			if scanErr := rows.Scan(&option.ID, &option.Text.En, &option.Text.Ru, &option.Text.Kz, &option.Correct); scanErr != nil {
				err = scanErr
				return
			}
			options = append(options, option)
		}
		err = rows.Err()
		return
	}

	rows, queryErr := db.Query(ctx, `SELECT answer_text_en, answer_text_ru, answer_text_kz FROM problem_accepted_answers WHERE problem_id = $1 ORDER BY id`, id)
	if queryErr != nil {
		err = queryErr
		return
	}
	defer rows.Close()
	for rows.Next() {
		var answer locale.Text
		if scanErr := rows.Scan(&answer.En, &answer.Ru, &answer.Kz); scanErr != nil {
			err = scanErr
			return
		}
		acceptedAnswers = append(acceptedAnswers, answer)
	}
	err = rows.Err()
	return
}
