package ws

import (
	"encoding/json"
	"time"
)

// Message types
const (
	MessageTypeJoin         = "join"
	MessageTypeStateSync    = "state_sync"
	MessageTypePresence     = "presence"
	MessageTypeError        = "error"
	MessageTypePing         = "ping"
	MessageTypePong         = "pong"
	MessageTypeChooseLevel  = "choose_level"
	MessageTypeSubmitAnswer = "submit_answer"
	MessageTypeTurnStarted  = "turn_started"
	MessageTypeTurnEnded    = "turn_ended"
	MessageTypeRollResolved = "roll_resolved"
)

// Message is the standard WebSocket JSON frame wrapper.
type Message struct {
	Version int             `json:"v,omitempty"`
	Type    string          `json:"type"`
	RoomID  string          `json:"room_id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// ChooseLevelPayload represents the payload to select question difficulty.
type ChooseLevelPayload struct {
	Difficulty string `json:"difficulty"`
}

// SubmitAnswerPayload represents the payload to submit an answer.
type SubmitAnswerPayload struct {
	ProblemID string          `json:"problem_id,omitempty"`
	Answer    json.RawMessage `json:"answer,omitempty"`
}

// JoinPayload represents the payload sent by a client to join a room.
type JoinPayload struct {
	Name   string `json:"name"`
	RoomID string `json:"room_id,omitempty"`
}

// PlayerInfo represents player state sent over WebSocket.
type PlayerInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	RoomID      string    `json:"room_id"`
	JoinedAt    time.Time `json:"joined_at"`
	IsConnected bool      `json:"is_connected"`
}

// StateSyncPayload represents the complete player list for a room.
type StateSyncPayload struct {
	RoomID  string       `json:"room_id"`
	Players []PlayerInfo `json:"players"`
}

// PresencePayload represents a presence broadcast (join/leave).
type PresencePayload struct {
	Event  string     `json:"event"` // "joined" or "left"
	Player PlayerInfo `json:"player"`
}

// NewMessage creates a serialized Message.
func NewMessage(msgType string, roomID string, payload interface{}) ([]byte, error) {
	var payloadBytes json.RawMessage
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		payloadBytes = b
	}

	msg := Message{
		Version: 1,
		Type:    msgType,
		RoomID:  roomID,
		Payload: payloadBytes,
	}

	return json.Marshal(msg)
}

// NewErrorMessage creates a serialized error Message.
func NewErrorMessage(errMsg string) []byte {
	msg := Message{
		Version: 1,
		Type:    MessageTypeError,
		Error:   errMsg,
	}
	b, _ := json.Marshal(msg)
	return b
}
