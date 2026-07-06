# Authentication — Architecture

One-line purpose: how authentication fits into the existing system — the backend domain, middleware, token lifecycle, and the flows that tie them together.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

`architecture.md` describes the *shape* of the solution: which components exist,
which are added, how requests flow, and how tokens live and die. It stays at the
structural level — request/response field detail belongs in [`api.md`](api.md),
schema in [`database.md`](database.md), and control rationale in
[`security.md`](security.md). It exists so an implementer can build each piece in
isolation while trusting the seams.

---

## Existing implementation

Current system context (from Discovery):

```
 Ingestion/OSM ──► PostgreSQL/PostGIS ──► Go/Gin Discovery API ──► Next.js Web (App Router)
                                              │  public, read-only        │
                                              └──── {data,error} envelope ─┘
```

- **Layering (ADR-0002):** `Handler → Service → Repository → Postgres`, one-way.
  Handlers parse/validate and never touch the DB; services hold logic and never
  import `gin`; repositories persist only.
- **Domain pattern:** `internal/discovery` is the template — a domain package with a
  service, a repository interface, an in-memory fallback repo, and a Postgres repo.
- **Middleware chain:** `recovery → requestID → requestLogging`. No auth, no CORS,
  no rate limiting.
- **Composition:** a manual dependency-injection `Container` wires repos → services →
  handlers; handlers self-register routes under `/api/v1`.
- **Response envelope (ADR-0003):** `{data,error}` — immutable.
- **Config:** environment-based, no secret manager; secrets never read from `.env`.

## Proposed changes

### New backend domain: `internal/auth`

A domain package mirroring `internal/discovery`, with these **responsibilities**
(described, not coded):

| Component | Responsibility |
|-----------|----------------|
| Auth service | Orchestrates register / login / refresh / logout / Google login; enforces domain rules; returns uniform errors. No `gin`, no SQL. |
| User repository | Persistence for users (create, find by email / id / google-sub, link google-sub). Interface + in-memory + Postgres impls. |
| Refresh-token repository | Persistence for refresh tokens (store hash, find by hash, revoke, revoke-all). Interface + in-memory + Postgres impls. |
| Password component | One-way hashing and verification. |
| Token component | Issue/parse access JWT; generate/hash refresh-token values. |
| Google OAuth component | Build the authorization URL; exchange the code and fetch the verified profile (Authorization-Code flow). Isolated behind an interface for testing. |

### New HTTP surface & middleware

```
recovery → requestID → requestLogging → corsMiddleware → [routes]
                                                          └─ protected group: authMiddleware → handler
```

- **`AuthHandler`** self-registers `/api/v1/auth/*` (see [`api.md`](api.md)). It
  validates input, calls the service, maps sentinel errors → envelope codes, and
  manages the refresh cookie.
- **`authMiddleware`** validates the `Bearer` access token and injects the user id
  into the request context for protected routes (e.g. `/auth/me`).
- **`corsMiddleware`** enables credentialed cross-origin requests for an explicit
  allowed-origin list (D8).
- Discovery routes remain **outside** the protected group — public and unchanged (D7).

### Where auth checks live

- **Handler:** shape/format validation, cookie read/write, error mapping.
- **Middleware:** authentication (is this a valid token? who is the user?).
- **Service:** authorization/business rules (credentials valid? token rotatable?).
- **Repository:** no policy — persistence only.

### Token lifecycle

```
 register/login ─► issue access JWT (~15m)  ─────────────► client memory
                └► create refresh token, store HASH ─────► httpOnly Secure cookie
 protected call ─► authMiddleware verifies access JWT ──► inject user id
 access expired ─► client calls /refresh (cookie) ─────► rotate:
                       revoke old refresh ─► issue new access + new refresh
 logout ────────► revoke current refresh ─► clear cookie (idempotent)
```

Refresh tokens are stored **hashed**; the raw value lives only in the cookie. The set
of a user's non-revoked, unexpired refresh tokens *is* their set of active sessions.

### Sequence flows (step lists)

- **Register:** validate → ensure email unused → hash password → create user →
  issue access + refresh → store refresh hash → set cookie → return user + access.
- **Login:** validate → find user → verify password (uniform error on failure) →
  issue tokens → store refresh hash → set cookie → return.
- **Refresh:** read cookie → look up by hash → check not expired/revoked → revoke old
  → issue new pair → reset cookie → return new access.
- **Logout:** read cookie → revoke by hash (idempotent) → clear cookie.
- **Google SSO:** `/google/login` sets a `state` cookie and redirects to Google →
  Google redirects to `/google/callback` → verify `state` → exchange code → fetch
  verified profile → find-by-google-sub, else link verified email, else create →
  issue tokens → set cookie → redirect back to the web app.

### Composition & config

- The `Container` gains the auth repos, components, service, handler, and middleware.
- New config keys (JWT secret, token TTLs, Google client id/secret/redirect, allowed
  origins, cookie-secure toggle) are loaded from the environment — never `.env`.
- A new **ADR-0007 (authentication approach)** records these decisions alongside
  ADR-0002 (layering) and ADR-0003 (envelope).

## Open questions

1. **Access-token TTL vs UX** — is ~15 min right, or tune idle/absolute lifetimes?
2. **Refresh-token reuse detection** — on reuse of a rotated token, auto-revoke the
   whole session family (breach response) now, or defer? See [`security.md`](security.md).
3. **Stateless vs cached user lookup** in `authMiddleware` — pure stateless now, or a
   short-TTL cache later to support instant global logout / role changes?
4. **Rate limiting** — add brute-force protection in this feature or a follow-up?
5. **Multi-device session management** — surface active sessions to users now, or keep
   internal? See [`database.md`](database.md).
6. **Google linking policy** — auto-link verified email vs require explicit
   confirmation. See [`security.md`](security.md).
7. **CI gating** — add build/test gating alongside this feature? See [`release.md`](release.md).

## Related repository documentation

- [`overview.md`](overview.md) · [`specification.md`](specification.md) · [`api.md`](api.md) · [`database.md`](database.md) · [`security.md`](security.md).
- [`docs/architecture.md`](../../architecture.md) — system architecture (current state & roadmap).
- [`docs/adr/`](../../adr/) — ADR-0002 (Handler/Service/Repository layering), ADR-0003 ({data,error} envelope); proposed ADR-0007 (authentication approach).
- [`apps/api/CLAUDE.md`](../../../apps/api/CLAUDE.md) — layering, DI Container, handler naming, error codes.
- [`docs/standards/backend/`](../../standards/backend/) — architecture, patterns, request/response, errors.
- [`/CLAUDE.md`](../../../CLAUDE.md) — anti-goals and approval gates for auth/migrations/external calls.
- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — backend DoD gates.
