# Phần 6 — Context Engineering

> ✅ **Trạng thái: Đầy đủ (qua link).** Quy tắc chi tiết sống ở
> [`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md) — một artifact vận hành, dùng
> chung cho cả người và AI. Phần này tóm tắt và trỏ tới đó.

Xây trên thứ tự context-loading ở [Phần 1 §1.4](01-ai-operating-system.md#14-context-loading).
**Nguyên tắc:** nạp *tối thiểu* ngữ cảnh đủ để làm đúng.

## Trình tự

1. **Core luôn đọc:** [`/CLAUDE.md`](../../CLAUDE.md) → [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md)
   → [`REPOSITORY-MAP.md`](../REPOSITORY-MAP.md).
2. **Theo loại task:** đọc đúng hàng trong ma trận ở
   [`CONTEXT-LOADING.md` §2](../ai/CONTEXT-LOADING.md) (backend / frontend / prompt /
   migration / …) — đọc `CLAUDE.md` của app **trước** khi sửa app đó.
3. **Process skill trước** (xem [Phần 7](07-superpowers-integration.md)).

## Khi nào bỏ context / tách sub-agent

- **Bỏ context** khi một file đã xong nhiệm vụ — đừng giữ file lớn "phòng khi".
- **Tách sub-agent** (`Explore` / `general-purpose`) khi câu hỏi phải quét nhiều file mà
  chỉ cần kết luận; giữ tóm tắt, bỏ file dump. Chạy song song các điều tra độc lập.
- **git worktree** cho việc cần cô lập (skill `using-git-worktrees`).
- **Chia nhỏ task** khi một task chạm nhiều vùng — làm từng vùng, một task mỗi phiên.

Chi tiết + checklist: [`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md).

---

Tiếp theo: [Phần 7 — Superpowers Integration](07-superpowers-integration.md).
