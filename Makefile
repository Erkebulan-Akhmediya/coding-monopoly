.PHONY: dev-server dev-client db-up migrate-up migrate-down

MIGRATION_DIR=migrations
DB_CONN=postgres://postgres:postgres@localhost:5432/monopoly?sslmode=disable

dev-server:
	cd server && go run ./cmd/server

dev-client:
	cd client && npm install && npm run dev

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
