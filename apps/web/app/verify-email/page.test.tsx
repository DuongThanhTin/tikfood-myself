import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import VerifyEmailPage from "./page";

let currentToken: string | null = "good-token";
vi.mock("next/navigation", () => ({
  useSearchParams: () => new URLSearchParams(currentToken ? `token=${currentToken}` : "")
}));

const verifyEmail = vi.fn();
vi.mock("../../lib/auth", async () => {
  const actual = await vi.importActual<typeof import("../../lib/auth")>("../../lib/auth");
  return { ...actual, verifyEmail: (token: string) => verifyEmail(token) };
});

afterEach(() => {
  vi.clearAllMocks();
  currentToken = "good-token";
});

describe("VerifyEmailPage", () => {
  it("verifies with the token from the URL and shows success", async () => {
    verifyEmail.mockResolvedValue({ id: "u1", email: "a@b.com" });
    render(<VerifyEmailPage />);

    await waitFor(() => expect(screen.getByText("Email đã được xác minh")).toBeInTheDocument());
    expect(verifyEmail).toHaveBeenCalledWith("good-token");
  });

  it("shows an error when verification fails", async () => {
    verifyEmail.mockRejectedValue(new Error("boom"));
    render(<VerifyEmailPage />);

    await waitFor(() => expect(screen.getByText("Xác minh thất bại")).toBeInTheDocument());
  });

  it("shows an error and does not call the API when the token is missing", async () => {
    currentToken = null;
    render(<VerifyEmailPage />);

    await waitFor(() => expect(screen.getByText("Xác minh thất bại")).toBeInTheDocument());
    expect(verifyEmail).not.toHaveBeenCalled();
  });
});
