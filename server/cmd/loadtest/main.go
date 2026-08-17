// Command loadtest measures one full turn rotation for a realistic class size.
//
// Turns are sequential, so the useful metric is wall-clock time for every
// player to take one turn at realistic answer speeds — not concurrent RPS.
//
// Usage:
//
//	go run ./cmd/loadtest
//	go run ./cmd/loadtest -players 24 -correct-rate 0.8
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"server/internal/room"
	"server/internal/ws"
)

type thinkProfile struct {
	difficulty string
	think      time.Duration
	weight     int // relative weight for random pick
}

// Realistic (not worst-case) think times: well under the server deadlines
// of 30s/45s/60s. Mix skewed toward easy, as students usually prefer it.
var profiles = []thinkProfile{
	{difficulty: "easy", think: 12 * time.Second, weight: 60},
	{difficulty: "medium", think: 20 * time.Second, weight: 30},
	{difficulty: "hard", think: 30 * time.Second, weight: 10},
}

type fixedProvider struct {
	mu sync.Mutex
	n  int
}

func (p *fixedProvider) AssignQuestion(difficulty string) (room.Question, error) {
	p.mu.Lock()
	p.n++
	id := fmt.Sprintf("load-q-%d", p.n)
	p.mu.Unlock()

	return room.Question{
		ID:         id,
		Type:       "mcq",
		Difficulty: difficulty,
		Prompt:     "load-test question",
		Options: []room.QuestionOption{
			{ID: "opt-correct", Text: "right", Correct: true},
			{ID: "opt-wrong", Text: "wrong", Correct: false},
		},
	}, nil
}

type bot struct {
	name   string
	conn   *websocket.Conn
	inbox  chan ws.Message
	done   chan struct{}
	player string // server-assigned player id once known from state_sync
}

func main() {
	players := flag.Int("players", 24, "class size (players in one room)")
	correctRate := flag.Float64("correct-rate", 0.80, "fraction of answers that are correct")
	seed := flag.Int64("seed", 42, "RNG seed for reproducible difficulty mix")
	flag.Parse()

	if *players < 2 {
		log.Fatal("need at least 2 players")
	}
	if *correctRate < 0 || *correctRate > 1 {
		log.Fatal("correct-rate must be in [0,1]")
	}

	rng := rand.New(rand.NewSource(*seed))

	hub := ws.NewHub(&fixedProvider{})
	go hub.Run()
	defer hub.Stop()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	roomID := "loadtest-class"
	if err := hub.CreateRoom(roomID); err != nil {
		log.Fatalf("CreateRoom: %v", err)
	}

	bots := make([]*bot, *players)
	for i := 0; i < *players; i++ {
		name := fmt.Sprintf("Student%02d", i+1)
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			log.Fatalf("dial %s: %v", name, err)
		}
		b := &bot{
			name:  name,
			conn:  conn,
			inbox: make(chan ws.Message, 256),
			done:  make(chan struct{}),
		}
		go b.readLoop()
		bots[i] = b
		if err := b.sendJoin(roomID); err != nil {
			log.Fatalf("join %s: %v", name, err)
		}
		// Wait until this bot's own player id appears in a state_sync.
		b.drainUntil(5*time.Second, func(m ws.Message) bool {
			if m.Type != ws.MessageTypeStateSync {
				return false
			}
			var payload ws.StateSyncPayload
			if json.Unmarshal(m.Payload, &payload) != nil {
				return false
			}
			for _, p := range payload.Players {
				if p.Name == b.name {
					b.player = p.ID
					return true
				}
			}
			return false
		})
		if b.player == "" {
			log.Fatalf("could not resolve player id for %s after join", name)
		}
	}

	byID := make(map[string]*bot, len(bots))
	for _, b := range bots {
		byID[b.player] = b
	}

	fmt.Fprintf(os.Stderr, "loadtest: %d players joined room %q\n", *players, roomID)
	fmt.Fprintf(os.Stderr, "profile: easy~12s(60%%) medium~20s(30%%) hard~30s(10%%), correct-rate=%.0f%%\n", *correctRate*100)

	// Discover who currently has the turn (first joiner).
	current := bots[0].player
	start := time.Now()
	var (
		totalThink   time.Duration
		correctCount int
		easyN        int
		mediumN      int
		hardN        int
	)

	for turn := 0; turn < *players; turn++ {
		active := byID[current]
		if active == nil {
			log.Fatalf("turn %d: unknown active player %s", turn, current)
		}

		prof := pickProfile(rng)
		switch prof.difficulty {
		case "easy":
			easyN++
		case "medium":
			mediumN++
		case "hard":
			hardN++
		}
		totalThink += prof.think

		choosePayload, _ := json.Marshal(ws.ChooseLevelPayload{Difficulty: prof.difficulty})
		if err := active.conn.WriteJSON(ws.Message{
			Type:    ws.MessageTypeChooseLevel,
			RoomID:  roomID,
			Payload: choosePayload,
		}); err != nil {
			log.Fatalf("%s choose_level: %v", active.name, err)
		}

		// Wait for private question_started on active bot.
		qMsg := active.waitType(ws.MessageTypeQuestionStarted, 5*time.Second)
		var q room.QuestionStartedPayload
		if err := json.Unmarshal(qMsg.Payload, &q); err != nil {
			log.Fatalf("decode question: %v", err)
		}

		// Realistic think time (the metric we care about).
		time.Sleep(prof.think)

		correct := rng.Float64() < *correctRate
		answerIDs := []string{"opt-wrong"}
		if correct {
			answerIDs = []string{"opt-correct"}
			correctCount++
		}
		answerJSON, _ := json.Marshal(answerIDs)
		submitPayload, _ := json.Marshal(map[string]any{
			"problem_id": q.ProblemID,
			"answer":     json.RawMessage(answerJSON),
		})
		if err := active.conn.WriteJSON(ws.Message{
			Type:    ws.MessageTypeSubmitAnswer,
			RoomID:  roomID,
			Payload: submitPayload,
		}); err != nil {
			log.Fatalf("%s submit: %v", active.name, err)
		}

		// Drain until the next player's turn_started (or turn_ended then turn_started).
		nextID := ""
		deadline := time.Now().Add(10 * time.Second)
		for nextID == "" && time.Now().Before(deadline) {
			msg, ok := active.recv(time.Until(deadline))
			if !ok {
				break
			}
			if msg.Type == ws.MessageTypeTurnStarted {
				var payload struct {
					ActivePlayerID string `json:"active_player_id"`
				}
				_ = json.Unmarshal(msg.Payload, &payload)
				if payload.ActivePlayerID != "" && payload.ActivePlayerID != current {
					nextID = payload.ActivePlayerID
				}
			}
		}
		if nextID == "" {
			log.Fatalf("turn %d (%s): did not observe next turn_started", turn, active.name)
		}

		// Flush other bots' inboxes so they don't back up.
		for _, b := range bots {
			if b == active {
				continue
			}
			b.drainUntil(50*time.Millisecond, func(ws.Message) bool { return false })
		}

		current = nextID
		fmt.Fprintf(os.Stderr, "  turn %2d/%d %-10s difficulty=%-6s correct=%v think=%s → next\n",
			turn+1, *players, active.name, prof.difficulty, correct, prof.think)
	}

	elapsed := time.Since(start)
	avgThink := totalThink / time.Duration(*players)
	overhead := elapsed - totalThink

	fmt.Println()
	fmt.Println("=== Full-rotation load test results ===")
	fmt.Printf("players:            %d\n", *players)
	fmt.Printf("difficulty mix:     easy=%d medium=%d hard=%d\n", easyN, mediumN, hardN)
	fmt.Printf("correct answers:    %d / %d (%.0f%%)\n", correctCount, *players, 100*float64(correctCount)/float64(*players))
	fmt.Printf("sum of think times: %s\n", totalThink.Round(time.Millisecond))
	fmt.Printf("avg think / turn:   %s\n", avgThink.Round(time.Millisecond))
	fmt.Printf("server overhead:    %s (WS + grading + dice + turn advance)\n", overhead.Round(time.Millisecond))
	fmt.Printf("FULL ROTATION:      %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("per-player average: %s\n", (elapsed / time.Duration(*players)).Round(time.Millisecond))
	fmt.Println()
	fmt.Println("Sanity check vs class period:")
	fmt.Printf("  50-minute period → ~%.1f full rotations\n", (50 * time.Minute).Seconds()/elapsed.Seconds())
	fmt.Printf("  90-minute period → ~%.1f full rotations\n", (90 * time.Minute).Seconds()/elapsed.Seconds())
	if elapsed > 15*time.Minute {
		fmt.Println("  WARN: one rotation exceeds 15 minutes — consider fewer players or shorter think times for a single period.")
	} else {
		fmt.Println("  OK: one rotation fits comfortably in a normal class period with time left for multiple laps / discussion.")
	}

	for _, b := range bots {
		close(b.done)
		_ = b.conn.Close()
	}
}

func pickProfile(rng *rand.Rand) thinkProfile {
	total := 0
	for _, p := range profiles {
		total += p.weight
	}
	n := rng.Intn(total)
	for _, p := range profiles {
		if n < p.weight {
			return p
		}
		n -= p.weight
	}
	return profiles[0]
}

func (b *bot) readLoop() {
	defer close(b.inbox)
	for {
		select {
		case <-b.done:
			return
		default:
		}
		_ = b.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		var msg ws.Message
		if err := b.conn.ReadJSON(&msg); err != nil {
			return
		}
		select {
		case b.inbox <- msg:
		case <-b.done:
			return
		}
	}
}

func (b *bot) sendJoin(roomID string) error {
	payload, _ := json.Marshal(ws.JoinPayload{Name: b.name, RoomID: roomID})
	return b.conn.WriteJSON(ws.Message{
		Type:    ws.MessageTypeJoin,
		RoomID:  roomID,
		Payload: payload,
	})
}

func (b *bot) recv(timeout time.Duration) (ws.Message, bool) {
	select {
	case msg, ok := <-b.inbox:
		return msg, ok
	case <-time.After(timeout):
		return ws.Message{}, false
	}
}

func (b *bot) waitType(want string, timeout time.Duration) ws.Message {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msg, ok := b.recv(time.Until(deadline))
		if !ok {
			break
		}
		if msg.Type == want {
			return msg
		}
	}
	log.Fatalf("%s: timed out waiting for %s", b.name, want)
	return ws.Message{}
}

func (b *bot) drainUntil(timeout time.Duration, stop func(ws.Message) bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msg, ok := b.recv(time.Until(deadline))
		if !ok {
			return
		}
		if stop(msg) {
			return
		}
	}
}
