# Mở một terminal / tab mới

> Dùng khi: bạn vừa mở một tab terminal mới và muốn Claude Code (hoặc chính bạn) hiểu
> repo ngay và làm tiếp được. Xem thêm [`new-feature.md`](new-feature.md) nếu bắt đầu
> việc mới, [`worktrees.md`](worktrees.md) nếu muốn chạy song song.

## Cái gì tự nạp (nhắc lại nhanh)

Session mới trong repo tự đọc `CLAUDE.md` + `apps/*/CLAUDE.md`, AI-OS spine
(`docs/ai/AI-CONTRACT.md` → `docs/REPOSITORY-MAP.md` → `docs/ai/CONTEXT-LOADING.md`) và
auto-memory. Bạn **không cần** dán lại các file này — chỉ cần trỏ đúng việc.

## Checklist khi mở tab mới

```bash
# 1. Vào đúng thư mục repo (CLAUDE.md tự nạp từ đây)
cd /Users/tinduong245/Documents/Myself/tikfood-ai-automation-starter

# 2. Xem mình đang ở branch nào, có gì dở dang không
git status --short
git branch --show-current

# 3. (nếu cần) cập nhật base
git checkout dev && git pull
```

Sau đó gõ câu mở đầu tùy mục đích (mẫu đầy đủ trong
[`templates/session-kickoff.md`](templates/session-kickoff.md)).

## Câu mở đầu mẫu

**A. Tiếp tục việc cũ / hỏi trạng thái** (session mới, muốn nắm lại):
> Đọc `CLAUDE.md` và `docs/REPOSITORY-MAP.md`. Tóm tắt trạng thái hiện tại của repo và
> việc đang dở (dựa trên git + memory), rồi hỏi tôi muốn làm gì tiếp.

**B. Bắt đầu ngay một việc cụ thể:**
> Đọc `CLAUDE.md`, `docs/REPOSITORY-MAP.md`, `docs/ai/CONTEXT-LOADING.md`.
> Việc: **\<một câu>**. Nếu là UI, theo `docs/design/`. Brainstorm/plan trước khi code.

## Ví dụ mẫu

**Bối cảnh:** Sáng thứ Hai, mở tab mới, muốn tiếp tục nhưng quên đang dở gì.

```text
Bạn dán:
> Đọc CLAUDE.md + docs/REPOSITORY-MAP.md. Dựa trên git status/log và memory,
> tóm tắt việc đang dở và đề xuất bước tiếp theo.

Claude sẽ:
1. Đọc các file được trỏ + auto-memory.
2. Chạy git status/log để xem branch & commit gần nhất.
3. Trả về: "Branch ai/... đang có commit X (design system docs). Việc dở: Y.
   Bạn muốn (a) mở PR, (b) làm feature mới, (c) …?"
```

→ Bạn chỉ cần chọn hướng; không phải kể lại lịch sử.

## Lưu ý

- **Đừng làm việc mới ngay trên branch đang dở** của người khác/việc khác — xem
  [`new-feature.md`](new-feature.md) để tạo branch sạch.
- Nếu 2 tab cùng sửa 1 file → dễ đụng nhau. Muốn song song thật sự thì dùng
  [worktree](worktrees.md).

## Liên quan

- [`new-feature.md`](new-feature.md) · [`worktrees.md`](worktrees.md) · [`templates/session-kickoff.md`](templates/session-kickoff.md)
- [`README.md`](README.md) (bản đồ quyết định) · [`../../CLAUDE.md`](../../CLAUDE.md)
