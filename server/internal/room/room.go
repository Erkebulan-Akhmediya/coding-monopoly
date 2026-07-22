package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotActivePlayer    = errors.New("not active player's turn")
	ErrPlayerNotFound     = errors.New("player not found in room")
	ErrInvalidDifficulty  = errors.New("invalid difficulty level, must be 'easy', 'medium', or 'hard'")
	ErrQuestionInProgress = errors.New("a question is already in progress")
	ErrNoQuestion         = errors.New("no question is assigned for this turn")
)

// Broadcaster provides an interface for Room to push messages to connected clients.
type Broadcaster interface {
	BroadcastRoom(roomID string, msgType string, payload any)
	SendError(clientID string, errMsg string)
}

// PrivateBroadcaster is implemented by transports that can send a message to
// one player's connection without exposing it to spectators.
type PrivateBroadcaster interface {
	SendToPlayer(roomID string, clientID string, msgType string, payload any)
}

// ExcludingBroadcaster can broadcast to a room while omitting one connection.
// It keeps the redacted question metadata away from the active player.
type ExcludingBroadcaster interface {
	BroadcastRoomExcept(roomID string, excludedClientID string, msgType string, payload any)
}

// QuestionOption is an option shown to the active player. Correct is kept out
// of JSON so it can never leak through a question_started payload.
type QuestionOption struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"-"`
}

// Question is the server-side representation of an assigned question.
type Question struct {
	ID              string
	Type            string
	Difficulty      string
	Prompt          string
	Options         []QuestionOption
	AcceptedAnswers []string
}

// QuestionProvider assigns one published question for a difficulty.
type QuestionProvider interface {
	AssignQuestion(difficulty string) (Question, error)
}

// QuestionStartedPayload is sent privately with content, and broadcast in
// redacted form to every other connection.
type QuestionStartedPayload struct {
	ProblemID  string           `json:"problem_id,omitempty"`
	Type       string           `json:"type,omitempty"`
	Difficulty string           `json:"difficulty"`
	Deadline   time.Time        `json:"deadline"`
	Prompt     string           `json:"prompt,omitempty"`
	Options    []QuestionOption `json:"options,omitempty"`
}

// AnswerResultPayload is broadcast after grading. CorrectAnswer is populated
// only in the active player's private copy.
type AnswerResultPayload struct {
	PlayerID      string       `json:"player_id"`
	Correct       bool         `json:"correct"`
	TimedOut      bool         `json:"timed_out"`
	Rolls         []RollResult `json:"rolls,omitempty"`
	CorrectAnswer any          `json:"correct_answer,omitempty"`
}

type activeTurn struct {
	question Question
	deadline time.Time
	timer    *time.Timer
	resolved bool // protected by Room.mu; this is the single-resolution guard
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
	ID                string
	mu                sync.RWMutex
	players           []*Player
	playerMap         map[string]*Player
	turnIdx           int
	activePlayerID    string
	board             []BoardCell
	diceRng           *rand.Rand
	broadcaster       Broadcaster
	questionProvider  QuestionProvider
	currentTurn       *activeTurn
	deadlineDurations map[string]time.Duration
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
		deadlineDurations: map[string]time.Duration{
			"easy":   30 * time.Second,
			"medium": 45 * time.Second,
			"hard":   60 * time.Second,
		},
	}
}

// NewRoomWithQuestionProvider creates a room backed by published questions.
func NewRoomWithQuestionProvider(id string, broadcaster Broadcaster, provider QuestionProvider) *Room {
	r := NewRoom(id, broadcaster)
	r.questionProvider = provider
	return r
}

// SetQuestionProvider configures the source used when a level is selected.
// It is intended to be called before the room starts receiving game actions.
func (r *Room) SetQuestionProvider(provider QuestionProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.questionProvider = provider
}

// SetDeadlineDurations overrides the standard deadlines, primarily for fast
// deterministic tests. Production defaults remain 30/45/60 seconds.
func (r *Room) SetDeadlineDurations(easy, medium, hard time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deadlineDurations["easy"] = easy
	r.deadlineDurations["medium"] = medium
	r.deadlineDurations["hard"] = hard
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
		r.cancelCurrentTurnLocked()
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
	if r.currentTurn != nil && !r.currentTurn.resolved {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "A question is already in progress for this turn")
		}
		return ErrQuestionInProgress
	}

	player := r.playerMap[clientID]
	if player == nil {
		return ErrPlayerNotFound
	}
	player.ChosenDifficulty = difficulty

	// The nil-provider path preserves the phase-3 in-memory engine behavior.
	// The production hub always installs a database-backed provider.
	if r.questionProvider == nil {
		return nil
	}

	question, err := r.questionProvider.AssignQuestion(difficulty)
	if err != nil {
		player.ChosenDifficulty = ""
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Unable to assign a question")
		}
		return err
	}
	if question.Type != "mcq" && question.Type != "text" {
		player.ChosenDifficulty = ""
		return errors.New("assigned question has invalid type")
	}

	duration := r.deadlineDurations[difficulty]
	deadline := time.Now().Add(duration)
	turn := &activeTurn{question: question, deadline: deadline}
	r.currentTurn = turn
	turn.timer = time.AfterFunc(duration, func() {
		r.resolveTimeout(turn, clientID)
	})

	// Only the active player receives prompt/options. Everyone else receives
	// only the difficulty and deadline so they can follow the countdown.
	if r.broadcaster != nil {
		redacted := QuestionStartedPayload{Difficulty: difficulty, Deadline: deadline}
		if excluding, ok := r.broadcaster.(ExcludingBroadcaster); ok {
			excluding.BroadcastRoomExcept(r.ID, clientID, "question_started", redacted)
		} else {
			r.broadcaster.BroadcastRoom(r.ID, "question_started", redacted)
		}
		if private, ok := r.broadcaster.(PrivateBroadcaster); ok {
			private.SendToPlayer(r.ID, clientID, "question_started", QuestionStartedPayload{
				ProblemID:  question.ID,
				Type:       question.Type,
				Difficulty: difficulty,
				Deadline:   deadline,
				Prompt:     question.Prompt,
				Options:    question.Options,
			})
		}
	}
	return nil
}

// SubmitAnswer handles answer submission from the active player.
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

	if r.questionProvider != nil {
		if r.currentTurn == nil || r.currentTurn.resolved {
			return nil, ErrNoQuestion
		}
		return r.resolveAnswerLocked(r.currentTurn, clientID, payload, false), nil
	}

	// Legacy phase-3 behavior for rooms without a question provider.
	return r.rollAndEndTurnLocked(player, clientID), nil
}

func (r *Room) resolveTimeout(turn *activeTurn, clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.currentTurn != turn || turn.resolved || r.activePlayerID != clientID {
		return
	}
	r.resolveAnswerLocked(turn, clientID, nil, true)
}

// resolveAnswerLocked is called by either submit_answer or the deadline
// callback while Room.mu is held. Marking resolved before grading makes the
// first caller win the submit-vs-timeout race; the loser becomes a no-op.
func (r *Room) resolveAnswerLocked(turn *activeTurn, clientID string, payload json.RawMessage, timedOut bool) []RollResult {
	if turn.resolved || r.currentTurn != turn || r.activePlayerID != clientID {
		return nil
	}
	turn.resolved = true
	if !timedOut && turn.timer != nil {
		turn.timer.Stop()
	}

	correct := gradeQuestion(turn.question, payload, timedOut)
	player := r.playerMap[clientID]
	var rolls []RollResult
	if correct {
		rolls = r.rollPlayerLocked(player)
	}

	// The public result intentionally contains no answer value.
	if r.broadcaster != nil {
		r.broadcaster.BroadcastRoom(r.ID, "answer_result", AnswerResultPayload{
			PlayerID: clientID,
			Correct:  correct,
			TimedOut: timedOut,
			Rolls:    rolls,
		})
		if private, ok := r.broadcaster.(PrivateBroadcaster); ok {
			private.SendToPlayer(r.ID, clientID, "answer_result", AnswerResultPayload{
				PlayerID:      clientID,
				Correct:       correct,
				TimedOut:      timedOut,
				Rolls:         rolls,
				CorrectAnswer: correctAnswerFor(turn.question),
			})
		}
	}

	r.endTurnLocked(clientID)
	return rolls
}

func gradeQuestion(question Question, payload json.RawMessage, timedOut bool) bool {
	if timedOut {
		return false
	}
	var envelope struct {
		ProblemID string          `json:"problem_id"`
		Answer    json.RawMessage `json:"answer"`
		OptionIDs []string        `json:"option_ids"`
	}
	answer := payload
	if len(payload) > 0 && json.Unmarshal(payload, &envelope) == nil && (envelope.Answer != nil || envelope.OptionIDs != nil) {
		if envelope.OptionIDs != nil {
			answer, _ = json.Marshal(envelope.OptionIDs)
		} else {
			answer = envelope.Answer
			if question.Type == "mcq" {
				var nested struct {
					OptionIDs []string `json:"option_ids"`
				}
				if json.Unmarshal(answer, &nested) == nil && nested.OptionIDs != nil {
					answer, _ = json.Marshal(nested.OptionIDs)
				}
			}
		}
	}

	if question.Type == "mcq" {
		var submitted []string
		if json.Unmarshal(answer, &submitted) != nil {
			var one string
			if json.Unmarshal(answer, &one) != nil {
				return false
			}
			submitted = []string{one}
		}
		correctIDs := make([]string, 0)
		for _, option := range question.Options {
			if option.Correct {
				correctIDs = append(correctIDs, option.ID)
			}
		}
		return sameStringSet(submitted, correctIDs)
	}

	var submitted string
	if json.Unmarshal(answer, &submitted) != nil {
		return false
	}
	for _, accepted := range question.AcceptedAnswers {
		if strings.EqualFold(strings.TrimSpace(submitted), strings.TrimSpace(accepted)) {
			return true
		}
	}
	return false
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int, len(a))
	for _, value := range a {
		counts[value]++
	}
	for _, value := range b {
		counts[value]--
		if counts[value] < 0 {
			return false
		}
	}
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}

func correctAnswerFor(question Question) any {
	if question.Type == "mcq" {
		ids := make([]string, 0)
		for _, option := range question.Options {
			if option.Correct {
				ids = append(ids, option.ID)
			}
		}
		return ids
	}
	return question.AcceptedAnswers
}

func (r *Room) rollPlayerLocked(player *Player) []RollResult {

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

	return rollResults
}

func (r *Room) rollAndEndTurnLocked(player *Player, clientID string) []RollResult {
	rolls := r.rollPlayerLocked(player)
	r.endTurnLocked(clientID)
	return rolls
}

func (r *Room) endTurnLocked(clientID string) {
	if r.broadcaster != nil {
		r.broadcaster.BroadcastRoom(r.ID, "turn_ended", TurnEndedPayload{PlayerID: clientID})
	}
	if player := r.playerMap[clientID]; player != nil {
		player.ChosenDifficulty = ""
	}
	r.currentTurn = nil
	r.advanceTurnLocked()
}

func (r *Room) cancelCurrentTurnLocked() {
	if r.currentTurn != nil && r.currentTurn.timer != nil {
		r.currentTurn.timer.Stop()
	}
	r.currentTurn = nil
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
