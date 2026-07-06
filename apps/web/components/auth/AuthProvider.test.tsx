import { act, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { AuthProvider, useAuth } from "./AuthProvider";

// Mock the auth client so no real network happens.
vi.mock("../../lib/auth", () => ({
  refresh: vi.fn(),
  getMe: vi.fn(),
  login: vi.fn(),
  register: vi.fn(),
  logout: vi.fn(),
  googleLoginUrl: () => "http://api.local/api/v1/auth/google/login"
}));

import { getMe, login, logout, refresh } from "../../lib/auth";

const user = { id: "u1", email: "a@b.com", display_name: "A", email_verified: true, created_at: "", updated_at: "" };

function Probe() {
  const auth = useAuth();
  return (
    <div>
      <span data-testid="status">{auth.status}</span>
      <span data-testid="email">{auth.user?.email ?? "none"}</span>
      <button onClick={() => auth.login("a@b.com", "password123")}>login</button>
      <button onClick={() => auth.logout()}>logout</button>
    </div>
  );
}

beforeEach(() => {
  vi.clearAllMocks();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("AuthProvider", () => {
  it("bootstraps to authenticated when the refresh cookie is valid", async () => {
    vi.mocked(refresh).mockResolvedValue({ access_token: "t", access_expires_at: "" });
    vi.mocked(getMe).mockResolvedValue(user);

    render(
      <AuthProvider>
        <Probe />
      </AuthProvider>
    );

    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("authenticated"));
    expect(screen.getByTestId("email").textContent).toBe("a@b.com");
  });

  it("bootstraps to anonymous when refresh fails", async () => {
    vi.mocked(refresh).mockRejectedValue(new Error("no session"));

    render(
      <AuthProvider>
        <Probe />
      </AuthProvider>
    );

    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("anonymous"));
    expect(screen.getByTestId("email").textContent).toBe("none");
  });

  it("login updates the user", async () => {
    vi.mocked(refresh).mockRejectedValue(new Error("no session"));
    vi.mocked(login).mockResolvedValue(user);

    render(
      <AuthProvider>
        <Probe />
      </AuthProvider>
    );
    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("anonymous"));

    await act(async () => {
      screen.getByText("login").click();
    });

    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("authenticated"));
    expect(screen.getByTestId("email").textContent).toBe("a@b.com");
  });

  it("logout clears the user", async () => {
    vi.mocked(refresh).mockResolvedValue({ access_token: "t", access_expires_at: "" });
    vi.mocked(getMe).mockResolvedValue(user);
    vi.mocked(logout).mockResolvedValue();

    render(
      <AuthProvider>
        <Probe />
      </AuthProvider>
    );
    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("authenticated"));

    await act(async () => {
      screen.getByText("logout").click();
    });

    await waitFor(() => expect(screen.getByTestId("status").textContent).toBe("anonymous"));
    expect(screen.getByTestId("email").textContent).toBe("none");
  });
});
