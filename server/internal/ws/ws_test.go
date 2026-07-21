package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupTestServer(t *testing.T, clientOpts ...ClientOptions) (*Hub, *httptest.Server) {
	t.Helper()
	hub := NewHub()
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
		Type:   MessageTypeJoin,
		RoomID: roomID,
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

func TestWS_JoinFlowAndStateSync(t *testing.T) {
	_, server := setupTestServer(t)

	connA := connectClient(t, server)
	defer connA.Close()

	sendJoin(t, connA, "Alice", "room-1")

	// Expect presence broadcast for Alice
	msg1 := readMessageTimeout(t, connA, 2*time.Second)
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
	msg2 := readMessageTimeout(t, connA, 2*time.Second)
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
	_ = readMessageTimeout(t, connA, 2*time.Second) // presence A
	_ = readMessageTimeout(t, connA, 2*time.Second) // state_sync [Alice]

	// Connect Client B
	connB := connectClient(t, server)
	defer connB.Close()
	sendJoin(t, connB, "Bob", "game-room")

	// Client A should receive presence for Bob joining
	msgA_presence := readMessageTimeout(t, connA, 2*time.Second)
	if msgA_presence.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence for Bob, got %s", msgA_presence.Type)
	}
	var presenceB PresencePayload
	_ = json.Unmarshal(msgA_presence.Payload, &presenceB)
	if presenceB.Event != "joined" || presenceB.Player.Name != "Bob" {
		t.Errorf("Client A received wrong presence: %+v", presenceB)
	}

	// Client A should receive updated state_sync containing both Alice and Bob
	msgA_sync := readMessageTimeout(t, connA, 2*time.Second)
	if msgA_sync.Type != MessageTypeStateSync {
		t.Fatalf("Client A expected state_sync, got %s", msgA_sync.Type)
	}
	var stateSync1 StateSyncPayload
	_ = json.Unmarshal(msgA_sync.Payload, &stateSync1)
	if len(stateSync1.Players) != 2 {
		t.Errorf("Expected 2 players in state_sync for Client A, got %d", len(stateSync1.Players))
	}

	// Client B should also receive presence and state_sync
	msgB_presence := readMessageTimeout(t, connB, 2*time.Second)
	if msgB_presence.Type != MessageTypePresence {
		t.Fatalf("Client B expected presence message, got %s", msgB_presence.Type)
	}
	msgB_sync := readMessageTimeout(t, connB, 2*time.Second)
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
	_ = readMessageTimeout(t, connA, 2*time.Second) // presence
	_ = readMessageTimeout(t, connA, 2*time.Second) // state_sync

	connB := connectClient(t, server)
	sendJoin(t, connB, "Bob", "room-disconnect")
	_ = readMessageTimeout(t, connA, 2*time.Second) // presence B
	_ = readMessageTimeout(t, connA, 2*time.Second) // state_sync [A, B]
	_ = readMessageTimeout(t, connB, 2*time.Second) // presence B
	_ = readMessageTimeout(t, connB, 2*time.Second) // state_sync [A, B]

	// Graceful disconnect Client B
	_ = connB.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "leaving"))
	_ = connB.Close()

	// Client A should receive presence ("left") for Bob
	msgA_left := readMessageTimeout(t, connA, 2*time.Second)
	if msgA_left.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence message on disconnect, got %s", msgA_left.Type)
	}
	var presenceLeft PresencePayload
	_ = json.Unmarshal(msgA_left.Payload, &presenceLeft)
	if presenceLeft.Event != "left" || presenceLeft.Player.Name != "Bob" {
		t.Errorf("Unexpected leave presence: %+v", presenceLeft)
	}

	// Client A should receive updated state_sync with only Alice
	msgA_sync := readMessageTimeout(t, connA, 2*time.Second)
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
	_ = readMessageTimeout(t, connA, 2*time.Second) // presence
	_ = readMessageTimeout(t, connA, 2*time.Second) // state_sync

	connB := connectClient(t, server)
	sendJoin(t, connB, "Bob", "room-ungraceful")
	_ = readMessageTimeout(t, connA, 2*time.Second) // presence B
	_ = readMessageTimeout(t, connA, 2*time.Second) // state_sync [A, B]
	_ = readMessageTimeout(t, connB, 2*time.Second) // presence B
	_ = readMessageTimeout(t, connB, 2*time.Second) // state_sync [A, B]

	// Simulate ungraceful disconnect by abruptly closing underlying TCP connection without sending close frame
	_ = connB.UnderlyingConn().Close()

	// Client A should receive presence ("left") and state_sync automatically after pong timeout
	msgA_left := readMessageTimeout(t, connA, 2*time.Second)
	if msgA_left.Type != MessageTypePresence {
		t.Fatalf("Client A expected presence message on ungraceful disconnect, got %s", msgA_left.Type)
	}
	var presenceLeft PresencePayload
	_ = json.Unmarshal(msgA_left.Payload, &presenceLeft)
	if presenceLeft.Event != "left" || presenceLeft.Player.Name != "Bob" {
		t.Errorf("Unexpected leave presence: %+v", presenceLeft)
	}

	msgA_sync := readMessageTimeout(t, connA, 2*time.Second)
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
