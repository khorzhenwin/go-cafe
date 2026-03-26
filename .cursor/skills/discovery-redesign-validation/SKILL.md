---
name: discovery-redesign-validation
description: Validate major UX or fullstack redesign work in go-cafe. Use when redesigning routes, discovery flows, backend contracts, or when the user asks to self validate before calling the task complete.
---

# Discovery Redesign Validation

Use this skill when changes affect product structure, frontend UX, backend contracts, or cross-stack behavior in `go-cafe`.

## Validation standard

Never mark the task complete until all applicable checks pass or a blocker is explicitly surfaced to the user.

## Required workflow

1. Confirm whether the change touched:
   - frontend routes or components
   - frontend API helpers
   - backend handlers, services, repositories, models, or migrations
   - README architecture or contract docs
2. Run the matching validations immediately after implementation:
   - Frontend: `cd frontend && npm run lint` and `cd frontend && npm run build`
   - Backend unit scope: `cd backend && go test ./internal/cafelisting ./internal/rating ./internal/server ./internal/user`
   - Broader backend verification when handlers or models changed: use the closest relevant `go test` packages, and run integration coverage if local DB/env is available
3. If validation fails:
   - fix the issue
   - rerun the same validation
   - do not proceed until it passes or you have a concrete blocker
4. If API routes, schema, or product flows changed, update `README.md` before finishing.
5. Manually verify the redesigned product flow when route UX changed.

## Manual flow checklist

Use this checklist for the current redesign architecture:

- `/` communicates discovery-first intent and links cleanly into the rest of the product
- `/map` loads discovery results, filtering works, and a cafe can be opened from the map/list
- `/cafes/[id]` shows place context and the save action behaves correctly
- `/my-places` supports add, status update, and delete flows
- `/reviews` supports create and delete review flows for visited cafes
- `/auth` supports login/register and session persistence

## Repo-specific guardrails

- Preserve the Next proxy route under `frontend/app/api/backend/[...path]/route.js`
- Keep `frontend/lib/api` organized by domain rather than returning to one giant client file
- Treat `visit_status` as the user-facing saved/visited state unless the task explicitly changes the model
- Do not leave README contract sections stale after backend or flow changes

## Completion rule

Only report completion once:

- validations passed
- any changed docs are synced
- residual risks, if any, are explicitly called out
