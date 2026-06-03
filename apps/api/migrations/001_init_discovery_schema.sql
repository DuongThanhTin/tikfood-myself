create extension if not exists postgis;
create extension if not exists pgcrypto;

create table if not exists venues (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  slug text not null unique,
  short_description text,
  about text,
  address text,
  city text,
  district text,
  ward text,
  area_name text,
  latitude double precision,
  longitude double precision,
  location geography(Point, 4326),
  phone text,
  website_url text,
  instagram_url text,
  tiktok_url text,
  photos text[] not null default '{}',
  price_level int,
  avg_price_min_vnd int,
  avg_price_max_vnd int,
  currency text not null default 'VND',
  social_video_count int not null default 0,
  social_post_count int not null default 0,
  total_view_count bigint not null default 0,
  total_like_count bigint not null default 0,
  verified_at timestamptz,
  source text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists dishes (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  normalized_name text not null,
  slug text not null unique,
  short_description text,
  about text,
  category text,
  cuisine text,
  aliases text[] not null default '{}',
  photos text[] not null default '{}',
  typical_price_min_vnd int,
  typical_price_max_vnd int,
  currency text not null default 'VND',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists venue_dishes (
  venue_id uuid not null references venues(id) on delete cascade,
  dish_id uuid not null references dishes(id) on delete cascade,
  price_min_vnd int,
  price_max_vnd int,
  currency text not null default 'VND',
  price_source text,
  price_updated_at timestamptz,
  confidence_score numeric,
  mention_count int not null default 0,
  video_count int not null default 0,
  view_count bigint not null default 0,
  trend_score numeric,
  last_seen_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  primary key (venue_id, dish_id)
);

create table if not exists tags (
  id uuid primary key default gen_random_uuid(),
  slug text not null unique,
  label text not null,
  type text not null,
  created_at timestamptz not null default now()
);

create table if not exists venue_tags (
  venue_id uuid not null references venues(id) on delete cascade,
  tag_id uuid not null references tags(id) on delete cascade,
  confidence_score numeric,
  source text,
  created_at timestamptz not null default now(),
  primary key (venue_id, tag_id)
);

create table if not exists dish_tags (
  dish_id uuid not null references dishes(id) on delete cascade,
  tag_id uuid not null references tags(id) on delete cascade,
  confidence_score numeric,
  source text,
  created_at timestamptz not null default now(),
  primary key (dish_id, tag_id)
);

create table if not exists venue_opening_hours (
  id uuid primary key default gen_random_uuid(),
  venue_id uuid not null references venues(id) on delete cascade,
  day_of_week int not null check (day_of_week between 0 and 6),
  open_time time,
  close_time time,
  is_closed bool not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists trend_scores (
  id uuid primary key default gen_random_uuid(),
  venue_id uuid references venues(id) on delete cascade,
  dish_id uuid references dishes(id) on delete cascade,
  scope text not null,
  score numeric not null,
  signal_count int,
  time_window text,
  calculated_at timestamptz not null default now()
);

create table if not exists ai_summaries (
  id uuid primary key default gen_random_uuid(),
  target_type text not null,
  target_id uuid not null,
  summary text not null,
  model text,
  prompt_version text,
  source_signal_count int,
  generated_at timestamptz not null default now(),
  expires_at timestamptz
);

create index if not exists idx_venues_location on venues using gist (location);
create index if not exists idx_venues_city_district on venues (city, district);
create index if not exists idx_venues_slug on venues (slug);
create index if not exists idx_venues_price_max on venues (avg_price_max_vnd);
create index if not exists idx_venues_social_video_count on venues (social_video_count);
create index if not exists idx_dishes_normalized_name on dishes (normalized_name);
create index if not exists idx_dishes_slug on dishes (slug);
create index if not exists idx_tags_type_slug on tags (type, slug);
create index if not exists idx_venue_opening_hours_venue_day on venue_opening_hours (venue_id, day_of_week);
create index if not exists idx_trend_scores_venue_calculated_at on trend_scores (venue_id, calculated_at desc);
create index if not exists idx_trend_scores_dish_calculated_at on trend_scores (dish_id, calculated_at desc);
create index if not exists idx_ai_summaries_target on ai_summaries (target_type, target_id, generated_at desc);
