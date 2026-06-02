# ai-code-runner

This is an MVP skeleton for the TikFood AI automation runner. It is intentionally not production-ready.

OpenAI API calls are intentionally deferred for now to avoid API cost. The current runner does not need `OPENAI_API_KEY` to validate requests, prepare Git workspaces, and read safe repo context.

Implemented now:

- `GET /health`
- `POST /jobs/feature`
- Basic JSON parsing
- Basic required-field validation
- MVP anti-goal detection
- Git workspace preparation: clone, checkout base branch, create local `ai/*` branch, read status
- Safe repo context reading for root docs/config and `docs/**`
- Honest failure response at `stage: "model"` because model planning, edits, review, commit, and push are still TODO

TODO:

- Later, call OpenAI with `packages/prompts/coding-agent.md` when API budget is available.
- Run allowlisted checks.
- Inspect diffs and protected paths.
- Run `packages/prompts/reviewer.md`.
- Commit and push `ai/*` branches.
- Return full success JSON matching `packages/schemas/openapi.yaml`.

Run locally:

```bash
npm install
npm run build
npm start
```

The service listens on `PORT`, default `8080`.
