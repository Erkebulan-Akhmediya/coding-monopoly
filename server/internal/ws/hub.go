package ws

import (
	"log"
	"sync"
)

// Authoritative Connection Set Concurrency Model:
// The Hub holds the authoritative connection set (all active clients and room memberships)
// managed exclusively by a single owning goroutine executing the Hub.Run() event loop.
// All client state mutations (register, unregister, join flow, presence updates, and state broadcasts)
// are received over channels and executed sequentially inside the single owning goroutine.
// This guarantees race-free, synchronized hub operations without needing coarse mutex locks.

type joinRequest struct {
	client *Client
	name   string
	roomID string
}

type broadcastRequest struct {
	roomID string
	data   []byte
}

// Hub maintains the set of active connections and broadcasts messages to clients.
type Hub struct {
	// Registered clients: map of active client connections.
	clients map[*Client]bool

	// Rooms map: roomID -> map of joined clients in that room.
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Join requests from clients.
	join chan *joinRequest

	// Channel to broadcast messages to clients in a room or all.
	broadcast chan *broadcastRequest

	// Channel to stop hub.
	stopChan chan struct{}

	// Optional mutex for read-only external inspection methods.
	mu sync.RWMutex
}

// NewHub creates and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		join:       make(chan *joinRequest),
		broadcast:  make(chan *broadcastRequest),
		stopChan:   make(chan struct{}),
	}
}

// Run starts the Hub event loop. It MUST be executed in its own single owning goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[WS Hub] Client registered: %s", client.id)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.CloseSendChannel()

				roomID := client.GetRoomID()
				wasJoined := client.IsJoined()

				if roomID != "" {
					if roomClients, exists := h.rooms[roomID]; exists {
						delete(roomClients, client)
						if len(roomClients) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}

				log.Printf("[WS Hub] Client unregistered: %s (name: %s, room: %s)", client.id, client.GetName(), roomID)

				if wasJoined && roomID != "" {
					playerInfo := client.ToPlayerInfo(false)
					h.broadcastPresence(roomID, "left", playerInfo)
					h.broadcastStateSync(roomID)
				}
			}

		case req := <-h.join:
			c := req.client
			if _, ok := h.clients[c]; !ok {
				log.Printf("[WS Hub] Join rejected: client %s not registered", c.id)
				continue
			}

			roomID := req.roomID
			if roomID == "" {
				roomID = "default"
			}

			// If client was previously in a different room, clean up old room
			oldRoomID := c.GetRoomID()
			if oldRoomID != "" && oldRoomID != roomID {
				if oldRoom, exists := h.rooms[oldRoomID]; exists {
					delete(oldRoom, c)
					if len(oldRoom) == 0 {
						delete(h.rooms, oldRoomID)
					}
				}
				if c.IsJoined() {
					playerInfo := c.ToPlayerInfo(false)
					h.broadcastPresence(oldRoomID, "left", playerInfo)
					h.broadcastStateSync(oldRoomID)
				}
			}

			// Update client state
			c.SetJoined(req.name, roomID)

			// Add to new room
			if h.rooms[roomID] == nil {
				h.rooms[roomID] = make(map[*Client]bool)
			}
			h.rooms[roomID][c] = true

			log.Printf("[WS Hub] Client %s (%s) joined room: %s", c.id, req.name, roomID)

			playerInfo := c.ToPlayerInfo(true)
			// Broadcast presence (joined) to room
			h.broadcastPresence(roomID, "joined", playerInfo)

			// Broadcast state_sync to all clients in room (including the joining client)
			h.broadcastStateSync(roomID)

		case req := <-h.broadcast:
			if req.roomID != "" {
				if roomClients, ok := h.rooms[req.roomID]; ok {
					for c := range roomClients {
						c.SendBytes(req.data)
					}
				}
			} else {
				for c := range h.clients {
					c.SendBytes(req.data)
				}
			}

		case <-h.stopChan:
			log.Printf("[WS Hub] Shutting down...")
			for c := range h.clients {
				c.CloseSendChannel()
			}
			return
		}
	}
}

// Stop stops the hub loop.
func (h *Hub) Stop() {
	close(h.stopChan)
}

// RegisterClient queues a client for registration in the hub.
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

// UnregisterClient queues a client for unregistration in the hub.
func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}

// JoinRoom queues a join request for a client.
func (h *Hub) JoinRoom(c *Client, name string, roomID string) {
	h.join <- &joinRequest{
		client: c,
		name:   name,
		roomID: roomID,
	}
}

// BroadcastRoom queues a message to be sent to all joined clients in a room.
func (h *Hub) BroadcastRoom(roomID string, data []byte) {
	h.broadcast <- &broadcastRequest{
		roomID: roomID,
		data:   data,
	}
}

// Internal helper to broadcast presence to a room (must be called inside single owning goroutine loop).
func (h *Hub) broadcastPresence(roomID string, event string, player PlayerInfo) {
	payload := PresencePayload{
		Event:  event,
		Player: player,
	}
	data, err := NewMessage(MessageTypePresence, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating presence message: %v", err)
		return
	}

	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			c.SendBytes(data)
		}
	}
}

// Internal helper to broadcast state_sync to a room (must be called inside single owning goroutine loop).
func (h *Hub) broadcastStateSync(roomID string) {
	var players []PlayerInfo
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			if c.IsJoined() {
				players = append(players, c.ToPlayerInfo(true))
			}
		}
	}

	payload := StateSyncPayload{
		RoomID:  roomID,
		Players: players,
	}

	data, err := NewMessage(MessageTypeStateSync, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating state_sync message: %v", err)
		return
	}

	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			c.SendBytes(data)
		}
	}
}

// GetRoomPlayers returns a snapshot of connected players in a room (thread-safe utility).
func (h *Hub) GetRoomPlayers(roomID string) []PlayerInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var players []PlayerInfo
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			if c.IsJoined() {
				players = append(players, c.ToPlayerInfo(true))
			}
		}
	}
	return players
}
