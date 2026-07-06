import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { RequireAuth } from "./RequireAuth";

const replace = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace })
}));

let currentStatus = "loading";
vi.mock("./AuthProvider", () => ({
  useAuth: () => ({ status: currentStatus })
}));

afterEach(() => {
  vi.clearAllMocks();
  currentStatus = "loading";
});

describe("RequireAuth", () => {
  it("redirects anonymous users to /login and hides children", () => {
    currentStatus = "anonymous";
    render(
      <RequireAuth>
        <p>secret</p>
      </RequireAuth>
    );
    expect(replace).toHaveBeenCalledWith("/login");
    expect(screen.queryByText("secret")).not.toBeInTheDocument();
  });

  it("renders children when authenticated", () => {
    currentStatus = "authenticated";
    render(
      <RequireAuth>
        <p>secret</p>
      </RequireAuth>
    );
    expect(screen.getByText("secret")).toBeInTheDocument();
    expect(replace).not.toHaveBeenCalled();
  });

  it("renders nothing while loading", async () => {
    currentStatus = "loading";
    render(
      <RequireAuth>
        <p>secret</p>
      </RequireAuth>
    );
    expect(screen.queryByText("secret")).not.toBeInTheDocument();
    await waitFor(() => expect(replace).not.toHaveBeenCalled());
  });
});
