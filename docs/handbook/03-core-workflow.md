# Phần 3 — Core Workflow

> ✅ **Trạng thái: Đầy đủ (qua link).** Template chuẩn + 15 bước đã chốt; mỗi bước trỏ
> tới artifact thật trong repo (recipes, thinking, verification) thay vì lặp lại nội dung.

Quy trình cốt lõi từ khi nhận request đến khi hoàn thành. Nguyên tắc: **hiểu trước khi
code, một task mỗi phiên, verification bắt buộc, review tách khỏi implementation**.

Phần này là **bản đồ** quy trình. Bộ **prompt sẵn-dùng theo từng session**
(Discovery/Planning/Implementation/…) nằm ở
[`Claude_Code_Work_Session_Playbook.md`](../../Claude_Code_Work_Session_Playbook.md) —
dùng nó khi cần prompt cụ thể để dán.

## Template cho mỗi bước

```markdown
### <Tên bước>
🎯 Goal — bước này để đạt điều gì.
📥 Input — cần gì trước khi bắt đầu.
📤 Output — tạo ra artifact/kết quả gì.
📚 Documents — link tới nguồn sự thật liên quan.
🧠 Skill/Prompt — skill hoặc prompt dùng cho bước này.
✅ Exit Criteria — điều kiện coi là xong.
⚠️ Common Mistakes — lỗi hay gặp.
```

## 15 bước — ánh xạ tới artifact

| # | Bước | Skill | Artifact/nguồn trong repo |
| --- | --- | --- | --- |
| 1 | Task Classification | (routing) | [`docs/thinking/README.md`](../thinking/README.md) §1 |
| 2 | Context Loading | — | [`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md) (Phần 6) |
| 3 | Requirement Analysis | `brainstorming` | [`docs/thinking/README.md`](../thinking/README.md) |
| 4 | Research | `Explore` / `deep-research` | [`CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md) §4 (delegate) |
| 5 | Brainstorm | `brainstorming` | [Phần 7](07-superpowers-integration.md) |
| 6 | Architecture Design | `brainstorming` → spec | [`docs/adr/`](../adr/), [`docs/architecture.md`](../architecture.md) |
| 7 | Planning | `writing-plans` | `docs/superpowers/plans/` |
| 8 | Task Breakdown | `writing-plans` | idem |
| 9 | Coding | `test-driven-development` | [`docs/recipes/`](../recipes/) (recipe đúng loại việc) |
| 10 | Refactoring | `code-simplifier` | [recipe: refactor](../recipes/refactor.md) |
| 11 | Testing | `test-driven-development` | [`docs/standards/testing.md`](../standards/testing.md) |
| 12 | Verification | `verification-before-completion` | [`definition-of-done.md`](../verification/definition-of-done.md) |
| 13 | Review | `requesting-code-review` | [`review-guide.md`](../verification/review-guide.md) |
| 14 | Documentation | — | [`definition-of-done.md`](../verification/definition-of-done.md) → *Docs to update* |
| 15 | Retrospective | — | [`docs/adr/`](../adr/) (quyết định có đáng ghi ADR?) |

**Đường A (Claude Code tương tác)** dùng bảng trên. **Đường B (ai-code-runner tự động)**
chạy pipeline prompt ở [`packages/prompts/`](../../packages/prompts/) — xem
[Phần 5](05-prompt-library.md).

## Recipe theo domain

Với từng loại việc cụ thể (Add API, Fix Bug, Refactor…), đừng tự dựng lại 15 bước — mở
recipe tương ứng ở [Phần 4](04-domain-workflows.md) → [`docs/recipes/`](../recipes/).

Phần này giữ vai trò **bản đồ** 15 bước. Muốn **chạy** trọn request → PR theo từng loại
task (feature/bug/refactor/gated), dùng [`docs/workflows/`](../workflows/) — nó ghép các
bước trên với recipe + skill + gate bằng link (không lặp nội dung).

---

Tiếp theo: [Phần 4 — Domain Workflows](04-domain-workflows.md).
