# Service Contract — `apps/web` (Discovery Frontend)

> **Why this doc exists:** a single at-a-glance contract for the Next.js frontend — what
> it renders, what it depends on, and the rules for changing it — so UI work respects
> the API contract and the frontend conventions. Binding rules live in
> [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) and
> [`docs/standards/frontend-architecture.md`](../standards/frontend-architecture.md).

## Purpose

Deliver the TikFood discovery UX: a map-first, dish-first browsing experience over the
venues served by `apps/api`. Discovery only.

## Interface

This service exposes **UI**, not an API. Its "contract" is what it consumes and the
boundaries it honors.

- **Entry:** `app/page.tsx` (server component) fetches initial venues and passes
  `initialVenues` to the client `components/DiscoveryExperience.tsx`.
- **Surfaces:** map (MapLibre GL, lazy-imported), venue list/rail, venue detail,
  filters; **auth** — `/login`, `/register`, `/forgot-password` (placeholder),
  `/auth/google/callback`.
- **Styling:** plain CSS in `app/globals.css`
  ([ADR-0005](../adr/0005-plain-css-frontend.md)).

### Authentication ([ADR-0007](../adr/0007-authentication-approach.md))

- **`components/auth/AuthProvider.tsx`** — session context wrapped around the app in
  `app/layout.tsx`; exposes `useAuth()` with `{ user, status, login, register, logout,
  loginWithGoogleUrl }`. Silent bootstrap via the refresh cookie on mount.
- **`lib/auth.ts`** — the auth client (`register`/`login`/`refresh`/`logout`/`getMe`/
  `googleLoginUrl`); keeps the access token in memory, sends `credentials: "include"`,
  attaches `Bearer`, and does one silent refresh + retry on 401. Discovery calls in
  `lib/api.ts` stay public/unauthenticated.
- **`lib/validation.ts`** + `components/auth/{FormField,AuthForm,GoogleSignInButton,
  AuthCta,RequireAuth}.tsx` — reuse existing tokens/classes; no new design language.

## Dependencies

- **`apps/api`** — the only data source. **All** access goes through
  [`lib/api.ts`](../../apps/web/lib/api.ts) (`fetchDiscoveryVenues`,
  `fetchVenueDetail`, `getDiscoveryVenues`); no raw `fetch` scattered in components.
  Respects the `{ data, error }` envelope. When `NEXT_PUBLIC_API_URL` is unset, falls
  back to `fallbackVenues`.
- **MapLibre GL** for map rendering. External routing/tiles (e.g. OSRM demo) are
  **not** production-safe — see `IMPROVEMENTS.md`.
- No React Query / React Hook Form / Zod today (approved to add only when real
  caching/mutations/forms appear).

## Consumers

End users (browser). No other service depends on `apps/web`.

## Change rules

- **Types mirror the API.** Keep `lib/api.ts` types in sync with the Go JSON tags
  (snake_case, e.g. `avg_price_min_vnd`, `social_videos`). If the API envelope changes,
  update here in the same change ([ADR-0003](../adr/0003-data-error-response-envelope.md)).
- **Route API calls through `lib/api.ts`** only.
- **Styling:** reuse existing `globals.css` classes; inline `style` only for computed
  values (e.g. marker positions); no new CSS framework without justification
  ([ADR-0005](../adr/0005-plain-css-frontend.md)).
- **Naming:** descriptive components (`VenueDetail`, `VenueMap`, `VenueRailCard`); never
  `Component1`/`DataCard2`.
- **Accessibility:** keep `aria-label` on icon-only buttons; label inputs. Vietnamese UI
  copy is the norm.
- **No anti-goal UI** (cart/checkout/booking/chat/etc.) — human approval required.
- Verify with `make verify-web` (typecheck + Vitest + build), or the individual
  `npm --workspace apps/web run {typecheck,test,build}`.

## Current state

`components/DiscoveryExperience.tsx` (~1350 lines) concentrates map/filters/detail/cards;
an approved split into `VenueMap`/`VenueDetail`/`VenueRailCard`/filter controls/format
utils is pending. A Vitest + React Testing Library harness is in place (auth client,
provider, forms, guard, plus `VenueList`). Some UI data is placeholder — see
`IMPROVEMENTS.md` (*Fabricated UI data*); do not treat ratings/reviews/photos as real
until wired end-to-end.
