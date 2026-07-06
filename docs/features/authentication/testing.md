# Authentication — Testing

One-line purpose: the test strategy for the authentication feature — what to test at each layer, how it is verified, and how test coverage maps back to requirements.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

Define the testing approach so every requirement in [`specification.md`](specification.md)
is covered by a named test, every PR has a clear verification gate, and reviewers
can trust that auth behaves correctly and securely before merge. It lists test
scenarios by name (not source) and maps them to requirements; it does not contain
test code.

---

## Existing implementation

- **Backend** (`apps/api`): standard-library `testing` with `net/http/httptest`,
  table-driven style. Service-level tests use a mock/in-memory `VenueRepository`;
  the in-memory fallback repo also backs tests when `DATABASE_URL` is unset.
  Example: `internal/http/router_test.go` (`TestHealth`). No external test
  frameworks (no testify/ginkgo).
- **Frontend** (`apps/web`): **zero tests, no harness** (no Vitest/Jest, no React
  Testing Library).
- **CI**: none — `.github/workflows/` does not exist; verification is local only.

## Proposed changes

### Backend — unit tests (`internal/auth`, in-memory repos, no DB)

| Area | Test scenarios (names) |
|------|------------------------|
| Password | `TestHashPassword_VerifyRoundTrip`, `TestVerifyPassword_WrongPasswordFails`, `TestHashPassword_RejectsOver72Bytes` |
| JWT | `TestIssueAndParseAccessToken_RoundTrip`, `TestParseAccessToken_ExpiredFails`, `TestParseAccessToken_TamperedFails`, `TestParseAccessToken_RejectsNonHS256` (alg-confusion) |
| Refresh token | `TestNewRefreshTokenValue_HashDeterministic`, `TestNewRefreshTokenValue_HighEntropy` |
| Auth service | `TestRegister_HappyPath`, `_DuplicateEmail`, `_WeakPassword`, `_InvalidEmail`; `TestLogin_HappyPath`, `_WrongPassword`, `_UnknownUserSameError`; `TestRefresh_RotatesAndOldRevoked`, `_ExpiredRejected`, `_RevokedRejected`; `TestLogout_Idempotent` |
| Google | `TestLoginWithGoogle_NewUser`, `_ExistingByGoogleSub`, `_LinksVerifiedEmail`, `_RejectsUnverifiedEmailLink` |

### Backend — HTTP tests (`internal/http`, httptest)

| Area | Test scenarios |
|------|----------------|
| Register | `TestRegister_201SetsCookie`, `TestRegister_DuplicateReturns422` |
| Login | `TestLogin_200`, `TestLogin_401OnBadCreds` |
| Refresh | `TestRefresh_200RotatesCookie`, `TestRefresh_401WithoutCookie` |
| Logout | `TestLogout_ClearsCookie` |
| Envelope | `TestAuthResponsesUseEnvelope` (every response is `{data,error}`) |
| Middleware | `TestAuthMiddleware_ValidTokenAllows`, `_MissingTokenReturns401`, `_ExpiredReturns401` |
| CORS | `TestCORS_AllowsConfiguredOrigin`, `TestCORS_RejectsUnknownOrigin` |
| Protected | `TestMe_ReturnsCurrentUser` |
| Regression | `TestDiscoveryRoutesStayPublic` (existing endpoints unaffected) |
| Google | `TestGoogleLogin_RedirectsToGoogle`, `TestGoogleCallback_StateMismatchRejected` (fake authenticator injected) |

### Frontend — new Vitest + React Testing Library harness

Added as devDeps in milestone M7 (see [`tasks.md`](tasks.md)); command `npm run web:test`.

| Area | Test scenarios |
|------|----------------|
| `lib/validation` | `validateEmail` boundaries, `validatePassword` min-length boundary, required-field |
| `lib/auth` (mock fetch) | `login stores access token`, `getMe attaches Bearer`, `api retries once after 401 then refresh`, `public venue fetch works with no token`, `refresh failure clears token` |
| `AuthProvider` | `bootstraps authenticated when refresh succeeds`, `bootstraps anonymous when refresh fails`, `login updates user`, `logout clears user` |
| Login/Register pages | `shows error on invalid email`, `submit disabled while loading`, `renders API error text`, `successful submit redirects` |
| Google | `button links to backend login url`, `callback success sets session`, `callback failure shows error` |
| Discovery integration | `authCta links to login when anonymous`, `shows logout when authenticated`, `RequireAuth redirects anonymous` |

### Manual / E2E (documented in PR, not in CI)

- Live Google OAuth round trip in a dev environment.
- Apply migrations `007`/`008` to a scratch DB; inspect `users` / `refresh_tokens`.
- Responsive checks at 375 / 768 / 1280 px; dark and light `data-theme`.

### Coverage matrix (requirement → test)

| Requirement (see specification.md) | Primary tests |
|---|---|
| Registration | `TestRegister_*`, login/register page tests |
| Login | `TestLogin_*`, `lib/auth login`, page tests |
| Refresh / session | `TestRefresh_*`, `api retries once after 401`, `AuthProvider bootstraps` |
| Logout | `TestLogout_*`, `logout clears user` |
| Google SSO | `TestLoginWithGoogle_*`, `TestGoogleCallback_*`, Google FE tests |
| JWT auth / middleware | JWT unit tests, `TestAuthMiddleware_*`, `TestMe_*` |
| Validation | `lib/validation`, `TestRegister_WeakPassword/_InvalidEmail` |
| No breaking changes | `TestDiscoveryRoutesStayPublic`, envelope tests |
| Security controls | alg-confusion, `_UnknownUserSameError`, CORS, state-mismatch |

### Per-PR verification gates

- Backend: `go build ./...`, `go vet ./...`, `go test ./...` all green.
- Frontend: `npm run web:typecheck`, `npm run web:build`, `npm run web:test` all green.
- Recommend a `/security-review` pass before merging backend auth PRs (M4–M6) and
  the frontend token-handling PR (M7).

## Open questions

1. **CI gating** — add a minimal `.github/workflows` (separate, human-gated PR) so
   auth tests gate merges, or rely on local verification for this feature?
2. **Frontend harness scope** — Vitest + RTL only, or also add Playwright for the
   Google round trip and cookie/session E2E?
3. **Integration DB tier** — keep HTTP tests on in-memory repos only, or add a real
   Postgres tier (e.g. testcontainers) to catch SQL/migration regressions?
4. **Coverage threshold** — enforce a minimum for `internal/auth` and `lib/auth`?
5. **Google E2E flake** — is a mocked OAuth provider acceptable for automation, or
   must the live provider run on a schedule outside PR gates?

## Related repository documentation

- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — the completion gates this doc operationalizes.
- [`docs/standards/testing.md`](../../standards/testing.md) — repository testing standard.
- [`apps/api/CLAUDE.md`](../../../apps/api/CLAUDE.md) — backend test conventions (table-driven, httptest, in-memory repo).
- [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — frontend stack reality (no harness yet).
- Sibling docs: [`specification.md`](specification.md), [`security.md`](security.md), [`api.md`](api.md), [`tasks.md`](tasks.md), [`release.md`](release.md).
