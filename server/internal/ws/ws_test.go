package ws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"server/internal/room"
)

func setupTestServer(t *testing.T, clientOpts ...ClientOptions) (*Hub, *httptest.Server) {
	return setupTestServerWithProvider(t, nil, clientOpts...)
}

func setupTestServerWithProvider(t *testing.T, provider room.QuestionProvider, clientOpts ...ClientOptions) (*Hub, *httptest.Server) {
	t.Helper()
	hub := NewHub(provider)
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r, clientOpts...)
	}))

	t.Cleanup(func() {
		server.Close()
		hub.Stop()
	})

	return hub, server
}

type fixedQuestionProvider struct {
	question room.Question
}

func (p fixedQuestionProvider) AssignQuestion(string) (room.Question, error) {
	return p.question, nil
}

func connectClient(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	return conn
}

func sendJoin(t *testing.T, conn *websocket.Conn, name string, roomID string) {
	t.Helper()
	joinPayload := JoinPayload{
		Name:   name,
		RoomID: roomID,
	}
	payloadBytes, _ := json.Marshal(joinPayload)
	msg := Message{
		Type:    MessageTypeJoin,
		RoomID:  roomID,
		Payload: payloadBytes,
	}
	if err := conn.WriteJSON(msg); err != nil {
		t.Fatalf("Failed to send join message: %v", err)
	}
}

func readMessageTimeout(t *testing.T, conn *websocket.Conn, timeout time.Duration) Message {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	var msg Message
	err := conn.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("Failed to read JSON message within deadline: %v", err)
	}
	return msg
}

func readRawMessageTimeout(t *testing.T, conn *websocket.Conn, timeout time.Duration) ([]byte, Message) {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read raw WebSocket message within deadline: %v", err)
	}
	var msg Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		t.Fatalf("Failed to decode raw WebSocket message: %v", err)
	}
	return raw, msg
}

func readUntilType(t *testing.T, conn *websocket.Conn, want string) ([]byte, Message) {
	t.Helper()
	for i := 0; i < 12; i++ {
		raw, msg := readRawMessageTimeout(t, conn, 2*time.Second)
		if msg.Type == want {
			return raw, msg
		}
	}
	t.Fatalf("did not receive message type %q", want)
	return nil, Message{}
}

func readTypes(t *testing.T, conn *websocket.Conn, wanted ...string) map[string]Message {
	t.Helper()
	remaining := make(map[string]bool, len(wanted))
	for _, want := range wanted {
		remaining[want] = true
	}
	result := make(map[string]Message, len(wanted))
	for len(remaining) > 0 {
		_, msg := readRawMessageTimeout(t, conn, 2*time.Second)
		if remaining[msg.Type] {
			result[msg.Type] = msg
			delete(remaining, msg.Type)
		}
	}
	return result
}

func TestWS_JoinFlowAndStateSync(t *testing.T) {
	_, server := setupTestServer(t)

	connA := connectClient(t, server)
	defer connA.Close()

	sendJoin(t, connA, "Alice", "room-1")

	// Expect presence broadcast for Alice
	_, msg1 := readUntilType(t, connA, MessageTypePresence)
	if msg1.Type != MessageTypePresence {
		t.Fatalf("Expected presence message, got %s", msg1.Type)
	}
	var presenceA PresencePayload
	if err := json.Unmarshal(msg1.Payload, &presenceA); err != nil {
		t.Fatalf("Failed to parse presence payload: %v", err)
	}
	if presenceA.Event != "joined" || presenceA.Player.Name != "Alice" {
		t.Errorf("Unexpected presence payload: %+v", presenceA)
	}

	// Expect state_sync message with 1 player
	_, msg2 := readUntilType(t, connA, MessageTypeStateSync)
	if msg2.Type != MessageTypeStateSync {
		t.Fatalf("Expected state_sync message, got %s", msg2.Type)
	}
	var stateSyncA StateSyncPayload
	if err := json.Unmarshal(msg2.Payload, &stateSyncA); err != nil {
		t.Fatalf("Failed to parse state_sync payload: %v", err)
	}
	if len(stateSyncA.Players) != 1 || stateSyncA.Players[0].Name != "Alice" {
		t.Errorf("Unexpected state_sync payload: %+v", stateSyncA)
	}
}

func TestWS_MultiClientPresenceAndStateSync(t *testing.T) {
	_, server := setupTestServer(t)

	// Connect Client A
	connA := connectClient(t, server)
	defer connA.Close()
	sendJoin(t, connA, "Alice", "game-room")
	_, _ = readUntilType(t, connA, MessageTypePresence)  // presence A
	_, _ = readUntilType(t, connA, MessageTypeStateSync) // state_sync [Alice]

	// Connect Client B
	connB := connectClient(t, server)
	defer connB.Close()
	sendJoin(t, connB, "Bob", "game-room")

	// Client A should receive presence for Bob joining
	_, msgA_presence := readUntilType(t, connA, MessageTypePresence)
	if msgA_presence.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence for Bob, got %s", msgA_presence.Type)
	}
	var presenceB PresencePayload
	_ = json.Unmarshal(msgA_presence.Payload, &presenceB)
	if presenceB.Event != "joined" || presenceB.Player.Name != "Bob" {
		t.Errorf("Client A received wrong presence: %+v", presenceB)
	}

	// Client A should receive updated state_sync containing both Alice and Bob
	_, msgA_sync := readUntilType(t, connA, MessageTypeStateSync)
	if msgA_sync.Type != MessageTypeStateSync {
		t.Fatalf("Client A expected state_sync, got %s", msgA_sync.Type)
	}
	var stateSync1 StateSyncPayload
	_ = json.Unmarshal(msgA_sync.Payload, &stateSync1)
	if len(stateSync1.Players) != 2 {
		t.Errorf("Expected 2 players in state_sync for Client A, got %d", len(stateSync1.Players))
	}

	// Client B should also receive presence and state_sync
	_, msgB_presence := readUntilType(t, connB, MessageTypePresence)
	if msgB_presence.Type != MessageTypePresence {
		t.Fatalf("Client B expected presence message, got %s", msgB_presence.Type)
	}
	_, msgB_sync := readUntilType(t, connB, MessageTypeStateSync)
	if msgB_sync.Type != MessageTypeStateSync {
		t.Fatalf("Client B expected state_sync message, got %s", msgB_sync.Type)
	}
	var stateSyncB StateSyncPayload
	_ = json.Unmarshal(msgB_sync.Payload, &stateSyncB)
	if len(stateSyncB.Players) != 2 {
		t.Errorf("Expected 2 players in state_sync for Client B, got %d", len(stateSyncB.Players))
	}
}

func TestWS_GracefulDisconnect(t *testing.T) {
	_, server := setupTestServer(t)

	connA := connectClient(t, server)
	defer connA.Close()
	sendJoin(t, connA, "Alice", "room-disconnect")
	_, _ = readUntilType(t, connA, MessageTypePresence)
	_, _ = readUntilType(t, connA, MessageTypeStateSync)

	connB := connectClient(t, server)
	sendJoin(t, connB, "Bob", "room-disconnect")
	_, _ = readUntilType(t, connA, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connA, MessageTypeStateSync) // state_sync [A, B]
	_, _ = readUntilType(t, connB, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connB, MessageTypeStateSync) // state_sync [A, B]

	// Graceful disconnect Client B
	_ = connB.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "leaving"))
	_ = connB.Close()

	// Client A should receive presence ("left") for Bob
	_, msgA_left := readUntilType(t, connA, MessageTypePresence)
	if msgA_left.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence message on disconnect, got %s", msgA_left.Type)
	}
	var presenceLeft PresencePayload
	_ = json.Unmarshal(msgA_left.Payload, &presenceLeft)
	if presenceLeft.Event != "left" || presenceLeft.Player.Name != "Bob" {
		t.Errorf("Unexpected leave presence: %+v", presenceLeft)
	}

	// Client A should receive updated state_sync with only Alice
	_, msgA_sync := readUntilType(t, connA, MessageTypeStateSync)
	if msgA_sync.Type != MessageTypeStateSync {
		t.Fatalf("Client A expected state_sync after disconnect, got %s", msgA_sync.Type)
	}
	var stateSyncAfter StateSyncPayload
	_ = json.Unmarshal(msgA_sync.Payload, &stateSyncAfter)
	if len(stateSyncAfter.Players) != 1 || stateSyncAfter.Players[0].Name != "Alice" {
		t.Errorf("Expected 1 player (Alice) after disconnect, got: %+v", stateSyncAfter.Players)
	}
}

func TestWS_UngracefulDisconnect(t *testing.T) {
	// Set fast pong timeout for test speed
	opts := ClientOptions{
		PongWait:   300 * time.Millisecond,
		PingPeriod: 100 * time.Millisecond,
	}
	_, server := setupTestServer(t, opts)

	connA := connectClient(t, server)
	defer connA.Close()
	sendJoin(t, connA, "Alice", "room-ungraceful")
	_, _ = readUntilType(t, connA, MessageTypePresence)
	_, _ = readUntilType(t, connA, MessageTypeStateSync)

	connB := connectClient(t, server)
	sendJoin(t, connB, "Bob", "room-ungraceful")
	_, _ = readUntilType(t, connA, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connA, MessageTypeStateSync) // state_sync [A, B]
	_, _ = readUntilType(t, connB, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connB, MessageTypeStateSync) // state_sync [A, B]

	// Simulate ungraceful disconnect by abruptly closing underlying TCP connection without sending close frame
	_ = connB.UnderlyingConn().Close()

	// Client A should receive presence ("left") and state_sync automatically after pong timeout
	_, msgA_left := readUntilType(t, connA, MessageTypePresence)
	if msgA_left.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence message on ungraceful disconnect, got %s", msgA_left.Type)
	}
	var presenceLeft PresencePayload
	_ = json.Unmarshal(msgA_left.Payload, &presenceLeft)
	if presenceLeft.Event != "left" || presenceLeft.Player.Name != "Bob" {
		t.Errorf("Unexpected leave presence: %+v", presenceLeft)
	}

	_, msgA_sync := readUntilType(t, connA, MessageTypeStateSync)
	if msgA_sync.Type != MessageTypeStateSync {
		t.Fatalf("Client A expected state_sync after ungraceful disconnect, got %s", msgA_sync.Type)
	}
	var stateSyncAfter StateSyncPayload
	_ = json.Unmarshal(msgA_sync.Payload, &stateSyncAfter)
	if len(stateSyncAfter.Players) != 1 || stateSyncAfter.Players[0].Name != "Alice" {
		t.Errorf("Expected 1 player after ungraceful disconnect, got: %+v", stateSyncAfter.Players)
	}
}

func TestWS_EmptyNameValidation(t *testing.T) {
	_, server := setupTestServer(t)

	conn := connectClient(t, server)
	defer conn.Close()

	sendJoin(t, conn, "   ", "room-val")

	msg := readMessageTimeout(t, conn, 2*time.Second)
	if msg.Type != MessageTypeError {
		t.Fatalf("Expected error message for empty name, got %s", msg.Type)
	}
	if !strings.Contains(msg.Error, "Name is required") {
		t.Errorf("Expected name required error message, got: %s", msg.Error)
	}
}

func TestWS_MessageEnvelopeVersionAndType(t *testing.T) {
	hub, server := setupTestServer(t)

	conn := connectClient(t, server)
	defer conn.Close()

	sendJoin(t, conn, "Alice", "room-envelope")

	// Read presence message
	rawBytes, _ := readUntilType(t, conn, MessageTypePresence)
	var raw map[string]interface{}
	if err := json.Unmarshal(rawBytes, &raw); err != nil {
		t.Fatalf("Failed to decode raw JSON message: %v", err)
	}

	if raw["v"] != float64(1) {
		t.Errorf("Expected version v=1 in envelope, got: %v", raw["v"])
	}
	if raw["type"] != MessageTypePresence {
		t.Errorf("Expected type=%s in envelope, got: %v", MessageTypePresence, raw["type"])
	}

	// Send unknown/future message type
	unknownMsg := Message{
		Version: 99,
		Type:    "future_msg_v99",
	}
	if err := conn.WriteJSON(unknownMsg); err != nil {
		t.Fatalf("Failed to send unknown message type: %v", err)
	}

	// Verify connection is still alive by asking for room players from hub
	time.Sleep(100 * time.Millisecond)
	players := hub.GetRoomPlayers("room-envelope")
	if len(players) != 1 || players[0].Name != "Alice" {
		t.Errorf("Client should remain active after sending unknown message type, got: %+v", players)
	}
}

func TestWS_ConcurrentConnectDisconnectRace(t *testing.T) {
	hub, server := setupTestServer(t)

	const numClients = 15
	done := make(chan struct{})

	for i := 0; i < numClients; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()

			roomID := "race-room-1"
			if id%2 == 0 {
				roomID = "race-room-2"
			}

			conn := connectClient(t, server)
			defer conn.Close()

			sendJoin(t, conn, strings.Repeat("User", 1)+string(rune('A'+id)), roomID)

			// Perform concurrent reads & state checks
			for j := 0; j < 5; j++ {
				_ = hub.GetRoomPlayers(roomID)
				bMsg, _ := NewMessage("test_event", roomID, map[string]string{"foo": "bar"})
				hub.BroadcastRoom(roomID, bMsg)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	for i := 0; i < numClients; i++ {
		<-done
	}
}

func TestWS_GoroutineLeakOnDisconnectAndShutdown(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	}))

	const clientCount = 5
	var conns []*websocket.Conn
	for i := 0; i < clientCount; i++ {
		conn := connectClient(t, server)
		sendJoin(t, conn, strings.Repeat("P", 1)+string(rune('1'+i)), "leak-room")
		conns = append(conns, conn)
	}

	time.Sleep(100 * time.Millisecond)

	// Close all connections
	for _, c := range conns {
		_ = c.Close()
	}

	time.Sleep(200 * time.Millisecond)

	// Stop Hub
	server.Close()
	hub.Stop()

	time.Sleep(100 * time.Millisecond)

	// Verify room is empty
	players := hub.GetRoomPlayers("leak-room")
	if len(players) != 0 {
		t.Errorf("Expected 0 players after disconnect and hub shutdown, got %d", len(players))
	}
}

func TestWS_SlowClientBackpressure(t *testing.T) {
	hub, server := setupTestServer(t)

	// Fast client
	connFast := connectClient(t, server)
	defer connFast.Close()
	sendJoin(t, connFast, "FastUser", "backpressure-room")
	_ = readMessageTimeout(t, connFast, 2*time.Second) // presence
	_ = readMessageTimeout(t, connFast, 2*time.Second) // state_sync

	// Slow client: custom tiny send channel buffer override isn't needed if we fill the default channel
	// Or we can broadcast many messages rapidly
	connSlow := connectClient(t, server)
	defer connSlow.Close()
	sendJoin(t, connSlow, "SlowUser", "backpressure-room")
	_ = readMessageTimeout(t, connFast, 2*time.Second) // presence SlowUser
	_ = readMessageTimeout(t, connFast, 2*time.Second) // state_sync [Fast, Slow]

	// Broadcast many messages without reading from connSlow
	bigPayload := strings.Repeat("x", 1024)
	for i := 0; i < 300; i++ {
		msg, _ := NewMessage("test_broadcast", "backpressure-room", map[string]string{"data": bigPayload})
		hub.BroadcastRoom("backpressure-room", msg)
	}

	// Fast client should receive messages without blocking
	recvCount := 0
	for i := 0; i < 50; i++ {
		_ = connFast.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		var m Message
		if err := connFast.ReadJSON(&m); err == nil {
			recvCount++
		} else {
			break
		}
	}

	if recvCount == 0 {
		t.Errorf("Fast client failed to receive broadcast messages during slow client backpressure")
	}
}

func TestWS_TurnEngineIntegration(t *testing.T) {
	_, server := setupTestServer(t)

	// Connect Alice
	connA := connectClient(t, server)
	defer connA.Close()
	sendJoin(t, connA, "Alice", "ws-turn-room")

	initialMessages := readTypes(t, connA, MessageTypePresence, MessageTypeStateSync, MessageTypeTurnStarted)
	msgStartA := initialMessages[MessageTypeTurnStarted]
	if msgStartA.Type != MessageTypeTurnStarted {
		t.Fatalf("Expected turn_started message, got %s", msgStartA.Type)
	}

	var turnStartedA struct {
		ActivePlayerID string `json:"active_player_id"`
	}
	_ = json.Unmarshal(msgStartA.Payload, &turnStartedA)
	if turnStartedA.ActivePlayerID == "" {
		t.Errorf("Empty active_player_id in turn_started: %+v", turnStartedA)
	}

	// Connect Bob
	connB := connectClient(t, server)
	defer connB.Close()
	sendJoin(t, connB, "Bob", "ws-turn-room")

	_, _ = readUntilType(t, connA, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connA, MessageTypeStateSync) // state_sync [A, B]

	_, _ = readUntilType(t, connB, MessageTypePresence)  // presence B
	_, _ = readUntilType(t, connB, MessageTypeStateSync) // state_sync [A, B]

	// Bob attempts choose_level (NOT active player) -> should be rejected with error
	chooseMsg := Message{
		Type:    MessageTypeChooseLevel,
		RoomID:  "ws-turn-room",
		Payload: json.RawMessage(`{"difficulty":"easy"}`),
	}
	if err := connB.WriteJSON(chooseMsg); err != nil {
		t.Fatalf("Failed to send choose_level from Bob: %v", err)
	}

	errMsg := readMessageTimeout(t, connB, 2*time.Second)
	if errMsg.Type != MessageTypeError {
		t.Fatalf("Expected error message for Bob choose_level, got %s", errMsg.Type)
	}
	if !strings.Contains(errMsg.Error, "Not your turn") {
		t.Errorf("Expected 'Not your turn' error, got: %s", errMsg.Error)
	}

	// Alice (active player) selects difficulty hard (3 rolls)
	chooseAlice := Message{
		Type:    MessageTypeChooseLevel,
		RoomID:  "ws-turn-room",
		Payload: json.RawMessage(`{"difficulty":"hard"}`),
	}
	if err := connA.WriteJSON(chooseAlice); err != nil {
		t.Fatalf("Failed to send choose_level from Alice: %v", err)
	}

	// Alice submits answer
	submitAlice := Message{
		Type:    MessageTypeSubmitAnswer,
		RoomID:  "ws-turn-room",
		Payload: json.RawMessage(`{}`),
	}
	if err := connA.WriteJSON(submitAlice); err != nil {
		t.Fatalf("Failed to send submit_answer from Alice: %v", err)
	}

	// Read 3 roll_resolved broadcasts on Alice connection
	for i := 1; i <= 3; i++ {
		rollMsg := readMessageTimeout(t, connA, 2*time.Second)
		if rollMsg.Type != MessageTypeRollResolved {
			t.Fatalf("Expected roll_resolved message %d, got %s", i, rollMsg.Type)
		}
	}

	// Read turn_ended broadcast on Alice connection
	endMsg := readMessageTimeout(t, connA, 2*time.Second)
	if endMsg.Type != MessageTypeTurnEnded {
		t.Fatalf("Expected turn_ended message, got %s", endMsg.Type)
	}

	// Read turn_started broadcast on Alice connection (now Bob's turn)
	nextStart := readMessageTimeout(t, connA, 2*time.Second)
	if nextStart.Type != MessageTypeTurnStarted {
		t.Fatalf("Expected turn_started message for Bob, got %s", nextStart.Type)
	}
}

func TestWS_QuestionContentAndCorrectAnswerStayOffSpectatorWire(t *testing.T) {
	provider := fixedQuestionProvider{question: room.Question{
		ID:     "q-wire-secret",
		Type:   "mcq",
		Prompt: "SECRET_PROMPT_ONLY_FOR_ALICE",
		Options: []room.QuestionOption{
			{ID: "SECRET_CORRECT_OPTION_ID", Text: "SECRET_CORRECT_OPTION_TEXT", Correct: true},
			{ID: "SECRET_WRONG_OPTION_ID", Text: "SECRET_WRONG_OPTION_TEXT"},
		},
	}}
	_, server := setupTestServerWithProvider(t, provider)

	connA := connectClient(t, server)
	defer connA.Close()
	connB := connectClient(t, server)
	defer connB.Close()

	sendJoin(t, connA, "Alice", "wire-room")
	_, _ = readUntilType(t, connA, MessageTypeStateSync)
	sendJoin(t, connB, "Bob", "wire-room")
	_, _ = readUntilType(t, connA, MessageTypeStateSync)
	_, _ = readUntilType(t, connB, MessageTypeStateSync)

	choose := Message{
		Type:    MessageTypeChooseLevel,
		RoomID:  "wire-room",
		Payload: json.RawMessage(`{"difficulty":"easy"}`),
	}
	if err := connA.WriteJSON(choose); err != nil {
		t.Fatalf("failed to choose level: %v", err)
	}

	_, activeQuestion := readUntilType(t, connA, MessageTypeQuestionStarted)
	if !bytes.Contains(activeQuestion.Payload, []byte("SECRET_PROMPT_ONLY_FOR_ALICE")) {
		t.Fatalf("active player did not receive the prompt: %s", activeQuestion.Payload)
	}
	spectatorRaw, spectatorQuestion := readUntilType(t, connB, MessageTypeQuestionStarted)
	for _, secret := range []string{
		"SECRET_PROMPT_ONLY_FOR_ALICE",
		"SECRET_CORRECT_OPTION_ID",
		"SECRET_CORRECT_OPTION_TEXT",
		"SECRET_WRONG_OPTION_ID",
		"SECRET_WRONG_OPTION_TEXT",
	} {
		if bytes.Contains(spectatorRaw, []byte(secret)) {
			t.Fatalf("spectator question_started wire payload leaked %q: %s", secret, spectatorRaw)
		}
	}
	var spectatorPayload map[string]any
	if err := json.Unmarshal(spectatorQuestion.Payload, &spectatorPayload); err != nil {
		t.Fatalf("invalid spectator question payload: %v", err)
	}
	for _, forbiddenField := range []string{"problem_id", "type", "prompt", "options"} {
		if _, present := spectatorPayload[forbiddenField]; present {
			t.Fatalf("spectator question payload contains forbidden field %q: %s", forbiddenField, spectatorQuestion.Payload)
		}
	}

	submit := Message{
		Type:    MessageTypeSubmitAnswer,
		RoomID:  "wire-room",
		Payload: json.RawMessage(`{"problem_id":"q-wire-secret","answer":["SECRET_CORRECT_OPTION_ID"]}`),
	}
	if err := connA.WriteJSON(submit); err != nil {
		t.Fatalf("failed to submit answer: %v", err)
	}

	spectatorResultRaw, spectatorResult := readUntilType(t, connB, MessageTypeAnswerResult)
	for _, secret := range []string{
		"correct_answer",
		"SECRET_CORRECT_OPTION_ID",
		"SECRET_CORRECT_OPTION_TEXT",
		"SECRET_PROMPT_ONLY_FOR_ALICE",
	} {
		if bytes.Contains(spectatorResultRaw, []byte(secret)) {
			t.Fatalf("spectator answer_result wire payload leaked %q: %s", secret, spectatorResultRaw)
		}
	}
	var publicResult map[string]any
	if err := json.Unmarshal(spectatorResult.Payload, &publicResult); err != nil {
		t.Fatalf("invalid spectator answer_result payload: %v", err)
	}
	if _, present := publicResult["correct_answer"]; present {
		t.Fatalf("spectator answer_result contains correct_answer: %s", spectatorResult.Payload)
	}

	// Alice receives the private review copy after the public result.
	_, _ = readUntilType(t, connA, MessageTypeAnswerResult)
	privateResultRaw, _ := readUntilType(t, connA, MessageTypeAnswerResult)
	if !bytes.Contains(privateResultRaw, []byte("SECRET_CORRECT_OPTION_ID")) {
		t.Fatalf("active player did not receive the private correct-answer review: %s", privateResultRaw)
	}
}
