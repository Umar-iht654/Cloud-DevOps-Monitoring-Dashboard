import { useEffect, useRef, useState, type FormEvent } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import { AuthLayout } from "../components/layout/AuthLayout";
import { ArrowRightIcon } from "../components/ui/Icons";
import { FormField } from "../components/ui/FormField";
import { FullPageLoader } from "../components/ui/FullPageLoader";
import { useAuth } from "../context/AuthContext";

export function RegisterPage() {
  const { register, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [fieldErrors, setFieldErrors] = useState<{
    name?: string;
    email?: string;
    password?: string;
    confirmPassword?: string;
  }>({});
  const [submitting, setSubmitting] = useState(false);
  const errorRef = useRef<HTMLDivElement>(null);
  const formRef = useRef<HTMLFormElement>(null);

  useEffect(() => {
    if (error) errorRef.current?.focus();
  }, [error]);

  if (isLoading) {
    return <FullPageLoader label="Restoring your session" />;
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (submitting) return;

    const normalizedName = name.trim();
    const normalizedEmail = email.trim();
    const normalizedPassword = password.trim();
    const nextErrors: typeof fieldErrors = {};

    if (!normalizedName) nextErrors.name = "Enter your full name.";
    if (!normalizedEmail) {
      nextErrors.email = "Enter your email address.";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalizedEmail)) {
      nextErrors.email = "Enter a valid email address.";
    }
    if (normalizedPassword.length < 7) {
      nextErrors.password = "Use at least 7 characters.";
    } else if (password !== normalizedPassword) {
      nextErrors.password = "Password cannot start or end with spaces.";
    }
    if (!confirmPassword) {
      nextErrors.confirmPassword = "Confirm your password.";
    } else if (password !== confirmPassword) {
      nextErrors.confirmPassword = "Passwords do not match.";
    }

    setFieldErrors(nextErrors);
    setError("");
    if (Object.keys(nextErrors).length > 0) {
      window.requestAnimationFrame(() => {
        formRef.current?.querySelector<HTMLInputElement>('[aria-invalid="true"]')?.focus();
      });
      return;
    }

    setSubmitting(true);
    try {
      await register(normalizedName, normalizedEmail, password);
      navigate("/verify-email", {
        replace: true,
        state: { email: normalizedEmail, registrationStarted: true },
      });
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, "Unable to create your account."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <AuthLayout>
      <div className="mb-8">
        <p className="mb-2 text-sm font-semibold text-cyan-600">Get started</p>
        <h1 className="text-3xl font-semibold tracking-tight text-slate-950">Create your workspace</h1>
        <p className="mt-3 text-sm leading-6 text-slate-500">
          Add your first service and start building a reliability history.
        </p>
      </div>

      {error && (
        <div
          ref={errorRef}
          role="alert"
          tabIndex={-1}
          className="notice-enter mb-5 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700 focus:outline-none"
        >
          {error}
        </div>
      )}

      <form
        ref={formRef}
        onSubmit={handleSubmit}
        noValidate
        aria-busy={submitting}
        className="space-y-4"
      >
        <FormField
          id="name"
          label="Full name"
          autoComplete="name"
          placeholder="Your name"
          required
          autoFocus
          disabled={submitting}
          error={fieldErrors.name}
          value={name}
          onChange={(event) => {
            setName(event.target.value);
            setFieldErrors((current) => ({ ...current, name: undefined }));
            setError("");
          }}
        />
        <FormField
          id="email"
          label="Email address"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          required
          disabled={submitting}
          error={fieldErrors.email}
          value={email}
          onChange={(event) => {
            setEmail(event.target.value);
            setFieldErrors((current) => ({ ...current, email: undefined }));
            setError("");
          }}
        />
        <FormField
          id="password"
          label="Password"
          type="password"
          autoComplete="new-password"
          placeholder="At least 7 characters"
          required
          minLength={7}
          disabled={submitting}
          error={fieldErrors.password}
          hint="Use at least 7 characters."
          passwordToggle
          value={password}
          onChange={(event) => {
            setPassword(event.target.value);
            setFieldErrors((current) => ({
              ...current,
              password: undefined,
              confirmPassword: undefined,
            }));
            setError("");
          }}
        />
        <FormField
          id="confirm-password"
          label="Confirm password"
          type="password"
          autoComplete="new-password"
          placeholder="Repeat your password"
          required
          disabled={submitting}
          error={fieldErrors.confirmPassword}
          passwordToggle
          value={confirmPassword}
          onChange={(event) => {
            setConfirmPassword(event.target.value);
            setFieldErrors((current) => ({ ...current, confirmPassword: undefined }));
            setError("");
          }}
        />

        <button
          type="submit"
          disabled={submitting}
          className="primary-action group flex w-full items-center justify-center gap-2 rounded-xl bg-[#07111f] px-4 py-3.5 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {submitting ? "Creating account…" : "Create account"}
          {!submitting && <ArrowRightIcon className="h-4 w-4 transition group-hover:translate-x-0.5" />}
        </button>
      </form>

      <p className="mt-7 text-center text-sm text-slate-500">
        Already have an account?{" "}
        <Link
          to="/login"
          className="rounded font-semibold text-cyan-700 hover:text-cyan-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-600 focus-visible:ring-offset-2"
        >
          Sign in
        </Link>
      </p>
    </AuthLayout>
  );
}
