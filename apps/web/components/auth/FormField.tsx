"use client";

// FormField is a thin wrapper over the existing `.field` + `.errorText` pattern so the
// six-odd auth inputs don't repeat label/input/error markup. It adds aria-describedby
// wiring for accessibility and is visually identical to `.field`. No new design tokens.

type FormFieldProps = {
  id: string;
  label: string;
  type: "email" | "password" | "text";
  value: string;
  onChange: (value: string) => void;
  error?: string | null;
  autoComplete?: string;
  placeholder?: string;
  disabled?: boolean;
};

export function FormField({
  id,
  label,
  type,
  value,
  onChange,
  error,
  autoComplete,
  placeholder,
  disabled
}: FormFieldProps) {
  const errorId = `${id}-error`;
  return (
    <label className="field authField" htmlFor={id}>
      <span>{label}</span>
      <input
        id={id}
        name={id}
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        autoComplete={autoComplete}
        placeholder={placeholder}
        disabled={disabled}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
      />
      {error ? (
        <span className="errorText" id={errorId} role="alert">
          {error}
        </span>
      ) : null}
    </label>
  );
}
