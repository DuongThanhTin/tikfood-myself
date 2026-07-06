// Authentication client. All auth network access goes through here (mirrors the
// lib/api.ts discovery wrappers). The access token is kept in memory only; the refresh
// token lives in an httpOnly cookie the browser sends automatically with
// `credentials: "include"`. Field names are snake_case to match the Go JSON tags.
import { getClientApiBaseUrl } from "./api";

export type AuthUser = {
  id: string;
  email: string;
  display_name: string;
  email_verified: boolean;
  created_at: string;
  updated_at: string;
};

export type TokenResponse = {
  access_token: string;
  access_expires_at: string;
};

type AuthSession = TokenResponse & { user: AuthUser };

type ApiEnvelope<T> = {
  data: T;
  error?: { code: string; message: string; details?: Record<string, unknown> };
};

// ApiError carries the server's error code + message so the UI can render error.message.
export class ApiError extends Error {
  code: string;
  constructor(code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.code = code;
  }
}

// In-memory access token. Lost on reload — AuthProvider rehydrates via the refresh cookie.
let accessToken: string | null = null;

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAccessToken(token: string | null): void {
  accessToken = token;
}

export function clearAccessToken(): void {
  accessToken = null;
}

function authUrl(path: string): string {
  return `${getClientApiBaseUrl()}${path}`;
}

async function readEnvelope<T>(response: Response): Promise<T> {
  const body = (await response.json().catch(() => null)) as ApiEnvelope<T> | null;
  if (!response.ok || !body) {
    const code = body?.error?.code ?? "internal_error";
    const message = body?.error?.message ?? "Đã xảy ra lỗi. Vui lòng thử lại.";
    throw new ApiError(code, message);
  }
  return body.data;
}

export async function register(email: string, password: string, displayName = ""): Promise<AuthUser> {
  const response = await fetch(authUrl("/api/v1/auth/register"), {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, display_name: displayName })
  });
  const session = await readEnvelope<AuthSession>(response);
  setAccessToken(session.access_token);
  return session.user;
}

export async function login(email: string, password: string): Promise<AuthUser> {
  const response = await fetch(authUrl("/api/v1/auth/login"), {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  const session = await readEnvelope<AuthSession>(response);
  setAccessToken(session.access_token);
  return session.user;
}

// refresh mints a new access token from the refresh cookie. On failure it clears the
// in-memory token and throws, so callers can treat the session as ended.
export async function refresh(): Promise<TokenResponse> {
  const response = await fetch(authUrl("/api/v1/auth/refresh"), {
    method: "POST",
    credentials: "include"
  });
  if (!response.ok) {
    clearAccessToken();
    await readEnvelope<TokenResponse>(response); // throws ApiError
  }
  const tokens = await readEnvelope<TokenResponse>(response);
  setAccessToken(tokens.access_token);
  return tokens;
}

export async function logout(): Promise<void> {
  try {
    await fetch(authUrl("/api/v1/auth/logout"), { method: "POST", credentials: "include" });
  } finally {
    clearAccessToken();
  }
}

export async function getMe(): Promise<AuthUser> {
  const response = await authFetch("/api/v1/auth/me", { method: "GET" });
  const body = await readEnvelope<{ user: AuthUser }>(response);
  return body.user;
}

export function googleLoginUrl(): string {
  return authUrl("/api/v1/auth/google/login");
}

// authFetch attaches the Bearer access token and, on a 401, performs a single silent
// refresh then retries once. The `allowRefresh` guard prevents an infinite loop.
export async function authFetch(path: string, init: RequestInit = {}, allowRefresh = true): Promise<Response> {
  const headers = new Headers(init.headers);
  if (accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  const response = await fetch(authUrl(path), { ...init, headers, credentials: "include" });
  if (response.status === 401 && allowRefresh) {
    const refreshed = await refresh().then(() => true).catch(() => false);
    if (refreshed) {
      return authFetch(path, init, false);
    }
  }
  return response;
}
