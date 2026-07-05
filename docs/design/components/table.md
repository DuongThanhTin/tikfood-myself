# Table

> **Status: Proposed — not implemented.** No table component exists in `apps/web`. The
> closest shipped patterns are the **row lists** in the detail panel (`.dishDetailItem`,
> `.openingRow`) — see [`../pages/restaurant.md`](../pages/restaurant.md). This spec is a
> forward proposal grounded in those existing row styles.

## Purpose

Display structured, columnar data. The realistic near-term needs are **opening hours**
and **dish price lists** — both currently rendered as two-column *rows*, not a semantic
table.

## Usage (proposed)

- Prefer semantic `<table>` for truly tabular data (headers + rows), styled to match the
  existing row language: `1px var(--outline)` cells, `border-radius:12px` container,
  `--glass` row background, 13px body / 12px muted secondary.
- For simple key→value pairs (like opening hours), the current
  `.openingRow` grid (`minmax(0,1fr) auto`) is adequate and may be kept.

## Examples (proposed shape)

```jsx
<table className="dataTable" aria-label="Giờ mở cửa">
  <tbody>
    {hours.map((h) => (
      <tr key={h.day_of_week}><th scope="row">{formatDay(h.day_of_week)}</th><td>{h.is_closed ? "Đóng cửa" : `${h.open_time}–${h.close_time}`}</td></tr>
    ))}
  </tbody>
</table>
```

## Rules (proposed)

- Use `<th scope>` for row/column headers; don't fake a table with `<div>`s when the data is tabular.
- Reuse the outline/glass/radius tokens from the existing row items.
- Right-align numerics (prices) as `.dishDetailItem`/`.openingRow` already do.

## Do / Don't

- **Do** reuse the current row styling; extend, don't reinvent.
- **Do** keep dish/hour data from the real API shapes (`VenueDish`, `OpeningHour`).
- **Don't** introduce a heavy data-grid; discovery data is small and read-only.
- **Don't** add sorting/pagination inside the detail panel (out of scope).

## Accessibility (proposed)

- Semantic table with `aria-label`/`<caption>`; `scope` on headers; no layout-only tables.

## Related components

- [Card](card.md) · [pages/restaurant.md](../pages/restaurant.md) (existing row lists) · [create-table playbook](../playbooks/create-table.md)

## Files

- **None yet** (Proposed). Existing analogues: `app/globals.css` `.dishDetailItem`/`.openingRow` `:1386-1415`.

## References

- `apps/web/lib/api.ts` — `VenueDish` `:39-54`, `OpeningHour` `:56-61`.

## Future improvements

- N/A until implemented. Decide per case whether a semantic `<table>` or the existing rows fit.
