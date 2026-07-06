# Authentication System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan milestone-by-milestone. Steps use checkbox (`- [ ]`) syntax for tracking. Each milestone is one Pull Request.

**Goal:** Add email/password + Google OAuth authentication (register, login, refresh, logout, session management) to the TikFood Go API and Next.js web app, reusing the existing design system without introducing a new design language.

**Architecture:** Stateless short-lived **access JWT** (~15 min) sent as a `Bearer` header, paired with a **rotating refresh token** persisted in Postgres and delivered via an httpOnly `Secure` cookie (this is the "session" — revocable on logout). Backend follows the existing `Handler → Service → Repository → Postgres` layering (ADR-0002) and the `{data,error}` envelope (ADR-0003). Frontend keeps the access token in memory inside an `AuthProvider`, refreshes silently via the cookie, and routes all calls through `lib/api.ts` / `lib/auth.ts`.

**Tech Stack:** Go 1.23 · Gin · pgx/`database/sql` · `github.com/golang-jwt/jwt/v5` · `golang.org/x/crypto/bcrypt` · `golang.org/x/oauth2` (+ `google` endpoint) · Next.js App Router + TypeScript · plain CSS (`globals.css` tokens) · Vitest + React Testing Library (new devDeps).

## Global Constraints

- **Human-approval gates already cleared** for auth, migrations, and Google (external network) — recorded via the ADR in Milestone 0. Any *new* gated area beyond these still requires approval.
- **Branch discipline:** all work on branches under `ai/`; never push to `main`/`master`; no force-push; no auto-merge. One PR per milestone.
- **Response envelope is immutable:** every API response is `{ "data": ..., "error": { "code", "message", "details?" } }`. No breaking changes to existing discovery endpoints; they stay public.
- **Layering (ADR-0002):** Handler validates + parses; Service holds business logic and never imports `gin`; Repository does persistence only. No forbidden edges.
- **Secrets:** never read `.env*`; never log passwords, tokens, hashes, or secrets. All secrets come from `internal/config` via env.
- **Naming:** verb-first (`CreateUser`, `FindByEmail`, `IssueTokens`); never `Do`/`Handle`/`Process`/`Component1`.
- **Frontend:** plain CSS only (reuse `globals.css` tokens/classes); all fetches through `lib/api.ts`/`lib/auth.ts`; frontend types mirror Go JSON tags in **snake_case**; UI copy in **Vietnamese**, identifiers in English; every icon-only button has `aria-label`; inputs wrapped in `<label>`.
- **Design system:** reuse existing tokens/components (`--primary`, `--danger`, `--surface-low`, `.primaryButton`, `.field`, `.errorText`, glass tokens). Do **not** introduce a new design language. Update `docs/design/` only when a genuinely new reusable pattern is added (Milestone 12).
- **Migrations** live in `apps/api/migrations/`, are **additive only**, applied by the existing PostGIS Docker entrypoint (no runner lib). Never alter existing discovery tables.
- **TDD:** write the failing test first, watch it fail, implement minimally, watch it pass, commit. Backend uses stdlib `testing` + `httptest` + in-memory repo (mirrors discovery). Frontend uses Vitest + RTL.
- **Verification floors per PR:** backend — `go build ./...`, `go vet ./...`, `go test ./...` all green; frontend — `npm run web:typecheck`, `npm run web:build`, `npm test` all green.

---

## Design Decisions (locked for this plan)

| # | Decision | Rationale |
|---|----------|-----------|
| D1 | Access JWT (~15 min) + rotating refresh token in DB | Reconciles "JWT auth" + "session management"; enables logout/revocation. |
| D2 | Refresh token in httpOnly `Secure` `SameSite=Lax` cookie; access token in memory | Mitigates XSS token theft; cookie enables silent refresh. |
| D3 | Dedicated `/login`, `/register` **pages** (+ `/forgot-password` placeholder) | Feature spec explicitly names Login/Register pages. |
| D4 | Controlled-state validation helper, no form library | Respects current frontend standard ("no RHF/Zod yet"); avoids new runtime dep. |
| D5 | Google OAuth **Authorization Code** flow (server-side callback) | Keeps client secret server-side; standard, secure. |
| D6 | New `internal/auth/` domain + in-memory fallback repo | Mirrors existing `internal/discovery/` pattern; enables tests without a DB. |
| D7 | Auth endpoints under `/api/v1/auth/*`; discovery endpoints unchanged & public | Additive, non-breaking. |
| D8 | CORS middleware added (credentialed) with an explicit allowed-origin list from config | Browser cookie auth across web/api origins requires it. |

---

## File Structure (created / modified across the whole feature)

**Backend — `apps/api/`**
- `internal/auth/` *(new domain)* — `model.go`, `service.go`, `repository.go` (interface), `memory_repository.go`, `password.go`, `token.go`, `oauth_google.go`
- `internal/storage/postgres/` — `user_repository.go`, `refresh_token_repository.go` *(new)*
- `internal/http/` — `auth_handler.go` *(new)*; modify `router.go`, `middleware.go`, `errors.go`
- `internal/config/config.go` — modify (new keys)
- `internal/app/container.go` — modify (wire auth deps)
- `migrations/007_auth_users.sql`, `migrations/008_auth_refresh_tokens.sql` *(new, additive)*
- `go.mod` / `go.sum` — modify (new deps)

**Frontend — `apps/web/`**
- `lib/auth.ts` *(new)*; modify `lib/api.ts`
- `lib/validation.ts` *(new)* — pure validators
- `components/auth/AuthProvider.tsx`, `components/auth/AuthForm.tsx`, `components/auth/FormField.tsx`, `components/auth/GoogleSignInButton.tsx` *(new)*
- `app/login/page.tsx`, `app/register/page.tsx`, `app/forgot-password/page.tsx` *(new)*
- `app/auth/google/callback/page.tsx` *(new)*
- modify `app/layout.tsx` (wrap `AuthProvider`), `components/DiscoveryExperience.tsx` (wire `.authCta`), `app/globals.css` (auth-form/page classes reusing tokens)
- `package.json` — modify (Vitest + RTL devDeps, `web:test` script), `vitest.config.ts` *(new)*

**Docs**
- `docs/adr/0007-authentication-approach.md` *(new)*
- `docs/services/api.md`, `docs/services/web.md` — modify (auth contract)
- `docs/design/components/modal.md` and/or auth-form playbook — modify **only if** a new reusable pattern is introduced
- `IMPROVEMENTS.md` — modify (resolve inert `.authCta` item)

---

## Milestones

Each milestone = one PR. Milestones are ordered so every PR builds and tests green on its own. Backend M0–M6 and frontend M7–M11 can partly overlap once M4's contract is merged.

---

### Milestone 0 — ADR + dependencies + config scaffolding

**Purpose:** Record the authentication decision (ADR-0007) and land the non-functional groundwork: new Go deps, new config keys (unused yet), and CORS config plumbing — so later milestones don't repeat setup.

**Dependencies:** None (first milestone).

**Affected services:** `apps/api` (config, deps), `docs/`.

**Files:**
- Create `docs/adr/0007-authentication-approach.md` — captures D1–D8 (method, token storage/transport, layering placement, schema outline, OAuth flow).
- Modify `apps/api/go.mod` / `go.sum` — add `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto`, `golang.org/x/oauth2`.
- Modify `apps/api/internal/config/config.go` — add fields + env loads with safe defaults: `JWTSecret` (`JWT_SECRET`), `AccessTokenTTL` (`ACCESS_TOKEN_TTL`, default `15m`), `RefreshTokenTTL` (`REFRESH_TOKEN_TTL`, default `720h`), `GoogleClientID` (`GOOGLE_CLIENT_ID`), `GoogleClientSecret` (`GOOGLE_CLIENT_SECRET`), `GoogleRedirectURL` (`GOOGLE_REDIRECT_URL`), `WebOrigin`/`AllowedOrigins` (`WEB_ORIGIN`, default `http://localhost:3000`).
- Modify `.env.example` (template only — never a real `.env`) with the new keys documented and blank/placeholder values.

**Verification:**
- `cd apps/api && go build ./...` succeeds; `go vet ./...` clean.
- Add `TestConfigDefaultsAuth` in `internal/config/config_test.go`: with no env set, TTLs equal defaults and `WebOrigin` defaults. Run `go test ./internal/config/... -v` → PASS.
- ADR renders and links from `docs/adr/README.md` index.

**Risks:** Weak/placeholder `JWT_SECRET` accidentally shipping — mitigate: service refuses to start (Milestone 2) if secret is empty in non-dev. `go.sum` churn — keep the deps minimal.

**Estimated complexity:** Low.

**Suggested PR boundary:** "chore(auth): ADR-0007 + auth deps & config scaffolding". No behavior change; safe to merge early.

---

### Milestone 1 — DB schema: users & refresh_tokens (migrations + repositories)

**Purpose:** Persist users and refresh tokens with correct constraints; provide repository interfaces + Postgres + in-memory implementations (data layer only, no business logic).

**Dependencies:** M0 (config), gated migration approval (granted in Discovery).

**Affected services:** `apps/api` (migrations, `internal/auth`, `internal/storage/postgres`), PostgreSQL.

**Files:**
- Create `apps/api/migrations/007_auth_users.sql` — `users(id uuid pk, email citext unique not null, password_hash text null, display_name text, google_sub text unique null, created_at, updated_at)`; index on `email`, on `google_sub`.
- Create `apps/api/migrations/008_auth_refresh_tokens.sql` — `refresh_tokens(id uuid pk, user_id uuid fk→users on delete cascade, token_hash text unique not null, expires_at timestamptz not null, revoked_at timestamptz null, created_at, user_agent text, ip inet null)`; index on `user_id`, on `token_hash`.
- Create `apps/api/internal/auth/model.go` — `User`, `RefreshToken` structs with snake_case JSON tags (password_hash **never** serialized: `json:"-"`).
- Create `apps/api/internal/auth/repository.go` — interfaces: `UserRepository` (`CreateUser`, `FindByEmail`, `FindByID`, `FindByGoogleSub`, `LinkGoogleSub`), `RefreshTokenRepository` (`StoreRefreshToken`, `FindRefreshTokenByHash`, `RevokeRefreshToken`, `RevokeAllForUser`).
- Create `apps/api/internal/auth/memory_repository.go` — in-memory impls for tests.
- Create `apps/api/internal/storage/postgres/user_repository.go`, `refresh_token_repository.go` — pgx/`database/sql` impls, parameterized queries only.

**Verification:**
- `docker compose config` valid (no compose change expected).
- `go build ./...`, `go vet ./...` clean.
- Table-driven tests against the **in-memory** repo: `TestMemoryUserRepository_CreateAndFindByEmail`, `_DuplicateEmailRejected`, `_FindByGoogleSub`; `TestMemoryRefreshTokenRepository_StoreFindRevoke`. Run `go test ./internal/auth/... ./internal/storage/... -v` → PASS.
- Manual (documented in PR, not automated): apply migrations to a scratch DB via the Docker entrypoint; `\d users` shows constraints.

**Risks:** `citext` extension may be absent — migration must `CREATE EXTENSION IF NOT EXISTS citext;` (needs superuser at container init — verify against `db/Dockerfile`). Cascade deletes on refresh tokens — intended. Email uniqueness case-sensitivity — solved by `citext`.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(auth): user & refresh-token schema + repositories". Ships data layer with tests; no HTTP surface yet.

---

### Milestone 2 — Password hashing + JWT token service (pure core)

**Purpose:** Provide pure, unit-testable primitives: bcrypt hash/verify and access/refresh token issue/parse — the security core, isolated from HTTP and DB.

**Dependencies:** M0 (config: secret + TTLs), M1 (model types).

**Affected services:** `apps/api` (`internal/auth`).

**Files:**
- Create `apps/api/internal/auth/password.go` — `HashPassword(plain string) (string, error)` (bcrypt, cost 12), `VerifyPassword(hash, plain string) error`.
- Create `apps/api/internal/auth/token.go` — `TokenIssuer` struct (holds secret, TTLs); `IssueAccessToken(user User) (string, error)` (HS256 JWT, claims: sub, email, exp, iat, iss=`tikfood`); `ParseAccessToken(raw string) (Claims, error)`; `NewRefreshTokenValue() (raw string, hash string, err error)` (random 32 bytes → hash stored, raw returned to client); `HashRefreshToken(raw string) string`.

**Verification:**
- Tests: `TestHashPassword_VerifyRoundTrip`, `TestVerifyPassword_WrongPasswordFails`, `TestIssueAndParseAccessToken_RoundTrip`, `TestParseAccessToken_ExpiredFails` (issue with negative TTL), `TestParseAccessToken_TamperedFails`, `TestNewRefreshTokenValue_HashDeterministic`. Run `go test ./internal/auth/... -run 'Password|Token' -v` → PASS.
- `go vet ./...` clean; ensure secret-empty construction returns an error (guard).

**Risks:** Wrong signing alg acceptance (alg-confusion) — `ParseAccessToken` must enforce HS256 explicitly. bcrypt 72-byte truncation — document password max length (validated in Milestone 3). Refresh token entropy — use `crypto/rand`.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(auth): password hashing + JWT token service". Pure functions, fully unit-tested; no wiring.

---

### Milestone 3 — Auth service (register / login / refresh / logout business logic)

**Purpose:** Orchestrate the primitives + repositories into the core use-cases, with domain validation and correct error semantics — no HTTP, no `gin`.

**Dependencies:** M1 (repos), M2 (password/token).

**Affected services:** `apps/api` (`internal/auth`).

**Files:**
- Create `apps/api/internal/auth/service.go` — `AuthService` with:
  - `Register(ctx, email, password, displayName) (User, TokenPair, error)` — validate email/password strength, reject duplicate, hash, create, issue tokens, store refresh hash.
  - `Login(ctx, email, password) (User, TokenPair, error)` — find, verify, issue+store; uniform "invalid credentials" error (no user-enumeration).
  - `Refresh(ctx, rawRefresh) (TokenPair, error)` — look up by hash, check not expired/revoked, **rotate** (revoke old, issue new).
  - `Logout(ctx, rawRefresh) error` — revoke by hash (idempotent).
  - `Me(ctx, userID) (User, error)`.
  - `TokenPair{ AccessToken string; RefreshTokenRaw string; AccessExpiresAt time.Time }`.
  - Sentinel errors: `ErrInvalidCredentials`, `ErrEmailTaken`, `ErrWeakPassword`, `ErrInvalidEmail`, `ErrRefreshInvalid`.

**Verification:**
- Tests with in-memory repos: `TestRegister_HappyPath`, `TestRegister_DuplicateEmail`, `TestRegister_WeakPassword`, `TestRegister_InvalidEmail`, `TestLogin_HappyPath`, `TestLogin_WrongPassword`, `TestLogin_UnknownUserSameError` (identical error to wrong-password), `TestRefresh_RotatesAndOldRevoked`, `TestRefresh_ExpiredRejected`, `TestRefresh_RevokedRejected`, `TestLogout_Idempotent`. Run `go test ./internal/auth/... -v` → PASS.

**Risks:** User enumeration via differing errors/timing — return identical `ErrInvalidCredentials` for unknown-user and wrong-password. Refresh replay after rotation — old token must be revoked atomically before new issued. Password policy too weak/strong — codify min length 8, cap ≤72 bytes.

**Estimated complexity:** Medium-High.

**Suggested PR boundary:** "feat(auth): auth service (register/login/refresh/logout)". Full business logic under test; still no HTTP.

---

### Milestone 4 — HTTP handlers + routes (email/password endpoints)

**Purpose:** Expose register/login/refresh/logout over HTTP with request validation, the `{data,error}` envelope, correct status codes, and refresh-token cookie handling. This milestone freezes the API contract the frontend consumes.

**Dependencies:** M3 (service), M0 (config), gated auth approval.

**Affected services:** `apps/api` (`internal/http`, `internal/app`).

**Files:**
- Create `apps/api/internal/http/auth_handler.go` — `AuthHandler` implementing `RouteRegistrar`; routes on `v1.Group("/auth")`:
  - `POST /register` → 201, body `{email,password,display_name}`; sets refresh cookie; returns `{data:{user, access_token, access_expires_at}}`.
  - `POST /login` → 200; same shape.
  - `POST /refresh` → 200; reads refresh cookie, rotates, resets cookie.
  - `POST /logout` → 200; revokes, clears cookie.
  - Request structs + validation in the handler (mirror `venue_request.go`), mapping to `422 domain_rejected` / `400 invalid_request`.
- Modify `apps/api/internal/http/errors.go` — add codes `unauthorized` (401), `forbidden` (403), and map sentinel service errors → codes/status. Never leak tokens/SQL.
- Modify `apps/api/internal/http/router.go` — register `AuthHandler`.
- Modify `apps/api/internal/app/container.go` — construct repos → service → handler.
- Add cookie helpers: `setRefreshCookie`, `clearRefreshCookie` (httpOnly, Secure, SameSite=Lax, Path=`/api/v1/auth`).

**Verification:**
- `httptest` table-driven tests: `TestRegister_201SetsCookie`, `TestRegister_DuplicateReturns422`, `TestLogin_200`, `TestLogin_401OnBadCreds`, `TestRefresh_200RotatesCookie`, `TestRefresh_401WithoutCookie`, `TestLogout_ClearsCookie`, `TestResponsesUseEnvelope`. Run `go test ./internal/http/... -v` → PASS.
- `go build ./...`, `go vet ./...` clean. Confirm existing discovery route tests still pass (`TestHealth`, venue tests).

**Risks:** Cookie flags wrong for local (Secure over http) — document dev toggle via config; default Secure=true, allow `COOKIE_SECURE=false` for local http. Envelope drift — assert in tests. Handler leaking service internals — map errors explicitly.

**Estimated complexity:** Medium-High.

**Suggested PR boundary:** "feat(auth): auth HTTP endpoints (register/login/refresh/logout)". **Contract-freezing PR** — frontend milestones depend on it.

---

### Milestone 5 — Auth middleware + CORS + protected `/me`

**Purpose:** Validate access tokens on protected routes, inject the user into context, add credentialed CORS, and ship the first protected endpoint (`GET /auth/me`) proving the chain end-to-end.

**Dependencies:** M4 (handlers/contract), M2 (token parse), M0 (allowed origins).

**Affected services:** `apps/api` (`internal/http`, `internal/app`).

**Files:**
- Modify `apps/api/internal/http/middleware.go` — add `authMiddleware(issuer)` (reads `Authorization: Bearer`, parses, sets `userID` in `gin.Context`, else 401 envelope) and `corsMiddleware(allowedOrigins)` (credentialed; echoes allowed origin; handles `OPTIONS`).
- Modify `apps/api/internal/http/router.go` — add `corsMiddleware` to the global chain (before routes); apply `authMiddleware` only to a protected group.
- Modify `apps/api/internal/http/auth_handler.go` — add `GET /auth/me` under the protected group → `{data:{user}}`.

**Verification:**
- Tests: `TestAuthMiddleware_ValidTokenAllows`, `TestAuthMiddleware_MissingTokenReturns401`, `TestAuthMiddleware_ExpiredReturns401`, `TestMe_ReturnsCurrentUser`, `TestCORS_AllowsConfiguredOrigin`, `TestCORS_RejectsUnknownOrigin`, `TestDiscoveryRoutesStayPublic`. Run `go test ./internal/http/... -v` → PASS.
- `go build/vet/test ./...` all green.

**Risks:** CORS misconfig blocks the browser or opens `*` with credentials (forbidden combo) — never `*` with credentials; use explicit list. Public discovery routes accidentally gated — assert they stay open. Bearer parsing edge cases — trim/validate scheme.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(auth): auth middleware + CORS + /me". Completes the backend session loop for email/password.

---

### Milestone 6 — Google OAuth SSO (backend)

**Purpose:** Add server-side Authorization-Code Google login: start URL + callback that provisions/links a user and issues the same token pair.

**Dependencies:** M3 (service), M4/M5 (issuing + cookies), M0 (Google config), gated external-call approval.

**Affected services:** `apps/api` (`internal/auth`, `internal/http`), Google (external).

**Files:**
- Create `apps/api/internal/auth/oauth_google.go` — `GoogleAuthenticator`: `AuthCodeURL(state) string`; `ExchangeAndFetchProfile(ctx, code) (GoogleProfile{Sub,Email,Name}, error)` using `golang.org/x/oauth2` + `google` endpoint.
- Modify `apps/api/internal/auth/service.go` — `LoginWithGoogle(ctx, profile) (User, TokenPair, error)`: find by `google_sub` → else find by email + link → else create user (nil password_hash).
- Modify `apps/api/internal/http/auth_handler.go`:
  - `GET /auth/google/login` → sets signed/opaque `state` cookie, 302 to `AuthCodeURL`.
  - `GET /auth/google/callback` → verifies `state`, exchanges code, `LoginWithGoogle`, sets refresh cookie, 302 to `WEB_ORIGIN/auth/google/callback#access=...` (or a short-lived one-time code the SPA exchanges).

**Verification:**
- Tests with a **fake** authenticator (interface-injected; no live Google call): `TestLoginWithGoogle_NewUser`, `_ExistingByGoogleSub`, `_LinksExistingEmail`, `TestGoogleCallback_StateMismatchRejected`, `TestGoogleLogin_RedirectsToGoogle`. Run `go test ./internal/auth/... ./internal/http/... -v` → PASS.
- `go build/vet ./...` clean. Live Google flow verified manually in a dev environment (documented in PR; not in CI).

**Risks:** CSRF on callback — require `state` cookie match. Account-takeover via unverified email linking — only link when Google reports `email_verified`. Secret exposure — client secret server-side only, never logged. External dependency flakiness — isolate behind interface for tests.

**Estimated complexity:** High.

**Suggested PR boundary:** "feat(auth): Google OAuth SSO (backend)". Backend feature-complete after this.

---

### Milestone 7 — Frontend auth API client + token handling

**Purpose:** Add the typed client for auth calls and teach `lib/api.ts` to attach the in-memory access token and perform one silent refresh on 401 — without breaking existing public calls.

**Dependencies:** M4 (contract), M5 (/me), M0 (CORS/cookies). Frontend test harness set up here.

**Affected services:** `apps/web` (`lib/`).

**Files:**
- Modify `apps/web/package.json` — add devDeps `vitest`, `@testing-library/react`, `@testing-library/jest-dom`, `jsdom`; add script `"web:test": "vitest run"`. Create `apps/web/vitest.config.ts`.
- Create `apps/web/lib/auth.ts` — types (`AuthUser`, `TokenResponse` — snake_case matching Go), functions `register`, `login`, `refresh`, `logout`, `getMe`, `googleLoginUrl`; in-memory `accessToken` accessor (`setAccessToken`/`getAccessToken`). All via `credentials: 'include'`.
- Modify `apps/web/lib/api.ts` — inject `Authorization: Bearer` when a token is set; on 401 from a protected call, attempt one `refresh()` then retry; preserve `{data,error}` handling; public discovery calls unchanged when no token.

**Verification:**
- Vitest tests (mock `fetch`): `login sends credentials and stores access token`, `getMe attaches bearer header`, `api retries once after 401 then refresh`, `public venue fetch works with no token`, `refresh failure clears token`. Run `cd apps/web && npm run web:test` → PASS.
- `npm run web:typecheck`, `npm run web:build` → PASS.

**Risks:** Infinite refresh loop on repeated 401 — guard with a single-retry flag. Token in memory lost on reload — acceptable; `AuthProvider` rehydrates via refresh cookie (M8). CORS/cookie mismatch — validated against M5.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(web): auth API client + token refresh + test harness".

---

### Milestone 8 — AuthProvider (session context + bootstrap)

**Purpose:** Provide a single React context for session state, silent bootstrap on load (via refresh cookie), and `login/register/logout` actions the UI consumes.

**Dependencies:** M7 (client).

**Affected services:** `apps/web` (`components/auth`, `app/layout.tsx`).

**Files:**
- Create `apps/web/components/auth/AuthProvider.tsx` — context `{ user, status: 'loading'|'authenticated'|'anonymous', login, register, logout, loginWithGoogleUrl }`; on mount calls `refresh()`→`getMe()` to rehydrate; exposes `useAuth()` hook.
- Modify `apps/web/app/layout.tsx` — wrap children in `<AuthProvider>`.

**Verification:**
- Vitest + RTL: `AuthProvider bootstraps to authenticated when refresh succeeds`, `bootstraps to anonymous when refresh fails`, `login updates user`, `logout clears user`. Run `npm run web:test` → PASS.
- `npm run web:typecheck`, `npm run web:build` → PASS.

**Risks:** Hydration/SSR mismatch (App Router) — `AuthProvider` is a client component; render a stable `loading` state on first paint. Bootstrap request storm — run once on mount.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(web): AuthProvider session context".

---

### Milestone 9 — Login & Register pages (forms, validation, states, responsive)

**Purpose:** Ship the `/login` and `/register` pages reusing existing design tokens/components, with client-side validation, loading and error states, and responsive layout.

**Dependencies:** M8 (useAuth), M7 (client).

**Affected services:** `apps/web` (`app/`, `components/auth`, `globals.css`).

**Files:**
- Create `apps/web/lib/validation.ts` — pure `validateEmail`, `validatePassword` (min 8), `validateRequired`; returns Vietnamese messages.
- Create `apps/web/components/auth/FormField.tsx` — reuses `.field` + native input, renders `.errorText` when invalid, `<label>` wrapping, `aria-describedby`.
- Create `apps/web/components/auth/AuthForm.tsx` — shared controlled-state form shell (email + password [+ display_name on register]); submit via `useAuth`; loading via `disabled`+label swap; maps API `error.message`; success → redirect to `/`.
- Create `apps/web/app/login/page.tsx`, `apps/web/app/register/page.tsx` — compose `AuthForm`, cross-link, include `GoogleSignInButton` slot (wired in M10).
- Modify `apps/web/app/globals.css` — add `.authPage`, `.authCard` classes built **only** from existing tokens (glass surface, `--surface-low` inputs, 24px→16px gutter, `.filterGrid` responsive pattern). No new colors/fonts.

**Verification:**
- Vitest tests in `lib/validation.ts` (`validateEmail`, `validatePassword` boundaries) and RTL: `login page shows error on invalid email`, `submit disabled while loading`, `renders API error text`, `successful login redirects`. Run `npm run web:test` → PASS.
- `npm run web:typecheck`, `npm run web:build` → PASS. Manual responsive check at 375/768/1280 documented in PR.

**Risks:** Introducing a competing design language — mitigate by reusing only documented tokens/classes; design self-review before PR. A11y regressions — labels + `aria-describedby` + focus management on error. Copy not Vietnamese — review strings.

**Estimated complexity:** Medium-High.

**Suggested PR boundary:** "feat(web): login & register pages".

---

### Milestone 10 — Google Sign-In (frontend) + callback route

**Purpose:** Wire the Google button to the backend start URL and handle the callback that lands the session, reusing existing button language (no foreign branded widget breaking the design language, within Google brand rules).

**Dependencies:** M6 (backend Google), M8 (context), M9 (pages).

**Affected services:** `apps/web` (`components/auth`, `app/auth/google/callback`).

**Files:**
- Create `apps/web/components/auth/GoogleSignInButton.tsx` — reuses a glass/secondary button variant + `aria-label`; navigates to `auth.googleLoginUrl()`.
- Create `apps/web/app/auth/google/callback/page.tsx` — reads the returned one-time access hand-off (or triggers `refresh()`), sets session via `useAuth`, redirects to `/` or shows `.errorText` on failure.
- Modify `apps/web/app/login/page.tsx` / `register/page.tsx` — mount the button.

**Verification:**
- RTL: `google button links to backend login url`, `callback success sets session and redirects`, `callback failure shows error`. Run `npm run web:test` → PASS.
- `npm run web:typecheck`, `npm run web:build` → PASS. End-to-end Google login verified manually in dev (documented).

**Risks:** Hand-off token in URL fragment leaking to history — use short-lived one-time code and clear the fragment after read. Design drift on the Google button — keep within existing variants + Google brand guidance.

**Estimated complexity:** Medium.

**Suggested PR boundary:** "feat(web): Google Sign-In + callback".

---

### Milestone 11 — Forgot-password placeholder + wire `.authCta` + optional guard

**Purpose:** Add the future forgot-password placeholder page, wire the existing inert `.authCta` button to `/login` (resolving accessibility gap A6), and add a minimal route guard/redirect for any authenticated-only surface.

**Dependencies:** M8 (context), M9 (pages).

**Affected services:** `apps/web` (`app/`, `components/DiscoveryExperience.tsx`).

**Files:**
- Create `apps/web/app/forgot-password/page.tsx` — static placeholder, clearly marked "future" (TODO), Vietnamese copy, links back to `/login`.
- Modify `apps/web/components/DiscoveryExperience.tsx` — make `.authCta` a real link/button to `/login` (or show user + logout when authenticated via `useAuth`); add `aria-label`.
- Modify `apps/web/lib/auth.ts` or add `components/auth/RequireAuth.tsx` — small client guard that redirects `anonymous` users to `/login` (only used if/when a protected page exists; ship the util + a test).
- Modify `IMPROVEMENTS.md` — mark the inert `.authCta` item resolved.

**Verification:**
- RTL: `authCta links to login when anonymous`, `shows logout when authenticated`, `RequireAuth redirects anonymous`. Run `npm run web:test` → PASS.
- `npm run web:typecheck`, `npm run web:build` → PASS.

**Risks:** Guard causing redirect loops on `/login` itself — exclude auth routes. Placeholder mistaken for working feature — mark clearly (TODO + disabled submit).

**Estimated complexity:** Low-Medium.

**Suggested PR boundary:** "feat(web): forgot-password placeholder + wire auth CTA".

---

### Milestone 12 — Documentation & design-system reconciliation

**Purpose:** Update service contracts and design docs to match what was built; promote any genuinely new reusable pattern (e.g., auth form/page) from Proposed to as-built. No code behavior change.

**Dependencies:** All prior milestones merged.

**Affected services:** `docs/`.

**Files:**
- Modify `docs/services/api.md` — document `/api/v1/auth/*` endpoints, envelopes, status codes, cookie behavior.
- Modify `docs/services/web.md` — document `AuthProvider`, `lib/auth.ts`, routes.
- Modify/create `docs/design/` — add the auth-form/page pattern **only if** new reusable classes were introduced in M9 (reuse-audit first); update `components/modal.md` only if a modal was actually built (it wasn't in this page-based plan — leave as Proposed).
- Modify `docs/adr/0007-authentication-approach.md` — mark "Accepted / Implemented" with final endpoint list.

**Verification:**
- Docs link-check (no broken relative links); `docs/REPOSITORY-MAP.md` still accurate.
- Cross-check each endpoint in `docs/services/api.md` against `internal/http/auth_handler.go` route registrations.

**Risks:** Docs drifting from code — write docs from the merged code, not the plan. Over-documenting (adding a design language) — reuse-audit gate.

**Estimated complexity:** Low.

**Suggested PR boundary:** "docs(auth): service contracts + design-system reconciliation".

---

## Dependency Graph (milestone order)

```
M0 ─┬─> M1 ─> M2 ─> M3 ─> M4 ─> M5 ─> M6 (backend complete)
    │                        │
    └────────────────────────┴─> M7 ─> M8 ─> M9 ─┬─> M10
                                                  └─> M11
M6 ──────────────────────────────────────────────> M10 (needs backend Google)
(all) ──────────────────────────────────────────> M12 (docs)
```

Backend track: M0→M1→M2→M3→M4→M5→M6. Frontend track can start once M4 (contract) merges: M7→M8→M9→{M10 needs M6, M11}. M12 last.

## Cross-Cutting Risks (whole feature)

1. **Security correctness** — token storage, rotation, CORS/cookie flags, user-enumeration, alg-confusion. Mitigated per-milestone; recommend a `/security-review` pass before merging M4–M6 and M7.
2. **No CI yet** — `.github/workflows` absent; verification is local per PR. Consider a small CI addition (separate, gated PR) so auth PRs are gated — out of scope here but flagged.
3. **Design-language drift** — every frontend PR runs a reuse/design self-review; only M12 may extend design docs.
4. **Scope creep toward anti-goals** — auth must not silently enable follow-graph/monetization features; keep to sessions only.
5. **Secret management in prod** — dev uses env; production secret store is out of scope and must be handled by infra (human-gated) before deploy.

## Self-Review (against the feature spec)

- Backend register/login/refresh/logout → M3 (logic) + M4 (HTTP). ✅
- Google OAuth SSO → M6 (backend) + M10 (frontend). ✅
- JWT authentication → M2 (tokens) + M5 (middleware). ✅
- Session management → M1 (refresh table) + M3 (rotation/revoke) + M5. ✅
- Login/Register pages → M9. ✅ Google Sign-In → M10. ✅ Forgot-password placeholder → M11. ✅
- Responsive layout, form validation, loading/error states → M9 (+ M7 client states). ✅
- Reuse design system / no new language / update docs if new patterns → constraints + M9 self-review + M12. ✅
- Reuse existing components → M9 (FormField from `.field`), `.authCta` wiring (M11). ✅
```
