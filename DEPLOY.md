# LAN deployment

One instructor machine runs a single Go binary that embeds the Vue frontend and
serves HTTP + WebSockets. Student devices only need a browser pointed at that
machine's LAN address.

## Build the binary

From the repo root (needs Node + Go + Docker for Postgres):

```bash
make db-up
make migrate-up
make seed          # board + sample questions
make build         # builds Vue → embeds into bin/monopoly-server
```

`make build` runs `npm run build`, copies assets into `server/cmd/server/dist`,
then `go build`s with `//go:embed`. The result is `bin/monopoly-server` —
no separate static file server and no Node process in production.

## Configure environment

Copy and edit `server/.env` (or export the same variables before launching):

| Variable | Purpose |
|---|---|
| `DATABASE_URL` | Postgres DSN for the **target** instance (see migrations below) |
| `LISTEN_ADDR` | Bind address, e.g. `0.0.0.0:8080` (default if unset) |
| `PORT` | Alternative to `LISTEN_ADDR`: port only, always on `0.0.0.0` |
| `ADMIN_PASSWORD` | Admin panel / spectator login (change from the scaffold default) |
| `ADMIN_TOKEN_SECRET` | HMAC secret for admin bearer tokens (long random string) |
| `ADMIN_TOKEN_TTL` | Token lifetime (default `15m`) |

`godotenv` autoloads `.env` from the **process working directory**. Either:

```bash
cd server && ../bin/monopoly-server
# or
set -a && source server/.env && set +a && ./bin/monopoly-server
```

## LAN IP binding

The server defaults to `0.0.0.0:8080` so it accepts connections on every
interface (loopback **and** the classroom NIC).

1. On the host, find the LAN IPv4 address:

   ```bash
   hostname -I
   # or: ip -4 addr show
   ```

2. Start the binary (example):

   ```bash
   LISTEN_ADDR=0.0.0.0:8080 ./bin/monopoly-server
   ```

3. Students open `http://<LAN-IP>:8080/` (e.g. `http://10.10.40.69:8080/`).
   Admin panel: `http://<LAN-IP>:8080/?admin=1`.

The production frontend derives `ws://<same-host>/ws` from `window.location`,
so students must use the LAN IP (or a local DNS name), **not** a bookmark that
still points at `localhost`.

To bind only one NIC:

```bash
LISTEN_ADDR=10.10.40.69:8080 ./bin/monopoly-server
```

## Firewall / port notes

| Port | Direction | Who | Notes |
|---|---|---|---|
| **8080/tcp** | inbound to instructor host | student browsers | HTTP + WebSocket (`/`, `/ws`, `/ws/admin`, `/admin`) |
| **5432/tcp** | instructor host → Postgres | server process only | Prefer localhost / private Docker network. **Do not** expose Postgres to the student LAN. |

Examples:

```bash
# ufw
sudo ufw allow 8080/tcp comment 'coding-monopoly'
sudo ufw deny 5432/tcp

# firewalld
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

Campus Wi‑Fi client isolation (AP/client isolation) will block students from
reaching the instructor PC even if the host firewall is open — turn isolation
off for the lab SSID, or put everyone on the same non-isolated VLAN.

WebSockets use the same TCP port as HTTP (upgrade on `/ws`). No second port.

## Migrations against the target Postgres

Point `DB_CONN` / `DATABASE_URL` at the **classroom** database before migrating.
Never migrate blindly against a shared prod instance you do not intend to change.

```bash
# Local docker-compose (default)
make db-up
make migrate-up
make seed

# Remote / lab Postgres
export DB_CONN='postgres://USER:PASS@LAB_HOST:5432/monopoly?sslmode=require'
export DATABASE_URL="$DB_CONN"
make migrate-up
cd server && go run ./cmd/seed
```

`make migrate-up` uses the `migrate` CLI if installed, otherwise runs
`migrate/migrate` via Docker with `--network host`.

Schema lives in `server/migrations/`. Apply **up** on a fresh database, then
`seed` once for the 32-cell board and starter question bank. Re-seeding
truncates game/content tables — do not re-seed mid-session.

Rollback (destroys schema):

```bash
make migrate-down
```

## Run checklist

1. Postgres healthy and migrated; questions published (seed or admin UI).
2. Strong `ADMIN_PASSWORD` / `ADMIN_TOKEN_SECRET` set.
3. `LISTEN_ADDR=0.0.0.0:8080` (or your LAN bind).
4. Host firewall allows **8080/tcp**; Postgres not reachable from students.
5. `./bin/monopoly-server` — confirm `/health` returns `ok`.
6. From a student laptop: open `http://<LAN-IP>:8080/`, join lobby, play.

## Full-rotation load test (measured)

Turns are **sequential**. The meaningful metric is wall-clock time for one
complete pass through every player at realistic answer speeds — not raw
concurrent request throughput.

```bash
make loadtest
# or: cd server && go run ./cmd/loadtest -players 24 -correct-rate 0.8
```

### Scenario

| Parameter | Value |
|---|---|
| Class size | **24** players (one room) |
| Think times | easy ~12s, medium ~20s, hard ~30s (under 30/45/60s deadlines) |
| Difficulty mix | ~60% easy / 30% medium / 10% hard (weighted RNG, seed 42) |
| Correct rate | 80% target |

### Measured results (2026-08-13, this repo)

| Metric | Value |
|---|---|
| Difficulty mix realized | easy=18, medium=5, hard=1 |
| Correct answers | 20 / 24 (83%) |
| Sum of think times | **5m46s** |
| Avg think / turn | **14.4s** |
| Server overhead (WS + grade + dice + advance) | **~28s** for the whole rotation |
| **Full rotation** | **6m13.8s** |
| Per-player average | **15.6s** |

### Sanity check vs a normal class period

| Period length | Full rotations that fit |
|---|---|
| 50 minutes | ~**8.0** |
| 90 minutes | ~**14.4** |

**Verdict:** one 24-player rotation (~6¼ minutes) fits comfortably in a normal
lab/lecture block, with headroom for multiple laps, discussion, and admin
pauses. Worst-case timeouts (every player burning the full 30–60s deadline)
would roughly double cycle time and are **not** the planning target — coach
students to pick a difficulty and submit when ready.

Re-run anytime with `make loadtest` after engine changes.

## Fairness / security notes (phase 11 pass)

- Question **prompt/options** go only to the active player's socket; spectators
  get difficulty + deadline. Correct answers are private review for the solver
  only (`json:"-"` on option correctness; redaction covered by WS tests).
- Admin spectator WS requires a validated token **before** upgrade; admin
  clients cannot `choose_level` / `submit_answer`.
- Pause blocks choose/submit until resumed.
- Production clients use same-origin WS URLs (no hardcoded `localhost`).
- Change scaffold admin secrets before a real class. Keep Postgres off the LAN.
