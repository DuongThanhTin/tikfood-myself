# Phần 2 — Repository Design

Mục tiêu: giúp bạn **định vị được chỗ cần sửa** và hiểu *vì sao* repo được tổ chức như
hiện tại. Nếu bạn chỉ nhớ một điều: mỗi vùng có một file `CLAUDE.md` hoặc chuẩn chi
phối nó — đọc file đó trước khi chạm vào vùng đó.

## 2.1 Monorepo map

| Vùng | Vai trò | Stack thật | File chi phối |
| --- | --- | --- | --- |
| [`apps/api`](../../apps/api) | Backend discovery API | Go + Gin, `database/sql` qua pgx (không GORM), slog, PostGIS | [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md), [`docs/standards/backend/`](../standards/backend/) |
| [`apps/web`](../../apps/web) | Frontend discovery UX | Next.js App Router + TS, MapLibre GL, plain CSS | [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md), [`docs/standards/frontend-architecture.md`](../standards/frontend-architecture.md) |
| [`apps/ai-code-runner`](../../apps/ai-code-runner) | Runner automation (n8n → runner → PR), **MVP skeleton** | TypeScript | [`docs/agents/`](../agents/), [`docs/runner-contract.md`](../runner-contract.md) |
| [`packages/`](../../packages) | Prompts, schemas, config policies dùng chung | — | [`packages/prompts/`](../../packages/prompts/), [`packages/schemas/`](../../packages/schemas/), [`packages/config/`](../../packages/config/) |
| [`workflows/n8n`](../../workflows) | Docs & ví dụ workflow n8n | — | `workflows/n8n/README.md` |
| [`docs/`](../) | Toàn bộ tài liệu (xem 2.2) | — | — |

Bản đồ quy tắc gọn nằm ở mục *Monorepo Map* của [`/CLAUDE.md`](../../CLAUDE.md).

### Điểm quan trọng về từng app

- **`apps/api`** — layering một chiều: `http (Handler) → discovery (Service) →
  Repository → postgres`. Cấm Handler chạm SQL, cấm Service dùng `gin.Context`.
  Response envelope thật là `{ "data", "error" }` — **không** phải `{ "success": ... }`.
- **`apps/web`** — mọi truy cập API qua `lib/api.ts`, không rải `fetch` trong
  component. Copy UI mặc định tiếng Việt. Chưa dùng React Query/Zod.
- **`apps/ai-code-runner`** — là skeleton; bước chưa xong đánh dấu `TODO`, không tuyên
  bố production-ready.

## 2.2 Giải phẫu `docs/`

| Thư mục | Nội dung | Đọc khi nào |
| --- | --- | --- |
| [`docs/ai/`](../ai/) | AI-OS spine: [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md), [`CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md), [`ROADMAP.md`](../ai/ROADMAP.md) | Đầu mỗi phiên |
| [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) | Bản đồ mọi thư mục → mục đích | Khi định vị chỗ cần sửa |
| [`docs/standards/`](../standards/) | Chuẩn kỹ thuật (backend/frontend/API/logging/AI rules/testing/security/performance/prompt) — **nguồn sự thật** | Trước khi đổi code app |
| [`docs/adr/`](../adr/) | Architecture Decision Records (quyết định + lý do) | Khi thiết kế / trước khi trái với một quyết định cũ |
| [`docs/services/`](../services/) | Hợp đồng dịch vụ (api/web/ai-code-runner) | Khi đổi thứ liên-dịch-vụ |
| [`docs/recipes/`](../recipes/) | Playbook từng bước theo loại việc | Khi bắt tay làm một task |
| [`docs/thinking/`](../thinking/) | Khung tư duy / heuristics | Khi phân loại & quyết định |
| [`docs/verification/`](../verification/) | Definition of Done + Review guide | Trước khi tuyên bố "xong" / review |
| [`docs/agents/`](../agents/) | Hợp đồng con người-đọc-được của từng agent runner | Khi làm việc với ai-code-runner |
| [`docs/tikfood/`](../tikfood/) | Định vị sản phẩm, MVP scope, anti-goals, [`domain-model.md`](../tikfood/domain-model.md) | Khi quyết định phạm vi feature |
| [`docs/superpowers/`](../superpowers/) | `specs/` và `plans/` thật của quy trình brainstorm→plan→execute | Khi bắt đầu/tiếp tục một feature |
| [`docs/handbook/`](.) | Chính bộ sổ tay này (giải thích, không phải nguồn sự thật) | Onboarding |
| Các file lẻ | [`architecture.md`](../architecture.md), [`security.md`](../security.md), [`roadmap.md`](../roadmap.md), [`runner-contract.md`](../runner-contract.md), [`local-development.md`](../local-development.md) | Theo chủ đề |

## 2.3 CLAUDE.md phân tầng

Repo dùng CLAUDE.md **phân tầng**, override từ ngoài vào trong:

```text
/CLAUDE.md                 ← quy tắc & anti-goal toàn workspace (đọc trước tiên)
  └── apps/api/CLAUDE.md    ← mô tả backend "như đang được build hôm nay"
  └── apps/web/CLAUDE.md    ← mô tả frontend "như đang được build hôm nay"
```

Mỗi file per-app mở đầu bằng "Read `/CLAUDE.md` first". File per-app mô tả **hiện
trạng** và có mục *Future Direction (approved, do not do preemptively)* — nghĩa là biết
trước hướng đi nhưng **không** làm sớm.

## 2.4 `.ai-agent.yaml` & `packages/config`

[`/.ai-agent.yaml`](../../.ai-agent.yaml) là policy máy-đọc-được cho `ai-code-runner`,
khai báo tường minh:

- `docs_required.before_coding` — danh sách docs bắt buộc đọc trước khi code.
- `protected_paths` — đường dẫn cấm chạm (`.env*`, secrets, `.github/workflows/**`,
  `docker-compose.yml`, `packages/schemas/**`...).
- `allowed_commands` / `blocked_commands` — allowlist lệnh shell.
- `coding_rules`, `testing_rules`, `security_rules`.

[`packages/config/`](../../packages/config/) chứa config theo phạm vi:
`default.ai-agent.yaml` (mặc định) và `tikfood.ai-agent.yaml` (định vị sản phẩm +
anti-goal cho chính app TikFood). Anti-goal phải được giữ trong `packages/prompts/**`,
`packages/config/tikfood.ai-agent.yaml`, `docs/tikfood/**`, `starters/tikfood/**`.

## 2.5 Các thư mục AI-OS (đã tạo)

Các thư mục từng được đề xuất cho tương lai nay **đã tồn tại**, được tạo theo
[`docs/ai/ROADMAP.md`](../ai/ROADMAP.md) — không còn là placeholder:

| Thư mục | Vai trò | Thay cho ý tưởng cũ |
| --- | --- | --- |
| [`docs/recipes/`](../recipes/) | Playbook từng bước theo loại việc | `recipes/` |
| [Phần 3 — Core Workflow](03-core-workflow.md) | Quy trình chung 15 bước | `playbooks/` |
| [`docs/verification/`](../verification/) | Definition of Done + Review guide | `verification/` |
| [`docs/thinking/`](../thinking/) | Khung tư duy / heuristics | `thinking/` |
| [`docs/adr/`](../adr/) | Architecture Decision Records | `adr/` |
| [`docs/services/`](../services/) | Hợp đồng dịch vụ per-service | (mới) |
| [`docs/ai/`](../ai/) | Spine: contract, context-loading, roadmap | (mới) |

Nguyên tắc giữ nguyên: không tạo thư mục rỗng — mỗi thư mục trên có nội dung thật và
**link** tới nguồn sự thật, không chép lại; "code wins" khi mâu thuẫn.

---

Tiếp theo: [Phần 3 — Core Workflow](03-core-workflow.md).
