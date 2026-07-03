# ADR-0005: Plain CSS for `apps/web`, no framework yet

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision predates this ADR)
- **Deciders:** Frontend owners
- **Type:** retroactive (documents an existing decision)

## Context

`apps/web` is an MVP discovery UI. A styling approach was needed. Adding a CSS
framework or design-system library early would introduce dependencies and conventions
before the UI's real needs are known.

## Decision

Use **plain CSS** in `app/globals.css` with reused class patterns. No CSS framework and
no design-system library. Inline `style` is allowed only for computed values (e.g. map
marker positions). A new CSS framework requires justification.

Note: `packages/config/tikfood.ai-agent.yaml` mentions Tailwind/shadcn as a *target*
stack for a future product repo; the **current** `apps/web` intentionally does not use
them — code wins over aspirational config.

## Consequences

- **Positive:** minimal dependencies; full control over styles; nothing to learn beyond
  CSS; fast MVP iteration.
- **Negative / cost:** `globals.css` is large and centralized; no utility-class
  ergonomics or design tokens; consistency relies on reusing existing classes.
- **Follow-ups:** revisit (with a superseding ADR) if/when a design system, theming, or
  component-library needs emerge. The approved split of `DiscoveryExperience.tsx` into
  focused components is a prerequisite for cleaner style ownership.

## Alternatives considered

- **Tailwind CSS** — deferred: valuable at scale, but adds tooling/conventions before
  the UI justifies them; remains the likely future choice per product config.
- **CSS Modules / styled-components** — not adopted; plain CSS is sufficient for MVP.

## Sources / links

- [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) → *Stack Reality* / *Styling* / *Future Direction*.
- [`docs/standards/frontend-architecture.md`](../standards/frontend-architecture.md).
