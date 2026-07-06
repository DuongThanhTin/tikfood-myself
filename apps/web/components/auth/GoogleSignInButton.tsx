"use client";

import { useAuth } from "./AuthProvider";

// GoogleSignInButton is a thin wrapper over the existing .secondaryButton language. It
// is a plain link to the backend's server-side Google login URL (Authorization Code
// flow) — no foreign branded widget, no client secret.
export function GoogleSignInButton() {
  const { loginWithGoogleUrl } = useAuth();
  return (
    <a className="secondaryButton googleButton" href={loginWithGoogleUrl()} aria-label="Đăng nhập với Google">
      <span className="googleG" aria-hidden="true">
        G
      </span>
      Đăng nhập với Google
    </a>
  );
}
