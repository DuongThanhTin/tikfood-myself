# Playbook: Update a Component

> Change an existing component without breaking its callers or the shared base styles.

## When to use

Fixing, restyling, or extending a component that already ships.

## Prerequisites

- Read the component's doc under [`../components/`](../components) and the relevant
  foundations.
- Note that many classes are **shared groups** (e.g. the button base, badge base, glass
  base) — a change ripples to every user.

## Steps

1. **Map the blast radius.** Grep the class(es) across `globals.css` and
   `DiscoveryExperience.tsx`. Shared bases (`.primaryButton,…` `:505`; badge base `:211`;
   `.trendScoreCard, .glassAction` `:1297`) affect multiple components.
2. **Prefer a modifier** over changing the base — add `.thing--variant` rather than
   editing shared rules, unless the change is intended for all.
3. **Keep tokens.** If replacing a literal, prefer the matching token; if none exists,
   consider adding one (foundations) rather than spreading another literal.
4. **Preserve the response/data contract.** Don't change prop shapes or the `{ data,
   error }` handling; keep `lib/api.ts` types in sync with Go tags.
5. **Check responsive + both themes.** Verify the two breakpoints (1180/640) and
   `data-theme` light/dark still hold.
6. **Update the doc.** Reflect the change in the component doc; move any newly-noticed
   inconsistency into its Future Improvements.

## Checklist

- [ ] Blast radius mapped; shared-base edits are intentional.
- [ ] Modifier used where the change isn't global.
- [ ] Tokens preserved/added, not new literals.
- [ ] Both themes + both breakpoints verified.
- [ ] No breaking prop/response changes; tests updated if behavior changed.
- [ ] Component doc updated.
- [ ] Meets [`definition-of-done`](../../verification/definition-of-done.md).

## Related

- [create-component.md](create-component.md) · [review-ui.md](review-ui.md) · [../04-spacing-system.md](../04-spacing-system.md)

## References

- `apps/web/CLAUDE.md` (no breaking API changes; reuse; update tests).
