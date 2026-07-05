# Spacing System

> **Status: As-built.** Today there is **no spacing scale or token** — values are raw
> px chosen per rule. This doc records the current values and the resulting
> inconsistency; a proposed scale is linked, not applied.

## Current state (no tokens)

Extracted from `apps/web/app/globals.css`:

- **`gap` — 14 distinct values:** `2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 18, 22, 26px`.
  `8px` is the most common; the small end (`2/3/4/5/6/7`) is used almost
  interchangeably.
- **`padding` — ~25 distinct shorthands**, several near-identical (`3px 7px`,
  `4px 8px`, `5px 10px 5px 6px`, `6px 10px`, `6px 10px 6px 6px`).
- **Panel gutter:** most left-panel sections use `margin: 0 24px` (`:272-277`); the brand
  block uses `30px 32px 16px`; under 640px all gutters drop to `16px`.

## Radii (no tokens)

- Non-pill radii: `8, 10, 12, 14, 16, 18px` (6 distinct); pill `999px` (×17);
  `inherit` (×2). Rough convention:

| Radius | Applied to |
|--------|-----------|
| 8px | side tabs |
| 10–12px | selects, buttons, thumbnails |
| 14px | rail/video cards, mini-tabs, selected marker |
| 16–18px | search, legends, trend card, auth CTA |
| 999px | all pills |

## Rules (current conventions to honor)

- Use `24px` as the left-panel horizontal gutter (`16px` on mobile) — the closest thing
  to a standard.
- Use `8px` for inline control gaps (chips, icon+label) and `16–18px` for card/section
  gaps — matches the dominant existing usage.
- Use `999px` for pills, `14px` for cards, `12px` for buttons/inputs.

## Do / Don't

- **Do** reuse an existing spacing value near your intent rather than inventing a new one.
- **Do** keep the `24px` gutter consistent across left-panel sections.
- **Don't** add another sub-8px gap variant; the `2–7px` cluster is already redundant.
- **Don't** hardcode the gutter in two places; it changes at one breakpoint only.

## Accessibility

- Maintain adequate hit-target spacing; controls use `min-height: 38–54px`, which keeps
  tap targets usable — preserve those minimums when tightening gaps.

## Related

- [Grid system](05-grid-system.md) (layout-level spacing) · [Color](03-color-system.md)
- Proposed scale: [improvement proposals §B1](../archive/design-system/improvement-proposals.md#b1-inconsistent-spacing-relates-to-r6)

## Files

- `apps/web/app/globals.css` — spacing/radii throughout (no `:root` scale defined).

## References

- `apps/web/CLAUDE.md` (reuse existing patterns).

## Future improvements

- Define `--space-1…8` (`4,8,12,16,20,24,28,32`) and `--radius-sm/md/lg/pill`, then map
  existing values to the nearest step.
- Collapse the redundant `2–7px` gaps toward `4/8`.
