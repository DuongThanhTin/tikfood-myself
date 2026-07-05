# Trend Score Card

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

A detail-panel card that visualizes a venue's trend score as a percentage with a
gradient progress track, followed by the AI summary sentence.

## Usage

- Shares a glass base with `.glassAction`: `.trendScoreCard, .glassAction`
  (`globals.css:1297-1302`) — `1px rgba(255,183,127,0.2)`, `background:var(--glass)`,
  `blur(18px)`.
- `.trendScoreCard` (`:1304-1327`): grid, `gap:10px`, radius 18, `padding:16px`; header
  row (`:1311-1316`) `space-between` with a label `span` (14px/800) and a `strong`
  percentage (`--primary`, Plus Jakarta, 24px).
- `.trendTrack` (`:1329-1334`): 8px pill track in `--surface-highest`; its fill `span`
  (`:1336-1341`) is a `--secondary → --trend-fire` gradient at an inline width.
- Body `p` (`:1343-1349`): italic, muted, 13px — the AI summary.

## Examples

```jsx
{/* :1104-1113 — trendPercent is clamped 1..99 */}
<section className="trendScoreCard">
  <div>
    <span>Trend Score</span>
    <strong>{trendPercent}%</strong>
  </div>
  <div className="trendTrack">
    <span style={{ width: `${trendPercent}%` }} />
  </div>
  <p>{venue.ai_summary}</p>
</section>
```

Where `const trendPercent = Math.min(99, Math.max(1, venue.trend_score));` (`:1054`).

## Rules

- The fill width is the **only** inline style (`width:${trendPercent}%`) — computed values are allowed inline per `apps/web/CLAUDE.md`; everything else stays in CSS.
- Clamp the score to 1–99 before rendering so the track always shows a sliver and never overflows.
- The percentage `strong` and the fill gradient are the emphasis; the body `p` is supporting AI copy.
- `venue.ai_summary` and `venue.trend_score` are **real API fields** (`lib/api.ts:19,23`) — unlike most other detail media.

## Do

- Feed `trend_score` from the API and clamp it (existing helper).
- Keep the fill as a percentage width so it scales with the track.
- Use this card once per detail view, near the top.

## Don't

- Don't move layout styles inline; only the computed fill width belongs there.
- Don't show an unclamped 0% or >100% track.
- Don't confuse this with mocked media — the trend data here is genuine.

## Accessibility

- The percentage is shown as text (`{trendPercent}%`) so it's readable independent of the visual bar.
- The track is a presentational `div`/`span`; there is no `role="progressbar"`/`aria-valuenow`, so assistive tech sees only the adjacent text (acceptable, but a progressbar role would be richer).

## Related components

- [Button](button.md) — `glassAction` shares this card's glass base.
- [Venue Detail Panel](venue-detail-panel.md) — the host.
- [Icon](icon.md) — used elsewhere in the detail body.

## Files

- `apps/web/app/globals.css` — glass base `:1297-1302`; card `:1304-1327`; track `:1329-1341`; body `:1343-1349`.
- `apps/web/components/DiscoveryExperience.tsx` — `:1054, 1104-1113`.
- `apps/web/lib/api.ts` — `trend_score` `:19`, `ai_summary` `:23`.

## References

- [Component catalog — Cards](../03-components.md#cards)
- [Foundations — Color (gradients)](../01-foundations.md#color-tokens)

## Future improvements

- Add `role="progressbar"` with `aria-valuenow/min/max` to expose the score to assistive tech.
