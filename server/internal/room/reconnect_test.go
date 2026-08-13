package room

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

func TestRoom_ActiveDisconnectGraceThenForfeit(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("grace-forfeit", mock)
	r.SetDisconnectGrace(40 * time.Millisecond)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("expected Alice active")
	}

	r.DisconnectPlayer("c1")
	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("active player should remain Alice during grace, got %s", r.GetActivePlayerID())
	}
	alice := r.playerMap["c1"]
	if alice.IsConnected {
		t.Fatalf("Alice should be marked disconnected during grace")
	}

	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if r.GetActivePlayerID() == "c2" {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if r.GetActivePlayerID() != "c2" {
		t.Fatalf("expected Bob after grace forfeit, got %s", r.GetActivePlayerID())
	}
}

func TestRoom_ActiveDisconnectReconnectWithinGrace(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("grace-reconnect", mock)
	r.SetDisconnectGrace(200 * time.Millisecond)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	alice := r.playerMap["c1"]
	alice.Position = 11
	alice.XP = 120

	r.DisconnectPlayer("c1")
	p, isFirst := r.AddOrReconnectPlayer("c1", "Alice")
	if isFirst || !p.IsConnected {
		t.Fatalf("reconnect failed: first=%v connected=%v", isFirst, p.IsConnected)
	}
	if p.Position != 11 || p.XP != 120 {
		t.Fatalf("expected preserved pos/xp, got pos=%d xp=%d", p.Position, p.XP)
	}
	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("Alice should still be active after reconnect-in-grace, got %s", r.GetActivePlayerID())
	}

	time.Sleep(250 * time.Millisecond)
	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("grace timer must not forfeit after reconnect, got %s", r.GetActivePlayerID())
	}
}

func TestRoom_MidQuestionResumeKeepsDeadline(t *testing.T) {
	mock := &MockBroadcaster{}
	provider := &fixedProvider{q: Question{
		ID:         "q1",
		Type:       "text",
		Difficulty: "easy",
		Prompt:     "2+2?",
		AcceptedAnswers: []string{"4"},
	}}
	r := NewRoomWithQuestionProvider("mid-q", mock, provider)
	r.SetDeadlineDurations(2*time.Second, 2*time.Second, 2*time.Second)
	r.SetDisconnectGrace(300 * time.Millisecond)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	if err := r.ChooseLevel("c1", "easy"); err != nil {
		t.Fatalf("ChooseLevel: %v", err)
	}

	_, beforeActive, beforeDeadline := r.GetTurnState()
	if !beforeActive || beforeDeadline == nil {
		t.Fatalf("expected active question before disconnect")
	}
	originalDeadline := *beforeDeadline

	payloadBefore := r.GetActiveQuestionPayload("c1")
	if payloadBefore == nil || payloadBefore.Prompt == "" {
		t.Fatalf("expected full question payload before disconnect")
	}

	r.DisconnectPlayer("c1")
	// Mid-question must survive grace: deadline unchanged, content available on reclaim.
	_, stillActive, midDeadline := r.GetTurnState()
	if !stillActive || midDeadline == nil || !midDeadline.Equal(originalDeadline) {
		t.Fatalf("deadline must not reset on disconnect: before=%v mid=%v active=%v", originalDeadline, midDeadline, stillActive)
	}
	if r.GetActiveQuestionPayload("c1") == nil {
		t.Fatalf("active question should remain addressable by player id during grace")
	}

	r.AddOrReconnectPlayer("c1", "Alice")
	payloadAfter := r.GetActiveQuestionPayload("c1")
	if payloadAfter == nil {
		t.Fatalf("expected question payload after reconnect")
	}
	if !payloadAfter.Deadline.Equal(originalDeadline) {
		t.Fatalf("reconnect must resume genuine remaining time, got %v want %v", payloadAfter.Deadline, originalDeadline)
	}
	if payloadAfter.Prompt != "2+2?" {
		t.Fatalf("prompt mismatch after resume: %q", payloadAfter.Prompt)
	}
}

func TestRoom_CanReclaimPlayer(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("reclaim", mock)
	r.AddOrReconnectPlayer("c1", "Alice")

	if r.CanReclaimPlayer("c1") {
		t.Fatalf("connected player must not be reclaimable")
	}
	r.DisconnectPlayer("c1")
	if !r.CanReclaimPlayer("c1") {
		t.Fatalf("disconnected player must be reclaimable")
	}
	if r.CanReclaimPlayer("missing") {
		t.Fatalf("unknown id must not be reclaimable")
	}
}

func TestRoom_GameOverOnTargetXP(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("game-over", mock)
	r.SetTargetXP(50)
	r.SetRNG(rand.New(rand.NewSource(1)))

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")
	r.playerMap["c1"].ChosenDifficulty = "easy"
	// High enough that even an unlucky cell effect still clears the threshold.
	r.playerMap["c1"].XP = 200

	if _, err := r.SubmitAnswer("c1", json.RawMessage(`{}`)); err != nil {
		t.Fatalf("SubmitAnswer: %v", err)
	}
	if !r.IsFinished() {
		t.Fatalf("expected game over when XP >= target (xp=%d)", r.playerMap["c1"].XP)
	}
	goPayload := r.GetGameOver()
	if goPayload == nil || goPayload.WinnerID != "c1" || goPayload.Reason != "target_xp" {
		t.Fatalf("unexpected game over payload: %+v", goPayload)
	}
	if r.GetActivePlayerID() != "" {
		t.Fatalf("no active player after game over")
	}
	var sawGameOver bool
	for _, b := range mock.Broadcasts {
		if b.MsgType == "game_over" {
			sawGameOver = true
		}
	}
	if !sawGameOver {
		t.Fatalf("expected game_over broadcast")
	}
}

func TestRoom_AdminEndGame(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("admin-end", mock)
	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")
	r.playerMap["c2"].XP = 80

	payload := r.AdminEndGame("")
	if payload == nil || payload.WinnerID != "c2" || payload.Reason != "admin" {
		t.Fatalf("unexpected admin end payload: %+v", payload)
	}
	if !r.IsFinished() {
		t.Fatalf("room should be finished")
	}
}

type fixedProvider struct {
	q Question
}

func (p *fixedProvider) AssignQuestion(string) (Question, error) {
	return p.q, nil
}
