import { useEffect, useRef, useState, type FormEvent } from "react";
import axios from "axios";
import { Link, useLocation, useSearchParams } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import { resendVerificationEmail, verifyEmail } from "../api/auth";
import { AuthLayout } from "../components/layout/AuthLayout";
import { AlertIcon, ArrowRightIcon, CheckIcon, PulseIcon } from "../components/ui/Icons";
import { FormField } from "../components/ui/FormField";

type VerificationStatus = "awaiting" | "verifying" | "verified" | "failed";

const RESEND_COOLDOWN_SECONDS = 120;

function formatCooldown(seconds: number) {
  const safeSeconds = Math.max(0, seconds);
  const minutes = Math.floor(safeSeconds / 60);
  const remainingSeconds = safeSeconds % 60;

  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}

interface VerificationNavigationState {
  email?: string;
  registrationStarted?: boolean;
}

export function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigationState = location.state as VerificationNavigationState | null;
  const token = searchParams.get("token")?.trim() ?? "";
  const attemptedTokenRef = useRef<string | null>(null);
  const [status, setStatus] = useState<VerificationStatus>(token ? "verifying" : "awaiting");
  const [verificationMessage, setVerificationMessage] = useState("");
  const [email, setEmail] = useState(navigationState?.email ?? "");
  const [emailError, setEmailError] = useState("");
  const [resendError, setResendError] = useState("");
  const [alreadyVerified, setAlreadyVerified] = useState(false);
  const [resending, setResending] = useState(false);
  const [cooldownDeadline, setCooldownDeadline] = useState<number | null>(
    !token && navigationState?.registrationStarted ? Date.now() + RESEND_COOLDOWN_SECONDS * 1000 : null,
  );
  const [, setCooldownTick] = useState(0);
  const errorRef = useRef<HTMLDivElement>(null);
  const resendFormRef = useRef<HTMLFormElement>(null);

  useEffect(() => {
    if (status === "failed" || resendError) {
      errorRef.current?.focus();
    }
  }, [resendError, status]);

  const cooldownRemaining = cooldownDeadline
    ? Math.max(0, Math.ceil((cooldownDeadline - Date.now()) / 1000))
    : 0;

  useEffect(() => {
    if (!cooldownDeadline) return undefined;

    // Recalculate against wall-clock time so backgrounded tabs do not extend the visible cooldown.
    const refreshCooldown = () => {
      if (cooldownDeadline <= Date.now()) {
        setCooldownDeadline(null);
        return;
      }
      setCooldownTick((currentTick) => currentTick + 1);
    };

    refreshCooldown();
    const timerId = window.setInterval(refreshCooldown, 1000);

    return () => window.clearInterval(timerId);
  }, [cooldownDeadline]);

  useEffect(() => {
    if (!token) {
      setStatus("awaiting");
      setVerificationMessage("");
      return;
    }

    if (attemptedTokenRef.current === token) return;
    attemptedTokenRef.current = token;
    setStatus("verifying");
    setVerificationMessage("");

    verifyEmail(token)
      .then(({ data }) => {
        setStatus("verified");
        setVerificationMessage(data.message);
      })
      .catch((requestError) => {
        setStatus("failed");
        setVerificationMessage(
          getApiErrorMessage(requestError, "We could not verify this email link. Request another one and try again."),
        );
      });
  }, [token]);

  const handleResend = async (event: FormEvent) => {
    event.preventDefault();
    if (resending || cooldownRemaining > 0 || alreadyVerified) return;

    const normalizedEmail = email.trim().toLowerCase();
    if (!normalizedEmail) {
      setEmailError("Enter your email address.");
      setResendError("");
      requestAnimationFrame(() => resendFormRef.current?.querySelector<HTMLInputElement>("input")?.focus());
      return;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalizedEmail)) {
      setEmailError("Enter a valid email address.");
      setResendError("");
      requestAnimationFrame(() => resendFormRef.current?.querySelector<HTMLInputElement>("input")?.focus());
      return;
    }

    setResending(true);
    setEmailError("");
    setResendError("");
    setAlreadyVerified(false);

    try {
      await resendVerificationEmail(normalizedEmail);
      setEmail(normalizedEmail);
      setCooldownDeadline(Date.now() + RESEND_COOLDOWN_SECONDS * 1000);
    } catch (requestError) {
      const responseCode = axios.isAxiosError(requestError)
        ? (requestError.response?.data as { code?: string } | undefined)?.code
        : undefined;

      if (responseCode === "EMAIL_ALREADY_VERIFIED") {
        setEmail(normalizedEmail);
        setAlreadyVerified(true);
        setCooldownDeadline(null);
        setResendError("");
        return;
      }

      setResendError(getApiErrorMessage(requestError, "Unable to request another verification email."));
    } finally {
      setResending(false);
    }
  };

  const showResendForm = status === "awaiting" || status === "failed";
  const isRegistrationStart = status === "awaiting" && navigationState?.registrationStarted;
  const isCooldownActive = cooldownRemaining > 0;
  const resendButtonDisabled = resending || isCooldownActive || alreadyVerified;
  const resendButtonLabel = resending
    ? "Requesting link..."
    : isCooldownActive
      ? `Resend available in ${formatCooldown(cooldownRemaining)}`
      : "Send another verification link";

  return (
    <AuthLayout>
      {status === "verifying" && (
        <section className="py-8 text-center" aria-busy="true">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-50 text-cyan-700 ring-1 ring-cyan-100">
            <PulseIcon className="h-7 w-7 animate-pulse" />
          </div>
          <p className="mt-6 text-sm font-semibold text-cyan-700">Email verification</p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-slate-950">Verifying your email</h1>
          <p role="status" aria-live="polite" className="mt-3 text-sm leading-6 text-slate-500">
            This only takes a moment. Please keep this page open.
          </p>
        </section>
      )}

      {status === "verified" && (
        <section className="py-3 text-center">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-emerald-50 text-emerald-700 ring-1 ring-emerald-100">
            <CheckIcon className="h-7 w-7" />
          </div>
          <p className="mt-6 text-sm font-semibold text-emerald-700">Email verified</p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight text-slate-950">Your account is ready</h1>
          <p role="status" aria-live="polite" className="mt-3 text-sm leading-6 text-slate-500">
            {verificationMessage || "Email verified successfully. You can now log in."}
          </p>
          <Link
            to="/login"
            state={{ verificationSuccess: true }}
            className="primary-action mt-8 flex min-h-12 w-full items-center justify-center gap-2 rounded-xl bg-[#07111f] px-4 py-3.5 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20"
          >
            Continue to sign in
            <ArrowRightIcon className="h-4 w-4" />
          </Link>
        </section>
      )}

      {showResendForm && (
        <section>
          <div className="mb-7">
            <p className="mb-2 text-sm font-semibold text-cyan-700">
              {status === "failed" ? "Verification link unavailable" : "Email verification"}
            </p>
            <h1 className="text-3xl font-semibold tracking-tight text-slate-950">
              {status === "failed" ? "This link can’t be used" : "Check your inbox"}
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-500">
              {status === "failed"
                ? "The link may have expired or already been used. Request another one to continue."
                : isRegistrationStart
                  ? "We sent a verification link to your email. Open it to finish creating your account."
                  : "Enter your email address to request another verification link."}
            </p>
          </div>

          {status === "failed" && (
            <div
              ref={errorRef}
              role="alert"
              tabIndex={-1}
              className="notice-enter mb-5 flex gap-3 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm leading-6 text-amber-900 focus:outline-none"
            >
              <AlertIcon className="mt-0.5 h-5 w-5 shrink-0" />
              <span>{verificationMessage}</span>
            </div>
          )}

          {resendError && (
            <div
              ref={errorRef}
              role="alert"
              tabIndex={-1}
              className="notice-enter mb-5 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700 focus:outline-none"
            >
              {resendError}
            </div>
          )}

          {alreadyVerified && (
            <div
              role="status"
              aria-live="polite"
              className="notice-enter mb-5 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm leading-6 text-emerald-800"
            >
              This email has already been verified.{" "}
              <Link
                to="/login"
                className="rounded font-semibold text-emerald-900 underline decoration-emerald-300 underline-offset-4 hover:text-emerald-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-600 focus-visible:ring-offset-2"
              >
                Sign in
              </Link>
            </div>
          )}

          <form ref={resendFormRef} onSubmit={handleResend} noValidate aria-busy={resending} className="space-y-4">
            <FormField
              id="verification-email"
              label="Email address"
              type="email"
              autoComplete="email"
              autoFocus={!token}
              placeholder="you@example.com"
              required
              disabled={resending}
              error={emailError}
              value={email}
              onChange={(event) => {
                setEmail(event.target.value);
                setEmailError("");
                setResendError("");
                setAlreadyVerified(false);
              }}
            />
            {!alreadyVerified && (
              <button
                type="submit"
                disabled={resendButtonDisabled}
                className="primary-action flex min-h-12 w-full items-center justify-center gap-2 rounded-xl bg-[#07111f] px-4 py-3.5 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {resendButtonLabel}
                {!resending && !isCooldownActive && <ArrowRightIcon className="h-4 w-4" />}
              </button>
            )}
            {!alreadyVerified && isCooldownActive && (
              <p role="status" aria-live="polite" className="text-sm leading-6 text-slate-500">
                If the address is eligible for another verification email, it will be sent.
                You can request again when the timer ends.
              </p>
            )}
          </form>

          <p className="mt-7 text-center text-sm text-slate-500">
            Already verified?{" "}
            <Link
              to="/login"
              className="rounded font-semibold text-cyan-700 hover:text-cyan-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-600 focus-visible:ring-offset-2"
            >
              Sign in
            </Link>
          </p>
        </section>
      )}
    </AuthLayout>
  );
}
