# Design System — Overview

> **Canonical home** for the TikFood `apps/web` design system. This tree
> (`docs/design/`) supersedes the earlier `docs/design-system/` audit set, now archived
> at [`docs/archive/design-system/`](../archive/design-system/).

## What this is

A reverse-engineered, as-built design system for the TikFood discovery web app
(Next.js App Router, plain CSS in `app/globals.css`, MapLibre GL, no CSS framework).
It describes **what exists in the code today** and, where a requested artifact does not
yet exist, marks it clearly as **Proposed**.

## Status & conventions

- **Source of truth:** the code. When a doc and `apps/web` disagree, the code wins —
  flag the mismatch (see [`09-accessibility.md`](09-accessibility.md) and the
  consistency findings referenced throughout).
- **As-built vs Proposed:** every doc states its status.
  - *As-built* — documents shipped code (line-referenced).
  - **Proposed** — the component/page/pattern does **not exist yet**; the doc is a
    forward spec grounded in current tokens, banner-labeled `> **Status: Proposed —
    not implemented.**`. Nothing about existing code is invented.
- **No duplication:** tokens and foundations live in `00–09`; components, patterns, and
  pages *reference* them rather than restating values.
- **Bilingual copy:** UI strings are Vietnamese; identifiers/technical terms are English
  (project convention, `apps/web/CLAUDE.md`).

## Structure

```
docs/design/
├─ 00-overview.md            ← you are here
├─ 01-design-principles.md   observed principles that guide the UI
├─ 02-brand.md               name, wordmark, voice, identity
├─ 03-color-system.md        color tokens + theming (canonical)
├─ 04-spacing-system.md      spacing / radii (as-built: no scale yet)
├─ 05-grid-system.md         the app-shell grid, panels, breakpoints
├─ 06-typography.md          families, scale, weights
├─ 07-icons.md               the unicode-glyph icon system
├─ 08-motion.md              transitions + map camera motion
├─ 09-accessibility.md       a11y patterns and gaps
├─ components/               button, card, input, map (+ Proposed: modal, table)
├─ patterns/                 search, filter, empty-state, loading, error-state (+ Proposed: pagination)
├─ pages/                    home (+ Proposed: restaurant, dish, profile, search)
└─ playbooks/                how-to guides (Proposed — follow-up pass)
```

## Reading order

1. This overview → [`01-design-principles.md`](01-design-principles.md) →
   [`02-brand.md`](02-brand.md) for the *why*.
2. Foundations `03–08` for the *what* (tokens).
3. `09-accessibility.md` for cross-cutting a11y.
4. `components/`, `patterns/`, `pages/` for applied guidance.
5. `playbooks/` when building something new.

## Where the code lives

- Styles: `apps/web/app/globals.css` (single stylesheet, all tokens + classes).
- UI: `apps/web/components/DiscoveryExperience.tsx` (the whole experience, with local
  sub-components `VenueRailCard`, `VenueMap`, `VenueDetail`, `VideoCard`, `Icon`).
- Root/fonts/metadata: `apps/web/app/layout.tsx`; data layer: `apps/web/lib/api.ts`.
- `apps/web/components/VenueList.tsx` is **orphaned** (imported nowhere, classes
  undefined) and is not treated as part of the live system.

## Migration note

Migration is **complete**. The original reverse-engineering audit (per-component docs,
the inconsistency register, and the consolidated **improvement proposals**) is preserved
as a frozen snapshot at [`docs/archive/design-system/`](../archive/design-system/). The
foundations docs here (`03–09`) embed the key consistency findings inline and link into the
archived proposals for detail.
The archived docs are historical and not maintained (their code line-references reflect the
pre-migration layout).

## Related

- [Design principles](01-design-principles.md) · [Brand](02-brand.md)
- Product framing: [`CLAUDE.md`](../../CLAUDE.md), [`docs/tikfood/`](../tikfood/)

## References

- `apps/web/CLAUDE.md` — frontend engineering rules.
- `docs/REPOSITORY-MAP.md` — where everything lives.

## Future improvements

- Add rendered visual examples per component (currently text-only).
- Flip Proposed docs (modal, table, pagination, restaurant-route, dish, profile, search)
  to As-built as they ship, dropping their Proposed banners.
