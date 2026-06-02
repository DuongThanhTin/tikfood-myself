# Architecture

Target architecture:

```text
Social ingestion
-> PostgreSQL/PostGIS
-> Trend scoring workers
-> AI summary workers
-> Go API
-> Next.js frontend
```

Keep trend scoring, AI summaries, ingestion, and HTTP handlers separated.
