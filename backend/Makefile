# go-cafe backend
# Run from backend dir or workspace root: make -C backend run

ROOT := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

.PHONY: run build test unit-test integration-test migrate-up migrate-down migrate-create teardown docker-up docker-down help

# Default target
help:
	@echo "Targets (run from backend or: make -C backend <target>):"
	@echo "  make run             - Run the API server (requires .env and migrations applied)"
	@echo "  make build           - Build the API binary to ./bin/api"
	@echo "  make test            - Run unit tests only (no DB required)"
	@echo "  make unit-test       - Same as test"
	@echo "  make integration-test - Run integration tests (requires DB in .env)"
	@echo "  make migrate-up      - Apply migrations (requires DB in .env)"
	@echo "  make migrate-down    - Rollback one migration"
	@echo "  make migrate-create  - Create a new migration (usage: make migrate-create name=my_migration)"
	@echo "  make docker-up       - Build and run the app via docker-compose"
	@echo "  make docker-down     - Stop docker-compose and remove containers"
	@echo "  make teardown        - Stop app, roll back migrations, remove build artifacts"

run:
	cd $(ROOT) && go run ./cmd/api

build:
	@mkdir -p $(ROOT)/bin
	cd $(ROOT) && go build -o bin/api ./cmd/api

# Unit tests: no build tags, no DB required (exclude server which has integration tests)
test: unit-test

unit-test:
	cd $(ROOT) && go test -v -count=1 ./internal/config/ ./internal/auth/ ./internal/user/ ./internal/cafelisting/ ./internal/rating/

# Integration tests: require -tags=integration and DB env
integration-test:
	cd $(ROOT) && go test -v -count=1 -tags=integration ./internal/server/...

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
docker-up:
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	else \
		docker-compose up --build; \
	fi

docker-down:
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose down --remove-orphans; \
	else \
		docker-compose down --remove-orphans; \
	fi

# Teardown: stop app, roll back migrations (drops gocafe_* tables) and remove bin/
teardown:
	@echo "Stopping any process on :8080 (if present)..."
	@if command -v lsof >/dev/null 2>&1; then \
		PIDS=$$(lsof -ti tcp:8080); \
		if [ -n "$$PIDS" ]; then kill -9 $$PIDS; fi; \
	fi
	@cd $(ROOT) && if docker compose version >/dev/null 2>&1; then \
		docker compose down --remove-orphans 2>/dev/null || true; \
	else \
		docker-compose down --remove-orphans 2>/dev/null || true; \
	fi
	cd $(ROOT) && go run ./cmd/migrate down 2>/dev/null || true
	@rm -rf $(ROOT)/bin
	@echo "Teardown complete: migrations rolled back, bin/ removed"
