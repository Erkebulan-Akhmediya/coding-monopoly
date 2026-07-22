package room

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

type testQuestionProvider struct {
	question Question
}

func (p testQuestionProvider) AssignQuestion(string) (Question, error) {
	return p.question, nil
}

type privateTestBroadcaster struct {
	mu         sync.Mutex
	broadcasts []BroadcastEvent
	private    []BroadcastEvent
}

func (b *privateTestBroadcaster) BroadcastRoom(roomID string, msgType string, payload any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.broadcasts = append(b.broadcasts, BroadcastEvent{RoomID: roomID, MsgType: msgType, Payload: payload})
}

func (b *privateTestBroadcaster) BroadcastRoomExcept(roomID, _ string, msgType string, payload any) {
	b.BroadcastRoom(roomID, msgType, payload)
}

func (b *privateTestBroadcaster) SendToPlayer(roomID, clientID, msgType string, payload any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.private = append(b.private, BroadcastEvent{RoomID: roomID, MsgType: msgType, Payload: payload})
}

func (b *privateTestBroadcaster) SendError(string, string) {}

func (b *privateTestBroadcaster) events(msgType string, private bool) []BroadcastEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	events := b.broadcasts
	if private {
		events = b.private
	}
	result := make([]BroadcastEvent, 0)
	for _, event := range events {
		if event.MsgType == msgType {
			result = append(result, event)
		}
	}
	return result
}

func TestRoom_AssignsAndGradesTextQuestionPrivately(t *testing.T) {
	b := &privateTestBroadcaster{}
	r := NewRoomWithQuestionProvider("question-room", b, testQuestionProvider{question: Question{
		ID:              "q-text",
		Type:            "text",
		Prompt:          "What does this print?",
		AcceptedAnswers: []string{"  Hello World  ", "greeting"},
	}})
	r.SetDeadlineDurations(200*time.Millisecond, 200*time.Millisecond, 200*time.Millisecond)
	r.AddOrReconnectPlayer("alice", "Alice")
	r.AddOrReconnectPlayer("bob", "Bob")

	if err := r.ChooseLevel("alice", "easy"); err != nil {
		t.Fatalf("ChooseLevel failed: %v", err)
	}
	public := b.events("question_started", false)
	if len(public) != 1 {
		t.Fatalf("expected one redacted question_started, got %d", len(public))
	}
	publicPayload := public[0].Payload.(QuestionStartedPayload)
	if publicPayload.Prompt != "" || publicPayload.Options != nil || publicPayload.Difficulty != "easy" || publicPayload.Deadline.IsZero() {
		t.Fatalf("question content leaked in public payload: %+v", publicPayload)
	}
	private := b.events("question_started", true)
	if len(private) != 1 || private[0].Payload.(QuestionStartedPayload).Prompt != "What does this print?" {
		t.Fatalf("active player did not receive full question: %+v", private)
	}

	rolls, err := r.SubmitAnswer("alice", json.RawMessage(`{"problem_id":"q-text","answer":"  HELLO WORLD "}`))
	if err != nil || len(rolls) != 1 {
		t.Fatalf("expected correct answer to produce one roll, rolls=%d err=%v", len(rolls), err)
	}
	publicResults := b.events("answer_result", false)
	privateResults := b.events("answer_result", true)
	if len(publicResults) != 1 || len(privateResults) != 1 {
		t.Fatalf("expected one public and one private answer_result, got %d/%d", len(publicResults), len(privateResults))
	}
	if publicResults[0].Payload.(AnswerResultPayload).CorrectAnswer != nil {
		t.Fatal("correct answer leaked in public answer_result")
	}
	if privateResults[0].Payload.(AnswerResultPayload).CorrectAnswer == nil {
		t.Fatal("active player did not receive correct answer review")
	}

	// The stopped timer must not resolve this already-finished turn again.
	time.Sleep(250 * time.Millisecond)
	if got := len(b.events("answer_result", false)); got != 1 {
		t.Fatalf("stale deadline resolved the turn a second time: %d results", got)
	}
}

func TestRoom_TimeoutWinsWithZeroRolls(t *testing.T) {
	b := &privateTestBroadcaster{}
	r := NewRoomWithQuestionProvider("timeout-room", b, testQuestionProvider{question: Question{
		ID:      "q-timeout",
		Type:    "mcq",
		Prompt:  "Pick one",
		Options: []QuestionOption{{ID: "a", Text: "A", Correct: true}, {ID: "b", Text: "B"}},
	}})
	r.SetDeadlineDurations(25*time.Millisecond, 25*time.Millisecond, 25*time.Millisecond)
	r.AddOrReconnectPlayer("alice", "Alice")
	r.AddOrReconnectPlayer("bob", "Bob")
	if err := r.ChooseLevel("alice", "easy"); err != nil {
		t.Fatalf("ChooseLevel failed: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	results := b.events("answer_result", false)
	if len(results) != 1 {
		t.Fatalf("expected one timeout result, got %d", len(results))
	}
	payload := results[0].Payload.(AnswerResultPayload)
	if payload.Correct || !payload.TimedOut || len(payload.Rolls) != 0 {
		t.Fatalf("timeout must be incorrect with zero rolls: %+v", payload)
	}
	if r.GetActivePlayerID() != "bob" {
		t.Fatalf("timeout did not advance the turn, active=%s", r.GetActivePlayerID())
	}
}

func TestRoom_MCQRequiresExactOptionSet(t *testing.T) {
	b := &privateTestBroadcaster{}
	r := NewRoomWithQuestionProvider("mcq-room", b, testQuestionProvider{question: Question{
		ID:      "q-mcq",
		Type:    "mcq",
		Options: []QuestionOption{{ID: "a", Correct: true}, {ID: "b", Correct: true}, {ID: "c"}},
	}})
	r.SetDeadlineDurations(time.Second, time.Second, time.Second)
	r.AddOrReconnectPlayer("alice", "Alice")
	r.AddOrReconnectPlayer("bob", "Bob")
	if err := r.ChooseLevel("alice", "easy"); err != nil {
		t.Fatal(err)
	}
	if rolls, err := r.SubmitAnswer("alice", json.RawMessage(`{"answer":["a"]}`)); err != nil || len(rolls) != 0 {
		t.Fatalf("partial MCQ answer should fail without rolls, rolls=%d err=%v", len(rolls), err)
	}
	result := b.events("answer_result", false)[0].Payload.(AnswerResultPayload)
	if result.Correct {
		t.Fatal("partial MCQ selection was graded correct")
	}
}
