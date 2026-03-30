SHELL := /bin/bash

REGISTRY ?= 192.168.100.8:5000
IMAGE_TAG ?= latest
FRONTEND_PLATFORM ?= linux/amd64
FRONTEND_IMAGE ?= $(REGISTRY)/go-cafe-frontend:$(IMAGE_TAG)
FRONTEND_K8S_MANIFEST ?= frontend/k8s/frontend.yaml
FRONTEND_DEPLOYMENT_NAME ?= go-cafe-frontend
K8S_NAMESPACE ?=
KUBECTL_NAMESPACE_ARGS := $(if $(K8S_NAMESPACE),-n $(K8S_NAMESPACE),)

.PHONY: help run run-backend run-frontend teardown docker-check k8s-check docker-build-frontend docker-push-frontend deploy-frontend

help:
	@echo "go-cafe root targets:"
	@echo "  make run           - Run backend + frontend together for local testing"
	@echo "  make run-backend   - Run backend only"
	@echo "  make run-frontend  - Run frontend only"
	@echo "  make docker-build-frontend - Build the frontend image for $(FRONTEND_PLATFORM)"
	@echo "  make docker-push-frontend  - Build and push the frontend image to $(REGISTRY)"
	@echo "  make deploy-frontend       - Build, push, and apply the frontend Kubernetes manifest"
	@echo "  make teardown      - Stop all services and clean backend/frontend artifacts"
	@echo ""
	@echo "Deployment overrides:"
	@echo "  make deploy-frontend IMAGE_TAG=v1.2.3"
	@echo "  make deploy-frontend REGISTRY=registry.example.com K8S_NAMESPACE=default"

run: up

up:
	@echo "Starting backend (:8080) and frontend (:3000)..."
	@$(MAKE) -C backend run & BACK_PID=$$!; \
	for i in {1..60}; do \
		if (echo > /dev/tcp/127.0.0.1/8080) >/dev/null 2>&1; then \
			break; \
		fi; \
		if ! kill -0 $$BACK_PID 2>/dev/null; then \
			echo "Backend process exited before becoming ready."; \
			wait $$BACK_PID; \
			exit $$?; \
		fi; \
		sleep 0.5; \
	done; \
	if ! (echo > /dev/tcp/127.0.0.1/8080) >/dev/null 2>&1; then \
		echo "Timed out waiting for backend on :8080"; \
		kill $$BACK_PID 2>/dev/null || true; \
		wait $$BACK_PID 2>/dev/null || true; \
		exit 1; \
	fi; \
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

k8s-check:
	@if ! command -v kubectl >/dev/null 2>&1; then \
		echo "kubectl not found. Install kubectl and try again."; \
		exit 1; \
	fi
	@kubectl version --client >/dev/null 2>&1 || { \
		echo "kubectl is installed but not working correctly."; \
		exit 1; \
	}
	@echo "kubectl is available."

docker-build-frontend: docker-check
	@echo "Building $(FRONTEND_IMAGE) for $(FRONTEND_PLATFORM)..."
	@docker build --platform $(FRONTEND_PLATFORM) -t $(FRONTEND_IMAGE) frontend

docker-push-frontend: docker-build-frontend
	@echo "Pushing $(FRONTEND_IMAGE)..."
	@docker push $(FRONTEND_IMAGE)

deploy-frontend: docker-push-frontend k8s-check
	@echo "Applying $(FRONTEND_K8S_MANIFEST)..."
	@kubectl apply $(KUBECTL_NAMESPACE_ARGS) -f $(FRONTEND_K8S_MANIFEST)
	@kubectl rollout restart deployment/$(FRONTEND_DEPLOYMENT_NAME) $(KUBECTL_NAMESPACE_ARGS)
	@kubectl rollout status deployment/$(FRONTEND_DEPLOYMENT_NAME) $(KUBECTL_NAMESPACE_ARGS) --timeout=180s

down:
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

teardown: down
