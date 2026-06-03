# Backend Standards

This folder is the backend source of truth for TikFood.

Any AI assistant or engineer changing `apps/api` must read these files before coding:

- `docs/standards/backend/framework.md`
- `docs/standards/backend/structure.md`
- `docs/standards/backend/architecture.md`
- `docs/standards/backend/patterns.md`
- `docs/standards/backend/dependencies.md`
- `docs/standards/backend/database.md`
- `docs/standards/backend/search-filtering.md`
- `docs/standards/backend/request-response.md`
- `docs/standards/backend/errors.md`
- `docs/standards/backend/logging.md`
- `docs/standards/backend/testing.md`
- `docs/standards/backend/scaling.md`

## Current Backend Position

TikFood backend is currently a Go MVP API in `apps/api`.

Current capabilities:

- `GET /health`
- `GET /api/v1/map/venues`
- Gin HTTP routing
- Structured request logging with Go `log/slog`
- PostgreSQL/PostGIS discovery persistence when `DATABASE_URL` is set
- In-memory discovery fallback when no database is configured
- No auth yet
- No background workers yet

## Backend Principles

- Keep transport, domain, and storage separate.
- Use Gin as the backend HTTP framework.
- Use Go `log/slog` for structured logs.
- Keep request validation at the HTTP boundary.
- Keep business rules inside domain services.
- Keep response contracts stable and typed.
- Log structured, non-sensitive context.
- Add tests with each behavior change.
- Design for PostgreSQL/PostGIS and workers later without adding premature infrastructure now.
- Use `about` on venues and dishes for rich detail pages; use short descriptions for cards/list results.

## MVP Product Guardrails

Do not implement backend logic for:

- Delivery
- Cart
- Orders
- Checkout
- Payment
- Booking or reservations
- In-app chat
- Social follow graph
- Creator monetization
- Livestream

If a request touches these areas, stop and ask for human approval.
