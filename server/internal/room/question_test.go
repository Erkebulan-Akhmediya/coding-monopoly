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

func (b *privateTestBroadcaster) BroadcastRoomExcept(roomID, excludedClientID string, msgType string, payload any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.broadcasts = append(b.broadcasts, BroadcastEvent{RoomID: roomID, MsgType: msgType, Payload: payload, ExcludedClientID: excludedClientID})
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
	if public[0].ExcludedClientID != "alice" {
		t.Fatalf("question_started was not excluded from active player: excluded=%q", public[0].ExcludedClientID)
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

func TestRoom_SubmitAndTimeoutRaceResolvesExactlyOnce(t *testing.T) {
	for attempt := 0; attempt < 100; attempt++ {
		b := &privateTestBroadcaster{}
		r := NewRoomWithQuestionProvider("race-room", b, testQuestionProvider{question: Question{
			ID:      "q-race",
			Type:    "mcq",
			Prompt:  "race prompt must stay private",
			Options: []QuestionOption{{ID: "correct", Text: "correct option", Correct: true}, {ID: "wrong", Text: "wrong option"}},
		}})
		r.SetDeadlineDurations(time.Hour, time.Hour, time.Hour)
		r.AddOrReconnectPlayer("alice", "Alice")
		r.AddOrReconnectPlayer("bob", "Bob")
		if err := r.ChooseLevel("alice", "easy"); err != nil {
			t.Fatalf("attempt %d: ChooseLevel failed: %v", attempt, err)
		}

		r.mu.RLock()
		turn := r.currentTurn
		r.mu.RUnlock()
		if turn == nil {
			t.Fatalf("attempt %d: expected active turn", attempt)
		}

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			r.resolveTimeout(turn, "alice")
		}()
		go func() {
			defer wg.Done()
			<-start
			_, _ = r.SubmitAnswer("alice", json.RawMessage(`{"problem_id":"q-race","answer":["correct"]}`))
		}()
		close(start)
		wg.Wait()

		publicResults := b.events("answer_result", false)
		privateResults := b.events("answer_result", true)
		if len(publicResults) != 1 || len(privateResults) != 1 {
			t.Fatalf("attempt %d: race produced public/private results %d/%d, want 1/1", attempt, len(publicResults), len(privateResults))
		}
		result := publicResults[0].Payload.(AnswerResultPayload)
		if result.TimedOut == result.Correct {
			t.Fatalf("attempt %d: invalid single winner result: %+v", attempt, result)
		}
		if result.Correct && len(result.Rolls) != 1 {
			t.Fatalf("attempt %d: correct winner applied %d rolls, want 1", attempt, len(result.Rolls))
		}
		if result.TimedOut && len(result.Rolls) != 0 {
			t.Fatalf("attempt %d: timeout winner applied rolls: %+v", attempt, result.Rolls)
		}
		if got := len(b.events("roll_resolved", false)); got != len(result.Rolls) {
			t.Fatalf("attempt %d: roll broadcasts=%d, result rolls=%d", attempt, got, len(result.Rolls))
		}
		if len(b.events("turn_ended", false)) != 1 {
			t.Fatalf("attempt %d: race double-ended the turn", attempt)
		}
		if r.GetActivePlayerID() != "bob" {
			t.Fatalf("attempt %d: race did not advance to Bob, active=%s", attempt, r.GetActivePlayerID())
		}
	}
}

func TestGradeQuestion_TextNormalizationIsBounded(t *testing.T) {
	question := Question{
		Type:            "text",
		AcceptedAnswers: []string{"Hello World"},
	}

	for _, answer := range []string{"  hello world  ", "\n\tHeLLo WoRLD\t"} {
		payload, _ := json.Marshal(answer)
		if !gradeQuestion(question, payload, false) {
			t.Errorf("normalized answer %q was rejected", answer)
		}
	}

	for _, answer := range []string{"hello worlds", "hello worl", "hello world!", "helloworld"} {
		payload, _ := json.Marshal(answer)
		if gradeQuestion(question, payload, false) {
			t.Errorf("wrong-but-similar answer %q was accepted", answer)
		}
	}
}
