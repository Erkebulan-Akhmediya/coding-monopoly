package ws

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var defaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow cross-origin requests for local area network (LAN) access.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GenerateID produces a random hex string identifier for client connections.
func GenerateID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("client-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// ServeWS handles WebSocket requests from clients.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, opts ...ClientOptions) {
	conn, err := defaultUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Handler] Upgrade error: %v", err)
		return
	}

	clientID := GenerateID()
	client := NewClient(clientID, hub, conn, opts...)

	hub.RegisterClient(client)

	// Allow registered goroutines to process reading/writing
	go client.WritePump()
	go client.ReadPump()
}

// Handler returns an http.HandlerFunc bound to the given hub.
func Handler(hub *Hub, opts ...ClientOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r, opts...)
	}
}

// TokenValidator is a function that verifies an admin bearer token.
// It is implemented by the admin package's Handler.ValidateToken method.
type TokenValidator func(token string) bool

// ServeAdminWS upgrades the connection only after the admin token (passed in
// the "token" query parameter) has been validated. On success the client is
// flagged as admin and joined to the requested room without creating a player
// record in the room engine, ensuring it can never trigger player actions.
func ServeAdminWS(hub *Hub, validate TokenValidator, w http.ResponseWriter, r *http.Request, opts ...ClientOptions) {
	token := r.URL.Query().Get("token")
	if token == "" || !validate(token) {
		http.Error(w, "admin authentication required", http.StatusUnauthorized)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	if strings.TrimSpace(roomID) == "" {
		http.Error(w, "room_id is required", http.StatusBadRequest)
		return
	}
	roomID = strings.TrimSpace(roomID)

	conn, err := defaultUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS AdminHandler] Upgrade error: %v", err)
		return
	}

	clientID := "admin-" + GenerateID()
	client := NewClient(clientID, hub, conn, opts...)
	client.SetAdmin(roomID)

	hub.RegisterAdminClient(client, roomID)

	go client.WritePump()
	go client.ReadPump()
}

// AdminHandler returns an http.HandlerFunc for the admin spectator WebSocket endpoint.
func AdminHandler(hub *Hub, validate TokenValidator, opts ...ClientOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeAdminWS(hub, validate, w, r, opts...)
	}
}
