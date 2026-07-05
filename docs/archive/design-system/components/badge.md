# Badge

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

Small pill labels that flag status or category. Three instances share one base style:
the guest-mode badge, the per-card trending badge, and the hero-overlay detail badges.

## Usage

Shared base `.guestBadge, .venueMetaRow em, .detailBadges span` (`globals.css:211-222`):
pill (`border-radius:999px`), `1px rgba(254,183,0,0.2)` border, text `--secondary`,
`background:rgba(254,183,0,0.1)`, 10px/800, `letter-spacing:0.08em`, `text-transform:uppercase`.

| Badge | Class | Distinct styles |
|-------|-------|-----------------|
| Guest | `guestBadge` | `padding:4px 8px` (`:224-226`) |
| Card trend | `venueMetaRow em` | `max-width:126px`, ellipsis, `padding:3px 7px` (`:657-664`) |
| Detail overlay | `detailBadges span` | white text on `rgba(255,255,255,0.12)`, `blur(16px)`, `padding:6px 10px` (`:1212-1217`) |
| Detail overlay (first) | `detailBadges span:first-child` | solid `--trend-fire` bg (`:1219-1222`) |

## Examples

```jsx
{/* Guest badge (:420) */}
<span className="guestBadge">Guest View</span>

{/* Card trend badge — content is mocked media, see R7 (:685) */}
<em>{media.badge}</em>

{/* Hero overlay badges — literal strings; first child is the fire-colored one (:1061-1064) */}
<div className="detailBadges">
  <span>HOT</span>
  <span>TRENDING</span>
</div>
```

## Rules

- Badges are non-interactive `<span>`/`<em>` — never buttons or links.
- Uppercase + `letter-spacing:0.08em` is intrinsic to the style; pass already-short text.
- The card badge (`em`) truncates with ellipsis at 126px; keep labels terse.
- The detail overlay's **first** badge is always the emphasis (fire) color; order accordingly.

## Do

- Reuse the shared base rather than re-styling a new pill.
- Keep badge text to 1–2 words (they truncate / wrap poorly otherwise).
- Use the card `em` badge for the venue's trend label and the overlay for hero flags.

## Don't

- Don't put interactive content in a badge.
- Don't treat badge text as data-truth — hero "HOT"/"TRENDING" and card badges come from the mocked `mediaByVenue`, not the API (see [R7](../05-inconsistencies.md#r7--presentation-data-is-mocked-not-api-driven)).
- Don't add long strings; the fixed sizing will clip them.

## Accessibility

- Badges are plain text spans, read inline by screen readers.
- They carry meaning by text alone (no icon-only badges), so no extra labelling is needed.
- Because content is decorative/marketing ("HOT"), it adds no functional state to announce.

## Related components

- [Venue Rail Card](venue-rail-card.md) — hosts the `venueMetaRow em` badge.
- [Venue Detail Panel](venue-detail-panel.md) — hosts the `detailBadges` overlay.

## Files

- `apps/web/app/globals.css` — base `:211-222`; guest `:224-226`; card `em` `:657-664`; detail overlay `:1203-1222`.
- `apps/web/components/DiscoveryExperience.tsx` — `:420, 685, 1061-1064`; mock media `:105-132`.

## References

- [Component catalog — Badges](../03-components.md#badges)
- [Inconsistency register — R7](../05-inconsistencies.md)

## Future improvements

- Drive badges from real API signals (trend tier, "open now") instead of the mocked `mediaByVenue` map (R7).
