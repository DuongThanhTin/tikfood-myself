# Contract — <<INTERFACE_NAME>>

> **Why this file exists:** the single, verified description of this interface's wire
> shape so producers and consumers agree. **Fill every field from the CODE**, not from
> memory — a contract written from memory drifts from what actually ships.
> Works for HTTP/REST, async/event, background-job, and CLI interfaces — keep the
> sections that apply, delete the rest.

## 1. Overview & scope
<!-- FILL: interface này là gì, ai gọi/ai phục vụ, đọc-only hay ghi. -->
<<One or two sentences. Who produces it, who consumes it.>>

## 2. Transport / trigger
<!-- FILL: chọn 1 dạng. HTTP: method+path. Event: tên event + bus/topic. Job: endpoint/queue. CLI: command. -->
<<`GET /api/v1/...`  |  event `<name>` on `<bus>`  |  `POST /jobs/<name>`  |  `mycli <cmd>`>>

## 3. Request / input
<!-- FILL: field bắt buộc + optional, kiểu, ràng buộc validation. Query vs body. -->
| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| <<field>> | <<type>> | yes/no | <<validation / default>> |

## 4. Success response / output
```json
<<{ "data": ... }>>
```

## 5. Error model
<!-- FILL: error thường là OBJECT, không phải string. Verify với code trước khi ghi. -->
```json
{ "data": null, "error": { "code": "<<code>>", "message": "<<human message>>", "details": {} } }
```
`error` is an **object** (`code`, `message`, optional `details`) — not a bare string.

**Error-code catalog** (source: <<where codes are defined in code>>):

| Code | When | HTTP / result |
| --- | --- | --- |
| `<<invalid_request>>` | <<bad input>> | 400 |
| `<<not_found>>` | <<missing resource>> | 404 |
| `<<internal_error>>` | <<unexpected>> | 500 |

## 6. Status / result states
<!-- FILL: HTTP status codes hoặc job states (queued/running/succeeded/failed). -->
<<200 read · 201 create · 204 delete · 400/404/409/422/500  |  job: queued → running → succeeded|failed>>

## 7. Versioning & compatibility
<!-- FILL: cách version (vd /api/v1). Cái gì là breaking change. -->
Versioned via <<e.g. /api/v1>>. **Breaking** = removing/renaming a field or changing its
type/meaning → needs a migration path + human approval. Additive fields are non-breaking.

## 8. Example (full round-trip)
<!-- EXAMPLE (Acme Notes) — thay bằng ví dụ thật của repo -->
Request: `GET /api/v1/notes?tag=work&limit=2`
```json
{ "data": [ { "id": "n_01", "title": "Standup", "tags": ["work"], "updated_at": "2026-01-01T09:00:00Z" } ] }
```
Error: `GET /api/v1/notes?limit=-1`
```json
{ "data": null, "error": { "code": "invalid_request", "message": "limit must be >= 1" } }
```
<!-- /EXAMPLE -->

## 9. Source of truth
<!-- FILL: nơi schema máy-đọc sống (OpenAPI / JSON Schema / proto). Doc này chỉ mô tả; schema + code thắng. -->
Machine-readable schema: `<<packages/schemas/...>>`. Implementation: `<<path in code>>`.
If this doc disagrees with either, **the code/schema wins** — fix this doc.

## 10. Change rules
- Coordinate producer **and** all consumers (e.g. backend + client types) in one change.
- Schema dir is a **protected path** → human approval before editing.
- Record envelope/error-model decisions as an [ADR](../adr/).

## 11. Checklist (adding/changing this contract)
- [ ] Request/response/error verified against the code.
- [ ] Error is an object with a catalogued `code`.
- [ ] Schema (source of truth) updated + consumers updated in the same change.
- [ ] Breaking change? → migration path + human approval.
- [ ] Example round-trip present and accurate.
