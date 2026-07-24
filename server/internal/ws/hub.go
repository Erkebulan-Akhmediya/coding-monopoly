package ws

import (
	"log"
	"sync"

	"server/internal/room"
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

// HubBroadcaster implements room.Broadcaster interface for WebSockets.
type HubBroadcaster struct {
	hub *Hub
}

func (hb *HubBroadcaster) BroadcastRoom(roomID string, msgType string, payload any) {
	data, err := NewMessage(msgType, roomID, payload)
	if err == nil {
		hb.hub.deliverToRoom(roomID, data)
	}
}

func (hb *HubBroadcaster) BroadcastRoomExcept(roomID string, excludedClientID string, msgType string, payload any) {
	data, err := NewMessage(msgType, roomID, payload)
	if err == nil {
		hb.hub.deliverToRoomExcept(roomID, excludedClientID, data)
	}
}

func (hb *HubBroadcaster) SendToPlayer(roomID string, clientID string, msgType string, payload any) {
	data, err := NewMessage(msgType, roomID, payload)
	if err != nil {
		return
	}
	hb.hub.mu.RLock()
	defer hb.hub.mu.RUnlock()
	for c := range hb.hub.clients {
		if c.GetID() == clientID && c.GetRoomID() == roomID && c.IsJoined() {
			c.SendBytes(data)
			return
		}
	}
}

func (hb *HubBroadcaster) SendError(clientID string, errMsg string) {
	hb.hub.mu.RLock()
	var targetClient *Client
	for c := range hb.hub.clients {
		if c.GetID() == clientID {
			targetClient = c
			break
		}
	}
	hb.hub.mu.RUnlock()

	if targetClient != nil {
		targetClient.sendError(errMsg)
	}
}

// Hub maintains the set of active connections and broadcasts messages to clients.
type Hub struct {
	// Registered clients: map of active client connections.
	clients map[*Client]bool

	// Rooms map: roomID -> map of joined clients in that room.
	rooms map[string]map[*Client]bool

	// Room engine instances: roomID -> *room.Room
	roomInstances    map[string]*room.Room
	questionProvider room.QuestionProvider

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
func NewHub(providers ...room.QuestionProvider) *Hub {
	var provider room.QuestionProvider
	if len(providers) > 0 {
		provider = providers[0]
	}
	return &Hub{
		clients:          make(map[*Client]bool),
		rooms:            make(map[string]map[*Client]bool),
		roomInstances:    make(map[string]*room.Room),
		questionProvider: provider,
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		join:             make(chan *joinRequest),
		broadcast:        make(chan *broadcastRequest),
		stopChan:         make(chan struct{}),
	}
}

// Run starts the Hub event loop. It MUST be executed in its own single owning goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[WS Hub] Client registered: %s", client.id)

		case client := <-h.unregister:
			h.mu.Lock()
			roomID := client.GetRoomID()
			_, ok := h.clients[client]
			if ok {
				delete(h.clients, client)
				if roomID != "" {
					if roomClients, exists := h.rooms[roomID]; exists {
						delete(roomClients, client)
						if len(roomClients) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}
			}
			h.mu.Unlock()

			if ok {
				client.CloseSendChannel()
				wasJoined := client.IsJoined()
				log.Printf("[WS Hub] Client unregistered: %s (name: %s, room: %s)", client.id, client.GetName(), roomID)

				if wasJoined && roomID != "" {
					playerInfo := client.ToPlayerInfo(false)
					h.broadcastPresence(roomID, "left", playerInfo)
					h.broadcastStateSync(roomID)

					if r, exists := h.roomInstances[roomID]; exists {
						r.DisconnectPlayer(client.id)
					}
				}
			}

		case req := <-h.join:
			c := req.client
			h.mu.RLock()
			_, registered := h.clients[c]
			h.mu.RUnlock()

			if !registered {
				log.Printf("[WS Hub] Join rejected: client %s not registered", c.id)
				continue
			}

			roomID := req.roomID
			if roomID == "" {
				roomID = "default"
			}

			// If client was previously in a different room, clean up old room
			oldRoomID := c.GetRoomID()
			wasJoined := c.IsJoined()

			h.mu.Lock()
			if oldRoomID != "" && oldRoomID != roomID {
				if oldRoom, exists := h.rooms[oldRoomID]; exists {
					delete(oldRoom, c)
					if len(oldRoom) == 0 {
						delete(h.rooms, oldRoomID)
					}
				}
				if r, exists := h.roomInstances[oldRoomID]; exists {
					r.DisconnectPlayer(c.id)
				}
			}

			// Add to new room
			if h.rooms[roomID] == nil {
				h.rooms[roomID] = make(map[*Client]bool)
			}
			h.rooms[roomID][c] = true
			h.mu.Unlock()

			// Update client state
			c.SetJoined(req.name, roomID)

			// Update room engine state
			r := h.GetRoomInstance(roomID)
			r.AddOrReconnectPlayer(c.id, req.name)

			if oldRoomID != "" && oldRoomID != roomID && wasJoined {
				playerInfo := c.ToPlayerInfo(false)
				h.broadcastPresence(oldRoomID, "left", playerInfo)
				h.broadcastStateSync(oldRoomID)
			}

			log.Printf("[WS Hub] Client %s (%s) joined room: %s", c.id, req.name, roomID)

			playerInfo := c.ToPlayerInfo(true)
			// Broadcast presence (joined) to room
			h.broadcastPresence(roomID, "joined", playerInfo)

			// Broadcast state_sync to all clients in room (including the joining client)
			h.broadcastStateSync(roomID)

		case req := <-h.broadcast:
			h.deliverToRoom(req.roomID, req.data)

		case <-h.stopChan:
			log.Printf("[WS Hub] Shutting down...")
			h.mu.Lock()
			for c := range h.clients {
				c.CloseSendChannel()
			}
			h.mu.Unlock()
			return
		}
	}
}

// Stop stops the hub loop safely.
func (h *Hub) Stop() {
	h.mu.Lock()
	select {
	case <-h.stopChan:
		// Already stopped
	default:
		close(h.stopChan)
	}
	h.mu.Unlock()
}

// RegisterClient queues a client for registration in the hub.
func (h *Hub) RegisterClient(c *Client) {
	select {
	case h.register <- c:
	case <-h.stopChan:
	}
}

// UnregisterClient queues a client for unregistration in the hub.
func (h *Hub) UnregisterClient(c *Client) {
	select {
	case h.unregister <- c:
	case <-h.stopChan:
	}
}

// JoinRoom queues a join request for a client.
func (h *Hub) JoinRoom(c *Client, name string, roomID string) {
	select {
	case h.join <- &joinRequest{
		client: c,
		name:   name,
		roomID: roomID,
	}:
	case <-h.stopChan:
	}
}

// BroadcastRoom queues a message to be sent to all joined clients in a room.
func (h *Hub) BroadcastRoom(roomID string, data []byte) {
	select {
	case h.broadcast <- &broadcastRequest{
		roomID: roomID,
		data:   data,
	}:
	case <-h.stopChan:
	}
}

func (h *Hub) deliverToRoom(roomID string, data []byte) {
	h.deliverToRoomExcept(roomID, "", data)
}

func (h *Hub) deliverToRoomExcept(roomID string, excludedClientID string, data []byte) {
	h.mu.RLock()
	var targetClients []*Client
	if roomID != "" {
		if roomClients, ok := h.rooms[roomID]; ok {
			for c := range roomClients {
				targetClients = append(targetClients, c)
			}
		}
	} else {
		for c := range h.clients {
			targetClients = append(targetClients, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targetClients {
		if excludedClientID != "" && c.GetID() == excludedClientID {
			continue
		}
		c.SendBytes(data)
	}
}

// Internal helper to broadcast presence to a room.
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

	h.mu.RLock()
	var targetClients []*Client
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			targetClients = append(targetClients, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targetClients {
		c.SendBytes(data)
	}
}

// Internal helper to broadcast state_sync to a room.
func (h *Hub) broadcastStateSync(roomID string) {
	h.mu.RLock()
	var players []PlayerInfo
	var targetClients []*Client
	var cells []room.BoardCell
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			if c.IsJoined() {
				players = append(players, c.ToPlayerInfo(true))
			}
			targetClients = append(targetClients, c)
		}
	}
	r, ok := h.roomInstances[roomID]
	if ok {
		cells = r.Board()
		roomPlayers := r.GetPlayers()
		playerMap := make(map[string]room.Player)
		for _, rp := range roomPlayers {
			playerMap[rp.ID] = rp
		}
		for i := range players {
			if rp, found := playerMap[players[i].ID]; found {
				players[i].Position = rp.Position
				players[i].XP = rp.XP
				players[i].InCodeFreeze = rp.InCodeFreeze
				players[i].SkipNextTurn = rp.SkipNextTurn
				players[i].DoubleXP = rp.DoubleXP
				players[i].FreePasses = rp.FreePasses
			}
		}
	}
	h.mu.RUnlock()

	payload := StateSyncPayload{
		RoomID:     roomID,
		Players:    players,
		BoardCells: cells,
	}

	data, err := NewMessage(MessageTypeStateSync, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating state_sync message: %v", err)
		return
	}

	for _, c := range targetClients {
		c.SendBytes(data)
	}
}

// GetRoomPlayers returns a snapshot of connected players in a room (thread-safe utility).
func (h *Hub) GetRoomPlayers(roomID string) []PlayerInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var players []PlayerInfo
	var playerMap map[string]room.Player
	if r, ok := h.roomInstances[roomID]; ok {
		playerMap = make(map[string]room.Player)
		for _, rp := range r.GetPlayers() {
			playerMap[rp.ID] = rp
		}
	}

	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			if c.IsJoined() {
				pi := c.ToPlayerInfo(true)
				if rp, found := playerMap[c.id]; found {
					pi.Position = rp.Position
					pi.XP = rp.XP
					pi.InCodeFreeze = rp.InCodeFreeze
					pi.SkipNextTurn = rp.SkipNextTurn
					pi.DoubleXP = rp.DoubleXP
					pi.FreePasses = rp.FreePasses
				}
				players = append(players, pi)
			}
		}
	}
	return players
}

// GetRoomInstance retrieves an existing room.Room instance or creates a new one thread-safely.
func (h *Hub) GetRoomInstance(roomID string) *room.Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	r, ok := h.roomInstances[roomID]
	if !ok {
		hb := &HubBroadcaster{hub: h}
		r = room.NewRoomWithQuestionProvider(roomID, hb, h.questionProvider)
		h.roomInstances[roomID] = r
	}
	return r
}
