# Brand

> **Status: As-built.** Describes the brand expression present in the code.

## Name & product statement

- **Name:** TikFood.
- **Product statement (metadata):** title "TikFood Discovery", description "Realtime
  social food discovery for trending dishes and venues." (`app/layout.tsx:5-8`).
- **Positioning (product docs):** "TikTok + Google Maps for food" — realtime social food
  discovery; discovery only (`CLAUDE.md`).

## Wordmark

- `.brandName` (`globals.css:190-196`): "TikFood" in **Plus Jakarta Sans 800**, color
  `--primary`, 24px.
- **Collapsed logo:** when the left panel collapses, the wordmark becomes a 44px rounded
  gradient chip (`--primary → --trend-fire`) showing only the first letter in `#2f1500`
  (`:198-209`). This is the compact brand mark.

```jsx
{/* DiscoveryExperience.tsx:417-421 */}
<div className="brandHeader">
  <span className="brandName">TikFood</span>
  {!leftCollapsed ? <span className="guestBadge">Guest View</span> : null}
</div>
```

## Color identity

Warm, food-forward palette (full tokens in [`03-color-system.md`](03-color-system.md)):

- **Primary** amber `#ffb77f`, **primary-strong** `#ff8a00` (orange).
- **Secondary** gold `#feb700` (ratings, badges).
- **Trend fire** `#ff4d00` (trend gradients, hot markers).
- Near-black surfaces in dark mode; warm cream in light mode.

## Typographic identity

Two families (details in [`06-typography.md`](06-typography.md)):

- **Plus Jakarta Sans** — display/brand (wordmark, headings, numerals).
- **Be Vietnam Pro** — body and UI copy; chosen for full Vietnamese diacritic support.

## Voice & tone

- **Language:** Vietnamese, warm and direct — e.g. "Chào mừng bạn!" / "Hôm nay bạn muốn
  đi đâu?" (`DiscoveryExperience.tsx:445-447`).
- **Register:** friendly, second-person, concise; action labels are short verbs ("Tìm
  kiếm", "Gần tôi", "Chỉ đường").
- **Scope discipline:** copy never implies ordering/booking (discovery-only, per
  anti-goals).

## Rules

- Use the full wordmark when space allows; use the gradient first-letter chip only in the
  collapsed rail.
- Keep the wordmark in `--primary` (it re-themes correctly in light mode via the token).
- UI copy is Vietnamese; don't mix in English user-facing strings.

## Do / Don't

- **Do** pair Plus Jakarta (brand/headings) with Be Vietnam Pro (body).
- **Do** keep brand color references on tokens so light mode works.
- **Don't** hardcode the wordmark color; use `--primary`.
- **Don't** introduce commerce-flavored copy (cart/checkout/booking).

## Accessibility

- The "Guest View" badge and location hint are text (readable inline).
- The collapsed logo hides the wordmark text visually (`color: transparent`) — the
  brand becomes a decorative gradient chip; ensure a text alternative if it ever conveys
  state (currently purely decorative).

## Related

- [Color system](03-color-system.md) · [Typography](06-typography.md) · [Design principles](01-design-principles.md)

## Files

- `apps/web/app/globals.css` — `.brandBlock`/`.brandName`/collapsed logo `:168-209`; `.guestBadge`/`.locationHint` `:211-238`.
- `apps/web/app/layout.tsx` — metadata `:5-8`.
- `apps/web/components/DiscoveryExperience.tsx` — brand header `:417-428`.

## References

- `CLAUDE.md`, `docs/tikfood/product-vision.md`.

## Future improvements

- No standalone logo asset exists (the mark is CSS-only); a real logo/favicon set would
  strengthen the brand. (Observation, not applied.)
