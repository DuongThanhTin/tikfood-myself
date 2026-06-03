# API Contract Standard

This standard defines request and response conventions for TikFood APIs.

## Base Principles

- API contracts must be explicit, typed, and documented.
- Backend structs, frontend types, examples, and schemas should agree.
- Use `snake_case` for JSON fields.
- Do not expose internal implementation details.
- Do not expose secrets, tokens, stack traces, or raw SQL errors.

## Response Envelope

Successful response:

```json
{
  "data": {}
}
```

List response:

```json
{
  "data": []
}
```

Error response:

```json
{
  "data": null,
  "error": "Human-readable error message"
}
```

Current API code uses this envelope shape. Keep it consistent.

## HTTP Status Codes

- `200`: successful read
- `201`: successful create
- `204`: successful delete with no body
- `400`: invalid request input
- `401`: unauthenticated
- `403`: unauthorized
- `404`: resource not found
- `409`: conflict
- `422`: valid JSON but invalid domain request
- `500`: internal server error

## Request Rules

- Query params are allowed for simple filters.
- JSON bodies are preferred for complex creates/updates.
- Validate request input at the HTTP boundary.
- Normalize and sanitize user-provided strings before domain logic.
- Do not accept payment, booking, order, cart, checkout, or delivery request shapes in MVP.

## Current Endpoint

```text
GET /api/v1/map/venues
```

Query params:

```text
district optional string
dish optional string
```

Response:

```json
{
  "data": [
    {
      "id": "venue_001",
      "name": "Banh Mi Hem",
      "address": "12 Nguyen Trai",
      "district": "District 1",
      "latitude": 10.7712,
      "longitude": 106.6899,
      "categories": ["street_food", "banh_mi"],
      "trend_score": 92,
      "trending_dishes": ["banh mi thit nuong"],
      "ai_summary": "Trending for late-night banh mi clips."
    }
  ]
}
```

## Versioning

Use versioned routes:

```text
/api/v1/...
```

Do not break existing response fields without adding a migration path.

## Frontend Type Sync

Frontend types live in:

```text
apps/web/lib/api.ts
```

When backend response fields change, update frontend types in the same change.

## Schema Direction

Public or runner-facing contracts should be documented under:

```text
packages/schemas/
```

For product API contracts, add OpenAPI docs when endpoints stabilize beyond MVP demo data.
