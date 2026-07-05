# Playbook: Create a Table

> **Target is Proposed** — no table component exists ([`../components/table.md`](../components/table.md)).
> The shipped analogues are the detail-panel **row lists** (`.dishDetailItem`,
> `.openingRow`). Prefer extending those before introducing a real `<table>`.

## When to use

Displaying small, read-only structured data (opening hours, dish price lists). For a
heavy interactive data-grid — you almost certainly don't need one here (discovery data is
small and read-only).

## Steps

1. **Prefer the existing rows** for key→value data: `.openingRow`/`.dishDetailItem`
   (`globals.css:1386-1415`) are grid `minmax(0,1fr) auto` with outline/glass/radius tokens.
2. **Use a semantic `<table>`** only when data is genuinely tabular (multiple columns with
   headers). Style cells to match the row language (outline borders, 13px body, 12px muted,
   right-aligned numerics).
3. **Headers**: `<th scope="row"|"col">`; add `<caption>`/`aria-label`.
4. **Data**: real API shapes (`VenueDish`, `OpeningHour`); format via `formatPrice`/`formatDay`.
5. **Responsive**: ensure the table/rows reflow or scroll within their container at 640px.
6. **Document** in `table.md`; flip its banner to As-built when built.

## Checklist

- [ ] Reused row styling; semantic `<table>` only when truly tabular.
- [ ] `scope` headers + caption/label; numerics right-aligned.
- [ ] Real API data; localized formatting.
- [ ] Reflows/scrolls at 640px.
- [ ] `table.md` updated; meets [DoD](../../verification/definition-of-done.md).

## Related

- [../components/table.md](../components/table.md) · [../pages/restaurant.md](../pages/restaurant.md) · [create-component.md](create-component.md)

## References

- `apps/web/app/globals.css` (`.dishDetailItem`/`.openingRow` `:1386-1415`) · `apps/web/lib/api.ts` (`:39-61`).
