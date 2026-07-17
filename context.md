## Context primer

Building a browser-based, LAN-only, real-time multiplayer educational board game for software engineering students. One machine runs the server (does not play, but can spectate/administer); others connect as clients over the local network.

**Stack**: server in Go, standard library only plus two minimal dependencies (`gorilla/websocket`, `jackc/pgx` — no framework). Client in Vue 3 + TypeScript, **Options API only** (no Composition API, no `<script setup>`). Database is Postgres. Real-time sync over WebSockets.

**Board**: 32 cells around a square perimeter (8 per side, corners shared). Four corners: `Deploy` (start/GO, lap bonus), `Code freeze` (jail-equivalent, solve an easy question or wait out turns to leave), `Coffee break` (neutral rest), `Deadline` (one big swing event). The other 28 cells carry effects: XP gain (S/M/L), XP loss (S/M), mystery/random event, teleport, skip-next-turn, double-XP, free pass, special bonus challenge.

**Turn flow — sequential.** Exactly one player is active at a time; everyone else watches live and waits. Active player picks a difficulty (easy/medium/hard, a UI picker) → server assigns a question of that difficulty → a server-authoritative countdown starts (**30s easy / 45s medium / 60s hard**) → player submits an answer before the deadline, or the deadline fires first → **exactly one attempt per turn** (no re-pick-and-retry — retrying would extend everyone else's wait, so failure or timeout simply ends the turn) → correct: roll dice N times (1/2/3 by difficulty), resolving each roll individually (move, apply that cell's effect, then roll again if any remain) → incorrect or timed out: turn ends immediately with zero rolls → advance to the next player in turn order (join order).

**Fairness rule**: the actual question content (MCQ options / text prompt) is sent only to the active player's connection. Everyone else receives redacted metadata only — active player, difficulty, deadline — so they can watch a live countdown without previewing content they might draw themselves later. The correct answer is shown privately to the player who answered (for their own learning) but is never broadcast to spectators, since the same question can be drawn again by someone else later in the game.

**Two problem types, no code execution anywhere in the system:**
- **MCQ** — options with one or more marked correct; server compares submitted option id(s) exactly.
- **Text** — free-text answer, used for "what's the output of this program," "which line has the bug," fill-in-the-missing-code, etc. Server compares the submission, normalized (trimmed, case-folded), against a stored list of accepted answers for that question.

There is no build/compile/sandbox subsystem in this version — that entire concern is gone.

**Admin panel** stays split: content management (REST CRUD over the mcq/text question bank, works with no game running) vs. live room control (WebSocket spectator view + start/pause/kick, tied to an active room).