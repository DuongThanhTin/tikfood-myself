# n8n Workflow

This folder documents the local MVP workflow that sends TikFood feature requests to `ai-code-runner`.

## URLs

- Local host runner URL: `http://localhost:8080/jobs/feature`
- Docker network runner URL for n8n: `http://ai-code-runner:8080/jobs/feature`
- n8n local URL: `http://localhost:5678`

## Manual Trigger Setup

Use the default Manual Trigger node for local testing.

## Set Node JSON

Use one of the example payloads from `examples/feature-requests/`, such as `tikfood-realtime-map.json`.

## HTTP Request Node

- Method: `POST`
- URL: `http://ai-code-runner:8080/jobs/feature`
- Body content type: JSON
- Body: current item JSON

## IF Condition

Branch on:

```text
success == true
```

The MVP runner skeleton returns `success: false` until the real implementation is wired.

## GitHub PR Node

When the runner eventually returns success:

- Base branch: use the configured base branch or `main`
- Head branch: `branch`
- Title: `pr_title`
- Body: `pr_body`

n8n must create PRs only. It must not merge.

## Notification Examples

Success notification should include repo, branch, PR title, changed files, tests run, and risks.

Failure notification should include feature id, stage, error, logs, and recommendation.

## Error Handling

If the HTTP Request node receives a failure response, route it to notification and manual review. Do not retry indefinitely and do not create a PR from failure output.
