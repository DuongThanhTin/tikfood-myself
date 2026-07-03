# Phần 4 — Domain Workflows

> ✅ **Trạng thái: Đầy đủ (qua link).** Các recipe cụ thể sống ở [`docs/recipes/`](../recipes/)
> — một artifact tái sử dụng được, không lặp lại trong handbook. Phần này chỉ giải thích
> *khi nào dùng recipe nào*.

Mỗi "recipe" là một playbook từng bước cho một loại việc thường gặp: mở recipe đúng loại,
làm từ trên xuống, dừng ở các cổng cần **human approval**, kết thúc theo
[Definition of Done](../verification/definition-of-done.md).

## Bản đồ recipe

| Loại việc | Recipe |
| --- | --- |
| Thêm/đổi endpoint `apps/api` | [add-api-endpoint](../recipes/add-api-endpoint.md) |
| Thêm/đổi UI `apps/web` | [frontend-component](../recipes/frontend-component.md) |
| Sửa bug | [fix-bug](../recipes/fix-bug.md) |
| Refactor (không đổi hành vi) | [refactor](../recipes/refactor.md) |
| Thêm/sửa runtime prompt | [add-runtime-prompt](../recipes/add-runtime-prompt.md) |
| Migration schema (**cần duyệt**) | [add-migration](../recipes/add-migration.md) |
| Điều tra hiệu năng | [performance-investigation](../recipes/performance-investigation.md) |
| Review bảo mật | [security-review](../recipes/security-review.md) |

Chưa có recipe cho việc bạn cần? Tạo mới theo format ở
[`docs/recipes/README.md`](../recipes/README.md) (Goal · When to use · Context · Skill ·
Steps · Verification · Exit · Common mistakes).

---

Tiếp theo: [Phần 5 — Prompt Library](05-prompt-library.md).
