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

func TestWS_AdminClientRejectsPlayerActions(t *testing.T) {
	hub := NewHub(nil)
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})
	mux.HandleFunc("/admin/ws", AdminHandler(hub, func(token string) bool { return token == "secret" }))

	server := httptest.NewServer(mux)
	defer server.Close()
	defer hub.Stop()

	// 1. connect normal player Alice
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	connAlice, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect Alice: %v", err)
	}
	defer connAlice.Close()
	sendJoin(t, hub, connAlice, "Alice", "admin-test-room")

	// Wait for presence
	_ = readMessageTimeout(t, connAlice, 2*time.Second)

	// Start game
	hub.GetRoomInstance("admin-test-room").AdminStartGame()

	// 2. connect Admin
	adminWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/admin/ws?token=secret&room_id=admin-test-room"
	connAdmin, _, err := websocket.DefaultDialer.Dial(adminWSURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect Admin: %v", err)
	}
	defer connAdmin.Close()

	// Admin tries to choose level
	chooseMsg := Message{
		Type:    MessageTypeChooseLevel,
		RoomID:  "admin-test-room",
		Payload: json.RawMessage(`{"difficulty":"easy"}`),
	}
	if err := connAdmin.WriteJSON(chooseMsg); err != nil {
		t.Fatalf("Failed to send choose_level from Admin: %v", err)
	}

	// Admin tries to submit answer
	submitMsg := Message{
		Type:    MessageTypeSubmitAnswer,
		RoomID:  "admin-test-room",
		Payload: json.RawMessage(`{"answer":"test"}`),
	}
	if err := connAdmin.WriteJSON(submitMsg); err != nil {
		t.Fatalf("Failed to send submit_answer from Admin: %v", err)
	}

	// Verify Admin receives no error or response from these player actions
	// (they should be silently dropped with a log). We might receive state_sync though.
	_ = connAdmin.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	for {
		var m Message
		if err := connAdmin.ReadJSON(&m); err == nil {
			if m.Type == MessageTypeError {
				t.Fatalf("Admin received error message when it should have been silently rejected: %s", m.Error)
			} else if m.Type != MessageTypeStateSync {
				t.Fatalf("Admin received unexpected message response: %s", m.Type)
			}
		} else {
			break
		}
	}

	// Double check room state hasn't changed to have an active question
	roomInstance := hub.GetRoomInstance("admin-test-room")
	activePlayerID, questionActive, _ := roomInstance.GetTurnState()
	if questionActive {
		t.Errorf("Question became active from admin command!")
	}
	if activePlayerID == "" {
		t.Errorf("Expected an active player, got empty string")
	}
}

func TestWS_AdminKickPlayerClosesConnection(t *testing.T) {
	hub := NewHub(nil)
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})
	mux.HandleFunc("/admin/ws", AdminHandler(hub, func(token string) bool { return token == "secret" }))

	server := httptest.NewServer(mux)
	defer server.Close()
	defer hub.Stop()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	connAlice, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("connect Alice: %v", err)
	}
	defer connAlice.Close()
	sendJoin(t, hub, connAlice, "Alice", "kick-room")

	_, joinedMsg := readUntilType(t, connAlice, MessageTypeJoined)
	var joined JoinedPayload
	_ = json.Unmarshal(joinedMsg.Payload, &joined)
	aliceID := joined.PlayerID
	_, _ = readUntilType(t, connAlice, MessageTypeStateSync)

	connBob, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("connect Bob: %v", err)
	}
	defer connBob.Close()
	sendJoin(t, hub, connBob, "Bob", "kick-room")
	_, _ = readUntilType(t, connBob, MessageTypeJoined)
	_, _ = readUntilType(t, connBob, MessageTypeStateSync)

	// Connect Admin
	adminWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/admin/ws?token=secret&room_id=kick-room"
	connAdmin, _, err := websocket.DefaultDialer.Dial(adminWSURL, nil)
	if err != nil {
		t.Fatalf("connect Admin: %v", err)
	}
	defer connAdmin.Close()

	// Admin kicks Alice
	kickPayload, _ := json.Marshal(AdminKickPayload{PlayerID: aliceID})
	_ = connAdmin.WriteJSON(Message{Type: MessageTypeAdminKick, RoomID: "kick-room", Payload: kickPayload})

	// Alice should receive kick error and connection should close
	_ = connAlice.SetReadDeadline(time.Now().Add(2 * time.Second))
	for {
		_, _, err := connAlice.ReadMessage()
		if err != nil {
			// Connection closed as expected
			break
		}
	}

	// Bob should receive state_sync where Alice is no longer in players list
	_, syncBob := readUntilType(t, connBob, MessageTypeStateSync)
	var payload StateSyncPayload
	_ = json.Unmarshal(syncBob.Payload, &payload)
	for _, p := range payload.Players {
		if p.ID == aliceID {
			t.Errorf("Alice still in players list after kick!")
		}
	}
}

func TestWS_AdminPauseAndResumeSyncs(t *testing.T) {
	hub := NewHub(nil)
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})
	mux.HandleFunc("/admin/ws", AdminHandler(hub, func(token string) bool { return token == "secret" }))

	server := httptest.NewServer(mux)
	defer server.Close()
	defer hub.Stop()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	connAlice, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	defer connAlice.Close()
	sendJoin(t, hub, connAlice, "Alice", "pause-test-room")
	_, _ = readUntilType(t, connAlice, MessageTypeJoined)
	_, _ = readUntilType(t, connAlice, MessageTypeStateSync)

	adminWSURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/admin/ws?token=secret&room_id=pause-test-room"
	connAdmin, _, _ := websocket.DefaultDialer.Dial(adminWSURL, nil)
	defer connAdmin.Close()

	// Start game via Admin WS
	_ = connAdmin.WriteJSON(Message{Type: MessageTypeAdminStart, RoomID: "pause-test-room", Payload: json.RawMessage(`{}`)})
	_, syncStarted := readUntilType(t, connAlice, MessageTypeStateSync)
	var payloadStarted StateSyncPayload
	_ = json.Unmarshal(syncStarted.Payload, &payloadStarted)
	if !payloadStarted.IsStarted {
		t.Fatalf("expected is_started=true after admin_start")
	}

	// Pause game via Admin WS
	_ = connAdmin.WriteJSON(Message{Type: MessageTypeAdminPause, RoomID: "pause-test-room", Payload: json.RawMessage(`{}`)})
	_, syncPaused := readUntilType(t, connAlice, MessageTypeStateSync)
	var payloadPaused StateSyncPayload
	_ = json.Unmarshal(syncPaused.Payload, &payloadPaused)
	if !payloadPaused.IsPaused {
		t.Fatalf("expected is_paused=true after admin_pause")
	}

	// Player cannot choose level while paused
	_ = connAlice.WriteJSON(Message{Type: MessageTypeChooseLevel, RoomID: "pause-test-room", Payload: json.RawMessage(`{"difficulty":"easy"}`)})
	_, errMsg := readUntilType(t, connAlice, MessageTypeError)
	if !strings.Contains(errMsg.Error, "paused") {
		t.Fatalf("expected paused error, got: %s", errMsg.Error)
	}

	// Resume game
	_ = connAdmin.WriteJSON(Message{Type: MessageTypeAdminPause, RoomID: "pause-test-room", Payload: json.RawMessage(`{}`)})
	_, syncResumed := readUntilType(t, connAlice, MessageTypeStateSync)
	var payloadResumed StateSyncPayload
	_ = json.Unmarshal(syncResumed.Payload, &payloadResumed)
	if payloadResumed.IsPaused {
		t.Fatalf("expected is_paused=false after resume")
	}
}
