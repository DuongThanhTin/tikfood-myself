import Link from "next/link";

import { AuthForm } from "../../components/auth/AuthForm";
import { GoogleSignInButton } from "../../components/auth/GoogleSignInButton";

export const metadata = {
  title: "Tạo tài khoản · TikFood"
};

export default function RegisterPage() {
  return (
    <main className="authPage">
      <section className="authCard">
        <h1>Tạo tài khoản</h1>
        <p className="authSubtitle">Tạo tài khoản để lưu và theo dõi các địa điểm đang thịnh hành.</p>

        <AuthForm
          mode="register"
          footer={
            <>
              <div className="authDivider">
                <span>hoặc</span>
              </div>
              <GoogleSignInButton />
              <p className="authSwitch">
                Đã có tài khoản? <Link href="/login">Đăng nhập</Link>
              </p>
            </>
          }
        />
      </section>
    </main>
  );
}
