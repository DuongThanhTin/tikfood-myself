# 03 — Component Catalog

> Every reusable visual unit in `apps/web`, its defining CSS class(es), variants,
> states, and consuming JSX. Sources: `globals.css` (styles) and
> `DiscoveryExperience.tsx` (usage). There are no separate component files — all
> pieces live inside `DiscoveryExperience.tsx` as local sub-components
> (`VenueRailCard`, `VenueMap`, `VenueDetail`, `VideoCard`, `Icon`).

## Buttons

All buttons share the base group `.primaryButton, .secondaryButton, .glassButton,
.iconGlassButton, .glassAction` (`:505-513`): `1px solid var(--outline)`,
`border-radius: 12px`, `font-weight: 800`.

| Component | Class | Shape / key styles | States | Used |
|-----------|-------|--------------------|--------|------|
| Primary | `.primaryButton` (`:515`) | min-h 40, transparent border, text `#4e2600`, bg `--primary`, soft amber shadow | `:disabled { cursor:wait; opacity:.7 }` | Search action, auth CTA, bottom action |
| Primary large | `.primaryButton.large` (`:524`) | full width, min-h 54, `grid-column:1/-1`, Plus Jakarta 16px | — | "Xem chi tiết & Menu" (`:1173`) |
| Secondary | `.secondaryButton` (`:532`) | min-h 40, text `--on-surface`, bg `--glass` | disabled as above | "Reset" (`:538`) |
| Glass | `.glassButton` (`:773`) | inline-flex, min-h 40, bg `--glass`, `blur(20px)`, text `--on-muted` | `.solid` → min-w 108, bg `--surface-high` | Topbar "Bộ lọc", "Gần tôi" (`:589,602`) |
| Theme toggle | `.glassButton.themeToggle` (`:784`) | glass button, bg `--glass-strong`, text `--on-surface` | swaps sun/moon glyph + label | Topbar (`:594`) |
| Icon glass | `.iconGlassButton` (`:802`) | 42×42 circle, glass, blur | — | Help button (`:605`) |
| Glass action | `.glassAction` (`:1473`) | min-h 48, glass + rgba-primary border, blur18, `text-decoration:none` | `:disabled { opacity:.58 }` | Detail actions (route/share/save) |
| Side tab | `.sideTab` (`:247`) | min-h 42, radius 8, transparent, muted, 12px/800 | `:hover` → glass; `.active` → primary + rgba-primary bg | Left-panel nav (`:433,437`) |
| Mini tab | `.miniTab` (`:714`) | 48×48, radius 14, glass | `.active` primary; `:disabled { opacity:.45; not-allowed }` | Collapsed nav (`:574-582`) |
| Section link | `.sectionHeader button` (`:578`) | borderless, primary text, 12px/800 | — | "Tất cả" (`:550`) |
| Bottom action | `.mapBottomAction button` (`:1144`) | pill, min-h 52, padding 0 28, text `#4e2600`, bg primary, weight 900 | `:hover { scale(1.03) }` | "Tìm quanh đây" (`:625`) |
| Close | `.closeButton` (`:1163`) | 42×42 circle, glass, blur, white glyph | — | Detail panel (`:634`) |
| Route clear | `.routeSummary button` (`:1107`) | borderless, blue `#2f80ff`, 12px/900 | — | "Xóa chỉ đường" (`:1025`) |
| Auth | `.authCta button` (`:698`) | min-h 46, radius 12, text `#4e2600`, bg primary | — | "Đăng nhập ngay" (`:569`) |

## Chips & tags

- **Base** `.contextChip, .tag` (`:408-423`): inline-flex, min-h 38, pill, `1px
  var(--outline)`, bg `--glass`, 12px/800, `blur(18px)`. Hover → rgba-primary border,
  `--surface-highest` bg, `translateY(-1px)`.
- **Trend variant** `.contextChip.trend` (`:432`): borderless, white text,
  gradient `--secondary → --trend-fire`, fire-tinted shadow; carries a leading
  `fire` glyph (`:467`).
- **Active tag** `.tag.active` (`:475`): rgba-primary border + text `--primary` +
  rgba-primary bg.
- Chip container rows share `display:flex; flex-wrap:wrap; gap:8px`
  (`.contextChips, .tagBar, .actions, .filterActions, .categoryList, .dishList`, `:397-406`).
- Data: `contextChips` and `tagOptions` arrays (`DiscoveryExperience.tsx:88-103`) —
  Vietnamese labels, `kind` of `trend | query | price` drives `applyChip` behavior.

## Form controls

| Control | Class | Styles |
|---------|-------|--------|
| Field wrapper (label) | `.field` (`:455`) | grid gap 7, muted, 12px/800 |
| Select | `.field select`, `.sortSelect` (`:463`) | width 100%, min-h 42, radius 10, `--outline` border, bg `--surface-low`; `.sortSelect` max-w 150 |
| Filter grid | `.filterGrid` (`:449`) | 2-col grid gap 10 (→ 1 col under 640px) |
| Toggle | `.toggle` (`:486`) | flex gap 8, 13px/800; native checkbox 18×18, `accent-color: var(--primary)` |
| Hero search | `.heroSearch` (`:325`) | relative grid; blurred gradient glow via `::before` (opacity .32 → .62 on focus-within); `input` min-h 64, `2px` rgba-primary border, radius 16, padding `0 54px`, bg `--surface-highest`, 16px/700; leading `search` glyph (primary, 24px), trailing `spark` glyph |

## Badges

Shared group `.guestBadge, .venueMetaRow em, .detailBadges span` (`:211-222`): pill,
rgba-secondary border + bg, text `--secondary`, 10px/800, uppercase, `letter-spacing:0.08em`.

- `.guestBadge` — "Guest View" (`:420`).
- `.venueMetaRow em` — the per-card badge (e.g. "TOP 5 TRENDING"), max-w 126 ellipsis.
- `.detailBadges span` — hero overlay badges, white text on `rgba(255,255,255,0.12)` +
  `blur(16px)`; **first child** overrides to solid `--trend-fire` (`:1219`). JSX hardcodes
  "HOT" + "TRENDING" (`:1062-1063`).

## Cards

### Venue rail card — `.venueRailCard` (`:591-693` CSS / `662-694` JSX)
2-col grid (96px thumb + summary), borderless button, radius 14. Hover/selected turns
the title `--primary` and zooms the thumb (`scale(1.08)`). Contains:
- `.venueThumb` (96×96, radius 12, `--surface-high`, `object-fit:cover`; 82px under 640px).
- `.venueSummary` — `strong` (15px/800, ellipsis) · `small` (cuisine · district) ·
  `.venueMetaRow` (badge `em` + star rating span).

### Video card — `.videoCard` (`:1424-1461`)
`aspect-ratio: 9/16`, radius 14, `--surface-high`; image zoom on hover; dark gradient
`::after`; caption `span` (bottom-left, white 11px/900) with play glyph + platform + views.

### Trend-score card — `.trendScoreCard` (`:1297-1327`)
glass + rgba-primary border, radius 18, blur18. Header row (label + `strong` percentage
in `--primary`, Plus Jakarta 24px) over `.trendTrack` (8px pill track, `--surface-highest`)
whose fill `span` is a `--secondary → --trend-fire` gradient at inline `width:${percent}%`
(`:1109-1110`). Body `p` is the AI summary (italic, muted).

### Auth CTA — `.authCta` (`:679-705`)
Centered glass card, radius 18, rgba-primary border, blur18, vertical margins 32px.

### Map overlays (also card-like)
- `.mapCountLegend` (`:1043`) — glass-strong, radius 16, count `strong` in primary
  (Plus Jakarta 24px) + muted label.
- `.routeSummary` (`:1075`) — glass-strong, radius 18, **blue** (`#2f80ff`) border/numerals,
  distance + duration + clear button.

## Detail-panel building blocks

| Block | Class | Notes |
|-------|-------|-------|
| Hero | `.detailHero` (`:1183`) | 320px (240px under 640px), `object-fit:cover`, dark gradient `::after` |
| Body | `.detailBody` (`:1224`) | grid, gap 26, padding `28 28 32` |
| Title row | `.detailTitleRow` (`:1230`) | grid `1fr auto` (→ 1 col under 640px); h2 30px Plus Jakarta |
| Compact actions | `.compactDetailActions` (`:1251`) | 3-col grid of `.glassAction`, min-h 36, 11px/**850** |
| Rating | `.ratingBlock` (`:1276`) | right-aligned; `span` secondary, Plus Jakarta 22px + star; `small` review count |
| Category / dish chips | `.category`, `.dish` (`:1363-1377`) | pills; category = primary text on rgba-primary; dish = on-surface on glass |
| Dish rows / hours | `.dishDetailItem`, `.openingRow` (`:1386-1415`) | grid `1fr auto`, outline border, radius 12, glass; opening `strong` in primary |
| Status line | `.detailStatus` (`:1356`) | muted 12px/800, negative `margin:-8px 0` |
| Route hint | `.routeHint` (`:1484`) | muted 12px/700; `.error` variant → `--danger` |

## Map markers (two systems)

Base `.mapMarker, .fallbackMarker` (`:874-882`): 2px primary border, rgba-primary bg,
fire-tinted shadow, 12px/900.

| Marker | Class | Form |
|--------|-------|------|
| Cluster | `.mapMarker.cluster` (`:892`) | 34×34 pill, shows a count `strong` |
| Named | `.mapMarker.named` (`:897`) | borderless bubble + pin, max-w 168 |
| Bubble | `.markerBubble` (`:908`) | pill, rgba-primary border, `rgba(20,20,20,.78)` bg, blur16; inner `em` icon chip + ellipsised name |
| Pin | `.markerPin` (`:950`) | 17×17 circle, 3px white border, `--trend-fire` bg, double shadow, white `::after` dot |
| User location | `.userLocationMarker`, `.fallbackUserLocation` (`:995`) | 24×24, 3px white border, **blue `#2f80ff`** bg, double shadow, white dot |
| Fallback | `.fallbackMarker` (`:975`) | absolute, 36×36, `translate(-50%,-50%)`, hover `scale(1.08)` — used pre-map-ready, positioned by `toMapX/toMapY` |
| Selected | `.mapMarker.selected`, `.fallbackMarker.selected` (`:1022-1134`) | expands (min-w 138), radius 14, dark bg, shows a 36×36 circular thumb + short name |

> The MapLibre markers are built as raw DOM in `DiscoveryExperience.tsx:824-851` via
> `element.innerHTML`; the fallback layer is JSX (`:975-1014`). See register **R8**.

## The `Icon` component

`Icon` (`DiscoveryExperience.tsx:1344-1350`) renders a unicode glyph inside
`.uiIcon` (`:55-63`: 1em inline-grid box, 20px, weight 900, `aria-hidden="true"`).
Optional `className` (e.g. `.heroSearchSpark`, `:391`) tweaks position/color. Full
glyph mapping is in [`04-patterns.md`](04-patterns.md#icon-glyph-system).

## Orphaned component

`components/VenueList.tsx` defines an alternate venue-card layout (`.venueCard`,
`.venueGrid`, `.venueName`, `.meta`, `.score`, …) but **none of those classes exist in
`globals.css`** and the component is imported nowhere. It is dead code — see register
**R1**. It is intentionally *not* documented as part of the live system.
