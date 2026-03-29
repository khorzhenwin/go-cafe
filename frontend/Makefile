.PHONY: run-frontend-helper teardown-helper up down

run-frontend-helper:
	@rm -rf .next/dev .next/server
	exec npm run dev

teardown-helper:
	@rm -rf .next
	@echo "Frontend teardown complete: removed .next/"

up:
	docker compose up -d --build

down:
	docker compose down
