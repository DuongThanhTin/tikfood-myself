# Backend Pattern Standard

This file defines the backend patterns that AI agents must follow.

## Primary Pattern

TikFood backend uses:

```text
Handler -> Service -> Repository
```

With explicit dependency wiring in `cmd/server`.

## Request Flow

```text
Gin route
-> middleware
-> handler
-> request DTO validation
-> domain service input
-> domain service
-> repository interface
-> storage implementation
-> domain result
-> response DTO/envelope
```

## Handler Pattern

Handlers are transport code.

Allowed:

- read `gin.Context`
- parse query/path/body
- call service methods
- map errors to HTTP status and error code
- return JSON envelope

Not allowed:

- direct database access
- direct external API access
- non-trivial business rules
- trend ranking algorithms
- AI prompt construction

## Service Pattern

Services own business use cases.

Example:

```text
discovery.VenueService.List(ctx, search)
```

Services may:

- orchestrate repositories
- apply domain rules
- normalize domain input
- call other domain services through interfaces when needed

Services must not:

- import Gin
- import `internal/http`
- know HTTP status codes
- return response envelopes

## Repository Pattern

Repositories own persistence.

Example:

```text
discovery.VenueRepository
postgres.DiscoveryRepository
```

Rules:

- Define repository interfaces close to the domain that consumes them.
- Implement database repositories under `internal/storage/postgres`.
- Keep SQL explicit.
- Return domain models or storage models that are mapped before leaving the storage layer.
- Wrap database errors with operation context.
- Never return raw SQL text or internal database errors directly to clients.

## DTO Pattern

Use separate DTOs when the API contract diverges from the domain model.

MVP may return domain models directly only when the domain model is intentionally the public API shape.

When APIs grow, prefer:

```text
request DTO -> domain input
domain result -> response DTO
```

## Error Pattern

Domain errors should be stable and mappable.

HTTP layer maps errors to:

- HTTP status
- stable error code
- safe client message

Internal errors are logged, not returned to clients.

## Dependency Injection Pattern

Use constructor injection.

Good:

```text
service := discovery.NewVenueServiceWithRepository(repo)
router := apihttp.NewRouterWithLogger(service, logger)
```

Avoid global mutable dependencies.

## Transaction Pattern

When a use case needs multiple writes:

- start the transaction in storage or unit-of-work boundary
- pass `context.Context`
- commit only after all writes succeed
- rollback on error

Do not manage SQL transactions inside Gin handlers.

## Context Pattern

Every public service/repository method that can block must accept:

```go
ctx context.Context
```

Use request context from Gin:

```go
c.Request.Context()
```
