# Bắt đầu một feature mới

> Dùng khi: bạn có một ý tưởng/feature và muốn làm từ đầu tới PR một cách sạch sẽ,
> để bất kỳ session mới nào cũng hiểu và làm tiếp được.
>
> **Đây là quickstart cho người.** Bản quy trình *chạy-được từng bước* cho AI agent
> (đủ input/artifact/exit/gate mỗi bước) ở [`docs/workflows/feature.md`](../workflows/feature.md).

## Quy trình 6 bước

### 1. Tạo branch sạch (đừng làm trên branch đang dở)
```bash
git checkout dev && git pull            # base = dev (integration). PR thì nhắm main.
git checkout -b ai/<ten-feature>         # PR branch LUÔN dưới tiền tố ai/
```
> Quy tắc từ `CLAUDE.md`: PR branch dưới `ai/`; **không** push `main`/`master`, không
> force-push, không auto-merge. Muốn làm song song với việc khác → dùng
> [worktree](worktrees.md).

### 2. Viết feature brief (bắt buộc — để session sau hiểu)
Dùng [`templates/feature-brief.md`](templates/feature-brief.md). Ngắn gọn: **mục tiêu 1
câu · scope · out-of-scope · acceptance · ràng buộc**. Lưu vào:
- việc nhỏ → đầu PR hoặc một comment;
- việc lớn/nhiều bước → một **spec** trong `docs/superpowers/specs/<ngày>-<feature>-design.md`
  (nếp spec-driven của repo).

### 3. Kickoff Claude
Dán [`templates/session-kickoff.md`](templates/session-kickoff.md) (biến thể "feature mới")
và **yêu cầu brainstorm/plan trước khi code**. Với việc nhiều bước, để Claude viết plan
ra file trước.

### 4. Chọn đúng "đường ray"
- **UI** (`apps/web`) → theo [`docs/design/`](../design/00-overview.md) + playbook phù hợp trong
  [`docs/design/playbooks/`](../design/playbooks/README.md) (`create-page`, `create-component`,
  `create-form`, `create-modal`…).
- **Backend/việc theo domain** → [`docs/recipes/`](../recipes).
- Cần **human approval**: auth, migration, infra, external network call, hoặc chạm
  anti-goal (`CLAUDE.md`).

### 5. Làm — reuse, token, test
Reuse trước; dùng token thay hardcode; data qua `lib/api.ts`; copy tiếng Việt; cập nhật
test khi đổi hành vi.

### 6. Kết thúc & PR
Chạy **`make verify`** (api + web + runner — hoặc `make verify-web` / `verify-api` /
`verify-runner` cho một surface) và đối chiếu
[`docs/verification/definition-of-done.md`](../verification/definition-of-done.md).
Commit (message kết bằng dòng `Co-Authored-By:` theo repo). Mở PR bằng
[`templates/pr-description.md`](templates/pr-description.md) (body kết bằng dòng
`🤖 Generated with Claude Code`).

## Ví dụ mẫu — thêm "empty-state có nút reset" cho venue rail (UI)

```bash
git checkout dev && git pull
git checkout -b ai/rail-empty-state-reset
```

**Brief (dán cho Claude):**
> Feature: khi rail không có venue nào, hiện empty-state kèm nút "Xoá bộ lọc".
> Scope: chỉ `apps/web` rail. Out-of-scope: đổi API, phân trang.
> Acceptance: rỗng → thấy thông báo + nút; bấm nút gọi `clearFilters` và search lại.
> Ràng buộc: discovery-only; theo docs/design; dùng token; không component mới nếu tái dùng được.

**Đường ray:** đọc [`docs/design/patterns/empty-state.md`](../design/patterns/empty-state.md)
+ [`docs/design/playbooks/update-component.md`](../design/playbooks/update-component.md)
(vì đây là mở rộng pattern có sẵn, không phải component mới).

**Kết:** đối chiếu DoD → commit → PR theo template.

## Ví dụ mẫu — feature lớn (backend + web) → tách tab

Nếu feature đụng cả `apps/api` và `apps/web`, cân nhắc **2 worktree/branch** để 2 tab làm
song song không đụng file (xem [`worktrees.md`](worktrees.md)), rồi ghép qua 2 PR.

## Liên quan

- [`new-terminal.md`](new-terminal.md) · [`worktrees.md`](worktrees.md) · [`templates/`](templates/)
- [`docs/design/playbooks/`](../design/playbooks/README.md) · [`docs/recipes/`](../recipes) · [`docs/verification/definition-of-done.md`](../verification/definition-of-done.md)
