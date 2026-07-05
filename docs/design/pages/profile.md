# Page: Profile

> **Status: Proposed — not implemented, and anti-goal-constrained.** There is no profile,
> no auth, and no user model today. The UI runs as **"Guest View"** (`.guestBadge`), and
> "Đăng nhập ngay" is an inert button. Any profile work must respect the product
> anti-goals. This is a cautious forward proposal.

## Anti-goal guardrails (read first)

`CLAUDE.md` blocks: social follow graph, creator monetization, in-app chat, and all
commerce (cart/order/checkout/payment/booking/reservation). A profile must therefore
**not** become a social/creator/commerce surface. Require human approval before building
anything auth-related.

## Purpose (proposed, minimal)

The only anti-goal-safe scope is a lightweight **saved/bookmarked venues** view for a
signed-in guest — matching the existing (disabled) bookmark affordances
(`miniTab` bookmark is `disabled`; detail "Lưu" is inert).

## Usage (proposed)

- Route `app/profile/page.tsx` (Proposed), gated behind auth (**requires approval**).
- Reuse [Venue rail card](../components/card.md) to list saved venues; reuse
  [empty-state](../patterns/empty-state.md) when none are saved.
- Sign-in likely surfaces via the Proposed [modal](../components/modal.md).

## Rules (proposed)

- Discovery-only: saved venues + basic guest prefs (e.g. theme) — nothing social/commerce.
- Auth, and any external identity/network calls, need human approval (`CLAUDE.md`).
- Reuse existing cards/empty-state; no new visual language.

## Do / Don't

- **Do** keep scope to saved places / preferences.
- **Do** flag auth work for human approval before implementing.
- **Don't** add followers, DMs, creator tools, or purchases.
- **Don't** ship the currently-inert "Lưu"/bookmark as real without a persistence + auth plan.

## Accessibility (proposed)

- Standard landmark/labeling; the saved list reuses the accessible rail-card pattern.

## Related

- [components/card.md](../components/card.md) · [components/modal.md](../components/modal.md) · [patterns/empty-state.md](../patterns/empty-state.md) · [Brand — Guest View](../02-brand.md)

## Files

- **None yet** (Proposed). Analogues: `.guestBadge` (`globals.css:224`), inert bookmark/save (`DiscoveryExperience.tsx:580, 1081`).

## References

- `CLAUDE.md` (anti-goals, security/approval) · `docs/tikfood/anti-goals.md`.

## Future improvements

- N/A until a human-approved auth + persistence design exists.
