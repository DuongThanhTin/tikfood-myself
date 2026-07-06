# AI-OS Blueprint — bộ docs tối thiểu cho một project mới

> Đúc kết từ workspace TikFood: những gì một repo mới **cần** để Claude (CLI/agent)
> làm việc tốt nhất — và những gì **không nên** viết sớm. Nguyên tắc xuyên suốt:
> ít file, mỗi rule một nguồn, code wins, link thay vì copy, một ngôn ngữ mỗi doc.

## Tier 1 — Ngày đầu tiên (bắt buộc, ~5 file)

| File | Nội dung | Vì sao |
| --- | --- | --- |
| `CLAUDE.md` (root) | Project là gì (3 câu) · anti-goals · security/git rules (không đọc `.env`, không push main, branch `ai/*`) · map thư mục ngắn · rule "code wins over docs" | File duy nhất Claude luôn đọc — mọi thứ khác trỏ từ đây |
| `docs/REPOSITORY-MAP.md` | Bảng: thư mục → mục đích → đọc gì trước khi sửa | Trả lời "X ở đâu, được đụng không" trong 1 file |
| `docs/ai/AI-CONTRACT.md` | **Index** rules trỏ về nguồn (không copy nội dung) + thứ tự ưu tiên khi mâu thuẫn | Chống drift: rule sống ở 1 chỗ, contract chỉ trỏ |
| `docs/verification/definition-of-done.md` | Checklist "xong" theo surface (backend/frontend) + lệnh verify cụ thể | Không có cái này agent tự nhận "done" |
| `Makefile` | `make verify` mỗi app (lint + test + build) | 1 lệnh chung cho agent, người, và CI — TikFood đang thiếu |

## Tier 2 — Tuần đầu (khi bắt đầu có feature thật)

| File | Nội dung |
| --- | --- |
| `docs/workflows/feature.md`, `bug.md`, `gated-change.md` | Flow request → PR theo loại task, mỗi bước: goal · action · artifact · exit · 🚦 gate (dừng chờ human) |
| `docs/adr/` (README + TEMPLATE) | Ghi quyết định kiến trúc — chỉ khi quyết định "đắt để đảo ngược" |
| `apps/*/CLAUDE.md` | Rule riêng từng app (layering, naming, contract), override root trong phạm vi app |
| CI (`.github/workflows/verify.yml`) | Chạy đúng `make verify` — enforce thay vì tin agent tự giác |
| Test setup từng app | Không có test runner thì rule "update tests" vô nghĩa |

## Tier 3 — Chỉ viết khi đau (không viết trước)

- `docs/recipes/<task>.md` — khi một loại việc lặp ≥2 lần và agent làm sai/lệch cách.
- `docs/thinking/` — khi standards có nhiều chỗ "tuỳ judgment" cần heuristics.
- `docs/handbook/` — khi onboarding người mới thật sự diễn ra.
- Design system docs — khi UI đủ lớn để có inconsistency.

## Anti-patterns (TikFood đã dính, tránh lặp lại)

1. **Viết docs trước khi có code/nhu cầu** → stale ngay, agent đọc nhầm spec cũ.
2. **Hai bản ngôn ngữ song song** → drift; chọn 1 canonical/doc (xem ADR-0008 TikFood).
3. **Copy rule vào nhiều file** → mâu thuẫn; index + link (kiểu AI-CONTRACT).
4. **Map/index không cập nhật khi tree đổi** → agent tìm thư mục ma.
5. **Docs nhiều, enforce ít** → có DoD nhưng không CI = rule chỉ trên giấy.
6. **Xoá file nhưng không sửa mọi reference** → broken links; chạy link-check trước commit.

## Trình tự setup gợi ý (paste cho Claude CLI)

```
1. Viết CLAUDE.md root: <mô tả project>, anti-goals: <...>, security/git rules
   chuẩn (không .env, không push main, branch ai/*, gated: auth/migration/
   infra/external), rule code-wins.
2. Tạo docs/REPOSITORY-MAP.md từ tree hiện tại.
3. Tạo docs/ai/AI-CONTRACT.md dạng index trỏ về CLAUDE.md + các nguồn.
4. Tạo docs/verification/definition-of-done.md với lệnh verify thật
   chạy được của repo này.
5. Tạo Makefile: make verify = lint + test + build từng app.
6. Tạo docs/workflows/{README,feature,bug,gated-change}.md — orchestrate
   bằng link, mỗi bước có exit criteria và gate.
7. Tạo docs/adr/{README,TEMPLATE}.md.
8. Thêm CI chạy make verify (cần tôi approve trước khi tạo).
Ràng buộc: không copy nội dung giữa docs — chỉ link; mỗi doc một ngôn ngữ;
mỗi file ≤ 100 dòng; không viết docs cho thứ chưa tồn tại trong code.
```

## Template sẵn có

Bộ template điền-vào-chỗ-trống cho toàn bộ file trên: [`ai-os-kit/`](ai-os-kit/README.md)
— copy folder, đổi `*.tmpl` → tên thật, điền theo comment `FILL`, xoá block `EXAMPLE`.
Contract dùng [`ai-os-kit/docs/contracts/CONTRACT.template.md`](ai-os-kit/docs/contracts/CONTRACT.template.md)
(đa loại: HTTP/event/job/CLI, có error-code catalog + source-of-truth).
