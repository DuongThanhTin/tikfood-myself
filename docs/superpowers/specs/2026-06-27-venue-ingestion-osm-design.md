# Design: Venue Ingestion from OpenStreetMap / Overpass (Sub-project A)

Date: 2026-06-27
Status: Approved (pending spec review)
Active source: OpenStreetMap / Overpass.
Future option (retained, not active): the Google Places spec/plan
(`2026-06-27-venue-ingestion-google-places-design.md` and its plan) are kept as a
documented future direction — revisit when there is budget for a Places billing
account. They are NOT deleted.

## Context

TikFood's core differentiation (social trend discovery) needs real venue data. The
larger effort splits into four independent sub-projects:

- **A. Venue ingestion** (this spec) — populate `venues`/`locations`/`tags`.
- B. Social video ingestion (TikTok/Instagram → `social_videos`).
- C. Entity resolution (link videos ↔ venues/dishes).
- D. Trend scoring + AI summary workers.

Build order A → B → C → D. This spec covers **A only**.

The Postgres schema already models venues, locations, location aliases, tags, and
opening hours; what is missing is the pipeline that fills them.

## Source decision: OpenStreetMap via the Overpass API

- **Free, no API key, no credit card / billing account.** Google Places was the
  original plan but requires a billing account; the user has no budget for it.
- **Legal:** OSM data is ODbL-licensed. The product must display an OSM attribution
  ("© OpenStreetMap contributors") wherever this data is shown.
- **Overpass usage policy:** send a descriptive `User-Agent`, keep queries modest
  (one district per run), and tolerate transient `429`/`504` with backoff. Default
  endpoint `https://overpass-api.de/api/interpreter` with a mirror fallback.
- `source = 'openstreetmap'`, `source_external_id = '<type>/<osm_id>'` (e.g.
  `node/1770305118`).

### What OSM provides (verified against a live Quận 1 query, ~808 venues)

- latitude/longitude: 100% (every node) — the core for the map.
- name: ~77% (rows without a name are skipped).
- amenity (restaurant/cafe/fast_food): 100%; cuisine: ~45%.
- address (`addr:housenumber` + `addr:street`): ~37%.
- `opening_hours`: ~15% (sparse; free-text OSM syntax, needs a parser).
- phone (`phone`/`contact:phone`): ~13%; website (`website`/`contact:website`): ~9%.
- **No rating, no photos** — OSM has neither. The fabricated-rating P0 issue is NOT
  solved by this source; rating/photo columns stay NULL until a future source
  (e.g. Google) is added behind the same client interface.

## Goal

A re-runnable Go CLI (`apps/api/cmd/ingest`) that ingests Ho Chi Minh City
restaurants/cafes from Overpass **one district at a time**, upserting idempotently so
a district can be re-run safely after fixing mapping errors — no full re-migration.

## Key decisions

1. **Source: OpenStreetMap / Overpass.** No key, no cost. ODbL attribution required
   in the UI.
2. **Scope v1: Ho Chi Minh City, district-by-district.** A run targets one district
   (`--district=quan-1`) and filters `amenity` in `restaurant|cafe|fast_food`.
3. **Idempotent upsert keyed on `(source, source_external_id)`** where
   `source='openstreetmap'` and `source_external_id='<type>/<id>'`.
4. **Fields captured:** name, address, latitude, longitude, osm id, opening hours
   (when present, parsed), `amenity`+`cuisine` → tags, phone, website. No rating, no
   photos, no dishes/prices.
5. **Worker shape:** Go CLI `apps/api/cmd/ingest`
   (`--district`, `--category`, `--limit`, `--dry-run`), cron/n8n-schedulable later.
   Reuses existing config + DB layer + domain model.
6. **Resolve-or-create district in `locations`.** Only Quận 1 and Quận 3 are seeded;
   the service ensures a `locations` row exists before linking `district_location_id`.
7. **Per-district tracking via `ingestion_runs` audit table.** Each run records
   district, category, timestamps, fetched/upserted/failed counts, and status.

## Architecture

Respects `apps/api/CLAUDE.md` layering. Dependencies point one way:

```
cmd/ingest (CLI)
  → internal/ingest (Service: orchestration, mapping, location/tag resolution)
      → internal/ingest/overpass (Overpass API client: query one district)
      → VenueRepository / IngestionRepository (postgres: upsert)
```

The source is behind a `VenueSource` interface (`SearchDistrict(ctx, district,
category, limit) ([]SourcePlace, error)`), so a future Google client can be added
without touching the service, repository, migration, or CLI.

## Components

- **`apps/api/cmd/ingest/main.go`** — CLI entry. Flags `--district` (required),
  `--category` (default `restaurant`; maps to Overpass amenity set), `--limit`,
  `--dry-run`. Loads config, builds the Overpass client + Service + repositories,
  runs one district, prints a summary.
- **`internal/ingest/overpass`** — Overpass client. Builds and POSTs an Overpass QL
  query for the district's bounding area, decodes the JSON `elements` into
  `[]SourcePlace`. Sends a descriptive `User-Agent`; backoff on `429`/`504`; mirror
  fallback. The JSON→`SourcePlace` decode is a pure, testable function.
- **`internal/ingest/openinghours`** — pure parser for the common OSM
  `opening_hours` patterns (`Mo-Su HH:MM-HH:MM`, `Mo-Fr ...; Sa-Su ...`, `HH:MM-HH:MM`
  applied to all days, `24/7`). Unrecognized/complex expressions yield no hours
  (logged), never a crash.
- **`internal/ingest`** — `Service.IngestDistrict(ctx, opts)`. For each source place:
  map → domain `Venue` (+ hours, tags) → resolve-or-create district → upsert.
  Per-place errors are logged and skipped. A pure mapper (`SourcePlace` → domain) is
  isolated for table-driven testing.
- **Repository (postgres)** — `UpsertVenueFromSource`, `ReplaceVenueOpeningHours`,
  `ResolveOrCreateDistrict`, `UpsertVenueTags`, `RecordIngestionRun`,
  `FinishIngestionRun` (same interface a future Google client would reuse).
- **Config** — no API key needed; add `OverpassEndpoint string` (env
  `OVERPASS_ENDPOINT`, default `https://overpass-api.de/api/interpreter`) so the
  endpoint/mirror is configurable.

## Schema changes — migration `006_venue_ingestion_schema.sql`

- `alter table venues add column google_rating numeric`, `google_rating_count int`
  — kept for a future Google source; OSM leaves them NULL.
- `alter table venues add constraint venues_source_external_unique unique (source,
  source_external_id)` (required for upsert `ON CONFLICT`).
- **Convert `venues.location` to a generated column** from `longitude`/`latitude`
  (folds in the HIGH-priority geo-sync item from the DB review; ingestion writes
  coordinates and distance queries depend on `location`).
- New table `ingestion_runs` (id, source, district_slug, category, status,
  fetched_count, upserted_count, failed_count, error, started_at, finished_at).

## Data flow (one district)

1. CLI: `ingest --district=quan-1 --category=restaurant`.
2. Service records an `ingestion_runs` row (status `running`).
3. Overpass client queries the district area for the amenity set → `[]SourcePlace`.
4. For each place: map to `Venue`; parse `opening_hours` if present; map
   `amenity`+`cuisine` → tag slugs.
5. Resolve-or-create the district in `locations`.
6. `--dry-run`: log what would be written and stop. Otherwise upsert venue + hours +
   tags; increment counts.
7. Finish the `ingestion_runs` row with counts and `completed`/`failed`.
8. CLI prints: `district=quan-1 fetched=N upserted=N failed=N`.

## Error handling & ToS

- Per-place failures are logged with the osm id and skipped; one bad place never
  aborts the district. Fatal errors (DB down) exit non-zero and mark the run `failed`.
- Overpass `429`/`504` trigger bounded backoff + retry; a mirror endpoint is tried if
  the primary fails.
- ODbL: store OSM data and surface "© OpenStreetMap contributors" attribution in the
  UI wherever venue data is shown (tracked as a UI follow-up).
- Re-run friendly: `source_external_id` stored permanently; fields refreshed on
  re-run (`source_updated_at`).

## Testing

- No network in tests/CI. Use saved Overpass JSON fixtures.
- Table-driven unit tests for: the Overpass JSON decoder (`elements` → `SourcePlace`),
  the `opening_hours` parser (common patterns + unrecognized → empty), and the mapper
  (`SourcePlace` → `Venue` incl. cuisine/amenity → tags, Vietnamese slug correctness).
- Service test with a fake `VenueSource` + mock repository: upsert per place,
  resolve-or-create district, `--dry-run` writes nothing, per-place error skipped,
  run counts recorded.
- The project has no DB-backed test harness today; idempotent `ON CONFLICT` upsert is
  verified manually during a district run (re-run a district, confirm counts update
  not duplicate). A DB-backed repository test is deferred to the CI/test-harness
  improvement (P1).

## Out of scope

- Dishes/prices, social videos (B), entity resolution (C), trend scoring/AI (D).
- Rating and photos (no current free source provides them); the `apps/web` UI swap
  away from the fabricated `mediaByVenue` data (still P0 in IMPROVEMENTS.md) — venue
  name/geo/address/tags become real, but rating/photos remain placeholder.
- Provinces other than Ho Chi Minh City.
- A purge/rollback command (idempotent re-run covers re-do).
- A migration-runner tool (P1); migration 006 follows the existing idempotent SQL
  convention.
- Adding Google Places as a second source (future; the `VenueSource` interface keeps
  this cheap).
