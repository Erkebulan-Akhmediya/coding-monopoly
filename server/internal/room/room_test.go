package room

import (
	"encoding/json"
	"math/rand"
	"testing"
)

// MockBroadcaster captures room broadcasts and error messages for test assertions.
type MockBroadcaster struct {
	Broadcasts []BroadcastEvent
	Errors     []ErrorEvent
}

type BroadcastEvent struct {
	RoomID  string
	MsgType string
	Payload interface{}
}

type ErrorEvent struct {
	ClientID string
	ErrMsg   string
}

func (m *MockBroadcaster) BroadcastRoom(roomID string, msgType string, payload interface{}) {
	m.Broadcasts = append(m.Broadcasts, BroadcastEvent{
		RoomID:  roomID,
		MsgType: msgType,
		Payload: payload,
	})
}

func (m *MockBroadcaster) SendError(clientID string, errMsg string) {
	m.Errors = append(m.Errors, ErrorEvent{
		ClientID: clientID,
		ErrMsg:   errMsg,
	})
}

func TestRoom_JoinOrderAndActivePlayer(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room", mock)

	// Add Alice
	pA, isFirstA := r.AddOrReconnectPlayer("c1", "Alice")
	if !isFirstA || pA.Name != "Alice" {
		t.Fatalf("Expected Alice to be first connected player")
	}
	if r.GetActivePlayerID() != "c1" {
		t.Errorf("Expected active player ID to be c1, got %s", r.GetActivePlayerID())
	}

	// Add Bob
	pB, isFirstB := r.AddOrReconnectPlayer("c2", "Bob")
	if isFirstB || pB.Name != "Bob" {
		t.Fatalf("Bob should not be first player")
	}

	// Add Charlie
	r.AddOrReconnectPlayer("c3", "Charlie")

	// Verify join order
	players := r.GetPlayers()
	if len(players) != 3 {
		t.Fatalf("Expected 3 players, got %d", len(players))
	}
	if players[0].ID != "c1" || players[1].ID != "c2" || players[2].ID != "c3" {
		t.Errorf("Players not in join order: %+v", players)
	}

	// Verify initial turn_started broadcast for Alice
	if len(mock.Broadcasts) < 1 {
		t.Fatalf("Expected turn_started broadcast on first join")
	}
	firstEvent := mock.Broadcasts[0]
	if firstEvent.MsgType != "turn_started" {
		t.Errorf("Expected turn_started message, got %s", firstEvent.MsgType)
	}
	payload, ok := firstEvent.Payload.(TurnStartedPayload)
	if !ok || payload.ActivePlayerID != "c1" {
		t.Errorf("Unexpected turn_started payload: %+v", firstEvent.Payload)
	}
}

func TestRoom_ActivePlayerEnforcement(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-enforce", mock)

	r.AddOrReconnectPlayer("c1", "Alice") // Active player
	r.AddOrReconnectPlayer("c2", "Bob")   // Inactive player

	// Bob attempts choose_level
	err := r.ChooseLevel("c2", "easy")
	if err != ErrNotActivePlayer {
		t.Errorf("Expected ErrNotActivePlayer when non-active player chooses level, got: %v", err)
	}

	// Bob attempts submit_answer
	_, err = r.SubmitAnswer("c2", json.RawMessage("{}"))
	if err != ErrNotActivePlayer {
		t.Errorf("Expected ErrNotActivePlayer when non-active player submits answer, got: %v", err)
	}

	// Verify error messages were sent to Bob
	if len(mock.Errors) != 2 {
		t.Fatalf("Expected 2 error events for Bob, got %d", len(mock.Errors))
	}
	if mock.Errors[0].ClientID != "c2" || mock.Errors[1].ClientID != "c2" {
		t.Errorf("Errors not sent to Bob: %+v", mock.Errors)
	}

	// Alice (active player) attempts choose_level
	err = r.ChooseLevel("c1", "medium")
	if err != nil {
		t.Errorf("Active player choose_level failed: %v", err)
	}
}

func TestRoom_DiceRollingAndDifficultyResolution(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-dice", mock)
	// Fixed deterministic RNG for testing
	r.SetRNG(rand.New(rand.NewSource(42)))

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	// Set difficulty: medium (2 rolls)
	if err := r.ChooseLevel("c1", "medium"); err != nil {
		t.Fatalf("ChooseLevel failed: %v", err)
	}

	// Submit answer
	results, err := r.SubmitAnswer("c1", json.RawMessage("{}"))
	if err != nil {
		t.Fatalf("SubmitAnswer failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 roll results for medium difficulty, got %d", len(results))
	}

	for _, res := range results {
		if res.DieRoll < 1 || res.DieRoll > 6 {
			t.Errorf("Die roll out of range 1..6: %d", res.DieRoll)
		}
		if res.PlayerID != "c1" {
			t.Errorf("Roll result player ID mismatch: %s", res.PlayerID)
		}
	}

	// Check turn advancement to Bob ("c2")
	if r.GetActivePlayerID() != "c2" {
		t.Errorf("Expected turn to advance to Bob (c2), got %s", r.GetActivePlayerID())
	}
}

func TestRoom_AllCatalogEffects(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-effects", mock)
	p := NewPlayer("c1", "Tester")

	// Test 1: xp_gain
	cellGain := BoardCell{Index: 1, Name: "Gain", Type: "xp_gain", Params: map[string]interface{}{"amount": 20}}
	resGain := r.ApplyCellEffect(p, cellGain)
	if resGain.XPDelta != 20 || p.XP != 20 {
		t.Errorf("xp_gain failed: expected +20 XP, got delta %d, total %d", resGain.XPDelta, p.XP)
	}

	// Test 2: double_xp + xp_gain
	cellDouble := BoardCell{Index: 5, Name: "Double", Type: "double_xp", Params: map[string]interface{}{}}
	resDouble := r.ApplyCellEffect(p, cellDouble)
	if !p.DoubleXP || resDouble.EffectType != "double_xp" {
		t.Errorf("double_xp failed: DoubleXP flag not set")
	}
	resGainDoubled := r.ApplyCellEffect(p, cellGain)
	if resGainDoubled.XPDelta != 40 || p.XP != 60 || p.DoubleXP {
		t.Errorf("double_xp effect failed: expected +40 XP, got delta %d, total %d", resGainDoubled.XPDelta, p.XP)
	}

	// Test 3: xp_loss
	cellLoss := BoardCell{Index: 2, Name: "Loss", Type: "xp_loss", Params: map[string]interface{}{"amount": 15}}
	resLoss := r.ApplyCellEffect(p, cellLoss)
	if resLoss.XPDelta != -15 || p.XP != 45 {
		t.Errorf("xp_loss failed: expected -15 XP, got delta %d, total %d", resLoss.XPDelta, p.XP)
	}

	// Test 4: teleport
	cellTele := BoardCell{Index: 12, Name: "Teleport", Type: "teleport", Params: map[string]interface{}{"target_position": 16}}
	resTele := r.ApplyCellEffect(p, cellTele)
	if p.Position != 16 || resTele.NewPosition != 16 {
		t.Errorf("teleport failed: expected pos 16, got %d", p.Position)
	}

	// Test 5: skip_next
	cellSkip := BoardCell{Index: 10, Name: "Skip", Type: "skip_next", Params: map[string]interface{}{}}
	r.ApplyCellEffect(p, cellSkip)
	if !p.SkipNextTurn {
		t.Errorf("skip_next failed: SkipNextTurn flag not set")
	}

	// Test 6: free_pass
	cellPass := BoardCell{Index: 7, Name: "Pass", Type: "free_pass", Params: map[string]interface{}{}}
	r.ApplyCellEffect(p, cellPass)
	if p.FreePasses != 1 {
		t.Errorf("free_pass failed: FreePasses count %d != 1", p.FreePasses)
	}

	// Test 7: special_challenge
	cellChallenge := BoardCell{Index: 13, Name: "Challenge", Type: "special_challenge", Params: map[string]interface{}{"bonus": 30}}
	resChall := r.ApplyCellEffect(p, cellChallenge)
	if resChall.XPDelta != 30 || p.XP != 75 {
		t.Errorf("special_challenge failed: expected +30 XP, got delta %d, total %d", resChall.XPDelta, p.XP)
	}

	// Test 8: mystery
	cellMystery := BoardCell{Index: 3, Name: "Mystery", Type: "mystery", Params: map[string]interface{}{}}
	resMyst := r.ApplyCellEffect(p, cellMystery)
	if resMyst.EffectType != "mystery" {
		t.Errorf("mystery failed: unexpected effect type %s", resMyst.EffectType)
	}
}

func TestRoom_DisconnectAndReconnectTurnSkip(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-disc", mock)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")
	r.AddOrReconnectPlayer("c3", "Charlie")

	if r.GetActivePlayerID() != "c1" {
		t.Fatalf("Expected Alice to start active")
	}

	// Disconnect Bob ("c2")
	r.DisconnectPlayer("c2")

	// Verify Bob's slot is retained in players list
	players := r.GetPlayers()
	if len(players) != 3 {
		t.Fatalf("Players count should remain 3 despite Bob disconnecting")
	}
	if players[1].IsConnected {
		t.Errorf("Bob should be marked disconnected")
	}

	// Alice completes turn
	_, err := r.SubmitAnswer("c1", json.RawMessage("{}"))
	if err != nil {
		t.Fatalf("Alice turn failed: %v", err)
	}

	// Active player should skip disconnected Bob ("c2") and go to Charlie ("c3")!
	if r.GetActivePlayerID() != "c3" {
		t.Errorf("Expected turn to skip disconnected Bob and advance to Charlie (c3), got %s", r.GetActivePlayerID())
	}

	// Bob reconnects!
	pBob, isFirst := r.AddOrReconnectPlayer("c2", "Bob")
	if !pBob.IsConnected || isFirst {
		t.Errorf("Bob reconnect failed: isConnected=%t, isFirst=%t", pBob.IsConnected, isFirst)
	}

	// Charlie completes turn
	_, err = r.SubmitAnswer("c3", json.RawMessage("{}"))
	if err != nil {
		t.Fatalf("Charlie turn failed: %v", err)
	}

	// Turn advances to Alice ("c1")
	if r.GetActivePlayerID() != "c1" {
		t.Errorf("Expected turn to advance to Alice (c1), got %s", r.GetActivePlayerID())
	}

	// Alice completes turn
	_, err = r.SubmitAnswer("c1", json.RawMessage("{}"))
	if err != nil {
		t.Fatalf("Alice second turn failed: %v", err)
	}

	// Reconnected Bob ("c2") should now get their turn!
	if r.GetActivePlayerID() != "c2" {
		t.Errorf("Expected reconnected Bob (c2) to receive turn, got %s", r.GetActivePlayerID())
	}
}

func TestRoom_SkipNextTurnModifier(t *testing.T) {
	mock := &MockBroadcaster{}
	r := NewRoom("test-room-skip", mock)

	r.AddOrReconnectPlayer("c1", "Alice")
	r.AddOrReconnectPlayer("c2", "Bob")

	// Set Bob's SkipNextTurn flag
	bob := r.playerMap["c2"]
	bob.SkipNextTurn = true

	// Alice completes turn
	_, err := r.SubmitAnswer("c1", json.RawMessage("{}"))
	if err != nil {
		t.Fatalf("Alice turn failed: %v", err)
	}

	// Turn should skip Bob (due to SkipNextTurn) and return to Alice!
	if r.GetActivePlayerID() != "c1" {
		t.Errorf("Expected turn to skip Bob (SkipNextTurn) and return to Alice, got %s", r.GetActivePlayerID())
	}

	// Verify Bob's SkipNextTurn flag was cleared
	if bob.SkipNextTurn {
		t.Errorf("Bob's SkipNextTurn flag should have been cleared after skipping")
	}
}
