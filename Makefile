SHELL := /bin/bash

.PHONY: help run run-backend run-frontend teardown

help:
	@echo "go-cafe root targets:"
	@echo "  make run           - Run backend + frontend together for local testing"
	@echo "  make run-backend   - Run backend only"
	@echo "  make run-frontend  - Run frontend only"
	@echo "  make teardown      - Stop all services and clean backend/frontend artifacts"

run:
	@echo "Starting backend (:8080) and frontend (:3000)..."
	@cd backend && go run ./cmd/api & BACK_PID=$$!; \
	$(MAKE) -C frontend run-frontend-helper & FRONT_PID=$$!; \
	trap 'kill $$BACK_PID $$FRONT_PID 2>/dev/null || true' INT TERM EXIT; \
	while kill -0 $$BACK_PID 2>/dev/null && kill -0 $$FRONT_PID 2>/dev/null; do sleep 1; done; \
	if ! kill -0 $$BACK_PID 2>/dev/null; then \
		wait $$BACK_PID; STATUS=$$?; \
	else \
		wait $$FRONT_PID; STATUS=$$?; \
	fi; \
	kill $$BACK_PID $$FRONT_PID 2>/dev/null || true; \
	wait $$BACK_PID $$FRONT_PID 2>/dev/null || true; \
	exit $$STATUS

run-backend:
	@$(MAKE) -C backend run

run-frontend:
	@$(MAKE) -C frontend run-frontend-helper

teardown:
	@echo "Stopping local processes on :3000 and :8080 (if present)..."
	@if command -v lsof >/dev/null 2>&1; then \
		PIDS=$$((lsof -ti tcp:3000 2>/dev/null; lsof -ti tcp:8080 2>/dev/null) | sort -u); \
		if [ -n "$$PIDS" ]; then \
			kill $$PIDS 2>/dev/null || true; \
			sleep 1; \
			kill -9 $$PIDS 2>/dev/null || true; \
		fi; \
	fi
	@$(MAKE) -C backend teardown
	@$(MAKE) -C frontend teardown-helper
	@echo "Root teardown complete."
