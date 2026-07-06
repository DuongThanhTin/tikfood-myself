import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { fetchDiscoveryVenues } from "./api";
import {
  ApiError,
  clearAccessToken,
  getAccessToken,
  getMe,
  login,
  refresh,
  setAccessToken
} from "./auth";

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" }
  });
}

const AUTH_HEADER = (call: unknown[]): string | null => {
  const init = call[1] as RequestInit | undefined;
  return new Headers(init?.headers).get("Authorization");
};

beforeEach(() => {
  clearAccessToken();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("auth client", () => {
  it("login sends credentials and stores the access token", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, {
        data: {
          user: { id: "u1", email: "a@b.com", display_name: "A", email_verified: false, created_at: "", updated_at: "" },
          access_token: "access-abc",
          access_expires_at: "2026-01-01T00:00:00Z"
        }
      })
    );
    vi.stubGlobal("fetch", fetchMock);

    const user = await login("a@b.com", "password123");

    expect(user.email).toBe("a@b.com");
    expect(getAccessToken()).toBe("access-abc");
    const init = fetchMock.mock.calls[0][1] as RequestInit;
    expect(init.credentials).toBe("include");
  });

  it("getMe attaches the Bearer header", async () => {
    setAccessToken("token-123");
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, { data: { user: { id: "u1", email: "a@b.com", display_name: "", email_verified: true, created_at: "", updated_at: "" } } })
    );
    vi.stubGlobal("fetch", fetchMock);

    await getMe();

    expect(AUTH_HEADER(fetchMock.mock.calls[0])).toBe("Bearer token-123");
  });

  it("retries once after a 401 by refreshing the token", async () => {
    setAccessToken("stale-token");
    const fetchMock = vi
      .fn()
      // 1) /me with stale token -> 401
      .mockResolvedValueOnce(jsonResponse(401, { error: { code: "unauthorized", message: "expired" } }))
      // 2) /refresh -> new token
      .mockResolvedValueOnce(jsonResponse(200, { data: { access_token: "fresh-token", access_expires_at: "" } }))
      // 3) /me retried with fresh token -> ok
      .mockResolvedValueOnce(jsonResponse(200, { data: { user: { id: "u1", email: "a@b.com", display_name: "", email_verified: true, created_at: "", updated_at: "" } } }));
    vi.stubGlobal("fetch", fetchMock);

    const user = await getMe();

    expect(user.id).toBe("u1");
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(AUTH_HEADER(fetchMock.mock.calls[2])).toBe("Bearer fresh-token");
    expect(getAccessToken()).toBe("fresh-token");
  });

  it("clears the token when refresh fails", async () => {
    setAccessToken("doomed");
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(401, { error: { code: "unauthorized", message: "session expired" } })
    );
    vi.stubGlobal("fetch", fetchMock);

    await expect(refresh()).rejects.toBeInstanceOf(ApiError);
    expect(getAccessToken()).toBeNull();
  });

  it("does not send a token on public discovery fetches", async () => {
    vi.stubEnv("NEXT_PUBLIC_API_URL", "http://localhost:18081");
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { data: [] }));
    vi.stubGlobal("fetch", fetchMock);
    // Public discovery works with no auth token set.
    const venues = await fetchDiscoveryVenues({ q: "pho" });
    expect(Array.isArray(venues)).toBe(true);
    expect(AUTH_HEADER(fetchMock.mock.calls[0])).toBeNull();
  });
});
