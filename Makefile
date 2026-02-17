SHELL := /bin/bash

.PHONY: help run run-backend run-frontend

help:
	@echo "go-cafe root targets:"
	@echo "  make run           - Run backend + frontend together for local testing"
	@echo "  make run-backend   - Run backend only"
	@echo "  make run-frontend  - Run frontend only"

run:
	@echo "Starting backend (:8080) and frontend (:3000)..."
	@cd backend && go run ./cmd/api & BACK_PID=$$!; \
	cd frontend && npm run dev & FRONT_PID=$$!; \
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
