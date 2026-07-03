# Application Security Standard

> **Why this doc exists:** [`docs/security.md`](../security.md) governs the **automation
> runner** (secrets, git, shell allowlist, human-approval gates for AI-generated PRs).
> It does not cover **product** application security — how `apps/api` and `apps/web`
> should handle input, authorization, data exposure, and abuse. That guidance was
> scattered across the API/errors/logging standards. This doc consolidates it for the
> product apps. It complements, and does not replace, `docs/security.md`.

**Related:** [`docs/security.md`](../security.md) (runner/automation security) ·
[`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) (hard rules) ·
[`docs/standards/backend/errors.md`](backend/errors.md) ·
[`docs/standards/logging-observability.md`](logging-observability.md).

Scope note: TikFood MVP is **discovery-only and read-only** with **no user accounts**.
Several controls below are therefore stated as *principles + when-to-apply*, not
requirements to build now. Anything that adds auth, writes, or external calls requires
**human approval** ([`AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §4).

## 1. Input validation & injection

- Validate at the boundary (handler): length, range, enum, type — as
  `internal/http/venue_request.go` already does. Reject invalid input with
  `invalid_request` (400); never trust client input downstream.
- **SQL:** always use parameterized queries via `database/sql`/pgx
  ([ADR-0001](../adr/0001-gin-and-pgx-not-gorm.md)). Never build SQL by string
  concatenation of user input.
- **Frontend:** rely on React's default escaping; avoid `dangerouslySetInnerHTML`. Treat
  any venue/social text as untrusted when rendering.

## 2. Output & data exposure

- Never leak SQL, stack traces, internal hostnames, tokens, or raw error strings to
  clients — map to safe codes (see [`backend/errors.md`](backend/errors.md)).
- Return only fields the client needs; do not echo internal/debug fields.
- Do not surface fabricated data as real (see `IMPROVEMENTS.md` — placeholder
  ratings/reviews must be clearly marked until backed by real data).

## 3. Authentication & authorization (when introduced)

- None today (public, read-only discovery). **Adding auth/authz requires human
  approval** and its own ADR.
- When added: deny-by-default; check authorization in the **service** layer (not the
  handler alone); never trust a client-supplied identity/role.

## 4. Transport, CORS & rate limiting

- Serve over HTTPS in any deployed environment; do not disable TLS verification.
- **CORS:** allow only known frontend origins; do not use `*` with credentials.
- **Rate limiting / abuse:** discovery endpoints are unauthenticated — when deployed,
  put rate limiting at the edge/gateway. Bound expensive queries (the `limit` cap of
  100 and `radius_m` cap of 50000 are examples of this discipline).

## 5. Secrets & configuration

- Secrets come from env/`internal/config` only — never hardcoded, never committed,
  never logged. Never read `.env*` (see [`docs/security.md`](../security.md)).
- Third-party keys (maps, model providers) live in env and are injected, not embedded in
  the client bundle unless explicitly public (`NEXT_PUBLIC_*`).

## 6. Third-party & client-side calls

- Browser-side calls to public/unmetered services (e.g. the OSRM demo, tile servers) are
  not production-safe — route through a controlled backend/managed provider. Track in
  `IMPROVEMENTS.md`. New external network calls require **human approval**.

## 7. Privacy & platform compliance

- Social ingestion must respect platform terms, rate limits, privacy, and legal
  constraints; scraping/collection changes require human review
  ([`docs/security.md`](../security.md) → *Social Platform Compliance*).
- Minimize personal data; do not log location/user data in a way that identifies
  individuals.

## 8. Dependencies

- Follow [`dependency-management.md`](dependency-management.md): justify new deps, prefer
  well-maintained ones, review for known vulnerabilities before adding.

## Checklist for a security-relevant change

- [ ] Inputs validated & bounded at the boundary.
- [ ] Queries parameterized; no string-built SQL.
- [ ] No secrets/SQL/internal detail in responses or logs.
- [ ] External calls / auth / migrations → human approval obtained.
- [ ] Untrusted text safely rendered on the frontend.
- [ ] Change reflected in tests and, if a decision, an ADR.
