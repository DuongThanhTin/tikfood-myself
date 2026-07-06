# Authentication — API Contract

HTTP contract for the authentication endpoints: routes, request/response shapes, error codes, and cookie behavior.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

## Purpose of this document

This document defines the **external HTTP contract** for authentication so frontend and backend can be built against one agreed surface. It fixes route paths, request fields and their validation intent, success/error response shapes, status codes, and the refresh-cookie behavior. It describes the contract as design (field tables), not implementation — handler logic lives in [`architecture.md`](architecture.md) and validation rules in [`specification.md`](specification.md).

## Existing implementation

- The API is served by Go/Gin under a versioned group `/api/v1`, with an unversioned `GET /health`.
- Existing routes are discovery-only and **public**: `GET /api/v1/discovery/venues`, `GET /api/v1/discovery/venues/:slug`, `GET /api/v1/map/venues`.
- Every response already uses the immutable envelope `{ "data": ..., "error": { "code", "message", "details?" } }`.
- Error codes today: `invalid_request` (400), `not_found` (404), `internal_error` (500). There is **no** auth endpoint, no `Authorization` handling, and no cookies.

## Proposed changes

All new endpoints live under `/api/v1/auth`. Discovery endpoints stay public and unchanged. Access is via `Authorization: Bearer <access_jwt>`; the refresh token travels only in an httpOnly cookie.

### Shared response envelope

Success: `{ "data": <object> }`. Failure: `{ "error": { "code": <string>, "message": <string>, "details?": <object> } }`. No field is ever removed or renamed.

### Refresh cookie

| Property | Value |
|----------|-------|
| Name | `tikfood_refresh` |
| Contents | Opaque refresh token (raw value; only its hash is stored server-side) |
| Flags | `HttpOnly`, `Secure` (toggleable to false for local http via config), `SameSite=Lax` |
| Path | `/api/v1/auth` |
| Max-Age | Refresh-token TTL (~30 days) |
| Set on | `register`, `login`, `refresh` (rotated), Google callback |
| Cleared on | `logout` (Max-Age=0) |

### Endpoints

**POST /api/v1/auth/register** — auth: none — create an account and start a session.

Request:

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `email` | string | yes | Valid email format; unique (case-insensitive) |
| `password` | string | yes | Min 8 chars; max 72 bytes |
| `display_name` | string | no* | Non-empty if provided (see Open questions) |

Success `201`: `{ "data": { "user": <User>, "access_token": string, "access_expires_at": RFC3339 } }` + Set-Cookie refresh. Errors: `422 domain_rejected` (validation / email taken), `400 invalid_request` (malformed body), `500 internal_error`.

**POST /api/v1/auth/login** — auth: none — authenticate with email/password.

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `email` | string | yes | Valid email format |
| `password` | string | yes | Non-empty |

Success `200`: `{ "data": { "user": <User>, "access_token": string, "access_expires_at": RFC3339 } }` + Set-Cookie refresh. Errors: `401 unauthorized` (uniform "invalid credentials" for unknown user **or** wrong password — no enumeration), `400 invalid_request`, `500 internal_error`.

**POST /api/v1/auth/refresh** — auth: refresh cookie — rotate the session and mint a new access token.

Request: no body; reads `tikfood_refresh` cookie. Success `200`: `{ "data": { "access_token": string, "access_expires_at": RFC3339 } }` + Set-Cookie (rotated) refresh. Errors: `401 unauthorized` (missing/expired/revoked/reused refresh token), `500 internal_error`.

**POST /api/v1/auth/logout** — auth: refresh cookie (optional) — revoke the current session.

Request: no body; reads `tikfood_refresh` cookie. Success `200`: `{ "data": { "logged_out": true } }` + Set-Cookie clearing refresh. Idempotent (still `200` if already logged out). Errors: `500 internal_error`.

**GET /api/v1/auth/me** — auth: Bearer (protected) — return the current user.

Request: `Authorization: Bearer <access_jwt>`. Success `200`: `{ "data": { "user": <User> } }`. Errors: `401 unauthorized` (missing/expired/invalid token), `500 internal_error`.

**GET /api/v1/auth/google/login** — auth: none — begin Google SSO.

Behavior: sets a short-lived opaque `state` cookie and returns `302` redirect to Google's consent URL. No JSON body on the happy path.

**GET /api/v1/auth/google/callback** — auth: none — complete Google SSO.

| Query param | Type | Required | Rules |
|-------------|------|----------|-------|
| `code` | string | yes | Google authorization code |
| `state` | string | yes | Must match the `state` cookie (CSRF) |

Behavior: verifies `state`, exchanges the code, provisions/links the user, sets the refresh cookie, and `302`-redirects to the web app's callback page with a short-lived one-time hand-off (mechanism confirmed in [`frontend.md`](frontend.md)/[`security.md`](security.md)). Errors: `401 unauthorized` (state mismatch / exchange failure), `403 forbidden` (unverified Google email during link), `500 internal_error`.

### User object (shape returned by the API)

| Field | Type | Notes |
|-------|------|-------|
| `id` | string (uuid) | |
| `email` | string | |
| `display_name` | string | May be empty |
| `created_at` | RFC3339 | |

`password_hash` and `google_sub` are **never** serialized.

### Shared error-code table

| HTTP | `error.code` | Meaning |
|------|--------------|---------|
| 400 | `invalid_request` | Malformed/unparseable request |
| 401 | `unauthorized` | Missing/invalid/expired credentials or token *(new)* |
| 403 | `forbidden` | Authenticated but not permitted *(new)* |
| 404 | `not_found` | Resource does not exist |
| 422 | `domain_rejected` | Validation failure / business rule (e.g., weak password, email taken) |
| 500 | `internal_error` | Unexpected server error |

## Open questions

1. **Email-taken status** — `422 domain_rejected` (proposed) vs a new `409 conflict` code?
2. **Google hand-off** — one-time exchange code (preferred) vs URL-fragment access token? See [`security.md`](security.md).
3. **`logout-all` endpoint** — expose `POST /api/v1/auth/logout-all` now or defer?
4. **`me` scope** — include auth method / last-login beyond core `User` fields?
5. **Forgot-password/verify endpoints** — confirmed out of scope this iteration (UI placeholder only)?
6. **Cookie name/path** — confirm `tikfood_refresh` and `Path=/api/v1/auth` suit the deployment/proxy topology.

## Related repository documentation

- [`docs/contracts/api.md`](../../contracts/api.md) — API contract rules.
- [`docs/standards/backend/request-response.md`](../../standards/backend/request-response.md) — request parsing & envelope.
- [`docs/standards/backend/errors.md`](../../standards/backend/errors.md) — error codes & mapping.
- [`docs/adr/0003-data-error-response-envelope.md`](../../adr/0003-data-error-response-envelope.md) — the immutable `{data,error}` envelope.
- [`/CLAUDE.md`](../../../CLAUDE.md) — "no breaking API changes; keep response format unchanged."
- Siblings: [`overview.md`](overview.md) · [`specification.md`](specification.md) · [`architecture.md`](architecture.md) · [`database.md`](database.md) · [`security.md`](security.md) · [`frontend.md`](frontend.md) · [`ui.md`](ui.md) · [`component-reuse.md`](component-reuse.md) · [`testing.md`](testing.md) · [`release.md`](release.md) · [`tasks.md`](tasks.md)
