# Curated CLAUDE.md Guidance + Improvement Roadmap Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace two generic "ideal enterprise" CLAUDE.md drafts with reality-based AI-agent guidance (one root + two per-app files) plus a prioritized improvement roadmap.

**Architecture:** Four standalone Markdown documents. `/CLAUDE.md` holds shared rules and points to per-app files; `apps/api/CLAUDE.md` and `apps/web/CLAUDE.md` describe each app exactly as it is built today; `IMPROVEMENTS.md` holds the P0/P1/P2 review findings. No code changes, no dependency changes.

**Tech Stack:** Markdown only. The documents describe an existing Go (Gin + database/sql + slog) backend and a Next.js (App Router + MapLibre GL + plain CSS) frontend.

## Global Constraints

- Describe reality, not aspiration. The response envelope is `{ "data": ..., "error": { "code", "message", "details?" } }` — never document `{success,...}`, `meta`, or `trace_id` as if they exist today.
- Do NOT modify any file under `docs/standards/`.
- Do NOT add or recommend installing GORM, React Query, React Hook Form, or Zod as current rules.
- Backend layer names are `Handler` / `Service` / `Repository` (not "Controller").
- "Future direction (approved)" items are limited to: split `DiscoveryExperience.tsx`, add `trace_id` to the response envelope, error-code catalog.
- Source spec: `docs/superpowers/specs/2026-06-26-claude-md-curation-design.md`.
- Work happens on the `dev` branch.

---

### Task 1: Root `/CLAUDE.md`

**Files:**
- Create: `/CLAUDE.md`

**Interfaces:**
- Produces: the shared rule set and the "code wins over stale docs/standards" rule that the two per-app files assume.

- [ ] **Step 1: Create `/CLAUDE.md` with this exact content**

```markdown
# CLAUDE.md — TikFood Workspace

TikFood is realtime social food discovery (TikTok + Google Maps for food):
dish-first, map-first, social-proof-driven, trend-scored, geo-aware,
AI-summary-assisted. Discovery only.

## MVP Anti-Goals (block or require human approval)

Do not implement: delivery, cart, order, checkout, payment, booking,
reservation, in-app chat, social follow graph, creator monetization, livestream.

## Monorepo Map

- `apps/api` — Go backend (discovery API). See `apps/api/CLAUDE.md`.
- `apps/web` — Next.js frontend (discovery UX). See `apps/web/CLAUDE.md`.
- `apps/ai-code-runner` — automation runner skeleton (n8n → runner → PR).
- `packages/` — prompts, schemas, config policies.
- `docs/standards/` — engineering standards (background reference).
- `IMPROVEMENTS.md` — prioritized improvement roadmap.

## Source of Truth

When `docs/standards/*` conflicts with the actual code, the code wins. Flag the
mismatch instead of silently following stale docs. Read the relevant
`apps/*/CLAUDE.md` before changing that app.

## Security & Git

- Never read `.env`, `.env.*`, secrets, credentials, private keys, or tokens.
- Never push to `main` or `master`. Never force-push. Never auto-merge.
- Create PR branches under `ai/`. Active working branch is `dev`.
- Require human approval for auth, migrations, infra, external network calls,
  and any MVP anti-goal area.

## Naming (shared)

Use clear verb names: `CreateUser`, `FindByEmail`, `VenueDetail`.
Avoid `Do`, `Handle`, `Process`, `ExecuteX`, `Component1`.

## Before Finishing (shared)

- Reuse existing code; search before adding a new helper/service/component.
- No breaking API changes. Keep the response format unchanged.
- Update tests when you change behavior.
```

- [ ] **Step 2: Verify accuracy against the repo**

Run: `git ls-files | grep -E '^apps/(api|web|ai-code-runner)/' | head` and confirm the three app paths exist.
Run: `grep -n "ai/" README.md` and confirm the `ai/` branch + security rules match what is written.
Expected: app paths exist; security rules match README.

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "Add root CLAUDE.md with shared TikFood agent rules"
```

---

### Task 2: `apps/api/CLAUDE.md`

**Files:**
- Create: `apps/api/CLAUDE.md`

**Interfaces:**
- Consumes: the "Source of Truth" and shared rules from `/CLAUDE.md` (Task 1).
- Produces: backend-specific rules referenced by reviewers of `apps/api` changes.

- [ ] **Step 1: Create `apps/api/CLAUDE.md` with this exact content**

```markdown
# CLAUDE.md — apps/api (Go Backend)

Read `/CLAUDE.md` first. This file describes the backend as it is built today.

## Stack Reality

- Go + **Gin**.
- **`database/sql` via the pgx stdlib driver** — NOT GORM.
- **slog** JSON logging; manual dependency injection via `internal/app/Container`.
- PostGIS-backed Postgres; an in-memory fallback repository runs when
  `DATABASE_URL` is unset.

## Layering

`http (Handler) → discovery (Service) → Repository → postgres`. Dependencies
point one way only.

Forbidden: Handler → repository, Handler → SQL, Service → `gin.Context`,
Repository → business logic.

## Handler (`internal/http`)

Handlers are named `*Handler` (e.g. `VenueHandler`) — not "Controller". They:
parse and validate the request (`venue_request.go`), call the service, and
respond with the `respondWith*` helpers. Handlers never contain business logic,
touch the database, or manage transactions.

## Service (`internal/discovery`)

Services hold all business logic and input normalization and call repositories.
Services never import or reference HTTP / `gin.Context`.

## Repository

`VenueRepository` is an interface with two implementations: `postgres`
(`storage/postgres`) and the in-memory fallback (`discovery`). Repositories do
persistence only — no business rules, no HTTP.

## Response Envelope (real)

`{ "data": ..., "error": { "code", "message", "details?" } }`.
Use `respondWithData`, `respondWithError`, `respondWithBadRequest`,
`respondWithInternalServerError`. Do NOT invent a `{ "success": ... }` format,
and do not add `meta` or `trace_id` fields that do not exist yet.

## Error Handling

Use the predefined codes in `internal/http/errors.go` (`invalid_request`,
`not_found`, `internal_error`). Map domain errors (e.g. `ErrVenueNotFound`) to
HTTP status. Never leak SQL or raw error strings to the client.

## Validation

Query/DTO validation lives in the handler (`venue_request.go`). Business rules
live in the service.

## Config & DI

Configuration comes from `internal/config` (env). Never hardcode URLs, secrets,
or credentials. Use constructor injection through `Container`; no globals.

## Testing

Table-driven tests; follow `internal/http/router_test.go`. Mock
`VenueRepository` for service-level tests.

## Future Direction (approved, do not do preemptively)

- Add `trace_id` to the response envelope (`requestIDMiddleware` already exists).
- Consolidate error codes into a documented catalog as they grow.
- Transactions live in the Service layer once write endpoints exist (discovery
  is read-only today).
```

- [ ] **Step 2: Verify accuracy against the code**

Run: `grep -rn "gorm" apps/api/go.mod apps/api/go.sum || echo "no gorm (correct)"`
Expected: prints "no gorm (correct)".
Run: `grep -n "func respondWith" apps/api/internal/http/response.go`
Expected: shows `respondWithData`, `respondWithError`, `respondWithBadRequest`, `respondWithInternalServerError`.
Run: `grep -n "ErrVenueNotFound\|invalid_request\|not_found\|internal_error" apps/api/internal/http/errors.go apps/api/internal/discovery/model.go`
Expected: confirms the error codes and `ErrVenueNotFound` exist.

- [ ] **Step 3: Commit**

```bash
git add apps/api/CLAUDE.md
git commit -m "Add apps/api CLAUDE.md describing the Go backend as built"
```

---

### Task 3: `apps/web/CLAUDE.md`

**Files:**
- Create: `apps/web/CLAUDE.md`

**Interfaces:**
- Consumes: the "Source of Truth" and shared rules from `/CLAUDE.md` (Task 1).
- Produces: frontend-specific rules referenced by reviewers of `apps/web` changes.

- [ ] **Step 1: Create `apps/web/CLAUDE.md` with this exact content**

```markdown
# CLAUDE.md — apps/web (Next.js Frontend)

Read `/CLAUDE.md` first. This file describes the frontend as it is built today.

## Stack Reality

- Next.js **App Router** + TypeScript.
- **MapLibre GL** for the map (lazy-imported).
- **Plain CSS** in `app/globals.css` — no CSS framework, no design-system library.
- Data through `lib/api.ts` `fetch` wrappers.
- No React Query, no React Hook Form, no Zod today.

## Structure

`app/` (layout, page) → `components/` (`DiscoveryExperience.tsx`,
`VenueList.tsx`) → `lib/api.ts`. The server `page.tsx` passes `initialVenues`
to the client `DiscoveryExperience`.

## API Access

Always go through `lib/api.ts` (`fetchDiscoveryVenues`, `fetchVenueDetail`,
`getDiscoveryVenues`). Do not scatter raw `fetch` calls in components. Respect
the `{ data, error }` envelope. When `NEXT_PUBLIC_API_URL` is unset, fall back
to `fallbackVenues`.

## Types

Keep frontend types in `lib/api.ts` in sync with the Go JSON tags
(snake_case field names such as `avg_price_min_vnd`, `social_videos`).

## Styling

Use plain CSS classes in `globals.css` and reuse existing class patterns. Inline
`style` only for computed positions (e.g. map markers). No new CSS framework
without justification.

## Accessibility & Copy

Keep `aria-label` on icon-only buttons (existing pattern) and label inputs.
Vietnamese UI copy is the norm.

## Naming

Descriptive component names (`VenueDetail`, `VenueMap`, `VenueRailCard`). Never
`Component1` / `DataCard2`.

## Future Direction (approved, do not do preemptively)

- Split `components/DiscoveryExperience.tsx` (~1350 lines) into focused pieces:
  `VenueMap`, `VenueDetail`, `VenueRailCard`, filter controls, map-style helpers,
  and format utilities, under `components/`.
- Introduce React Query / React Hook Form + Zod ONLY when real
  caching/mutations/forms appear.
```

- [ ] **Step 2: Verify accuracy against the code**

Run: `grep -n "react-query\|react-hook-form\|zod\|tailwind" apps/web/package.json || echo "none present (correct)"`
Expected: prints "none present (correct)".
Run: `grep -n "export async function fetch\|export async function get\|fallbackVenues" apps/web/lib/api.ts | head`
Expected: shows `getDiscoveryVenues`, `fetchDiscoveryVenues`, `fetchVenueDetail`, `fallbackVenues`.
Run: `wc -l apps/web/components/DiscoveryExperience.tsx`
Expected: ~1350 lines (confirms the split-justification claim).

- [ ] **Step 3: Commit**

```bash
git add apps/web/CLAUDE.md
git commit -m "Add apps/web CLAUDE.md describing the Next.js frontend as built"
```

---

### Task 4: `IMPROVEMENTS.md`

**Files:**
- Create: `IMPROVEMENTS.md`

**Interfaces:**
- Consumes: referenced from `/CLAUDE.md` monorepo map (Task 1).
- Produces: the standalone prioritized roadmap. No other file depends on it.

- [ ] **Step 1: Create `IMPROVEMENTS.md` with this exact content**

```markdown
# IMPROVEMENTS.md — TikFood Roadmap

Prioritized review findings. This is a roadmap, not a set of mandatory rules.
CLAUDE.md only references the relevant constraints (e.g. do not show fabricated
ratings).

## P0 — Product Credibility

### Fabricated UI data
`apps/web/components/DiscoveryExperience.tsx` hardcodes `mediaByVenue` with fake
ratings ("4.8"), review counts ("128 reviews"), badges, and expiring
`lh3.googleusercontent.com/aida-public/...` image URLs keyed to seed UUIDs. The
data model has no rating/review fields, and `venues.photos[]` exists in the
schema but the UI ignores it. Decision needed: either surface real `photos` and
add rating fields end to end, or remove the rating/review UI until real data
exists. Until then, fake data must be clearly marked as placeholder.

### Trend score & AI summary are static seed values
`trend_scores` and `ai_summaries` tables exist (migration 001) but no worker
populates them. The core differentiation — realtime trend scoring and AI "why
trending" summaries — does not exist yet. Roadmap: social ingestion →
trend-scoring worker → AI-summary worker.

## P1 — Engineering Robustness

### No CI
There is no `.github/workflows`. Add GitHub Actions: `go test` / `go vet` /
`go build` for `apps/api`, and `typecheck` / `lint` / `build` for `apps/web`.
This enforces the "tests updated" rule.

### Thin test coverage
Only `apps/api/internal/http/router_test.go` exists. Add table-driven tests for
`normalizeSearch` and alias matching (`apps/api/internal/discovery/search.go`)
and for the in-memory fallback repository filter. The long list SQL is a
regression risk.

### Location-alias logic duplicated in three places / two paradigms
`apps/api/internal/discovery/search.go` hardcodes `quan-1`/`quan-3`; the DB has
`locations` / `location_aliases` (migration 003) used by the SQL query; and
`apps/web/lib/api.ts` re-implements `normalizeLocationAlias`. Consolidate on the
DB as the source of truth; mark in-memory/fallback alias logic as dev-only.

## P2 — API / Scaling Polish

### No pagination metadata
The list endpoint returns a bare array in `data` with only `limit`. This is
where a `meta` field (total / cursor) is legitimately useful — reserve it for
pagination. (Future.)

### OSRM public demo server
`apps/web/components/DiscoveryExperience.tsx` calls `router.project-osrm.org`
directly from the browser — rate-limited, no SLA, not production-safe. Replace
with self-hosted or managed routing.

### Shallow health check
`GET /health` returns `{ "ok": true }` without pinging the DB. Add a readiness
check that verifies the database before orchestration relies on it.

### Repo mixes automation tooling and product app
`apps/ai-code-runner` and the n8n workflow docs live beside the product app.
`starters/tikfood` hints the product graduates to its own repo. A conscious
decision for later.
```

- [ ] **Step 2: Verify accuracy against the code**

Run: `ls .github/workflows 2>/dev/null || echo "no CI (correct)"`
Expected: prints "no CI (correct)".
Run: `grep -rn "router.project-osrm.org" apps/web/components/DiscoveryExperience.tsx`
Expected: confirms the OSRM call exists.
Run: `grep -n "ok.*true" apps/api/internal/http/router.go`
Expected: confirms the shallow health check.

- [ ] **Step 3: Commit**

```bash
git add IMPROVEMENTS.md
git commit -m "Add IMPROVEMENTS.md prioritized roadmap"
```

---

## Self-Review (completed by plan author)

- **Spec coverage:** Deliverable 1 → Task 1; Deliverable 2 → Task 2; Deliverable 3
  → Task 3; Deliverable 4 → Task 4. The "describe reality" philosophy, future-
  direction list, and "docs/standards = reference only" rule are all reflected.
- **Placeholder scan:** No TBD/TODO; every file body is complete and embedded.
- **Type/name consistency:** Helper names (`respondWith*`), error codes, and
  `lib/api.ts` function names are used identically across tasks and match the code
  verified during brainstorming.
- **Out of scope honored:** no `docs/standards/` edits; no GORM/React Query/RHF/Zod;
  envelope unchanged.
```
