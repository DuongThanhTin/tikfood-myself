"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "../../../../components/auth/AuthProvider";

const HOME_PATH = "/";

// Landing page for the Google OAuth redirect. The backend has already set the httpOnly
// refresh cookie, so AuthProvider's bootstrap (refresh -> getMe) authenticates the
// session — no access token is carried in the URL. We simply react to the resulting
// status: redirect home on success, show a retry path on failure.
export default function GoogleCallbackPage() {
  const { status } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    // Redirect home once authenticated, but never onto the page we are already on.
    if (status === "authenticated" && pathname !== HOME_PATH) {
      router.replace(HOME_PATH);
    }
  }, [status, router, pathname]);

  if (status === "anonymous") {
    return (
      <main className="authPage">
        <section className="authCard">
          <h1>Đăng nhập Google thất bại</h1>
          <p className="authSubtitle">Không thể hoàn tất đăng nhập bằng Google. Vui lòng thử lại.</p>
          <p className="authSwitch">
            <Link href="/login">Quay lại đăng nhập</Link>
          </p>
        </section>
      </main>
    );
  }

  return (
    <main className="authPage">
      <section className="authCard">
        <p className="authSubtitle" role="status">
          Đang hoàn tất đăng nhập…
        </p>
      </section>
    </main>
  );
}
