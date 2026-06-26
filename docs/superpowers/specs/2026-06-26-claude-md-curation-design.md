# Design: AI Agent Guidance (CLAUDE.md) + Improvement Roadmap

Date: 2026-06-26
Status: Approved (pending spec review)

## Goal

The user has two generic "ideal enterprise" CLAUDE.md drafts (one Go backend, one
Next.js frontend). Much of their content does not match the actual stack and would
make AI agents fight the real code. This work curates those drafts into guidance
that **describes the project as it actually is today**, plus a separate prioritized
improvement roadmap.

## Key decisions (from brainstorming)

1. **Philosophy: describe reality.** CLAUDE.md documents what the code actually does
   now. A small "Future direction (approved)" section captures a few high-value
   upgrades — without forcing them to be done now.
2. **Structure: one root + two per-app files.** `/CLAUDE.md` (shared, thin),
   `apps/api/CLAUDE.md` (Go backend), `apps/web/CLAUDE.md` (Next.js frontend).
   Claude Code loads the file matching the directory being worked in.
3. **Future-direction items to keep:** split `DiscoveryExperience.tsx`, add
   `trace_id` to the response envelope, error-code catalog. **Dropped from the
   drafts:** React Query, React Hook Form, Zod, GORM, the `{success,...}` envelope,
   "Controller" naming, Design System library, apiClient-as-mandatory — none of
   these exist today and YAGNI applies.
4. **docs/standards: reference only, do not modify** in this work. CLAUDE.md is the
   primary source of truth and includes the rule: *when docs/standards conflicts
   with the actual code, the code wins — flag the mismatch, do not silently follow
   stale docs.*
5. **Review findings go into a separate `IMPROVEMENTS.md`** (prioritized P0/P1/P2),
   not mixed into CLAUDE.md.

## Stack reality (verified against code)

- **Backend (`apps/api`):** Go + **Gin**, **`database/sql` via pgx stdlib (NOT
  GORM)**, **slog** JSON logging, manual DI via `Container`. Layering:
  `http (Handler) → discovery (Service) → Repository → postgres`. Handlers are named
  `*Handler` (e.g. `VenueHandler`), not "Controller".
- **Response envelope (real):** `{ "data": ..., "error": { "code", "message",
  "details?" } }` via `respondWithData/respondWithError/respondWithBadRequest/
  respondWithInternalServerError`. There is **no** `success`, `meta`, or `trace_id`
  field today.
- **Error codes (real):** `invalid_request`, `not_found`, `internal_error` in
  `errors.go`; domain error `ErrVenueNotFound` mapped to HTTP 404.
- **Data layer is strong:** PostGIS `geography(Point,4326)` + GiST index, normalized
  schema, `locations`/`location_aliases` identity tables, `trend_scores`/
  `ai_summaries` tables, `unaccent`+`pg_trgm`. Filtering is pushed down to SQL.
- **Frontend (`apps/web`):** Next.js App Router + TypeScript, **MapLibre GL**, plain
  CSS in `globals.css`, data through `lib/api.ts` `fetch` wrappers. **No** React
  Query / RHF / Zod / design-system library. One large client component
  `DiscoveryExperience.tsx` (~1350 lines). Server `page.tsx` passes `initialVenues`
  to the client component. Routing uses the public OSRM demo server from the browser.

---

## Deliverable 1 — `/CLAUDE.md` (root, thin, shared)

Sections:

- **What TikFood is** (1–2 lines) and **MVP anti-goals** (block: delivery, cart,
  order, checkout, payment, booking, reservation, in-app chat, social follow graph,
  creator monetization, livestream).
- **Monorepo map** + pointer: read `apps/*/CLAUDE.md` for per-app detail;
  `docs/standards/*` is background.
- **Source-of-truth rule:** when `docs/standards` conflicts with the actual code,
  the code wins; flag the mismatch instead of silently following stale docs.
- **Security & Git** (from README): never read `.env`/secrets; never push
  `main`/`master`; never force-push; never auto-merge; PR branches under `ai/`;
  human approval for auth/migration/infra/external-network/anti-goal changes.
  Current active working branch is `dev`.
- **Naming (shared):** clear verb names (`CreateUser`, `FindByEmail`); avoid
  `Do`/`Handle`/`Process`/`ExecuteX`.
- **Before finishing (shared):** reuse existing code; no breaking API changes;
  update tests; keep the response format unchanged.

## Deliverable 2 — `apps/api/CLAUDE.md` (Go backend, reality)

- **Stack reality:** Gin, `database/sql` via pgx (no GORM), slog JSON, manual DI via
  `Container`.
- **Layering (real names):** `http (Handler) → discovery (Service) → Repository →
  postgres`. One-way dependency. Forbidden: Handler→repo, Handler→SQL,
  Service→`gin.Context`, Repository→business logic.
- **Handler:** parse + validate request (`venue_request.go`), call service, respond
  via `respondWith*`. Never business logic / DB / transactions.
- **Service** (`discovery/service.go`): all business logic, input normalization,
  calls repository. Never touches HTTP / `gin.Context`.
- **Repository:** `VenueRepository` interface with two implementations (postgres +
  in-memory fallback). Persistence only.
- **Response envelope (real):** `{ data, error{code,message,details?} }`. Use the
  existing `respondWith*` helpers. Do NOT invent a `{success,...}` format.
- **Error handling:** use predefined codes in `errors.go`; map domain errors
  (`ErrVenueNotFound`) to HTTP; never leak SQL or raw error strings.
- **Validation:** query/DTO validation in the handler (`venue_request.go`); business
  rules in the service.
- **Config:** via `internal/config` (env); never hardcode URLs/secrets.
- **DI:** constructor injection through `Container`; no globals.
- **Testing:** table-driven; follow `router_test.go`; mock `VenueRepository` for
  service tests.
- **Future direction (approved):** (a) add `trace_id` to the response envelope —
  `requestIDMiddleware` already exists; (b) consolidate error codes into a documented
  catalog as they grow; (c) transactions live in the Service layer once write
  endpoints exist (discovery is read-only today).

## Deliverable 3 — `apps/web/CLAUDE.md` (Next.js frontend, reality)

- **Stack reality:** Next.js App Router + TS, MapLibre GL, plain CSS (`globals.css`),
  data via `lib/api.ts`. No React Query / RHF / Zod / design-system library today.
- **Current structure:** `app/` (layout, page) → `components/`
  (`DiscoveryExperience.tsx`, `VenueList.tsx`) → `lib/api.ts`. Server `page.tsx`
  passes `initialVenues` to the client `DiscoveryExperience`.
- **API access:** always go through `lib/api.ts` functions
  (`fetchDiscoveryVenues`, `fetchVenueDetail`, `getDiscoveryVenues`); don't scatter
  raw `fetch` in components. Respect the `{data,error}` envelope; fall back to
  `fallbackVenues` when `NEXT_PUBLIC_API_URL` is unset.
- **Types:** keep frontend types in `lib/api.ts` in sync with the Go JSON tags
  (snake_case fields).
- **Styling:** plain CSS classes in `globals.css`; reuse existing class patterns;
  inline `style` only for computed positions (map markers). No new CSS framework
  without justification.
- **Accessibility:** keep `aria-label` on icon buttons (existing pattern); label
  inputs. Vietnamese UI copy is the norm.
- **Naming:** descriptive component names (`VenueDetail`, `VenueMap`); never
  `Component1`/`DataCard2`.
- **Future direction (approved):** split `DiscoveryExperience.tsx` into smaller
  pieces (`VenueMap`, `VenueDetail`, `VenueRailCard`, filter controls, map-style
  helpers, format utils) under `components/`; introduce React Query / RHF+Zod ONLY
  when real caching/mutations/forms appear.

## Deliverable 4 — `IMPROVEMENTS.md` (prioritized roadmap)

A standalone doc capturing the architecture/design review. CLAUDE.md only references
the relevant rules (e.g. "do not show fabricated ratings"; "mark fake/seed data
explicitly").

### P0 — Product credibility
- **Fake UI data:** `mediaByVenue` hardcodes ratings ("4.8"), review counts, badges,
  and expiring `lh3.googleusercontent.com/aida-public/...` image URLs keyed to seed
  UUIDs. The data model has no rating/review fields; `venues.photos[]` exists but the
  UI ignores it. Decision: either surface real `photos` + add rating fields end to
  end, or remove rating/review UI until real data exists. Until then, fake data must
  be clearly marked.
- **Trend score & AI summary are static seed values.** `trend_scores`/`ai_summaries`
  tables exist but no worker populates them. The core product differentiation
  (realtime trend scoring, AI "why trending" summaries) does not exist yet. Roadmap:
  ingestion → trend-scoring worker → AI-summary worker.

### P1 — Engineering robustness
- **No CI.** Add GitHub Actions: `go test`/`vet`/`build` for the API; `typecheck`/
  `lint`/`build` for web. Enforces the "tests updated" rule.
- **Thin test coverage.** Only `router_test.go`. Add table-driven tests for
  `normalizeSearch`, alias matching (`search.go`), and the fallback repository
  filter. The long list SQL is a regression risk.
- **Location-alias logic duplicated in three places / two paradigms:** Go
  `search.go` hardcodes `quan-1/quan-3`; the DB has `locations`/`location_aliases`
  (migration 003) used by the SQL query; the frontend re-implements
  `normalizeLocationAlias`. Consolidate on the DB as the source of truth; mark
  in-memory/fallback alias logic as dev-only approximations.

### P2 — API / scaling polish
- **No pagination metadata** on the list endpoint (bare array in `data`, `limit`
  only). This is where the original draft's `meta` field is legitimately useful —
  reserve `meta` for pagination/total. (Future.)
- **OSRM public demo server** is called directly from the browser — rate-limited, no
  SLA, not production-safe. Replace with self-hosted/managed routing.
- **Shallow health check** (`{ok:true}`) does not ping the DB. Add a readiness check.
- **Repo mixes automation tooling and the product app.** `ai-code-runner`/n8n live
  beside the product; `starters/tikfood` hints the product graduates to its own repo.
  A conscious decision for later.

## Out of scope
- Modifying `docs/standards/*`.
- Implementing any of the IMPROVEMENTS items (this work only documents them).
- Adding GORM, React Query, RHF, Zod, or changing the response envelope.

## Testing / validation
- These are documentation files; validation is human review for accuracy against the
  code. No automated tests are added by this work itself (CI is a P1 roadmap item,
  not part of this deliverable).
