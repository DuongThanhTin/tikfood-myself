# Backend Search And Filtering Standard

TikFood search is a blend of text, tags, geo, price, opening hours, and social proof.

This document defines how queries like these should be handled:

```text
pizza
Italian restaurant
Japanese restaurant
Korean food
open now
near me
most videos
date under 500k in District 1
rooftop bar
```

## Search API

Current venue discovery endpoint:

```text
GET /api/v1/discovery/venues
```

Query params:

```text
q optional string
lat optional float
lng optional float
radius_m optional int, max 50000
district optional string
city optional string
open_now optional bool
max_price_vnd optional int
tags optional comma-separated string
dish optional string
sort optional string
limit optional int
```

Allowed `sort` values:

```text
trending
videos
distance
price
```

Compatibility endpoint:

```text
GET /api/v1/map/venues
```

This route uses the same filter handler for now. New clients should use `/api/v1/discovery/venues`.

## Query Interpretation

Example:

```text
date under 500k in District 1 rooftop bar
```

Normalized query:

```json
{
  "occasion": ["date"],
  "max_price_vnd": 500000,
  "city_location_id": "uuid:ho-chi-minh",
  "district_location_id": "uuid:quan-1",
  "amenity": ["rooftop"],
  "venue_type": ["bar"]
}
```

Example:

```text
pizza Italian near me open now most videos
```

Normalized query:

```json
{
  "dish": "pizza",
  "cuisine": ["italian"],
  "near_me": true,
  "open_now": true,
  "sort": "videos"
}
```

MVP can start with explicit UI filters. Natural-language query parsing can be added later.

## Location Normalization

City, district, ward, and area input must be resolved before filtering venues.

Do not treat raw strings such as `Ho Chi Minh City`, `HCM`, `Quận 1`, or `District 1` as canonical database values. They are aliases for canonical rows in `locations`.

Examples:

```text
Ho Chi Minh City
Ho Chi Minh
HCM
Hồ Chí Minh
Tp HCM
Tp Hồ Chí Minh
TP. Hồ Chí Minh
-> locations.slug = ho-chi-minh

District 1
Quận 1
Quan 1
Q1
Q. 1
Q.1
-> locations.slug = quan-1
```

Resolution flow:

```text
raw input
-> trim
-> lowercase
-> unaccent
-> exact match location_aliases.normalized_alias
-> fallback trigram match when confidence is high enough
-> use locations.id in venue filters
```

API handlers may accept raw `city` and `district` strings for MVP compatibility, but repository queries should translate them to canonical location IDs or match through `location_aliases`.

Preferred future API shape:

```text
city_id optional uuid
city_slug optional string
district_id optional uuid
district_slug optional string
area_id optional uuid
area_slug optional string
```

UI should prefer sending slugs or IDs from a location lookup endpoint instead of arbitrary display text.

## Filter Mapping

Dish search:

```text
q=pizza or dish=pizza
-> dishes.normalized_name
-> dish_tags
-> venue_dishes
-> venues
```

Cuisine:

```text
cuisine=italian
-> tags.type = cuisine
-> venue_tags or dish_tags
```

Venue type:

```text
venue_type=bar
venue_type=restaurant
venue_type=rooftop
-> tags.type = venue_type
-> venue_tags
```

Occasion:

```text
occasion=date
occasion=family
occasion=friends
-> tags.type = occasion
-> venue_tags
```

Amenities:

```text
amenity=rooftop
amenity=outdoor
amenity=parking
-> tags.type = amenity
-> venue_tags
```

Price:

```text
max_price_vnd=500000
-> venues.avg_price_max_vnd <= 500000
```

Dish-specific price:

```text
dish=pizza&max_price_vnd=300000
-> venue_dishes.price_max_vnd <= 300000
```

Open now:

```text
open_now=true
-> venue_opening_hours
```

Near me:

```text
lat=10.77&lng=106.69&radius_m=3000
-> PostGIS ST_DWithin
```

City/district:

```text
district=Q1
district=Quận 1
district=District 1
-> normalize with unaccent/lowercase
-> location_aliases.normalized_alias
-> venues.district_location_id
```

Most videos:

```text
sort=videos
-> venues.social_video_count desc
```

Dish-specific video count:

```text
dish=pizza&sort=videos
-> venue_dishes.video_count desc
```

## SQL Direction

Geo filter:

```sql
where ST_DWithin(
  venues.location,
  ST_SetSRID(ST_MakePoint(:lng, :lat), 4326)::geography,
  :radius_m
)
```

Geo sort:

```sql
order by ST_Distance(
  venues.location,
  ST_SetSRID(ST_MakePoint(:lng, :lat), 4326)::geography
)
```

Tag filter:

```sql
exists (
  select 1
  from venue_tags
  join tags on tags.id = venue_tags.tag_id
  where venue_tags.venue_id = venues.id
    and tags.type = 'occasion'
    and tags.slug = 'date'
)
```

Location alias filter:

```sql
exists (
  select 1
  from location_aliases district_alias
  where district_alias.location_id = venues.district_location_id
    and district_alias.normalized_alias = lower(unaccent(:district))
)
```

Dish-at-venue filter:

```sql
exists (
  select 1
  from venue_dishes
  join dishes on dishes.id = venue_dishes.dish_id
  where venue_dishes.venue_id = venues.id
    and dishes.normalized_name = :dish
)
```

## Response Shape

Discovery search response:

```json
{
  "data": [
    {
      "type": "venue",
      "venue": {
        "id": "7b7e7ab7-c7b5-4ef8-8cb9-6330b9d2cf55",
        "name": "Rooftop Pizza Bar",
        "slug": "rooftop-pizza-bar-district-1",
        "short_description": "Rooftop Italian bar for pizza and date nights.",
        "about": "A rooftop Italian-inspired venue known for pizza, city views, evening dates, and social video buzz.",
        "district": "Quận 1",
        "latitude": 10.77,
        "longitude": 106.69,
        "price_level": 3,
        "avg_price_min_vnd": 250000,
        "avg_price_max_vnd": 500000,
        "is_open_now": true,
        "social_video_count": 128,
        "trend_score": 91
      },
      "matched_dishes": [
        {
          "id": "a25a4a0e-0e9e-4d5a-8c16-34b3c7793948",
          "name": "Pizza Margherita",
          "slug": "pizza-margherita",
          "price_min_vnd": 180000,
          "price_max_vnd": 260000,
          "video_count": 42
        }
      ],
      "match_reasons": [
        "Matches pizza",
        "Italian cuisine",
        "Rooftop bar",
        "Good for date",
        "Under 500k",
        "Open now"
      ]
    }
  ],
  "page": {
    "next_cursor": null,
    "has_more": false
  }
}
```

## Ranking Direction

MVP score can combine:

```text
text_match_score
tag_match_score
distance_score
trend_score
social_video_score
price_fit_score
open_now_boost
```

Start simple:

```text
filter first
sort by requested sort
fallback sort by trend_score desc
```

Later:

```text
weighted ranking
personalized saved/collection boosts
semantic search
OpenSearch
```

## Backend Code Placement

Future files:

```text
apps/api/internal/discovery/search_query.go
apps/api/internal/discovery/search_result.go
apps/api/internal/discovery/search_service.go
apps/api/internal/http/discovery_search_handler.go
apps/api/internal/storage/discovery_repository.go
```

Rules:

- Parse HTTP query in `internal/http`.
- Normalize filter input before domain calls.
- Keep ranking logic in `internal/discovery`.
- Keep SQL in `internal/storage`.
- Keep response DTO mapping explicit.
