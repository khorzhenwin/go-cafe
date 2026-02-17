# go-cafe backend
# Run from backend dir or workspace root: make -C backend run

ROOT := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

.PHONY: run auth build swagger test unit-test integration-test migrate-up migrate-down migrate-create teardown teardown-reset docker-up docker-down docker-check help

# Default target
help:
	@echo "Targets (run from backend or: make -C backend <target>):"
	@echo "  make run             - Run the API server (requires .env and migrations applied)"
	@echo "  make auth            - Print a JWT (registers or logs in a user at localhost:8080)"
	@echo "  make swagger         - Generate Swagger docs to ./docs"
	@echo "  make build           - Generate Swagger docs and build API binary to ./bin/api"
	@echo "  make test            - Run unit tests only (no DB required)"
	@echo "  make unit-test       - Same as test"
	@echo "  make integration-test - Run integration tests (requires DB in .env)"
	@echo "  make migrate-up      - Apply migrations (requires DB in .env)"
	@echo "  make migrate-down    - Rollback one migration"
	@echo "  make migrate-create  - Create a new migration (usage: make migrate-create name=my_migration)"
	@echo "  make docker-check    - Verify Docker CLI and daemon status"
	@echo "  make docker-up       - Build and run the app via docker-compose"
	@echo "  make docker-down     - Stop docker-compose and remove containers"
	@echo "  make teardown        - Stop app/containers and remove build artifacts (non-destructive)"
	@echo "  make teardown-reset  - DESTRUCTIVE: teardown + rollback one migration"

run:
	cd $(ROOT) && go run ./cmd/api

# Print JWT for Swagger Auth button.
# Optional overrides:
#   make auth email=you@example.com password=secret123 name="Your Name" url=http://localhost:8080
auth:
	@cd $(ROOT) && \
	EMAIL=$${email:-"dev+$$(date +%s)@example.com"}; \
	PASSWORD=$${password:-secret123}; \
	NAME=$${name:-"Dev User"}; \
	URL=$${url:-http://localhost:8080}; \
	REG_RESP=$$(curl -s -X POST "$$URL/api/v1/auth/register" \
		-H "Content-Type: application/json" \
		-d "{\"email\":\"$$EMAIL\",\"name\":\"$$NAME\",\"password\":\"$$PASSWORD\"}"); \
	TOKEN=$$(printf '%s' "$$REG_RESP" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p'); \
	if [ -z "$$TOKEN" ]; then \
		LOGIN_RESP=$$(curl -s -X POST "$$URL/api/v1/auth/login" \
			-H "Content-Type: application/json" \
			-d "{\"email\":\"$$EMAIL\",\"password\":\"$$PASSWORD\"}"); \
		TOKEN=$$(printf '%s' "$$LOGIN_RESP" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p'); \
	fi; \
	if [ -z "$$TOKEN" ]; then \
		echo "Failed to obtain JWT. Is the API running at $$URL?"; \
		echo "Try: make docker-up"; \
		exit 1; \
	fi; \
	echo "$$TOKEN"

swagger:
	cd $(ROOT) && go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d ./cmd/api,./internal -o ./docs --parseInternal

build: swagger
	@mkdir -p $(ROOT)/bin
	cd $(ROOT) && go build -o bin/api ./cmd/api

# Unit tests: no build tags, no DB required (exclude server which has integration tests)
test: unit-test

unit-test:
	cd $(ROOT) && go test -v -count=1 ./internal/config/ ./internal/auth/ ./internal/user/ ./internal/cafelisting/ ./internal/rating/

# Integration tests: require -tags=integration and DB env
integration-test:
	cd $(ROOT) && set -a && [ -f .env ] && . ./.env || true; set +a; go test -v -count=1 -tags=integration ./internal/server/...

# Migrations (use DB_* from .env); run from ROOT so .env and migrations/ are found
migrate-up:
	cd $(ROOT) && go run ./cmd/migrate up

migrate-down:
	cd $(ROOT) && go run ./cmd/migrate down

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=my_migration"; exit 1; fi; \
	n=$$(ls $(ROOT)/migrations/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	n=$$(printf "%06d" $$((n + 1))); \
	touch "$(ROOT)/migrations/$${n}_$(name).up.sql" "$(ROOT)/migrations/$${n}_$(name).down.sql"; \
	echo "Created $(ROOT)/migrations/$${n}_$(name).up.sql and .down.sql"

# Optional: run all tests (unit + integration) when DB available
test-all: unit-test integration-test

# Docker Compose
docker-check:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "Docker CLI not found. Install Docker Desktop first."; \
		exit 1; \
	fi
	@if ! docker info >/dev/null 2>&1; then \
		echo "Docker daemon is not running. Start Docker Desktop and try again."; \
		exit 1; \
	fi
	@echo "Docker is available and daemon is running."

docker-up: docker-check
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up --build; \
	else \
		echo "docker compose/docker-compose not found, using plain docker build+run fallback"; \
		docker build -t go-cafe-backend . && \
		docker rm -f go-cafe-backend >/dev/null 2>&1 || true; \
		docker run --name go-cafe-backend --env-file .env -p 8080:8080 -d go-cafe-backend; \
	fi

docker-down: docker-check
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose down --remove-orphans; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down --remove-orphans; \
	else \
		docker rm -f go-cafe-backend >/dev/null 2>&1 || true; \
	fi

up: docker-up
down: docker-down

# Teardown: stop app/containers and remove bin/ (non-destructive; keeps DB tables/data)
teardown:
	@echo "Stopping any process on :8080 (if present)..."
	@if command -v lsof >/dev/null 2>&1; then \
		PIDS=$$(lsof -ti tcp:8080); \
		if [ -n "$$PIDS" ]; then kill -9 $$PIDS; fi; \
	fi
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose down --remove-orphans 2>/dev/null || true; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down --remove-orphans 2>/dev/null || true; \
	else \
		docker rm -f go-cafe-backend >/dev/null 2>&1 || true; \
	fi
	@rm -rf $(ROOT)/bin
	@echo "Teardown complete: containers stopped, bin/ removed, DB preserved"

# Destructive reset for local-only workflows.
teardown-reset: teardown
	@echo "WARNING: rolling back one migration (destructive depending on migration)"
	cd $(ROOT) && go run ./cmd/migrate down 2>/dev/null || true
	@echo "Teardown reset complete: migration rolled back"
