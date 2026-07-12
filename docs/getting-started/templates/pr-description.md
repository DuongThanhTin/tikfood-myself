# Template: PR Description

> Dùng khi mở PR (branch `ai/<feature>` → `main`/`dev`). Không auto-merge; chờ human review.
> Cách dùng: [`../new-feature.md`](../new-feature.md).

---

```md
## Mục tiêu
<1–2 câu: feature này làm gì / vì sao>

## Thay đổi
- <thay đổi chính 1>
- <thay đổi chính 2>

## Scope & Out-of-scope
- Thuộc: apps/web | apps/api | …
- KHÔNG bao gồm: <…>

## Acceptance / Cách kiểm thử
- [ ] <bước kiểm thử 1 / kết quả mong đợi>
- [ ] <bước 2>

## Verification (evidence)   ← BẮT BUỘC, dán output THẬT (không phải checkbox)
<!-- Dán đuôi output của `make verify-*` đã chạy: lệnh + kết quả + exit code.
     verify-evidence-guard.py chặn commit/PR nếu surface bị đổi mà chưa verify tươi. -->
```
$ make verify-api
ok  github.com/tikfood/api/...   12 passed
exit 0
```
<!-- Nếu là BUG: dán regression test đỏ-trước / xanh-sau + tên test -->
- Regression test: `<TestName>` — RED trước fix, GREEN sau fix (dán 2 lần chạy ở trên).
- Skipped surface (nếu có): <ví dụ verify-web bỏ vì không có harness> — nói rõ, không giấu.

## Ràng buộc đã tuân thủ
- [ ] Discovery-only, không đụng anti-goals
- [ ] Reuse/token/`lib/api.ts`; không breaking API/response
- [ ] Copy tiếng Việt; a11y (label icon-only, alt ảnh)
- [ ] Cần human approval? <không / có — lý do>

## Ảnh / ghi chú (nếu UI)
<screenshot hoặc mô tả>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

---

### Nhắc quan trọng (luật repo)

- PR branch dưới `ai/`; **không** push thẳng `main`/`master`, **không** force-push, **không** auto-merge.
- Commit message kết bằng dòng: `Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>`
- PR body kết bằng dòng: `🤖 Generated with [Claude Code](https://claude.com/claude-code)`
