# Typography

> **Status: As-built. Canonical type reference.** Source: `globals.css` (`@import` on
> line 1 loads the fonts).

## Families

| Family | Weights loaded | Used for |
|--------|----------------|----------|
| **Be Vietnam Pro** | 400, 500, 700, 800 | Body/default (`body`, `:38`), `h1 span`, `.sectionHeader h2` |
| **Plus Jakarta Sans** | 600, 700, 800 | Display: `h1`, `h2`, `.brandName`, `.detailTitleRow h2`, numerals, `.primaryButton.large` |

Body fallback stack: `"Be Vietnam Pro", Arial, Helvetica, sans-serif`. Be Vietnam Pro is
chosen for full Vietnamese diacritic coverage.

## Type scale (as-built — no size tokens)

| Element / class | Size | Weight | Family | Notes |
|-----------------|------|--------|--------|-------|
| `h1` (`:298`) | 28px | 800 | Plus Jakarta | `line-height:1.18`; 25px ≤640px |
| `h1 span` (`:305`) | 18px | 400 | Be Vietnam | Sub-headline, muted |
| `h2` (`:314`) | 16px | 800 | Plus Jakarta | Base heading |
| `.sectionHeader h2` (`:570`) | 12px | — | Be Vietnam | Uppercase, `0.08em`, muted |
| `.detailTitleRow h2` (`:1237`) | 30px | 800 | Plus Jakarta | Venue name, `line-height:1.14` |
| `h3` (`:319`) | 15px | 800 | Be Vietnam | Detail subheads |
| `.brandName` (`:190`) | 24px | 800 | Plus Jakarta | Primary color |
| Body / input (`:357`) | 16px | 700 | Be Vietnam | |
| Controls (chips/tags/fields) | 12px | 800 | Be Vietnam | Dominant control size |
| Badges (`:213`) | 10px | 800 | — | Uppercase, `0.08em` |
| Numerals (rating/legend/route/trend) | 22–24px | 800 | Plus Jakarta | |

- **Font sizes present:** `10,11,12,13,14,15,16,18,20,22,24,25,28,30px` (14 distinct, no
  scale). `12px` dominates (×14); `11/12/13` are near-identical and all in use.
- **Line-heights:** `1, 1.14, 1.18, 1.35, 1.45, 1.5, 1.55` (7 ad-hoc values).
- **Letter-spacing:** only `0` (heading resets) and `0.08em` (uppercase labels).
- **`text-transform: uppercase`** only on badges and `.sectionHeader h2`.

## Weights (caution)

- Present in CSS: `400` (×1), `700` (×2), `800` (×22 — the norm), **`850` (×1)**,
  **`900` (×6)**.
- **`850` and `900` are not loaded** (faces are `400/500/700/800` and `600/700/800`) →
  they silently fall back to 800 or synthetic bold.
- `500` / `600` are loaded but unused.

## Rules

- Headings use Plus Jakarta; body/labels use Be Vietnam Pro.
- Reset margins are global (`h1,h2,h3,p { margin:0 }`, `:285-290`); manage spacing via the parent layout.
- Prefer weight `800` for emphasis (the established norm); avoid `850/900` until loaded.

## Do / Don't

- **Do** keep heading families consistent (Plus Jakarta display, Be Vietnam body).
- **Do** pick an existing size near your intent rather than a new px value.
- **Don't** rely on `font-weight` above 800 — it isn't loaded (falls back).
- **Don't** re-style `h2` to a new size ad hoc; it already means three different things.

## Accessibility

- Sizes are in px (not rem); they don't scale with the user's root font-size preference —
  a known limitation to weigh when adjusting type.
- Maintain contrast for muted text (`--on-muted`) on its surfaces.

## Related

- [Brand](02-brand.md) · [Color](03-color-system.md) · [Icons](07-icons.md)
- Consistency findings: [improvement proposals §B2](../archive/design-system/improvement-proposals.md#b2-inconsistent-typography-relates-to-r4)

## Files

- `apps/web/app/globals.css` — `@import` `:1`; headings `:285-323`; per-component sizes throughout.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Define a type scale (size + line-height + weight per role) and map the 14 sizes onto ~6 steps.
- Decouple heading semantics from visual size (role classes, not restyled `h2`).
- Load or drop `850/900`; drop unused `500/600` from the import.
- Consider rem-based sizing for user scaling.
