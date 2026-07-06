# TikFood Web — Design System (as-built)

> **Reverse-engineered documentation of the *current* `apps/web` design system.**
> Every value here is traceable to source. Where the code is inconsistent, the
> inconsistency is documented, not fixed. No application code was changed to produce
> this set. See [`ROADMAP.md`](ROADMAP.md) for scope and method.

## What the system actually is

TikFood's frontend has **no CSS framework and no design-system library**. The entire
visual system is:

- **One stylesheet** — `apps/web/app/globals.css` (1,606 lines) holds every token and
  every styled class.
- **CSS custom properties** on `:root` (dark) and `.appShell[data-theme="light"]`
  (light) for color/theming.
- **Two Google Fonts** imported via CSS `@import`: *Be Vietnam Pro* (body) and
  *Plus Jakarta Sans* (display).
- **One large client component** — `components/DiscoveryExperience.tsx` — that renders
  the whole three-column discovery experience and owns the MapLibre map, the icon
  glyph system, and mock presentation data.
- **Plain CSS class names**, BEM-ish but not strict (`.venueRailCard`,
  `.venueRailCard.selected`), applied via `className` string joins in TSX.

The aesthetic is a **dark-first, glassmorphism** look: warm amber/orange primary on
near-black surfaces, translucent `backdrop-filter: blur()` panels, pill-shaped chips,
and large soft shadows. A light theme (warm cream) is a full token override scoped to
the app shell.

## How to read this set

| Doc | Read it for |
|-----|-------------|
| [`01-foundations.md`](01-foundations.md) | Tokens: color, theming, typography, spacing, radii, shadows, borders, glass, z-index, motion |
| [`02-layout.md`](02-layout.md) | The `.appShell` grid, the three panels, scroll model, responsive breakpoints |
| [`03-components.md`](03-components.md) | Catalog of every reusable visual unit, its classes, variants, and states |
| [`04-patterns.md`](04-patterns.md) | Page/navigation/map patterns, the icon glyph table, UI states, copy & a11y conventions |
| [`05-inconsistencies.md`](05-inconsistencies.md) | The register of divergences, dead styles, off-scale values, and mocked data (R1–R…) |
| [`components/`](components/README.md) | Per-component reference docs (one per component, fixed template) — the detailed companion to `03-components.md` |
| [`improvement-proposals.md`](improvement-proposals.md) | Review of this doc set + the system: duplicated/missing docs, missing components, and spacing/type/color/responsive inconsistencies — **proposals only**, nothing applied |

## Source-of-truth statement

When these docs and the code disagree, **the code wins** — flag the mismatch here
rather than trusting the doc. Authoritative sources, in order: `app/globals.css`,
`components/DiscoveryExperience.tsx`, `app/layout.tsx` + `app/page.tsx`, `lib/api.ts`,
`components/VenueList.tsx` (legacy — see [R1](05-inconsistencies.md)).

## Related

- Frontend engineering rules: [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md)
- Frontend architecture standard: [`docs/standards/frontend-architecture.md`](../../standards/frontend-architecture.md)
- Repository map: [`docs/REPOSITORY-MAP.md`](../../REPOSITORY-MAP.md)
