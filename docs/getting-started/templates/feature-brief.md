# Template: Feature Brief

> Điền form ngắn này trước khi bắt đầu. Việc nhỏ: để ở đầu PR. Việc lớn: lưu thành
> `docs/superpowers/specs/<ngày>-<feature>-design.md`. Cách dùng: [`../new-feature.md`](../new-feature.md).

---

```md
# Feature: <tên ngắn>

- **Mục tiêu (1 câu):** <người dùng làm được gì / vấn đề gì được giải quyết>
- **Thuộc:** apps/web | apps/api | cả hai | packages/…
- **Scope (làm gì):**
  - <…>
- **Out-of-scope (KHÔNG làm):**
  - <…>
- **Acceptance (coi là xong khi):**
  - [ ] <điều kiện quan sát được 1>
  - [ ] <điều kiện 2>
- **Ràng buộc / rủi ro:**
  - Discovery-only; không đụng anti-goals (cart/order/checkout/payment/booking/chat/follow/monetization/livestream).
  - Cần human approval? (auth / migration / infra / external network call) → <có/không>
  - Reuse trước; token thay hardcode; data qua lib/api.ts; copy tiếng Việt.
- **Đường ray:** UI → docs/design/ + playbook <…>; backend/domain → docs/recipes/<…>.
- **Branch:** ai/<feature>  (base: dev)
```

---

### Ví dụ đã điền

```md
# Feature: Empty-state có nút reset cho venue rail

- **Mục tiêu:** Khi rail rỗng, người dùng thấy thông báo + nút "Xoá bộ lọc" để quay lại kết quả.
- **Thuộc:** apps/web
- **Scope:** thêm nút reset vào empty-state của rail; gọi clearFilters + search lại.
- **Out-of-scope:** đổi API; phân trang; illustration.
- **Acceptance:**
  - [ ] Rỗng → hiện `.emptyText` + nút.
  - [ ] Bấm nút → clearFilters() chạy, list nạp lại.
  - [ ] Hoạt động ở cả light/dark và 2 breakpoint.
- **Ràng buộc:** discovery-only; theo docs/design; không component mới nếu tái dùng được.
- **Đường ray:** docs/design/patterns/empty-state.md + playbooks/update-component.md
- **Branch:** ai/rail-empty-state-reset (base: dev)
```
