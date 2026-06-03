# Frontend Architecture Standard

Frontend code lives in `apps/web`.

TikFood frontend is a Next.js App Router app for dish-first and map-first discovery. It should feel like an operational discovery product, not a delivery or checkout app.

## Current Frontend Baseline

- Framework: Next.js App Router
- Language: TypeScript
- Styling: plain CSS for MVP
- API access: `apps/web/lib/api.ts`
- Main route: `apps/web/app/page.tsx`

## Structure

```text
apps/web/app/          Routes, layouts, global styles
apps/web/components/   Reusable UI components
apps/web/lib/          API access, typed helpers, shared utilities
apps/web/public/       Static files
```

As the app grows, feature modules may be introduced:

```text
apps/web/features/discovery/
apps/web/features/venue/
apps/web/features/search/
```

Do not create feature folders before they carry real complexity.

## Component Rules

- Keep page files focused on route composition.
- Move repeated UI into `components`.
- Keep data fetching helpers in `lib`.
- Use typed props.
- Include loading and error states for networked client UI.
- Keep text readable on mobile and desktop.
- Avoid marketing-only pages when building product workflows.

## Data Access

API client functions should:

- Return typed data.
- Hide raw fetch details from components.
- Handle fallback or error behavior explicitly.
- Avoid leaking backend response envelope handling into many components.

Current pattern:

```text
apps/web/lib/api.ts
```

Future pattern with TanStack Query:

```text
features/<feature>/queries.ts
features/<feature>/types.ts
```

## Styling

Current MVP uses plain CSS in `app/globals.css`.

Future preferred stack:

- Tailwind CSS
- shadcn/ui
- MapLibre GL or Mapbox GL for real map rendering

Rules:

- Do not add a design dependency without a feature needing it.
- Keep cards at 8px border radius or less unless the design system changes.
- Avoid one-note palettes dominated by a single hue.
- Make layout stable across mobile and desktop.

## Product Guardrails

Do not add UI for:

- Delivery
- Cart
- Orders
- Checkout
- Payment
- Booking or reservations
- In-app chat
- Social follow graph
- Creator monetization
- Livestream

If a request asks for these areas, block or request human approval.

## Testing And Checks

Required for frontend changes:

```bash
npm run web:typecheck
```

For production build validation:

```bash
npm run web:build
```

Add component or integration tests when user workflows become non-trivial.
