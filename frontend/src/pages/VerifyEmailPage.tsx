import { useEffect, useRef, useState, type FormEvent } from "react";
import axios from "axios";
import { Link, useLocation, useNavigate, useSearchParams } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import {
  checkVerificationSessionStatus,
  resendVerificationEmail,
  VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY,
  VERIFICATION_SESSION_TOKEN_STORAGE_KEY,
  verifyEmail,
} from "../api/auth";
import { AuthLayout } from "../components/layout/AuthLayout";
import { AlertIcon, ArrowRightIcon, CheckIcon, PulseIcon } from "../components/ui/Icons";
import { FormField } from "../components/ui/FormField";
import { useAuth } from "../context/AuthContext";

type VerificationStatus = "awaiting" | "verifying" | "verified" | "failed" | "expired";

const RESEND_COOLDOWN_SECONDS = 180;
const VERIFICATION_POLL_INTERVAL_MS = 3000;
const EXPIRED_RECOVERY_MESSAGE =
  "This verification link and sign-in session have expired. Request a new verification link to continue.";

function formatCooldown(seconds: number) {
  const safeSeconds = Math.max(0, seconds);
  const minutes = Math.floor(safeSeconds / 60);
  const remainingSeconds = safeSeconds % 60;

  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}

function isRecoverablePollingError(error: unknown) {
  if (!axios.isAxiosError(error)) return false;

  return !error.response || error.response.status >= 500;
}

interface VerificationNavigationState {
  email?: string;
  registrationStarted?: boolean;
}

export function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { completeLogin } = useAuth();
  const navigationState = location.state as VerificationNavigationState | null;
  const token = searchParams.get("token")?.trim() ?? "";
  const attemptedTokenRef = useRef<string | null>(null);
  const pollingRequestInFlightRef = useRef(false);
  const [status, setStatus] = useState<VerificationStatus>(token ? "verifying" : "awaiting");
  const [verificationMessage, setVerificationMessage] = useState("");
  const [verificationSessionToken, setVerificationSessionToken] = useState(() =>
    token ? "" : (sessionStorage.getItem(VERIFICATION_SESSION_TOKEN_STORAGE_KEY) ?? ""),
  );
  const [verificationSessionExpiresAt, setVerificationSessionExpiresAt] = useState(() =>
    token ? "" : (sessionStorage.getItem(VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY) ?? ""),
  );
  const [email, setEmail] = useState(navigationState?.email ?? "");
  const [emailError, setEmailError] = useState("");
  const [resendError, setResendError] = useState("");
  const [alreadyVerified, setAlreadyVerified] = useState(false);
  const [linkAlreadySent, setLinkAlreadySent] = useState(false);
  const [showLinkAlreadySentNotice, setShowLinkAlreadySentNotice] = useState(false);
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

  useEffect(() => {
    if (!showLinkAlreadySentNotice) return undefined;

    const timerId = window.setTimeout(() => {
      setShowLinkAlreadySentNotice(false);
    }, 5000);

    return () => window.clearTimeout(timerId);
  }, [showLinkAlreadySentNotice]);

  const cooldownRemaining = cooldownDeadline
    ? Math.max(0, Math.ceil((cooldownDeadline - Date.now()) / 1000))
    : 0;

  const clearVerificationSession = () => {
    sessionStorage.removeItem(VERIFICATION_SESSION_TOKEN_STORAGE_KEY);
    sessionStorage.removeItem(VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY);
    setVerificationSessionToken("");
    setVerificationSessionExpiresAt("");
  };

  useEffect(() => {
    if (!cooldownDeadline) return undefined;

    // Recalculate against wall-clock time so backgrounded tabs do not extend the visible cooldown.
    const refreshCooldown = () => {
      if (cooldownDeadline <= Date.now()) {
        setCooldownDeadline(null);
        setLinkAlreadySent(false);
        return;
      }
      setCooldownTick((currentTick) => currentTick + 1);
    };

    refreshCooldown();
    const timerId = window.setInterval(refreshCooldown, 1000);

    return () => window.clearInterval(timerId);
  }, [cooldownDeadline]);

  useEffect(() => {
    if (token || status !== "awaiting" || !verificationSessionToken || !verificationSessionExpiresAt) {
      return undefined;
    }

    const expireSession = () => {
      clearVerificationSession();
      setStatus("expired");
      setVerificationMessage(EXPIRED_RECOVERY_MESSAGE);
    };

    const expiryTime = new Date(verificationSessionExpiresAt).getTime();
    const expiryDelay = Math.max(0, expiryTime - Date.now());
    const timerId = window.setTimeout(expireSession, expiryDelay);

    return () => window.clearTimeout(timerId);
  }, [status, token, verificationSessionExpiresAt, verificationSessionToken]);

  useEffect(() => {
    if (token || status !== "awaiting" || !verificationSessionToken) {
      return undefined;
    }

    let isActive = true;

    const pollVerificationSession = async () => {
      if (pollingRequestInFlightRef.current) return;
      pollingRequestInFlightRef.current = true;

      try {
        const { data } = await checkVerificationSessionStatus(verificationSessionToken);
        if (!isActive) return;

        if (data.status === "pending") {
          if (data.expiresAt) {
            sessionStorage.setItem(VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY, data.expiresAt);
            setVerificationSessionExpiresAt(data.expiresAt);
          }
          return;
        }

        if (data.status === "verified" && data.token && data.user) {
          clearVerificationSession();
          completeLogin({
            message: data.message ?? "Login successful",
            token: data.token,
            user: data.user,
          });
          navigate("/dashboard", { replace: true });
          return;
        }

        clearVerificationSession();
        setStatus("expired");
        setVerificationMessage(
          data.code === "VERIFICATION_SESSION_EXPIRED"
            ? EXPIRED_RECOVERY_MESSAGE
            : (data.message ?? "This verification session can no longer be used. Request a new verification link to continue."),
        );
      } catch (requestError) {
        if (!isActive) return;

        if (isRecoverablePollingError(requestError)) {
          return;
        }

        clearVerificationSession();
        setStatus("expired");
        setVerificationMessage(
          getApiErrorMessage(requestError, "This verification session can no longer be used. Request a new verification link to continue."),
        );
      } finally {
        pollingRequestInFlightRef.current = false;
      }
    };

    pollVerificationSession();
    const intervalId = window.setInterval(pollVerificationSession, VERIFICATION_POLL_INTERVAL_MS);

    return () => {
      isActive = false;
      window.clearInterval(intervalId);
    };
  }, [completeLogin, navigate, status, token, verificationSessionToken]);

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
    setLinkAlreadySent(false);
    setShowLinkAlreadySentNotice(false);

    try {
      const { data } = await resendVerificationEmail(normalizedEmail);
      setEmail(normalizedEmail);
      setStatus("awaiting");
      setVerificationMessage("");

      if (data.code === "VERIFICATION_LINK_ALREADY_SENT") {
        setLinkAlreadySent(true);
        setShowLinkAlreadySentNotice(true);
        setCooldownDeadline(Date.now() + (data.retryAfterSeconds ?? RESEND_COOLDOWN_SECONDS) * 1000);

        if (verificationSessionToken && data.expiresAt) {
          sessionStorage.setItem(VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY, data.expiresAt);
          setVerificationSessionExpiresAt(data.expiresAt);
        }

        return;
      }

      if (data.code === "VERIFICATION_RESEND_COOLDOWN") {
        setCooldownDeadline(Date.now() + (data.retryAfterSeconds ?? RESEND_COOLDOWN_SECONDS) * 1000);
        return;
      }

      setCooldownDeadline(Date.now() + RESEND_COOLDOWN_SECONDS * 1000);

      if (data.verificationSessionToken) {
        sessionStorage.setItem(VERIFICATION_SESSION_TOKEN_STORAGE_KEY, data.verificationSessionToken);
        setVerificationSessionToken(data.verificationSessionToken);
      }

      if (data.verificationSessionExpiresAt) {
        sessionStorage.setItem(VERIFICATION_SESSION_EXPIRES_AT_STORAGE_KEY, data.verificationSessionExpiresAt);
        setVerificationSessionExpiresAt(data.verificationSessionExpiresAt);
      }
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

  const showResendForm = status === "awaiting" || status === "failed" || status === "expired";
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
              {status === "failed" || status === "expired" ? "Verification link unavailable" : "Email verification"}
            </p>
            <h1 className="text-3xl font-semibold tracking-tight text-slate-950">
              {status === "expired" ? "Request a new verification link" : status === "failed" ? "This link can’t be used" : "Check your inbox"}
            </h1>
            <p className="mt-3 text-sm leading-6 text-slate-500">
              {status === "expired"
                ? (verificationMessage || EXPIRED_RECOVERY_MESSAGE)
                : status === "failed"
                ? "The link may have expired or already been used. Request another one to continue."
                : isRegistrationStart || linkAlreadySent
                  ? "We sent a verification link to your email. Open it to finish creating your account."
                  : "Enter your email address to request another verification link."}
            </p>
            {status === "awaiting" && (
              <div className="mt-4 flex gap-3 rounded-xl border border-yellow-300 bg-yellow-100 px-4 py-3 text-sm font-medium leading-6 text-black dark:border-yellow-300 dark:bg-yellow-100 dark:text-black">
                <AlertIcon className="mt-0.5 h-5 w-5 shrink-0 text-black" />
                <span>
                  {verificationSessionToken
                    ? "This verification link expires in 3 minutes."
                    : "Verification links expire in 3 minutes."}
                </span>
              </div>
            )}
            {showLinkAlreadySentNotice && (
              <div
                role="status"
                aria-live="polite"
                className="notice-enter mt-3 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium leading-6 text-rose-800"
              >
                A verification link has already been sent to this email address.
              </div>
            )}
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
                setLinkAlreadySent(false);
                setShowLinkAlreadySentNotice(false);
              }}
            />
            {status === "awaiting" && (
              <p className="-mt-2 text-xs leading-5 text-slate-500">
                Can&apos;t find the email? Check your spam or junk folder.
              </p>
            )}
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
            {!alreadyVerified && isCooldownActive && !linkAlreadySent && (
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
