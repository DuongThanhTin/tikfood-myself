# Playbook: Review UI

> Review a UI change against the design system before merge. Pairs with the repo's
> [`review-guide`](../../verification/review-guide.md) and
> [`definition-of-done`](../../verification/definition-of-done.md).

## When to use

Reviewing any PR that touches `apps/web` visuals, components, or pages.

## What to check

### Consistency with the system
- [ ] Colors reference tokens ([`03`](../03-color-system.md)) — no new `#hex`/`rgba(255,183,127,…)` literals where a token exists.
- [ ] Spacing/radii reuse existing values ([`04`](../04-spacing-system.md)); no needless new magnitudes.
- [ ] Type uses the established families/sizes/weights ([`06`](../06-typography.md)); **no `font-weight` >800** (not loaded).
- [ ] Motion reuses the standard durations/easing ([`08`](../08-motion.md)); no gratuitous keyframes.

### Reuse & structure
- [ ] Reuses existing components/classes; no near-duplicate of an existing piece.
- [ ] Descriptive names; semantic markup; inline `style` only for computed values.
- [ ] Data flows through `lib/api.ts`; no raw `fetch` in components.

### Accessibility ([`09`](../09-accessibility.md))
- [ ] Icon-only controls have `aria-label`; decoration is `aria-hidden`; images have `alt`.
- [ ] State conveyed beyond color (`aria-pressed`/`aria-current` where relevant).
- [ ] Focus is visible; no enabled-but-inert buttons.

### Responsive & theme
- [ ] Works at both breakpoints (1180/640) and in light + dark.
- [ ] Sibling grids collapse consistently (watch the known 640px inconsistency).

### Product guardrails
- [ ] Vietnamese UI copy; English identifiers.
- [ ] No anti-goal surface (cart/order/checkout/payment/booking/chat/follow/monetization/livestream).

## Output

Cite `file:line`. Separate **must-fix** (broken a11y, anti-goal, breaking API/response
change, new literal where a token exists) from **nice-to-have** (token/scale cleanup).
Prefer linking a foundations rule over restating it.

## Related

- [../../verification/review-guide.md](../../verification/review-guide.md) · [../../verification/definition-of-done.md](../../verification/definition-of-done.md) · [create-component.md](create-component.md) · [update-component.md](update-component.md)

## References

- `apps/web/CLAUDE.md`.
