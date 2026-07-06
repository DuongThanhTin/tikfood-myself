# AI-OS Kit — bộ template docs/contract để bootstrap Claude cho project mới

> Rút ra từ workspace này. Mục tiêu: copy folder → điền → có ngay bộ docs + contract
> chuẩn giúp Claude (CLI/agent) phát triển project mới nhanh, nhất quán, ít sai.
> Đây là phần **"cần gì + làm thế nào"**; danh sách file tối thiểu theo tier xem
> [`../ai-os-blueprint.md`](../ai-os-blueprint.md).

## 1. Dùng như thế nào (2 phút)

**Cách khuyến nghị (Claude điền từ code):**
1. Copy folder `ai-os-kit/` vào repo mới, đổi tên các file `*.tmpl` → tên thật
   (`CLAUDE.md.tmpl` → `CLAUDE.md`, `ai-agent.yaml.tmpl` → `.ai-agent.yaml`,
   `ci-verify.yml.tmpl` → `.github/workflows/verify.yml` — file CI là **gated**).
2. Điền phần **"con người biết"** ở §3 vào một đoạn ngắn.
3. Dán **bootstrap prompt** (§5) cho Claude → nó đọc code repo, trích lệnh verify +
   shape contract **từ code**, điền các `<<PLACEHOLDER>>`, xoá block `EXAMPLE`.
4. Chạy `make verify` + link-check → mở PR.

Link trong template đã trỏ tới **tên file cuối** (không có `.tmpl`), nên sau khi đổi tên là dùng được.

## 2. Sáu quy tắc meta (giữ nguyên cho mọi project)
1. **Index, đừng copy** — mỗi rule một nguồn; chỗ khác chỉ link (kiểu AI-CONTRACT).
2. **Code wins over docs** — doc lệch code thì flag + sửa doc.
3. **Một ngôn ngữ / doc** — không bản song song.
4. **Đừng viết docs cho thứ chưa có trong code** — stale ngay.
5. **≤ ~100 dòng / file**, mỗi file mở đầu bằng "> vì sao tồn tại".
6. **Enforce bằng CI**, không chỉ DoD trên giấy; **link-check 0 broken** trước commit.

## 3. Cần THÔNG TIN gì (intake) — 2 loại, đừng lẫn

**A. Chỉ con người biết → bạn cung cấp**
1. Định danh & mục đích (1–3 câu) · 2. Anti-goals / scope · 3. Security & git policy
(không `.env`, không push `main`, branch `ai/*`, cái gì cần approval) · 4. Quyết định đã
chốt + lý do (→ ADR) · 5. Quy ước ngôn ngữ doc · 6. Chuẩn team / PR flow / có CI không.

**B. Phải TRÍCH TỪ CODE → Claude tự đọc repo (đừng đoán)**
7. Layout repo (→ REPOSITORY-MAP) · 8. Tech stack thật từng app (→ apps/*/CLAUDE.md) ·
9. **Lệnh verify thật, chạy được** (→ Makefile/DoD) · 10. Layering & naming rules ·
11. **Contract shape thật** cho mỗi interface (→ docs/contracts, error là *object*) ·
12. Source-of-truth priority (code > CLAUDE.md > app CLAUDE.md > config > docs).

> ⚠️ Mục 9 và 11 là chỗ dễ sai nhất: **luôn verify ngược lại code**, đừng viết từ trí nhớ.

## 4. File nào → thành gì

| Template | Thành | Vai trò |
| --- | --- | --- |
| `CLAUDE.md.tmpl` | `/CLAUDE.md` | File Claude luôn đọc; mọi thứ trỏ từ đây |
| `ai-agent.yaml.tmpl` | `/.ai-agent.yaml` | Guardrails máy-đọc (protected paths, commands) |
| `Makefile.tmpl` | `/Makefile` | `make verify` mỗi app (một lệnh cho agent/người/CI) |
| `ci-verify.yml.tmpl` | `.github/workflows/verify.yml` | CI chạy `make verify` (**gated** — chờ approve) |
| `docs/REPOSITORY-MAP.md.tmpl` | như tên | Dir → mục đích → đọc gì trước khi sửa |
| `docs/ai/AI-CONTRACT.md.tmpl` | như tên | Index rules + thứ tự ưu tiên |
| `docs/ai/CONTEXT-LOADING.md.tmpl` | như tên | Đọc gì, theo loại task |
| `docs/verification/definition-of-done.md.tmpl` | như tên | Gate "xong" theo surface |
| `docs/workflows/*.tmpl` | như tên | Flow request → PR (feature/bug/gated) |
| `docs/adr/{README,TEMPLATE}.md.tmpl` | như tên | Ghi quyết định đắt-để-đảo-ngược |
| `docs/contracts/README.md.tmpl` | như tên | Index mọi contract |
| `docs/contracts/CONTRACT.template.md` | `docs/contracts/<iface>.md` | Template contract dày (HTTP/event/job/CLI) |
| `docs/recipes/RECIPE.template.md` | `docs/recipes/<task>.md` | Playbook cho 1 loại việc lặp lại — **chỉ viết khi đau** |

**Cố ý KHÔNG có (Tier 3 — chỉ viết khi đau):** `handbook/`, `thinking/`, design-system.
Recipe thì đã có sẵn *template* (viết recipe cụ thể khi một loại việc lặp ≥2 lần và bị lệch).
Xem [`../ai-os-blueprint.md`](../ai-os-blueprint.md) §Tier-3.

## 5. Bootstrap prompt (dán cho Claude ở repo mới)

```
Bối cảnh tôi cung cấp (loại A):
- Project: <1-3 câu>   Anti-goals: <...>   Gated: auth/migration/infra/external
- Quyết định đã chốt: <X vì Y>   Ngôn ngữ doc: <...>   PR: branch ai/*, không push main

Việc của bạn — TRƯỚC KHI viết, hãy TRÍCH từ code (loại B, đừng đoán):
1. Đọc repo → điền docs/REPOSITORY-MAP.md (dir → mục đích → đọc gì trước khi sửa).
2. Mỗi app: đọc code → apps/*/CLAUDE.md "Stack Reality" + layering + naming THẬT.
3. Tìm lệnh build/test/lint thật → điền Makefile (make verify mỗi app) → CHẠY THỬ cho pass.
4. Mỗi interface ra ngoài: trích request/success/error shape TỪ CODE → tạo
   docs/contracts/<iface>.md từ CONTRACT.template.md; error phải là object; trỏ schema là source of truth.
Rồi mới điền: CLAUDE.md, AI-CONTRACT, CONTEXT-LOADING, definition-of-done, workflows, adr.
CI (.github/workflows/verify.yml): CHUẨN BỊ nội dung nhưng DỪNG chờ tôi approve (protected path).
Ràng buộc: index đừng copy · code wins · 1 ngôn ngữ/doc · ≤100 dòng/file ·
không viết docs cho thứ chưa tồn tại · xoá mọi block EXAMPLE · link-check 0 broken trước commit.
```

## 6. Quy ước điền (legend)
- `<<PLACEHOLDER: mô tả>>` — giá trị cần điền.
- `<!-- FILL: ... -->` — hướng dẫn điền gì / lấy từ đâu (Tiếng Việt), xoá sau khi điền.
- `<!-- EXAMPLE (Acme Notes) ... --> … <!-- /EXAMPLE -->` — ví dụ mẫu, **xoá sau khi điền**.
- Thân doc = **English** (chuẩn reference-doc); comment FILL + README này = **Tiếng Việt**.
