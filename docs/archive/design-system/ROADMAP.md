# Design System Audit — Roadmap

> **⚠️ ARCHIVED — superseded by [`docs/design/`](../../design/00-overview.md).** Frozen
> snapshot of the original audit plan; not maintained. Use `docs/design/` for the current
> design system.

> **What this is.** A plan for documenting the *current* `apps/web` design system by
> reverse-engineering the shipped code. This roadmap is produced **before** the
> documentation itself (deliverable order agreed with the maintainer).
>
> **What this is not.** Not a redesign, not a proposal, not a wishlist. We describe
> what exists. Where the code is inconsistent, we **document the inconsistency** — we
> do not fix it here. No application code changes are made by this audit.

## 1. Role & rules of engagement

- **Role:** Senior Design System Architect performing a read-only reverse-engineering audit.
- **Read-only:** no changes to `apps/web/**` source. The only writes are these docs under `docs/design-system/`.
- **Don't invent:** every documented value, class, or pattern must be traceable to a line in the source. No aspirational tokens, no "should be".
- **Inconsistencies are findings, not bugs to fix:** they go in the register (§5), each with a source reference.

## 2. Sources of truth (ranked)

| Rank | Source | What it authoritatively defines |
|------|--------|--------------------------------|
| 1 | `apps/web/app/globals.css` (1,606 lines) | All design tokens (CSS custom properties), every styled class, dark/light theming, responsive rules, motion |
| 2 | `apps/web/components/DiscoveryExperience.tsx` (~1,350 lines) | The real component tree, icon system, page/navigation patterns, inline styles, mock media, map behavior |
| 3 | `apps/web/app/layout.tsx`, `app/page.tsx` | Root shell, font imports (via CSS `@import`), metadata, server→client data handoff |
| 4 | `apps/web/lib/api.ts` | Data shapes (`Venue`, `SocialVideo`, `VenueDish`, `OpeningHour`) that drive what the UI renders |
| 5 | `apps/web/components/VenueList.tsx` | **Legacy/orphaned** component — see register item R1 |

When these disagree, the **code wins** and the disagreement is logged in §5.

## 3. Deliverable structure (planned files)

All under `docs/design-system/`:

| File | Contents |
|------|----------|
| `README.md` | Index + how to read this set + one-paragraph system summary + source-of-truth statement |
| `01-foundations.md` | **Tokens.** Color palette (dark + light), theming mechanism (`data-theme` on `.appShell`), typography (2 families, weights, the h1/h2/h3 scale), the (implicit) spacing values, radii, shadows, borders/outlines, glass/blur, z-index layering, motion/transition durations |
| `02-layout.md` | **Layout.** The `.appShell` CSS-grid shell and its 4 column states (default / left-collapsed / detail-open / both), the three panels (left/map/right), scroll model (`body{overflow:hidden}`), responsive breakpoints (1180px, 640px) and how each collapses the grid |
| `03-components.md` | **Component catalog.** Every reusable visual unit with its class(es), variants, states, and where it's used: buttons (primary/secondary/glass/icon-glass/glass-action), chips & tags, form fields & selects, toggle, venue rail card, section header, badges, detail hero, trend-score card, dish/opening rows, video card, map markers (bubble/pin/cluster/user/fallback), route summary, count legend, close button, auth CTA |
| `04-patterns.md` | **Patterns.** Page layout pattern (three-column discovery), navigation (side tabs / collapsed mini-tabs / map topbar), the `Icon` glyph system + full glyph table, map interaction patterns (marker modes, flyTo, routing overlay, dark/light basemap tuning), UI states (loading/error/empty/disabled), bilingual copy convention (Vietnamese UI + English technical), accessibility patterns (`aria-label` on icon buttons, `aria-hidden` on decoration) |
| `05-inconsistencies.md` | **Register.** Every divergence, dead style, off-scale value, hardcoded color, and duplicated helper — each with a source reference and a note on impact. Fixes are *out of scope*; this is the catalog. |

## 4. Method (per file)

1. Extract raw values/classes directly from source (already gathered — see §2).
2. Group by concern; present as tables where values are enumerable (tokens, glyphs, breakpoints).
3. For each component/pattern, cite the defining CSS lines and the consuming JSX lines.
4. Cross-check JSX class usage against CSS definitions; anything used-but-undefined or defined-but-unused → §5.
5. No interpretation beyond what the code states; label anything inferred as "observed behavior".

## 5. Inconsistency register — preview (confirmed from source)

These are already verified against the code and will be expanded in `05-inconsistencies.md`.
IDs are stable so later docs can reference them.

| ID | Finding | Source |
|----|---------|--------|
| R1 | **`VenueList.tsx` is orphaned.** It is imported nowhere (`page.tsx` renders `DiscoveryExperience`), and every class it uses — `.venues`, `.sectionTitle`, `.resultCount`, `.emptyState`, `.venueGrid`, `.venueCard`, `.venueHeader`, `.venueName`, `.meta`, `.score` — is **undefined** in `globals.css`. It would render unstyled. | `components/VenueList.tsx`; absent from `globals.css` |
| R2 | **`.mapCountLegendLive` used but undefined.** JSX adds this class for the live cluster legend but only base `.mapCountLegend` exists in CSS. | `DiscoveryExperience.tsx:1016`; `globals.css` |
| R3 | **Duplicate, divergent format helpers.** `formatPrice`/`formatDistance` exist in *both* `VenueList.tsx` and `DiscoveryExperience.tsx` with different rounding rules and units (VenueList: `k` only; Discovery: `k`/`tr`, km with variable precision). | `VenueList.tsx:53-65`; `DiscoveryExperience.tsx:1295-1320` |
| R4 | **Font weights used that aren't loaded.** CSS uses `font-weight: 900` (icons, markers, section headers, route buttons) and `850` (`.compactDetailActions .glassAction`), but the imported faces are Be Vietnam Pro `400;500;700;800` and Plus Jakarta Sans `600;700;800`. 850/900 fall back to 800 (or synthetic bold). | `globals.css:1` (import); `:55,882,584,1114,1264` |
| R5 | **Blue accent has no token.** All warm colors are CSS variables, but the routing/user-location blue `#2f80ff` (and `#2f271f`, button-text `#4e2600`, `#2f1500`) is hardcoded and repeated across CSS + TSX. | `globals.css:1003,1094,1111,948`; `DiscoveryExperience.tsx:948` |
| R6 | **No spacing / radius scale tokens.** Spacing uses raw px (6,7,8,10,12,14,16,18,22,24,26,28,30,32) and radii (8,10,12,14,16,18,999) with no `--space-*` / `--radius-*` variables — magnitudes are chosen ad hoc per rule. | `globals.css` throughout |
| R7 | **Mocked presentation data.** Ratings, review counts, "HOT"/"TRENDING" badges, cuisine labels, and all imagery come from a hardcoded `mediaByVenue` map keyed to only 2 venue IDs; every other venue reuses venue #1's media via `defaultMedia`. None of this is API-driven. | `DiscoveryExperience.tsx:105-132,1209-1211` |
| R8 | **Two parallel marker systems.** MapLibre DOM markers (`.mapMarker*`) and a CSS fallback layer (`.fallbackMarker*`, positioned by the linear `toMapX/toMapY` approximation) have overlapping but non-identical style rules. | `globals.css:874-1134`; `DiscoveryExperience.tsx:975-1014,1336-1342` |
| R9 | **Theme scoped to `.appShell`.** `:root` is `color-scheme: dark` and light-mode tokens only apply under `.appShell[data-theme="light"]`; anything rendered outside the shell stays dark. | `globals.css:3-22,72-91` |

(Additional smaller items — e.g. redundant `letter-spacing:0` resets, `.detailStatus{margin:-8px 0}` negative margin, `!important` uses on `.sidebarToggle`/`.heroSearchSpark` — will be catalogued in the full register.)

## 6. Sequencing & checkpoints

1. ✅ **Fact-gathering** — all five sources read in full.
2. ✅ **Roadmap** — this document.
3. ⏳ **Checkpoint:** confirm file placement (`docs/design-system/`) and structure (§3) before writing the set.
4. ⏳ **Write docs** in order: `README` → `01-foundations` → `02-layout` → `03-components` → `04-patterns` → `05-inconsistencies`.
5. ⏳ **Cross-check pass** — verify every cited line number resolves; move any newly found divergence into R-register.
6. ⏳ Update `docs/REPOSITORY-MAP.md` with the new `docs/design-system/` entry.

## 7. Explicitly out of scope

- Any edit to `apps/web` source (components, CSS, helpers) — including "obvious" fixes like deleting `VenueList.tsx` or tokenizing the blue.
- Proposing a target/ideal design system or refactor plan.
- Visual/screenshot capture or running the app.
- Accessibility remediation (we document current a11y patterns only).
