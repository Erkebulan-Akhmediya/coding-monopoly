## Final phased plan

| Phase | Goal | Key deliverables |
|---|---|---|
| 0 | Scaffolding | Go module + Vite/Vue3/TS project, docker-compose Postgres, migration tooling |
| 1 | Schema & seed data | `games`, `players`, `board_cells`, `problems` (type: mcq/text), `problem_options`, `problem_accepted_answers`, `submissions`, `game_events`; seeded board + questions |
| 2 | WS transport & lobby | Connect, join by name, presence broadcast, `state_sync` — no game logic yet |
| 3 | Core game engine (sequential turns) | Turn order, active-player enforcement, per-roll movement, cell effect dispatch, basic disconnect skip |
| 4 | Admin content API | `internal/admin`: token auth, mcq/text CRUD, per-type validation, publish/draft |
| 5 | Turn timer & answer grading | Server-authoritative countdown, race-safe submit-vs-timeout resolution, mcq/text grading, redacted broadcasts |
| 6 | Frontend board & lobby | `LobbyView`, `BoardView`, turn indicator, waiting-vs-active state |
| 7 | Frontend problem UI | Level picker (active player only), MCQ/text panels, synced countdown, redacted spectator view |
| 8 | Live spectator/admin view | Read-only board, event feed, start/pause/kick, manual skip-turn |
| 9 | Admin panel frontend | Login, question list + filters, mcq/text forms, publish toggle |
| 10 | Polish & resilience | Reconnect/resume (incl. mid-question), disconnect grace period, animations, leaderboard, end screen |
| 11 | LAN deployment hardening | Embedded binary, deployment doc, full-rotation load test, final fairness/security pass |

Same rule as before: one phase = one agent session, implement then review in a fresh session, fix before moving on. 5 is the highest-stakes phase now that there's no judge — the timer race and content redaction are the two things worth the most scrutiny in this entire plan.