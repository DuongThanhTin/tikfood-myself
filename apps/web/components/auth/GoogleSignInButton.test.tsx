import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { GoogleSignInButton } from "./GoogleSignInButton";

vi.mock("./AuthProvider", () => ({
  useAuth: () => ({
    loginWithGoogleUrl: () => "http://api.local/api/v1/auth/google/login"
  })
}));

describe("GoogleSignInButton", () => {
  it("links to the backend Google login URL", () => {
    render(<GoogleSignInButton />);
    const link = screen.getByRole("link", { name: "Đăng nhập với Google" });
    expect(link).toHaveAttribute("href", "http://api.local/api/v1/auth/google/login");
  });
});
