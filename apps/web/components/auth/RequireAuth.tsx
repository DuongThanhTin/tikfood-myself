"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "./AuthProvider";

// RequireAuth is a minimal client route guard for authenticated-only surfaces. It
// redirects anonymous users to /login and renders nothing until the session resolves.
// (No protected page ships yet; this is the reusable primitive + its test.)
export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { status } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (status === "anonymous") {
      router.replace("/login");
    }
  }, [status, router]);

  if (status !== "authenticated") {
    return null;
  }
  return <>{children}</>;
}
