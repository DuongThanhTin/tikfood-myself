# Template: Câu mở đầu cho tab mới

> Dán một trong các mẫu dưới vào một session Claude Code mới. Thay phần `<…>`.
> Xem cách dùng: [`../new-terminal.md`](../new-terminal.md), [`../new-feature.md`](../new-feature.md).

---

## A. Nắm lại trạng thái / tiếp tục việc cũ

```text
Đọc CLAUDE.md và docs/REPOSITORY-MAP.md trước.
Dựa trên git status/log và memory, hãy tóm tắt:
- branch hiện tại + commit gần nhất,
- việc đang dở,
- 2–3 bước tiếp theo khả dĩ.
Rồi hỏi tôi chọn hướng nào. Chưa sửa gì cả.
```

## B. Bắt đầu một feature mới

```text
Đọc CLAUDE.md, docs/REPOSITORY-MAP.md, docs/ai/CONTEXT-LOADING.md.
Nếu là UI: đọc thêm docs/design/ và playbook phù hợp trong docs/design/playbooks/.

Feature: <mục tiêu 1 câu>
Scope: <thuộc apps/web hay apps/api; phần nào>
Out-of-scope: <những gì KHÔNG làm>
Acceptance: <điều kiện coi là xong>
Ràng buộc: discovery-only, không đụng anti-goals; reuse trước; dùng token; data qua lib/api.ts.

Hãy brainstorm/plan trước khi code. Việc nhiều bước thì viết plan ra file trước.
```

## C. Bắt đầu từ một brief/spec đã viết sẵn

```text
Đọc CLAUDE.md rồi đọc <đường-dẫn brief hoặc docs/superpowers/specs/…>.
Xác nhận bạn hiểu đúng, nêu rủi ro/điểm mơ hồ, rồi đề xuất plan.
Chưa code cho tới khi tôi duyệt plan.
```

## D. Review một thay đổi UI

```text
Đọc docs/design/playbooks/review-ui.md.
Review diff hiện tại theo checklist đó (consistency với token, reuse, a11y, responsive,
anti-goals). Trả về must-fix vs nice-to-have, kèm file:line.
```

---

### Mẹo

- Luôn nhắc "**plan trước khi code**" cho việc >1 bước.
- Nếu làm song song nhiều tab → tạo branch/worktree riêng trước (xem
  [`../worktrees.md`](../worktrees.md)).
