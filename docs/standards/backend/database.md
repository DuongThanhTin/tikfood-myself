# Backend Database Standard

TikFood database direction is PostgreSQL + PostGIS.

This schema is the target discovery MVP schema. The current API still uses in-memory data, but future persistence work should follow this document.

## Core Principles

- Model venues, dishes, social proof, search filters, user saves, and collections explicitly.
- Keep delivery, cart, order, checkout, payment, booking, reservation, chat, monetization, and livestream out of the MVP schema.
- Prefer normalized tags because TikFood search/filtering is tag-heavy.
- Store price at the correct level: venue, dish concept, and dish-at-venue.
- Generate AI summaries asynchronously later; do not generate expensive summaries inside map/search requests.

## Entity Overview

```text
venues
dishes
venue_dishes

tags
venue_tags
dish_tags

venue_opening_hours

users
user_saved_items
collections
collection_items

creators
social_posts
social_post_dishes

trend_scores
ai_summaries
hidden_gem_scores
```

## venues

Restaurants, cafes, bars, street-food places, and other food discovery locations.

```sql
venues
- id uuid primary key
- name text not null
- slug text not null unique
- short_description text
- about text
- address text
- city text
- district text
- ward text
- area_name text
- latitude double precision
- longitude double precision
- location geography(Point, 4326)
- phone text
- website_url text
- instagram_url text
- tiktok_url text
- photos text[]
- price_level int
- avg_price_min_vnd int
- avg_price_max_vnd int
- currency text default 'VND'
- social_video_count int default 0
- social_post_count int default 0
- total_view_count bigint default 0
- total_like_count bigint default 0
- verified_at timestamptz
- source text
- created_at timestamptz not null
- updated_at timestamptz not null
```

Field meaning:

- `slug`: URL-friendly identifier for a specific venue, for example `banh-mi-hem-nguyen-trai-district-1`.
- `short_description`: short text for cards and list results.
- `about`: richer restaurant introduction for venue detail pages. This should answer what the place is known for, atmosphere, signature dishes, and why users may want to visit.
- `price_level`: coarse display value from `1` to `4`.
- `avg_price_min_vnd` and `avg_price_max_vnd`: expected price range for a typical person at the venue.

Indexes:

```sql
create index idx_venues_location on venues using gist (location);
create index idx_venues_city_district on venues (city, district);
create index idx_venues_slug on venues (slug);
create index idx_venues_price_max on venues (avg_price_max_vnd);
create index idx_venues_social_video_count on venues (social_video_count);
```

## dishes

Dish concepts independent of any single venue.

```sql
dishes
- id uuid primary key
- name text not null
- normalized_name text not null
- slug text not null unique
- short_description text
- about text
- category text
- cuisine text
- aliases text[]
- photos text[]
- typical_price_min_vnd int
- typical_price_max_vnd int
- currency text default 'VND'
- created_at timestamptz not null
- updated_at timestamptz not null
```

Field meaning:

- `slug`: URL-friendly dish concept, for example `banh-mi-thit-nuong`.
- `normalized_name`: search/matching key, not a URL.
- `about`: richer dish explanation for dish detail pages, including taste profile, ingredients, origin/context, and how users usually find it.

Indexes:

```sql
create index idx_dishes_normalized_name on dishes (normalized_name);
create index idx_dishes_slug on dishes (slug);
```

## venue_dishes

Dish availability and signals for a dish at a specific venue.

```sql
venue_dishes
- venue_id uuid references venues(id)
- dish_id uuid references dishes(id)
- price_min_vnd int
- price_max_vnd int
- currency text default 'VND'
- price_source text
- price_updated_at timestamptz
- confidence_score numeric
- mention_count int default 0
- video_count int default 0
- view_count bigint default 0
- trend_score numeric
- last_seen_at timestamptz
- created_at timestamptz not null
- updated_at timestamptz not null

primary key (venue_id, dish_id)
```

Use this table for queries like:

- pizza near me
- bun bo under 100k
- most-mentioned dish at this venue
- exact saved item: dish at a specific venue

## tags

Normalized tags for scalable filtering.

```sql
tags
- id uuid primary key
- slug text not null unique
- label text not null
- type text not null
- created_at timestamptz not null
```

Allowed `type` values:

```text
cuisine
venue_type
occasion
amenity
vibe
dish_type
ingredient
taste
dietary
area
```

Examples:

```text
italian, japanese, korean
restaurant, bar, rooftop, cafe
date, family, friends, solo
outdoor, parking, air_conditioning
romantic, cozy, lively, quiet
pizza, sushi, bbq, noodle
```

## venue_tags

```sql
venue_tags
- venue_id uuid references venues(id)
- tag_id uuid references tags(id)
- confidence_score numeric
- source text
- created_at timestamptz not null

primary key (venue_id, tag_id)
```

## dish_tags

```sql
dish_tags
- dish_id uuid references dishes(id)
- tag_id uuid references tags(id)
- confidence_score numeric
- source text
- created_at timestamptz not null

primary key (dish_id, tag_id)
```

## venue_opening_hours

Use a normalized table for `open_now`.

```sql
venue_opening_hours
- id uuid primary key
- venue_id uuid references venues(id)
- day_of_week int not null
- open_time time
- close_time time
- is_closed bool default false
- created_at timestamptz not null
- updated_at timestamptz not null
```

`day_of_week`:

```text
0 Sunday
1 Monday
2 Tuesday
3 Wednesday
4 Thursday
5 Friday
6 Saturday
```

Handle overnight opening by allowing `close_time < open_time`.

## users

Application users for saving, sharing, and collections.

```sql
users
- id uuid primary key
- email text unique
- display_name text
- avatar_url text
- created_at timestamptz not null
- updated_at timestamptz not null
```

This supports bookmarks and collections. It does not imply a social follow graph.

## user_saved_items

User bookmarks for venues, dishes, or a dish at a specific venue.

```sql
user_saved_items
- id uuid primary key
- user_id uuid references users(id)
- target_type text not null
- target_id uuid not null
- note text
- created_at timestamptz not null
```

Allowed `target_type`:

```text
venue
dish
venue_dish
```

Use `venue_dish` when the user wants to save "this dish at this place".

## collections

User-created collections.

```sql
collections
- id uuid primary key
- owner_user_id uuid references users(id)
- title text not null
- slug text not null
- description text
- visibility text not null
- collection_type text not null
- city text
- district text
- area_name text
- dish_id uuid references dishes(id)
- cover_image_url text
- created_at timestamptz not null
- updated_at timestamptz not null
```

Allowed `visibility`:

```text
private
public
unlisted
```

Allowed `collection_type`:

```text
district
area
dish
custom
```

Examples:

- Best date spots in District 1
- Pizza places under 500k
- Rooftop bars to try
- Korean food around Thao Dien

## collection_items

```sql
collection_items
- id uuid primary key
- collection_id uuid references collections(id)
- target_type text not null
- target_id uuid not null
- note text
- sort_order int default 0
- created_at timestamptz not null
```

Allowed `target_type`:

```text
venue
dish
venue_dish
```

## creators

Social creators or accounts.

```sql
creators
- id uuid primary key
- platform text not null
- platform_user_id text
- username text
- display_name text
- profile_url text
- follower_count int
- created_at timestamptz not null
- updated_at timestamptz not null

unique (platform, platform_user_id)
```

## social_posts

Social proof posts.

```sql
social_posts
- id uuid primary key
- platform text not null
- platform_post_id text
- creator_id uuid references creators(id)
- venue_id uuid references venues(id)
- post_url text
- caption text
- posted_at timestamptz
- like_count int
- comment_count int
- share_count int
- view_count int
- raw_engagement_score numeric
- created_at timestamptz not null

unique (platform, platform_post_id)
```

## social_post_dishes

Detected dish mentions in social posts.

```sql
social_post_dishes
- social_post_id uuid references social_posts(id)
- dish_id uuid references dishes(id)
- confidence_score numeric
- source text

primary key (social_post_id, dish_id)
```

## trend_scores

Time-based trend snapshots.

```sql
trend_scores
- id uuid primary key
- venue_id uuid references venues(id)
- dish_id uuid references dishes(id)
- scope text not null
- score numeric not null
- signal_count int
- time_window text
- calculated_at timestamptz not null
```

Allowed `scope`:

```text
venue
dish
venue_dish
district
city
```

## ai_summaries

Pre-generated summaries.

```sql
ai_summaries
- id uuid primary key
- target_type text not null
- target_id uuid not null
- summary text not null
- model text
- source_signal_count int
- generated_at timestamptz not null
- expires_at timestamptz
```

Allowed `target_type`:

```text
venue
dish
venue_dish
collection
```

## hidden_gem_scores

```sql
hidden_gem_scores
- id uuid primary key
- venue_id uuid references venues(id)
- score numeric not null
- reason text
- calculated_at timestamptz not null
```

## Explicitly Not In MVP Schema

Do not add:

```text
orders
carts
payments
checkout_sessions
bookings
reservations
delivery_jobs
chat_messages
creator_payouts
livestreams
user_follows
```
