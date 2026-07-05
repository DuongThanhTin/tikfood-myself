# Icon

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

A single lightweight icon primitive. Icons are **unicode glyphs**, not SVGs or an icon
font — a name→character map rendered inside a fixed 1em box. Used as decoration inside
buttons, chips, markers, and headers.

## Usage

- `.uiIcon` (`globals.css:55-63`): `display:inline-grid`, `1em × 1em`, `place-items:center`,
  `font-size:20px`, `font-weight:900`, `line-height:1`. Callers scale/recolor via
  `font-size`/`color` on the surrounding rule (e.g. hero search sets 24px + primary).
- `Icon` component (`DiscoveryExperience.tsx:1344-1350`): takes `name: IconName` and an
  optional `className`, renders `iconGlyphs[name]` inside `.uiIcon` with `aria-hidden="true"`.
- Glyph table (`iconGlyphs`, `:66-86`):

| Name | Glyph | Name | Glyph | Name | Glyph |
|------|-------|------|-------|------|-------|
| bookmark | ▣ | help | ? | share | ↗ |
| chevronLeft | ‹ | location | ⌖ | spark | ✦ |
| chevronRight | › | map | ▦ | star | ★ |
| close | × | moon | ◐ | sun | ☼ |
| explore | ◆ | play | ▶ | target | ◎ |
| fire | ● | route | ↱ | tune | ☰ |
| search | ⌕ | | | | |

Separately, map-marker category icons use **emoji** (`☕`/`🍽`) via `getVenueMapIcon`
(`:1259-1265`) — not part of `iconGlyphs`.

## Examples

```jsx
{/* Standalone / inside a button (:414, 605) */}
<Icon name={leftCollapsed ? "chevronRight" : "chevronLeft"} />
<Icon name="help" />

{/* With a positioning className (:456) */}
<Icon name="spark" className="heroSearchSpark" />
```

## Rules

- Icons are **decorative**: `.uiIcon` is always `aria-hidden`. Never rely on an icon to convey the accessible name — pair it with text or an `aria-label` on the parent control.
- Add new icons by extending the `IconName` union **and** `iconGlyphs` together (the type keeps them in sync).
- Size/color come from the context rule (`font-size`, `color`), not from `.uiIcon` itself, aside from its 20px default.
- Weight is fixed at 900 in `.uiIcon`, which isn't a loaded font weight (see R4) — glyph rendering is font-agnostic so this mostly affects the character metrics.

## Do

- Reuse an existing glyph name; check the table before adding one.
- Give the **parent** interactive element the `aria-label` when the icon is the only content.
- Let the parent set size and color.

## Don't

- Don't use an icon as the sole meaning of a control without an accompanying label.
- Don't hand a raw unicode character to markup; go through the `Icon` component so it stays typed and `aria-hidden`.
- Don't mix the emoji marker icons into `iconGlyphs`; they serve a different (map) purpose.

## Accessibility

- Every glyph is `aria-hidden="true"`, so it's skipped by screen readers by design.
- Meaning must come from adjacent text or the parent's `aria-label` (the pattern buttons already follow).
- Because glyphs are text, they inherit color/contrast from the parent — verify contrast when recoloring (e.g. white on gradient chips).

## Related components

- [Button](button.md), [Chip & Tag](chip-tag.md), [Navigation Tabs](navigation-tabs.md), [Map Marker](map-marker.md) — all host icons.
- [Hero Search](hero-search.md) — leading/trailing glyphs.

## Files

- `apps/web/app/globals.css` — `.uiIcon` `:55-63`.
- `apps/web/components/DiscoveryExperience.tsx` — `IconName` `:45-64`; `iconGlyphs` `:66-86`; `Icon` `:1344-1350`; `getVenueMapIcon` `:1259-1265`.

## References

- [Patterns — Icon glyph system](../04-patterns.md#icon-glyph-system)
- [Foundations — Iconography](../01-foundations.md#iconography-summary)
- [Inconsistency register — R4](../05-inconsistencies.md)

## Future improvements

- Glyph icons vary in visual weight/baseline across platforms; an SVG set would render consistently. Recorded as an option, not a proposal.
