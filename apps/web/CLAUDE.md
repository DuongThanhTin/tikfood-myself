# CLAUDE.md — apps/web (Next.js Frontend)

Read `/CLAUDE.md` first. This file describes the frontend as it is built today.

## Stack Reality

- Next.js **App Router** + TypeScript.
- **MapLibre GL** for the map (lazy-imported).
- **Plain CSS** in `app/globals.css` — no CSS framework, no design-system library.
- Data through `lib/api.ts` `fetch` wrappers.
- No React Query, no React Hook Form, no Zod today.

## Structure

`app/` (layout, page) → `components/` (`DiscoveryExperience.tsx`,
`VenueList.tsx`) → `lib/api.ts`. The server `page.tsx` passes `initialVenues`
to the client `DiscoveryExperience`.

## API Access

Always go through `lib/api.ts` (`fetchDiscoveryVenues`, `fetchVenueDetail`,
`getDiscoveryVenues`). Do not scatter raw `fetch` calls in components. Respect
the `{ data, error }` envelope. When `NEXT_PUBLIC_API_URL` is unset, fall back
to `fallbackVenues`.

## Types

Keep frontend types in `lib/api.ts` in sync with the Go JSON tags
(snake_case field names such as `avg_price_min_vnd`, `social_videos`).

## Styling

Use plain CSS classes in `globals.css` and reuse existing class patterns. Inline
`style` only for computed positions (e.g. map markers). No new CSS framework
without justification.

## Accessibility & Copy

Keep `aria-label` on icon-only buttons (existing pattern) and label inputs.
Vietnamese UI copy is the norm.

## Naming

Descriptive component names (`VenueDetail`, `VenueMap`, `VenueRailCard`). Never
`Component1` / `DataCard2`.

## Future Direction (approved, do not do preemptively)

- Split `components/DiscoveryExperience.tsx` (~1350 lines) into focused pieces:
  `VenueMap`, `VenueDetail`, `VenueRailCard`, filter controls, map-style helpers,
  and format utilities, under `components/`.
- Introduce React Query / React Hook Form + Zod ONLY when real
  caching/mutations/forms appear.
