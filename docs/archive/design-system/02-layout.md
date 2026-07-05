# 02 — Layout

> Source: `apps/web/app/globals.css` and `components/DiscoveryExperience.tsx`.

## The app shell

The whole UI is one CSS-grid container, `.appShell`
(`globals.css:65-70`), rendered as `<main>` in
`DiscoveryExperience.tsx:398-406`:

```
.appShell { display: grid; grid-template-columns: 420px minmax(0, 1fr); min-height: 100vh; }
```

Columns are **left panel · map stage · (optional) right detail panel**. The grid
template is swapped by two boolean classes that mirror React state:

| Shell state | Classes | `grid-template-columns` | Trigger |
|-------------|---------|-------------------------|---------|
| Default | `appShell` | `420px  minmax(0,1fr)` | initial |
| Left collapsed | `appShell leftCollapsed` | `80px  minmax(0,1fr)` (`:93`) | `leftCollapsed` state (sidebar toggle) |
| Detail open | `appShell detailOpen` | `420px  minmax(0,1fr)  420px` (`:97`) | a venue is selected |
| Both | `appShell leftCollapsed detailOpen` | `80px  minmax(0,1fr)  420px` (`:101`) | both |

Class composition (`DiscoveryExperience.tsx:401-405`):

```jsx
["appShell", leftCollapsed ? "leftCollapsed" : "", selectedVenue ? "detailOpen" : ""]
  .filter(Boolean).join(" ")
```

`data-theme={theme}` on the same element drives light/dark (see
[foundations](01-foundations.md#theming-mechanism)).

## The three panels

### Left panel — `.leftPanel` (`:105-135`)

- `position: relative; z-index: 4; height: 100vh; overflow: auto`.
- `background: var(--surface)`, `border-right: 1px solid var(--outline)`.
- `transition: width 0.3s ease` for the collapse animation.
- Thin custom scrollbar: `scrollbar-width: thin` + 4px WebKit thumb in
  `--surface-highest` (`:116-130`).
- **Internal stack (JSX order):** `.brandBlock` → `.sideTabs` → `.searchFocus`
  (headline + hero search + context chips) → `.compactFilters` (filter grid + tags +
  toggle/sort + actions) → `.venueRail` (recommendations) → `.authCta`.
- When collapsed, this whole stack is replaced by `.collapsedTabs` (3 icon buttons)
  — a conditional render, not a CSS hide (`DiscoveryExperience.tsx:430-584`).

### Map stage — `.mapStage` (`:735-740`)

- `position: relative; min-height: 100vh; overflow: hidden`.
- Contains, absolutely stacked: `.mapFallbackGrid` → `.mapCanvas` (MapLibre) →
  `.mapShade` → marker layers → overlays (`.mapCanvasWrap` children, `:970-1028`).
- `.mapTopbar` (`:742-755`) is `position: absolute` across the top with
  `pointer-events: none`; its direct children re-enable `pointer-events: auto`
  (`:757`) so the map stays draggable between controls.
- `.mapBottomAction` (`:1136`) is a centered floating pill at `bottom: 32px`.

### Right detail panel — `.rightPanel` (`:105-144`)

- Same base as left panel but `border-left` + a strong left shadow
  `-24px 0 80px rgba(0,0,0,0.44)`.
- Only mounts when `selectedVenue` is truthy (`DiscoveryExperience.tsx:632-657`);
  contains the `.closeButton` and `<VenueDetail>`.

### Sidebar toggle — `.sidebarToggle` (`:146-166`)

- `position: fixed; z-index: 8; left: 400px; top: 50%` — a 40px circle straddling the
  left-panel edge.
- Collapsed state slides it with `transform: translateX(-340px) …` using `!important`
  (`:165`). Hidden entirely below 1180px.

## Scroll model

- **Desktop:** `body { overflow: hidden }` (`:35`) — the page itself never scrolls.
  The layout is a locked full-viewport frame; the left and right panels scroll
  **internally** (`overflow: auto`, `height: 100vh`), and the map fills its column.
- **Below 1180px:** `body { overflow: auto }` (`:1497`) — the page scrolls normally
  and panels become `height: auto` (see below).

## Responsive breakpoints

Only **two** media queries exist.

### `@media (max-width: 1180px)` (`:1495-1540`) — collapse to single column

- Body becomes scrollable.
- **All four** shell states collapse to one column:
  `grid-template-columns: minmax(0,1fr)` (`:1500-1505`) — i.e. panels stack vertically.
- `.leftPanel`/`.rightPanel`: `height: auto; min-height: 0; border: 0`.
- `.sidebarToggle { display: none }` (collapse is desktop-only).
- Collapsed brand realigns to the start; `.collapsedTabs` becomes a 3×48px row.
- `.mapStage { min-height: 620px }` (map gets a fixed floor since it no longer fills a
  viewport column).

### `@media (max-width: 640px)` (`:1542-1605`) — phone tuning

- Content gutters shrink from 24px to **16px** (`:1543-1560`).
- `h1` → 25px.
- `.filterGrid`, `.detailTitleRow`, `.detailActions` → single column (`:1566-1570`).
- `.venueRailCard` thumbnail: 96px → **82px** (`:1572-1579`).
- `.mapTopbar` stacks vertically (`flex-direction: column`, auto height); `.mapStats`
  goes full-width `space-between`.
- Selected markers widen to `min-width: 150px`; `.detailHero` height 320px → **240px**.

> There is no explicit tablet/`min-width` query and no container queries. Layout is
> desktop-first with two `max-width` downshifts.
