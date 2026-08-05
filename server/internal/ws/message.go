package ws

import (
	"encoding/json"
	"time"

	"server/internal/room"
)

// Message types
const (
	MessageTypeJoin            = "join"
	MessageTypeStateSync       = "state_sync"
	MessageTypePresence        = "presence"
	MessageTypeError           = "error"
	MessageTypePing            = "ping"
	MessageTypePong            = "pong"
	MessageTypeChooseLevel     = "choose_level"
	MessageTypeSubmitAnswer    = "submit_answer"
	MessageTypeTurnStarted     = "turn_started"
	MessageTypeTurnEnded       = "turn_ended"
	MessageTypeRollResolved    = "roll_resolved"
	MessageTypeQuestionStarted = "question_started"
	MessageTypeAnswerResult    = "answer_result"
	MessageTypeStateRequest    = "state_request"

	// Admin-only message types (sent by admin spectator clients).
	MessageTypeAdminJoin      = "admin_join"      // authenticates the spectator WS connection
	MessageTypeAdminStart     = "admin_start"     // start the game
	MessageTypeAdminPause     = "admin_pause"     // pause / resume the game
	MessageTypeAdminKick      = "admin_kick"      // remove a player from the room
	MessageTypeAdminSkipTurn  = "admin_skip_turn" // manually advance past the active player's turn

	// Server → admin event feed entry.
	MessageTypeGameEvent = "game_event"
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
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RoomID       string    `json:"room_id"`
	JoinedAt     time.Time `json:"joined_at"`
	IsConnected  bool      `json:"is_connected"`
	Position     int       `json:"position"`
	XP           int       `json:"xp"`
	InCodeFreeze bool      `json:"in_code_freeze,omitempty"`
	SkipNextTurn bool      `json:"skip_next_turn,omitempty"`
	DoubleXP     bool      `json:"double_xp,omitempty"`
	FreePasses   int       `json:"free_passes,omitempty"`
}

// StateSyncPayload represents the complete player list for a room.
type StateSyncPayload struct {
	RoomID            string           `json:"room_id"`
	Players           []PlayerInfo     `json:"players"`
	BoardCells        []room.BoardCell `json:"board_cells"`
	CurrentTurnPlayer string           `json:"current_turn_player,omitempty"`
	QuestionActive    bool             `json:"question_active"`
	Deadline          *time.Time       `json:"deadline,omitempty"`
}

// PresencePayload represents a presence broadcast (join/leave).
type PresencePayload struct {
	Event  string     `json:"event"` // "joined" or "left"
	Player PlayerInfo `json:"player"`
}

// AdminJoinPayload is used by the admin WS client to identify itself.
type AdminJoinPayload struct {
	Token  string `json:"token"`
	RoomID string `json:"room_id,omitempty"`
}

// AdminKickPayload identifies the player to be kicked.
type AdminKickPayload struct {
	PlayerID string `json:"player_id"`
}

// AdminSkipTurnPayload optionally carries the target player ID (defaults to
// the current active player if omitted).
type AdminSkipTurnPayload struct {
	PlayerID string `json:"player_id,omitempty"`
}

// GameEventPayload is an entry in the live event feed pushed to admin spectators.
type GameEventPayload struct {
	Kind      string    `json:"kind"`       // e.g. "turn_started", "answer_result", "admin_action"
	Message   string    `json:"message"`    // human-readable description
	Timestamp time.Time `json:"timestamp"`
	Meta      any       `json:"meta,omitempty"` // optional raw payload for rich display
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
