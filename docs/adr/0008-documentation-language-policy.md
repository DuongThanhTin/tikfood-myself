# ADR-0008: Documentation language policy (one canonical language per doc)

- **Status:** Accepted
- **Date:** 2026-07-06
- **Deciders:** tin.duong (user) + AI agent
- **Type:** decision
- **Supersedes:** ADR-0006

## Context

Documentation language was drifting. `docs/handbook/` had a bilingual convention
([ADR-0006](0006-bilingual-handbook.md)) and `docs/standards|adr|services|ai` were
English-only, but the rest of `docs/` was unspecified. This produced two problems:

- A **parallel English translation** was maintained under `docs/getting-started/en/`
  (7 files mirroring the Vietnamese originals) — exactly the kind of duplicate that
  drifts and doubles maintenance.
- New machine-facing docs (`docs/workflows/`) were written in Vietnamese, inconsistent
  with the English reference corpus agents consume.

ADR-0006 only scoped the handbook; there was no repo-wide rule stating *which single
language each doc uses* and *that no doc keeps a parallel translation*.

## Decision

**Each document has exactly one canonical language. We do not maintain parallel
translations of the same document.**

- **Machine / agent-facing docs → English.** Applies to `docs/ai/`, `docs/recipes/`,
  `docs/workflows/`, `docs/standards/`, `docs/adr/`, `docs/verification/`,
  `docs/services/`, `docs/tikfood/`, and — as the repo-wide default for any
  otherwise-unclassified reference area — `docs/thinking/`, `docs/design/`,
  `docs/agents/`, `docs/features/`, and top-level `docs/*.md`.
- **Human onboarding docs → Vietnamese.** Applies to `docs/getting-started/` and
  `docs/handbook/`. The handbook keeps the ADR-0006 register — **Vietnamese prose,
  English technical terms** (prompts, code, file/concept names, exit-criteria); that
  rule is preserved here, not revoked.
- **Language ≠ content.** Only prose is in the canonical language. Code, identifiers,
  file paths, config keys, and **sample UI copy** (e.g. Vietnamese strings a screen
  actually renders) stay verbatim regardless of the doc's language.
- **Exemptions:** `docs/archive/**` (frozen history) and dated working artifacts under
  `docs/superpowers/specs/` and `docs/superpowers/plans/` are **not** retranslated —
  they are point-in-time records, near-archival.

## Consequences

- **Positive:** no translation drift; a single source of truth per doc; the English
  reference corpus stays uniformly reusable by agents; Vietnamese onboarding stays
  approachable for the team.
- **Negative / cost:** non-Vietnamese readers rely on English for everything except
  onboarding; the Vietnamese team reads machine docs in English (mitigated by the
  handbook's bilingual register for learning).
- **Follow-ups (this ADR's rollout):** delete `docs/getting-started/en/` and fix its
  one inbound link; translate `docs/workflows/*` to English (and restructure for
  readability); record this policy in [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md).

## Alternatives considered

- **Fully bilingual everywhere (two files or two registers per doc)** — rejected:
  heaviest to maintain, guarantees drift; this is the problem we are removing.
- **English-only everywhere (incl. onboarding)** — rejected: less accessible for team
  onboarding, the reason ADR-0006 chose a Vietnamese register for the handbook.
- **Leave ADR-0006 as-is and add nothing** — rejected: it left most of `docs/`
  unspecified and did not forbid parallel translations.

## Sources / links

- [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) → *Documentation language* (index entry).
- [ADR-0006](0006-bilingual-handbook.md) — superseded; its handbook register is retained above.
- [`/CLAUDE.md`](../../CLAUDE.md) · [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md).
