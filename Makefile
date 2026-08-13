.PHONY: dev-server dev-client db-up migrate-up migrate-down seed build-client build loadtest smoke

MIGRATION_DIR=server/migrations
DB_CONN ?= postgres://postgres:postgres@localhost:5432/monopoly?sslmode=disable
LISTEN_ADDR ?= 0.0.0.0:8080

dev-server:
	cd server && go run ./cmd/server

seed:
	cd server && go run ./cmd/seed

dev-client:
	cd client && npm install && npm run dev

# Build Vue into server/cmd/server/dist so go:embed can pack it into the binary.
# Unset VITE_WS_BASE_URL so a stray client/.env cannot bake localhost into LAN builds
# (production JS ignores it anyway; this keeps the embed free of localhost strings).
build-client:
	cd client && npm install && env -u VITE_WS_BASE_URL npm run build
	rm -rf server/cmd/server/dist
	mkdir -p server/cmd/server/dist
	cp -r client/dist/. server/cmd/server/dist/
	@test -f server/cmd/server/dist/index.html || (echo "missing embedded index.html" && exit 1)
	@js=$$(sed -n 's/.*src="\/assets\/\([^"]*\)".*/\1/p' server/cmd/server/dist/index.html); \
		test -n "$$js" || (echo "index.html has no /assets/*.js ref" && exit 1); \
		test -f "server/cmd/server/dist/assets/$$js" || (echo "missing $$js — stale or incomplete client build" && exit 1)
	@! grep -R -n -E 'localhost:8080|ws://localhost|wss://localhost' server/cmd/server/dist/assets --include='*.js' \
		|| (echo "embedded JS contains localhost WS URL — refuse to ship" && exit 1)
	@echo "Client assets copied to server/cmd/server/dist (same-origin WS)."

# Single deployable binary with the frontend embedded.
build: build-client
	cd server && go build -o ../bin/monopoly-server ./cmd/server
	@echo "Built bin/monopoly-server (frontend embedded)."

# Re-verify phase-5 fairness + embed sanity against the deploy mux / embedded assets.
smoke: build-client
	cd server && go test ./cmd/server -count=1 -run 'TestDeployMux_|TestStatic'

db-up:
	docker compose up -d
	@echo "Waiting for PostgreSQL..."
	@until [ "$$(docker compose ps -q db | xargs docker inspect -f '{{.State.Health.Status}}')" = "healthy" ]; do \
		sleep 1; \
	done
	@echo "PostgreSQL is ready."

migrate-up:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path=$(MIGRATION_DIR) -database "$(DB_CONN)" up; \
	else \
		docker run --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations --network host migrate/migrate -path=/migrations -database "$(DB_CONN)" up; \
	fi

migrate-down:
	@if command -v migrate >/dev/null 2>&1; then \
		echo "y" | migrate -path=$(MIGRATION_DIR) -database "$(DB_CONN)" down; \
	else \
		echo "y" | docker run -i --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations --network host migrate/migrate -path=/migrations -database "$(DB_CONN)" down; \
	fi

# Full-rotation class load test (sequential turns, realistic answer speeds).
loadtest:
	cd server && go run ./cmd/loadtest -players 24 -correct-rate 0.8
