# 01 — Foundations (Tokens)

> Source of truth: `apps/web/app/globals.css`. Line references are to that file
> unless noted. Only the blue accent has no token; all other statements below are
> literal values from the stylesheet.

## Color tokens

Colors are CSS custom properties. Dark is the default (`:root`); light is a full
override scoped to `.appShell[data-theme="light"]`.

| Token | Dark (`:root`, `:3-22`) | Light (`.appShell[data-theme="light"]`, `:72-91`) | Role (observed use) |
|-------|-------------------------|---------------------------------------------------|---------------------|
| `--background` | `#050505` | `#eee7dc` | Page / map-stage background |
| `--surface` | `#131313` | `#fbf3e8` | Left/right panel background |
| `--surface-low` | `#1c1b1b` | `#eadfce` | Select inputs |
| `--surface-mid` | `#201f1f` | `#e2d4c1` | (declared; no direct class use found) |
| `--surface-high` | `#2a2a2a` | `#fff7ea` | Thumbnails, hero placeholders, glass-solid button |
| `--surface-highest` | `#353534` | `#eadac6` | Search input bg, chip hover, trend track, scrollbar thumb |
| `--on-surface` | `#e5e2e1` | `#251b14` | Primary text |
| `--on-muted` | `#ddc1ae` | `#6b5746` | Secondary/muted text |
| `--outline` | `rgba(255,255,255,0.12)` | `rgba(74,55,40,0.22)` | 1px borders, dividers |
| `--glass` | `rgba(255,255,255,0.08)` | `rgba(250,239,224,0.78)` | Translucent surfaces |
| `--glass-strong` | `rgba(255,255,255,0.12)` | `rgba(248,233,214,0.92)` | Stronger translucent surfaces (legends, theme toggle) |
| `--primary` | `#ffb77f` (amber) | `#8f4700` | Brand, active states, primary buttons |
| `--primary-strong` | `#ff8a00` | `#d06b00` | (declared; gradients use literal `#ff8a00` — see note) |
| `--secondary` | `#feb700` (gold) | `#a26300` | Ratings, badges, trend meta |
| `--trend-fire` | `#ff4d00` (red-orange) | `#d94f00` | Trend gradients, marker pin, hot badge |
| `--danger` | `#ffb4ab` | `#922922` | Error text |
| `--shadow` | `0 24px 70px rgba(0,0,0,0.45)` | `0 24px 70px rgba(74,52,32,0.2)` | Shared large soft shadow |

### Colors that are **not** tokenized (used as literals)

- **Blue accent `#2f80ff`** — user-location marker, route line, route summary text
  (`globals.css:1003, 1094, 1111`; and `DiscoveryExperience.tsx:948`). No CSS
  variable exists for it, unlike every warm color. See register **R5**.
- `#4e2600` — text color on primary buttons (`:519, 702, 1152`).
- `#2f1500` — collapsed brand `::first-letter` (`:209`).
- `#080808` — fallback-grid base (`:829`).
- `#ffffff` / `white` — marker text, pin borders, close button, video captions.
- Primary/secondary/fire re-expressed as literal `rgba(...)` many times
  (`rgba(255,183,127,…)`, `rgba(254,183,0,…)`, `rgba(255,77,0,…)`) instead of
  referencing the token — e.g. tab active `rgba(255,183,127,0.1)` (`:269`).
- **Map basemap colors** — a large palette of road/water/park hex values applied
  imperatively in `DiscoveryExperience.tsx:145-216` (`applyMapLayerPreferences`),
  separate from the CSS token system entirely.

### Theming mechanism

- `<main className="appShell" data-theme={theme}>` (`DiscoveryExperience.tsx:398-406`).
- `theme` is React state, default `"dark"`, toggled by the topbar button
  (`:227, 597`). No persistence (not stored to `localStorage`/cookie).
- Light mode = re-declaring the same token names under
  `.appShell[data-theme="light"]`, plus ~15 component-level overrides for glass
  surfaces, hero search, markers, and map-canvas filters (`:72-91, 142-144, 365-381,
  789-800, 841-853, 863-867, 922-927, 969-973, 1038-1041`).
- `:root` sets `color-scheme: dark`; light sets `color-scheme: light` on the shell.
  Anything rendered **outside** `.appShell` stays dark — see register **R9**.

## Typography

### Families (imported `:1` via Google Fonts `@import`)

| Family | Weights loaded | Used for |
|--------|----------------|----------|
| **Be Vietnam Pro** | 400, 500, 700, 800 | Body/default (`body`, `:38`), `h1 span`, `.sectionHeader h2` |
| **Plus Jakarta Sans** | 600, 700, 800 | Display: `h1`, `h2`, `.brandName`, `.detailTitleRow h2`, rating/legend/route/trend numerals, `.primaryButton.large` |

> Fallback stack for body: `"Be Vietnam Pro", Arial, Helvetica, sans-serif`.
> Weights **850 and 900 are used in CSS but not loaded** — see register **R4**.

### Type scale (as-built, no size tokens)

| Element / class | Size | Weight | Family | Notes |
|-----------------|------|--------|--------|-------|
| `h1` (`:298`) | 28px | 800 | Plus Jakarta | `line-height:1.18`; 25px under 640px (`:1563`) |
| `h1 span` (`:305`) | 18px | 400 | Be Vietnam | Sub-headline, muted, block |
| `h2` (`:314`) | 16px | 800 | Plus Jakarta | Base heading |
| `.sectionHeader h2` (`:570`) | 12px | (inherits) | Be Vietnam | Uppercase, `letter-spacing:0.08em`, muted |
| `.detailTitleRow h2` (`:1237`) | 30px | (800) | Plus Jakarta | `line-height:1.14`; venue name |
| `h3` (`:319`) | 15px | 800 | (inherits Be Vietnam) | Section subheads in detail |
| `.brandName` (`:190`) | 24px | 800 | Plus Jakarta | Primary color |
| Body / inputs | 16px | 700 | Be Vietnam | `.heroSearch input` (`:357`) |
| Chips/tags/fields | 12px | 800 | Be Vietnam | pervasive control label size |
| Badges (`:213`) | 10px | 800 | — | uppercase, `letter-spacing:0.08em` |
| Numerals (rating/legend/route/trend) | 22–24px | 800 | Plus Jakarta | e.g. `.ratingBlock>span` 22px, `.mapCountLegend strong` 24px |

- **Font sizes present:** 10, 11, 12, 13, 14, 15, 16, 18, 22, 24, 25, 28, 30 px — no
  size scale variable; each rule sets px directly.
- **Line-heights present:** 1, 1.14, 1.18, 1.35, 1.45, 1.5, 1.55.
- **Weights present in CSS:** 400, 700, 800, **850**, **900** (500/600 are loaded but
  not obviously applied in class rules).
- `letter-spacing` is only ever `0` (explicit reset on `h1/h2/.brandName`, `:196,296`)
  or `0.08em` (uppercase badges/section headers).
- `text-transform: uppercase` only on badges and `.sectionHeader h2`.
- Global resets: `h1,h2,h3,p { margin:0 }` (`:285-290`); `button,input,select { font:inherit }` (`:41-45`).

## Spacing

There is **no spacing scale / token** (register **R6**). Values are raw px chosen per
rule. Observed `gap`/`padding`/`margin` magnitudes:

`2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 18, 22, 24, 26, 28, 30, 32` px.

- Panel content gutters standardize on **`margin: 0 24px`** (`.searchFocus`,
  `.compactFilters`, `.venueRail`, `.authCta`, `:272-277`), dropping to `16px` under
  640px (`:1542-1560`).
- One negative margin exists: `.detailStatus { margin: -8px 0 }` (`:1357`).

## Radii

No radius tokens (register **R6**). Values: `8, 10, 12, 14, 16, 18` px and `999px`
(pill). `border-radius: inherit` is used on `.trendTrack span` (`:1339`) and
`.markerPin::after` (`:963`). Rough convention observed:

| Radius | Applied to |
|--------|-----------|
| 8px | side tabs, small chips grid |
| 10–12px | selects, buttons, thumbnails |
| 14px | rail cards, mini-tabs, video cards, selected marker |
| 16–18px | search input, legends, hero-search glow, trend card, auth CTA |
| 999px | all pills: chips, tags, icon buttons, markers, toggles, trend track |

## Shadows & elevation

`--shadow` (`0 24px 70px …`) is the shared token, used on the sidebar toggle, count
legend, route summary. Most other shadows are **bespoke literals** per component:

| Component | Shadow (`:line`) |
|-----------|------------------|
| `.rightPanel` | `-24px 0 80px rgba(0,0,0,0.44)` (`:139`) |
| `.heroSearch input` | `0 22px 60px rgba(0,0,0,0.36)` (`:355`) |
| `.primaryButton` | `0 12px 28px rgba(255,183,127,0.12)` (`:521`) |
| `.contextChip.trend` | `0 12px 32px rgba(255,77,0,0.2)` (`:436`) |
| `.mapMarker` | `0 10px 24px rgba(255,77,0,0.2)` (`:879`) |
| `.markerBubble` | `0 12px 28px rgba(0,0,0,0.24)` (`:917`) |
| `.markerPin` | `0 0 0 6px rgba(255,77,0,0.18), 0 10px 22px rgba(0,0,0,0.24)` (`:958`) |
| `.userLocationMarker` | `0 0 0 8px rgba(47,128,255,0.22), 0 18px 38px rgba(0,0,0,0.28)` (`:1004`) |
| `.mapBottomAction button` | `0 22px 60px rgba(255,183,127,0.18)` (`:1154`) |

Light theme swaps several of these to warm-brown-tinted shadows (`:794, 926, 972`).

## Borders

- Standard hairline: `1px solid var(--outline)`.
- Heavier borders: `2px` on `.heroSearch input` (`:350`) and base `.mapMarker` (`:876`);
  `3px solid #ffffff` on `.markerPin` and `.userLocationMarker` pins (`:955, 1001`).
- Panel dividers: `.leftPanel { border-right }`, `.rightPanel { border-left }` (`:133, 138`).

## Glass / blur

Glassmorphism is a core motif. `backdrop-filter: blur(<n>px)` values: **12, 16, 18, 20** px.
Applied to chips, glass buttons, marker bubbles, legends, route summary, auth CTA,
close button, detail badges. Backed by `--glass` / `--glass-strong` fills. The hero
search uses a separate blurred **glow pseudo-element** (`.heroSearch::before`,
`filter: blur(12px)`, `:331-345`) rather than a backdrop filter.

## Z-index layering

Observed stacking order (higher = closer to viewer):

| z-index | Layer(s) |
|---------|----------|
| `-1` | `.heroSearch::before` glow (`:333`) |
| `1` | `.mapCanvas` (`:834`) |
| `2` | `.mapShade`, `.detailBadges`, `.videoCard>span`, `.heroSearch>.uiIcon` |
| `3` | `.fallbackMarkerLayer` (`:871`) |
| `4` | `.leftPanel`, `.rightPanel`, `.fallbackUserLocation` |
| `5` | `.mapCountLegend`, `.mapBottomAction` |
| `6` | `.mapTopbar`, `.routeSummary` |
| `7` | `.closeButton` |
| `8` | `.sidebarToggle` (`:148`) |

## Motion

- **No `@keyframes`** exist — there are no looping/entrance animations. All motion is
  CSS `transition` on hover/state, plus JS map camera moves.
- **Transition durations:** `0.18s, 0.2s, 0.25s, 0.3s, 0.45s`; easing is always `ease`.
  - `0.3s` — panel width collapse, sidebar toggle transform (`:135, 161`).
  - `0.45s` — image zoom on hover (`.venueThumb img`, `.videoCard img`, `:621, 1436`).
  - `0.25s` — hero-search glow opacity (`:340`).
  - `0.18s` — marker bubble / fallback marker (`:919, 984`).
  - `0.2s` — chips, tabs, buttons, thumbnail color.
- **Hover transforms:** `translateY(-1px)` (chips, marker bubble), `scale(1.03–1.08)`
  (bottom action, thumbnails, fallback marker).
- **JS camera motion** (`DiscoveryExperience.tsx`): `flyTo` 500ms on venue select
  (`:867`), 700ms on user location (`:901`); `fitBounds` 700ms on route (`:956`).

## Iconography (summary)

Icons are **unicode glyphs**, not SVG/icon-font — mapped in `iconGlyphs`
(`DiscoveryExperience.tsx:66-86`) and rendered by the `Icon` component through
`.uiIcon` (1em box, 20px, weight 900, `aria-hidden`). Full glyph table in
[`04-patterns.md`](04-patterns.md).
