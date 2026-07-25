import { useState, type InputHTMLAttributes, type ReactNode } from "react";

interface FormFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  hint?: string;
  icon?: ReactNode;
  error?: string;
  passwordToggle?: boolean;
}

export function FormField({
  label,
  hint,
  icon,
  error,
  passwordToggle = false,
  id,
  type,
  "aria-describedby": ariaDescribedBy,
  ...props
}: FormFieldProps) {
  const [passwordVisible, setPasswordVisible] = useState(false);
  const canTogglePassword = passwordToggle && type === "password";
  const hintId = id && hint && !error ? `${id}-hint` : undefined;
  const errorId = id && error ? `${id}-error` : undefined;
  const describedBy = [ariaDescribedBy, errorId, hintId].filter(Boolean).join(" ") || undefined;

  return (
    <div className="form-field block">
      <label htmlFor={id} className="form-label mb-2 block text-sm font-semibold text-slate-700">
        {label}
      </label>
      <span className="relative block">
        {icon && (
          <span className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-slate-400">
            {icon}
          </span>
        )}
        <input
          id={id}
          type={canTogglePassword && passwordVisible ? "text" : type}
          aria-invalid={error ? true : undefined}
          aria-describedby={describedBy}
          {...props}
          className={`form-input w-full rounded-xl border bg-white px-3.5 py-3 text-sm text-slate-900 outline-none transition placeholder:text-slate-500 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-500 ${
            error ? "border-rose-400" : "border-slate-300"
          } ${icon ? "pl-10" : ""} ${canTogglePassword ? "pr-12" : ""}`}
        />
        {canTogglePassword && (
          <button
            type="button"
            onClick={() => setPasswordVisible((visible) => !visible)}
            disabled={props.disabled}
            aria-label={passwordVisible ? `Hide ${label.toLowerCase()}` : `Show ${label.toLowerCase()}`}
            aria-controls={id}
            aria-pressed={passwordVisible}
            className="absolute inset-y-0 right-0 flex w-11 items-center justify-center rounded-r-xl text-slate-500 transition hover:text-cyan-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-cyan-600 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {passwordVisible ? (
              <svg
                aria-hidden="true"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="1.8"
                className="h-5 w-5"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M3 3l18 18" />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M10.6 10.7a2 2 0 002.7 2.7M9.9 4.3A10.7 10.7 0 0112 4c5.2 0 8.5 4.2 9.4 5.6a4.3 4.3 0 010 4.8 14.2 14.2 0 01-2.2 2.7M6.2 6.2a14.6 14.6 0 00-3.6 3.4 4.3 4.3 0 000 4.8C3.5 15.8 6.8 20 12 20a10.8 10.8 0 004.1-.8"
                />
              </svg>
            ) : (
              <svg
                aria-hidden="true"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="1.8"
                className="h-5 w-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M2.6 9.6a4.3 4.3 0 000 4.8C3.5 15.8 6.8 20 12 20s8.5-4.2 9.4-5.6a4.3 4.3 0 000-4.8C20.5 8.2 17.2 4 12 4S3.5 8.2 2.6 9.6z"
                />
                <circle cx="12" cy="12" r="3" />
              </svg>
            )}
          </button>
        )}
      </span>
      {error && (
        <span id={errorId} className="mt-1.5 block text-xs font-medium leading-5 text-rose-700">
          {error}
        </span>
      )}
      {hint && !error && (
        <span id={hintId} className="mt-1.5 block text-xs leading-5 text-slate-500">
          {hint}
        </span>
      )}
    </div>
  );
}
