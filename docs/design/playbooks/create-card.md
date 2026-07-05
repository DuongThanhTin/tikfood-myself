# Playbook: Create a Card

> Build a new card variant. Cards **exist** today — see [`../components/card.md`](../components/card.md)
> (venue rail, trend score, video). Extend that language; don't invent a new one.

## When to use

You need a content container that groups a unit (a venue, a metric, a media item).

## Steps

1. **Pick the closest existing variant** and reuse it: the rail card (`.venueRailCard`),
   the glass info card (`.trendScoreCard`/`.glassAction` base, `globals.css:1297-1302`),
   or the media card (`.videoCard`).
2. **Reuse tokens/shape**: radius **14px** (cards), glass surfaces via `--glass`,
   `blur(18px)`, `--shadow`; spacing per [`../04-spacing-system.md`](../04-spacing-system.md).
3. **If interactive**, make the whole card a single `<button>` (like the rail card) with
   inner non-interactive `<span>`s — one clear tap target.
4. **Images**: `object-fit:cover`, descriptive `alt`; decorative images `alt=""`; enforce
   ratios with `aspect-ratio` (like `.videoCard`).
5. **Data**: use real API fields; if you need media that the API lacks, flag it rather
   than hardcoding a mock (the existing `mediaByVenue` mock is a known issue).
6. **Document** the new variant in `../components/card.md`.

## Checklist

- [ ] Reused an existing card variant/base; radius/glass tokens used.
- [ ] Single tap target if interactive; semantic inner markup.
- [ ] `aspect-ratio` for media; `alt` correct.
- [ ] Real data (no new hardcoded media).
- [ ] `card.md` updated; meets [DoD](../../verification/definition-of-done.md).

## Related

- [../components/card.md](../components/card.md) · [create-component.md](create-component.md) · [../03-color-system.md](../03-color-system.md)

## References

- `apps/web/app/globals.css` (`:558-677, 1297-1349, 1417-1461`).
