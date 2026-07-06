# Authentication — Release

One-line purpose: how the authentication feature is rolled out — PR sequencing, migrations, configuration, rollback, and operational readiness.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

Describe the path from merged code to a working, operable feature: the order PRs
land, what configuration and migrations each requires, how to roll back safely, and
what to check before and after launch. It complements [`tasks.md`](tasks.md) (the
milestone breakdown) with the operational view.

---

## Existing implementation

- **Git flow**: work on branches under `ai/`; never push to `main`/`master`; no
  force-push; no auto-merge; humans review and merge (per [`/CLAUDE.md`](../../../CLAUDE.md)
  and [`docs/security.md`](../../security.md)).
- **CI/CD**: none — `.github/workflows/` is absent; verification is local.
- **Local runtime**: `docker-compose.yml` brings up web, api, n8n, runner (a
  **protected path**; changes need human approval). Postgres/PostGIS via the DB image.
- **Config/secrets**: loaded from environment only (`os.Getenv`); no secret manager;
  `.env*` files are never read.

## Proposed changes

### PR sequencing (mirrors the 12 milestones)

| PR | Milestone | Depends on | Gated? |
|----|-----------|-----------|--------|
| 1 | M0 ADR + deps + config scaffolding | — | — |
| 2 | M1 users & refresh_tokens schema + repos | M0 | migration approval |
| 3 | M2 password + JWT core | M1 | — |
| 4 | M3 auth service | M2 | — |
| 5 | M4 auth HTTP endpoints (**contract freeze**) | M3 | auth approval |
| 6 | M5 auth middleware + CORS + /me | M4 | — |
| 7 | M6 Google OAuth (backend) | M5 | external-call approval |
| 8 | M7 web auth client + refresh + test harness | M4 | — |
| 9 | M8 AuthProvider | M7 | — |
| 10 | M9 login & register pages | M8 | — |
| 11 | M10 Google Sign-In (frontend) | M6, M9 | — |
| 12 | M11 forgot-password placeholder + wire authCta | M8 | — |
| 13 | M12 docs & design reconciliation | all | — |

Frontend PRs (8–12) may proceed in parallel once PR 5 (contract) merges.

### Migration deployment

- Migrations `007_auth_users.sql`, `008_auth_refresh_tokens.sql` are **additive**,
  applied by the existing PostGIS Docker entrypoint (no migration runner library).
- Apply in numeric order after the current highest migration; never modify existing
  discovery tables. Human approval required (granted in Discovery).
- Confirm the `citext` (or fallback) and UUID-generation extension prerequisites are
  present in the DB image before applying (see [`database.md`](database.md)).

### Required environment / configuration checklist

| Key | Purpose | Dev default |
|-----|---------|-------------|
| `JWT_SECRET` | Access-token signing key | required (no default; service refuses empty) |
| `ACCESS_TOKEN_TTL` | Access-token lifetime | `15m` |
| `REFRESH_TOKEN_TTL` | Refresh-token lifetime | `720h` (~30d) |
| `GOOGLE_CLIENT_ID` | OAuth client id | required for Google |
| `GOOGLE_CLIENT_SECRET` | OAuth client secret | required for Google |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | dev callback |
| `WEB_ORIGIN` | Allowed CORS origin | `http://localhost:3000` |
| `COOKIE_SECURE` | `Secure` flag on refresh cookie | `false` in local http, `true` elsewhere |

### Dev vs prod configuration differences

- **Cookies**: `COOKIE_SECURE=false` acceptable for local http; **must** be `true`
  in any deployed environment (https).
- **Origins**: dev allows `localhost:3000`; prod uses the real web origin(s).
- **Google redirect URIs**: register per-environment URIs in the Google console.
- **Secrets**: dev via env/`.env` (never read by agents); prod via an infra-managed
  secret store (out of scope here — human/infra-gated).

### Rollback strategy

- The feature is **additive and non-breaking**: discovery endpoints are untouched,
  new tables are independent. Reverting a PR removes its code cleanly.
- Migrations are additive; a table drop (if ever needed) is a manual, human-approved
  operation because there is no migration runner (see Open questions).

### Observability

- Structured `slog` JSON logs for auth events (login success/failure counts,
  token revocation) — **never** log passwords, tokens, hashes, or secrets.
- Reuse the existing `requestID` correlation; consider auth-event fields for audit.

### Pre-launch checklist

- All PR gates green (build, vet/typecheck, tests); security-review done on M4–M7.
- Migrations applied and verified on the target DB.
- All required config keys set; `COOKIE_SECURE=true`; Google redirect URIs registered.
- Discovery smoke test passes for anonymous users.

### Post-launch monitoring

- Watch auth error rates (401/422), refresh failures, and Google callback failures.
- Verify no secret-like values appear in logs.
- Track refresh-token table growth (informs the cleanup decision).

## Open questions

1. **CI before launch?** Add a minimal test-gating workflow first, or accept local
   per-PR verification for this feature?
2. **Prod secret store** — which mechanism and owner for `JWT_SECRET` + Google creds?
3. **Feature flag** — an env flag to enable/disable auth routes for staged rollout,
   or is PR-level sequencing sufficient?
4. **Refresh-token cleanup** — scheduled purge of expired/revoked rows, and where
   (worker vs SQL cron)? See [`database.md`](database.md).
5. **Rate limiting at launch** — blocker or fast-follow? See [`security.md`](security.md).
6. **Down-migration policy** — with no migration runner, how are rollbacks executed
   operationally if a drop is ever required?

## Related repository documentation

- [`/CLAUDE.md`](../../../CLAUDE.md) — git/security rules (branches under `ai/`, no push to main).
- [`docs/security.md`](../../security.md) — human-approval-required operations.
- [`docs/local-development.md`](../../local-development.md) — local setup.
- [`docker-compose.yml`](../../../docker-compose.yml) — local stack (protected path).
- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — merge gates.
- The approved roadmap: `docs/superpowers/plans/2026-07-05-authentication-system.md`.
- Sibling docs: [`tasks.md`](tasks.md), [`database.md`](database.md), [`security.md`](security.md), [`testing.md`](testing.md).
