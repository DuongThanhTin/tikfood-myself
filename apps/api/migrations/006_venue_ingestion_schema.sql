-- venues: rating fields kept for a future Google Places source (OSM leaves them NULL)
alter table venues add column if not exists google_rating numeric;
alter table venues add column if not exists google_rating_count int;

-- venues: upsert key for source ingestion
alter table venues drop constraint if exists venues_source_external_unique;
alter table venues add constraint venues_source_external_unique unique (source, source_external_id);

-- venues: keep PostGIS location in sync with lat/lng via a generated column.
-- Existing data is seed-only, so dropping and re-adding is safe.
drop index if exists idx_venues_location;
alter table venues drop column if exists location;
alter table venues add column location geography(Point, 4326)
  generated always as (
    case
      when longitude is null or latitude is null then null
      else ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography
    end
  ) stored;
create index if not exists idx_venues_location on venues using gist (location);

-- ingestion audit log (one row per district run)
create table if not exists ingestion_runs (
  id uuid primary key default gen_random_uuid(),
  source text not null,
  district_slug text not null,
  category text not null,
  status text not null check (status in ('running', 'completed', 'failed')),
  fetched_count int not null default 0,
  upserted_count int not null default 0,
  failed_count int not null default 0,
  error text,
  started_at timestamptz not null default now(),
  finished_at timestamptz
);

create index if not exists idx_ingestion_runs_district on ingestion_runs (district_slug, started_at desc);
