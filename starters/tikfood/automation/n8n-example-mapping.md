# n8n Example Mapping

Send feature request JSON to:

```text
http://ai-code-runner:8080/jobs/feature
```

On success, map `branch`, `pr_title`, and `pr_body` into a GitHub Create Pull Request node.

On failure, notify a human with `stage`, `error`, and `recommendation`.
