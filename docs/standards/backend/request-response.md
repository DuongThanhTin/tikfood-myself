# Backend Request And Response Standard

This file defines backend request parsing and response shape.

## JSON Field Style

Use `snake_case` for API JSON.

Good:

```json
{
  "trend_score": 92,
  "trending_dishes": []
}
```

Avoid:

```json
{
  "trendScore": 92,
  "trendingDishes": []
}
```

Go structs should use explicit JSON tags.

## Response Envelope

All JSON responses use this envelope:

Success object:

```json
{
  "data": {}
}
```

Success list:

```json
{
  "data": []
}
```

Error:

```json
{
  "data": null,
  "error": {
    "code": "invalid_request",
    "message": "District is invalid."
  }
}
```

The backend must use the object shape above for errors.

HTTP handlers must use response helpers from `internal/http/response.go` instead of writing envelopes manually:

```text
respondWithData
respondWithBadRequest
respondWithError
respondWithInternalServerError
```

## Error Object

Standard error shape:

```json
{
  "code": "invalid_request",
  "message": "Human-readable message.",
  "details": {}
}
```

Rules:

- `code` is stable and machine-readable.
- `message` is safe for client display.
- `details` is optional and must not contain secrets or internals.

## Status Codes

Use:

- `200 OK`: successful read
- `201 Created`: successful create
- `204 No Content`: successful delete or empty success
- `400 Bad Request`: malformed or invalid request
- `401 Unauthorized`: missing or invalid authentication
- `403 Forbidden`: authenticated but not allowed
- `404 Not Found`: missing resource
- `409 Conflict`: duplicate or conflicting state
- `422 Unprocessable Entity`: valid request shape but invalid domain action
- `500 Internal Server Error`: unexpected server failure

## Query Parameters

Use query params for simple filters:

```text
GET /api/v1/discovery/venues?q=pho&district=District%201&tags=pho&max_price_vnd=120000
```

Rules:

- Parse query params in HTTP handlers.
- Normalize strings before passing to domain services.
- Validate allowed values when finite.
- Do not let query params directly construct SQL.

## Request Bodies

Use JSON bodies for create/update operations.

Rules:

- Limit request body size before decoding when public endpoints are added.
- Decode once.
- Reject unknown fields for stable public APIs when appropriate.
- Validate required fields explicitly.
- Convert request DTOs into domain input structs.

## Response DTOs

Do not expose database rows directly as API responses once storage exists.

Preferred flow:

```text
storage row
-> domain model
-> response DTO
```

For MVP in-memory data, direct domain response is acceptable if the domain model is already the API shape.

## Pagination

When list endpoints grow beyond MVP:

Request:

```text
limit optional int, max 100
cursor optional string
```

Response:

```json
{
  "data": [],
  "page": {
    "next_cursor": "string",
    "has_more": true
  }
}
```

Do not add pagination until an endpoint can realistically return large data.

## Current Endpoint Contract

```text
GET /api/v1/discovery/venues
```

Query params:

- `q`: optional text search across venue, dish, cuisine, and tags
- `city`: optional string, MVP compatibility alias such as `HCM`, `Ho Chi Minh`, `Tp Hồ Chí Minh`
- `district`: optional string, MVP compatibility alias such as `District 1`, `Quận 1`, `Q1`
- `dish`: optional string
- `tags`: optional comma-separated tag slugs
- `platform`: optional comma-separated social platform slugs such as `tiktok`, `instagram`
- `lat`: optional float, requires `lng`
- `lng`: optional float, requires `lat`
- `radius_m`: optional int, max 50000
- `min_price_vnd`: optional int
- `max_price_vnd`: optional int
- `open_now`: optional bool
- `sort`: optional enum, one of `trending`, `videos`, `distance`, `price`
- `limit`: optional int, 1-100

Compatibility endpoint:

```text
GET /api/v1/map/venues
```

It currently uses the same handler as `/api/v1/discovery/venues`.

List response data item:

```json
{
  "id": "7b7e7ab7-c7b5-4ef8-8cb9-6330b9d2cf55",
  "name": "Banh Mi Hem",
  "slug": "banh-mi-hem-nguyen-trai-district-1",
  "short_description": "Late-night banh mi spot trending on social video.",
  "about": "A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz.",
  "address": "12 Nguyen Trai",
  "district": "Quận 1",
  "latitude": 10.7712,
  "longitude": 106.6899,
  "categories": ["street_food", "banh_mi"],
  "price_level": 1,
  "avg_price_min_vnd": 30000,
  "avg_price_max_vnd": 80000,
  "currency": "VND",
  "social_video_count": 42,
  "social_videos": [
    {
      "id": "0ec82711-4541-4680-866f-4462ce4e7d38",
      "platform": "tiktok",
      "url": "https://www.tiktok.com/@tikfood/video/banh-mi-hem-1",
      "creator_handle": "@tikfood",
      "caption": "Late-night banh mi with grilled pork near Nguyen Trai.",
      "thumbnail_url": "",
      "view_count": 120000,
      "like_count": 8200,
      "published_at": "2026-05-20T10:00:00Z"
    }
  ],
  "trend_score": 92,
  "trending_dishes": ["banh mi thit nuong"],
  "ai_summary": "Trending for late-night banh mi clips.",
  "distance_meters": 720.4
}
```

Detail endpoint:

```text
GET /api/v1/discovery/venues/:slug
```

Use this endpoint when a user opens the right panel/detail view. The list endpoint can remain lighter for map performance, while detail includes richer venue data.

Detail response data item adds:

```json
{
  "dishes": [
    {
      "id": "f57581ad-4144-4f85-a995-78f729d35e6b",
      "name": "Banh mi thit nuong",
      "slug": "banh-mi-thit-nuong",
      "short_description": "Vietnamese grilled pork banh mi.",
      "about": "A crispy baguette filled with grilled pork, pickles, herbs, pate, and sauces.",
      "category": "sandwich",
      "cuisine": "vietnamese",
      "price_min_vnd": 35000,
      "price_max_vnd": 60000,
      "currency": "VND",
      "mention_count": 24,
      "video_count": 18,
      "view_count": 120000,
      "trend_score": 92
    }
  ],
  "opening_hours": [
    {
      "day_of_week": 0,
      "open_time": "08:00",
      "close_time": "23:30",
      "is_closed": false
    }
  ]
}
```

If the slug does not exist, return:

```json
{
  "data": null,
  "error": {
    "code": "not_found",
    "message": "Venue was not found."
  }
}
```

`id` is a database-generated UUID. Do not expose manually configured identifiers as the primary ID. If the product needs a stable non-UUID public identifier, add `public_id`.

Current MVP clients may send raw `district` text. Backend repositories must normalize that value through `location_aliases` before filtering. Future clients should prefer location slugs or IDs:

```text
district_slug=quan-1
city_slug=ho-chi-minh
```

## Contract Change Rules

If response shape changes:

- Update backend structs.
- Update frontend types in `apps/web/lib/api.ts`.
- Update docs.
- Update examples if relevant.
- Add or update tests.
