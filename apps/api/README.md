# TikFood API

Go/Gin backend skeleton for TikFood discovery.

Implemented:

- `GET /health`
- `GET /api/v1/map/venues`
- District and dish filters
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
