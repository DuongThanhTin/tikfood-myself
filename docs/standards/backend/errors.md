# Backend Error Handling Standard

Errors must be explicit, safe, and useful.

## Error Principles

- Do not panic for request-level failures.
- Do not expose raw internal errors to clients.
- Do not log secrets.
- Use stable error codes.
- Preserve original errors internally where useful.

## Error Categories

Use these stable codes:

```text
invalid_request
not_found
conflict
unauthorized
forbidden
domain_rejected
internal_error
service_unavailable
```

MVP anti-goal requests should use:

```text
domain_rejected
```

## HTTP Mapping

```text
invalid_request      -> 400
unauthorized         -> 401
forbidden            -> 403
not_found            -> 404
conflict             -> 409
domain_rejected      -> 422
internal_error       -> 500
service_unavailable  -> 503
```

## Error Response

Future standard:

```json
{
  "data": null,
  "error": {
    "code": "invalid_request",
    "message": "District is invalid."
  }
}
```

Optional details:

```json
{
  "data": null,
  "error": {
    "code": "invalid_request",
    "message": "Request validation failed.",
    "details": {
      "field": "district"
    }
  }
}
```

Do not include:

- stack traces
- SQL strings
- database hostnames
- tokens
- request headers containing auth
- internal filesystem paths

## Domain Errors

Domain services should return meaningful errors that HTTP can map.

Example future shape:

```go
type Code string

const (
    CodeInvalidRequest Code = "invalid_request"
    CodeDomainRejected Code = "domain_rejected"
)

type AppError struct {
    Code    Code
    Message string
    Err     error
}
```

Do not introduce this abstraction until there are enough error paths to justify it.

## Logging Errors

Log internal error context server-side, but sanitize:

Allowed:

- request id
- route
- status
- stable error code
- safe message

Not allowed:

- secrets
- `.env` values
- full credentials
- raw auth headers
- raw third-party payloads

## Testing Errors

Add tests for:

- invalid query params
- invalid request body
- domain rejection
- not found cases
- storage error mapping when storage exists
