# Input

> **Status: As-built.** Tokens: [`../03-color-system.md`](../03-color-system.md),
> [`../04-spacing-system.md`](../04-spacing-system.md).

## Purpose

Form controls for discovery: the hero **search** field, labelled **selects** (district,
max price, sort), and the open-now **toggle**. All are native controls, styled with
plain CSS — no custom widgets.

## Usage

### Hero search — `.heroSearch` (`globals.css:325-395`)
A `<label>` wrapping an `<input>`. Blurred gradient glow via `::before`
(opacity `.32 → .62` on `:focus-within`); `input` `min-height:64px`, `2px` rgba-primary
border, radius 16, `padding:0 54px` (room for the two glyphs), bg `--surface-highest`,
16px/700. Leading `search` glyph + trailing `spark` glyph (`.heroSearchSpark`, uses
`!important`).

### Fields & selects — `.field` / `.field select` / `.sortSelect` (`:455-503`)
`.field` is a `<label>` grid (caption `<span>` + control). Selects: `width:100%`,
`min-height:42px`, radius 10, `1px var(--outline)`, bg `--surface-low`. `.filterGrid`
pairs two (2-col → 1-col ≤640px). `.sortSelect` is standalone (`max-width:150px`,
`aria-label="Sort venues"`).

### Toggle — `.toggle` / `.toggle input` (`:486-499`)
A `<label>` wrapping a native checkbox styled only via `accent-color: var(--primary)`,
18×18.

## Examples

```jsx
{/* Search (:449-457) */}
<label className="heroSearch">
  <Icon name="search" />
  <input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="VD: date dưới 500k ở Quận 1" />
  <Icon name="spark" className="heroSearchSpark" />
</label>

{/* Select (:476-483) */}
<label className="field">
  <span>Quận</span>
  <select value={district} onChange={(e) => setDistrict(e.target.value)}>
    <option value="">Tất cả</option><option value="District 1">Quận 1</option><option value="District 3">Quận 3</option>
  </select>
</label>

{/* Toggle (:512-519) */}
<label className="toggle"><input type="checkbox" checked={openNow} onChange={(e) => setOpenNow(e.target.checked)} /><span>Đang mở cửa</span></label>
```

## Rules

- Wrap captioned controls in `<label>`; a caption-less control (sort) needs `aria-label`.
- Keep the search input's `padding:0 54px` in sync with the two absolute glyphs.
- Option values map to `VenueSearchParams` (`lib/api.ts:72-87`); the empty option means "unset" and serializes away (`toSearchParams`).
- Search executes on the explicit button, not on keystroke (typing only updates `query`).
- Theme the checkbox with `accent-color` only; keep it a native checkbox.

## Do / Don't

- **Do** use native `<select>`/`<input>` for built-in keyboard + a11y.
- **Do** keep captions short and Vietnamese.
- **Don't** build a custom dropdown or switch.
- **Don't** add more `!important` (the `heroSearchSpark` one is already flagged).

## Accessibility

- Captioned controls are named by their wrapping `<label>`; the sort select uses `aria-label`.
- Glyphs are decorative (`aria-hidden`).
- Focus shows via the search glow + native focus; there is no dedicated `:focus-visible` — see [`../09-accessibility.md`](../09-accessibility.md) (A1).

## Related components

- [Button](button.md) (search/reset actions) · [patterns/search.md](../patterns/search.md) · [patterns/filter.md](../patterns/filter.md) · [Icon](../07-icons.md)

## Files

- `apps/web/app/globals.css` — hero search `:325-395`; field/select `:449-503`; toggle `:486-499`.
- `apps/web/components/DiscoveryExperience.tsx` — `:449-532`.
- `apps/web/lib/api.ts` — `VenueSearchParams` `:72-87`; `toSearchParams` `:229-241`.

## References

- [Color](../03-color-system.md) · [Spacing](../04-spacing-system.md)

## Future improvements

- Add a shared `:focus-visible`; remove the `!important` via cleaner specificity.
- The `spark` glyph implies AI search that isn't implemented — align affordance with capability.
