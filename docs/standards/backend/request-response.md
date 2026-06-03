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

The current MVP code still uses a string `error` field. Future backend changes should migrate to the object shape above when introducing the first real error response.

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
GET /api/v1/map/venues?district=District%201&dish=pho
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
GET /api/v1/map/venues
```

Query params:

- `district`: optional string
- `dish`: optional string

Response data item:

```json
{
  "id": "venue_001",
  "name": "Banh Mi Hem",
  "slug": "banh-mi-hem-nguyen-trai-district-1",
  "short_description": "Late-night banh mi spot trending on social video.",
  "about": "A compact street-food venue known for grilled pork banh mi, quick service, and strong late-night local buzz.",
  "address": "12 Nguyen Trai",
  "district": "District 1",
  "latitude": 10.7712,
  "longitude": 106.6899,
  "categories": ["street_food", "banh_mi"],
  "trend_score": 92,
  "trending_dishes": ["banh mi thit nuong"],
  "ai_summary": "Trending for late-night banh mi clips."
}
```

## Contract Change Rules

If response shape changes:

- Update backend structs.
- Update frontend types in `apps/web/lib/api.ts`.
- Update docs.
- Update examples if relevant.
- Add or update tests.
