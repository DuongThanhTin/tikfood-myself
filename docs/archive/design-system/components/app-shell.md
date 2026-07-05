# App Shell

> As-built. The full grid mechanics live in [`../02-layout.md`](../02-layout.md); this
> doc is the component view (usage/rules/a11y) and links there rather than repeating it.

## Purpose

The single top-level layout container (`<main className="appShell">`) that arranges the
whole app as a three-column CSS grid — left panel · map · optional detail panel — and
carries the theme. It also owns the collapse sidebar toggle.

## Usage

- `.appShell` (`globals.css:65-70`): `display:grid; grid-template-columns:420px minmax(0,1fr); min-height:100vh`.
- Two boolean classes swap the grid template: `.leftCollapsed` (`:93`, 80px rail) and
  `.detailOpen` (`:97`, adds a 420px third column). See the full state table in
  [`../02-layout.md`](../02-layout.md#the-app-shell).
- `data-theme={theme}` on the same element drives dark/light tokens (see
  [Foundations](../01-foundations.md#theming-mechanism)).
- `.sidebarToggle` (`:146-166`): fixed circular toggle straddling the left edge;
  translates when collapsed (uses `!important` — R12); hidden below 1180px.
- Panels: `.leftPanel` / `.rightPanel` (`:105-144`) scroll internally; `.mapStage` is
  the middle column (see [Venue Map](venue-map.md)).
- Responsive: everything collapses to one column under 1180px; gutters shrink under
  640px (`:1495-1605`).

## Examples

```jsx
{/* :398-406 — class + theme composition */}
<main
  data-theme={theme}
  className={["appShell", leftCollapsed ? "leftCollapsed" : "", selectedVenue ? "detailOpen" : ""]
    .filter(Boolean).join(" ")}
>
  <aside className="leftPanel" aria-label="Discovery controls">…</aside>
  <section className="mapStage" aria-label="Restaurant map">…</section>
  {selectedVenue ? <aside className="rightPanel" aria-label="Venue detail">…</aside> : null}
</main>
```

## Rules

- The grid template is driven **only** by the `.leftCollapsed` / `.detailOpen` classes, which mirror React state (`leftCollapsed`, `selectedVenue`) — never set columns inline.
- The right column mounts conditionally with the detail panel; it isn't a hidden column.
- Theme is set via `data-theme` on this element; all light tokens are scoped under it (anything outside `.appShell` stays dark — see [R9](../05-inconsistencies.md#r9--theme-is-scoped-to-appshell-only)).
- On desktop the page itself doesn't scroll (`body{overflow:hidden}`); panels scroll internally. Below 1180px the body scrolls and panels become `height:auto`.
- The sidebar toggle is desktop-only (hidden under 1180px).

## Do

- Compose the shell classes from state with the existing `filter(Boolean).join(" ")` idiom.
- Put the theme on `.appShell` so tokens cascade correctly.
- Let each panel manage its own internal scroll.

## Don't

- Don't render content that needs theming outside `.appShell` (it won't get light tokens — R9).
- Don't add more `!important`; the toggle's existing one is already flagged (R12).
- Don't hardcode column widths in components; the shell owns the grid.

## Accessibility

- Each region is a labelled landmark: `<aside aria-label="Discovery controls">`, `<section aria-label="Restaurant map">`, `<aside aria-label="Venue detail">` — good for screen-reader navigation.
- The collapse toggle has an `aria-label` that reflects its action ("Expand/Collapse left panel").
- Theme is not persisted and resets to dark on reload (R9) — no impact on semantics, but worth noting for user preference.

## Related components

- [Navigation Tabs](navigation-tabs.md) — the left panel's tab set (and its collapsed form).
- [Venue Map](venue-map.md) — the middle column.
- [Venue Detail Panel](venue-detail-panel.md) — the conditional right column.

## Files

- `apps/web/app/globals.css` — shell/states `:65-103`; panels `:105-144`; toggle `:146-166`; responsive `:1495-1605`.
- `apps/web/components/DiscoveryExperience.tsx` — `:398-660`; theme/collapse state `:226-227`.
- `apps/web/app/layout.tsx` — root html/body + font & maplibre CSS imports.

## References

- [Layout — full grid mechanics](../02-layout.md)
- [Foundations — Theming](../01-foundations.md#theming-mechanism)
- [Inconsistency register — R9, R12](../05-inconsistencies.md)

## Future improvements

- Persist theme (localStorage/cookie) and set it before paint to avoid a flash (R9).
- Extend theming to cover content rendered outside the shell, or document the constraint (R9).
