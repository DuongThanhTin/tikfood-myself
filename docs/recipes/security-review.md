# Recipe — Security review

**Goal:** assess a change (or an area) for security issues before it ships.

**When to use:** reviewing a change that touches input handling, data exposure, external
calls, secrets, dependencies, or any gated area.

## Context to load
- [`docs/standards/application-security.md`](../standards/application-security.md) (product
  app) · [`docs/security.md`](../security.md) (runner/automation).
- [`docs/verification/review-guide.md`](../verification/review-guide.md) (lens + severity).
- [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) (hard rules / gates).

## Process skill
Review discipline (`superpowers:requesting-code-review`) with the security lens.

## Steps
1. **Map the attack surface** — what untrusted input enters, what data leaves, what
   external systems are touched.
2. **Input** — validated and bounded at the boundary? SQL parameterized (no string
   building)? Untrusted text safely rendered on the frontend?
3. **Output** — no SQL/stack/secret/internal detail leaked in responses or logs? Only
   necessary fields returned? No fabricated-as-real data?
4. **Secrets** — sourced from env only; never logged; no `.env*` reads; no keys in the
   client bundle (except intentional `NEXT_PUBLIC_*`).
5. **Gated areas** — auth, migrations, infra, new external network calls, model-cost
   changes → **human approval obtained?** If not, block.
6. **Transport/abuse** — HTTPS assumed; CORS restricted; expensive queries bounded; rate
   limiting considered for public endpoints.
7. **Dependencies** — new deps justified and checked for known vulnerabilities.
8. **Privacy/compliance** — social ingestion respects platform terms; minimal personal
   data.

## Verification
Run the [security checklist](../standards/application-security.md#checklist-for-a-security-relevant-change);
findings ranked by the [severity rubric](../verification/review-guide.md#severity-rubric)
with file:line and a concrete fix.

## Exit
Surface mapped, checklist passed or findings filed by severity, gated items approved or
blocked, no secrets/leaks.

## Common mistakes
- Reviewing code style but missing the trust boundary.
- Assuming internal input is safe.
- Letting a gated change through without explicit approval.
- Treating the runner security doc as covering the product app (it doesn't — use both).
