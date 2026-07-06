"use client";

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

import {
  type AuthUser,
  getMe,
  googleLoginUrl,
  login as apiLogin,
  logout as apiLogout,
  refresh,
  register as apiRegister
} from "../../lib/auth";

export type AuthStatus = "loading" | "authenticated" | "anonymous";

type AuthContextValue = {
  user: AuthUser | null;
  status: AuthStatus;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, displayName?: string) => Promise<void>;
  logout: () => Promise<void>;
  loginWithGoogleUrl: () => string;
};

const AuthContext = createContext<AuthContextValue | null>(null);

// AuthProvider holds session state for the app. On mount it silently bootstraps the
// session from the httpOnly refresh cookie (refresh -> getMe); it renders a stable
// "loading" status on first paint to avoid an SSR/hydration mismatch.
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [status, setStatus] = useState<AuthStatus>("loading");

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        await refresh();
        const me = await getMe();
        if (active) {
          setUser(me);
          setStatus("authenticated");
        }
      } catch {
        if (active) {
          setUser(null);
          setStatus("anonymous");
        }
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const authed = await apiLogin(email, password);
    setUser(authed);
    setStatus("authenticated");
  }, []);

  const register = useCallback(async (email: string, password: string, displayName = "") => {
    const authed = await apiRegister(email, password, displayName);
    setUser(authed);
    setStatus("authenticated");
  }, []);

  const logout = useCallback(async () => {
    await apiLogout();
    setUser(null);
    setStatus("anonymous");
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ user, status, login, register, logout, loginWithGoogleUrl: googleLoginUrl }),
    [user, status, login, register, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
