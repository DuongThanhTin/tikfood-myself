import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { ApiError } from "../../lib/auth";
import { AuthForm } from "./AuthForm";

const push = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push })
}));

const login = vi.fn();
const register = vi.fn();
vi.mock("./AuthProvider", () => ({
  useAuth: () => ({
    user: null,
    status: "anonymous",
    login,
    register,
    logout: vi.fn(),
    loginWithGoogleUrl: () => ""
  })
}));

function type(label: string, value: string) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

beforeEach(() => {
  vi.clearAllMocks();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("AuthForm (login)", () => {
  it("shows an error on an invalid email and does not submit", async () => {
    render(<AuthForm mode="login" />);
    type("Email", "not-an-email");
    type("Mật khẩu", "password123");
    fireEvent.click(screen.getByRole("button", { name: "Đăng nhập" }));

    expect(await screen.findByText("Email không hợp lệ.")).toBeInTheDocument();
    expect(login).not.toHaveBeenCalled();
  });

  it("clears a field's inline error once the user fixes it", async () => {
    render(<AuthForm mode="login" />);
    type("Email", "not-an-email");
    type("Mật khẩu", "password123");
    fireEvent.click(screen.getByRole("button", { name: "Đăng nhập" }));
    expect(await screen.findByText("Email không hợp lệ.")).toBeInTheDocument();

    type("Email", "user@example.com");
    expect(screen.queryByText("Email không hợp lệ.")).not.toBeInTheDocument();
  });

  it("disables the submit button while submitting", async () => {
    let resolveLogin: () => void = () => {};
    login.mockReturnValue(new Promise<void>((resolve) => (resolveLogin = resolve)));

    render(<AuthForm mode="login" />);
    type("Email", "user@example.com");
    type("Mật khẩu", "password123");
    fireEvent.click(screen.getByRole("button", { name: "Đăng nhập" }));

    await waitFor(() => expect(screen.getByRole("button", { name: "Đang đăng nhập…" })).toBeDisabled());
    resolveLogin();
  });

  it("renders the API error message on failure", async () => {
    login.mockRejectedValue(new ApiError("unauthorized", "Email hoặc mật khẩu không đúng."));

    render(<AuthForm mode="login" />);
    type("Email", "user@example.com");
    type("Mật khẩu", "password123");
    fireEvent.click(screen.getByRole("button", { name: "Đăng nhập" }));

    expect(await screen.findByText("Email hoặc mật khẩu không đúng.")).toBeInTheDocument();
    expect(push).not.toHaveBeenCalled();
  });

  it("redirects home after a successful login", async () => {
    login.mockResolvedValue(undefined);

    render(<AuthForm mode="login" />);
    type("Email", "user@example.com");
    type("Mật khẩu", "password123");
    fireEvent.click(screen.getByRole("button", { name: "Đăng nhập" }));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/"));
    expect(login).toHaveBeenCalledWith("user@example.com", "password123");
  });
});
