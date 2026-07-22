package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProblemInputValidation(t *testing.T) {
	validMCQ := problemInput{Type: "mcq", Difficulty: "easy", Title: "t", Prompt: "p", Options: []struct {
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}{{Text: "one", IsCorrect: true}, {Text: "two"}}}
	if err := validMCQ.validate(); err != nil {
		t.Fatalf("valid mcq rejected: %v", err)
	}
	validMCQ.Options = validMCQ.Options[:1]
	if err := validMCQ.validate(); err == nil {
		t.Fatal("mcq with one option was accepted")
	}
	validText := problemInput{Type: "text", Difficulty: "hard", Title: "t", Prompt: "p", AcceptedAnswers: []string{"answer"}}
	if err := validText.validate(); err != nil {
		t.Fatalf("valid text rejected: %v", err)
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
