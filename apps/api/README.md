# TikFood API

Go/Gin backend skeleton for TikFood discovery.

Architecture:

- `cmd/server`: thin process entrypoint.
- `internal/app`: dependency wiring and concrete adapter selection.
- `internal/http`: Gin routes, middleware, request parsing, response envelopes.
- `internal/discovery`: discovery domain models, service, repository interface, fallback repository, and pure search rules.
- `internal/storage/postgres`: PostgreSQL/PostGIS repository implementation.
- `internal/platform`: bounded reusable infrastructure helpers such as text normalization.

Implemented:

- `GET /health`
- `GET /api/v1/discovery/venues`
- `GET /api/v1/discovery/venues/:slug`
- `GET /api/v1/map/venues`
- Search, district, dish, tag, price, open-now, geo radius, social platform, and sort filters
- Venue detail response with dishes, opening hours, and social videos
- Structured request logging with Go `log/slog`
- PostgreSQL/PostGIS discovery persistence when `DATABASE_URL` is set
- In-memory fallback for discovery testing when no database is configured

This service does not implement delivery, cart, order, payment, booking, reservations, chat, monetization, or livestream logic.

Run locally:

```bash
go test ./...
go run ./cmd/server
```

Default URL: `http://localhost:18081`.
