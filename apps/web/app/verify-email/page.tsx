"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";

import { ApiError, verifyEmail } from "../../lib/auth";

type VerifyState = "verifying" | "success" | "error";

// Reads the ?token from the emailed link and calls the API to verify the email, then
// reports the outcome. Wrapped in Suspense because useSearchParams requires it.
function VerifyEmailInner() {
  const params = useSearchParams();
  const token = params.get("token") ?? "";
  const [state, setState] = useState<VerifyState>("verifying");
  const [message, setMessage] = useState("");

  useEffect(() => {
    let active = true;
    if (!token) {
      setState("error");
      setMessage("Liên kết xác minh không hợp lệ.");
      return;
    }
    (async () => {
      try {
        await verifyEmail(token);
        if (active) {
          setState("success");
        }
      } catch (err) {
        if (active) {
          setState("error");
          setMessage(err instanceof ApiError ? err.message : "Xác minh email thất bại. Vui lòng thử lại.");
        }
      }
    })();
    return () => {
      active = false;
    };
  }, [token]);

  if (state === "success") {
    return (
      <main className="authPage">
        <section className="authCard">
          <h1>Email đã được xác minh</h1>
          <p className="authSubtitle">Cảm ơn bạn. Tài khoản của bạn đã được xác minh.</p>
          <p className="authSwitch">
            <Link href="/login">Đến trang đăng nhập</Link>
          </p>
        </section>
      </main>
    );
  }

  if (state === "error") {
    return (
      <main className="authPage">
        <section className="authCard">
          <h1>Xác minh thất bại</h1>
          <p className="authSubtitle">{message}</p>
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
          Đang xác minh email…
        </p>
      </section>
    </main>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <main className="authPage">
          <section className="authCard">
            <p className="authSubtitle" role="status">
              Đang tải…
            </p>
          </section>
        </main>
      }
    >
      <VerifyEmailInner />
    </Suspense>
  );
}
