# Icons

> **Status: As-built.** Icons are **unicode glyphs**, not SVGs or an icon font.

## The system

- `.uiIcon` (`globals.css:55-63`): `display:inline-grid`, `1em × 1em`, `place-items:center`,
  `font-size:20px`, `font-weight:900`, `line-height:1`. Size/color are set by the
  surrounding rule (e.g. hero search overrides to 24px + primary).
- `Icon` component (`DiscoveryExperience.tsx:1344-1350`): `name: IconName` + optional
  `className`, renders `iconGlyphs[name]` inside `.uiIcon` with `aria-hidden="true"`.

## Glyph table (`iconGlyphs`, `:66-86`)

| Name | Glyph | Name | Glyph | Name | Glyph |
|------|-------|------|-------|------|-------|
| bookmark | ▣ | help | ? | share | ↗ |
| chevronLeft | ‹ | location | ⌖ | spark | ✦ |
| chevronRight | › | map | ▦ | star | ★ |
| close | × | moon | ◐ | sun | ☼ |
| explore | ◆ | play | ▶ | target | ◎ |
| fire | ● | route | ↱ | tune | ☰ |
| search | ⌕ | | | | |

## Map category icons (separate)

Map markers use **emoji**, not `iconGlyphs`, via `getVenueMapIcon` (`:1259-1265`):
`☕` for cafe/coffee names, `🍽` otherwise.

## Examples

```jsx
<Icon name={leftCollapsed ? "chevronRight" : "chevronLeft"} />
<Icon name="spark" className="heroSearchSpark" />
```

## Rules

- Icons are **decorative** (`aria-hidden`); meaning comes from adjacent text or the parent's `aria-label`.
- Add an icon by extending the `IconName` union **and** `iconGlyphs` together (type keeps them in sync).
- The parent sets size/color; `.uiIcon` only provides the 20px default + centering.

## Do / Don't

- **Do** reuse an existing glyph name; check the table first.
- **Do** put the accessible name on the parent control when the icon is the only content.
- **Don't** use an icon as a control's sole meaning without a label.
- **Don't** fold the emoji map icons into `iconGlyphs` (different purpose).

## Accessibility

- Every glyph is `aria-hidden`, so it's skipped by screen readers by design.
- Icon-only controls must carry `aria-label` (buttons do; note the collapsed mini-tabs
  currently don't — see [`09-accessibility.md`](09-accessibility.md)).
- Glyphs are text: they inherit color/contrast from the parent — verify contrast when recoloring.

## Related

- [Typography](06-typography.md) · [Accessibility](09-accessibility.md) · [components/button.md](components/button.md)

## Files

- `apps/web/app/globals.css` — `.uiIcon` `:55-63`.
- `apps/web/components/DiscoveryExperience.tsx` — `IconName` `:45-64`; `iconGlyphs` `:66-86`; `Icon` `:1344-1350`; `getVenueMapIcon` `:1259-1265`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Glyphs render inconsistently across platforms/fonts; an inline SVG set would be uniform.
- Give icon-only controls (collapsed mini-tabs) accessible names.
