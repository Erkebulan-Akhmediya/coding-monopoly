# Coding Monopoly — Context Primer

## 1. Project Overview & Status

**Coding Monopoly** is a browser-based, LAN-only, real-time multiplayer educational board game for software engineering students.
- **Topology**: One machine (e.g. instructor laptop) runs the monolithic server binary that hosts the Go backend, serves the embedded Vue SPA, and connects to Postgres. Student devices connect as clients over the local network via their browser.
- **Current Status**: All initial 11 phases (Phases 0–11) from `plan.md` are **complete**. The system has a fully functional game engine, WebSocket real-time transport, PostgreSQL schema and migrations, content management CRUD API, admin live spectator dashboard, reconnect resilience, audio-visual polish, and single-binary LAN deployment with full-rotation load testing. Current project focus is **bug fixes, edge-case hardening, performance tuning, and UX polishing**.

---

## 2. Tech Stack & Invariants

### Backend (Go)
- **Language & Runtime**: Go 1.22+ (standard library wherever possible).
- **External Dependencies (Minimal)**:
  - `github.com/gorilla/websocket` — WebSocket server & client transport.
  - `github.com/jackc/pgx/v5` (`pgxpool`) — High-performance PostgreSQL driver & connection pool.
  - `github.com/joho/godotenv` — Development `.env` loader.
- **Architecture**: No heavyweight web frameworks (pure `net/http` stdlib mux). Single monolithic process embedding frontend assets via `//go:embed`.

### Frontend (Vue 3 + TypeScript)
- **Framework**: Vue 3 + TypeScript built with Vite.
- **Strict Invariant**: **Options API only** (`defineComponent({ data, computed, methods, watch, mounted })`). **No Composition API, no `<script setup>`**.
- **State Management**: Reactive store modules (`src/store.ts` for players, `src/adminStore.ts` for admin spectator) built with Vue's `reactive({})` — Options API compatible.
- **Styling**: Modern, responsive dark-themed CSS with CSS Grid for the 32-cell perimeter board.
- **Libraries**:
  - `roll-a-die` — 3D canvas dice rolling animations.
  - Web Audio API — sound effects (`soundService.ts`) with custom synthesized audio fallback.
- **Localization (`src/i18n`)**: Custom lightweight i18n supporting English (`en`), Kazakh (`kk`), and Russian (`ru`).

### Database (PostgreSQL)
- **Schema Management**: `golang-migrate` (`server/migrations/000001_init.up.sql`).
- **Tables**: `games`, `players`, `board_cells`, `problems`, `problem_options`, `problem_accepted_answers`, `submissions`, `game_events`.

---

## 3. Board Layout & Cell Effects

The board contains **32 cells** around a square perimeter (8 per side, corners shared, indices 0..31):

### Corner Cells
- **0 — `Deploy` (Start / GO)**: +100 XP bonus upon landing directly. Passing GO during any roll movement grants a +50 XP lap bonus.
- **8 — `Code Freeze` (Jail Equivalent)**: Player enters frozen state (`in_code_freeze = true`).
- **16 — `Coffee Break` (Rest)**: Neutral rest spot (+5 XP boost).
- **24 — `Deadline` (Swing Event)**: 50% chance of +50 XP bonus, 50% chance of -20 XP penalty.

### 28 Perimeter Effect Cells
- **`xp_gain`**: Small (+10 XP), Medium (+25 XP), or Large (+50 XP).
- **`xp_loss`**: Small (-10 XP) or Medium (-25 XP). Cannot drop XP below 0.
- **`double_xp`**: Multiplies the next XP gain by 2x (consumed upon use).
- **`free_pass`**: Grants an inventory pass ticket (`free_passes++`).
- **`skip_next`**: Player's next sequential turn is skipped (`skip_next_turn = true`).
- **`teleport`**: Relocates player token directly to a target cell (e.g. cell 0 `Deploy`, cell 8 `Code Freeze`, or cell 16 `Coffee Break`).
- **`special_challenge`**: Instant bonus challenge reward (+30 to +35 XP).
- **`mystery`**: Random roll from the effect catalog (XP gain/loss, double XP, free pass, teleport, or special challenge).

---

## 4. Game Rules & Sequential Turn Flow

### Sequential Turn Lifecycle
Exactly one player is active at any time; all other players watch live as spectators.
1. **Turn Start (`turn_started`)**: Active player is determined in circular join order (skipping disconnected players or players with active `skip_next_turn`).
2. **Difficulty Selection (`choose_level`)**: Active player picks question difficulty via UI:
   - **Easy**: 30s timer deadline, 1 dice roll reward if correct.
   - **Medium**: 45s timer deadline, 2 dice rolls reward if correct.
   - **Hard**: 60s timer deadline, 3 dice rolls reward if correct.
3. **Question Assignment & Redacted Broadcast (`question_started`)**:
   - DB provider selects a random published question matching the chosen difficulty.
   - **Active Player Private Channel**: Receives full problem ID, prompt, and options.
   - **Spectator Broadcast Channel**: Receives only redacted metadata (`difficulty`, `deadline`) for countdown display — prompt and options are hidden.
4. **Answer Submission & Grading (`submit_answer` / Timeout)**:
   - Server enforces **exactly one attempt per turn**.
   - Server-authoritative timer (`time.AfterFunc`) races against `submit_answer`.
   - Thread-safe resolution guarded by `Room.mu` and `turn.resolved = true`. First caller (submission or timeout) wins; subsequent events are no-ops.
5. **Rolls & Movement Resolution (`roll_resolved`)**:
   - **Correct Answer**: Player rolls the die $N$ times (1, 2, or 3). Each roll is resolved individually (1–6 steps movement → check passed GO lap bonus → apply landed cell effect → broadcast roll result).
   - **Incorrect / Timed Out**: 0 rolls; turn ends immediately.
6. **Turn End (`turn_ended`)**: Turn ends, active turn state is cleared, and turn advances to the next player.
7. **Win Condition (`game_over`)**: First player to reach **500 XP** (`DefaultTargetXP`) triggers match conclusion with final standings, or admin manually triggers game over.

---

## 5. Question Types & Fairness Architecture

### Problem Types (No Code Sandbox / Execution)
1. **MCQ (`mcq`)**:
   - Multiple choice options with one or more marked correct.
   - Server compares submitted option UUIDs against correct option UUIDs as sets.
2. **Text (`text`)**:
   - Free-form text questions (bug finding, code output, missing syntax).
   - Server normalizes strings (trims whitespace, case-folding) against stored accepted answers.

### Fairness & Privacy Guarantees
- Option correctness flags (`option.Correct`) are marked `json:"-"` in Go structs and never serialized to wire payloads.
- Spectators receive zero prompt/options content over WebSocket before/during the turn.
- Correct answers are sent privately to the active player only upon turn conclusion for learning feedback, never broadcast to spectators (allowing the question to be recycled later).

---

## 6. Real-Time WebSocket & Reconnect Architecture

### Concurrency Model
- **Authoritative Hub**: `internal/ws/Hub` runs an event loop in a dedicated goroutine managing client registrations, unregistrations, join requests, and room broadcasts without coarse locks.
- **Room Engine**: `internal/room/Room` uses `sync.RWMutex` to protect all gameplay mutations, turn progression, timers, and standings.
- **Interfaces**:
  - `Broadcaster`: Broadcasts to all clients in a room.
  - `PrivateBroadcaster`: Sends private messages to a single client socket.
  - `ExcludingBroadcaster`: Broadcasts to room while omitting a specific client socket (used for redacted question metadata).

### Reconnect Resilience
- Client stores identity (`playerId`, `playerName`, `roomId`) in browser `sessionStorage`.
- Reconnecting client sends `player_id` with `join` message:
  - If recognized and disconnected, the slot is reclaimed (`resumed = true`) without resetting position, XP, or modifiers.
  - If the player is the active player mid-question, the server re-delivers the active question payload and countdown.
- **Disconnect Grace Period**: Active player has a **5-second grace period** (`DefaultDisconnectGrace`) on disconnect before turn forfeiture, preventing brief network blips from ruining a turn.
- Mid-game joins by new players are rejected (`ErrGameInProgress`) once a game has started; only reclaims of existing player slots are permitted.

---

## 7. Admin & Spectator Subsystem

The admin subsystem is accessed via `?admin=1` in the client URL:

### 1. Question Bank Management (REST API)
- Authentication: `POST /admin/login` returns an HMAC-SHA256 bearer token with configurable TTL (default 15m).
- Endpoints under `/admin/problems`:
  - `GET /admin/problems` — Filter by `type`, `difficulty`, `is_published`, with pagination.
  - `POST /admin/problems` — Create new MCQ or text problem.
  - `GET /admin/problems/{id}` — Get single problem details.
  - `PUT /admin/problems/{id}` — Update problem title, prompt, options, or accepted answers.
  - `DELETE /admin/problems/{id}` — Delete problem.
  - `PATCH /admin/problems/{id}/publish` — Toggle published state.
- Multi-Room API: `GET /admin/rooms` and `POST /admin/rooms` for creating and monitoring game rooms.

### 2. Live Room Spectator & Control (WebSocket `/ws/admin`)
- Upgrades after verifying bearer token.
- Real-time mirrored board and player status.
- Chronological event feed (`GameEvent`, capped at 200 items in `adminStore.ts`).
- Administrative actions:
  - `admin_start` — Manually start the game match from lobby.
  - `admin_pause` — Pause/resume game timer and player actions.
  - `admin_kick` — Forcibly remove a player from the room.
  - `admin_skip_turn` — Skip current active player's turn (e.g. if player is AFK before choosing level).
  - `admin_end_game` — Forcibly conclude match and trigger game over summary.

---

## 8. Directory & Codebase Structure

```
coding-monopoly/
├── Makefile                     # Build, run, test, migration, seed, smoke, loadtest automation
├── docker-compose.yml           # Local PostgreSQL service container
├── plan.md                      # Phased roadmap (Phases 0-11 complete)
├── context.md                   # This context primer
├── DEPLOY.md                    # Classroom/LAN deployment documentation & test benchmarks
├── server/
│   ├── cmd/
│   │   ├── server/              # Main server binary entry point
│   │   │   ├── main.go          # CLI entry, env loading, DB pool init, HTTP server
│   │   │   ├── mux.go           # Route registration (REST, WS, static assets)
│   │   │   ├── static.go        # Embedded frontend filesystem (//go:embed dist)
│   │   │   ├── static_test.go   # Static delivery tests
│   │   │   └── deploy_fairness_test.go # End-to-end fairness, race, and deployment tests
│   │   ├── seed/
│   │   │   └── main.go          # Database seed script (board cells + problem bank)
│   │   └── loadtest/
│   │       └── main.go          # 24-player sequential rotation load test simulator
│   ├── internal/
│   │   ├── admin/
│   │   │   ├── admin.go         # REST API handler, HMAC auth, problem CRUD, room manager
│   │   │   ├── admin_test.go    # Unit tests for admin API
│   │   │   └── admin_room_test.go # Room manager unit tests
│   │   ├── problems/
│   │   │   └── selector.go      # Random published problem selection queries
│   │   ├── room/
│   │   │   ├── room.go          # Core game engine: turns, timers, grading, admin controls
│   │   │   ├── board.go         # 32-cell board definition
│   │   │   ├── effect.go        # Cell effect dispatch table and handlers
│   │   │   ├── player.go        # Player state model
│   │   │   ├── room_test.go     # Turn engine & movement tests
│   │   │   ├── question_test.go # Question grading and timer race tests
│   │   │   ├── reconnect_test.go# Reconnect & grace period tests
│   │   │   └── room_admin_test.go # Admin pause/kick/skip tests
│   │   └── ws/
│   │       ├── hub.go           # Authoritative WebSocket connection hub & event loop
│   │       ├── client.go        # WebSocket connection wrapper, pump loops
│   │       ├── handler.go       # HTTP -> WS upgrade handlers (/ws and /ws/admin)
│   │       ├── message.go       # JSON message protocol types & payloads
│   │       ├── question_provider.go # PostgreSQL-backed QuestionProvider bridge
│   │       ├── ws_test.go       # Hub & protocol integration tests
│   │       └── ws_admin_test.go # Admin WebSocket integration tests
│   └── migrations/
│       ├── 000001_init.up.sql   # Initial schema migration
│       └── 000001_init.down.sql # Migration rollback
└── client/
    ├── package.json             # Frontend dependencies (Vue 3, Vite, TypeScript)
    ├── src/
    │   ├── App.vue              # Root component: routes LobbyView vs BoardView vs AdminView
    │   ├── main.ts              # Vue application bootstrap & i18n mounting
    │   ├── style.css            # Global dark-theme styles, board grid layout, animations
    │   ├── store.ts             # Player reactive state store
    │   ├── adminStore.ts        # Admin spectator reactive state store
    │   ├── i18n/                # Localization module
    │   │   ├── index.ts         # t() helper, language state, reactive switch
    │   │   └── locales/         # en.ts, kk.ts, ru.ts translation dictionaries
    │   ├── services/
    │   │   ├── websocketService.ts      # Player WebSocket client & event router
    │   │   ├── adminWebsocketService.ts # Admin WebSocket client & event router
    │   │   ├── adminApiService.ts       # Admin REST API client & validators
    │   │   ├── soundService.ts          # Web Audio & synthesizers
    │   │   ├── tokenMovement.ts         # Step-by-step token hop animations
    │   │   └── serverUrls.ts            # Dynamic same-origin URL resolver
    │   └── components/
    │       ├── LobbyView.vue            # Player join screen & room lobby
    │       ├── BoardView.vue            # Main 32-cell board view & center stage
    │       ├── ProblemPanel.vue         # Active player / spectator problem view
    │       ├── LevelPicker.vue          # Easy / Medium / Hard selector
    │       ├── McqPanel.vue             # Multiple choice answer component
    │       ├── TextAnswerPanel.vue      # Free-text answer component
    │       ├── DiceOverlay.vue          # 3D dice animation overlay
    │       ├── LandedCellPreview.vue    # Magnified destination cell preview
    │       ├── EffectToastStack.vue     # Pop-up effect notification toasts
    │       ├── Leaderboard.vue          # Live player leaderboard
    │       ├── PauseOverlay.vue         # Game paused indicator
    │       ├── EndGameSummary.vue       # Match podium & final standings modal
    │       ├── PlayerToken.vue          # Player avatar marker on cells
    │       ├── LanguageSwitcher.vue     # Language picker dropdown
    │       ├── AdminView.vue            # Admin dashboard container (Tabs + Login)
    │       ├── AdminQuestionList.vue    # Question bank management view
    │       ├── AdminQuestionFormModal.vue # Question editor modal
    │       └── AdminSpectatorView.vue   # Live room spectator & room controls
```

---

## 9. Common Workflows & Commands

```bash
# Start Postgres database container
make db-up

# Run database migrations
make migrate-up

# Seed board cells and starter questions
make seed

# Run backend Go unit & integration test suites
make test

# Run deploy smoke test (fairness, submit-race, embedded asset verification)
make smoke

# Run 24-player full-rotation sequential load test
make loadtest

# Build complete single binary with embedded frontend
make build

# Start server binary
LISTEN_ADDR=0.0.0.0:8080 ADMIN_PASSWORD=secret ADMIN_TOKEN_SECRET=devsecret ./bin/monopoly-server
```

---

## 10. Key Guidelines for Bug Fixing & Polishing

1. **Maintain Options API Strictness**: All new or modified Vue components in `client/src/components/` must use Vue 3 Options API (`defineComponent({ ... })`). Never introduce Composition API `<script setup>` or `ref()` inside SFC script blocks.
2. **Preserve Concurrency Guarantees**:
   - All room state modifications must acquire `Room.mu.Lock()`.
   - Client connection registration/unregistration must go through `Hub.Run()` channels.
   - Timer vs Submission races must be protected by the `turn.resolved` check.
3. **Preserve Content Redaction**: Never send question prompts, options, or correctness to public spectator broadcasts. Only send redacted metadata (`difficulty`, `deadline`).
4. **Dynamic Host Resolution**: In frontend services, always resolve API and WebSocket URLs dynamically from `window.location` (via `src/services/serverUrls.ts`), avoiding hardcoded `localhost` to ensure seamless LAN functionality.
5. **Preserve i18n Translations**: Any new user-visible text in components must have corresponding translation keys added in `src/i18n/locales/en.ts`, `kk.ts`, and `ru.ts`.