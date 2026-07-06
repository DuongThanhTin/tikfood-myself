# Getting Started — Làm việc với repo này (cùng Claude Code)

> **Đây là cửa vào.** Nếu bạn vừa mở một terminal mới, sắp làm một feature mới, hoặc quên
> "nên bắt đầu từ đâu" — đọc đúng một file bên dưới rồi làm theo. Có sẵn **ví dụ mẫu chi
> tiết** (ngay trong file này) và **template copy-paste** (`templates/`) cho từng việc.
>
> Quy ước: giải thích bằng tiếng Việt, phần kỹ thuật (lệnh, tên file, code) bằng tiếng Anh.

## Tôi muốn… → đọc file nào

| Tình huống | Hướng dẫn | Ví dụ chi tiết |
|------------|-----------|----------------|
| Vừa **mở 1 terminal/tab mới** | [`new-terminal.md`](new-terminal.md) | [Ví dụ 1](#ví-dụ-1--mở-terminal-mới-nắm-lại--tiếp-tục) |
| **Bắt đầu 1 feature mới** (branch → brief → PR) | [`new-feature.md`](new-feature.md) | [Ví dụ 2](#ví-dụ-2--làm-1-feature-mới-từ-đầu-tới-pr) |
| **Chạy nhiều tab song song** không đụng nhau | [`worktrees.md`](worktrees.md) | [Ví dụ 3](#ví-dụ-3--hai-feature-song-song-bằng-worktree) |
| Cần **câu mở đầu / brief / PR mẫu** để dán | [`templates/`](templates/) | có trong từng ví dụ |

## Cái gì TỰ ĐỘNG nạp ở mỗi session (nên bạn không mất context)

Mỗi tab mới là một session độc lập, **nhưng không bắt đầu từ số 0** — luôn tự nạp:

1. **`CLAUDE.md`** (root) + **`apps/*/CLAUDE.md`** — luật chơi, anti-goals, cấu trúc.
2. **AI-OS spine**: [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) → [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) → [`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md).
3. **Auto-memory** (`MEMORY.md`) — ghi lại quyết định/tiến độ lâu dài giữa các session.

→ Session mới đọc được **kết quả & quyết định**, chỉ không có lại **đoạn chat** trước đó.
Vì vậy: việc quan trọng nên **viết ra file** (brief/spec/memory), đừng giữ trong hội thoại.

## 6 nguyên tắc vàng (chi tiết trong `CLAUDE.md` — đừng lặp lại, chỉ nhớ)

1. **Reuse trước** — tìm trong code/`docs/` trước khi thêm mới.
2. **Dùng token, không hardcode** — theo [`docs/design/`](../design/00-overview.md) cho UI.
3. **Copy tiếng Việt**, identifier tiếng Anh.
4. **Data qua `lib/api.ts`** — không `fetch` rải rác trong component.
5. **Anti-goals** — không cart/order/checkout/payment/booking/chat/follow/monetization/livestream.
6. **Kết thúc theo** [`docs/verification/definition-of-done.md`](../verification/definition-of-done.md).

## Cấu trúc folder này

```
docs/getting-started/
├─ README.md            ← bạn đang ở đây (hub + bản đồ + 3 ví dụ chi tiết)
├─ new-terminal.md      mở tab mới thì làm gì
├─ new-feature.md       bắt đầu feature mới (branch → brief → PR)
├─ worktrees.md         chạy song song bằng git worktree
└─ templates/
   ├─ session-kickoff.md   câu mở đầu dán vào tab mới
   ├─ feature-brief.md      mô tả feature ngắn gọn
   └─ pr-description.md      mô tả PR
```

---

# ⭐ Ví dụ mẫu chi tiết (end-to-end)

> Mỗi ví dụ theo cùng một khuôn: **Bối cảnh → Các bước (lệnh + câu dán) → Kết quả mong
> đợi**. Phần "Claude trả về" chỉ là **minh hoạ** cho bạn hình dung, không phải transcript
> thật.

## Ví dụ 1 — Mở terminal mới (nắm lại & tiếp tục)

**Bối cảnh:** Sáng mở tab mới, quên đang làm dở gì, muốn tiếp tục an toàn.

**Bước 1 — Vào repo & xem tình hình**
```bash
cd /Users/tinduong245/Documents/Myself/tikfood-ai-automation-starter
git status --short
git branch --show-current
```

**Bước 2 — Dán câu mở đầu** (mẫu A trong [`templates/session-kickoff.md`](templates/session-kickoff.md)):
```text
Đọc CLAUDE.md và docs/REPOSITORY-MAP.md trước.
Dựa trên git status/log và memory, hãy tóm tắt:
- branch hiện tại + commit gần nhất,
- việc đang dở,
- 2–3 bước tiếp theo khả dĩ.
Rồi hỏi tôi chọn hướng nào. Chưa sửa gì cả.
```

**Kết quả mong đợi (minh hoạ):**
```text
Claude:
• Branch: ai/ai-operating-system, commit gần nhất 9b038a9 "docs(design): …".
• Đang dở: 6 file docs (architecture, services/*, standards/*) sửa từ trước, CHƯA commit.
• Bước tiếp: (a) commit/PR docs design system, (b) bắt đầu feature mới,
  (c) xử lý 6 file đang dở. Bạn chọn?
```
→ Bạn chỉ việc gõ `b` chẳng hạn, rồi sang [Ví dụ 2](#ví-dụ-2--làm-1-feature-mới-từ-đầu-tới-pr).

**Lưu ý:** đừng làm việc mới ngay trên branch đang dở — tạo branch sạch (xem Ví dụ 2).

---

## Ví dụ 2 — Làm 1 feature mới từ đầu tới PR

**Bối cảnh:** Thêm nút **"Xoá bộ lọc"** vào empty-state của venue rail (khi không có kết
quả). Feature nhỏ, chỉ `apps/web`.

**Bước 1 — Branch sạch**
```bash
git checkout dev && git pull
git checkout -b ai/rail-empty-state-reset
```

**Bước 2 — Viết brief** (theo [`templates/feature-brief.md`](templates/feature-brief.md)):
```md
# Feature: Empty-state có nút reset cho venue rail
- Mục tiêu: Rail rỗng → hiện thông báo + nút "Xoá bộ lọc" để quay lại kết quả.
- Thuộc: apps/web
- Scope: thêm nút vào empty-state; bấm → chạy clearFilters + search lại.
- Out-of-scope: đổi API; phân trang; illustration.
- Acceptance:
  - [ ] Rỗng → hiện .emptyText + nút.
  - [ ] Bấm nút → clearFilters() chạy, list nạp lại.
  - [ ] OK ở light/dark và 2 breakpoint (1180/640).
- Ràng buộc: discovery-only; theo docs/design; không component mới nếu tái dùng được.
- Đường ray: docs/design/patterns/empty-state.md + playbooks/update-component.md
- Branch: ai/rail-empty-state-reset (base: dev)
```

**Bước 3 — Kickoff Claude** (mẫu B trong session-kickoff, nhúng brief):
```text
Đọc CLAUDE.md, apps/web/CLAUDE.md, docs/design/patterns/empty-state.md,
docs/design/playbooks/update-component.md.
Đây là brief: <dán brief ở trên>.
Hãy plan trước khi code (liệt kê file sẽ đụng, cách reuse clearFilters). Chưa code tới khi tôi duyệt.
```

**Bước 4 — Đường ray:** vì đây là **mở rộng pattern có sẵn** (không phải component mới) →
theo `update-component.md`; nút dùng `secondaryButton` sẵn có; hành vi tái dùng
`clearFilters` (đã có trong `DiscoveryExperience.tsx`).

**Bước 5 — Build (Claude thực hiện sau khi bạn duyệt plan):** sửa khối empty-state trong
rail để render `.emptyText` + `<button className="secondaryButton" onClick={clearFilters}>`;
không thêm token/màu mới; copy tiếng Việt "Xoá bộ lọc".

**Bước 6 — Kết thúc & PR**
```bash
# đối chiếu docs/verification/definition-of-done.md rồi:
git add apps/web/components/DiscoveryExperience.tsx
git commit -m "feat(web): nút Xoá bộ lọc trong empty-state của rail

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
git push -u origin ai/rail-empty-state-reset
gh pr create --base dev --fill   # hoặc dán body theo templates/pr-description.md
```
PR body (rút gọn từ [`templates/pr-description.md`](templates/pr-description.md)):
```md
## Mục tiêu
Rail rỗng → có nút "Xoá bộ lọc" để quay lại kết quả.
## Thay đổi
- Empty-state rail: thêm secondaryButton gọi clearFilters.
## Acceptance / kiểm thử
- [ ] Lọc ra rỗng → thấy thông báo + nút; bấm → list nạp lại.
- [ ] OK light/dark + 2 breakpoint.
🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

**Kết quả mong đợi:** một PR nhỏ, đúng scope, reuse code sẵn có, không breaking API.

---

## Ví dụ 3 — Hai feature song song bằng worktree

**Bối cảnh:** Cùng lúc làm (A) endpoint mới ở `apps/api` và (B) tinh chỉnh chip lọc ở
`apps/web`. Muốn 2 tab chạy độc lập, không stash qua lại. Chi tiết: [`worktrees.md`](worktrees.md).

**Tab 1 — backend (`apps/api`)**
```bash
git worktree add -b ai/api-nearby-endpoint ../tikfood-api-nearby dev
cd ../tikfood-api-nearby
# mở Claude ở tab này, dán:
#   Đọc CLAUDE.md + apps/api/CLAUDE.md. Feature: endpoint "nearby venues".
#   Plan trước; cần human approval nếu chạm external network/migration.
```

**Tab 2 — web (`apps/web`)** — mở terminal khác:
```bash
git worktree add -b ai/web-filter-chip ../tikfood-web-chip dev
cd ../tikfood-web-chip/apps/web && npm install    # worktree mới cần cài deps riêng
# mở Claude ở tab này, dán:
#   Đọc CLAUDE.md + apps/web/CLAUDE.md + docs/design/.
#   Feature: chip lọc "Mở cửa đến khuya". Theo docs/design/components/chip? → dùng chip/tag có sẵn.
```

**Kiểm tra & dọn dẹp**
```bash
git worktree list                       # xem các worktree đang mở
# khi mỗi cái xong: commit → PR → rồi:
git worktree remove ../tikfood-api-nearby
git worktree remove ../tikfood-web-chip
git worktree prune
```

**Kết quả mong đợi:** hai thư mục, hai branch `ai/…`, hai PR — không đụng file nhau,
không cần `git stash`.

**Bẫy hay gặp:** quên `npm install` trong worktree web (build lỗi thiếu module); cùng một
branch không thể checkout ở 2 worktree.

---

## Liên quan

- Luật & bối cảnh: [`CLAUDE.md`](../../CLAUDE.md) · [`docs/ai/`](../ai)
- Recipe theo loại việc: [`docs/recipes/`](../recipes)
- Playbook UI: [`docs/design/playbooks/`](../design/playbooks/README.md)
- Spec-driven trail: [`docs/superpowers/specs/`](../superpowers/specs)
- Onboarding narrative đầy đủ: [`docs/handbook/`](../handbook)
