create extension if not exists unaccent;
create extension if not exists pg_trgm;

create table if not exists locations (
  id uuid primary key default gen_random_uuid(),
  type text not null check (type in ('country', 'city', 'district', 'ward', 'area')),
  parent_id uuid references locations(id),
  official_name text not null,
  normalized_name text not null,
  code text,
  slug text not null,
  country_code text not null default 'VN',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (type, country_code, slug)
);

create table if not exists location_aliases (
  id uuid primary key default gen_random_uuid(),
  location_id uuid not null references locations(id) on delete cascade,
  alias text not null,
  normalized_alias text not null,
  locale text,
  confidence numeric not null default 1,
  created_at timestamptz not null default now(),
  unique (location_id, normalized_alias)
);

alter table venues
  add column if not exists public_id text unique,
  add column if not exists source_external_id text,
  add column if not exists source_updated_at timestamptz,
  add column if not exists address_raw text,
  add column if not exists address_display text,
  add column if not exists city_location_id uuid references locations(id),
  add column if not exists district_location_id uuid references locations(id),
  add column if not exists ward_location_id uuid references locations(id),
  add column if not exists area_location_id uuid references locations(id);

update venues
set
  address_raw = coalesce(address_raw, address),
  address_display = coalesce(address_display, address)
where address is not null
  and (address_raw is null or address_display is null);

insert into locations (id, type, parent_id, official_name, normalized_name, code, slug, country_code) values
  ('10000000-0000-0000-0000-000000000001', 'country', null, 'Việt Nam', 'viet nam', 'VN', 'viet-nam', 'VN'),
  ('10000000-0000-0000-0000-000000000100', 'city', '10000000-0000-0000-0000-000000000001', 'Thành phố Hồ Chí Minh', 'thanh pho ho chi minh', 'SG', 'ho-chi-minh', 'VN'),
  ('10000000-0000-0000-0000-000000000101', 'district', '10000000-0000-0000-0000-000000000100', 'Quận 1', 'quan 1', 'Q1', 'quan-1', 'VN'),
  ('10000000-0000-0000-0000-000000000103', 'district', '10000000-0000-0000-0000-000000000100', 'Quận 3', 'quan 3', 'Q3', 'quan-3', 'VN')
on conflict (type, country_code, slug) do update
set
  official_name = excluded.official_name,
  normalized_name = excluded.normalized_name,
  code = excluded.code,
  parent_id = excluded.parent_id,
  updated_at = now();

with aliases (location_id, alias, locale, confidence, priority) as (
  values
    ('10000000-0000-0000-0000-000000000100', 'Thành phố Hồ Chí Minh', 'vi', 1, 10),
    ('10000000-0000-0000-0000-000000000100', 'TP. Hồ Chí Minh', 'vi', 1, 20),
    ('10000000-0000-0000-0000-000000000100', 'Tp Hồ Chí Minh', 'vi', 1, 30),
    ('10000000-0000-0000-0000-000000000100', 'Hồ Chí Minh', 'vi', 1, 40),
    ('10000000-0000-0000-0000-000000000100', 'Ho Chi Minh City', 'en', 1, 50),
    ('10000000-0000-0000-0000-000000000100', 'Ho Chi Minh', 'en', 1, 60),
    ('10000000-0000-0000-0000-000000000100', 'HCM', null, 1, 70),
    ('10000000-0000-0000-0000-000000000100', 'TP HCM', 'vi', 1, 80),
    ('10000000-0000-0000-0000-000000000100', 'Tp HCM', 'vi', 1, 90),
    ('10000000-0000-0000-0000-000000000101', 'Quận 1', 'vi', 1, 10),
    ('10000000-0000-0000-0000-000000000101', 'Quan 1', 'vi', 1, 20),
    ('10000000-0000-0000-0000-000000000101', 'District 1', 'en', 1, 30),
    ('10000000-0000-0000-0000-000000000101', 'Q1', null, 1, 40),
    ('10000000-0000-0000-0000-000000000101', 'Q. 1', null, 1, 50),
    ('10000000-0000-0000-0000-000000000101', 'Q.1', null, 1, 60),
    ('10000000-0000-0000-0000-000000000103', 'Quận 3', 'vi', 1, 10),
    ('10000000-0000-0000-0000-000000000103', 'Quan 3', 'vi', 1, 20),
    ('10000000-0000-0000-0000-000000000103', 'District 3', 'en', 1, 30),
    ('10000000-0000-0000-0000-000000000103', 'Q3', null, 1, 40),
    ('10000000-0000-0000-0000-000000000103', 'Q. 3', null, 1, 50),
    ('10000000-0000-0000-0000-000000000103', 'Q.3', null, 1, 60)
),
deduped_aliases as (
  select distinct on (location_id, lower(unaccent(alias)))
    location_id::uuid,
    alias,
    lower(unaccent(alias)) as normalized_alias,
    locale,
    confidence
  from aliases
  order by location_id, lower(unaccent(alias)), priority
)
insert into location_aliases (location_id, alias, normalized_alias, locale, confidence)
select location_id, alias, normalized_alias, locale, confidence
from deduped_aliases
on conflict (location_id, normalized_alias) do update
set
  alias = excluded.alias,
  locale = excluded.locale,
  confidence = excluded.confidence;

update venues
set
  city_location_id = '10000000-0000-0000-0000-000000000100',
  district_location_id = case
    when lower(unaccent(coalesce(district, ''))) in ('district 1', 'quan 1', 'q1', 'q. 1', 'q.1') then '10000000-0000-0000-0000-000000000101'::uuid
    when lower(unaccent(coalesce(district, ''))) in ('district 3', 'quan 3', 'q3', 'q. 3', 'q.3') then '10000000-0000-0000-0000-000000000103'::uuid
    else district_location_id
  end
where city_location_id is null
  and lower(unaccent(coalesce(city, ''))) in (
    'thanh pho ho chi minh',
    'tp. ho chi minh',
    'tp ho chi minh',
    'ho chi minh city',
    'ho chi minh',
    'hcm',
    'tp hcm'
  );

create index if not exists idx_locations_type_slug on locations (type, slug);
create index if not exists idx_locations_parent_type on locations (parent_id, type);
create index if not exists idx_location_aliases_normalized on location_aliases (normalized_alias);
create index if not exists idx_location_aliases_trgm on location_aliases using gin (normalized_alias gin_trgm_ops);
create index if not exists idx_venues_public_id on venues (public_id);
create index if not exists idx_venues_source_external_id on venues (source, source_external_id);
create index if not exists idx_venues_city_location on venues (city_location_id);
create index if not exists idx_venues_district_location on venues (district_location_id);
