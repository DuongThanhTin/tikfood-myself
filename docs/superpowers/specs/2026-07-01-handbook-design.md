# Spec — TikFood AI Handbook (`docs/handbook/`)

Date: 2026-07-01
Status: Approved (blocks A/B/C approved by user)

## Mục tiêu

Xây một bộ **onboarding handbook** cho người: giải thích hệ thống AI vận hành thế nào
trong repo `tikfood-ai-automation-starter`, cấu trúc repo được thiết kế ra sao, quy
trình làm việc cốt lõi, các workflow theo domain, thư viện prompt, context engineering,
tích hợp Superpowers, và cách áp dụng vào chính repo TikFood.

Đây **không** phải nguồn sự thật mới về quy tắc — handbook **link** tới nguồn đã có
(`CLAUDE.md`, `.ai-agent.yaml`, `docs/standards`, `docs/agents`, `docs/tikfood`,
`docs/superpowers`) thay vì chép lại.

## Quy ước

- **Ngôn ngữ:** văn xuôi giải thích + tiêu đề tiếng Việt; prompt / spec / exit-criteria
  / code / tên file / tên khái niệm kỹ thuật giữ **tiếng Anh**.
- **Chống trùng lặp:** mỗi mục quy tắc trỏ tới file nguồn, không sao chép nội dung.
- **Trung thực với hiện trạng:** mô tả repo *như nó đang có*. Thư mục chưa tồn tại
  (`recipes/`, `playbooks/`, `verification/`, `thinking/`, `adr/`) được đánh dấu rõ là
  "đề xuất tương lai", không tạo thư mục rỗng. (Theo nguyên tắc CLAUDE.md: code wins,
  flag mismatch.)

## Kiến trúc thư mục

```text
docs/handbook/
├── README.md                     # Master index
├── 01-ai-operating-system.md     # PHẦN 1 — đầy đủ
├── 02-repository-design.md       # PHẦN 2 — đầy đủ
├── 03-core-workflow.md           # PHẦN 3 — skeleton + step template + link
├── 04-domain-workflows.md        # placeholder
├── 05-prompt-library.md          # placeholder
├── 06-context-engineering.md     # placeholder
├── 07-superpowers-integration.md # placeholder
└── 08-tikfood-ai-starter.md      # placeholder
```

## Phạm vi lần này

- Viết đầy đủ: README (index), Phần 1, Phần 2.
- Phần 3: skeleton + template chuẩn cho mỗi bước (Goal/Input/Output/Docs/Prompt/Exit/Mistakes).
- Phần 4–8: placeholder có dàn ý + "TODO (future cycle)".

## Nội dung Phần 1 — AI Operating System

1. Handbook này là gì / không là gì.
2. Mô hình 3 lớp: nguồn sự thật / quy trình (skills) / thực thi (Claude Code + ai-code-runner).
3. Luồng request → hoàn thành (2 đường: interactive vs automated qua n8n).
4. Context loading (trỏ Phần 6).
5. Skill orchestration (process trước, implementation sau; trỏ Phần 7).
6. Ranh giới an toàn (link `docs/security.md`, mục Security của CLAUDE.md).

## Nội dung Phần 2 — Repository Design

1. Monorepo map (apps/*, packages/, workflows/n8n, docs/) — vai trò + CLAUDE.md chi phối + link.
2. Giải phẫu `docs/` (standards, agents, tikfood, superpowers, handbook).
3. CLAUDE.md phân tầng (gốc vs per-app, thứ tự override).
4. `.ai-agent.yaml` & `packages/config` (policy cho ai-code-runner).
5. Thư mục đề xuất tương lai (recipes/playbooks/verification/thinking/adr) — chưa tạo.

## Exit criteria (lần này)

- `docs/handbook/` tồn tại với 9 file như cây trên.
- README index liệt kê 8 phần + quy ước ngôn ngữ + trạng thái từng phần.
- Phần 1, 2 đầy đủ, mọi tham chiếu trỏ tới file thật đang tồn tại.
- Không tạo thư mục rỗng cho các phần "future".
- Không có tuyên bố sai về hiện trạng (vd runner là skeleton phải nói rõ).
