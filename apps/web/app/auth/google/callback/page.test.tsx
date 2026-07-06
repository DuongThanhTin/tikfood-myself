import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import GoogleCallbackPage from "./page";

const replace = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace })
}));

let currentStatus = "loading";
vi.mock("../../../../components/auth/AuthProvider", () => ({
  useAuth: () => ({ status: currentStatus })
}));

afterEach(() => {
  vi.clearAllMocks();
  currentStatus = "loading";
});

describe("GoogleCallbackPage", () => {
  it("redirects home when the session becomes authenticated", async () => {
    currentStatus = "authenticated";
    render(<GoogleCallbackPage />);
    await waitFor(() => expect(replace).toHaveBeenCalledWith("/"));
  });

  it("shows an error and a retry link when authentication failed", () => {
    currentStatus = "anonymous";
    render(<GoogleCallbackPage />);
    expect(screen.getByText("Đăng nhập Google thất bại")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Quay lại đăng nhập" })).toHaveAttribute("href", "/login");
    expect(replace).not.toHaveBeenCalled();
  });

  it("shows a loading state while bootstrapping", () => {
    currentStatus = "loading";
    render(<GoogleCallbackPage />);
    expect(screen.getByRole("status")).toHaveTextContent("Đang hoàn tất đăng nhập…");
  });
});
