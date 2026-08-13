package ws

import (
	"log"
	"sync"
	"time"

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
					if r, exists := h.roomInstances[roomID]; exists {
						r.DisconnectPlayer(client.id)
					}
					playerInfo := client.ToPlayerInfo(false)
					h.broadcastPresence(roomID, "left", playerInfo)
					h.broadcastStateSync(roomID)
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

			// Mid-question resume: only the active player gets full question content,
			// using the original deadline (remaining time), not a reset countdown.
			if qPayload := r.GetActiveQuestionPayload(c.id); qPayload != nil {
				qData, err := NewMessage(MessageTypeQuestionStarted, roomID, qPayload)
				if err == nil {
					c.SendBytes(qData)
				}
			}

			if goPayload := r.GetGameOver(); goPayload != nil {
				data, err := NewMessage(MessageTypeGameOver, roomID, goPayload)
				if err == nil {
					c.SendBytes(data)
				}
			}

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
	payload, targetClients := h.buildStateSyncPayload(roomID)

	data, err := NewMessage(MessageTypeStateSync, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating state_sync message: %v", err)
		return
	}

	for _, c := range targetClients {
		c.SendBytes(data)
	}
}

// sendStateSyncToClient sends a state_sync message to a single client (used on reconnect).
func (h *Hub) sendStateSyncToClient(roomID string, target *Client) {
	payload, _ := h.buildStateSyncPayload(roomID)

	data, err := NewMessage(MessageTypeStateSync, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating state_sync message for client %s: %v", target.GetID(), err)
		return
	}

	target.SendBytes(data)

	// Resync the active question to the active player if a question is currently active.
	h.mu.RLock()
	r, ok := h.roomInstances[roomID]
	h.mu.RUnlock()

	if ok {
		if qPayload := r.GetActiveQuestionPayload(target.GetID()); qPayload != nil {
			qData, err := NewMessage(MessageTypeQuestionStarted, roomID, qPayload)
			if err == nil {
				target.SendBytes(qData)
			}
		}
	}
}

// buildStateSyncPayload assembles the StateSyncPayload and the list of connected
// clients for a room. Shared by broadcastStateSync and sendStateSyncToClient.
// Players come from the room engine (including disconnected slots) so reconnect
// clients recover exact position/XP/turn order.
func (h *Hub) buildStateSyncPayload(roomID string) (StateSyncPayload, []*Client) {
	h.mu.RLock()
	var targetClients []*Client
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			targetClients = append(targetClients, c)
		}
	}

	var players []PlayerInfo
	var cells []room.BoardCell
	var currentTurnPlayer string
	var questionActive bool
	var deadline *time.Time
	var targetXP int
	var gameOver *room.GameOverPayload

	r, ok := h.roomInstances[roomID]
	h.mu.RUnlock()

	if ok {
		cells = r.Board()
		roomPlayers := r.GetPlayers()
		players = make([]PlayerInfo, 0, len(roomPlayers))
		for _, rp := range roomPlayers {
			players = append(players, PlayerInfo{
				ID:           rp.ID,
				Name:         rp.Name,
				RoomID:       roomID,
				JoinedAt:     rp.JoinedAt,
				IsConnected:  rp.IsConnected,
				Position:     rp.Position,
				XP:           rp.XP,
				InCodeFreeze: rp.InCodeFreeze,
				SkipNextTurn: rp.SkipNextTurn,
				DoubleXP:     rp.DoubleXP,
				FreePasses:   rp.FreePasses,
			})
		}
		currentTurnPlayer, questionActive, deadline = r.GetTurnState()
		targetXP = r.GetTargetXP()
		if r.IsFinished() {
			gameOver = r.GetGameOver()
		}
	}

	payload := StateSyncPayload{
		RoomID:            roomID,
		Players:           players,
		BoardCells:        cells,
		CurrentTurnPlayer: currentTurnPlayer,
		QuestionActive:    questionActive,
		Deadline:          deadline,
		TargetXP:          targetXP,
		GameOver:          gameOver,
	}

	return payload, targetClients
}

// GetRoomPlayers returns a snapshot of all players in a room (including
// disconnected slots) from the room engine.
func (h *Hub) GetRoomPlayers(roomID string) []PlayerInfo {
	h.mu.RLock()
	r, ok := h.roomInstances[roomID]
	h.mu.RUnlock()
	if !ok {
		return nil
	}

	roomPlayers := r.GetPlayers()
	players := make([]PlayerInfo, 0, len(roomPlayers))
	for _, rp := range roomPlayers {
		players = append(players, PlayerInfo{
			ID:           rp.ID,
			Name:         rp.Name,
			RoomID:       roomID,
			JoinedAt:     rp.JoinedAt,
			IsConnected:  rp.IsConnected,
			Position:     rp.Position,
			XP:           rp.XP,
			InCodeFreeze: rp.InCodeFreeze,
			SkipNextTurn: rp.SkipNextTurn,
			DoubleXP:     rp.DoubleXP,
			FreePasses:   rp.FreePasses,
		})
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

// RegisterAdminClient registers an already-authenticated admin spectator client
// and adds it to the named room's client set so it receives all broadcasts.
// Unlike the regular JoinRoom flow, no player record is created in the room engine.
func (h *Hub) RegisterAdminClient(c *Client, roomID string) {
	h.mu.Lock()
	h.clients[c] = true
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][c] = true
	h.mu.Unlock()

	// Ensure a room engine exists so state_sync works.
	_ = h.GetRoomInstance(roomID)

	// Send an immediate state snapshot to the new admin spectator.
	h.sendStateSyncToClient(roomID, c)
	log.Printf("[WS Hub] Admin client %s registered in room %s", c.id, roomID)
}

// broadcastGameEvent sends a game_event message to every admin spectator watching a room.
func (h *Hub) broadcastGameEvent(roomID string, kind string, message string, meta any) {
	payload := GameEventPayload{
		Kind:      kind,
		Message:   message,
		Timestamp: timeNow(),
		Meta:      meta,
	}
	data, err := NewMessage(MessageTypeGameEvent, roomID, payload)
	if err != nil {
		log.Printf("[WS Hub] Error creating game_event message: %v", err)
		return
	}
	h.broadcastToAdmins(roomID, data)
}

// broadcastToAdmins delivers data only to admin spectator clients in a room.
func (h *Hub) broadcastToAdmins(roomID string, data []byte) {
	h.mu.RLock()
	var targets []*Client
	if roomClients, ok := h.rooms[roomID]; ok {
		for c := range roomClients {
			if c.IsAdmin() {
				targets = append(targets, c)
			}
		}
	}
	h.mu.RUnlock()
	for _, c := range targets {
		c.SendBytes(data)
	}
}

// timeNow is a replaceable clock for testing.
var timeNow = func() time.Time { return time.Now() }
