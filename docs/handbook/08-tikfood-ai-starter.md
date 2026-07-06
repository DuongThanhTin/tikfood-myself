# Phần 8 — TikFood AI Starter

> ✅ **Trạng thái: Đầy đủ (qua link).** Áp dụng toàn bộ hệ thống (Phần 1–7) vào từng khu
> vực repo: đọc gì, dùng recipe nào, ranh giới ở đâu. Mỗi ô trỏ tới artifact thật.

Bắt đầu bất kỳ task nào: đọc core ([`AI-CONTRACT.md`](../ai/AI-CONTRACT.md) +
[`REPOSITORY-MAP.md`](../REPOSITORY-MAP.md)), phân loại task
([`thinking`](../thinking/README.md)), nạp context
([`CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md)), rồi mở recipe đúng khu vực.

## Theo khu vực

| Khu vực | Đọc trước | Recipe | Ranh giới |
| --- | --- | --- | --- |
| [`apps/api`](../../apps/api) | [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md), [`services/api.md`](../services/api.md), [`standards/backend/`](../standards/backend/) | [add-api-endpoint](../recipes/add-api-endpoint.md), [fix-bug](../recipes/fix-bug.md) | Giữ layering; envelope `{data,error}`; migration cần duyệt |
| [`apps/web`](../../apps/web) | [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md), [`services/web.md`](../services/web.md) | [frontend-component](../recipes/frontend-component.md) | Data qua `lib/api.ts`; plain CSS; không data giả |
| [`apps/ai-code-runner`](../../apps/ai-code-runner) | [`services/ai-code-runner.md`](../services/ai-code-runner.md), [`contracts/runner.md`](../contracts/runner.md) | (mở rộng skeleton) | Chỉ branch `ai/*`; không push main; không đọc secret |
| [`packages/prompts`](../../packages/prompts) | [`prompt-engineering.md`](../standards/prompt-engineering.md), [`prompts/README.md`](../../packages/prompts/README.md) | [add-runtime-prompt](../recipes/add-runtime-prompt.md) | JSON-only; giữ đồng bộ với `docs/agents/` |
| [`packages/schemas`](../../packages/schemas) | [`contracts/runner.md`](../contracts/runner.md) | — | **Protected path** → cần duyệt |
| [`docs/`](../) | [`REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) | — | Cập nhật doc theo [DoD → Docs to update](../verification/definition-of-done.md) |
| [`workflows/n8n`](../../workflows) | [`workflows/n8n/README.md`](../../workflows/n8n/README.md) | — | n8n tạo PR, **không** merge |

## Cổng cần human approval (mọi khu vực)

Auth, migration, infra, external network call, thay đổi model-cost, protected path, và
mọi vùng MVP anti-goal → dừng, xin duyệt (xem [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md)).

## Kết
Đây là điểm cuối của handbook. Vòng lặp chuẩn: **Discovery → Planning → Implementation →
Verification → Review → Documentation**, mỗi bước có artifact tương ứng ở các phần trên.
