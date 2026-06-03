# Backend Standards

This folder is the backend source of truth for TikFood.

Any AI assistant or engineer changing `apps/api` must read these files before coding:

- `docs/standards/backend/structure.md`
- `docs/standards/backend/architecture.md`
- `docs/standards/backend/dependencies.md`
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
- In-memory discovery data
- No database yet
- No auth yet
- No background workers yet

## Backend Principles

- Keep transport, domain, and storage separate.
- Prefer Go standard library until complexity justifies dependencies.
- Keep request validation at the HTTP boundary.
- Keep business rules inside domain services.
- Keep response contracts stable and typed.
- Log structured, non-sensitive context.
- Add tests with each behavior change.
- Design for PostgreSQL/PostGIS and workers later without adding premature infrastructure now.

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
