# Authentication — Specification

One-line purpose: the authoritative list of functional and non-functional requirements the authentication feature must satisfy, with acceptance criteria.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

`specification.md` translates the intent in [`overview.md`](overview.md) into
testable requirements. Each requirement has a stable ID (FR-n / NFR-n) so that
[`testing.md`](testing.md) can map tests to it and [`tasks.md`](tasks.md) can map
milestones to it. It states *what* must be true, not *how* to build it — mechanism
lives in [`architecture.md`](architecture.md), [`api.md`](api.md), and
[`database.md`](database.md).

---

## Existing implementation

- The product is public, read-only discovery. There are **no** user accounts,
  credentials, sessions, or protected endpoints today.
- The backend already enforces a `{data,error}` response envelope and returns
  structured error codes; auth requirements below must conform to that envelope.
- The frontend has an inert `.authCta` ("Đăng nhập ngay") placeholder and no auth
  screens, validation, or session state.
- UI copy convention is Vietnamese; technical identifiers are English.
- No requirement below may break or gate any existing discovery behavior.

## Proposed changes

### Functional requirements

| ID | Requirement | Acceptance criteria |
|----|-------------|---------------------|
| FR-1 | **Register (email/password)** | Given a unique, valid email and a password ≥ 8 chars, registration creates a user and returns an authenticated session (access token + refresh cookie). Duplicate email is rejected with a domain error; the account is not created twice. |
| FR-2 | **Login (email/password)** | Valid credentials return an authenticated session. Invalid credentials return a single uniform error that does not reveal whether the email exists. |
| FR-3 | **Logout** | An authenticated user can log out; the current refresh token is revoked and its cookie cleared. Logout is idempotent (repeat calls do not error). |
| FR-4 | **Refresh / rotation** | A valid, unexpired, unrevoked refresh token yields a new access token and a **rotated** refresh token; the previous refresh token is revoked. Expired/revoked/absent tokens are rejected. |
| FR-5 | **Session management** | Sessions are represented by persisted refresh tokens; they can be individually revoked (logout) and, per open question, optionally revoked in bulk ("log out everywhere"). |
| FR-6 | **JWT authentication** | Protected endpoints accept a valid, unexpired access token (`Authorization: Bearer`) and reject missing/expired/tampered/wrong-algorithm tokens with `401`. |
| FR-7 | **Current user (`/me`)** | An authenticated request can retrieve the current user's non-sensitive profile; the password hash is never returned. |
| FR-8 | **Google OAuth SSO** | A user can authenticate via Google (Authorization-Code flow). A first-time Google user is provisioned; a returning Google user is matched; an existing verified email may be linked. The result is the same session shape as FR-1/FR-2. |
| FR-9 | **Input validation** | Email must be a syntactically valid address; password length 8–72 bytes (bcrypt limit). Invalid input returns a validation error with a field-scoped, Vietnamese message; no server error. |
| FR-10 | **Login & Register pages** | Dedicated `/login` and `/register` routes render responsive forms with the existing design language, cross-link each other, and offer Google Sign-In. |
| FR-11 | **Forgot-password placeholder** | A `/forgot-password` page exists, clearly marked as a future capability; it performs no backend action in this iteration. |
| FR-12 | **Loading states** | Every submit shows a loading state (control disabled + label swap) and prevents double submission. |
| FR-13 | **Error states** | API and validation errors render inline near the relevant field/form using the existing error pattern; errors clear on retry. |
| FR-14 | **Responsive layout** | Auth screens are usable at 375 / 768 / 1280 px, reusing the existing responsive grid behavior. |
| FR-15 | **Session persistence across reload** | After a full page reload, an authenticated user is silently restored via the refresh cookie without re-entering credentials. |
| FR-16 | **Anonymous access preserved** | Discovery remains fully usable without authentication; auth is additive and optional. |

### Non-functional requirements

| ID | Requirement | Acceptance criteria |
|----|-------------|---------------------|
| NFR-1 | **Password storage** | Passwords are stored only as strong one-way hashes; plaintext is never persisted or logged. |
| NFR-2 | **Token integrity** | Access tokens are signed and verified with a fixed algorithm; algorithm-confusion and tampering are rejected. |
| NFR-3 | **Token transport** | Refresh tokens travel only in an httpOnly `Secure` cookie; access tokens are not persisted to durable storage by the client. |
| NFR-4 | **No user enumeration** | Login and registration do not leak account existence via message or status differences. |
| NFR-5 | **No secret leakage** | Secrets, tokens, password hashes, and raw SQL never appear in logs, error responses, PRs, or summaries. |
| NFR-6 | **Envelope stability** | All auth responses use the immutable `{data,error}` envelope; no existing response shape changes. |
| NFR-7 | **Layering integrity** | Handlers validate/parse, services hold logic (no `gin`), repositories persist only — no forbidden edges. |
| NFR-8 | **CORS correctness** | Cross-origin credentialed requests are allowed only for an explicit origin list; wildcard-with-credentials is never used. |
| NFR-9 | **Performance** | Auth operations complete within typical API latency budgets; hashing cost is tuned to be secure yet responsive. |
| NFR-10 | **Observability** | Auth events are logged via the existing structured logger with request IDs and no sensitive fields. |
| NFR-11 | **Accessibility & i18n** | Inputs are labelled; icon-only controls have `aria-label`; all user-facing copy is Vietnamese. |

## Open questions

- **OQ-SPEC-1:** Password policy — is length-only (≥ 8) sufficient, or add character-class rules / breached-password checks?
- **OQ-SPEC-2:** Display name — required at registration or optional/derived from email or Google profile?
- **OQ-SPEC-3:** Confirm exact token lifetimes (~15 min access / ~30 day refresh).
- **OQ-SPEC-4:** Logout — offer "log out everywhere" (revoke all sessions) now, or single-session only?
- **OQ-SPEC-5:** Google linking — link to an existing account only when Google reports the email as verified (assumed) — confirm the takeover stance. See [`security.md`](security.md).
- **OQ-SPEC-6:** Rate limiting / lockout on login/register — in this feature or a fast-follow? See [`security.md`](security.md).

## Related repository documentation

- [`overview.md`](overview.md) — feature intent and scope.
- [`/CLAUDE.md`](../../../CLAUDE.md) — anti-goals, security/git, naming, "Before Finishing".
- [`apps/api/CLAUDE.md`](../../../apps/api/CLAUDE.md) — envelope, layering, error codes.
- [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — plain CSS, `lib/api.ts`, Vietnamese copy, a11y.
- [`docs/standards/backend/`](../../standards/backend/) — request/response, errors, testing standards.
- [`docs/standards/application-security.md`](../../standards/application-security.md) — product security (auth requires approval + ADR).
- [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md) — completion gates these requirements are checked against.
- [`docs/tikfood/mvp-scope.md`](../../tikfood/mvp-scope.md) · [`docs/tikfood/anti-goals.md`](../../tikfood/anti-goals.md) — scope boundaries.
