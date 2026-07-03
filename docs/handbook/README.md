# TikFood AI Handbook

Sổ tay onboarding: **hệ thống AI vận hành như thế nào bên trong repo này**.

Đây là tài liệu dành cho *người mới* (dev, PM, người vận hành automation) muốn hiểu
Claude / `ai-code-runner` làm việc ra sao trong `tikfood-ai-automation-starter`. Nó
**không** thay thế các nguồn sự thật hiện có — nó *giải thích và nối* chúng lại.

## Quy ước

- **Ngôn ngữ:** phần giải thích và tiêu đề viết bằng **tiếng Việt**. Phần kỹ thuật —
  prompt, spec, exit-criteria, đoạn code, tên file, tên khái niệm — giữ nguyên **tiếng
  Anh** để AI agent và phần còn lại của repo dùng lại được.
- **Không trùng lặp:** handbook **link** tới nguồn sự thật (xem bảng dưới), không chép.
  Khi nội dung mâu thuẫn với code, **code thắng** — báo lại chỗ lệch, đừng theo docs cũ.

## Nguồn sự thật (handbook chỉ trỏ tới, không chép)

| Chủ đề | File nguồn |
| --- | --- |
| **AI-OS spine** (hợp đồng, bản đồ, context) | [`docs/ai/`](../ai/): [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md), [`CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md), [`ROADMAP.md`](../ai/ROADMAP.md) · [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) |
| Quyết định kiến trúc (ADR) | [`docs/adr/`](../adr/) |
| Hợp đồng dịch vụ | [`docs/services/`](../services/) |
| Recipe theo domain | [`docs/recipes/`](../recipes/) |
| Khung tư duy & verification | [`docs/thinking/`](../thinking/), [`docs/verification/`](../verification/) |
| Quy tắc & anti-goal toàn workspace | [`/CLAUDE.md`](../../CLAUDE.md) |
| Quy tắc backend Go | [`/apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) |
| Quy tắc frontend Next.js | [`/apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) |
| Policy cho ai-code-runner | [`/.ai-agent.yaml`](../../.ai-agent.yaml), [`/packages/config`](../../packages/config) |
| Chuẩn kỹ thuật | [`docs/standards/`](../standards/) |
| Hợp đồng agent | [`docs/agents/`](../agents/), [`packages/prompts/`](../../packages/prompts/) |
| Định vị sản phẩm | [`docs/tikfood/`](../tikfood/) |
| Kiến trúc | [`docs/architecture.md`](../architecture.md) |
| An toàn | [`docs/security.md`](../security.md) |
| Spec/plan workflow thực tế | [`docs/superpowers/`](../superpowers/) |

## Mục lục

| # | Phần | Trạng thái |
| --- | --- | --- |
| 1 | [AI Operating System Overview](01-ai-operating-system.md) — Claude hoạt động thế nào trong repo | ✅ Đầy đủ |
| 2 | [Repository Design](02-repository-design.md) — vì sao repo được tổ chức như vậy | ✅ Đầy đủ |
| 3 | [Core Workflow](03-core-workflow.md) — 15 bước từ request tới hoàn thành | ✅ Đầy đủ (qua link) |
| 4 | [Domain Workflows](04-domain-workflows.md) — Add API / Fix Bug / Refactor... | ✅ Đầy đủ ([`docs/recipes/`](../recipes/)) |
| 5 | [Prompt Library](05-prompt-library.md) — bộ prompt hoàn chỉnh | ✅ Đầy đủ ([`packages/prompts/`](../../packages/prompts/)) |
| 6 | [Context Engineering](06-context-engineering.md) — đọc file nào, khi nào tách agent | ✅ Đầy đủ ([`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md)) |
| 7 | [Superpowers Integration](07-superpowers-integration.md) — skill nào dùng/override/tự viết | ✅ Đầy đủ |
| 8 | [TikFood AI Starter](08-tikfood-ai-starter.md) — áp dụng vào từng khu vực repo | ✅ Đầy đủ (qua link) |

## Cách đọc

- Người mới hoàn toàn → đọc Phần 1 rồi Phần 2.
- Sắp bắt tay làm task → đọc Phần 3 (quy trình), rồi Phần 4 (workflow đúng domain).
- Cần prompt cụ thể → nhảy thẳng Phần 5.
- Muốn tối ưu cách nạp ngữ cảnh → Phần 6.
