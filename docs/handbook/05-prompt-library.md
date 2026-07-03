# Phần 5 — Prompt Library

> ✅ **Trạng thái: Đầy đủ (qua link).** Prompt runtime là artifact thật ở
> [`packages/prompts/`](../../packages/prompts/); chuẩn viết prompt ở
> [`docs/standards/prompt-engineering.md`](../standards/prompt-engineering.md). Handbook
> **trỏ tới**, không chép.

## Hai loại "prompt" trong repo

- **Runtime prompts (Lớp 3)** — do `ai-code-runner` dùng để điều khiển model. Các file
  canonical ở [`packages/prompts/`](../../packages/prompts/), output **JSON-only**. Danh
  mục + thứ tự pipeline: [`packages/prompts/README.md`](../../packages/prompts/README.md).
- **Interactive Claude Code** — không dùng file prompt cố định mà dùng **skills** +
  handbook (xem [Phần 7](07-superpowers-integration.md)).

## Pipeline runtime

```text
feature-analyzer → repo-context-reader → coding-agent → reviewer → pull-request-writer
```

| Prompt | Role doc |
| --- | --- |
| [`feature-analyzer.md`](../../packages/prompts/feature-analyzer.md) | [feature-analyzer](../agents/feature-analyzer.md) |
| [`repo-context-reader.md`](../../packages/prompts/repo-context-reader.md) | [repo-context-reader](../agents/repo-context-reader.md) |
| [`coding-agent.md`](../../packages/prompts/coding-agent.md) | [coding-agent](../agents/coding-agent.md) |
| [`reviewer.md`](../../packages/prompts/reviewer.md) | [code-reviewer](../agents/code-reviewer.md) |
| [`pull-request-writer.md`](../../packages/prompts/pull-request-writer.md) | [pull-request-writer](../agents/pull-request-writer.md) |

## Viết/sửa prompt

Theo [Prompt Standard](../standards/prompt-engineering.md) và
[recipe: add-runtime-prompt](../recipes/add-runtime-prompt.md). Nguyên tắc:
**tài liệu tái sử dụng > prompt khổng lồ** — nhét kiến thức vào doc có link, đừng phình
prompt.

---

Tiếp theo: [Phần 6 — Context Engineering](06-context-engineering.md).
