# Authentication — Overview

One-line purpose: the entry point to the Authentication feature design package — what we are building, why, what already exists, and where every other design document lives.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

`overview.md` orients any reader (human or AI agent) before they open the deeper
documents. It states the feature's intent and boundaries, summarizes the current
state of the repository, indexes the rest of the package, and records the locked
design decisions so every downstream document shares one baseline. It deliberately
stays at the "what and why" altitude — endpoint shapes, schema, and code-level
detail belong in the specialized documents listed below.

This is the single canonical design package for the Authentication feature.

---

## Feature summary

Add authentication to TikFood, which today is a **public, read-only discovery app**
with no user accounts. The feature covers:

- **Backend (Go/Gin):** user registration, login, logout, refresh-token rotation,
  JWT-based authentication, session management, and Google OAuth SSO.
- **Frontend (Next.js):** login page, register page, Google Sign-In, a
  forgot-password placeholder (future), responsive layout, form validation, and
  loading/error states.
- **UI/UX:** reuse the existing design system, improve consistency without changing
  the design language, and reuse existing components wherever possible.

## Goals

- Provide secure email/password and Google SSO authentication.
- Preserve the existing discovery experience: all current endpoints stay **public
  and unchanged**; the app degrades gracefully for anonymous users.
- Reuse the established design language and components — introduce no new visual
  system and no new frontend framework.
- Keep the change additive and reversible, following the repository's layering,
  response-envelope, and security rules.

## Non-goals (and MVP anti-goal alignment)

Authentication is a **foundation**, not a gateway to out-of-scope product areas.
Per the MVP anti-goals, this feature will **not** introduce: delivery, cart, order,
checkout, payment, booking/reservation, in-app chat, a social follow graph, creator
monetization, or livestream. It also does **not** add password-reset email delivery,
email verification mail, or account/profile management screens in this iteration
(forgot-password is a UI placeholder only).

→ See [`docs/tikfood/anti-goals.md`](../../tikfood/anti-goals.md) and
[`docs/tikfood/mvp-scope.md`](../../tikfood/mvp-scope.md).

## Scope boundary (Existing vs Proposed)

### Existing implementation

- **Backend** `apps/api`: Go/Gin discovery API; layering
  `Handler → Service → Repository → Postgres`; `{data,error}` response envelope;
  middleware chain `recovery → requestID → requestLogging`. **No auth code, no user
  tables, no JWT, no sessions.** An inert `.authCta` ("Đăng nhập ngay") placeholder
  button exists in the web UI.
- **Frontend** `apps/web`: Next.js App Router, plain CSS design system, a single `/`
  route, all fetches via `lib/api.ts`. **No auth UI, no auth routes, no session
  state.**
- **Governance**: auth, database migrations, and new external network calls
  (Google) are **human-approval-gated** — approval was granted during Discovery.

### Proposed changes (high level)

- A new `internal/auth` backend domain, `/api/v1/auth/*` endpoints, auth + CORS
  middleware, two additive migrations (users, refresh tokens).
- A frontend `AuthProvider`, `lib/auth.ts`, `/login` `/register` `/forgot-password`
  pages, a Google Sign-In flow, and reuse of existing design tokens/components.
- Detailed designs live in the per-topic documents indexed below.

## Locked design decisions (baseline for the whole package)

| # | Decision |
|---|----------|
| D1 | Short-lived access JWT (~15 min, HS256, `Bearer`) + rotating refresh token persisted in Postgres = session management |
| D2 | Refresh token in an httpOnly `Secure` `SameSite=Lax` cookie; access token held in memory by the client |
| D3 | Dedicated `/login`, `/register` pages (+ `/forgot-password` placeholder) — not modal-first |
| D4 | Controlled-state validation helper on the frontend; **no** new form library |
| D5 | Google OAuth via server-side Authorization-Code flow |
| D6 | New `internal/auth` domain with an in-memory fallback repository (mirrors `internal/discovery`) |
| D7 | Auth endpoints under `/api/v1/auth/*`; discovery endpoints unchanged and public |
| D8 | Credentialed CORS with an explicit allowed-origin list from config |

## Package index (generated one at a time, each gated by approval)

| # | Document | Purpose |
|---|----------|---------|
| 1 | `overview.md` (this) | Feature intent, scope, decisions, index |
| 2 | [`specification.md`](specification.md) | Functional & non-functional requirements + acceptance criteria |
| 3 | [`architecture.md`](architecture.md) | Backend/frontend architecture, token lifecycle, flows |
| 4 | [`api.md`](api.md) | HTTP contract: endpoints, requests/responses, errors, cookies |
| 5 | [`database.md`](database.md) | Schema changes: users & refresh-token tables, migrations |
| 6 | [`security.md`](security.md) | Threat model and controls |
| 7 | [`frontend.md`](frontend.md) | Web architecture: provider, client, routes, state |
| 8 | [`ui.md`](ui.md) | Screens, wireframes, states, accessibility, copy |
| 9 | [`component-reuse.md`](component-reuse.md) | Reuse-vs-create audit against the design system |
| 10 | [`testing.md`](testing.md) | Test strategy and coverage matrix |
| 11 | [`release.md`](release.md) | Rollout, migrations, config, rollback |
| 12 | [`tasks.md`](tasks.md) | Milestones, PR boundaries, dependencies, verification, risks |

## Success criteria

- A user can register, log in (email/password or Google), stay signed in across a
  reload, refresh silently, and log out.
- Discovery continues to work unchanged for anonymous users.
- The auth UI is visually indistinguishable in language from the existing app and
  reuses existing components; no new design system is introduced.
- Every backend/frontend verification gate passes (build, vet/typecheck, tests).
- No secrets are logged; no breaking API changes.

## Open questions

- **OQ-1 (scope trigger):** What is the first authenticated user-scoped feature that
  justifies auth now, and does it stay clear of the anti-goals? Affects how much
  ships first. *(Default: build the foundation; discovery stays public.)*
- **OQ-2 (token lifetimes):** Confirm ~15 min access / ~30 day refresh.
- **OQ-3 (CI gating):** No CI exists today; add a minimal gating workflow (separate,
  gated PR) or rely on local per-PR verification? See [`release.md`](release.md).
- **OQ-4 (production secrets):** Who owns the production secret store for
  `JWT_SECRET` and Google credentials (out of scope here, infra/human-gated)?

## Related repository documentation

- [`/CLAUDE.md`](../../../CLAUDE.md) — workspace rules, anti-goals, security/git, naming.
- [`docs/ai/AI-CONTRACT.md`](../../ai/AI-CONTRACT.md) — rules of engagement; §1 anti-goals, §4 human-approval gates (auth/migrations/external calls).
- [`docs/ai/CONTEXT-LOADING.md`](../../ai/CONTEXT-LOADING.md) — what to read per task type.
- [`docs/REPOSITORY-MAP.md`](../../REPOSITORY-MAP.md) — directory → purpose map.
- [`docs/architecture.md`](../../architecture.md) — system architecture.
- [`docs/tikfood/`](../../tikfood/) — product vision, MVP scope, anti-goals.
- [`docs/design/`](../../design/) — canonical `apps/web` design system.
- [`docs/adr/`](../../adr/) — ADR-0002 (layering), ADR-0003 (envelope); a new ADR-0007 (authentication approach) is proposed.
- [`docs/standards/`](../../standards/) — backend & frontend engineering standards, `application-security.md`.
- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — completion gates.
- [`apps/api/CLAUDE.md`](../../../apps/api/CLAUDE.md) · [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — per-app rules.
- Milestone plan (companion): `docs/superpowers/plans/2026-07-05-authentication-system.md` — mirrored in this package's [`tasks.md`](tasks.md).
