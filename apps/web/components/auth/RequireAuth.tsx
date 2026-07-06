"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "./AuthProvider";

const LOGIN_PATH = "/login";

// RequireAuth is a minimal client route guard for authenticated-only surfaces. It
// redirects anonymous users to /login and renders nothing until the session resolves.
// (No protected page ships yet; this is the reusable primitive + its test.)
export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { status } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    // Guard against redirecting to the page we are already on: if /login itself were ever
    // wrapped in a protected subtree, an unconditional replace would loop.
    if (status === "anonymous" && pathname !== LOGIN_PATH) {
      router.replace(LOGIN_PATH);
    }
  }, [status, router, pathname]);

  if (status !== "authenticated") {
    return null;
  }
  return <>{children}</>;
}
