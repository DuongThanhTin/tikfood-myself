# ADR-0003: `{ data, error }` response envelope

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision predates this ADR)
- **Deciders:** Backend + frontend owners
- **Type:** retroactive (documents an existing decision)

## Context

The API and the Next.js frontend need a single, stable response shape so the client
can handle success and failure uniformly, and so contract changes are visible and
controlled. This is a public contract consumed by `apps/web/lib/api.ts`.

## Decision

All API responses use the envelope:

```json
{ "data": ..., "error": { "code": "...", "message": "...", "details?": ... } }
```

- Success sets `data`; error sets `error` with a predefined `code`.
- Handlers use the `respondWith*` helpers (`respondWithData`, `respondWithError`,
  `respondWithBadRequest`, `respondWithInternalServerError`).
- Error `code`s come from `internal/http/errors.go` (`invalid_request`, `not_found`,
  `internal_error`); domain errors map to HTTP status; raw SQL/error strings are never
  leaked.
- Do **not** invent a `{ "success": ... }` shape, and do **not** add `meta` or
  `trace_id` fields that do not exist yet.

## Consequences

- **Positive:** uniform client handling; stable contract; no breaking format changes.
- **Negative / cost:** every new endpoint must conform; envelope evolution requires a
  deliberate, coordinated change (backend + `apps/web/lib/api.ts`).
- **Follow-ups (approved, not yet done):** add `trace_id` to the envelope
  (`requestIDMiddleware` already exists); reserve a `meta` field for pagination when a
  list endpoint can realistically return large data; grow a documented error-code
  catalog as codes multiply.

## Alternatives considered

- **`{ success, data, error }`** — rejected: redundant with the presence of `error`.
- **Bare arrays/objects (no envelope)** — rejected: no uniform error channel.

## Sources / links

- [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) → *Response Envelope (real)* / *Error Handling*.
- [`docs/contracts/api.md`](../contracts/api.md), [`docs/standards/backend/request-response.md`](../standards/backend/request-response.md), [`docs/standards/backend/errors.md`](../standards/backend/errors.md).
- Consumer: [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) → *API Access*.
