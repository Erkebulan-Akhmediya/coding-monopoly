package room

import (
	"time"
)

// Player represents a participant in the Monopoly game room.
type Player struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	XP               int       `json:"xp"`
	Position         int       `json:"position"`
	IsConnected      bool      `json:"is_connected"`
	JoinedAt         time.Time `json:"joined_at"`
	SkipNextTurn     bool      `json:"skip_next_turn"`
	DoubleXP         bool      `json:"double_xp"`
	FreePasses       int       `json:"free_passes"`
	InCodeFreeze     bool      `json:"in_code_freeze"`
	ChosenDifficulty string    `json:"chosen_difficulty,omitempty"`
}

// NewPlayer creates a new player instance.
func NewPlayer(id string, name string) *Player {
	return &Player{
		ID:           id,
		Name:         name,
		XP:           0,
		Position:     0,
		IsConnected:  true,
		JoinedAt:     time.Now(),
		SkipNextTurn: false,
		DoubleXP:     false,
		FreePasses:   0,
		InCodeFreeze: false,
	}
}
