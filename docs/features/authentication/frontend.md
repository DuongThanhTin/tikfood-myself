# Authentication — Frontend

One-line purpose: how the Next.js web app gains authentication — session state, the auth API client, new routes, and how it all reuses existing conventions.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

This document describes the frontend architecture for authentication: the pieces
introduced in `apps/web`, their responsibilities, and how they fit the app's
existing patterns without changing them. It is the map a web engineer reads before
building. Visual/interaction design (wireframes, states, copy) lives in
[`ui.md`](ui.md); the reuse audit lives in [`component-reuse.md`](component-reuse.md);
the HTTP contract this consumes lives in [`api.md`](api.md). Components are described
**by responsibility**, not as code.

---

## Existing implementation

- **Framework:** Next.js **App Router**, TypeScript (strict), plain CSS in
  `app/globals.css` (dark-first glassmorphism design system). No Tailwind/shadcn.
- **Routing:** a single route, `/` (the discovery experience).
- **Components:** essentially two — `components/DiscoveryExperience.tsx` (large,
  client, heavy local `useState`) and `components/VenueList.tsx`.
- **Data access:** all fetches go through `lib/api.ts`, which speaks the
  `{data,error}` envelope and reads `NEXT_PUBLIC_API_URL`. **No auth header is sent.**
- **State/data libs:** none beyond React state — **no** React Query/SWR, **no**
  react-hook-form, **no** Zod.
- **Auth UI:** none. An **inert `.authCta`** ("Đăng nhập ngay") placeholder button
  exists in the discovery shell but does nothing (a known accessibility gap).
- **Rules of the app** ([`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md)): reuse
  `globals.css` classes, keep types in sync with Go JSON tags (snake_case), UI copy
  in Vietnamese, `aria-label` on icon-only buttons, no new CSS framework, and
  introduce React Query / RHF+Zod only when real caching/forms appear.

---

## Proposed changes

### New routes (App Router, D3)

| Route | Purpose | Access |
|---|---|---|
| `/login` | Email/password login + Google Sign-In | public |
| `/register` | Account creation + Google Sign-In | public |
| `/forgot-password` | **Placeholder** for a future flow (clearly marked, no backend) | public |
| `/auth/google/callback` | Completes the Google round trip, establishes the session | public |

Discovery (`/`) stays public and unchanged; it simply gains a working sign-in entry
point (the `.authCta` is wired to `/login`).

### `AuthProvider` — the single source of session state

- A client-side React context exposing: `user`, `status` (`loading` |
  `authenticated` | `anonymous`), and actions `login`, `register`, `logout`, plus
  `googleLoginUrl`.
- **Silent bootstrap on mount:** attempts a refresh (using the httpOnly cookie) and,
  if it succeeds, loads the current user (`getMe`) — so a reload restores the
  session without the user re-entering credentials.
- Consumed via a `useAuth()` hook. Wrapped around the app in `app/layout.tsx`.

### `lib/auth.ts` — the auth API client

Responsibilities (all via the existing fetch approach, `credentials: 'include'` so
the cookie flows):

- `register`, `login`, `refresh`, `logout`, `getMe`, and `googleLoginUrl`.
- Holds the **access token in memory** (module-level accessor: set/get/clear) — never
  in localStorage — matching the security posture in [`security.md`](security.md).
- Exposes typed shapes that mirror the Go JSON tags in **snake_case** (e.g.
  `access_token`, `display_name`).

### `lib/api.ts` — minimal, backward-compatible changes

- **Attach `Authorization: Bearer`** when an access token is present; send nothing
  extra when it is not (so existing public discovery calls are unaffected).
- **One silent refresh on 401:** if a request to a protected endpoint returns 401,
  attempt a single `refresh()` and retry once. A **retry-once guard** prevents
  infinite refresh loops; a failed refresh clears the in-memory token and drops the
  session to `anonymous`.
- The `{data,error}` envelope handling is preserved unchanged.

### `lib/validation.ts` — pure validators (D4, no form library)

- Small pure functions: `validateEmail`, `validatePassword` (min 8, ≤72 bytes),
  `validateRequired`, returning **Vietnamese** messages.
- Used with controlled component state — **no** new form/validation dependency is
  introduced, honoring the current frontend standard.

### `components/auth/*` — responsibilities

| Component | Responsibility |
|---|---|
| `AuthProvider` | Session context + silent bootstrap (above) |
| `AuthForm` | Shared controlled-state form shell (email/password [+ display name on register]); submit via `useAuth`; loading + error handling |
| `FormField` | A single labelled field built from the existing `.field` class + native input; renders inline `.errorText`; wires `<label>` + `aria-describedby` |
| `GoogleSignInButton` | Navigates to the backend Google login URL, styled within existing button variants |
| `RequireAuth` (util) | Redirects `anonymous` users to `/login`; **excludes** the auth routes themselves to avoid redirect loops |

### State, loading, and error conventions

- **Controlled state** for inputs (value + onChange), mirroring the discovery
  filters.
- **Loading** shown by disabling the submit control and swapping its label (e.g.
  "Đăng nhập" → "Đang đăng nhập"), the existing pattern — no new spinner system.
- **Errors** rendered inline with the existing `.errorText` treatment (danger color),
  preferring the API `error.message` with a Vietnamese fallback; errors clear on
  retry.

### SSR / hydration note

- `AuthProvider` and the auth pages are **client components**. On first paint the
  provider renders a stable `loading` state to avoid a hydration mismatch and a flash
  of the wrong (anonymous vs authenticated) UI; it resolves after the bootstrap
  refresh completes.

### Wiring the existing placeholder

- The inert `.authCta` becomes a real entry point: it links to `/login` when
  anonymous and shows the user's display name + a logout control when authenticated —
  all within the existing design language (see [`ui.md`](ui.md)).

---

## Open questions

1. **Post-login destination:** always redirect to `/`, or return the user to the
   page they came from (needs a validated `redirect` param to prevent open-redirect)?
2. **Session-restore UX:** is a brief "restoring session…" state acceptable, or do we
   need a dedicated splash to avoid any flash of anonymous UI on reload?
3. **Guarded surfaces:** are there any authenticated-only pages in this MVP, or is
   `RequireAuth` shipped only as a utility for future use? (Discovery stays public.)
4. **Google callback hand-off:** confirm the exact mechanism the backend chooses
   (one-time code vs URL fragment) so the callback page can be finalized — see
   [`security.md`](security.md) and [`api.md`](api.md).
5. **Testing harness:** Vitest + React Testing Library are proposed as new devDeps
   (the repo has zero frontend tests today) — confirm acceptable. See
   [`testing.md`](testing.md).
6. **Logout UX in the discovery shell:** display name + logout link vs a compact
   avatar/menu — kept within the existing language.

---

## Related repository documentation

- [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — web app rules (plain CSS, lib/api.ts, snake_case types, Vietnamese copy).
- [`docs/standards/frontend-architecture.md`](../../standards/frontend-architecture.md) — frontend architecture standard.
- [`docs/design/`](../../design/) — canonical design system (tokens, components, patterns).
- [`/CLAUDE.md`](../../../CLAUDE.md) — naming, before-finishing rules.
- [`docs/ai/AI-CONTRACT.md`](../../ai/AI-CONTRACT.md) — change discipline, no breaking API changes.
- Sibling docs: [`ui.md`](ui.md), [`component-reuse.md`](component-reuse.md), [`api.md`](api.md), [`security.md`](security.md), [`testing.md`](testing.md).
