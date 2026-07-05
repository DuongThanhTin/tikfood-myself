# Grid & Layout System

> **Status: As-built. Canonical layout reference.** Source: `globals.css` +
> `DiscoveryExperience.tsx`.

## The app shell

The whole UI is one CSS-grid container, `.appShell` (`globals.css:65-70`), rendered as
`<main>`:

```
.appShell { display: grid; grid-template-columns: 420px minmax(0, 1fr); min-height: 100vh; }
```

Columns are **left panel · map · optional detail panel**, swapped by two state classes:

| State | Classes | `grid-template-columns` | Trigger |
|-------|---------|-------------------------|---------|
| Default | `appShell` | `420px  minmax(0,1fr)` | initial |
| Left collapsed | `appShell leftCollapsed` | `80px  minmax(0,1fr)` (`:93`) | sidebar toggle |
| Detail open | `appShell detailOpen` | `420px  minmax(0,1fr)  420px` (`:97`) | venue selected |
| Both | `appShell leftCollapsed detailOpen` | `80px  minmax(0,1fr)  420px` (`:101`) | both |

```jsx
{/* DiscoveryExperience.tsx:401-405 */}
["appShell", leftCollapsed ? "leftCollapsed" : "", selectedVenue ? "detailOpen" : ""]
  .filter(Boolean).join(" ")
```

## Panels

- **Left** `.leftPanel` (`:105-135`): `z-index:4`, `height:100vh`, `overflow:auto`,
  `--surface`, `border-right`, `transition: width .3s`. Hosts brand → nav → search →
  filters → rail → auth CTA (or the collapsed icon rail).
- **Map** `.mapStage` (`:735-740`): relative, `min-height:100vh`, `overflow:hidden`
  (see [`components/map.md`](components/map.md)).
- **Right** `.rightPanel` (`:137-144`): mirror of left with `border-left` + strong shadow;
  mounts only when a venue is selected.
- **Sidebar toggle** `.sidebarToggle` (`:146-166`): fixed circular edge toggle; slides on
  collapse via `transform … !important`; hidden below 1180px.

## Scroll model

- **Desktop:** `body { overflow: hidden }` (`:35`) — the page is a locked full-viewport
  frame; panels scroll internally.
- **≤1180px:** `body { overflow: auto }` — page scrolls, panels become `height:auto`.

## Breakpoints (only two, desktop-first)

- **`@media (max-width: 1180px)`** (`:1495-1540`): all four shell states collapse to a
  single `minmax(0,1fr)` column (panels stack); toggle hidden; `.collapsedTabs` → 3×48px;
  `.mapStage` gets `min-height: 620px`.
- **`@media (max-width: 640px)`** (`:1542-1605`): gutters 24→16px; `h1` 28→25px;
  `.filterGrid`/`.detailTitleRow`/`.detailActions` → 1 column; rail thumb 96→82px;
  topbar stacks; `.detailHero` 320→240px.

> Note: some sibling grids do **not** collapse at 640px (`.compactDetailActions`,
> `.videoGrid`) — see [improvement proposals §B4](../archive/design-system/improvement-proposals.md#b4-inconsistent-responsive-behavior).

## Rules

- Drive columns only via the `.leftCollapsed`/`.detailOpen` classes (mirror React state); never set grid columns inline.
- The right column mounts conditionally; it's not a hidden column.
- Let panels own their internal scroll; don't scroll the desktop body.

## Do / Don't

- **Do** compose shell classes with the `filter(Boolean).join(" ")` idiom.
- **Do** keep new full-height UI inside a panel that manages its own scroll.
- **Don't** hardcode column widths in components — the shell owns the grid.
- **Don't** add a third global scroll mode; the two-state model is intentional.

## Accessibility

- Regions are labelled landmarks: `aside[aria-label="Discovery controls"]`,
  `section[aria-label="Restaurant map"]`, `aside[aria-label="Venue detail"]`.
- The collapse toggle's `aria-label` reflects its action.

## Related

- [Spacing](04-spacing-system.md) · [components/map.md](components/map.md) · [pages/home.md](pages/home.md)

## Files

- `apps/web/app/globals.css` — shell/states `:65-103`; panels `:105-144`; toggle `:146-166`; responsive `:1495-1605`.
- `apps/web/components/DiscoveryExperience.tsx` — `:398-660`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Named breakpoint tokens + an intermediate tablet range (avoid the single 1180px cliff).
- Cap content width on ultra-wide screens (the map column is currently unbounded).
- Make sibling grids collapse consistently at 640px.
