# Navigation Tabs

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The left-panel primary navigation. Two forms of the same nav, chosen by panel state:
- **Side tabs** (`sideTab`) — full-width labelled tabs when the panel is expanded.
- **Mini tabs** (`miniTab`) — icon-only square tabs when the panel is collapsed.

There is no routing; these are in-panel controls (some trigger actions like locating
the user). See [Patterns — Navigation](../04-patterns.md#navigation-patterns).

## Usage

- `.sideTabs` (`globals.css:240-245`): 2-column grid, `gap:6px`, `padding:8px 16px 18px`.
- `.sideTab` (`:247-260`): inline-flex, `min-height:42px`, centered, `gap:8px`,
  `border:0`, `border-radius:8px`, `--on-muted`, transparent bg, 12px/800.
  `:hover` → `--on-surface` on `--glass` (`:262`); `.active` → `--primary` on
  `rgba(255,183,127,0.1)` (`:267`).
- `.collapsedTabs` (`:707-712`): grid, centered, `gap:14px`; becomes a 3×48px row under
  1180px (`:1531-1535`).
- `.miniTab` (`:714-723`): 48×48 grid, `1px var(--outline)`, `border-radius:14px`,
  `--on-muted`, `--glass`. `.active` (`:725`) → primary; `:disabled` (`:730`) →
  `opacity:.45; cursor:not-allowed`.

## Examples

```jsx
{/* Expanded — side tabs (:432-441) */}
<nav className="sideTabs" aria-label="Primary navigation">
  <button className="sideTab active" type="button">
    <Icon name="explore" />
    Khám phá
  </button>
  <button className="sideTab" type="button" onClick={useCurrentLocation}>
    <Icon name="map" />
    Bản đồ
  </button>
</nav>

{/* Collapsed — mini tabs; bookmark is disabled (:573-583) */}
<div className="collapsedTabs" aria-label="Collapsed navigation">
  <button className="miniTab active" type="button"><Icon name="explore" /></button>
  <button className="miniTab" type="button" onClick={useCurrentLocation}><Icon name="target" /></button>
  <button className="miniTab" type="button" disabled><Icon name="bookmark" /></button>
</div>
```

## Rules

- Exactly one tab carries `active` at a time (mirrors the current view).
- Side tabs and mini tabs are alternate renders of the same nav — the collapsed state swaps the whole block, not a CSS hide (`:430-584`).
- A tab that is not yet available uses the native `disabled` attribute (e.g. the bookmark mini tab).
- Side tabs include an icon + Vietnamese label; mini tabs are icon-only.
- The containing `<nav>`/`<div>` carries an `aria-label`.

## Do

- Keep `active` in sync with the current section.
- Give icon-only mini tabs a real destination or mark them `disabled` (the bookmark tab is correctly disabled today).
- Use `<nav aria-label>` around the tab group.

## Don't

- Don't set `active` on more than one tab.
- Don't leave an icon-only mini tab enabled with no handler (unlike the disabled bookmark, that would be a silent no-op — cf. [R10](../05-inconsistencies.md#secondary-findings)).
- Don't rely on color alone to show `active` for accessibility (see below).

## Accessibility

- Groups use `aria-label` ("Primary navigation" / "Collapsed navigation").
- Mini tabs are icon-only; their glyphs are `aria-hidden`, so **they currently lack an accessible name** — a real gap for the collapsed nav (side tabs are fine, they have text).
- `active` state is conveyed only by color/background; there is no `aria-current` or `aria-selected`.
- Disabled tabs use the native attribute, which is announced.

## Related components

- [Icon](icon.md) — tab glyphs.
- [Button](button.md) — the base button conventions these extend.
- [App Shell](app-shell.md) — the collapse state that switches the two forms.

## Files

- `apps/web/app/globals.css` — `.sideTabs`/`.sideTab` `:240-270`; `.collapsedTabs`/`.miniTab` `:707-733`; responsive `:1531-1535`.
- `apps/web/components/DiscoveryExperience.tsx` — `:430-441, 573-584`.

## References

- [Patterns — Navigation](../04-patterns.md#navigation-patterns)
- [Component catalog — Buttons](../03-components.md#buttons)
- [Inconsistency register — R10](../05-inconsistencies.md)

## Future improvements

- Add `aria-current="page"` (or `aria-selected`) to the active tab.
- Give collapsed mini tabs an `aria-label` so they have accessible names.
