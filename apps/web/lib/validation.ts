// Pure, dependency-free field validators for the auth forms. Each returns a
// Vietnamese error message, or null when the value is valid. No form library (per the
// frontend standard).

export const MIN_PASSWORD_LENGTH = 8;

export function validateEmail(value: string): string | null {
  const email = value.trim();
  if (!email) {
    return "Vui lòng nhập email.";
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    return "Email không hợp lệ.";
  }
  return null;
}

export function validatePassword(value: string): string | null {
  if (!value) {
    return "Vui lòng nhập mật khẩu.";
  }
  if (value.length < MIN_PASSWORD_LENGTH) {
    return `Mật khẩu phải có ít nhất ${MIN_PASSWORD_LENGTH} ký tự.`;
  }
  return null;
}

export function validateRequired(value: string, label: string): string | null {
  if (!value.trim()) {
    return `Vui lòng nhập ${label}.`;
  }
  return null;
}
