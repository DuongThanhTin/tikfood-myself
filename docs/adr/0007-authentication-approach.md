# ADR-0007: Authentication approach (JWT access + rotating refresh, Google SSO)

- **Status:** Accepted
- **Date:** 2026-07-06
- **Deciders:** Repository owner (auth, migrations, and Google external-call gates cleared in [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md))
- **Type:** decision

## Context

TikFood is a public, read-only discovery product today: discovery endpoints are
anonymous and there is an inert `.authCta` in the web UI. To let users sign in (and
later save places), the API and web app need authentication that fits the existing
`Handler → Service → Repository` layering ([ADR-0002](0002-handler-service-repository-layering.md))
and the `{data,error}` envelope ([ADR-0003](0003-data-error-response-envelope.md)),
without breaking any existing discovery contract.

The full design lives in [`docs/features/authentication/`](../features/authentication/overview.md);
this ADR records the load-bearing decisions (D1–D8) that constrain the implementation.

## Decision

We add email/password + Google OAuth authentication built on a **stateless short-lived
access JWT** paired with a **rotating refresh token persisted in Postgres**.

- **D1 — Token model:** access JWT (~15 min, HS256) sent as a `Bearer` header; a
  rotating refresh token stored in the DB is the revocable "session".
- **D2 — Token storage/transport:** refresh token in an httpOnly `Secure`
  `SameSite=Lax` cookie scoped to `/api/v1/auth`; access token kept in memory on the
  client. `COOKIE_SECURE=false` is allowed for local http only.
- **D3 — Pages:** dedicated `/login` and `/register` pages (+ `/forgot-password`
  placeholder).
- **D4 — Validation:** controlled-state validation helper, no form library / no new
  runtime dependency.
- **D5 — Google flow:** OAuth **Authorization Code** flow with a server-side callback,
  keeping the client secret server-side only.
- **D6 — Domain placement:** a new `internal/auth/` domain with an in-memory fallback
  repository, mirroring `internal/discovery/`.
- **D7 — Routing:** auth endpoints under `/api/v1/auth/*`; discovery endpoints stay
  unchanged and public.
- **D8 — CORS:** credentialed CORS middleware with an explicit allowed-origin list from
  config (never `*` with credentials).

Config keys (env, never hardcoded, never logged): `JWT_SECRET`, `ACCESS_TOKEN_TTL`
(default `15m`), `REFRESH_TOKEN_TTL` (default `720h`), `GOOGLE_CLIENT_ID`,
`GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`, `WEB_ORIGIN` (default
`http://localhost:3000`), `ALLOWED_ORIGINS`, `COOKIE_SECURE` (default `true`).

New Go dependencies: `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`,
`golang.org/x/oauth2` (+ `google` endpoint).

## Consequences

- **Positive:** logout/revocation works (server-side refresh store); XSS token theft
  mitigated (httpOnly cookie + in-memory access token); additive and non-breaking to
  discovery; unit-testable via in-memory repos.
- **Negative / cost:** two new tables + migrations; rotation/revocation logic to get
  right; CORS/cookie flags are environment-sensitive; a JWT secret must be managed
  (prod secret store is infra/human-gated, out of scope here).
- **Follow-ups:** implemented across milestones M0–M12
  ([`docs/features/authentication/tasks.md`](../features/authentication/tasks.md)); this
  ADR is updated to the final endpoint list in M12. Forgot-password is a placeholder
  only (no email delivery in scope).

## Alternatives considered

- **Stateless JWT only (no refresh store)** — rejected: no server-side revocation, so
  logout can't invalidate a session.
- **Opaque server-session cookie only (no JWT)** — rejected: the design explicitly
  calls for JWT auth and stateless access checks for API calls.
- **Client-side Google implicit/token flow** — rejected: would expose or complicate
  client-secret handling; Authorization Code with a server callback is standard.
- **Form library (React Hook Form / Zod)** — rejected: adds a runtime dependency the
  frontend standard defers (D4).

## Sources / links

- Design package: [`docs/features/authentication/`](../features/authentication/overview.md)
  (overview, api, database, security, frontend, ui, testing, release, tasks).
- Plan: [`docs/superpowers/plans/2026-07-05-authentication-system.md`](../superpowers/plans/2026-07-05-authentication-system.md).
- Layering [ADR-0002](0002-handler-service-repository-layering.md); envelope [ADR-0003](0003-data-error-response-envelope.md); in-memory fallback [ADR-0004](0004-in-memory-fallback-repository.md).
- Approval gates: [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md).
