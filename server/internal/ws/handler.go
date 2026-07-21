package ws

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
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
