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
    discovery/
      venue.go
    http/
      router.go
      router_test.go
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
    config/                 Runtime config parsing
    logging/                Logger setup and request logging helpers
    http/                   Router, middleware, request parsing, response writing
    discovery/              Map/dish discovery domain
    trend/                  Trend scoring domain
    summary/                AI summary domain
    ingestion/              Social signal ingestion domain
    geo/                    Geo filtering/ranking helpers
    storage/                Database interfaces and PostgreSQL implementations
    platform/               Cross-cutting infrastructure adapters

  migrations/               Future SQL migrations
  testdata/                 Test fixtures
```

## Ownership Rules

`cmd/server`:

- Starts the process.
- Reads config.
- Creates logger.
- Wires dependencies.
- Starts HTTP server.
- Must not contain business logic.

`internal/http`:

- Owns routing, request parsing, response envelopes, status codes.
- May call domain services.
- Must not contain scoring, storage, ingestion, or summary business logic.

`internal/<domain>`:

- Owns domain models, service methods, validation rules, and pure business logic.
- Must not know about HTTP request/response types unless explicitly needed for a small MVP.

`internal/storage`:

- Owns database access when PostgreSQL/PostGIS is added.
- Should expose interfaces or concrete repositories consumed by domain services.

`internal/platform`:

- Owns external integrations, clients, clocks, IDs, and shared infrastructure adapters.

## File Naming

Use simple names:

```text
service.go
model.go
repository.go
handler.go
router.go
middleware.go
```

Avoid names like:

```text
utils.go
common.go
manager.go
processor.go
```

unless the file has a very clear bounded responsibility.

## Package Rules

- Keep packages small and cohesive.
- Do not create circular dependencies.
- Domain packages must not import `internal/http`.
- Storage packages must not import `internal/http`.
- HTTP can import domain packages.
- `cmd/server` can import all wiring-level packages.

## When Adding A New Backend Feature

1. Identify the domain.
2. Add or update domain model/service under `internal/<domain>`.
3. Add HTTP route/handler under `internal/http`.
4. Add storage only if persistence is required.
5. Add tests for domain and/or HTTP behavior.
6. Update API docs or standards if response shape changes.
