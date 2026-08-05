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
	sendJoin(t, connAlice, "Alice", "admin-test-room")

	// Wait for presence
	_ = readMessageTimeout(t, connAlice, 2*time.Second)

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
