# TikFood API

Go backend skeleton for TikFood discovery.

Implemented:

- `GET /health`
- `GET /api/v1/map/venues`
- District and dish filters
- In-memory MVP data for discovery testing

This service does not implement delivery, cart, order, payment, booking, reservations, chat, monetization, or livestream logic.

Run locally:

```bash
go test ./...
go run ./cmd/server
```

Default URL: `http://localhost:18081`.
