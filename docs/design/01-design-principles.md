# Design Principles

> **Status: As-built (derived).** These principles are *observed* from the shipped UI
> and the product framing in `CLAUDE.md` / `docs/tikfood/`. There is no pre-existing
> written manifesto in the repo — each principle below cites where it shows up in code
> so it stays grounded, not invented.

## 1. Map-first discovery

The map is the center column and fills the viewport; the left rail and right detail
panel orbit it. Selecting anything flies the map to it.

- *Seen in:* `.appShell` grid puts `.mapStage` as the primary `minmax(0,1fr)` column
  (`globals.css:65-70`); `flyTo` on select (`DiscoveryExperience.tsx:863-867`).

## 2. Dish- and venue-first content

Content leads with the food and the place — venue name, trending dishes, categories —
not with accounts or feeds.

- *Seen in:* `Venue` shape (`lib/api.ts:1-25`) centers `trending_dishes`, `categories`,
  `trend_score`; the detail panel surfaces dishes and social videos.
- *Product framing:* "dish-first, map-first" (`CLAUDE.md`).

## 3. Social proof & trend signals

Trendiness is a first-class visual: trend score bar, "Đang hot" chip, video counts,
"HOT"/"TRENDING" badges.

- *Seen in:* `.trendScoreCard` + `.trendTrack` gradient (`globals.css:1304-1341`); the
  trend context chip (`DiscoveryExperience.tsx:88-95`).

## 4. Dark-first glassmorphism

The default theme is near-black with warm amber accents and translucent
`backdrop-filter: blur()` surfaces; a warm-cream light theme is a full token override.

- *Seen in:* `:root` dark tokens + `.appShell[data-theme="light"]` (`globals.css:3-91`);
  `--glass`/`--glass-strong` used across chips, buttons, overlays.

## 5. Bilingual, Vietnamese-facing

UI copy is Vietnamese; technical identifiers are English. Formatting (currency `k`/`tr`,
day abbreviations) is localized.

- *Seen in:* strings like "Chào mừng bạn!", "Tìm kiếm" (`DiscoveryExperience.tsx`);
  `formatPrice`/`formatDay` (`:1312-1334`). Convention in `apps/web/CLAUDE.md`.

## 6. Discovery-only (hard anti-goals)

No cart, order, checkout, payment, booking, reservation, chat, follow graph, or
livestream. The UI stops at discovery + directions.

- *Seen in:* inert "Đăng nhập"/"Lưu"/"Chia sẻ" (no commerce flows); routing is
  directions-only (OSRM). Anti-goals in `CLAUDE.md` / `docs/tikfood/anti-goals.md`.

## 7. Progressive disclosure

Browse in the rail → select → the right panel expands with full detail; the left panel
can collapse to an icon rail to give the map room.

- *Seen in:* `.detailOpen` / `.leftCollapsed` grid states (`globals.css:93-103`);
  conditional detail panel mount (`DiscoveryExperience.tsx:632-657`).

## 8. Graceful fallback

The app degrades instead of breaking: seed venues when the API is absent, a CSS grid +
fallback markers before MapLibre loads.

- *Seen in:* `fallbackVenues` (`lib/api.ts:89-169`); `.mapFallbackGrid` +
  `.fallbackMarker` layer (`globals.css:823-831, 975-989`).

## How to apply

When adding UI, check it against these: does it keep the map central, lead with
food/place, express trend/social proof, fit the dark-glass aesthetic, use Vietnamese
copy, avoid anti-goal commerce, disclose progressively, and degrade gracefully?

## Related

- [Overview](00-overview.md) · [Brand](02-brand.md) · [Color](03-color-system.md)

## References

- `CLAUDE.md`, `docs/tikfood/product-vision.md`, `docs/tikfood/anti-goals.md`.

## Future improvements

- Ratify these observed principles into an intentional, agreed manifesto (today they are
  reverse-engineered, not authored).
