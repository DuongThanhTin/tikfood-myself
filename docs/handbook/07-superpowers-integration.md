# Phần 7 — Superpowers Integration

> ✅ **Trạng thái: Đầy đủ.** Bảng phân loại skill dùng-nguyên-bản / đặc-thù-TikFood +
> quy ước `docs/superpowers/`.

Xây trên nguyên tắc ở [Phần 1 §1.5](01-ai-operating-system.md#15-skill-orchestration):
**process-skill trước, implementation-skill sau.** Skill cứng tuân theo chính xác; skill
pattern thì thích ứng.

## Dùng nguyên bản (native Superpowers)

| Skill | Dùng khi |
| --- | --- |
| `brainstorming` | Trước mọi việc "xây/thêm X" |
| `writing-plans` → `executing-plans` | Việc nhiều bước; lưu ở `docs/superpowers/plans/` |
| `systematic-debugging` | Trước khi sửa bất kỳ bug nào |
| `test-driven-development` | Khi viết feature/bugfix |
| `verification-before-completion` | Trước khi tuyên bố "xong" |
| `requesting-code-review` / `receiving-code-review` | Quanh bước review |
| `using-git-worktrees` | Việc cần cô lập workspace |
| `frontend-design` | UI mới cần thẩm mỹ có chủ đích |

## Lớp đặc thù TikFood (tự viết, trong repo)

Repo **không override** skill native; thay vào đó bổ sung một lớp tài liệu tái sử dụng
đặc thù TikFood mà skill trỏ tới:

| Nhu cầu TikFood | Artifact trong repo |
| --- | --- |
| Recipe theo domain (Add API, Fix Bug…) | [`docs/recipes/`](../recipes/) |
| Khung tư duy / heuristics | [`docs/thinking/README.md`](../thinking/README.md) |
| Chuẩn "done" + rubric review | [`docs/verification/`](../verification/) |
| Chuẩn viết prompt runtime | [`docs/standards/prompt-engineering.md`](../standards/prompt-engineering.md) |

Ví dụ: `brainstorming` → chốt spec ở `docs/superpowers/specs/`; `writing-plans` → plan ở
`docs/superpowers/plans/`; khi code thì mở [recipe](../recipes/) đúng loại việc.

## Quy ước `docs/superpowers/`

- `specs/` — design/spec đã brainstorm (ví dụ `2026-07-01-handbook-design.md`).
- `plans/` — implementation plan chi tiết.
- Đặt tên: `YYYY-MM-DD-<slug>-design.md` (spec) / `YYYY-MM-DD-<slug>.md` (plan).
- `.superpowers/sdd/` giữ trạng thái làm việc (task brief/report/progress).

---

Tiếp theo: [Phần 8 — TikFood AI Starter](08-tikfood-ai-starter.md).
