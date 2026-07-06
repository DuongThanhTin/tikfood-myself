import { describe, expect, it } from "vitest";

import { validateEmail, validatePassword, validateRequired } from "./validation";

describe("validateEmail", () => {
  it("rejects empty and malformed emails", () => {
    expect(validateEmail("")).not.toBeNull();
    expect(validateEmail("  ")).not.toBeNull();
    expect(validateEmail("not-an-email")).not.toBeNull();
    expect(validateEmail("a@b")).not.toBeNull();
  });

  it("accepts a well-formed email", () => {
    expect(validateEmail("user@example.com")).toBeNull();
    expect(validateEmail("  user@example.com  ")).toBeNull();
  });
});

describe("validatePassword", () => {
  it("rejects short or empty passwords", () => {
    expect(validatePassword("")).not.toBeNull();
    expect(validatePassword("1234567")).not.toBeNull(); // 7 chars
  });

  it("accepts a password at the 8-char boundary", () => {
    expect(validatePassword("12345678")).toBeNull();
  });
});

describe("validateRequired", () => {
  it("rejects blank values with the field label", () => {
    expect(validateRequired("", "tên")).toBe("Vui lòng nhập tên.");
  });

  it("accepts non-blank values", () => {
    expect(validateRequired("Alice", "tên")).toBeNull();
  });
});
