# Backend Architecture Standard

TikFood backend uses a layered architecture.

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

Domain logic should be testable without running an HTTP server.

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

Future package:

```text
internal/platform
```

Responsibilities:

- External API clients
- Clock/ID helpers
- Queue clients
- Model provider clients
- Object storage clients if ever needed

## Dependency Direction

Allowed:

```text
cmd/server -> internal/http -> internal/<domain> -> internal/storage
cmd/server -> internal/config
cmd/server -> log/slog
```

Avoid:

```text
internal/<domain> -> internal/http
internal/storage -> internal/http
internal/<domain> -> cmd/server
```

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
