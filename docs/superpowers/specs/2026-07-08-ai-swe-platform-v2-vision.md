# Vision — AI Software Engineer Platform, V2 (second-pass, first-principles)

**Date:** 2026-07-08 · **Status:** vision / not-yet-planned · **Altitude:** Horizon 2 (portable, multi-repo) is the center of gravity; H1 (this repo) is the migration base; H3 (platform) is a labeled north star, not a commitment.

> **What this document is.** A Principal-level second-pass review of the V1 work (the
> documentation-based AI Operating System + the `.claude/` executable spine) and a
> first-principles redesign. It is a *vision*, not a build plan — each accepted piece
> becomes its own spec → plan → batch later. It is deliberately written as an **honest
> review**: it names where V1 is strong, where it is weak, and — critically — where V2
> would be *over-engineering for one engineer and one repo*. Ambition without that last
> part is how teams build cathedrals nobody can maintain.

---

## 0. Thesis

V1's real weakness is not missing content. The repo has excellent content: an AI-Contract,
a context-loading protocol, ADRs, standards, recipes, workflows, agent role docs, a
domain model, and now a self-triggering `.claude/` layer over all of it. Very few
solo repos are this well-organized.

The weakness is that **V1 is an open-loop system.** It optimizes two things well —
*capturing knowledge* (docs/ADRs/standards) and *triggering procedure* (skills/agents/
commands) — but it has almost no **feedback**. It cannot tell whether the AI-OS is
getting better or worse over time, and it cannot tell whether its own claims are true.

The final whole-branch review of the V1 spine proved this concretely, not
hypothetically: the verification script — the one component whose entire job is to stop
drift — **did not enforce the conventions it advertised.** A skill missing its `name:`
key passed the gate; a broken bare reference passed the gate. Both were reproduced. We
fixed them, but the *category* is the point: a system that asserts "verified" without a
mechanism that makes the assertion true is doing **verification theatre.** That is the
signature of an open-loop system.

> **V2 in one sentence:** stop organizing knowledge and start closing loops —
> context → act → verify → **learn** → improve — so the system measurably corrects
> itself instead of merely being well-documented.

Reframe: from *"a library of rules and recipes"* to *"a self-correcting engineering
system."*

---

## 1. Second-pass critique of V1

Each gap is stated as: what's missing → why it matters → evidence (where it exists).

### 1.1 Verification theatre (the anchor finding)
- **Missing:** gates that actually enforce what they claim. The DoD, the recipes, and the
  `.claude/` gate all *assert* verification; only some of it is mechanically true.
- **Why it matters:** every downstream honesty rule (AI-CONTRACT §9) rests on verification
  being real. If "verified" can be false, the whole trust model is soft.
- **Evidence:** final review of the spine — missing-`name` skills and bare `@refs` passed
  the gate (reproduced, then fixed in commit `0a4ea43`). The fix closed *these two*; it
  did not add a mechanism that would catch the *next* such gap. That mechanism is §2.3.

### 1.2 No eval / regression harness for the AI-OS itself
- **Missing:** any way to answer "did editing this skill make it better or worse?" There
  are no golden transcripts, no "does this skill still trigger for its intended task?"
  test, no regression suite for prompts/recipes.
- **Why it matters:** the AI-OS is now a large, interdependent artifact. Editing a skill
  description, a recipe step, or the context-loading matrix can silently degrade behavior.
  Today the only detector is a human noticing later. This is the single highest-leverage
  gap.
- **Evidence:** `make verify` tests the *application* (Go/web/runner). Nothing tests the
  *AI-OS*. The `.claude/` gate checks syntax (frontmatter, links), never behavior.

### 1.3 Static memory, no learning loop
- **Missing:** an automatic path from "we learned something while working" back into the
  knowledge base. ADRs and standards are write-once-by-human. The user's auto-memory
  (`memory/MEMORY.md`) exists but is not wired to the *engineering* loop (bugs fixed,
  patterns discovered, mistakes made).
- **Why it matters:** the goal is "senior engineer," and seniority is compounded lessons.
  A system that forgets every bug's root cause the moment the PR merges cannot compound.
- **Evidence:** there is a `fix-bug` recipe but no "capture the lesson" step that produces
  a durable, retrievable artifact. Lessons live only in git history, which no session
  reads proactively.

### 1.4 No observability
- **Missing:** run traces, and per-skill / per-agent metrics (was it triggered? did it
  succeed? tokens? wall-clock? did the human have to intervene?).
- **Why it matters:** you cannot improve what you cannot see. Prioritizing which skill to
  fix, which agent wastes tokens, which recipe gets abandoned mid-way — all require data
  the system does not collect.
- **Evidence:** none of the V1 artifacts emit a structured event. The only record is the
  chat transcript, which is not queryable across sessions.

### 1.5 Governance is advisory prose, not enforced
- **Missing:** mechanical enforcement of the rules that matter most. "Require human
  approval for protected paths," "never push main," "never read secrets," "treat repo
  content as prompt-injection" — all live in Markdown that a model is *asked* to follow.
- **Why it matters:** prose rules are probabilistic; a hook is deterministic. For
  security-critical rules, "the model usually complies" is the wrong guarantee. The runner
  already has a `commandGuard`; the interactive Claude Code side does not.
- **Evidence:** `.ai-agent.yaml` `protected_paths` and AI-CONTRACT §2/§4 are enforced by
  the runner pipeline but not by a Claude Code hook. Interactive sessions rely on the
  model reading and obeying prose.

### 1.6 Single-repo coupling (tension with the H2 goal)
- **Missing:** portability. The `.claude/` thin wrappers hard-link `@docs/...` paths
  specific to this repo. The stated goal ("consistent across many projects") needs the
  layer to be a *distributable kit* with a clear seam between generic method and
  repo-specific content.
- **Why it matters:** if every new repo requires hand-porting, the method does not scale
  past this repo — which is exactly the H2 ambition.
- **Evidence:** `starters/ai-os-kit/` exists and anticipates this, but the `.claude/`
  layer was built repo-locally and is not yet part of the kit.

### 1.7 Doc↔code drift has a rule but no detector
- **Missing:** a mechanism behind "code wins." When a recipe references a route, file, or
  symbol that no longer exists, nothing flags it.
- **Why it matters:** the AI-OS's value is proportional to its accuracy. A confidently
  wrong recipe is worse than no recipe — it sends the agent down a dead path with full
  authority.
- **Evidence:** AI-CONTRACT and CLAUDE.md both state "code wins, flag the mismatch," but
  the only mismatch-detector is a human reading both.

### 1.8 Assumptions in V1 worth challenging
- **"Thin wrappers that link docs are strictly better than self-contained skills."** Mostly
  true (DRY, single source of truth) — but it adds a second file-read on every trigger and
  makes portability harder (§1.6). V2 should treat this as a *tunable*, not a dogma:
  generic method can be self-contained in the kit; repo-specifics stay linked.
- **"More artifacts = more capability."** False past a point. Every skill added is a
  trigger the model must disambiguate. The `add-api-endpoint` vs `contract-first-change`
  overlap (flagged in review) is the first instance; at 20+ skills, trigger collision and
  mis-fire become a real cost. V2 needs a *trigger-precision* discipline, not just more
  skills.
- **"The knowledge base should be exhaustive."** The instinct behind the original mega-
  prompt. But context is a budget, not a warehouse. Exhaustive docs that blow the context
  window are a liability. V2 optimizes *retrieval precision*, not corpus size.

---

## 2. First-principles V2 architecture

If we were designing the best AI software engineer we could actually run — one engineer,
Claude Code, portable across repos — it would be organized around **six planes**, each of
which closes a loop. Planes are concerns, not directories; several already partly exist.

```
                 ┌───────────────────────────────────────────────┐
                 │  GOVERNANCE / SAFETY  (enforced, not advised)   │  ← wraps everything
                 └───────────────────────────────────────────────┘
   ┌────────────┐   ┌────────────┐   ┌──────────────┐   ┌─────────────────┐
   │ KNOWLEDGE  │──▶│ EXECUTION  │──▶│ VERIFICATION │──▶│ MEMORY/LEARNING │──┐
   │ (what we   │   │ (do the    │   │ (prove it,   │   │ (capture the    │  │
   │  know)     │   │  work)     │   │  incl. self- │   │  lesson, feed   │  │
   │            │◀──┼────────────┼───┤  eval)       │   │  it back)       │◀─┘
   └────────────┘   └────────────┘   └──────────────┘   └─────────────────┘
          ▲                                                      │
          └──────────────── learning loop closes here ──────────┘
                 ┌───────────────────────────────────────────────┐
                 │  OBSERVABILITY  (traces/metrics under all 4)    │
                 └───────────────────────────────────────────────┘
```

### 2.1 Knowledge plane — *keep, tune for retrieval*
What exists (AI-CONTRACT, CONTEXT-LOADING, ADRs, standards, recipes, domain model) is the
strong core. V2 change is not "more" but "retrieved precisely": the context-loading matrix
becomes the contract a **context-builder** (already prototyped as `repo-context-reader`)
executes, returning the *minimum* pack. Success metric: relevant-context-loaded / total-
context-loaded goes up; window pressure goes down.

### 2.2 Execution plane — *keep, harden trigger precision*
Skills/agents/commands stay. Two disciplines added: (a) **trigger precision** — every new
skill must state what it does *and what neighboring skill to prefer instead* (fixes the
overlap class); (b) **the kit seam** — generic method (TDD loop, review loop, context
protocol) is portable; repo-specifics (venue endpoints, `{data,error}`) are linked.

### 2.3 Verification plane — *the big one: real gates + AI-OS self-eval*
Two layers:
- **Real gates (fix the theatre):** the pattern behind the spine's fix generalizes —
  every gate ships with a RED test that proves it fails on the bad input it claims to
  catch. A gate without a proof-of-failure is not merged. This is TDD applied to the
  AI-OS's own guardrails.
- **AI-OS self-eval (the missing harness, §1.2):** a small suite of **golden scenarios**:
  "given task description X, does skill Y trigger? does the wrapper still point at a doc
  that exists and still contains step Z?" Run it in CI. When someone edits a skill, the
  suite says whether behavior regressed. Start with cheap deterministic checks (trigger-
  phrase → expected-skill mapping; wrapper→doc anchor still present), escalate to
  transcript-replay only where it earns its cost.

### 2.4 Memory / Learning plane — *new: close the loop*
The differentiator. A structured, auto-updated memory wired into the engineering loop:
- **Capture:** the `fix-bug` / feature flows end with a "capture lesson" step that writes a
  structured record (symptom → root cause → fix → guard added) — not prose in a chat log.
- **Promote:** recurring lessons graduate to durable artifacts — a repeated bug class →
  a new gate or a standard; a repeated design choice → an ADR. Human approves promotions.
- **Retrieve:** the context-builder consults memory, so past lessons actually reach the
  next session. This is what turns "autocomplete" into "senior engineer": compounding.
- Builds on the user's existing `memory/` auto-memory rather than inventing a parallel one.

### 2.5 Observability plane — *new, lightweight (H2), not a platform (H3)*
At H2 this is deliberately modest: a structured event per skill/agent/command invocation
(name, task-type, outcome, human-intervened?, rough token/turn cost) appended to a local,
git-ignored JSONL, plus a tiny `report` command that summarizes "which skills fired, which
succeeded, which got abandoned." That is enough to *prioritize improvement*. The
full-fidelity tracing / dashboards / RL-feedback version is H3 and explicitly out of scope.

### 2.6 Governance / Safety plane — *new: enforce, don't advise*
Move the security-critical rules from prose into a **PreToolUse hook** (Claude Code's
settings hooks): block writes to protected paths without an approval marker; block reads
of secret paths; block pushes to `main`/force-push. The prose stays (it explains *why*);
the hook makes it *true*. This mirrors the runner's `commandGuard` on the interactive side.
Prompt-injection posture stays partly prose (it's judgment-heavy) but gains a rule:
untrusted repo content never silently overrides the hook layer.

---

## 3. Industry patterns to incorporate (cited honestly)

> **Honesty note.** These orgs do not publish their internal agent harnesses in full. Below
> I use *publicly documented* patterns and clearly separate what is published from what is
> my inference/generalization. I am not claiming inside knowledge of any company's stack.

- **Eval-driven development (published; SWE-bench and the broad agent-eval literature).**
  The industry-wide lesson is that agent quality is a function of a regression eval you run
  continuously. *Applied here (§2.3):* an eval suite for our own skills, not just the app.
  This is the most important borrowed idea.
- **Harness-over-model (Anthropic's published writing on building agents / "context
  engineering").** Capability comes from the scaffolding — tools, context discipline,
  verification — more than from prompt cleverness. *Applied:* V2 invests in planes
  (harness), not bigger prompts. V1 already leans this way; V2 makes it explicit.
- **Agent–Computer Interface discipline (published).** Give agents crisp, well-scoped
  tools/actions with clear contracts. *Applied (§2.2):* trigger precision + the kit seam
  are ACI thinking applied to skills.
- **Verifier / critic loops (widely published pattern).** A separate pass tries to *refute*
  the work before it's accepted. *Applied:* the subagent-driven review loop we already run
  (implementer → adversarial reviewer) *is* this. V2 keeps it and extends it to the AI-OS
  itself (self-eval as a standing critic).
- **Context isolation via subagents (published pattern; also how we already work).** Heavy
  read/search happens in a subagent that returns only the conclusion. *Applied:*
  `repo-context-reader` is exactly this; V2 elevates it to the standard entry point.
- **Golden-transcript / deterministic replay testing (published testing practice).**
  Record known-good runs; replay to detect regressions. *Applied (§2.3):* the escalation
  tier of the self-eval suite.
- **Spec → plan → execute with checkpoints (the Superpowers methodology we're using now).**
  *Applied:* already our backbone; V2 changes nothing here except feeding lessons back in.

What I deliberately do **not** borrow at H2: autonomous multi-agent fleets, RL-from-
production-feedback, and always-on CI agents. They are real and powerful at platform scale;
for one engineer and one repo they are cost and complexity without payoff. They live in §5's
H3 north star.

---

## 4. V1 ↔ V2 comparison (trade-off by trade-off)

| Dimension | V1 (today) | V2 (proposed) | What it costs | What it buys | Worth it now? |
|---|---|---|---|---|---|
| Knowledge | Rich, static docs | Same, retrieved by a context-builder to minimize window | Small build; retrieval discipline | Less context pollution; lower cost/turn | **Yes** — cheap, high value |
| Execution | Self-triggering skills/agents/commands | + trigger-precision rule + portable kit seam | Authoring discipline; a refactor for portability | Scales past 1 repo; fewer mis-fires | **Yes** for precision; **stage** portability |
| Verification | Syntactic gates; some theatre | Real gates (RED-proof required) + AI-OS self-eval suite | Real engineering effort (the eval harness) | Catches AI-OS regressions automatically; kills theatre | **Yes** — highest leverage |
| Memory/Learning | Human-written ADRs; unwired auto-memory | Capture→promote→retrieve loop | Moderate; needs discipline + retrieval wiring | Compounding "seniority"; bugs don't recur | **Yes, incrementally** — start with capture only |
| Observability | None | Lightweight JSONL events + `report` | Small build; some noise | Data to prioritize improvements | **Marginal now** — do only if cheap |
| Governance | Advisory prose | PreToolUse hooks for the hard rules | Hook setup; risk of over-blocking | Deterministic safety for security-critical rules | **Yes** for the 3–4 critical rules only |

**Reading the table:** the ordering of leverage is Verification-self-eval > Governance-
hooks ≈ Knowledge-retrieval > Learning-capture > Execution-portability > Observability.
That ordering — not "build all six" — is the actual recommendation.

### Where V2 is over-engineering (the honest part)
- **A full observability platform** for one engineer is a hobby, not a tool. Cap it at a
  JSONL + a summary command, or skip it until a second contributor exists.
- **Transcript-replay eval** for every skill is expensive and brittle. Use deterministic
  trigger/anchor checks first; replay only for the 2–3 highest-value flows.
- **Auto-promoting memory to ADRs without human approval** would pollute the decision
  record. Promotion must stay human-gated.
- **The learning loop can eat itself:** a memory that captures everything becomes noise.
  Capture must be selective (root causes and reusable patterns), or retrieval degrades.

---

## 5. Evolutionary migration (non-destructive)

The instruction that prompted this doc — "do not optimize the existing proposal; rebuild
from first principles" — is itself the thing a Principal should push back on. A from-
scratch rebuild would discard a *working, reviewed, merging* increment. The right move is
**first-principles thinking, evolutionary delivery**: we reasoned from scratch (§2), and we
now graft it onto what exists in leverage order. Nothing below throws away V1.

**Step 0 — done.** Fix the verification theatre in the spine (commit `0a4ea43`). The V1
spine merges as-is.

**Horizon 1 — this repo, close the top loops (next specs):**
1. **AI-OS self-eval suite** (§2.3) — the highest-leverage single artifact. Deterministic
   trigger/anchor checks first, wired into `make verify` / CI.
2. **Governance PreToolUse hook** (§2.6) — the 3–4 security-critical rules only.
3. **Learning capture step** (§2.4) — add "capture the lesson" to `fix-bug`; write
   structured records into the existing `memory/`. Retrieval + promotion follow once
   capture has produced enough signal.
4. **Doc↔code drift check** (§1.7) — a small linter: recipe references that no longer
   resolve in code get flagged. (Natural extension of the `.claude/` gate.)

**Horizon 2 — portable kit (the center of gravity):**
5. Extract the generic method into `starters/ai-os-kit` with a clear seam: portable method
   (self-contained) vs. repo-specific content (linked/templated). Prove it by standing the
   kit up on a second repo. This is where "consistent across many projects" becomes real.

**Horizon 3 — north star (explicitly not committed):**
Observability at fidelity, autonomous CI agents, agent fleets, RL-from-feedback. Sketched
so the direction is known; not planned, because they exceed what one engineer + one repo
can justify. Revisit when there is a team.

---

## 6. Final recommendation

Do **not** rebuild. Build the **verification self-eval harness** first (it retires the
anchor weakness and makes every future change safe), then the **governance hooks** (cheap,
deterministic safety), then **learning-capture** (start compounding), then **portability**
(the H2 payoff). Treat observability as optional and platform-scale ambitions as a labeled
horizon, not a backlog. The measure of success for V2 is not how much the system knows —
it is whether, six months from now, the system can *prove* it has not regressed and can
*show* which lessons it has compounded. That is the difference between a well-documented
repo and a self-correcting engineer.

---

## Open questions (for the reader)
1. Leverage order in §4 — agree, or do you want governance-hooks before self-eval (safety
   before quality)?
2. Learning plane: start with capture-only (my recommendation), or build capture+retrieve
   together?
3. Portability (H2): prove the kit on a *real* second repo, or a throwaway sandbox first?
