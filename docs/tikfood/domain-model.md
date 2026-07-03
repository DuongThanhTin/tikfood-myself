# Product Domain Model & Glossary

> **Why this doc exists:** `docs/tikfood/` explains the product *vision, scope, and
> anti-goals*, but the **entities and vocabulary** (what a Venue is, how Dishes, Social
> Videos, Trend Scores, AI Summaries, and Location aliases relate) lived only implicitly
> in the Go models and SQL. Shared, precise vocabulary — a *ubiquitous language* —
> prevents drift (e.g. the location-alias logic that currently exists in three places).
> This is the reference for those terms. Field-level truth is the code
> (`apps/api/internal/discovery/model.go`) and migrations — code wins.

**Related:** [`product-vision.md`](product-vision.md) · [`mvp-scope.md`](mvp-scope.md) ·
[`anti-goals.md`](anti-goals.md) · [`docs/services/api.md`](../services/api.md) ·
[`docs/standards/backend/database.md`](../standards/backend/database.md).

## Core entities

```text
Venue ─┬─ SocialVideo (many)      ← social proof
       ├─ VenueDish (many)        ← dish-first discovery
       ├─ OpeningHour (many)      ← open-now filter
       ├─ TrendScore (1, roadmap) ← "why trending"
       └─ AISummary (1, roadmap)  ← AI "why go here"
Location ── LocationAlias (many)  ← city/district normalization
```

### Venue
A physical food place — the primary unit of discovery. Source of truth:
`discovery.Venue`. Key fields: `id`, `name`, `slug`, `short_description`, `about`,
`address`, `city`, `district`, `latitude`/`longitude`, `categories`, `price_level`,
`avg_price_min_vnd`/`avg_price_max_vnd`, `currency`, `social_video_count`,
`trend_score`, `trending_dishes`, `ai_summary`, `distance_meters` (computed on the
Postgres path only). Discovery is **read-only** over venues.

### SocialVideo
A short-form video referencing a venue — the *social proof* signal. Fields: `platform`
(`tiktok|instagram|youtube|facebook|other`), `url`, `creator_handle`, `caption`,
`thumbnail_url`, `view_count`, `like_count`, `published_at`. Backed by the social-videos
schema (migration 004).

### VenueDish
A dish associated with a venue — enables *dish-first* search. Fields include `name`,
`slug`, `category`, `cuisine`, `price_min_vnd`/`price_max_vnd`, `mention_count`,
`video_count`, `view_count`, `trend_score`. Note: `Venue.trending_dishes` is a lightweight
string list; `VenueDish` is the richer entity.

### OpeningHour
Per-day hours (`day_of_week`, `open_time`, `close_time`, `is_closed`) powering the
`open_now` filter.

### TrendScore *(roadmap)*
The realtime "how hot is this" signal. The `trend_scores` table exists (migration 001)
and `Venue.trend_score` is surfaced, **but no worker computes it yet** — current values
are static seed data. The trend-scoring worker is roadmap
([`docs/architecture.md`](../architecture.md) §2.2, `IMPROVEMENTS.md`).

### AISummary *(roadmap)*
The AI-generated "why this place / why trending" text. `ai_summaries` table exists and
`Venue.ai_summary` is surfaced, **but no worker populates it yet** — static seed today.

### Location & LocationAlias
Canonical places and their aliases for normalizing `city`/`district` queries (e.g.
`quan-1` ↔ "Quận 1"). Backed by the location-identity schema (migration 003). **The DB
is the intended source of truth** for alias resolution; in-memory (`search.go`) and
frontend (`lib/api.ts`) alias logic are dev-only approximations to be consolidated (see
`IMPROVEMENTS.md`).

## Glossary

| Term | Meaning |
| --- | --- |
| **Discovery** | Browsing/searching venues, dishes, and social proof. The entire MVP scope. |
| **Dish-first** | Users find food by dish, not only by restaurant. |
| **Map-first** | The map is the primary browse surface. |
| **Social proof** | SocialVideos + counts that show a place is talked about. |
| **Trend score** | Realtime popularity signal per venue/dish (roadmap worker). |
| **AI summary** | AI "why go / why trending" text per venue (roadmap worker). |
| **Hidden gem** | A high-quality, low-visibility venue the product aims to surface. |
| **Alias** | An alternate name for a city/district resolved to a canonical Location. |
| **Anti-goal** | A feature deliberately excluded from the MVP — see [anti-goals.md](anti-goals.md). |

## Ingestion note
Venues can be ingested from OpenStreetMap via `apps/api/internal/ingest` + `cmd/ingest`
(migration 006), currently a CLI path not exposed over HTTP. Ingestion must respect
platform terms and privacy ([`docs/security.md`](../security.md)).
