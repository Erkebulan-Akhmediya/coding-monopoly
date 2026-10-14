package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"server/internal/locale"
	"server/internal/room"
)

func TestProblemInputValidation(t *testing.T) {
	validMCQ := problemInput{Type: "mcq", Difficulty: "easy", Topic: "Go", Title: locale.FromEn("t"), Prompt: locale.FromEn("p"), Options: []struct {
		Text      locale.Text `json:"text"`
		IsCorrect bool        `json:"is_correct"`
	}{{Text: locale.FromEn("one"), IsCorrect: true}, {Text: locale.FromEn("two")}}}
	if err := validMCQ.validate(); err != nil {
		t.Fatalf("valid mcq rejected: %v", err)
	}
	validMCQ.Options = validMCQ.Options[:1]
	if err := validMCQ.validate(); err == nil {
		t.Fatal("mcq with one option was accepted")
	}
	validText := problemInput{Type: "text", Difficulty: "hard", Topic: "Go", Title: locale.FromEn("t"), Prompt: locale.FromEn("p"), AcceptedAnswers: []locale.Text{locale.FromEn("answer")}}
	if err := validText.validate(); err != nil {
		t.Fatalf("valid text rejected: %v", err)
	}
	missingTopic := problemInput{Type: "text", Difficulty: "hard", Title: locale.FromEn("t"), Prompt: locale.FromEn("p"), AcceptedAnswers: []locale.Text{locale.FromEn("answer")}}
	if err := missingTopic.validate(); err == nil {
		t.Fatal("problem without topic was accepted")
	}
}

func TestTokenAuthorizationRejectsExpiredAndAcceptsValid(t *testing.T) {
	now := time.Now()
	h := &Handler{config: Config{TokenSecret: "secret"}, now: func() time.Time { return now }}
	token, err := h.signToken(now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/admin/problems", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	if !h.authorized(request) {
		t.Fatal("valid token rejected")
	}
	expired, err := h.signToken(now.Add(-time.Second))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+expired)
	if h.authorized(request) {
		t.Fatal("expired token accepted")
	}
}

func TestBoardInputValidation(t *testing.T) {
	cells := room.DefaultBoard()
	cells[0], cells[31] = cells[31], cells[0]
	input := boardInput{Name: "  Hard Mode  ", Cells: cells}
	if err := input.validate(); err != nil {
		t.Fatalf("valid board rejected: %v", err)
	}
	if input.Name != "Hard Mode" || input.Cells[0].Type != "deploy" || input.Cells[31].Index != 31 {
		t.Fatalf("board not normalized: name=%q cell0=%s cell31=%d", input.Name, input.Cells[0].Type, input.Cells[31].Index)
	}

	if err := (&boardInput{Name: " ", Cells: room.DefaultBoard()}).validate(); err == nil {
		t.Fatal("board without a name was accepted")
	}
	if err := (&boardInput{Name: strings.Repeat("b", 129), Cells: room.DefaultBoard()}).validate(); err == nil {
		t.Fatal("board with an overlong name was accepted")
	}
	if err := (&boardInput{Name: "Short", Cells: room.DefaultBoard()[:10]}).validate(); err == nil {
		t.Fatal("board with 10 cells was accepted")
	}
}
