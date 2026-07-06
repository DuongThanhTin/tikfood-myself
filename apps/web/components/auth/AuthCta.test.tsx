import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { AuthCta } from "./AuthCta";

const logout = vi.fn();
let ctx: { status: string; user: { display_name: string; email: string } | null } = { status: "anonymous", user: null };

vi.mock("./AuthProvider", () => ({
  useAuth: () => ({ ...ctx, logout })
}));

afterEach(() => {
  vi.clearAllMocks();
  ctx = { status: "anonymous", user: null };
});

describe("AuthCta", () => {
  it("links to /login when anonymous", () => {
    ctx = { status: "anonymous", user: null };
    render(<AuthCta />);
    expect(screen.getByRole("link", { name: "Đăng nhập" })).toHaveAttribute("href", "/login");
  });

  it("shows the user and a logout action when authenticated", () => {
    ctx = { status: "authenticated", user: { display_name: "Alice", email: "a@b.com" } };
    render(<AuthCta />);
    expect(screen.getByText(/Xin chào, Alice/)).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Đăng xuất" }));
    expect(logout).toHaveBeenCalled();
  });
});
