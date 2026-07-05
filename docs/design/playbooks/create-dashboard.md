# Playbook: Create a Dashboard

> **Target is Proposed, and scope-sensitive.** No dashboard exists, and there is **no
> dashboard component/page doc** — dashboards aren't part of the current discovery
> product. Read this before assuming one belongs here.

## Scope caution (read first)

TikFood is a consumer **discovery** app (map-first, dish-first, discovery-only). A
user-facing analytics/creator/merchant dashboard risks crossing anti-goals
(creator monetization, merchant tooling). **Confirm the audience and get human approval**
before building any dashboard-like surface (`CLAUDE.md`).

## When it might be legitimate

- An **internal/ops** view (not shipped to consumers) — likely lives outside `apps/web`.
- A small, discovery-flavored **stats strip** (e.g. "N venues in area", trend distribution)
  — which is really just cards/overlays, not a true dashboard.

## Steps (if approved)

1. **Confirm audience + approval.** If consumer-facing and analytics/creator-oriented →
   stop; it's likely an anti-goal.
2. **Compose from existing pieces**: [cards](../components/card.md) for metrics, the
   [table](create-table.md) rows for breakdowns, the [map](../components/map.md) for geo.
   Do **not** add a charting library without justification and approval.
3. **Layout** with the existing grid tokens; keep the dark-glass language.
4. **Data** through `lib/api.ts`; no raw `fetch`.
5. **Document** as a new page/pattern only after it's real.

## Checklist

- [ ] Audience confirmed; human approval obtained; not an anti-goal surface.
- [ ] Built from existing cards/rows/map; no unjustified charting dependency.
- [ ] Grid/tokens reused; data via `lib/api.ts`.
- [ ] Meets [DoD](../../verification/definition-of-done.md).

## Related

- [create-card.md](create-card.md) · [create-table.md](create-table.md) · [create-map.md](create-map.md) · [../01-design-principles.md](../01-design-principles.md)

## References

- `CLAUDE.md` (anti-goals, approvals) · `docs/tikfood/anti-goals.md`.
