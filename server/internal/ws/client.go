package ws

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	defaultWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	defaultPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	defaultPingPeriod = (defaultPongWait * 9) / 10

	// Maximum message size allowed from peer.
	defaultMaxMessageSize = 4096
)

// Client represents an active WebSocket connection and player identity.
type Client struct {
	id     string
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	closed bool

	mu        sync.RWMutex
	name      string
	roomID    string
	isJoined  bool
	joinedAt  time.Time

	// Timing configurations (can be customized for testing)
	writeWait      time.Duration
	pongWait       time.Duration
	pingPeriod     time.Duration
	maxMessageSize int64
}

// ClientOptions allows overriding default websocket parameters (useful in tests).
type ClientOptions struct {
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int64
}

// NewClient constructs a new Client instance.
func NewClient(id string, hub *Hub, conn *websocket.Conn, opts ...ClientOptions) *Client {
	writeWait := defaultWriteWait
	pongWait := defaultPongWait
	pingPeriod := defaultPingPeriod
	maxMessageSize := int64(defaultMaxMessageSize)

	if len(opts) > 0 {
		if opts[0].WriteWait > 0 {
			writeWait = opts[0].WriteWait
		}
		if opts[0].PongWait > 0 {
			pongWait = opts[0].PongWait
		}
		if opts[0].PingPeriod > 0 {
			pingPeriod = opts[0].PingPeriod
		}
		if opts[0].MaxMessageSize > 0 {
			maxMessageSize = opts[0].MaxMessageSize
		}
	}

	return &Client{
		id:             id,
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		writeWait:      writeWait,
		pongWait:       pongWait,
		pingPeriod:     pingPeriod,
		maxMessageSize: maxMessageSize,
	}
}

// GetID returns client unique ID.
func (c *Client) GetID() string {
	return c.id
}

// SetJoined sets player name and room ID when joining.
func (c *Client) SetJoined(name string, roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.name = name
	c.roomID = roomID
	c.isJoined = true
	c.joinedAt = time.Now()
}

// GetName returns player name.
func (c *Client) GetName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name
}

// GetRoomID returns client room ID.
func (c *Client) GetRoomID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.roomID
}

// IsJoined returns whether client has completed join flow.
func (c *Client) IsJoined() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isJoined
}

// ToPlayerInfo converts Client state into a PlayerInfo struct.
func (c *Client) ToPlayerInfo(isConnected bool) PlayerInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return PlayerInfo{
		ID:          c.id,
		Name:        c.name,
		RoomID:      c.roomID,
		JoinedAt:    c.joinedAt,
		IsConnected: isConnected,
	}
}

// SendBytes places a raw byte slice into the client's send channel.
func (c *Client) SendBytes(data []byte) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return
	}
	select {
	case c.send <- data:
	default:
		log.Printf("[WS Client] Client %s send buffer full, dropping message and closing connection", c.id)
		go func() {
			_ = c.conn.Close()
		}()
	}
}

// CloseSendChannel safely closes the send channel once.
func (c *Client) CloseSendChannel() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		close(c.send)
	}
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// Application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(c.maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("[WS Client] Read error for client %s: %v", c.id, err)
			}
			break
		}

		c.handleIncomingMessage(messageData)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				// Hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleIncomingMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		c.sendError("Invalid JSON format")
		return
	}

	switch msg.Type {
	case MessageTypeJoin:
		c.handleJoinMessage(msg)
	case MessageTypePing:
		pong, _ := NewMessage(MessageTypePong, c.GetRoomID(), nil)
		c.SendBytes(pong)
	default:
		log.Printf("[WS Client] Received unhandled message type '%s' from client %s", msg.Type, c.id)
	}
}

func (c *Client) handleJoinMessage(msg Message) {
	var payload JoinPayload

	// Support both nested payload struct and top-level fields
	if len(msg.Payload) > 0 {
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			c.sendError("Invalid join payload")
			return
		}
	} else {
		// Fallback parse directly from msg if top-level name was provided
		var topLevel struct {
			Name   string `json:"name"`
			RoomID string `json:"room_id"`
		}
		if err := json.Unmarshal(msg.Payload, &topLevel); err == nil {
			payload.Name = topLevel.Name
			payload.RoomID = topLevel.RoomID
		}
	}

	// Validate name
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		c.sendError("Name is required to join")
		return
	}

	if len(name) > 30 {
		c.sendError("Name exceeds maximum allowed length of 30 characters")
		return
	}

	roomID := strings.TrimSpace(payload.RoomID)
	if roomID == "" {
		roomID = msg.RoomID
	}
	if roomID == "" {
		roomID = "default"
	}

	c.hub.JoinRoom(c, name, roomID)
}

func (c *Client) sendError(errMsg string) {
	errBytes := NewErrorMessage(errMsg)
	c.SendBytes(errBytes)
}
