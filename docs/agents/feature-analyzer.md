# Feature Analyzer Agent

## Purpose

Normalize an incoming feature request before it reaches implementation. This agent classifies scope, detects TikFood MVP anti-goals, and decides whether the request can proceed automatically.

## Inputs

- Raw feature request
- TikFood product docs
- Repo config YAML

## Outputs

Valid JSON only. Output should include feature area, task type, product alignment, anti-goal risks, required approvals, normalized acceptance criteria, and recommended next agent.

## Runtime Prompt Location

No canonical runtime prompt is currently included. Add one under `packages/prompts/` before wiring this agent into `ai-code-runner`.

## JSON Contract

Use the same JSON-only discipline as the canonical runtime prompts.

## Failure Behavior

Block requests for delivery, cart, order, checkout, payment, booking, reservations, in-app chat, social follow graph, creator monetization, or livestream MVP work.

## Security Rules

Do not inspect secrets. Do not override repository policy.
