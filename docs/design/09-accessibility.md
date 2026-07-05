# Accessibility

> **Status: As-built.** Documents the current a11y patterns and the concrete gaps. No
> remediation is applied here.

## Patterns in place

- **Landmark labels:** each region is a labelled landmark —
  `aside[aria-label="Discovery controls"]`, `section[aria-label="Restaurant map"]`,
  `aside[aria-label="Venue detail"]`, plus `nav[aria-label]` and section-level labels
  (`DiscoveryExperience.tsx:407, 432, 443, 459, 474, 547, 587, 633`).
- **Icon-only buttons carry `aria-label`:** sidebar toggle, close, help, markers
  (`aria-label="Select <name>"`), sort select (`aria-label="Sort venues"`).
- **Decorative content is `aria-hidden`:** every `.uiIcon` glyph, `.mapFallbackGrid`,
  `.mapShade`, marker pins (`:972, 974, 840, 1346`).
- **Native controls:** selects and the open-now checkbox are native `<select>`/
  `<input type=checkbox>`, wrapped in `<label>` for implicit naming.
- **Alt text:** venue/hero/video images use descriptive `alt`; decorative marker images
  use `alt=""`.

## Known gaps

| # | Gap | Where |
|---|-----|-------|
| A1 | **No `:focus-visible` style anywhere** — keyboard focus relies on UA default | all interactive elements |
| A2 | **Collapsed mini-tabs lack an accessible name** — icon-only, glyphs `aria-hidden`, no `aria-label` | `.miniTab` (`:573-583`) |
| A3 | **Toggle/selection state not announced** — tags/rail cards/tabs show state by color only (no `aria-pressed`/`aria-current`/`aria-selected`) | tags, rail card, tabs |
| A4 | **No `aria-live` for async updates** — search results, detail load, route changes update silently | list/detail/route |
| A5 | **Loading uses only `disabled`** — no `aria-busy` on loading buttons | primary/glass actions |
| A6 | **Inert enabled buttons** — "Bộ lọc"/"Chia sẻ"/"Lưu"/help/auth look actionable but do nothing (neither wired nor `disabled`) | topbar/detail/auth |
| A7 | **No `prefers-reduced-motion`** — hover scale + 700ms camera moves aren't reduced | motion |
| A8 | **px type doesn't scale** with root font-size preference | typography |

## Rules

- Every icon-only control needs an `aria-label` (extend this to the mini-tabs).
- Keep decorative layers `aria-hidden` and provide `alt` on meaningful images.
- Prefer native controls; wrap inputs in `<label>`.

## Do / Don't

- **Do** add `aria-label` to any new icon-only button.
- **Do** convey state with `aria-pressed`/`aria-current`, not color alone.
- **Don't** ship an enabled button with no action (mark it `disabled` until wired).
- **Don't** rely on the glow/color as the only focus cue.

## Related

- [Icons](07-icons.md) · [Motion](08-motion.md) · [Typography](06-typography.md)
- Prioritized fixes: [improvement proposals §A3/§B](../archive/design-system/improvement-proposals.md#a3-missing-components-gaps-in-the-system-not-just-the-docs)

## Files

- `apps/web/components/DiscoveryExperience.tsx` — labels/aria throughout.
- `apps/web/app/globals.css` — (no focus styles defined).

## References

- `apps/web/CLAUDE.md` (a11y & copy section).

## Future improvements

- Add one shared `:focus-visible` treatment (highest-value, lowest-effort fix — A1).
- Name the mini-tabs (A2); add `aria-pressed`/`aria-current` (A3); add `aria-live`/`aria-busy` (A4/A5).
- Wire or disable the inert buttons (A6); add `prefers-reduced-motion` (A7).
