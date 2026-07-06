"use client";

import Link from "next/link";

import { useAuth } from "./AuthProvider";

// AuthCta wires the previously-inert `.authCta` block: it links to /login when
// anonymous and shows the signed-in user + a logout action when authenticated.
export function AuthCta() {
  const { status, user, logout } = useAuth();

  if (status === "authenticated") {
    return (
      <section className="authCta">
        <p>Xin chào, {user?.display_name?.trim() || user?.email}</p>
        <button type="button" onClick={() => logout()} aria-label="Đăng xuất">
          Đăng xuất
        </button>
      </section>
    );
  }

  return (
    <section className="authCta">
      <p>Đăng nhập để lưu lại những địa điểm bạn yêu thích</p>
      <Link className="authCtaButton" href="/login" aria-label="Đăng nhập">
        Đăng nhập ngay
      </Link>
    </section>
  );
}
