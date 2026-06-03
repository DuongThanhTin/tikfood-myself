# Backend Dependency Standard

Backend dependency policy is conservative by default.

## Current Backend Dependencies

The current backend uses:

- Go standard library
- `net/http`
- `encoding/json`
- `testing`

There are no third-party Go dependencies yet.

## Dependency Decision Rules

Before adding a Go module, answer:

1. What concrete complexity does this remove?
2. Can the Go standard library handle this cleanly?
3. Is this dependency actively maintained?
4. Does it introduce network, auth, codegen, reflection, or runtime magic?
5. Does it affect startup, build, deployment, or observability?
6. Will AI agents understand the pattern after reading local code?

If the answer is unclear, do not add the dependency.

## Preferred Backend Dependency Path

Routing:

- Keep `net/http` for MVP.
- Add `chi` only when middleware, route groups, or path params become painful.
- Add `Gin` only if the project explicitly chooses a Gin convention.
- Do not mix routers.

Database:

- Add PostgreSQL driver only when real persistence is needed.
- Prefer explicit SQL before ORM.
- Do not add an ORM by default.
- For PostGIS, keep geospatial behavior explicit and tested.

Logging:

- Start with standard `log` or a small internal structured logger.
- Add `slog` from standard library when structured logging is needed.
- Do not add heavy logging frameworks unless there is a clear operational reason.

Validation:

- Prefer small explicit validation functions.
- Do not add reflection-heavy validation unless request models become large and repetitive.

Testing:

- Use standard `testing` first.
- Add assertion libraries only if tests become noisy enough to justify them.

## go.mod Rules

- `apps/api/go.mod` owns backend Go dependencies.
- Do not edit `go.sum` manually.
- Run:

```bash
cd apps/api
go mod tidy
go test ./...
```

From root, prefer:

```bash
npm run api:test
```

## Dependency Change Checklist

Any PR adding a backend dependency must include:

- Why it is needed
- Why standard library is insufficient
- How it affects architecture
- Tests proving the new behavior
- Any security or operational risk

## Disallowed Dependency Motivation

Do not add dependencies because:

- "It might be useful later"
- "Most projects use it"
- "The AI generated it"
- "It makes one tiny helper shorter"
- "We may need it for delivery/payment/booking later"
