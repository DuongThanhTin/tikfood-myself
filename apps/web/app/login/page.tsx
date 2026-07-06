import Link from "next/link";

import { AuthForm } from "../../components/auth/AuthForm";

export const metadata = {
  title: "Đăng nhập · TikFood"
};

export default function LoginPage() {
  return (
    <main className="authPage">
      <section className="authCard">
        <h1>Đăng nhập</h1>
        <p className="authSubtitle">Đăng nhập để lưu lại những địa điểm bạn yêu thích.</p>

        <AuthForm
          mode="login"
          footer={
            <p className="authSwitch">
              Chưa có tài khoản? <Link href="/register">Tạo tài khoản</Link>
              <br />
              <Link href="/forgot-password">Quên mật khẩu?</Link>
            </p>
          }
        />
      </section>
    </main>
  );
}
