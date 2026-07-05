# Toggle

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

A labelled boolean control. The only instance today is the "Đang mở cửa" (open-now)
filter. It is a **native checkbox** styled only via `accent-color`, not a custom
switch.

## Usage

- `.toggle` (`globals.css:486-493`): `display:flex; align-items:center; gap:8px`,
  text `--on-surface`, 13px/800.
- `.toggle input` (`:495-499`): `18px × 18px`, `accent-color: var(--primary)`.

The control is a `<label className="toggle">` wrapping the checkbox and its caption.

## Examples

```jsx
{/* :512-519 */}
<label className="toggle">
  <input
    type="checkbox"
    checked={openNow}
    onChange={(event) => setOpenNow(event.target.checked)}
  />
  <span>Đang mở cửa</span>
</label>
```

## Rules

- Use a native `<input type="checkbox">`; theme it with `accent-color` only.
- Wrap the input + caption in `<label className="toggle">` for implicit labelling and a larger hit target.
- `checked` is controlled by React state (`openNow`); reflect changes through `onChange`.
- The checkbox value feeds the `open_now` API param (`lib/api.ts`).

## Do

- Keep it a native checkbox for built-in keyboard and screen-reader support.
- Use `accent-color: var(--primary)` so the checked state matches the brand.
- Keep captions short and in Vietnamese.

## Don't

- Don't rebuild this as a styled `div` switch; you'd lose native semantics.
- Don't separate the caption from the input (breaks the implicit label).

## Accessibility

- Native checkbox → full keyboard operability and an announced checked/unchecked state.
- The wrapping `<label>` provides the accessible name and enlarges the click target.
- No `aria-*` needed; the native role and state are sufficient.

## Related components

- [Form Field & Select](form-field-select.md) — shares the filter actions row (sort select).
- [Button](button.md) — search/reset actions in the same filter block.

## Files

- `apps/web/app/globals.css` — `:486-499`.
- `apps/web/components/DiscoveryExperience.tsx` — `:511-532`; `openNow` state `:233`.
- `apps/web/lib/api.ts` — `open_now` param `:84`.

## References

- [Component catalog — Form controls](../03-components.md#form-controls)

## Future improvements

- If a switch UI is ever desired, layer it on the native checkbox (`::before`/`::after`) rather than replacing it, to preserve semantics.
