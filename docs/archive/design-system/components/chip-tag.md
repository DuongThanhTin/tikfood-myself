# Chip & Tag

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

Two pill-shaped controls that share a base style but differ in behavior:

- **Chip** (`contextChip`) — a quick-discovery prompt. Tapping it applies a preset
  query/price/sort and does not stay "pressed".
- **Tag** (`tag`) — a toggle filter that keeps an `active` state until re-tapped.
- **Trend chip** (`contextChip trend`) — a special gradient chip that highlights the
  "Đang hot" trending shortcut.

## Usage

Base rule `.contextChip, .tag` (`globals.css:408-423`): inline-flex, `min-height:38px`,
pill (`border-radius:999px`), `1px var(--outline)`, `background:var(--glass)`,
12px/800, `backdrop-filter:blur(18px)`. Rows are laid out by the shared flex container
group `.contextChips, .tagBar, …` (`:397-406`, `gap:8px`, wrap).

| Control | Class | State model | Data source |
|---------|-------|-------------|-------------|
| Chip | `contextChip` | stateless (fires `applyChip`) | `contextChips` array `DiscoveryExperience.tsx:88-95` |
| Trend chip | `contextChip trend` | stateless | same array, `kind:"trend"` |
| Tag | `tag` / `tag active` | toggled via `activeTag` state | `tagOptions` array `:97-103` |

## Examples

```jsx
{/* Chips — the trend kind gets a leading fire icon (:459-471) */}
{contextChips.map((chip) => (
  <button
    key={chip.label}
    className={chip.kind === "trend" ? "contextChip trend" : "contextChip"}
    type="button"
    onClick={() => applyChip(chip)}
  >
    {chip.kind === "trend" ? <Icon name="fire" /> : null}
    {chip.label}
  </button>
))}

{/* Tags — toggle active (:498-509) */}
{tagOptions.map((tag) => (
  <button
    key={tag.value}
    type="button"
    className={activeTag === tag.value ? "tag active" : "tag"}
    onClick={() => setActiveTag(activeTag === tag.value ? "" : tag.value)}
  >
    {tag.label}
  </button>
))}
```

## Rules

- Chips and tags are `<button type="button">`, never links.
- Hover (`:425-430`) lifts with `translateY(-1px)` and switches bg to `--surface-highest`; keep them on a surface where that reads.
- `tag.active` (`:475-479`) is the only persistent-selection state; chips never carry `active`.
- The trend chip (`:432-437`) drops its border for a `--secondary → --trend-fire` gradient and white text — reserve it for the single trending shortcut.
- `applyChip` branches on `chip.kind` (`price` → sets max price, `query` → sets query, else → sort trending) — a chip's effect is data-driven, not per-button code.

## Do

- Use a **tag** when the choice must stay selected (a filter); use a **chip** for a fire-and-forget preset.
- Keep exactly one trend chip.
- Add the `fire` icon only to the trend chip (existing pattern).
- Let the container group (`.contextChips`/`.tagBar`) own spacing and wrapping.

## Don't

- Don't give a plain chip an `active` class — there's no styling contract for it and it implies toggle behavior it doesn't have.
- Don't hardcode the amber again as `rgba(255,183,127,…)` for hover/active — those literals already duplicate `--primary` (see [R5](../05-inconsistencies.md#r5--blue-accent-has-no-token-warm-colors-re-expressed-as-literals)).
- Don't nest interactive elements inside a chip.

## Accessibility

- Each chip/tag exposes its accessible name via visible text.
- The trend chip's leading `<Icon name="fire">` is decorative (`aria-hidden` via `.uiIcon`), so the label alone conveys meaning.
- Toggle state is conveyed **only visually** (the `active` class) — there is no `aria-pressed`. Screen-reader users don't hear the selected state (accessibility gap).

## Related components

- [Button](button.md) — general-purpose (non-pill) actions.
- [Icon](icon.md) — the fire glyph on the trend chip.
- [Form Field & Select](form-field-select.md) — the other filter inputs the tag bar sits among.

## Files

- `apps/web/app/globals.css` — base `:408-423`; hover `:425-430`; trend `:432-441`; `tag.active` `:475-479`; containers `:397-406`.
- `apps/web/components/DiscoveryExperience.tsx` — data `:88-103`; `applyChip` `:354-364`; render `:459-471, 498-509`.

## References

- [Component catalog — Chips & tags](../03-components.md#chips--tags)
- [Foundations — Color / Glass](../01-foundations.md#glass--blur)
- [Inconsistency register — R5](../05-inconsistencies.md)

## Future improvements

- Add `aria-pressed` to tags so toggle state is announced.
- Reference `--primary`-derived values instead of repeating amber `rgba(...)` literals (R5).
