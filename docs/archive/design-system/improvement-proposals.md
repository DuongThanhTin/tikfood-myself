# Design System — Improvement Proposals

> **⚠️ ARCHIVED** (superseded by [`docs/design/`](../../design/00-overview.md)), but **still
> referenced**: the canonical foundations link to the §A/§B findings below for detail. The
> findings remain valid; the doc-set critique reflects the pre-migration `docs/design-system/**` layout.

> **Proposals only.** This is a review of the generated design-system documentation
> (`docs/design-system/**`) and the underlying `apps/web` code it describes. It
> **proposes** changes; it does not apply any. **No application code and no existing
> docs were modified** to produce this file.
>
> **Method.** Re-read the doc set and re-extracted value distributions directly from
> `apps/web/app/globals.css` (font sizes, weights, radii, gaps, colors, breakpoints).
> Findings link to the existing register [`05-inconsistencies.md`](05-inconsistencies.md)
> (R-IDs) instead of restating it.
>
> **Scope note.** "Consistency" findings describe the *code*; the audit's rule is to
> document, not fix. These proposals are therefore recommendations for a future,
> separately-approved change — not part of the read-only audit.

## How to read this

Each section: **Findings** (data-backed) → **Proposals** (P-IDs, concrete). A
prioritized summary is at the [end](#prioritized-summary).

---

## A. Documentation review

### A1. Duplicated documentation

**Findings**

| # | Duplication | Locations |
|---|-------------|-----------|
| D1 | The full 19-row **icon glyph table** appears twice, verbatim | [`04-patterns.md`](04-patterns.md#icon-glyph-system) and [`components/icon.md`](components/icon.md#usage) |
| D2 | **Two component indexes** list the same components | [`03-components.md`](03-components.md) (catalog tables) and [`components/README.md`](components/README.md) |
| D3 | Per-component descriptions + line refs exist in **both** the catalog and the detail docs (buttons, chips, badges, cards, markers, detail blocks) | `03-components.md` ↔ `components/*.md` |
| D4 | The **orphaned `VenueList`** note is stated in four places | `README.md`, `03-components.md`, `components/README.md`, `05-inconsistencies.md` (R1) |
| D5 | The **color-literal enumeration** (blue, `#4e2600`, rgba primaries) is listed in both | `01-foundations.md` (R5 bullet) and `05-inconsistencies.md` (R5) |
| D6 | The **R1–R9 inconsistency preview** in the roadmap duplicates the register | `ROADMAP.md#5` ↔ `05-inconsistencies.md` |

**Proposals**

- **P-D1** — Make [`components/icon.md`](components/icon.md) the single source of the glyph table; replace the copy in `04-patterns.md` with a one-line summary + link. *(low effort, high clarity)*
- **P-D2 / P-D3** — Pick one home per fact: keep `components/README.md` as the **only** component index, and slim `03-components.md` to a short narrative overview (families, shared base group) that links into `components/*.md` for detail — removing the duplicated per-component tables. Alternatively, deprecate `03-components.md` and redirect it. *(medium effort)*
- **P-D4** — State the `VenueList` orphan **once** in R1; have the other three spots link to R1 with a single clause. *(low)*
- **P-D5** — Keep the color-literal list only in R5; `01-foundations.md` should reference R5 rather than re-list. *(low)*
- **P-D6** — Mark `ROADMAP.md §5` as a historical snapshot and point it at the live register, or trim it to the IDs only. *(low)*

### A2. Missing documentation

**Findings** — styled, in-use classes with **no** dedicated coverage, and doc types absent:

| # | Gap |
|---|-----|
| M1 | **Feedback text family** — `.statusText`, `.errorText`, `.emptyText`, `.detailStatus`, `.routeHint(.error)` (`globals.css:545-556, 1356-1361, 1484-1493`) are only mentioned inside patterns/detail; no component doc |
| M2 | **Section Header** — `.sectionHeader` + its `h2`/`button` (`:562-584`) is a reusable pattern with no doc |
| M3 | **Brand / location header** — `.brandBlock`, `.brandName`, `.locationHint`, `.guestBadge` grouping (`:168-238`) is undocumented as a unit |
| M4 | **Auth CTA** — `.authCta` (`:679-705`) appears in the catalog but has no per-component doc despite the template set covering smaller pieces |
| M5 | **Category / dish pills** — `.category`, `.dish` (`:1363-1377`) are documented only inside the detail panel, though they're a reusable pill akin to tags |
| M6 | **Light theme** is documented only in `01-foundations.md`; no per-component light-mode notes or a dedicated theming doc |
| M7 | **Interaction/focus states** — no doc covers focus, hover, disabled, and the (absent) `:focus-visible` treatment across components |
| M8 | **Contributing guide** — no "how to add a component / class-naming rules / where things live" doc |
| M9 | **Visual reference** — no screenshots/rendered examples (acknowledged out of scope in the roadmap, but still a gap for a design system) |

**Proposals**

- **P-M1** — Add `components/feedback-text.md` (status/error/empty/route-hint) using the standard template. *(low)*
- **P-M2 / P-M4 / P-M5** — Add short docs for Section Header, Auth CTA, and Category/Dish pills (or fold the pills into `chip-tag.md` as a "static pill" variant to avoid a thin doc). *(low–medium)*
- **P-M3** — Add a "Brand & Panel Header" doc, or a lightweight "left-panel composition" doc. *(low)*
- **P-M6** — Add a `theming.md` covering the `data-theme` mechanism end-to-end and per-surface light overrides; link from each component's Accessibility/Usage. *(medium)*
- **P-M7** — Add a `states.md` (hover/active/selected/disabled/focus) and standardize a `:focus-visible` recommendation. *(medium)*
- **P-M8** — Add `CONTRIBUTING.md` for the design system (naming conventions, token usage, adding a glyph). *(low)*
- **P-M9** — When tooling allows, capture rendered examples per component (out of current read-only scope). *(high)*

### A3. Missing components (gaps in the *system*, not just the docs)

**Findings** — patterns a discovery UI usually needs that don't exist in code today:

| # | Missing capability | Evidence |
|---|--------------------|----------|
| C1 | **Loading / skeleton** states | Loading is only a label swap ("Đang tìm") + a text line; no skeletons/spinners (`DiscoveryExperience.tsx:535, 1115`) |
| C2 | **Focus-visible treatment** | No custom focus style anywhere in `globals.css`; keyboard focus relies on UA default |
| C3 | **Dialog / modal** | Auth ("Đăng nhập ngay") and menu CTA are inert buttons (R10); no modal primitive |
| C4 | **Toast / inline-error surface** | Errors are ad-hoc text nodes; no reusable notification component |
| C5 | **Pagination / "load more"** | `limit:20` is fixed (`:253`); no way to page results |
| C6 | **Illustrated empty state** | Empty = a single `.emptyText` line; no empty-state component |
| C7 | **Creator/avatar chip** | `SocialVideo.creator_handle` exists (`lib/api.ts:31`) but is never rendered |
| C8 | **Tooltip** | Icon-only controls have `aria-label` but no visible tooltip affordance |

**Proposals**

- **P-C1** — Introduce a skeleton/spinner primitive and an `aria-busy` convention for the three loading paths (search/detail/route). *(medium)*
- **P-C2** — Add a single shared `:focus-visible` outline token + rule; highest-value accessibility win. *(low, high impact)*
- **P-C3/P-C4** — Add Dialog and Toast primitives before wiring the currently-inert auth/menu/share/save actions (R10). *(medium–high)*
- **P-C5** — Add a "load more"/pagination control tied to `limit`. *(medium)*
- **P-C6/P-C7/P-C8** — Backlog: empty-state illustration, creator chip (data already present), tooltip. *(low–medium)*

---

## B. Design-system consistency review

### B1. Inconsistent spacing (relates to [R6](05-inconsistencies.md#r6--no-spacing-or-radius-scale-tokens))

**Findings** (from `globals.css`)

- **14 distinct `gap` values**: `2,3,4,5,6,7,8,10,12,14,16,18,22,26px` — no scale token.
  Small gaps cluster indistinguishably (`2/3/4/5/6/7/8` all in use).
- **~25 distinct `padding` shorthands**, e.g. `3px 7px`, `4px 8px`, `5px 10px 5px 6px`,
  `6px 10px`, `6px 10px 6px 6px` — several near-identical.
- **Panel gutter is inconsistent**: `margin: 0 24px` on most left-panel sections but the
  brand block uses `30px 32px 16px`; mobile drops all to `16px`.
- **6 distinct non-pill radii** (`8,10,12,14,16,18px`) plus `999px` (17×) and `inherit` (2×).

**Proposals**

- **P-S1** — Define a spacing scale as tokens (e.g. `--space-1…8` → `4,8,12,16,20,24,28,32`) and a radius scale (`--radius-sm/md/lg/pill`), then map existing values to the nearest step in a follow-up. *(medium; unlocks consistency everywhere)*
- **P-S2** — Consolidate the near-duplicate small gaps (`2/3/4/5/6/7`) toward `4/8`. *(low, after P-S1)*
- **P-S3** — Standardize the panel gutter to one token so brand/section/mobile agree. *(low)*

### B2. Inconsistent typography (relates to [R4](05-inconsistencies.md#r4--font-weights-used-that-arent-loaded))

**Findings**

- **14 distinct font sizes** (`10–30px`) with no scale; `11/12/13px` are near-identical
  and all in heavy use (12px ×14, 13px ×6, 11px ×6).
- **`h2` is redefined three ways** on the same tag: base `16px/800`
  (`:314`), `.sectionHeader h2` `12px` uppercase (`:570`), `.detailTitleRow h2` `30px`
  (`:1237`) — semantic heading level decoupled from visual size.
- **Weights out of the loaded set**: `900` (×6) and `850` (×1) aren't imported
  (Be Vietnam Pro `400/500/700/800`, Plus Jakarta `600/700/800`) → silent fallback (R4).
  `700` is used only twice, against an `800` norm (×22).
- `500`/`600` weights are **loaded but never used**.
- **7 ad-hoc line-heights** (`1, 1.14, 1.18, 1.35, 1.45, 1.5, 1.55`).

**Proposals**

- **P-T1** — Define a type scale as tokens (size + line-height + weight per role: display/title/heading/body/caption/label) and map the 14 sizes onto ~6 steps. *(medium)*
- **P-T2** — Decouple heading semantics from size: use a role class (e.g. `.text-title`) rather than restyling `h2` three ways, so heading level stays semantic. *(medium)*
- **P-T3** — Either load `850/900` (adjust the `@import`) or stop using them; drop unused `500/600` from the import if they stay unused. *(low; resolves R4)*
- **P-T4** — Replace the lone `700` usages with the `800` norm (or codify when `700` is intentional). *(low)*

### B3. Inconsistent colors (relates to [R5](05-inconsistencies.md#r5--blue-accent-has-no-token-warm-colors-re-expressed-as-literals) / [R11](05-inconsistencies.md#secondary-findings))

**Findings**

- **`--primary` is re-expressed as 17 raw `rgba(255,183,127,…)` literals across 11 alpha
  levels** (`0.08 … 0.68`) — the token exists but is bypassed for every translucent
  tint. Same pattern for `--secondary` (`rgba(254,183,0,…)` ×2) and `--trend-fire`
  (`rgba(255,77,0,…)` ×4).
- **Un-tokenized accents**: blue `#2f80ff` (×3 CSS + ×1 TSX), button text `#4e2600` (×3),
  `#2f271f` (×2), `#2f1500`, `#080808`, and **`#ffffff` ×10** (no white token).
- **Declared-but-unused tokens**: `--primary-strong` (**0** refs) and `--surface-mid`
  (**0** refs). Gradients use literal `#ff8a00` instead of `--primary-strong`.
- The **map basemap palette** (`applyMapLayerPreferences`) is a separate hardcoded color
  set, disconnected from tokens.

**Proposals**

- **P-C-1** — Introduce alpha-tint tokens (or a `color-mix(...)`/rgba-from-var pattern)
  so translucent primaries reference `--primary` instead of hardcoding `255,183,127`.
  *(medium; eliminates 17+ literals)*
- **P-C-2** — Add tokens for the blue accent, `--on-primary` (currently `#4e2600`), and a
  white/`--on-inverse` token; replace literals. *(medium; resolves R5)*
- **P-C-3** — Either use `--primary-strong`/`--surface-mid` or remove them (R11). *(low)*
- **P-C-4** — Extract the map palette into a named, theme-aware config so it lives beside
  the token system. *(medium)*

### B4. Inconsistent responsive behavior

**Findings**

- **Only 2 breakpoints** (`max-width: 1180px`, `640px`), desktop-first; no tablet range
  and no `min-width`/large-screen handling.
- **Abrupt layout jump**: left panel is a fixed `420px` column until 1180px, then the
  whole grid collapses to one column — nothing intermediate.
- **Uncapped large screens**: the map column is `minmax(0,1fr)` with no max-width; panels
  stay `420px`, so the map grows unbounded on ultra-wide displays.
- **Selective collapse is inconsistent**: at 640px, `.detailActions` collapses to 1 column
  but the sibling 3-column `.compactDetailActions` (`:1251`) and the 2-column `.videoGrid`
  (`:1417`) **do not** — same context, different behavior.
- **Magic numbers**: `.mapStage { min-height: 620px }` at 1180px; markers bump
  `min-width` `138→150px` only at 640px.
- **Scroll model flips** globally: `body{overflow:hidden}` → `auto` at 1180px.

**Proposals**

- **P-R1** — Define named breakpoint tokens and add an intermediate tablet range so the
  1180px collapse isn't a single hard cliff. *(medium)*
- **P-R2** — Cap the content width on large screens (max-width on the map column or the
  shell) to avoid an unbounded map. *(low)*
- **P-R3** — Make sibling grids collapse consistently at 640px (include
  `.compactDetailActions` / `.videoGrid` or document why they differ). *(low)*
- **P-R4** — Replace responsive magic numbers with tokens/derived values and document the
  scroll-model switch. *(low)*

---

## Prioritized summary

Ordered by value ÷ effort. "Type" = Doc or Code (code items are proposals for a future,
separately-approved change; this audit changes neither).

| Priority | Proposal | Type | Why first |
|----------|----------|------|-----------|
| 1 | **P-C2** shared `:focus-visible` | Code | Biggest accessibility gap, tiny change |
| 2 | **P-D1/P-D4/P-D5** de-duplicate glyph table, orphan note, color list | Doc | Cheap; removes maintenance drift |
| 3 | **P-T3** fix/remove unloaded `850/900` weights (R4) | Code | Removes silent rendering fallback |
| 4 | **P-C-3** use or drop unused tokens (R11) | Code | Trivial; dead tokens mislead |
| 5 | **P-M1/P-M7** feedback-text + states docs | Doc | Closes the clearest doc gaps |
| 6 | **P-S1 / P-T1** spacing & type scale tokens | Code | Foundational; unblocks most other consistency fixes |
| 7 | **P-C-1/P-C-2** tokenize rgba primaries + blue/white (R5) | Code | Large literal cleanup once scales exist |
| 8 | **P-D2/P-D3** single component index; slim catalog | Doc | Medium effort; structural |
| 9 | **P-R1…R4** responsive consistency | Code | Needs design input on breakpoints |
| 10 | **P-C1/P-C3/P-C5** skeletons, dialog, pagination | Code | New components; larger builds |

### Explicitly not done here

- No `apps/web` source was changed.
- No existing `docs/design-system/**` file was edited to fix the duplications above — they
  are proposed (P-D*), not applied.
- The only new artifact is this file (+ a one-line index link from the design-system README
  for discoverability).
