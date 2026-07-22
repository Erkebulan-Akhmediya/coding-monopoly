package ws

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"server/internal/problems"
	"server/internal/room"
)

// DBQuestionProvider assigns only published questions and keeps their answer
// data in the room engine for server-side grading.
type DBQuestionProvider struct {
	db *pgxpool.Pool
}

func NewDBQuestionProvider(db *pgxpool.Pool) room.QuestionProvider {
	return &DBQuestionProvider{db: db}
}

func (p *DBQuestionProvider) AssignQuestion(difficulty string) (room.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, problemType, prompt, options, accepted, err := problems.SelectPublishedQuestion(ctx, p.db, difficulty)
	if err != nil {
		return room.Question{}, err
	}
	question := room.Question{
		ID:              id,
		Type:            problemType,
		Difficulty:      difficulty,
		Prompt:          prompt,
		AcceptedAnswers: accepted,
	}
	for _, option := range options {
		question.Options = append(question.Options, room.QuestionOption{
			ID:      option.ID,
			Text:    option.Text,
			Correct: option.Correct,
		})
	}
	return question, nil
}
