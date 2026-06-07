# Backend Structure Standard

Backend code lives in:

```text
apps/api/
```

## Current Structure

```text
apps/api/
  cmd/server/
    main.go
  internal/
    app/
      app.go
      container.go
    config/
      config.go
    discovery/
      model.go
      service.go
      repository.go or fallback_repository.go
      search.go
    http/
      dependencies.go
      errors.go
      response.go
      router.go
      router_test.go
      venue_handler.go
      venue_request.go
    platform/
      textutil/
        normalize.go
    storage/
      postgres/
        discovery_repository.go
  Dockerfile
  go.mod
  README.md
```

## Target Structure As Backend Grows

```text
apps/api/
  cmd/server/
    main.go                 Process entrypoint

  internal/
    app/                    Composition root and dependency wiring
    config/                 Runtime config parsing
    http/                   Gin router, middleware, handlers, request parsing, response writing
    discovery/              Map/dish discovery domain
    trend/                  Trend scoring domain
    summary/                AI summary domain
    ingestion/              Social signal ingestion domain
    geo/                    Geo filtering/ranking helpers
    storage/                Database interfaces and PostgreSQL implementations
    platform/               Cross-cutting infrastructure adapters and bounded reusable helpers

  migrations/               Future SQL migrations
  testdata/                 Test fixtures
```

## Ownership Rules

`cmd/server`:

- Starts the process.
- Loads config and calls `internal/app`.
- Starts HTTP server.
- Must not contain business logic.
- Must not directly choose storage implementations.

`internal/app`:

- Owns dependency wiring and process-level composition.
- Creates logger, repositories, services, router, and HTTP server.
- Chooses concrete adapters such as Postgres or in-memory fallback.
- Must not contain business rules or request parsing.

`internal/http`:

- Owns routing, request parsing, response envelopes, status codes.
- May call domain services.
- Must not contain scoring, storage, ingestion, or summary business logic.
- Must use Gin as the only HTTP framework.
- `router.go` should only create the Gin router, attach middleware, health routes, and call route registration helpers.
- `dependencies.go` defines HTTP dependency structs and route registrar composition helpers.
- Domain-specific handlers should live in files such as `venue_handler.go`, `collection_handler.go`, or `profile_handler.go`.
- Request parsing should live in files such as `venue_request.go`.
- Shared API response helpers belong in `response.go`.
- Stable HTTP error codes/messages belong in `errors.go`.

`internal/<domain>`:

- Owns domain models, service methods, validation rules, and pure business logic.
- Must not know about HTTP request/response types unless explicitly needed for a small MVP.

`internal/storage`:

- Owns database access when PostgreSQL/PostGIS is added.
- Should expose interfaces or concrete repositories consumed by domain services.

`internal/platform`:

- Owns external integrations, clients, clocks, IDs, and bounded shared helpers.
- Shared helpers must live in named packages such as `textutil`, not generic `utils`.

## File Naming

Use simple names:

```text
service.go
model.go
repository.go
handler.go
router.go
middleware.go
response.go
errors.go
```

Avoid names like:

```text
utils.go
common.go
manager.go
processor.go
```

unless the file has a very clear bounded responsibility.

Preferred helper naming:

```text
internal/platform/textutil/normalize.go
internal/platform/idgen/generator.go
internal/platform/clock/clock.go
```

## Package Rules

- Keep packages small and cohesive.
- Do not create circular dependencies.
- Domain packages must not import `internal/http`.
- Storage packages must not import `internal/http`.
- HTTP can import domain packages.
- `internal/app` can import config, HTTP, domain, and storage packages.
- `cmd/server` should import only config/app plus standard library process packages.

## When Adding A New Backend Feature

1. Identify the domain.
2. Add or update domain model/service under `internal/<domain>`.
3. Add HTTP route/handler under `internal/http`.
4. Add storage only if persistence is required.
5. Add tests for domain and/or HTTP behavior.
6. Update API docs or standards if response shape changes.
