package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type BoardCellData struct {
	Type   string
	Params map[string]interface{}
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/monopoly?sslmode=disable"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	log.Println("Starting database seeding...")

	// Truncate tables to ensure idempotency
	_, err = conn.Exec(ctx, "TRUNCATE games, board_cells, problems, problem_options, problem_accepted_answers, submissions, game_events CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v\n", err)
	}
	log.Println("Database truncated successfully.")

	// 1. Seed a default game
	defaultGameID := "00000000-0000-0000-0000-000000000000"
	_, err = conn.Exec(ctx, "INSERT INTO games (id, status) VALUES ($1, $2)", defaultGameID, "lobby")
	if err != nil {
		log.Fatalf("Failed to insert default game: %v\n", err)
	}
	log.Println("Seeded default game.")

	// 2. Define the 32 board cells
	// Index 0: Deploy, Index 8: Code Freeze, Index 16: Coffee Break, Index 24: Deadline.
	// The rest are effects.
	boardCells := make(map[int]BoardCellData)
	boardCells[0] = BoardCellData{Type: "deploy", Params: map[string]interface{}{"name": "Deploy", "bonus": 200}}
	boardCells[8] = BoardCellData{Type: "code_freeze", Params: map[string]interface{}{"name": "Code Freeze"}}
	boardCells[16] = BoardCellData{Type: "coffee_break", Params: map[string]interface{}{"name": "Coffee Break"}}
	boardCells[24] = BoardCellData{Type: "deadline", Params: map[string]interface{}{"name": "Deadline"}}

	// Distribute the other 28 cells carrying effects:
	// XP gain (S/M/L), XP loss (S/M), mystery/random event, teleport, skip-next-turn, double-XP, free pass, special bonus challenge.
	effects := []BoardCellData{
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (S)", "size": "S", "amount": 50}},
		{Type: "mystery", Params: map[string]interface{}{"name": "Mystery Event"}},
		{Type: "xp_loss", Params: map[string]interface{}{"name": "XP Loss (S)", "size": "S", "amount": 30}},
		{Type: "double_xp", Params: map[string]interface{}{"name": "Double XP"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (M)", "size": "M", "amount": 100}},
		{Type: "skip_turn", Params: map[string]interface{}{"name": "Skip Next Turn"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (S)", "size": "S", "amount": 50}},
		// index 8 is code_freeze
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (M)", "size": "M", "amount": 100}},
		{Type: "teleport", Params: map[string]interface{}{"name": "Teleport", "target": 18}},
		{Type: "xp_loss", Params: map[string]interface{}{"name": "XP Loss (M)", "size": "M", "amount": 60}},
		{Type: "free_pass", Params: map[string]interface{}{"name": "Free Pass"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (L)", "size": "L", "amount": 200}},
		{Type: "mystery", Params: map[string]interface{}{"name": "Mystery Event"}},
		{Type: "bonus_challenge", Params: map[string]interface{}{"name": "Special Bonus Challenge"}},
		// index 16 is coffee_break
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (S)", "size": "S", "amount": 50}},
		{Type: "mystery", Params: map[string]interface{}{"name": "Mystery Event"}},
		{Type: "xp_loss", Params: map[string]interface{}{"name": "XP Loss (S)", "size": "S", "amount": 30}},
		{Type: "double_xp", Params: map[string]interface{}{"name": "Double XP"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (M)", "size": "M", "amount": 100}},
		{Type: "skip_turn", Params: map[string]interface{}{"name": "Skip Next Turn"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (S)", "size": "S", "amount": 50}},
		// index 24 is deadline
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (M)", "size": "M", "amount": 100}},
		{Type: "teleport", Params: map[string]interface{}{"name": "Teleport", "target": 2}},
		{Type: "xp_loss", Params: map[string]interface{}{"name": "XP Loss (M)", "size": "M", "amount": 60}},
		{Type: "free_pass", Params: map[string]interface{}{"name": "Free Pass"}},
		{Type: "xp_gain", Params: map[string]interface{}{"name": "XP Gain (L)", "size": "L", "amount": 200}},
		{Type: "mystery", Params: map[string]interface{}{"name": "Mystery Event"}},
		{Type: "bonus_challenge", Params: map[string]interface{}{"name": "Special Bonus Challenge"}},
	}

	effectIdx := 0
	for i := 0; i < 32; i++ {
		if i == 0 || i == 8 || i == 16 || i == 24 {
			continue
		}
		boardCells[i] = effects[effectIdx]
		effectIdx++
	}

	// Insert board cells
	for i := 0; i < 32; i++ {
		cell := boardCells[i]
		paramsJSON, err := json.Marshal(cell.Params)
		if err != nil {
			log.Fatalf("Failed to marshal cell params: %v\n", err)
		}
		_, err = conn.Exec(ctx, "INSERT INTO board_cells (game_id, cell_index, type, params) VALUES ($1, $2, $3, $4)", defaultGameID, i, cell.Type, paramsJSON)
		if err != nil {
			log.Fatalf("Failed to insert board cell %d: %v\n", i, err)
		}
	}
	log.Printf("Seeded 32 board cells for game %s.\n", defaultGameID)

	// 3. Seed problems (18 minimum, 3 per difficulty per type)
	type SeedOption struct {
		Text      string
		IsCorrect bool
	}
	type SeedProblem struct {
		Type            string
		Difficulty      string
		Title           string
		Prompt          string
		Options         []SeedOption
		AcceptedAnswers []string
	}

	problemsToSeed := []SeedProblem{
		// EASY MCQ
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      "Go Variable Declaration",
			Prompt:     "Which of the following is the correct short variable declaration in Go?",
			Options: []SeedOption{
				{Text: "x := 10", IsCorrect: true},
				{Text: "x = 10", IsCorrect: false},
				{Text: "var x = 10", IsCorrect: false},
				{Text: "let x = 10", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      "Slices vs Arrays in Go",
			Prompt:     "What is the key difference in declaration between a slice and an array in Go?",
			Options: []SeedOption{
				{Text: "An array has a fixed size specified in its type, while a slice has a dynamic size.", IsCorrect: true},
				{Text: "A slice is declared with parentheses `()` and an array with brackets `[]`.", IsCorrect: false},
				{Text: "Slices are value types, while arrays are reference types.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      "Go Package Main",
			Prompt:     "Which package must be declared for a Go file to contain the entry point of an executable program?",
			Options: []SeedOption{
				{Text: "main", IsCorrect: true},
				{Text: "root", IsCorrect: false},
				{Text: "start", IsCorrect: false},
				{Text: "global", IsCorrect: false},
			},
		},

		// MEDIUM MCQ
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      "Go Channel Directional Constraint",
			Prompt:     "How do you define a write-only channel of type int in a Go function signature?",
			Options: []SeedOption{
				{Text: "chan<- int", IsCorrect: true},
				{Text: "<-chan int", IsCorrect: false},
				{Text: "chan int", IsCorrect: false},
				{Text: "write chan int", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      "Defer Statement Execution Order",
			Prompt:     "In what order are deferred function calls executed in Go?",
			Options: []SeedOption{
				{Text: "LIFO (Last-In, First-Out)", IsCorrect: true},
				{Text: "FIFO (First-In, First-Out)", IsCorrect: false},
				{Text: "Random order", IsCorrect: false},
				{Text: "Concurrent order", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      "Interface Implementation in Go",
			Prompt:     "How does a custom type in Go implement an interface?",
			Options: []SeedOption{
				{Text: "Implicitly, by implementing all the methods declared in the interface.", IsCorrect: true},
				{Text: "Explicitly, using the `implements` keyword.", IsCorrect: false},
				{Text: "By inheriting from a base interface struct.", IsCorrect: false},
				{Text: "By registering the type with the interface package.", IsCorrect: false},
			},
		},

		// HARD MCQ
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      "Go Scheduler M:N Model",
			Prompt:     "In Go's concurrency scheduler (G-M-P model), what does the letter 'P' represent?",
			Options: []SeedOption{
				{Text: "Processor: a logical resource representing a context required to execute Go code.", IsCorrect: true},
				{Text: "Platform: the operating system thread abstraction.", IsCorrect: false},
				{Text: "Program: the main package memory boundary.", IsCorrect: false},
				{Text: "Process: a standard Unix system process.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      "Sync Map Use Case",
			Prompt:     "When is it recommended to use sync.Map instead of a regular map with a sync.Mutex or sync.RWMutex?",
			Options: []SeedOption{
				{Text: "When the entry set is stable and the write rate is very low compared to read rate, or when multiple goroutines read, write, and overwrite entries for disjoint sets of keys.", IsCorrect: true},
				{Text: "Whenever multiple goroutines access a map concurrently, regardless of the read/write ratio.", IsCorrect: false},
				{Text: "When you need to sort the keys of the map concurrently.", IsCorrect: false},
				{Text: "When the map stores a small number of keys (less than 100) that change frequently.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      "Goroutine Stack Allocation",
			Prompt:     "What is the initial stack size allocated for a goroutine in Go?",
			Options: []SeedOption{
				{Text: "2 KB", IsCorrect: true},
				{Text: "8 KB", IsCorrect: false},
				{Text: "1 MB", IsCorrect: false},
				{Text: "4 KB", IsCorrect: false},
			},
		},

		// EASY TEXT
		{
			Type:            "text",
			Difficulty:      "easy",
			Title:           "Map Lookup Zero Value",
			Prompt:          "What is the second, optional boolean value returned when retrieving an element from a Go map (e.g., val, ok := m[key])?",
			AcceptedAnswers: []string{"ok", "exists", "found", "present"},
		},
		{
			Type:            "text",
			Difficulty:      "easy",
			Title:           "Go Garbage Collection",
			Prompt:          "What style of garbage collector does Go use? (Hint: uses three colors)",
			AcceptedAnswers: []string{"tri-color mark and sweep", "tricolor mark and sweep", "tri-color", "tricolor"},
		},
		{
			Type:            "text",
			Difficulty:      "easy",
			Title:           "Pointer Zero Value",
			Prompt:          "What is the zero value of a pointer, slice, channel, or map in Go?",
			AcceptedAnswers: []string{"nil"},
		},

		// MEDIUM TEXT
		{
			Type:            "text",
			Difficulty:      "medium",
			Title:           "Context Cancellation Propagation",
			Prompt:          "Which standard Go library function is used to create a context that can be cancelled manually?",
			AcceptedAnswers: []string{"context.WithCancel", "WithCancel"},
		},
		{
			Type:            "text",
			Difficulty:      "medium",
			Title:           "String Length in Bytes",
			Prompt:          "What built-in Go function returns the length of a string in bytes?",
			AcceptedAnswers: []string{"len"},
		},
		{
			Type:            "text",
			Difficulty:      "medium",
			Title:           "Append Capacity",
			Prompt:          "If you append an item to a slice that is at full capacity, the capacity is typically doubled. What built-in Go function creates a new slice under the hood?",
			AcceptedAnswers: []string{"append"},
		},

		// HARD TEXT
		{
			Type:            "text",
			Difficulty:      "hard",
			Title:           "Go Heap Profiler Tool",
			Prompt:          "Which tool in the Go toolchain is used to analyze profiling data (e.g., heap allocations or CPU profiles)?",
			AcceptedAnswers: []string{"pprof", "go tool pprof"},
		},
		{
			Type:            "text",
			Difficulty:      "hard",
			Title:           "Select Non-blocking Default",
			Prompt:          "Which keyword can be added to a select block in Go to make it non-blocking?",
			AcceptedAnswers: []string{"default"},
		},
		{
			Type:            "text",
			Difficulty:      "hard",
			Title:           "Unbuffered Channel Behavior",
			Prompt:          "If you send a value to an unbuffered channel when no goroutine is waiting to receive it, what state does the sending goroutine enter?",
			AcceptedAnswers: []string{"blocked", "blocking", "block"},
		},
	}

	for _, p := range problemsToSeed {
		var problemID string
		err = conn.QueryRow(ctx, `
			INSERT INTO problems (type, difficulty, title, prompt, is_published)
			VALUES ($1, $2, $3, $4, true)
			RETURNING id
		`, p.Type, p.Difficulty, p.Title, p.Prompt).Scan(&problemID)
		if err != nil {
			log.Fatalf("Failed to insert problem %q: %v\n", p.Title, err)
		}

		if p.Type == "mcq" {
			for _, opt := range p.Options {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_options (problem_id, text, is_correct)
					VALUES ($1, $2, $3)
				`, problemID, opt.Text, opt.IsCorrect)
				if err != nil {
					log.Fatalf("Failed to insert option %q for problem %q: %v\n", opt.Text, p.Title, err)
				}
			}
		} else if p.Type == "text" {
			for _, ans := range p.AcceptedAnswers {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_accepted_answers (problem_id, answer_text)
					VALUES ($1, $2)
				`, problemID, ans)
				if err != nil {
					log.Fatalf("Failed to insert accepted answer %q for problem %q: %v\n", ans, p.Title, err)
				}
			}
		}
	}

	log.Printf("Successfully seeded %d problems.\n", len(problemsToSeed))
	log.Println("Seeding completed successfully!")
}
