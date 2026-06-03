# Backend Logging Standard

Backend logs must be structured, small, and safe.

## Current State

Current API uses Go standard library `log/slog` with JSON output.

Gin request logging is implemented as middleware in `internal/http`.

## Log Format

Preferred production format: JSON.

Example:

```json
{
  "time": "2026-06-03T10:00:00Z",
  "level": "INFO",
  "service": "api",
  "request_id": "req_01",
  "method": "GET",
  "path": "/api/v1/map/venues",
  "status": 200,
  "duration_ms": 12
}
```

## Required Request Log Fields

- `service`
- `request_id`
- `method`
- `path`
- `status`
- `duration_ms`
- `user_id` only after auth exists and only if safe

## Request ID

- Accept `X-Request-ID` if provided and valid.
- Generate a request id if missing.
- Return `X-Request-ID` in response headers.
- Include request id in logs.

Do not add a request ID dependency.

## Log Levels

Use:

- `DEBUG`: local development detail only
- `INFO`: normal request lifecycle and process events
- `WARN`: recoverable unexpected behavior
- `ERROR`: failed operations requiring attention

Avoid noisy logs in hot paths.

## What Not To Log

Never log:

- `.env` values
- API keys
- tokens
- passwords
- cookies
- Authorization headers
- private keys
- raw social platform payloads without approval
- full request bodies by default

## Domain Logging

Domain services should not spam logs.

Log at domain boundaries only when useful:

- trend scoring job started/completed
- ingestion batch summary
- summary generation result
- storage operation failure

## Error Logging

For errors, log:

- stable error code
- route
- request id
- safe message
- wrapped internal error only if it contains no sensitive data

Do not return internal error text to clients.

## Metrics Direction

Do not add metrics dependencies yet.

Future metrics:

- request count
- request latency
- response status count
- venue query count
- trend scoring job duration
- summary job duration
- ingestion failure count
