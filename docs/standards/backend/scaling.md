# Backend Scaling Standard

This standard defines how TikFood backend should grow without premature complexity.

## Scaling Philosophy

Scale in layers:

```text
simple local MVP
-> persistent MVP
-> worker-backed MVP
-> observable production system
-> high-scale data platform
```

Do not introduce later-stage infrastructure before product behavior needs it.

## Stage 1: Current Local MVP

Current:

- Go API
- `net/http`
- in-memory discovery data
- Next.js frontend fallback data
- no auth
- no database
- no workers

Good for:

- product shape
- API contract iteration
- frontend discovery UX
- AI automation workflow testing

## Stage 2: Persistent MVP

Add:

- PostgreSQL
- PostGIS
- SQL migrations
- `internal/storage`
- config package
- request logging middleware

Keep:

- explicit SQL
- thin handlers
- domain services
- simple deployment

Avoid:

- ORM by default
- distributed tracing before request volume exists
- event sourcing
- microservices

## Stage 3: Worker-Backed MVP

Add workers for:

- social ingestion
- trend scoring
- AI summaries

Recommended structure:

```text
cmd/api/
cmd/worker/
internal/ingestion/
internal/trend/
internal/summary/
internal/storage/
```

Workers should share domain/storage packages with the API where appropriate.

## Stage 4: Production Observability

Add:

- structured request logs
- metrics
- dashboards
- alerting
- audit logs for automation
- error tracking if needed

## Stage 5: Higher-Scale Data

Only after MVP pressure exists:

- queue system
- cache
- OpenSearch
- ClickHouse
- vector search
- feature store

These are not default MVP dependencies.

## Database Scaling Direction

MVP:

```text
PostgreSQL + PostGIS
```

Later:

```text
PostgreSQL read replicas
OpenSearch for discovery search
ClickHouse for analytics/trends
object storage for media-derived artifacts
```

## API Scaling Direction

Keep API contracts stable.

Add:

- pagination
- cursor-based lists
- bounded query limits
- request timeouts
- context-aware database calls
- graceful shutdown

## AI Summary Scaling Direction

AI summaries should be asynchronous when they become expensive.

Flow:

```text
social signal changes
-> summary job queued
-> worker generates summary
-> summary stored
-> API returns latest available summary
```

Do not generate expensive summaries synchronously inside map query endpoints.
