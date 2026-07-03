# ADR-0006: Bilingual handbook (VN prose + EN technical)

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision made in the handbook design spec, 2026-07-01)
- **Deciders:** Handbook author + user (blocks approved in the design spec)
- **Type:** retroactive (documents an existing decision)

## Context

`docs/handbook/` is human onboarding for a Vietnamese-speaking team, but it must stay
usable by AI agents and consistent with the rest of the repo, which keeps technical
artifacts (prompts, specs, code, file/concept names) in English.

## Decision

In `docs/handbook/**`: write **explanatory prose and section headings in Vietnamese**,
and keep **technical content in English** — prompts, specs, exit-criteria, code, file
names, and concept names.

Scope: this bilingual convention applies to the **handbook only**. Reference material
outside the handbook — `docs/standards/**`, `docs/adr/**`, `docs/services/**`,
`docs/ai/**`, and other standard/contract docs — is **English-only**, matching existing
`docs/standards/` and keeping reference docs uniformly reusable.

## Consequences

- **Positive:** approachable onboarding for the team; technical terms stay copy-pasteable
  and stable for AI reuse; no translation drift on rules/prompts.
- **Negative / cost:** two registers in one file; contributors must know which parts stay
  English; non-VN readers rely on the English technical anchors.
- **Follow-ups:** new handbook parts (3,5,6,7,8) follow this convention; new reference
  docs do not (they are English-only) — this split is recorded in the AI-OS roadmap.

## Alternatives considered

- **English-only everywhere** — rejected for the handbook: less accessible for team
  onboarding.
- **Fully bilingual everywhere (incl. standards/ADRs)** — rejected: heavier to maintain
  and diverges from `docs/standards/`.

## Sources / links

- [`docs/superpowers/specs/2026-07-01-handbook-design.md`](../superpowers/specs/2026-07-01-handbook-design.md) → *Quy ước*.
- [`docs/handbook/README.md`](../handbook/README.md) → *Quy ước* (language convention).
- [`docs/ai/ROADMAP.md`](../ai/ROADMAP.md) → *Guiding principles* (English reference docs).
