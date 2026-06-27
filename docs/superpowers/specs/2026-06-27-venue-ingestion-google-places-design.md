# Design: Venue Ingestion from Google Places (Sub-project A)

Date: 2026-06-27
Status: DEFERRED (future option — no budget for a Google Places billing account).
The active sub-project A uses OpenStreetMap instead — see
`2026-06-27-venue-ingestion-osm-design.md`. This document is retained on purpose so
the team can develop the Google Places direction later (richer data: rating, photos,
structured opening hours).

## Context

TikFood's core differentiation (social trend discovery) needs real venue data. The
larger effort splits into four independent sub-projects:

- **A. Venue ingestion** (this spec) — populate `venues`/`locations`/`tags` from a
  directory source.
- B. Social video ingestion (TikTok/Instagram → `social_videos`).
- C. Entity resolution (link videos ↔ venues/dishes).
- D. Trend scoring + AI summary workers.

Build order A → B → C → D. This spec covers **A only**.

The Postgres schema already models venues, locations, location aliases, tags, and
opening hours; what is missing is the pipeline that fills them.

## Goal

A re-runnable Go CLI that ingests Ho Chi Minh City restaurants/cafes from the Google
Places API **one district at a time**, upserting them idempotently so a district can
be re-run safely after fixing mapping errors — no full re-migration.

## Key decisions (from brainstorming)

1. **Source: Google Places API (New).** Legal and structured. API key via env
   `GOOGLE_PLACES_API_KEY`. No dish/price data is available from Places — dishes are
   out of scope for sub-project A.
2. **Scope v1: Ho Chi Minh City, district-by-district.** Discovery uses Text Search
   per district + category (`restaurant`/`cafe`). City is fixed to Ho Chi Minh.
3. **Idempotent upsert keyed on `place_id`.** `ON CONFLICT (source,
   source_external_id)` updates existing rows instead of duplicating. Re-running a
   district is safe and cheap; this is what makes incremental ingestion low-risk.
4. **Fields captured in v1:** core (name, address, latitude/longitude, place_id,
   opening hours) plus Google rating + rating count, photo references, Places
   `types` → tags, phone + website.
5. **Worker shape:** Go CLI `apps/api/cmd/ingest`, run manually now
   (`--district`, `--category`, `--limit`, `--dry-run`), cron/n8n-schedulable later.
   Reuses existing config + DB layer + domain model.
6. **Resolve-or-create district in `locations`.** Only Quận 1 and Quận 3 are seeded
   today. When ingesting a new district, the service ensures a `locations` row (and
   basic alias) exists before linking `district_location_id`.
7. **Per-district tracking via `ingestion_runs` audit table.** Each run records
   district, category, timestamps, fetched/upserted/failed counts, and status, so
   each district can be verified before moving on. No purge command — idempotent
   re-run covers re-do.

## Architecture

Respects `apps/api/CLAUDE.md` layering. Dependencies point one way:

```
cmd/ingest (CLI)
  → internal/ingest (Service: orchestration, mapping, location/tag resolution)
      → internal/places (Google Places API client: TextSearch + PlaceDetails)
      → VenueRepository / IngestionRepository (postgres: upsert)
```

The ingest Service never imports HTTP/`gin.Context`. The Places client is an
external-API adapter behind an interface so the Service can be tested with a fake.

## Components

- **`apps/api/cmd/ingest/main.go`** — CLI entry. Parses flags
  (`--district` required, `--category` default `restaurant`, `--limit`,
  `--dry-run`), loads `config`, builds the Places client + Service + repositories,
  runs one district, prints a summary, exits non-zero on fatal errors.
- **`internal/places`** — Google Places API (New) client. `TextSearch(ctx, query,
  opts)` and `PlaceDetails(ctx, placeID, fieldMask)`. Uses field masks to limit cost.
  Reads `GOOGLE_PLACES_API_KEY`. Never logs the key. Backoff on quota/rate errors.
  Defined behind an interface (`PlacesClient`) for testing.
- **`internal/ingest`** — `Service` with `IngestDistrict(ctx, district, category,
  opts)`. For each search result: fetch details → map Places payload → domain
  `Venue` (+ opening hours, tags) → resolve-or-create the district `locations` row →
  upsert. Per-place errors are logged and skipped; the run continues. Returns a
  summary (counts). A pure mapper function (Places JSON → `Venue`) is isolated for
  table-driven testing.
- **Repository (postgres)** — new methods:
  - `UpsertVenueFromSource(ctx, venue, source, externalID) (venueID, error)` —
    `ON CONFLICT (source, source_external_id) DO UPDATE`, sets `source_updated_at`.
  - `ReplaceVenueOpeningHours(ctx, venueID, hours)`.
  - `ResolveOrCreateDistrict(ctx, citySlug, districtName) (locationID, error)`.
  - `UpsertVenueTags(ctx, venueID, tagSlugs)` (creates missing `tags`, links
    `venue_tags`).
  - `RecordIngestionRun(ctx, run) (runID, error)` / `FinishIngestionRun(ctx, runID,
    counts, status)`.
- **Config** — add `GooglePlacesAPIKey string` (env `GOOGLE_PLACES_API_KEY`) to
  `internal/config/config.go`.

## Schema changes — migration `006_venue_ingestion_schema.sql`

- `alter table venues add column google_rating numeric`, `google_rating_count int`.
- `alter table venues add constraint venues_source_external_unique unique (source,
  source_external_id)` (required for upsert `ON CONFLICT`).
- **Convert `venues.location` to a generated column** derived from
  `longitude`/`latitude`:
  `location geography(Point,4326) GENERATED ALWAYS AS
  (ST_SetSRID(ST_MakePoint(longitude, latitude),4326)::geography) STORED`.
  (Existing data is seed-only, so dropping and re-adding the column is safe. This
  folds in the HIGH-priority geo-sync item from the DB review, because ingestion
  writes coordinates and distance queries depend on `location` staying in sync.)
- New table `ingestion_runs`:
  `id uuid pk, source text, district_slug text, category text, status text
  (running|completed|failed), fetched_count int, upserted_count int,
  failed_count int, error text, started_at timestamptz, finished_at timestamptz`.

## Data flow (one district)

1. CLI: `ingest --district=quan-1 --category=restaurant`.
2. Service records an `ingestion_runs` row (status `running`).
3. `places.TextSearch("restaurant in Quận 1, Hồ Chí Minh", …)` → list of place IDs
   (paginated up to `--limit`).
4. For each place: `places.PlaceDetails(placeID, fieldMask)` → map to `Venue`.
5. Resolve-or-create the district in `locations`; map Places `types` → tag slugs.
6. `--dry-run`: log what would be written and stop. Otherwise upsert venue + hours +
   tags; increment counts.
7. Finish the `ingestion_runs` row with counts and `completed`/`failed`.
8. CLI prints summary: `district=quan-1 fetched=N upserted=N failed=N`.

## Error handling & ToS

- Per-place failures are logged with the place_id and skipped; one bad place never
  aborts the district. Fatal errors (missing API key, DB down) exit non-zero and
  mark the run `failed`.
- Google quota/rate-limit responses trigger bounded backoff + retry.
- API key only from env; never hardcoded, never logged.
- ToS compliance: `place_id` (as `source_external_id`) is stored permanently; other
  fields are refreshed on re-run (`source_updated_at`). Photos are stored as Places
  **photo references**, not re-hosted media.

## Testing

- No network in tests/CI. Use saved Places JSON fixtures.
- Table-driven unit tests for the mapper (Places payload → `Venue`, including
  opening-hours and types→tags mapping, missing-field handling).
- Service test with a fake `PlacesClient` + mock repository: verifies upsert calls,
  resolve-or-create district, `--dry-run` writes nothing, per-place error is skipped,
  and run counts are recorded.
- The project has no DB-backed test harness today (existing tests use the in-memory
  fallback). The idempotent `ON CONFLICT` upsert is therefore verified manually
  during district runs (re-run a district, confirm counts update not duplicate); a
  DB-backed repository test is deferred to the CI/test-harness improvement (P1) and
  is out of scope here.

## Out of scope

- Dishes/prices (not in Places), social videos (sub-project B), entity resolution
  (C), trend scoring/AI summaries (D).
- Provinces other than Ho Chi Minh City.
- A purge/rollback command (idempotent re-run covers re-do).
- The `apps/web` UI swap to show real Google rating/photos instead of the fabricated
  `mediaByVenue` data — enabled by this work but tracked separately (P0 in
  IMPROVEMENTS.md).
- A migration-runner tool (still tracked as a P1 improvement); migration 006 follows
  the existing `if not exists` / idempotent SQL convention.
