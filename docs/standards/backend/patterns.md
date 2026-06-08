# Backend Pattern Standard

This file defines the backend patterns that AI agents must follow.

## Primary Pattern

TikFood backend uses:

```text
Handler -> Service -> Repository
```

With explicit dependency wiring in `internal/app`; `cmd/server` remains a thin process entrypoint.

Architecture style:

```text
Layered architecture + light ports-and-adapters
```

The domain owns repository interfaces. Storage implements those interfaces. HTTP depends on domain services, not storage.

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

File rules:

- Keep `router.go` focused on router construction and high-level route registration.
- Put each domain handler in its own file, for example `venue_handler.go`.
- Put request parsing/validation helpers in `<domain>_request.go`.
- Do not add every handler to `router.go`; it becomes unreadable as models grow.
- Register domain routes through `RouteRegistrar` implementations.

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

HTTP error codes and common messages should be constants in `internal/http/errors.go`.

Common response writing should use helpers from `internal/http/response.go`:

```text
respondWithData
respondWithBadRequest
respondWithError
respondWithInternalServerError
```

Avoid repeating full response envelopes in every handler.

## Dependency Injection Pattern

Use constructor injection.

Good:

```text
application, err := app.New(config.Load())
container, err := app.NewContainer(config.Load())
router := apihttp.NewRouter(apihttp.RouterDependencies{
  Logger: container.Logger,
  RouteRegistrars: container.RouteRegistrars,
})
```

Avoid global mutable dependencies.

Composition root rules:

- Put concrete adapter selection in `internal/app`.
- Use `internal/app/container.go` to build repositories, services, and route registrars.
- Keep domain packages dependent on interfaces, not concrete storage.
- Keep `cmd/server/main.go` small enough to read in one screen.
- Do not initialize database connections inside HTTP handlers or domain services.
- Do not keep adding service-specific parameters to `NewRouter`; pass `RouteRegistrar` values instead.

When adding a new handler/service/repository:

1. Add the new domain package files.
2. Add the new HTTP handler implementing `RegisterRoutes`.
3. Wire the repository/service/registrar in the composition root.
4. Do not edit `router.go` unless router-wide behavior changes.

## Reusable Helper Pattern

Do not hide unrelated helpers inside service files.

Good:

```text
internal/platform/textutil.Normalize
internal/platform/textutil.Slugish
```

Avoid:

```text
internal/discovery/venue.go with normalizeVietnameseText mixed into service logic
```

Rules:

- Create a helper package only when at least one function has clear reuse potential.
- Give helper packages bounded names such as `textutil`, `clock`, or `idgen`.
- Avoid generic `utils`, `common`, or `helpers` packages.
- Keep helpers pure whenever possible.

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
