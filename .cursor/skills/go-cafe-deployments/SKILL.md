---
name: go-cafe-deployments
description: Build, tag, push, and apply go-cafe frontend deployment artifacts. Use when the user asks to create or update Dockerfiles, docker compose files, Kubernetes manifests, deployment Makefile commands, or to build/push/deploy the frontend image to the local registry or cluster.
---

# Go Cafe Deployments

## Use this skill when

- The user asks to containerize or deploy the frontend.
- The user asks to build, tag, push, or apply a frontend image.
- The user wants deployment steps automated in a `Makefile`.
- The user asks to update `frontend/Dockerfile`, `frontend/docker-compose.yml`, or `frontend/k8s/frontend.yaml`.

## Current project conventions

- Frontend image builds from `frontend/`.
- Frontend manifest lives at `frontend/k8s/frontend.yaml`.
- Default registry is `192.168.100.8:5000`.
- Default frontend image name is `go-cafe-frontend`.
- Default image tag is `latest`.
- Frontend images should be built for `linux/amd64`.
- The frontend manifest uses `imagePullPolicy: Always` so repeated deploys with `:latest` will pull the newest image.

## Preferred workflow

Prefer root `Makefile` commands over ad hoc shell commands:

```bash
make docker-build-frontend
make docker-push-frontend
make deploy-frontend
```

Common overrides:

```bash
make docker-build-frontend IMAGE_TAG=v1.2.3
make docker-push-frontend REGISTRY=registry.example.com IMAGE_TAG=v1.2.3
make deploy-frontend K8S_NAMESPACE=default IMAGE_TAG=v1.2.3
```

## What each command does

- `make docker-build-frontend`: builds `frontend/` with `docker build --platform linux/amd64`.
- `make docker-push-frontend`: builds and pushes the tagged frontend image.
- `make deploy-frontend`: builds, pushes, applies `frontend/k8s/frontend.yaml`, restarts the deployment, and waits for rollout.

## Editing guidance

- Keep image names and registry values aligned between `Makefile` defaults and `frontend/k8s/frontend.yaml`.
- If the deployment flow changes, update the `Makefile` first and then sync the skill instructions.
- When using `:latest`, keep `imagePullPolicy: Always` or use a unique tag/digest.
- Keep `API_BASE_URL` in the frontend manifest pointed at the in-cluster backend service.

## Validation

After changing deployment files:

1. Run `npm run build` in `frontend` if the Dockerfile or Next.js config changed.
2. Run `make help` to verify the new targets are visible.
3. If the user asked for an actual deployment, use the root `Makefile` targets instead of repeating long Docker and `kubectl` commands manually.
