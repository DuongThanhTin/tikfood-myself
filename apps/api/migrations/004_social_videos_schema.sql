create table if not exists social_videos (
  id uuid primary key default gen_random_uuid(),
  venue_id uuid not null references venues(id) on delete cascade,
  dish_id uuid references dishes(id) on delete set null,
  platform text not null check (platform in ('tiktok', 'instagram', 'youtube', 'facebook', 'other')),
  source_url text not null unique,
  source_external_id text,
  creator_handle text,
  creator_display_name text,
  caption text,
  thumbnail_url text,
  view_count bigint not null default 0,
  like_count bigint not null default 0,
  comment_count bigint not null default 0,
  share_count bigint not null default 0,
  published_at timestamptz,
  fetched_at timestamptz not null default now(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_social_videos_venue_platform on social_videos (venue_id, platform);
create index if not exists idx_social_videos_dish_platform on social_videos (dish_id, platform);
create index if not exists idx_social_videos_view_count on social_videos (view_count desc);
create index if not exists idx_social_videos_published_at on social_videos (published_at desc);
