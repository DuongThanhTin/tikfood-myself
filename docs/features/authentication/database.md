# Authentication — Database

Schema changes required for authentication: the users and refresh-token tables, their constraints, and migration approach.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

## Purpose of this document

This document specifies the **persistence design** for authentication: the two new tables, their columns/constraints/indexes, how refresh-token rotation and revocation map onto rows, retention of expired tokens, and how migrations are numbered and applied. It presents the schema as field tables (design intent); the final SQL is authored during implementation. Migrations are a **human-approval-gated** area (approval granted during Discovery).

## Existing implementation

- **Engine:** PostgreSQL 16 + PostGIS 3.4. Access via raw `database/sql` with the pgx driver — **no ORM, no query builder**.
- **Existing tables (discovery domain only):** `venues`, `dishes`, `venue_dishes`, `tags`, `venue_tags`, `dish_tags`, `locations`, `location_aliases`, `venue_opening_hours`, `social_videos`, `trend_scores`, `ai_summaries`, `ingestion_runs`. **No user, session, or auth tables exist.**
- **Migrations:** plain SQL files in `apps/api/migrations/` (`001_…`, `003_…`, `004_…`, `006_…`), applied by the PostGIS **Docker entrypoint at container init** — there is **no migration runner library** and no built-in down-migration mechanism.
- Queries are parameterized; PostGIS functions (`ST_Distance`, etc.) are used for geo.

## Proposed changes

Two **additive** migrations, next in sequence. No existing table is altered.

### Migration `007_auth_users.sql` — `users`

| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| `id` | uuid | PK | Generation strategy — see Open questions |
| `email` | citext | UNIQUE, NOT NULL | Case-insensitive uniqueness |
| `password_hash` | text | NULL | Null for OAuth-only accounts |
| `display_name` | text | | May be empty |
| `google_sub` | text | UNIQUE, NULL | Google subject id when linked |
| `created_at` | timestamptz | NOT NULL, default now() | |
| `updated_at` | timestamptz | NOT NULL, default now() | Maintenance — see Open questions |

Indexes: unique on `email`; unique on `google_sub` (partial, where not null). Requires the `citext` extension (`CREATE EXTENSION IF NOT EXISTS citext;`).

### Migration `008_auth_refresh_tokens.sql` — `refresh_tokens`

| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| `id` | uuid | PK | |
| `user_id` | uuid | FK → `users(id)` ON DELETE CASCADE, NOT NULL | |
| `token_hash` | text | UNIQUE, NOT NULL | **Hash** of the raw token; raw value never stored |
| `expires_at` | timestamptz | NOT NULL | Absolute expiry |
| `revoked_at` | timestamptz | NULL | Set on rotation/logout/revoke-all |
| `created_at` | timestamptz | NOT NULL, default now() | |
| `user_agent` | text | | Captured for audit |
| `ip` | inet | NULL | Captured for audit |

Indexes: unique on `token_hash`; index on `user_id`.

### Illustrative schema sketch (illustrative, not final)

```
-- ILLUSTRATIVE ONLY — final SQL authored during implementation
-- CREATE EXTENSION IF NOT EXISTS citext;
-- users(id uuid pk, email citext unique not null, password_hash text null,
--       display_name text, google_sub text unique null, created_at, updated_at)
-- refresh_tokens(id uuid pk, user_id uuid fk->users(id) on delete cascade,
--       token_hash text unique not null, expires_at timestamptz not null,
--       revoked_at timestamptz null, created_at, user_agent text, ip inet null)
```

### Rotation & revocation semantics (data level)

- **Issue:** on register/login/Google/refresh, store a new row with the token hash and `expires_at`.
- **Rotate (refresh):** verify the presented hash is unexpired and `revoked_at IS NULL`; set `revoked_at = now()` on the old row and insert a fresh row atomically.
- **Logout:** set `revoked_at = now()` for the presented token's row (idempotent).
- **Revoke-all (optional):** set `revoked_at = now()` for all rows of a `user_id`.
- **Validity check:** a token is usable only when `revoked_at IS NULL AND expires_at > now()`.

### Retention / cleanup

Expired/revoked rows accumulate. Cleanup options (decision deferred — see Open questions): a scheduled SQL job, a Go background sweep, opportunistic deletion on login, or no cleanup for MVP (rows are harmless once `expires_at` passes).

### Migration numbering & rollback stance

- Next numbers are **007** and **008**, continuing the existing gapped sequence.
- **Additive only** — no destructive changes to discovery tables; the feature is reversible by not applying the migrations. Because there is no runner/down-migration tooling, any future table removal is a manual, human-gated operation (see [`release.md`](release.md)).

## Open questions

1. **UUID generation** — DB-side (`gen_random_uuid()` via `pgcrypto`) vs app-side generation in Go? Determines which extension the migration enables.
2. **`citext` availability** — confirm it is installable in the PostGIS image; fallback is `text` + a unique index on `lower(email)`.
3. **Token cleanup** — scheduled job, background sweep, opportunistic, or deferred for MVP?
4. **Session listing** — expose stored `user_agent`/`ip` as a "your sessions" feature, or capture-only for audit? (Borders on scope beyond MVP.)
5. **`updated_at` maintenance** — DB trigger vs application-managed on write?
6. **`email_verified` column** — add now to support the OAuth-linking control in [`security.md`](security.md), or defer alongside password-flow email verification?
7. **Password history** — out of scope for MVP; confirm no `password_history` table is needed.

## Related repository documentation

- [`docs/standards/backend/database.md`](../../standards/backend/database.md) — database standards (required reading for migrations).
- [`apps/api/migrations/`](../../../apps/api/migrations/) — existing migration files and numbering.
- [`docs/REPOSITORY-MAP.md`](../../REPOSITORY-MAP.md) — notes migrations are a **human-approval** area.
- [`docs/adr/0001-gin-and-pgx-not-gorm.md`](../../adr/0001-gin-and-pgx-not-gorm.md) — why raw pgx, no ORM.
- [`/CLAUDE.md`](../../../CLAUDE.md) · [`apps/api/CLAUDE.md`](../../../apps/api/CLAUDE.md) — migration approval & layering rules.
- Siblings: [`overview.md`](overview.md) · [`specification.md`](specification.md) · [`architecture.md`](architecture.md) · [`api.md`](api.md) · [`security.md`](security.md) · [`frontend.md`](frontend.md) · [`ui.md`](ui.md) · [`component-reuse.md`](component-reuse.md) · [`testing.md`](testing.md) · [`release.md`](release.md) · [`tasks.md`](tasks.md)
