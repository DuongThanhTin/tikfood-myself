# 04 — Patterns

> Page, navigation, map, state, copy, and accessibility patterns.
> Source: `components/DiscoveryExperience.tsx`, `lib/api.ts`, `globals.css`.

## Page pattern: three-column discovery

The app is a single screen — a **filter/list rail · map · detail** triptych:

```
┌─────────────┬───────────────────────────┬─────────────┐
│ leftPanel   │        mapStage           │ rightPanel  │
│ brand       │  topbar (filters/theme)   │ close ✕     │
│ nav tabs    │                           │ hero image  │
│ hero search │      MapLibre map         │ title+rating│
│ chips       │      + markers            │ trend score │
│ filters     │                           │ videos      │
│ venue rail  │  ┌ count legend ┐         │ dishes      │
│ auth CTA    │  bottom "Tìm quanh đây"   │ hours       │
└─────────────┴───────────────────────────┴─────────────┘
```

The right panel is conditional (only when a venue is selected); the left panel
collapses to an icon rail. See [`02-layout.md`](02-layout.md) for the grid mechanics.

## Navigation patterns

There is **no router-based navigation** — it's a single page. "Navigation" is
in-panel state:

| Surface | Items | Behavior |
|---------|-------|----------|
| Side tabs `.sideTabs` (`:432-441`) | "Khám phá" (active), "Bản đồ" | Bản đồ triggers `useCurrentLocation` — it's an action, not a route |
| Collapsed mini-tabs `.collapsedTabs` (`:573-583`) | explore (active), target (locate), bookmark (**disabled**) | Shown only when left panel is collapsed |
| Map topbar `.mapTopbar` (`:588-609`) | "Bộ lọc" (no handler), theme toggle, "Gần tôi", help (no handler) | Floating controls over the map |

> "Bộ lọc", "Chia sẻ", "Lưu", the help button, and the auth button are rendered but
> have **no `onClick`** — they are visual-only in the current build.

## Icon glyph system

Icons are **unicode glyphs**, not SVG or an icon font. `iconGlyphs`
(`DiscoveryExperience.tsx:66-86`) maps 19 names → characters, rendered via `.uiIcon`
(1em box, weight 900, `aria-hidden`):

| Name | Glyph | Name | Glyph | Name | Glyph |
|------|-------|------|-------|------|-------|
| bookmark | ▣ | help | ? | share | ↗ |
| chevronLeft | ‹ | location | ⌖ | spark | ✦ |
| chevronRight | › | map | ▦ | star | ★ |
| close | × | moon | ◐ | sun | ☼ |
| explore | ◆ | play | ▶ | target | ◎ |
| fire | ● | route | ↱ | tune | ☰ |
| search | ⌕ | | | | |

Separately, **emoji** are used for map-marker category icons via `getVenueMapIcon`
(`:1259-1265`): `☕` for cafe/coffee names, `🍽` otherwise.

## Map interaction patterns

- **Engine:** MapLibre GL, **lazy-imported** (`import("maplibre-gl")`, `:730`) so it's
  off the initial bundle. Basemaps are CARTO styles: `dark-matter-gl` (dark) /
  `positron-gl` (light) (`:136-139`).
- **Initial camera:** center `[106.683, 10.778]` (HCMC), zoom 12.8 (`:737-738`).
- **Marker mode by zoom:** `getMapMarkerMode` (`:218-220`) → `cluster` below zoom 12.2,
  else `venue`. (`getVenueClusterCount` currently always returns 1 — clustering is a
  stub, `:1253-1257`.)
- **Camera moves:** `flyTo` zoom 14.5 on venue select (`:863-867`); zoom 15 on user
  location (`:897-901`); `fitBounds` with asymmetric padding on route (`:956-964`).
- **Basemap tuning:** `applyMapLayerPreferences` (`:145-216`) hides small place labels
  and repaints water/parks/roads with a bespoke warm palette per theme — applied
  imperatively after style load, outside the CSS token system (register **R5/R8**).
- **Routing:** `fetchRoute` calls the public OSRM demo server
  (`router.project-osrm.org`, `:1213-1251`); the line is drawn as two stacked layers
  (casing + blue line, `:921-952`). Requires user location first, else a Vietnamese
  hint is shown.
- **Fallback layer:** before the map is ready (`!mapReady`), a CSS grid backdrop
  (`.mapFallbackGrid`) plus HTML markers positioned by the linear projection
  `toMapX/toMapY` (`:1336-1342`) stand in. The map canvas also has a CSS filter for
  saturation/contrast and a `.mapShade` vignette (`:837-867`).

## UI states

| State | Pattern |
|-------|---------|
| Loading (search) | `isLoading` → button label "Đang tìm", buttons disabled (`:535-540`) |
| Loading (detail) | `isDetailLoading` → `.detailStatus` "Đang tải chi tiết quán..." (`:1115`) |
| Loading (route) | `isRouting` → action label "Đang tải", disabled (`:1073`) |
| Error | `.errorText` (search) / `.routeHint.error` (route + detail errors) in `--danger` |
| Empty | `.emptyText` "Không có địa điểm phù hợp bộ lọc." (`:554`) |
| Disabled | primary/secondary `opacity:.7 cursor:wait`; glassAction `.58`; miniTab `.45 not-allowed`; sort `select` disabled when `nearUser` is set (`:524`) |
| Selected | `.selected` modifier on rail card + markers; drives detail-panel mount |

## Data-flow pattern

- **Server → client handoff:** `app/page.tsx` awaits `getDiscoveryVenues()` and passes
  `initialVenues` into the client `DiscoveryExperience` (`page.tsx:4-7`).
- **All data through `lib/api.ts`** (per `apps/web/CLAUDE.md`): `getDiscoveryVenues`
  (server, `revalidate:30`), `fetchDiscoveryVenues` / `fetchVenueDetail` (client).
- **`{ data, error }` envelope** (`ApiResponse<T>`, `lib/api.ts:63-70`) unwrapped to
  `.data`; error `message` surfaced to the UI.
- **Fallback data:** when `NEXT_PUBLIC_API_URL` is unset, `fallbackVenues` (2 seeded
  venues) is used and filtered client-side (`filterFallbackVenues`, `:243-282`).
- **Lazy detail:** selecting a venue fetches full detail only if `dishes`/
  `opening_hours` are absent (`DiscoveryExperience.tsx:255-285`).
- **Presentation vs data:** ratings, reviews, badges, and imagery are **not** from the
  API — they come from the hardcoded `mediaByVenue` map (register **R7**).

## Copy convention (bilingual)

Per `apps/web/CLAUDE.md`, **UI copy is Vietnamese**; technical identifiers are English.
Observed in the code: user-facing strings ("Chào mừng bạn!", "Tìm kiếm", "Gần tôi",
"Đang mở cửa", "Chỉ đường", error messages) are Vietnamese; component/type/function
names, class names, and API params are English. Day labels are Vietnamese abbreviations
(`formatDay` → CN/T2…T7, `:1332-1334`); currency formats as `k`/`tr` VND
(`formatPrice`, `:1312-1320`).

## Accessibility patterns

- **`aria-label` on icon-only buttons** and landmark regions: sidebar toggle, close,
  help, mini-tabs, markers (`aria-label="Select ${name}"`), and each `<aside>`/`<nav>`/
  `<section>` region (`:407, 432, 443, 459, 474, 547, 587, 633`).
- **`aria-hidden="true"`** on all decorative layers: every `.uiIcon` glyph, the
  fallback grid, the map shade, marker pins (`:972, 974, 840, 1346`).
- **Native labels wrap inputs** — filters/search use `<label class="field">`/
  `<label class="heroSearch">` around their controls; the sort `<select>` carries an
  explicit `aria-label="Sort venues"` (`:526`).
- **`alt` text on images** — venue thumbs/hero/videos use descriptive alt built from
  the venue name; purely decorative marker images use `alt=""` (`:839, 1002, 1060`).

> These are the *current* patterns, documented as-is. Gaps (e.g. no-op buttons lacking
> `disabled`, glyph icons conveying meaning only via adjacent text) are noted in the
> register, not remediated here.
