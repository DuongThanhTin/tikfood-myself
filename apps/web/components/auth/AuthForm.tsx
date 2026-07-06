"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { ApiError } from "../../lib/auth";
import { validateEmail, validatePassword, validateRequired } from "../../lib/validation";
import { useAuth } from "./AuthProvider";
import { FormField } from "./FormField";

type AuthFormProps = {
  mode: "login" | "register";
  // Optional slot rendered under the submit button (e.g. Google button, links).
  footer?: React.ReactNode;
};

const COPY = {
  login: { submit: "Đăng nhập", loading: "Đang đăng nhập…" },
  register: { submit: "Tạo tài khoản", loading: "Đang tạo tài khoản…" }
} as const;

export function AuthForm({ mode, footer }: AuthFormProps) {
  const router = useRouter();
  const auth = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");

  const [errors, setErrors] = useState<{ email?: string | null; password?: string | null; displayName?: string | null }>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const copy = COPY[mode];

  // Update a field's value and clear its stale inline error (and aria-invalid) as the
  // user edits, so a fixed field stops advertising an error before the next submit.
  function changeField(setter: (value: string) => void, key: "email" | "password" | "displayName") {
    return (value: string) => {
      setter(value);
      setErrors((prev) => (prev[key] ? { ...prev, [key]: null } : prev));
    };
  }

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault();
    setFormError(null);

    const emailError = validateEmail(email);
    const passwordError = validatePassword(password);
    const displayNameError = mode === "register" ? validateRequired(displayName, "tên hiển thị") : null;

    if (emailError || passwordError || displayNameError) {
      setErrors({ email: emailError, password: passwordError, displayName: displayNameError });
      return;
    }
    setErrors({});

    setSubmitting(true);
    try {
      if (mode === "login") {
        await auth.login(email, password);
      } else {
        await auth.register(email, password, displayName);
      }
      router.push("/");
    } catch (error) {
      setFormError(error instanceof ApiError ? error.message : "Đã xảy ra lỗi. Vui lòng thử lại.");
      setSubmitting(false);
    }
  }

  return (
    <form className="authForm" onSubmit={handleSubmit} noValidate>
      {mode === "register" ? (
        <FormField
          id="display_name"
          label="Tên hiển thị"
          type="text"
          value={displayName}
          onChange={changeField(setDisplayName, "displayName")}
          error={errors.displayName}
          autoComplete="name"
          disabled={submitting}
        />
      ) : null}

      <FormField
        id="email"
        label="Email"
        type="email"
        value={email}
        onChange={changeField(setEmail, "email")}
        error={errors.email}
        autoComplete="email"
        disabled={submitting}
      />

      <FormField
        id="password"
        label="Mật khẩu"
        type="password"
        value={password}
        onChange={changeField(setPassword, "password")}
        error={errors.password}
        autoComplete={mode === "login" ? "current-password" : "new-password"}
        disabled={submitting}
      />

      {formError ? (
        <p className="errorText" role="alert">
          {formError}
        </p>
      ) : null}

      <button className="primaryButton large" type="submit" disabled={submitting}>
        {submitting ? copy.loading : copy.submit}
      </button>

      {footer}
    </form>
  );
}
