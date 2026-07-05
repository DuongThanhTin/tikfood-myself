# Chạy song song bằng git worktree

> Dùng khi: bạn muốn mở **nhiều tab terminal làm nhiều việc cùng lúc** mà không giẫm chân
> nhau (không phải liên tục `git stash`/đổi branch trong cùng một thư mục).

## Vấn đề & cách giải

Một thư mục repo chỉ ở **một branch** tại một thời điểm. Hai tab cùng thư mục sửa hai việc
khác nhau sẽ đụng nhau. **`git worktree`** cho phép một repo có **nhiều thư mục làm việc**,
mỗi thư mục ở một branch riêng → mỗi tab một worktree, độc lập hoàn toàn.

## Lệnh cốt lõi

```bash
# Tạo worktree mới + branch mới cho feature (base = dev)
git worktree add -b ai/<feature> ../tikfood-<feature> dev

# Liệt kê các worktree hiện có
git worktree list

# Xoá khi xong (sau khi đã merge/không cần)
git worktree remove ../tikfood-<feature>
git worktree prune          # dọn metadata thừa
```

> Ghi chú: tất cả worktree **chia sẻ chung một .git/history**. Một branch chỉ được
> checkout ở **một** worktree tại một thời điểm.

## Quy ước cho repo này

- Đặt worktree **cạnh** repo, không lồng bên trong: `../tikfood-<feature>`.
- Branch vẫn theo tiền tố **`ai/<feature>`** (luật `CLAUDE.md`); base từ `dev`.
- Mỗi worktree cần cài phụ thuộc riêng (không dùng chung `node_modules`, `.next`, build của Go):
  ```bash
  cd ../tikfood-<feature>/apps/web && npm install
  ```
- **Chia việc để không overlap file** giữa các tab (vd: tab A = `apps/api`, tab B = `apps/web`).

## Ví dụ mẫu — 2 feature song song

**Mục tiêu:** cùng lúc làm (1) một endpoint discovery mới ở `apps/api`, và (2) một tinh
chỉnh UI ở `apps/web`.

```bash
# Tab 1 — backend
git worktree add -b ai/api-nearby-endpoint ../tikfood-api-nearby dev
cd ../tikfood-api-nearby
# mở Claude ở tab này: "Đọc CLAUDE.md + apps/api/CLAUDE.md. Feature: <...>. Plan trước."

# Tab 2 — web (mở terminal khác)
git worktree add -b ai/web-filter-chip ../tikfood-web-chip dev
cd ../tikfood-web-chip/apps/web && npm install
# mở Claude ở tab này: "Đọc CLAUDE.md + apps/web/CLAUDE.md + docs/design/. Feature: <...>."
```

→ Hai tab, hai thư mục, hai branch, hai PR — không stash, không đụng nhau.
Khi xong mỗi cái: commit → PR → `git worktree remove ../tikfood-...`.

## Trong một session Claude

Bạn cũng có thể nhờ Claude tự làm trong worktree cô lập (nó có sẵn cơ chế worktree cho
tác vụ dễ xung đột). Nhưng để **bạn tự mở nhiều tab tay**, dùng `git worktree` CLI ở trên
là đủ và rõ ràng nhất.

## Lỗi hay gặp

- *"fatal: '<branch>' is already checked out"* → branch đó đang mở ở worktree khác; đổi
  branch khác hoặc `git worktree list` để tìm.
- Quên `npm install` trong worktree mới → web build lỗi thiếu module.
- Xoá thư mục worktree bằng tay mà quên `git worktree prune` → metadata rác.

## Liên quan

- [`new-terminal.md`](new-terminal.md) · [`new-feature.md`](new-feature.md) · [`README.md`](README.md)
- Luật branch/commit: [`../../CLAUDE.md`](../../CLAUDE.md) (Security & Git)
