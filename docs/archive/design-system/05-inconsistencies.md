# 05 — Inconsistency Register

> A catalog of divergences, dead styles, off-scale values, and mocked data found while
> reverse-engineering the current design system. **These are findings, not a fix list.**
> Remediation is explicitly out of scope for this audit (see [`ROADMAP.md`](ROADMAP.md#7-explicitly-out-of-scope)).
> IDs are stable so other docs can cite them.

## Primary findings

### R1 — `VenueList.tsx` is orphaned; its classes don't exist
`components/VenueList.tsx` is imported **nowhere** (`app/page.tsx` renders
`DiscoveryExperience`). Every class it applies — `.venues`, `.sectionTitle`,
`.resultCount`, `.emptyState`, `.venueGrid`, `.venueCard`, `.venueHeader`,
`.venueName`, `.meta`, `.score`, `.categoryList`/`.dishList`/`.category`/`.dish`
(the last four *do* exist) — is largely **undefined** in `globals.css`. If mounted it
would render mostly unstyled.
**Impact:** dead code that misrepresents the component inventory; a maintainer may
assume a working list view exists.
**Source:** `components/VenueList.tsx`; absent from `globals.css`.

### R2 — `.mapCountLegendLive` used but undefined
The live cluster legend adds class `mapCountLegend mapCountLegendLive`, but only base
`.mapCountLegend` is defined; `.mapCountLegendLive` has no rule.
**Impact:** the "live" legend is styled identically to the fallback one; the modifier
is a no-op.
**Source:** `DiscoveryExperience.tsx:1016`; `globals.css:1043` (no `Live` rule).

### R3 — Duplicate, divergent format helpers
`formatPrice` and `formatDistance` are defined **twice** with different logic:
- `VenueList.tsx:53-65` — price divides by 1000 above 1000 (`k` only); distance km at 1 dp.
- `DiscoveryExperience.tsx:1295-1320` — price uses `k`/`tr` with variable precision and
  an "Đang cập nhật" zero-case; distance km precision varies by magnitude.

**Impact:** two sources of truth for the same formatting; divergent output if both
were ever used. (Mitigated in practice by R1 — the VenueList copy is dead.)

### R4 — Font weights used that aren't loaded
CSS uses `font-weight: 900` widely (`.uiIcon`, markers, `.markerBubble>span`,
`.sectionHeader` context, `.routeSummary button`, `.videoCard>span`,
`.mapBottomAction button`) and `font-weight: 850` once
(`.compactDetailActions .glassAction`), but the `@import` loads only Be Vietnam Pro
`400;500;700;800` and Plus Jakarta Sans `600;700;800`.
**Impact:** 850/900 fall back to the nearest loaded weight (800) or synthetic bold;
intended "heaviest" emphasis isn't truly rendered. `850` is not a valid CSS keyword
step and reads as an arbitrary numeric weight.
**Source:** `globals.css:1` (import) vs `:61, 881, 945, 1114, 1155, 1460` (900) and `:1264` (850).

### R5 — Blue accent has no token; warm colors re-expressed as literals
Every warm color is a CSS variable, but the routing/user-location **blue `#2f80ff`**
(and companions `#2f271f`, button text `#4e2600`, brand `#2f1500`, grid `#080808`) is
hardcoded and repeated across CSS **and** TSX. Conversely, `--primary`/`--secondary`/
`--trend-fire` are frequently re-written as literal `rgba(255,183,127,…)` /
`rgba(254,183,0,…)` / `rgba(255,77,0,…)` instead of referencing the token. The map
palette in `applyMapLayerPreferences` is an entirely separate hardcoded color set.
**Impact:** the blue accent can't be themed (it stays blue in light mode); token
values and their rgba duplicates can drift independently.
**Source:** `globals.css:1003, 1094, 1111, 269, 436, 145-216`; `DiscoveryExperience.tsx:948` (route line), `:934` (casing swaps `#050505`/`#ffffff`).

### R6 — No spacing or radius scale tokens
Spacing (`2…32px`) and radii (`8,10,12,14,16,18,999px`) are raw literals chosen
per-rule; no `--space-*` / `--radius-*` variables exist. Only color and one `--shadow`
are tokenized.
**Impact:** no enforced rhythm; values are internally consistent by convention only,
and a change means editing many rules.
**Source:** `globals.css` throughout.

### R7 — Presentation data is mocked, not API-driven
`mediaByVenue` (`DiscoveryExperience.tsx:105-132`) hardcodes imagery, `cuisine`,
`rating`, `reviews`, and `badge` for **only two venue IDs**. `getVenueMedia`
(`:1209-1211`) falls every other venue back to `defaultMedia` (venue #1's assets).
The hero badges "HOT"/"TRENDING" (`:1062-1063`) are literal strings.
**Impact:** ratings, review counts, badges, and photos shown in the UI are **not
real** and don't reflect the `Venue` API payload; all non-seed venues share one set of
photos. This is UI-critical to flag.
**Source:** `DiscoveryExperience.tsx:105-132, 1062-1063, 1209-1211`.

### R8 — Two parallel marker systems
MapLibre DOM markers (`.mapMarker*`, built via `element.innerHTML`,
`DiscoveryExperience.tsx:824-851`) and a CSS/JSX fallback layer (`.fallbackMarker*`,
positioned by the linear approximation `toMapX/toMapY`, `:975-1014, 1336-1342`) have
overlapping but non-identical style rules and duplicated markup (including the count
legend, rendered in both branches).
**Impact:** two code paths and two style surfaces to keep in sync; the fallback
projection is a rough clamp, not a real map projection.
**Source:** `globals.css:874-1134`; `DiscoveryExperience.tsx:824-851, 975-1014`.

### R9 — Theme is scoped to `.appShell` only
`:root` is `color-scheme: dark` and light tokens live under
`.appShell[data-theme="light"]`. Anything rendered outside the shell (the `<body>`
itself, portals, error boundaries) keeps dark tokens regardless of theme. Theme state
is also **not persisted** (resets to dark on reload).
**Source:** `globals.css:3-22, 33-39, 72-91`.

## Secondary findings

| ID | Finding | Source |
|----|---------|--------|
| R10 | **No-op controls.** "Bộ lọc", "Chia sẻ", "Lưu", the help button, and "Đăng nhập ngay" render as full buttons but have no `onClick` and aren't `disabled` — they look interactive but do nothing. | `DiscoveryExperience.tsx:589,605,569,1077-1084` |
| R11 | **Declared-but-unused tokens.** `--surface-mid` is never referenced by any class; `--primary-strong` is declared but gradients use the literal `#ff8a00` (or `--trend-fire`) instead. | `globals.css:8,17,336` |
| R12 | **`!important` used.** `.sidebarToggle` collapse transform and `.heroSearchSpark` (`left:auto`, `color`) rely on `!important`, signalling specificity conflicts rather than clean cascade. | `globals.css:165,393-394` |
| R13 | **Clustering is a stub.** `cluster` marker mode exists (zoom < 12.2) but `getVenueClusterCount` always returns `1`, so "clusters" are single venues. | `DiscoveryExperience.tsx:218-220,1253-1257` |
| R14 | **External routing call to a public demo server.** `fetchRoute` hits `router.project-osrm.org` directly — an un-configurable, rate-limited, third-party endpoint. (Also an external network call, a class flagged in root `CLAUDE.md`.) | `DiscoveryExperience.tsx:1213-1218` |
| R15 | **Raw `fetch` in a component.** `fetchRoute` bypasses the "all data through `lib/api.ts`" rule from `apps/web/CLAUDE.md`. | `DiscoveryExperience.tsx:1216` |
| R16 | **Redundant resets.** `letter-spacing: 0` is set explicitly on `h1/h2/.brandName` (the default is already `normal`≈0). | `globals.css:196,296` |
| R17 | **Negative margin.** `.detailStatus { margin: -8px 0 }` pulls the loading line into adjacent spacing rather than adjusting the gap. | `globals.css:1357` |
| R18 | **BEM-ish but inconsistent naming.** Mix of camelCase block+modifier (`.venueRailCard.selected`), space-modifier (`.contextChip trend`), and generic element selectors (`.trendScoreCard span`, `.sectionHeader h2`) that couple styles to tag structure. | `globals.css` throughout |
| R19 | **Two "count legend" renders.** The `<venues in area>` legend is duplicated in the fallback branch (`:1010`) and the live branch (`:1016`) with slightly different classes and identical copy. | `DiscoveryExperience.tsx:1010-1020` |

## Notes

- **Not documented as live:** `VenueList.tsx` (R1) is excluded from
  [`03-components.md`](03-components.md) except as an orphan note.
- **Cross-references:** foundations doc flags R4/R5/R6/R9 inline; components doc flags
  R1/R8; patterns doc flags R7/R10/R13/R14/R15.
- This register reflects the repository state read during the audit. If code changes,
  re-verify line references before acting on any entry.
