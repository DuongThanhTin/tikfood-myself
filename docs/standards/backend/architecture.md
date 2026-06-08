# Backend Architecture Standard

TikFood backend uses a layered architecture with a light ports-and-adapters style.

The practical rule is:

```text
Transport adapter -> domain service -> domain-owned port -> storage adapter
```

In current code this means:

- `internal/http`: Gin transport adapter.
- `internal/discovery`: domain models, service, search rules, and repository interface.
- `internal/storage/postgres`: PostgreSQL adapter implementing the domain repository interface.
- `internal/app`: composition root that builds concrete dependencies.

## Default Flow

```text
HTTP request
-> Gin router
-> middleware
-> handler
-> request validation
-> domain service
-> storage or external adapter
-> domain result
-> response DTO
-> JSON response
```

Process bootstrap flow:

```text
cmd/server/main.go
-> config.Load
-> internal/app.New
-> internal/app.NewContainer
-> concrete repository selection
-> domain service construction
-> HTTP route registrar construction
-> Gin router construction
-> http.Server
```

## Layer Responsibilities

### Transport Layer

Package:

```text
apps/api/internal/http
```

Responsibilities:

- Gin route registration
- Middleware composition
- Request parsing
- Query/path/body validation
- Auth boundary when added
- Calling domain services
- Mapping domain results to response DTOs
- Writing response envelopes

HTTP router construction must not know every service directly. It receives route registrars from `internal/app`.

Not allowed:

- Trend scoring algorithms
- Database queries
- Social ingestion logic
- AI summary generation
- Business rule branching beyond request orchestration

### Domain Layer

Package examples:

```text
internal/discovery
internal/trend
internal/summary
internal/ingestion
internal/geo
```

Responsibilities:

- Domain models
- Domain services
- Business rules
- Filtering and ranking rules
- Domain validation
- Interfaces needed by the domain

Domain logic should be testable without running an HTTP server, database, or Gin.

Domain package file shape should stay boring and predictable:

```text
model.go       Domain types and stable domain errors
service.go     Use cases and orchestration
search.go      Pure search/filter rules
repository.go  Domain-owned repository interfaces when useful
```

### Storage Layer

Future package:

```text
internal/storage
```

Responsibilities:

- PostgreSQL/PostGIS queries
- Transactions
- Mapping database rows to storage models or domain models
- Database error normalization

Storage should not return raw SQL errors directly to HTTP handlers.

### Platform Layer

Package:

```text
internal/platform
```

Responsibilities:

- External API clients
- Clock/ID helpers
- Queue clients
- Model provider clients
- Object storage clients if ever needed
- Bounded pure helpers that are not specific to one domain

## Dependency Direction

Allowed:

```text
cmd/server -> internal/app
internal/app -> internal/http
internal/app -> internal/<domain>
internal/app -> internal/storage
internal/http -> internal/<domain>
internal/storage -> internal/<domain>
internal/<domain> -> internal/platform/<helper>
```

Avoid:

```text
internal/<domain> -> internal/http
internal/storage -> internal/http
internal/<domain> -> cmd/server
```

## Dependency Wiring

Dependency wiring is centralized in:

```text
apps/api/internal/app/container.go
```

Rules:

- Build loggers, repositories, services, and route registrars in the container.
- Do not create repositories inside handlers.
- Do not create services inside `router.go`.
- Do not hide missing dependencies with fallback behavior in handlers.
- Constructor functions may validate required dependencies and fail fast during startup.

Adding a new domain should follow this shape:

```text
internal/<domain>/model.go
internal/<domain>/service.go
internal/<domain>/repository.go
internal/http/<domain>_handler.go
internal/http/<domain>_request.go
internal/storage/postgres/<domain>_repository.go
internal/app/container.go wires the new repository, service, and registrar
```

`router.go` should not need service-specific parameters when domains are added.

## Domain Boundaries

Discovery:

- Venue listing
- Dish search
- Map discovery result shape
- Hidden gem signals when added

Trend:

- Trend score calculation
- Ranking
- Signal weighting

Summary:

- AI-generated summaries
- Summary freshness
- Summary source attribution

Ingestion:

- Social content ingestion
- Creator/source references
- Raw signal normalization

Geo:

- Bounding boxes
- Distance calculations
- PostGIS query helpers later

Do not mix these domains just because an MVP endpoint returns combined data. Composition belongs in a service, not in random handlers.

## Configuration

Use environment variables for runtime config.

Rules:

- Config parsing belongs in `internal/config` when config grows.
- Never read `.env` inside app code.
- Docker/local tooling may provide env values.
- Secrets must come from environment or secret manager, never committed files.

## Scalability Path

Current:

```text
Gin
PostgreSQL/PostGIS when DATABASE_URL is set
in-memory fallback
single Go API service
```

Next:

```text
discovery search/filter endpoint
auth boundary
collection persistence
CORS and timeout middleware
```

Later:

```text
background workers
queue
trend scoring jobs
AI summary jobs
metrics
distributed tracing
```

Do not jump to later-stage infrastructure before a feature needs it.
