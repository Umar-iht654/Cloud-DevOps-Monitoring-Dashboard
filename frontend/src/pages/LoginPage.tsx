import { useState, type FormEvent } from "react";
import { Link, Navigate, useLocation, useNavigate } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import { AuthLayout } from "../components/layout/AuthLayout";
import { ArrowRightIcon } from "../components/ui/Icons";
import { FormField } from "../components/ui/FormField";
import { FullPageLoader } from "../components/ui/FullPageLoader";
import { useAuth } from "../context/AuthContext";

export function LoginPage() {
  const { login, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const navigationState = location.state as {
    registrationSuccess?: boolean;
    email?: string;
    from?: { pathname?: string };
  } | null;
  const [email, setEmail] = useState(navigationState?.email ?? "");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const registrationSuccess = Boolean(navigationState?.registrationSuccess);

  if (isLoading) {
    return <FullPageLoader label="Restoring your session" />;
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      await login(email, password);
      const destination = navigationState?.from?.pathname || "/dashboard";
      navigate(destination, { replace: true });
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, "Unable to sign in."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <AuthLayout>
      <div className="mb-8">
        <p className="mb-2 text-sm font-semibold text-cyan-600">Welcome back</p>
        <h1 className="text-3xl font-semibold tracking-tight text-slate-950">Sign in to your workspace</h1>
        <p className="mt-3 text-sm leading-6 text-slate-500">
          Continue monitoring your services and reliability history.
        </p>
      </div>

      {error && (
        <div role="alert" className="notice-enter mb-5 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
          {error}
        </div>
      )}

      {registrationSuccess && !error && (
        <div role="status" className="notice-enter mb-5 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
          Account created successfully. Sign in to continue.
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5">
        <FormField
          id="email"
          label="Email address"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          required
          value={email}
          onChange={(event) => setEmail(event.target.value)}
        />
        <FormField
          id="password"
          label="Password"
          type="password"
          autoComplete="current-password"
          placeholder="Enter your password"
          required
          value={password}
          onChange={(event) => setPassword(event.target.value)}
        />

        <button
          type="submit"
          disabled={submitting}
          className="primary-action group flex w-full items-center justify-center gap-2 rounded-xl bg-[#07111f] px-4 py-3.5 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {submitting ? "Signing in…" : "Sign in"}
          {!submitting && <ArrowRightIcon className="h-4 w-4 transition group-hover:translate-x-0.5" />}
        </button>
      </form>

      <p className="mt-7 text-center text-sm text-slate-500">
        New to the dashboard?{" "}
        <Link to="/register" className="font-semibold text-cyan-700 hover:text-cyan-600">
          Create an account
        </Link>
      </p>
    </AuthLayout>
  );
}
