# Backend Testing Standard

Backend tests must protect behavior, not implementation trivia.

## Required Command

From repo root:

```bash
npm run api:test
```

Equivalent:

```bash
cd apps/api
GOCACHE=$(pwd)/../../.cache/go-build go test ./...
```

## Test Types

### HTTP Tests

Use `httptest`.

Test:

- status code
- response envelope
- response fields
- filtering behavior
- error status and error code when errors exist

### Domain Tests

Test domain services without HTTP.

Examples:

- discovery filtering
- trend score calculations
- summary freshness logic
- hidden gem ranking

### Storage Tests

When PostgreSQL is added:

- Keep unit tests for query builders or mappers.
- Add integration tests only with controlled test database setup.
- Do not require production credentials.

## What To Avoid

- Do not test private helper details unless logic is complex.
- Do not depend on test ordering.
- Do not call external network services in unit tests.
- Do not require `.env`.
- Do not snapshot huge JSON payloads.

## Test Naming

Use behavior names:

```go
func TestMapVenuesFiltersByDistrict(t *testing.T)
func TestVenueServiceReturnsEmptyListForUnknownDish(t *testing.T)
```

Avoid vague names:

```go
func TestHandler(t *testing.T)
func TestService(t *testing.T)
```

## Fixtures

Use `testdata/` for larger fixtures.

Small in-memory fixtures inside tests are fine.

## Coverage Expectation

For each backend feature request:

- Add or update at least one relevant test.
- If no test is added, explain why.
- Always run `npm run api:test`.

## AI Agent Requirement

AI tools must report:

- tests added
- tests updated
- commands run
- skipped tests with reason
- known risks
