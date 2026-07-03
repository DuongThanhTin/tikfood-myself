# Claude Code Work Session Playbook

> Mục tiêu: Chuẩn hóa cách bắt đầu mọi feature hoặc bug để Claude luôn
> có đủ context và làm việc nhất quán.

------------------------------------------------------------------------

# Workflow tổng thể

``` text
Request
    ↓
Discovery
    ↓
Planning
    ↓
Implementation
    ↓
Verification
    ↓
Review
    ↓
Documentation
```

Không bỏ qua bước nào với feature vừa hoặc lớn.

------------------------------------------------------------------------

# Session 1 -- Discovery

## Goal

Hiểu bài toán trước khi code.

## Claude phải làm

-   Đọc README
-   Đọc CLAUDE.md
-   Đọc tài liệu liên quan
-   Xác định service bị ảnh hưởng
-   Tóm tắt yêu cầu
-   Liệt kê điểm chưa rõ
-   KHÔNG viết code

## Prompt

``` text
Your goal is understanding, not coding.

Read the relevant documentation and existing implementation.

Summarize:

- Problem
- Scope
- Impacted services
- Assumptions
- Missing requirements
- Risks

Do not write any code.
Wait for approval.
```

## Output

-   Requirement Summary
-   Questions
-   Impact Analysis

------------------------------------------------------------------------

# Session 2 -- Planning

## Goal

Tạo kế hoạch triển khai.

## Prompt

``` text
Create an implementation plan.

Break the work into small tasks.

For every task include:

- Goal
- Files
- Dependencies
- Risks
- Verification

Do not implement.
```

## Output

Ví dụ

-   Task 1 Schema
-   Task 2 Repository
-   Task 3 Usecase
-   Task 4 API
-   Task 5 Tests

------------------------------------------------------------------------

# Session 3 -- Implementation

## Goal

Chỉ làm một task.

## Prompt

``` text
Implement ONLY Task X.

Before editing:

- explain approach
- identify files

After implementation:

- summarize changes
- explain decisions
- verify
- stop
```

## Output

-   Code
-   Summary
-   Verification

------------------------------------------------------------------------

# Session 4 -- Verification

## Checklist

### API

-   Compile
-   Unit tests
-   Integration tests
-   Backward compatibility
-   Logging
-   Documentation

### Frontend

-   Build
-   Lint
-   Accessibility

### Prompt

-   Examples
-   Edge cases
-   Output validation

------------------------------------------------------------------------

# Session 5 -- Review

## Prompt

``` text
Review against:

- Requirements
- Architecture
- Coding standards
- Performance
- Security
- Maintainability

Rank findings:

Critical
Major
Minor
Nit

Do not rewrite unless requested.
```

------------------------------------------------------------------------

# Session 6 -- Documentation

Claude phải kiểm tra:

-   README
-   Architecture
-   ADR
-   API docs
-   Prompt docs

Nếu thay đổi ảnh hưởng tài liệu thì cập nhật.

------------------------------------------------------------------------

# Opening Checklist

Khi mở tab Claude mới

-   [ ] Xác định mục tiêu phiên
-   [ ] Xác định service liên quan
-   [ ] Đọc CLAUDE.md
-   [ ] Đọc architecture
-   [ ] Đọc standards
-   [ ] Không code ngay
-   [ ] Hoàn thành Discovery trước

------------------------------------------------------------------------

# Definition of Done

## Discovery

-   Hiểu yêu cầu
-   Không còn giả định lớn

## Planning

-   Có task nhỏ
-   Có verification

## Implementation

-   Chỉ một task
-   Có giải thích
-   Có verify

## Review

-   Có danh sách issue
-   Có mức độ ưu tiên

## Documentation

-   Tài liệu được cập nhật nếu cần

------------------------------------------------------------------------

# Golden Rules

1.  Hiểu trước khi code.
2.  Chỉ làm một task mỗi session.
3.  Verification là bắt buộc.
4.  Review tách khỏi implementation.
5.  Cập nhật tài liệu khi kiến trúc hoặc hành vi thay đổi.
