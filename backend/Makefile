# go-cafe backend

.PHONY: run build test unit-test integration-test migrate-up migrate-down migrate-create help

# Default target
help:
	@echo "Targets:"
	@echo "  make run             - Run the API server (requires .env and migrations applied)"
	@echo "  make build           - Build the API binary to ./bin/api"
	@echo "  make test            - Run unit tests only (no DB required)"
	@echo "  make unit-test       - Same as test"
	@echo "  make integration-test - Run integration tests (requires DB in .env)"
	@echo "  make migrate-up      - Apply migrations (requires DB in .env)"
	@echo "  make migrate-down    - Rollback one migration"
	@echo "  make migrate-create  - Create a new migration (usage: make migrate-create name=my_migration)"

run:
	go run ./cmd/api

build:
	@mkdir -p bin
	go build -o bin/api ./cmd/api

# Unit tests: no build tags, no DB required (exclude server which has integration tests)
test: unit-test

unit-test:
	go test -v -count=1 ./internal/config/ ./internal/auth/ ./internal/user/ ./internal/cafelisting/ ./internal/rating/

# Integration tests: require -tags=integration and DB env
integration-test:
	go test -v -count=1 -tags=integration ./internal/server/...

# Migrations (use DB_* from .env)
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=my_migration"; exit 1; fi; \
	n=$$(ls migrations/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	n=$$(printf "%06d" $$((n + 1))); \
	touch "migrations/$${n}_$(name).up.sql" "migrations/$${n}_$(name).down.sql"; \
	echo "Created migrations/$${n}_$(name).up.sql and .down.sql"

# Optional: run all tests (unit + integration) when DB available
test-all: unit-test integration-test
