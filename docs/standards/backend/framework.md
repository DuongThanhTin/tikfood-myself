# Backend Framework Standard

TikFood backend uses Go with Gin.

## Standard Choice

Framework:

```text
github.com/gin-gonic/gin
```

Gin is the standard HTTP framework for `apps/api`.

Do not add another router or HTTP framework unless this standard is explicitly changed.

## Why Gin

Gin is selected because TikFood needs:

- clear route groups for versioned APIs
- middleware for request id, recovery, logging, auth, and CORS
- stable `net/http` interoperability
- straightforward testing with `httptest`
- familiar patterns for Go backend engineers and AI coding agents

## Why Not Fiber By Default

Fiber is not the default because it is built on `fasthttp` instead of `net/http`.

TikFood prioritizes:

- Go ecosystem compatibility
- observability integration
- simple testability
- long-term maintainability

Raw router throughput is not expected to be the first bottleneck. Database queries, ranking, search, cache, and ingestion jobs are more likely to dominate system performance.

## Router Rules

Route registration belongs in:

```text
apps/api/internal/http
```

Use versioned route groups:

```text
/api/v1
/api/v1/discovery
/api/v1/users
/api/v1/collections
```

Current MVP compatibility route:

```text
/api/v1/map/venues
```

Long-term discovery routes should move toward:

```text
/api/v1/discovery/venues
/api/v1/discovery/dishes
```

## Middleware Rules

Required middleware:

- panic recovery
- request id
- structured request logging

Future middleware:

- CORS
- auth
- rate limit
- request body size limit
- timeout

Middleware must be small and framework-specific. Business logic must not live in middleware.

## Handler Rules

Gin handlers may:

- parse path/query/body values
- validate request DTOs
- call domain services
- map domain errors to HTTP errors
- write response envelopes

Gin handlers must not:

- run SQL
- calculate trend scores
- generate AI summaries
- call social APIs directly
- contain large business rule branches

## Testing Rules

HTTP tests should continue using:

```text
net/http/httptest
```

Do not require a live server for handler tests.
