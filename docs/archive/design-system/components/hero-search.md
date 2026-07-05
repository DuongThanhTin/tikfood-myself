# Hero Search

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The primary free-text discovery input at the top of the left panel — a large search
field with a blurred gradient glow, a leading search glyph, and a trailing "spark"
glyph hinting at AI-assisted search.

## Usage

- `.heroSearch` (`globals.css:325-329`): `position:relative; display:grid`.
- Glow `::before` (`:331-345`): absolute `inset:-2px`, `z-index:-1`, radius 18,
  gradient `90deg, var(--primary), #ff8a00`, `filter:blur(12px)`, `opacity:.32`; on
  `:focus-within` opacity rises to `.62` (`:343`).
- `input` (`:347-359`): `min-height:64px`, `2px solid rgba(255,183,127,0.18)`,
  radius 16, `padding:0 54px` (room for both glyphs), bg `--surface-highest`,
  shadow `0 22px 60px rgba(0,0,0,0.36)`, 16px/700.
- Leading glyph `> .uiIcon` (`:383-389`): absolute `left:20px`, primary, 24px.
- Trailing glyph `.heroSearchSpark` (`:391-395`): `right:20px; left:auto !important`,
  muted primary — `!important` here signals a specificity override (see R12).
- Light theme retunes border/bg/glow (`:365-381`).

The whole thing is a `<label>` wrapping the `<input>`, so the field is labelled by
association.

## Examples

```jsx
{/* :449-457 */}
<label className="heroSearch">
  <Icon name="search" />
  <input
    value={query}
    onChange={(event) => setQuery(event.target.value)}
    placeholder="VD: date dưới 500k ở Quận 1"
  />
  <Icon name="spark" className="heroSearchSpark" />
</label>
```

## Rules

- Keep the `padding:0 54px` on the input in sync with the two absolutely-positioned glyphs so text never overlaps them.
- The glow is decorative and behind the input (`z-index:-1`); it must not intercept pointer events.
- Placeholder copy is a Vietnamese example query (`VD: …`), not a generic prompt.
- Search executes on the explicit "Tìm kiếm" button, not on input change (typing only updates `query` state).

## Do

- Keep the field a `<label>` wrapper so no separate `<label for>` / `aria-label` is needed.
- Use the leading `search` glyph and trailing `spark` glyph as the recognizable affordance.
- Preserve the focus-within glow intensification as the focus cue.

## Don't

- Don't remove the padding that reserves space for the glyphs.
- Don't rely on the glow for focus indication alone for accessibility — see below.
- Don't add more `!important`; the existing one (`heroSearchSpark`) is already flagged (R12).

## Accessibility

- The `<label className="heroSearch">` wrapper associates the caption/icons with the input; the input has an accessible name via the label context and its `placeholder`.
- Both glyphs are decorative (`aria-hidden` through `.uiIcon`).
- Focus is shown by the glow (`:focus-within`) plus the input's own border; there is no separate high-contrast focus ring, which may be insufficient for some users.

## Related components

- [Icon](icon.md) — search and spark glyphs.
- [Chip & Tag](chip-tag.md) — the quick-prompt chips directly below.
- [Button](button.md) — the "Tìm kiếm" action that runs the search.

## Files

- `apps/web/app/globals.css` — `:325-395`; light overrides `:365-381`.
- `apps/web/components/DiscoveryExperience.tsx` — `:449-457`; `query` state `:228`.

## References

- [Component catalog — Hero search](../03-components.md#form-controls)
- [Foundations — Motion / Glass](../01-foundations.md#motion)
- [Inconsistency register — R12 (`!important`)](../05-inconsistencies.md#secondary-findings)

## Future improvements

- The `spark` glyph implies AI search; the current build applies no AI behavior — align the affordance with actual capability or document it as aspirational.
- Add an explicit `:focus-visible` outline for keyboard accessibility.
