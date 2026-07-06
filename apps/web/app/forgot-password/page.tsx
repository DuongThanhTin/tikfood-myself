import Link from "next/link";

export const metadata = {
  title: "Quên mật khẩu · TikFood"
};

// Placeholder only — password reset (email delivery) is out of scope for this iteration.
// The submit control is intentionally disabled so it is not mistaken for a working flow.
export default function ForgotPasswordPage() {
  return (
    <main className="authPage">
      <section className="authCard">
        <h1>Quên mật khẩu</h1>
        <p className="authSubtitle">
          Tính năng đặt lại mật khẩu sẽ sớm ra mắt. Hiện tại vui lòng quay lại trang đăng nhập.
        </p>
        <button className="primaryButton large" type="button" disabled aria-disabled="true">
          Sắp ra mắt
        </button>
        <p className="authSwitch">
          <Link href="/login">Quay lại đăng nhập</Link>
        </p>
      </section>
    </main>
  );
}
