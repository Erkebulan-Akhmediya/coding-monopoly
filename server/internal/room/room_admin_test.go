package room

import (
	"testing"
)

func TestRoom_AdminSkipTurn(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-admin-skip", mock)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("Expected Alice (c1) to be active, got %s", r.GetActivePlayerID())
	}

	skippedName := r.AdminSkipTurn("c1")
	if skippedName != "Alice" {
		t.Errorf("Expected skipped name to be Alice, got %s", skippedName)
	}

	if r.GetActivePlayerID() != "c2" {
		t.Fatalf("Expected Bob (c2) to be active, got %s", r.GetActivePlayerID())
	}

	hasAnswerResult := false
	hasTurnEnded := false
	for _, call := range mock.Broadcasts {
		if call.MsgType == "answer_result" {
			hasAnswerResult = true
			res := call.Payload.(AnswerResultPayload)
			if !res.TimedOut || res.Correct {
				t.Errorf("Expected synthetic timed-out answer result, got %+v", res)
			}
		}
		if call.MsgType == "turn_ended" {
			hasTurnEnded = true
		}
	}
	if !hasAnswerResult {
		t.Errorf("Expected answer_result broadcast during AdminSkipTurn")
	}
	if !hasTurnEnded {
		t.Errorf("Expected turn_ended broadcast during AdminSkipTurn")
	}
}

func TestRoom_AdminKickPlayer(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-admin-kick", mock)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("Expected Alice (c1) to be active, got %s", r.GetActivePlayerID())
	}

	kickedName := r.AdminKickPlayer("c1")
	if kickedName != "Alice" {
		t.Errorf("Expected kicked name to be Alice, got %s", kickedName)
	}

	if r.GetActivePlayerID() != "c2" {
		t.Fatalf("Expected Bob (c2) to be active after Alice is kicked, got %s", r.GetActivePlayerID())
	}

	alice := r.playerMap["c1"]
	if alice.IsConnected {
		t.Errorf("Expected Alice to be marked disconnected")
	}
}

func TestRoom_PausedBlocksChooseAndSubmit(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoomWithQuestionProvider("pause-room", mock, testQuestionProvider{question: Question{
		ID:     "q1",
		Type:   "mcq",
		Prompt: "paused?",
		Options: []QuestionOption{
			{ID: "a", Text: "yes", Correct: true},
		},
	}})
	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	if paused := r.AdminTogglePause(); !paused {
		t.Fatal("expected game to be paused")
	}

	if err := r.ChooseLevel("c1", "easy"); err != ErrGamePaused {
		t.Fatalf("ChooseLevel while paused: got %v, want ErrGamePaused", err)
	}

	if paused := r.AdminTogglePause(); paused {
		t.Fatal("expected game to resume")
	}
	if err := r.ChooseLevel("c1", "easy"); err != nil {
		t.Fatalf("ChooseLevel after resume: %v", err)
	}

	if paused := r.AdminTogglePause(); !paused {
		t.Fatal("expected game to be paused again")
	}
	if _, err := r.SubmitAnswer("c1", []byte(`{"answer":["a"]}`)); err != ErrGamePaused {
		t.Fatalf("SubmitAnswer while paused: got %v, want ErrGamePaused", err)
	}
}
