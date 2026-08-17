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
	ErrGamePaused         = errors.New("game is paused")
	ErrGameInProgress     = errors.New("game is already in progress; new players cannot join")
	ErrNameTaken          = errors.New("a player with this name already exists in this room")
	ErrGameNotStarted     = errors.New("game has not started yet")
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
	question   Question
	difficulty string
	deadline   time.Time
	timer      *time.Timer
	resolved   bool // protected by Room.mu; this is the single-resolution guard
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

// PlayerStanding is one row in the end-of-game summary.
type PlayerStanding struct {
	PlayerID    string `json:"player_id"`
	Name        string `json:"name"`
	XP          int    `json:"xp"`
	Position    int    `json:"position"`
	Rank        int    `json:"rank"`
	IsConnected bool   `json:"is_connected"`
}

// GameOverPayload is broadcast when the match ends (target XP or admin).
type GameOverPayload struct {
	WinnerID   string           `json:"winner_id"`
	WinnerName string           `json:"winner_name"`
	Reason     string           `json:"reason"` // "target_xp" | "admin"
	TargetXP   int              `json:"target_xp"`
	Standings  []PlayerStanding `json:"standings"`
}

const (
	// DefaultTargetXP is the first-to-XP win threshold.
	DefaultTargetXP = 500
	// DefaultDisconnectGrace is how long an active player may blip offline
	// before their turn is forcibly ended. Distinct from the question timer.
	DefaultDisconnectGrace = 5 * time.Second
)

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
	started           bool // set by AdminStartGame
	paused            bool // set by AdminTogglePause
	finished          bool
	gameOver          *GameOverPayload
	targetXP          int
	disconnectGrace   time.Duration
	graceTimers       map[string]*time.Timer
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
		targetXP:        DefaultTargetXP,
		disconnectGrace: DefaultDisconnectGrace,
		graceTimers:     make(map[string]*time.Timer),
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

// SetDisconnectGrace overrides the active-player disconnect grace period (tests).
func (r *Room) SetDisconnectGrace(d time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.disconnectGrace = d
}

// SetTargetXP overrides the first-to-XP win threshold (tests / config).
func (r *Room) SetTargetXP(xp int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.targetXP = xp
}

// IsStarted reports whether the game has been started by an admin.
func (r *Room) IsStarted() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.started
}

// IsFinished reports whether the game has ended.
func (r *Room) IsFinished() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.finished
}

// GetGameOver returns the end-of-game payload when the match is finished.
func (r *Room) GetGameOver() *GameOverPayload {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.gameOver == nil {
		return nil
	}
	cp := *r.gameOver
	cp.Standings = append([]PlayerStanding(nil), r.gameOver.Standings...)
	return &cp
}

// GetTargetXP returns the configured first-to-XP threshold.
func (r *Room) GetTargetXP() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.targetXP
}

// CanReclaimPlayer reports whether a reconnecting client may take over playerID.
// Allowed when the slot exists and is currently disconnected (or in grace).
func (r *Room) CanReclaimPlayer(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.playerMap[playerID]
	if !exists {
		return false
	}
	return !p.IsConnected
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

// GetTurnState returns the current active player ID, whether a question is in progress, and the deadline if active.
func (r *Room) GetTurnState() (activePlayerID string, questionActive bool, deadline *time.Time) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	activePlayerID = r.activePlayerID
	if r.currentTurn != nil && !r.currentTurn.resolved {
		questionActive = true
		d := r.currentTurn.deadline
		deadline = &d
	}
	return activePlayerID, questionActive, deadline
}

// GetActiveQuestionPayload returns the full question payload if the given client is the active player.
func (r *Room) GetActiveQuestionPayload(clientID string) *QuestionStartedPayload {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activePlayerID == clientID && r.currentTurn != nil && !r.currentTurn.resolved {
		return &QuestionStartedPayload{
			ProblemID:  r.currentTurn.question.ID,
			Type:       r.currentTurn.question.Type,
			Difficulty: r.currentTurn.difficulty,
			Deadline:   r.currentTurn.deadline,
			Prompt:     r.currentTurn.question.Prompt,
			Options:    r.currentTurn.question.Options,
		}
	}
	return nil
}

// AddOrReconnectPlayer handles player join or reconnect, maintaining strict join order.
// A reconnecting player keeps the same slot, position, XP, and modifiers.
// New players can only join before the game is started by an admin, and player names must be unique within the room.
func (r *Room) AddOrReconnectPlayer(clientID string, name string) (*Player, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if p, exists := r.playerMap[clientID]; exists {
		r.clearGraceTimerLocked(clientID)
		p.IsConnected = true
		if name != "" {
			p.Name = name
		}

		// If the game is started and had no eligible active player, resume turn order
		if r.started && r.activePlayerID == "" && !r.finished {
			if r.turnIdx >= 0 && r.turnIdx < len(r.players) && r.players[r.turnIdx].ID == clientID {
				r.activePlayerID = p.ID
				r.broadcastTurnStartedLocked()
			} else {
				r.advanceTurnLocked()
			}
		}

		return p, nil
	}

	trimmedName := strings.TrimSpace(name)

	// Check name uniqueness among existing players in the room (case-insensitive)
	for _, p := range r.players {
		if strings.EqualFold(strings.TrimSpace(p.Name), trimmedName) {
			return nil, ErrNameTaken
		}
	}

	if r.finished {
		return nil, errors.New("game is over")
	}

	// Mid-game join restriction: only reconnecting players can join a started game
	if r.started {
		return nil, ErrGameInProgress
	}

	// New player joining before game starts
	player := NewPlayer(clientID, name)
	r.players = append(r.players, player)
	r.playerMap[clientID] = player

	return player, nil
}

// DisconnectPlayer marks a player as disconnected without removing their slot.
// If they are the active player, a short grace period starts before the turn is
// forfeited — a brief network blip does not immediately end their turn.
func (r *Room) DisconnectPlayer(clientID string) {
	r.disconnectPlayer(clientID, false)
}

// DisconnectPlayerImmediate marks a player disconnected and forfeits their turn
// immediately (used for admin kick).
func (r *Room) DisconnectPlayerImmediate(clientID string) {
	r.disconnectPlayer(clientID, true)
}

func (r *Room) disconnectPlayer(clientID string, immediate bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.playerMap[clientID]
	if !exists {
		return
	}

	p.IsConnected = false
	r.clearGraceTimerLocked(clientID)

	if r.activePlayerID != clientID || r.finished {
		return
	}

	if immediate || r.disconnectGrace <= 0 {
		r.forfeitActiveTurnLocked(clientID)
		return
	}

	grace := r.disconnectGrace
	r.graceTimers[clientID] = time.AfterFunc(grace, func() {
		r.onDisconnectGraceExpired(clientID)
	})
}

func (r *Room) onDisconnectGraceExpired(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.graceTimers, clientID)

	p, exists := r.playerMap[clientID]
	if !exists || p.IsConnected {
		return
	}
	if r.activePlayerID != clientID || r.finished {
		return
	}
	r.forfeitActiveTurnLocked(clientID)
}

// forfeitActiveTurnLocked ends the active player's turn without grading (grace
// expiry / kick). The question timer is cancelled; mid-question content is dropped.
func (r *Room) forfeitActiveTurnLocked(clientID string) {
	r.cancelCurrentTurnLocked()
	if player := r.playerMap[clientID]; player != nil {
		player.ChosenDifficulty = ""
	}
	if r.broadcaster != nil {
		r.broadcaster.BroadcastRoom(r.ID, "answer_result", AnswerResultPayload{
			PlayerID: clientID,
			Correct:  false,
			TimedOut: true,
		})
		r.broadcaster.BroadcastRoom(r.ID, "turn_ended", TurnEndedPayload{PlayerID: clientID})
	}
	r.advanceTurnLocked()
}

func (r *Room) clearGraceTimerLocked(clientID string) {
	if t, ok := r.graceTimers[clientID]; ok {
		t.Stop()
		delete(r.graceTimers, clientID)
	}
}

// ChooseLevel handles the active player's choice of difficulty.
func (r *Room) ChooseLevel(clientID string, difficulty string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.finished {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game is over")
		}
		return errors.New("game is over")
	}

	if !r.started {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game has not started yet")
		}
		return ErrGameNotStarted
	}

	if r.paused {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game is paused")
		}
		return ErrGamePaused
	}

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
	turn := &activeTurn{question: question, difficulty: difficulty, deadline: deadline}
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

	if r.finished {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game is over")
		}
		return nil, errors.New("game is over")
	}

	if !r.started {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game has not started yet")
		}
		return nil, ErrGameNotStarted
	}

	if r.paused {
		if r.broadcaster != nil {
			r.broadcaster.SendError(clientID, "Game is paused")
		}
		return nil, ErrGamePaused
	}

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
	r.checkGameOverLocked()
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
	r.checkGameOverLocked()
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
	if !r.finished {
		r.advanceTurnLocked()
	}
}

func (r *Room) cancelCurrentTurnLocked() {
	if r.currentTurn != nil && r.currentTurn.timer != nil {
		r.currentTurn.timer.Stop()
	}
	r.currentTurn = nil
}

// advanceTurnLocked advances active player pointer to next connected player, skipping disconnected slots.
func (r *Room) advanceTurnLocked() {
	if r.finished {
		r.activePlayerID = ""
		return
	}

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
	if r.broadcaster != nil && r.activePlayerID != "" && !r.finished {
		r.broadcaster.BroadcastRoom(r.ID, "turn_started", TurnStartedPayload{
			ActivePlayerID: r.activePlayerID,
		})
	}
}

func (r *Room) checkGameOverLocked() {
	if r.finished {
		return
	}
	var best *Player
	for _, p := range r.players {
		if p.XP >= r.targetXP {
			if best == nil || p.XP > best.XP {
				best = p
			}
		}
	}
	if best == nil {
		return
	}
	r.finishGameLocked(best.ID, "target_xp")
}

func (r *Room) finishGameLocked(winnerID string, reason string) {
	if r.finished {
		return
	}
	r.finished = true
	r.cancelCurrentTurnLocked()
	r.clearAllGraceTimersLocked()
	r.activePlayerID = ""

	standings := r.buildStandingsLocked()
	winnerName := winnerID
	if p := r.playerMap[winnerID]; p != nil {
		winnerName = p.Name
	}
	payload := &GameOverPayload{
		WinnerID:   winnerID,
		WinnerName: winnerName,
		Reason:     reason,
		TargetXP:   r.targetXP,
		Standings:  standings,
	}
	r.gameOver = payload

	if r.broadcaster != nil {
		r.broadcaster.BroadcastRoom(r.ID, "game_over", payload)
	}
}

func (r *Room) clearAllGraceTimersLocked() {
	for id, t := range r.graceTimers {
		t.Stop()
		delete(r.graceTimers, id)
	}
}

func (r *Room) buildStandingsLocked() []PlayerStanding {
	sorted := make([]*Player, len(r.players))
	copy(sorted, r.players)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].XP > sorted[i].XP {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	out := make([]PlayerStanding, 0, len(sorted))
	for i, p := range sorted {
		out = append(out, PlayerStanding{
			PlayerID:    p.ID,
			Name:        p.Name,
			XP:          p.XP,
			Position:    p.Position,
			Rank:        i + 1,
			IsConnected: p.IsConnected,
		})
	}
	return out
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

func (r *Room) Board() []BoardCell {
	return r.board
}

// ---------------------------------------------------------------------------
// Admin control methods
// ---------------------------------------------------------------------------

// AdminStartGame starts the game in the room manually.
// Can only be started if the game has not already started or finished, and at least 1 connected player exists.
func (r *Room) AdminStartGame() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started || r.finished {
		return false
	}

	for i, p := range r.players {
		if p.IsConnected {
			r.turnIdx = i
			r.activePlayerID = p.ID
			r.started = true
			r.broadcastTurnStartedLocked()
			return true
		}
	}
	return false
}

// AdminTogglePause pauses or unpauses the game. Returns true if the game is
// now paused, false if it was just resumed.
// Pause: the current turn timer is cancelled; the active turn is preserved.
// Resume: the turn timer is restarted with the remaining deadline.
func (r *Room) AdminTogglePause() (paused bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.paused = !r.paused
	if r.paused {
		// Stop the running timer without marking the turn as resolved.
		if r.currentTurn != nil && r.currentTurn.timer != nil {
			r.currentTurn.timer.Stop()
		}
	} else {
		// Resume: re-arm the timer for the remaining time (or 1 ms if already past).
		if r.currentTurn != nil && !r.currentTurn.resolved {
			remaining := time.Until(r.currentTurn.deadline)
			if remaining <= 0 {
				remaining = time.Millisecond
			}
			turn := r.currentTurn
			activeID := r.activePlayerID
			r.currentTurn.timer = time.AfterFunc(remaining, func() {
				r.resolveTimeout(turn, activeID)
			})
		}
	}
	return r.paused
}

// IsPaused reports whether the game is currently paused.
func (r *Room) IsPaused() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.paused
}

// AdminKickPlayer disconnects and removes the named player from the room engine immediately.
// Returns the player's name for logging/event feed display.
func (r *Room) AdminKickPlayer(playerID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.playerMap[playerID]
	if !exists {
		return playerID
	}
	name := p.Name

	r.clearGraceTimerLocked(playerID)
	p.IsConnected = false

	// Remove from players slice and map
	newPlayers := make([]*Player, 0, len(r.players)-1)
	for _, pl := range r.players {
		if pl.ID != playerID {
			newPlayers = append(newPlayers, pl)
		}
	}
	r.players = newPlayers
	delete(r.playerMap, playerID)

	// If this was the active player, forfeit active turn and advance
	if r.activePlayerID == playerID && !r.finished {
		r.forfeitActiveTurnLocked(playerID)
	}

	return name
}

// AdminEndGame forcibly finishes the match and broadcasts standings.
// If winnerID is empty, the current highest-XP player wins.
func (r *Room) AdminEndGame(winnerID string) *GameOverPayload {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.finished {
		return r.gameOver
	}

	chosen := winnerID
	if chosen == "" || r.playerMap[chosen] == nil {
		var best *Player
		for _, p := range r.players {
			if best == nil || p.XP > best.XP {
				best = p
			}
		}
		if best == nil {
			r.finishGameLocked("", "admin")
			return r.gameOver
		}
		chosen = best.ID
	}
	r.finishGameLocked(chosen, "admin")
	return r.gameOver
}

// AdminSkipTurn forcibly ends the current active player's turn without
// grading — useful when the player has not yet chosen a difficulty and the
// normal answer-phase timer therefore never started.
// If playerID is empty it defaults to the current active player.
// Returns the name of the player whose turn was skipped.
func (r *Room) AdminSkipTurn(playerID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.finished {
		return ""
	}

	target := playerID
	if target == "" {
		target = r.activePlayerID
	}
	if target == "" || target != r.activePlayerID {
		return target
	}

	p := r.playerMap[target]
	name := target
	if p != nil {
		name = p.Name
	}

	r.forfeitActiveTurnLocked(target)
	return name
}
