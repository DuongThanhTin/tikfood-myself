# Authentication — Security

One-line purpose: the threat model and the controls that make TikFood's authentication safe to operate, plus what is deliberately out of scope.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

This document defines the security posture for the Authentication feature: the
assets we are protecting, the threats against them, and the concrete controls the
design commits to. It is the reference a reviewer uses to decide whether the auth
implementation is safe to merge, and the checklist a `/security-review` pass runs
against. It stays at the level of *what protection is required and why*; exact code,
libraries, and configuration values are described only where necessary to make a
control unambiguous, and otherwise deferred to [`architecture.md`](architecture.md),
[`api.md`](api.md), and [`database.md`](database.md).

Authentication is a **human-approval-gated** area. Approval to proceed (with an
accompanying ADR) was granted during Discovery; this document is part of the record
that justifies that approval.

---

## Existing implementation

- TikFood is today a **public, read-only discovery application**. There are no user
  accounts, no credentials, no sessions, and therefore no authentication attack
  surface beyond ordinary API abuse.
- **No secrets live in code.** Backend configuration is read from environment
  variables (`internal/config`); there is no secret manager and no `.env` reading by
  agents.
- The `{data,error}` response envelope already avoids leaking internals — errors
  return a stable `code` + human message, never raw SQL or stack traces.
- Database access uses parameterized queries via the pgx driver (no string-built
  SQL), so the existing surface is not SQL-injection prone.
- Governance already treats auth as sensitive: [`docs/security.md`](../../security.md)
  and [`docs/standards/application-security.md`](../../standards/application-security.md)
  both require **human approval + a dedicated ADR** before any auth/authz is added,
  and [`docs/ai/AI-CONTRACT.md`](../../ai/AI-CONTRACT.md) §4 lists auth, migrations,
  and new external network calls as approval gates.

There are no existing auth-specific controls to reuse — everything below is new.

---

## Proposed changes

### Credential handling

- **Password hashing:** bcrypt at cost 12. Passwords are capped/validated at **72
  bytes** (bcrypt's input limit) to avoid silent truncation; minimum length 8.
  Plaintext passwords are never stored, never logged, and never returned.
- **Password verification** uses a constant-time comparison (bcrypt's own compare);
  no early-exit string comparison.

### Tokens and sessions (D1, D2)

- **Access token:** short-lived (~15 min) HS256 JWT sent as `Authorization: Bearer`.
  Verification **explicitly enforces the HS256 algorithm** and rejects any other
  `alg` value — this closes the classic *algorithm-confusion* (`alg: none` / RS→HS)
  attack. Claims are minimal (subject, email, issuer, iat, exp).
- **Refresh token:** a high-entropy random value (from a CSPRNG). Only its **hash**
  is stored server-side; the raw value lives solely in the client cookie. On use it
  is **rotated** (old token revoked, new one issued) so a captured token has a short
  useful life.
- **Session revocation:** because refresh tokens are persisted, logout and admin
  revocation are real (delete/`revoked_at`). This is what makes "JWT + session
  management" coherent rather than contradictory.

### Cookie posture

- The refresh token is delivered in a cookie that is **httpOnly** (not readable by
  JS → mitigates XSS token theft), **Secure** (HTTPS only), **SameSite=Lax**, and
  **path-scoped** to `/api/v1/auth` so it is only sent to auth endpoints.
- The **access token is held in memory** on the client (never in localStorage) so an
  XSS payload cannot trivially exfiltrate a long-lived credential.

### Authentication flows

- **No user enumeration:** login and registration failures return a **uniform**
  "invalid credentials" / generic error; the API does not reveal whether an email
  exists. Timing is kept comparable (always run a hash comparison path).
- **Google OAuth (D5):** server-side Authorization-Code flow. A **`state`
  parameter** (opaque, single-use, bound to a short-lived cookie) protects the
  callback against CSRF. An existing account is linked to a Google identity **only
  when Google reports the email as verified**, preventing account takeover via an
  unverified address. The Google client secret stays server-side and is never sent
  to the browser or logged.

### Transport and origin controls (D8)

- **CORS** is credentialed and restricted to an **explicit allowed-origin list** from
  config. The wildcard `*` is never combined with credentials.
- All auth endpoints are served over the same `{data,error}` envelope; error bodies
  never include tokens, hashes, secrets, or raw driver errors.

### Secrets and logging

- `JWT_SECRET`, `GOOGLE_CLIENT_ID/SECRET`, TTLs, and allowed origins are read from
  the environment only. In non-dev, the service refuses to start with an empty
  signing secret.
- **Nothing sensitive is logged** — no passwords, tokens, hashes, cookies, or
  Authorization headers. Structured logs keep the existing request-id fields.

### Data-layer controls

- All new queries are parameterized (pgx). Email uniqueness is enforced at the DB
  (case-insensitive), and refresh tokens store only hashes with `expires_at` /
  `revoked_at` for validity checks. See [`database.md`](database.md).

### Control → risk map

| Control | Mitigates (risk) | OWASP-ish category |
|---|---|---|
| bcrypt cost 12, 72-byte cap, min 8 | Credential cracking, silent truncation | Cryptographic / Identification failures |
| HS256 with explicit alg enforcement | Token forgery, alg-confusion (`alg:none`) | Cryptographic failures |
| Short access TTL + rotating refresh | Stolen-token replay window | Identification & auth failures |
| httpOnly + Secure + SameSite=Lax, path-scoped cookie | XSS token theft, CSRF, over-broad exposure | XSS / CSRF |
| Access token in memory only | Persistent XSS exfiltration | XSS |
| Uniform login errors + comparable timing | User enumeration | Identification & auth failures |
| OAuth `state` param (single-use) | OAuth CSRF / login-CSRF | CSRF / SSRF-adjacent |
| Link Google only if `email_verified` | Account takeover via unverified email | Broken access control |
| Explicit CORS allow-list (no `*`+creds) | Cross-origin credential leakage | Security misconfiguration |
| Env-only secrets, no secret logging | Secret leakage via code/logs | Security misconfiguration / Logging |
| Parameterized queries (pgx) | SQL injection | Injection |
| Envelope hides internals | Info disclosure via errors | Security misconfiguration |
| Rate limiting (recommended — see below) | Brute force / credential stuffing | Identification & auth failures |

### Recommended (posture, decision pending)

- **Rate limiting / lockout** on login, register, and refresh. No rate-limiting
  middleware exists today; the design *recommends* it but treats the exact mechanism
  and whether it is a launch blocker as an open question.

### Out of scope (explicitly)

- **Production secret storage** (a real secret manager / KMS) is **infra and
  human-gated** and out of scope for this feature; dev uses env vars. See
  [`release.md`](release.md).
- Password-reset email delivery and email-verification mail (forgot-password is a UI
  placeholder only).
- MFA, WebAuthn/passkeys, and role/permission systems.

---

## Open questions

1. **Rate limiting:** in-process middleware vs a shared store (e.g. Redis); launch
   blocker or fast-follow; thresholds per endpoint?
2. **JWT secret rotation:** a single static `JWT_SECRET` for MVP, or a key-id (`kid`)
   scheme so the secret can rotate without invalidating all live sessions?
3. **CSRF beyond SameSite:** is `SameSite=Lax` + a non-idempotent refresh endpoint
   sufficient, or do we add a double-submit CSRF token for the cookie-based refresh?
4. **Refresh-token theft detection:** on reuse of an already-rotated (revoked) token,
   do we auto **revoke the whole token family** for that user (breach response), or
   just reject the single request?
5. **Email verification for the password flow:** required at MVP, or deferred while
   still enforcing the OAuth `email_verified` linking rule?
6. **Account lockout:** temporary lockout after N failed logins, or rely on rate
   limiting alone?
7. **Cookie `SameSite`/`Domain`:** `Lax` (assumed) vs `Strict`, and domain/path
   scoping given web and API may run on different origins/ports.
8. **Audit logging:** do we record auth events (login success/failure, revocation)
   to a dedicated audit sink, and where?

---

## Related repository documentation

- [`docs/security.md`](../../security.md) — human-approval requirements for auth.
- [`docs/standards/application-security.md`](../../standards/application-security.md) — product security standard; auth needs approval + ADR.
- [`docs/ai/AI-CONTRACT.md`](../../ai/AI-CONTRACT.md) — §2 secrets, §4 approval gates.
- [`/CLAUDE.md`](../../../CLAUDE.md) — Security & Git rules; never read `.env*`.
- [`docs/standards/backend/errors.md`](../../standards/backend/) — error codes / no leakage.
- [`docs/adr/`](../../adr/) — ADR-0003 (envelope); proposed ADR-0007 (auth approach).
- Sibling docs: [`architecture.md`](architecture.md), [`api.md`](api.md), [`database.md`](database.md), [`frontend.md`](frontend.md), [`testing.md`](testing.md), [`release.md`](release.md).
