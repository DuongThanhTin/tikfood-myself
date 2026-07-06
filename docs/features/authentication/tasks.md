# Authentication — Tasks (Milestone Breakdown)

One-line purpose: the approved implementation roadmap restated for this package — 12 milestones, each a single PR, with purpose, dependencies, files, verification, risks, and complexity.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

Turn the design into an ordered, reviewable build sequence. Each milestone is the
smallest unit that carries its own test cycle and PR. This document is the canonical
task list; the operational rollout view is in [`release.md`](release.md), and the
full step-level detail lives in the plan file referenced at the bottom.

---

## Existing implementation

No authentication tasks exist yet. The repository is at the state described in
[`overview.md`](overview.md): public read-only discovery, no auth code, no user
tables, an inert `.authCta` placeholder in the web UI.

## Proposed changes

Twelve milestones, `M0`–`M12`. Backend track `M0→M6`; frontend track `M7→M11`
(startable once `M4` freezes the contract; `M10` also needs `M6`); `M12` is docs.

### Dependency graph

```
M0 ─> M1 ─> M2 ─> M3 ─> M4 ─> M5 ─> M6          (backend)
                          │
                          └─> M7 ─> M8 ─> M9 ─┬─> M10   (M10 also needs M6)
                                              └─> M11
(all) ─────────────────────────────────────────> M12   (docs)
```

---

### M0 — ADR + dependencies + config scaffolding
- **Purpose:** record ADR-0007; add Go deps and config keys (unused yet).
- **Dependencies:** none.
- **Affected services:** `apps/api`, `docs/`.
- **Files:** `docs/adr/0007-authentication-approach.md`; `apps/api/go.mod`/`go.sum`; `apps/api/internal/config/config.go`; `.env.example`.
- **Verification:** `go build ./...`, `go vet ./...`; `TestConfigDefaultsAuth`.
- **Risks:** weak/empty `JWT_SECRET` — guard in M2.
- **Complexity:** Low.
- **PR boundary:** "chore(auth): ADR-0007 + deps & config scaffolding".

### M1 — Users & refresh-token schema + repositories
- **Purpose:** persist users and refresh tokens; repository interfaces + Postgres + in-memory impls.
- **Dependencies:** M0; migration approval.
- **Affected services:** `apps/api`, PostgreSQL.
- **Files:** `apps/api/migrations/007_auth_users.sql`, `008_auth_refresh_tokens.sql`; `apps/api/internal/auth/model.go`, `repository.go`, `memory_repository.go`; `apps/api/internal/storage/postgres/user_repository.go`, `refresh_token_repository.go`.
- **Verification:** `go build/vet`; `TestMemoryUserRepository_CreateAndFindByEmail`, `_DuplicateEmailRejected`, `_FindByGoogleSub`, `TestMemoryRefreshTokenRepository_StoreFindRevoke`; manual migration apply.
- **Risks:** `citext` availability; cascade deletes.
- **Complexity:** Medium.
- **PR boundary:** "feat(auth): user & refresh-token schema + repositories".

### M2 — Password hashing + JWT token service
- **Purpose:** pure, unit-testable security core.
- **Dependencies:** M0 (config), M1 (types).
- **Affected services:** `apps/api`.
- **Files:** `apps/api/internal/auth/password.go`, `token.go`.
- **Verification:** `TestHashPassword_VerifyRoundTrip`, `TestIssueAndParseAccessToken_RoundTrip`, `_ExpiredFails`, `_TamperedFails`, `_RejectsNonHS256`, `TestNewRefreshTokenValue_HashDeterministic`.
- **Risks:** alg-confusion; bcrypt 72-byte cap; entropy.
- **Complexity:** Medium.
- **PR boundary:** "feat(auth): password hashing + JWT token service".

### M3 — Auth service (register/login/refresh/logout)
- **Purpose:** orchestrate primitives + repos into use-cases; domain validation.
- **Dependencies:** M1, M2.
- **Affected services:** `apps/api`.
- **Files:** `apps/api/internal/auth/service.go`.
- **Verification:** `TestRegister_*`, `TestLogin_*` (incl `_UnknownUserSameError`), `TestRefresh_RotatesAndOldRevoked`, `_ExpiredRejected`, `_RevokedRejected`, `TestLogout_Idempotent`.
- **Risks:** user enumeration; refresh replay after rotation.
- **Complexity:** Medium-High.
- **PR boundary:** "feat(auth): auth service".

### M4 — HTTP handlers + routes (contract freeze)
- **Purpose:** expose register/login/refresh/logout with envelope + cookie handling.
- **Dependencies:** M3; auth approval.
- **Affected services:** `apps/api`.
- **Files:** `apps/api/internal/http/auth_handler.go`; modify `router.go`, `errors.go`, `internal/app/container.go`.
- **Verification:** `TestRegister_201SetsCookie`, `_DuplicateReturns422`, `TestLogin_200`, `_401OnBadCreds`, `TestRefresh_200RotatesCookie`, `_401WithoutCookie`, `TestLogout_ClearsCookie`, `TestAuthResponsesUseEnvelope`; discovery tests still pass.
- **Risks:** cookie flags for local http; envelope drift.
- **Complexity:** Medium-High.
- **PR boundary:** "feat(auth): auth HTTP endpoints" (**freezes the contract**).

### M5 — Auth middleware + CORS + /me
- **Purpose:** validate access tokens, inject user, add credentialed CORS, first protected route.
- **Dependencies:** M4, M2, M0.
- **Affected services:** `apps/api`.
- **Files:** modify `apps/api/internal/http/middleware.go`, `router.go`, `auth_handler.go`.
- **Verification:** `TestAuthMiddleware_ValidTokenAllows`, `_MissingTokenReturns401`, `_ExpiredReturns401`, `TestMe_ReturnsCurrentUser`, `TestCORS_AllowsConfiguredOrigin`, `_RejectsUnknownOrigin`, `TestDiscoveryRoutesStayPublic`.
- **Risks:** CORS misconfig; discovery accidentally gated.
- **Complexity:** Medium.
- **PR boundary:** "feat(auth): middleware + CORS + /me".

### M6 — Google OAuth SSO (backend)
- **Purpose:** server-side Authorization-Code login; provision/link user; issue tokens.
- **Dependencies:** M3, M4, M5, M0; external-call approval.
- **Affected services:** `apps/api`, Google (external).
- **Files:** `apps/api/internal/auth/oauth_google.go`; modify `service.go`, `auth_handler.go`.
- **Verification:** `TestLoginWithGoogle_NewUser`, `_ExistingByGoogleSub`, `_LinksVerifiedEmail`, `TestGoogleCallback_StateMismatchRejected`, `TestGoogleLogin_RedirectsToGoogle` (fake authenticator); live dev check.
- **Risks:** CSRF on callback; takeover via unverified email; secret exposure.
- **Complexity:** High.
- **PR boundary:** "feat(auth): Google OAuth SSO (backend)".

### M7 — Frontend auth client + token handling
- **Purpose:** typed auth calls; Bearer injection + one silent refresh on 401; add test harness.
- **Dependencies:** M4, M5, M0.
- **Affected services:** `apps/web`.
- **Files:** `apps/web/lib/auth.ts`; modify `lib/api.ts`; `package.json` (Vitest+RTL devDeps, `web:test`); `vitest.config.ts`.
- **Verification:** `login stores access token`, `getMe attaches Bearer`, `api retries once after 401 then refresh`, `public venue fetch works with no token`, `refresh failure clears token`; `npm run web:typecheck/build/test`.
- **Risks:** infinite refresh loop (single-retry guard).
- **Complexity:** Medium.
- **PR boundary:** "feat(web): auth client + refresh + test harness".

### M8 — AuthProvider (session context)
- **Purpose:** single React context; silent bootstrap via refresh cookie; login/register/logout actions.
- **Dependencies:** M7.
- **Affected services:** `apps/web`.
- **Files:** `apps/web/components/auth/AuthProvider.tsx`; modify `app/layout.tsx`.
- **Verification:** `bootstraps authenticated/anonymous`, `login updates user`, `logout clears user`; typecheck/build/test.
- **Risks:** SSR/hydration mismatch (stable loading state).
- **Complexity:** Medium.
- **PR boundary:** "feat(web): AuthProvider session context".

### M9 — Login & Register pages
- **Purpose:** `/login` + `/register` reusing design tokens; validation, loading/error, responsive.
- **Dependencies:** M8, M7.
- **Affected services:** `apps/web`.
- **Files:** `apps/web/lib/validation.ts`; `components/auth/FormField.tsx`, `AuthForm.tsx`; `app/login/page.tsx`, `app/register/page.tsx`; modify `app/globals.css`.
- **Verification:** `validateEmail/validatePassword` boundaries; `shows error on invalid email`, `submit disabled while loading`, `renders API error text`, `successful login redirects`; typecheck/build/test; manual responsive.
- **Risks:** design-language drift; a11y regressions.
- **Complexity:** Medium-High.
- **PR boundary:** "feat(web): login & register pages".

### M10 — Google Sign-In (frontend) + callback
- **Purpose:** wire Google button + handle callback landing the session.
- **Dependencies:** M6, M8, M9.
- **Affected services:** `apps/web`.
- **Files:** `apps/web/components/auth/GoogleSignInButton.tsx`; `app/auth/google/callback/page.tsx`; modify login/register pages.
- **Verification:** `google button links to backend login url`, `callback success sets session`, `callback failure shows error`; typecheck/build/test; live dev check.
- **Risks:** hand-off token in URL history; button design drift.
- **Complexity:** Medium.
- **PR boundary:** "feat(web): Google Sign-In + callback".

### M11 — Forgot-password placeholder + wire authCta + guard
- **Purpose:** placeholder page; wire the inert `.authCta`; minimal route guard util.
- **Dependencies:** M8, M9.
- **Affected services:** `apps/web`.
- **Files:** `apps/web/app/forgot-password/page.tsx`; modify `components/DiscoveryExperience.tsx`; `components/auth/RequireAuth.tsx`; modify `IMPROVEMENTS.md`.
- **Verification:** `authCta links to login when anonymous`, `shows logout when authenticated`, `RequireAuth redirects anonymous`; typecheck/build/test.
- **Risks:** redirect loop on `/login`; placeholder mistaken for real feature.
- **Complexity:** Low-Medium.
- **PR boundary:** "feat(web): forgot-password placeholder + wire auth CTA".

### M12 — Documentation & design reconciliation
- **Purpose:** update service contracts + design docs to match built code; promote new reusable patterns only if introduced.
- **Dependencies:** all prior.
- **Affected services:** `docs/`.
- **Files:** modify `docs/services/api.md`, `docs/services/web.md`, `docs/design/*` (only if new pattern), `docs/adr/0007-authentication-approach.md`.
- **Verification:** link-check; each endpoint in `docs/services/api.md` matches route registrations.
- **Risks:** docs drift; over-documenting a new language.
- **Complexity:** Low.
- **PR boundary:** "docs(auth): service contracts + design reconciliation".

## Open questions

1. **Sequencing** — accept the backend-first order, or start frontend scaffolding
   (M7 harness) before M4 merges using a stubbed contract?
2. **CI** — introduce test-gating before or after these milestones? See [`testing.md`](testing.md) / [`release.md`](release.md).
3. **Feature flag** — gate auth routes behind a runtime flag for staged rollout?

## Related repository documentation

- The approved plan (step-level detail): `docs/superpowers/plans/2026-07-05-authentication-system.md`.
- [`/CLAUDE.md`](../../../CLAUDE.md) — branch/PR discipline, "Before Finishing" rules.
- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — per-PR gates.
- [`docs/recipes/`](../../recipes/) — task recipes (add-API, refactor) referenced during execution.
- Sibling docs: [`overview.md`](overview.md), [`architecture.md`](architecture.md), [`testing.md`](testing.md), [`release.md`](release.md).
