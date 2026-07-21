package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrNotActivePlayer   = errors.New("not active player's turn")
	ErrPlayerNotFound    = errors.New("player not found in room")
	ErrInvalidDifficulty = errors.New("invalid difficulty level, must be 'easy', 'medium', or 'hard'")
)

// Broadcaster provides an interface for Room to push messages to connected clients.
type Broadcaster interface {
	BroadcastRoom(roomID string, msgType string, payload any)
	SendError(clientID string, errMsg string)
}

// RollResult represents the details and outcome of a single dice roll.
type RollResult struct {
	PlayerID    string       `json:"player_id"`
	RollIndex   int          `json:"roll_index"`
	TotalRolls  int          `json:"total_rolls"`
	DieRoll     int          `json:"die_roll"`
	OldPosition int          `json:"old_position"`
	NewPosition int          `json:"new_position"`
	PassedGO    bool         `json:"passed_go"`
	LapBonus    int          `json:"lap_bonus"`
	LandedCell  BoardCell    `json:"landed_cell"`
	Effect      EffectResult `json:"effect"`
	PlayerXP    int          `json:"player_xp"`
}

// TurnStartedPayload represents the broadcast payload when a turn starts.
type TurnStartedPayload struct {
	ActivePlayerID string `json:"active_player_id"`
}

// TurnEndedPayload represents the broadcast payload when a turn ends.
type TurnEndedPayload struct {
	PlayerID string `json:"player_id"`
}

// Room manages game state, connected players in join order, turn progression, and cell effect execution.
type Room struct {
	ID             string
	mu             sync.RWMutex
	players        []*Player
	playerMap      map[string]*Player
	turnIdx        int
	activePlayerID string
	board          []BoardCell
	diceRng        *rand.Rand
	broadcaster    Broadcaster
}

// NewRoom creates a new Room with default board and initialized RNG.
func NewRoom(id string, broadcaster Broadcaster) *Room {
	return &Room{
		ID:          id,
		players:     make([]*Player, 0),
		playerMap:   make(map[string]*Player),
		turnIdx:     0,
		board:       DefaultBoard(),
		diceRng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		broadcaster: broadcaster,
	}
}

// SetRNG Seed or custom RNG generator (useful for deterministic testing).
func (r *Room) SetRNG(rng *rand.Rand) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diceRng = rng
}

// RNGIntn returns a random integer using the room's RNG (thread-safe).
func (r *Room) RNGIntn(n int) int {
	return r.diceRng.Intn(n)
}

// GetCell returns the board cell at a given position.
func (r *Room) GetCell(pos int) BoardCell {
	idx := pos % len(r.board)
	if idx < 0 {
		idx += len(r.board)
	}
	return r.board[idx]
}

// GetPlayers returns a copy of the ordered player list.
func (r *Room) GetPlayers() []Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Player, len(r.players))
	for i, p := range r.players {
		result[i] = *p
	}
	return result
}

// GetActivePlayerID returns the current active player's ID.
func (r *Room) GetActivePlayerID() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.activePlayerID
}

// AddOrReconnectPlayer handles player join or reconnect, maintaining strict join order.
func (r *Room) AddOrReconnectPlayer(clientID string, name string) (*Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if p, exists := r.playerMap[clientID]; exists {
		// Reconnecting player
		p.IsConnected = true
		if name != "" {
			p.Name = name
		}

		// If no active player set, make reconnecting player active if turn index points to them
		if r.activePlayerID == "" {
			r.activePlayerID = p.ID
			r.broadcastTurnStartedLocked()
		}

		return p, false
	}

	// New player joining
	player := NewPlayer(clientID, name)
	r.players = append(r.players, player)
	r.playerMap[clientID] = player

	// First connected player becomes active player
	isFirst := len(r.players) == 1 || r.activePlayerID == ""
	if isFirst {
		r.turnIdx = len(r.players) - 1
		r.activePlayerID = player.ID
		r.broadcastTurnStartedLocked()
	}

	return player, isFirst
}

// DisconnectPlayer marks a player as disconnected without removing their slot.
func (r *Room) DisconnectPlayer(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.playerMap[clientID]
	if !exists {
		return
	}

	p.IsConnected = false

	// If the active player disconnected, advance turn to next connected player
	if r.activePlayerID == clientID {
		r.advanceTurnLocked()
	}
}

// ChooseLevel handles the active player's choice of difficulty.
func (r *Room) ChooseLevel(clientID string, difficulty string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if clientID != r.activePlayerID {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Not your turn: only the active player can choose difficulty level")
		}
		return ErrNotActivePlayer
	}

	if difficulty != "easy" && difficulty != "medium" && difficulty != "hard" {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Invalid difficulty: must be easy, medium, or hard")
		}
		return ErrInvalidDifficulty
	}

	player := r.playerMap[clientID]
	if player != nil {
		player.ChosenDifficulty = difficulty
	}
	return nil
}

// SubmitAnswer handles answer submission from the active player.
// In Phase 3, grading is stubbed as always-succeed.
func (r *Room) SubmitAnswer(clientID string, payload json.RawMessage) ([]RollResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if clientID != r.activePlayerID {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Not your turn: only the active player can submit an answer")
		}
		return nil, ErrNotActivePlayer
	}

	player := r.playerMap[clientID]
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	// Determine number of rolls N matching difficulty (easy=1, medium=2, hard=3)
	numRolls := 1
	switch player.ChosenDifficulty {
	case "easy":
		numRolls = 1
	case "medium":
		numRolls = 2
	case "hard":
		numRolls = 3
	default:
		numRolls = 1
	}

	var rollResults []RollResult

	// Resolve each roll individually
	for i := 1; i <= numRolls; i++ {
		dieRoll := r.diceRng.Intn(6) + 1
		oldPos := player.Position
		rawPos := oldPos + dieRoll

		passedGO := rawPos >= 32
		newPos := rawPos % 32
		player.Position = newPos

		lapBonus := 0
		if passedGO {
			lapBonus = 50
			player.XP += lapBonus
		}

		landedCell := r.board[newPos]
		effectResult := r.ApplyCellEffect(player, landedCell)

		res := RollResult{
			PlayerID:    player.ID,
			RollIndex:   i,
			TotalRolls:  numRolls,
			DieRoll:     dieRoll,
			OldPosition: oldPos,
			NewPosition: newPos,
			PassedGO:    passedGO,
			LapBonus:    lapBonus,
			LandedCell:  landedCell,
			Effect:      effectResult,
			PlayerXP:    player.XP,
		}

		rollResults = append(rollResults, res)

		// Broadcast individual roll resolution
		if r.broadcaster != nil {
			r.broadcaster.BroadcastRoom(r.ID, "roll_resolved", res)
		}
	}

	// Turn fully resolved: broadcast turn_ended
	if r.broadcaster != nil {
		r.broadcaster.BroadcastRoom(r.ID, "turn_ended", TurnEndedPayload{
			PlayerID: clientID,
		})
	}

	player.ChosenDifficulty = ""

	// Advance turn to next connected player in join order
	r.advanceTurnLocked()

	return rollResults, nil
}

// advanceTurnLocked advances active player pointer to next connected player, skipping disconnected slots.
func (r *Room) advanceTurnLocked() {
	totalPlayers := len(r.players)
	if totalPlayers == 0 {
		r.activePlayerID = ""
		return
	}

	// Loop through players in circular order
	for i := 1; i <= totalPlayers; i++ {
		candidateIdx := (r.turnIdx + i) % totalPlayers
		p := r.players[candidateIdx]

		if !p.IsConnected {
			// Skip disconnected player
			continue
		}

		if p.SkipNextTurn {
			// Clear skip flag and skip this player's turn once
			p.SkipNextTurn = false
			continue
		}

		// Eligible active player found
		r.turnIdx = candidateIdx
		r.activePlayerID = p.ID
		r.broadcastTurnStartedLocked()
		return
	}

	// No connected eligible player found
	r.activePlayerID = ""
}

func (r *Room) broadcastTurnStartedLocked() {
	if r.broadcaster != nil && r.activePlayerID != "" {
		r.broadcaster.BroadcastRoom(r.ID, "turn_started", TurnStartedPayload{
			ActivePlayerID: r.activePlayerID,
		})
	}
}

// FormatPlayerTurnSummary produces a human-readable summary of the room state.
func (r *Room) FormatPlayerTurnSummary() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summary := fmt.Sprintf("Room %s | Active Player: %s | Total Players: %d\n", r.ID, r.activePlayerID, len(r.players))
	for i, p := range r.players {
		activeMark := " "
		if p.ID == r.activePlayerID {
			activeMark = "*"
		}
		summary += fmt.Sprintf("[%s] Slot %d: %s (ID: %s, Pos: %d, XP: %d, Connected: %t)\n",
			activeMark, i, p.Name, p.ID, p.Position, p.XP, p.IsConnected)
	}
	return summary
}
