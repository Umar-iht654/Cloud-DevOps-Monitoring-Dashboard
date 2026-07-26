import { useEffect, useRef, useState, type FormEvent } from "react";
import { useBlocker } from "react-router-dom";
import { useUnsavedChanges } from "../../context/UnsavedChangesContext";
import type { ServiceInput } from "../../types/api";
import { ArrowRightIcon, GlobeIcon } from "../ui/Icons";
import { FormField } from "../ui/FormField";

interface ServiceFormProps {
  initialValues?: ServiceInput;
  submitting: boolean;
  submitLabel: string;
  submittingLabel?: string;
  error?: string;
  onSubmit: (values: ServiceInput) => Promise<boolean>;
  onSubmitSuccess: () => void;
  onCancel: () => void;
  onValuesChange?: () => void;
}

const defaults: ServiceInput = {
  name: "",
  url: "",
  expected_status_code: 200,
  slow_threshold_ms: 750,
  check_interval_seconds: 60,
};

// This matches the backend minimum interval.
const MIN_CHECK_INTERVAL_SECONDS = 45;

type FieldErrors = Partial<Record<keyof ServiceInput, string>>;

const serviceFields = Object.keys(defaults) as (keyof ServiceInput)[];

const fieldIds: Record<keyof ServiceInput, string> = {
  name: "service-name",
  url: "service-url",
  expected_status_code: "expected-status",
  slow_threshold_ms: "slow-threshold",
  check_interval_seconds: "check-interval",
};

const UNSAVED_CHANGES_MESSAGE = "You have unsaved changes. Leave this page and discard them?";
const SAVING_IN_PROGRESS_MESSAGE =
  "Please wait for the service to finish saving before leaving this page.";

function parseNumberInput(value: string) {
  return value === "" ? Number.NaN : Number(value);
}

function displayNumberInput(value: number) {
  return Number.isNaN(value) ? "" : value;
}

function validateService(values: ServiceInput): FieldErrors {
  const errors: FieldErrors = {};

  if (!values.name.trim()) {
    errors.name = "Enter a service name.";
  }

  const trimmedUrl = values.url.trim();
  if (!trimmedUrl) {
    errors.url = "Enter the service URL.";
  } else {
    try {
      const parsedUrl = new URL(trimmedUrl);
      if (!["http:", "https:"].includes(parsedUrl.protocol) || !parsedUrl.hostname) {
        errors.url = "Use a complete URL beginning with http:// or https://.";
      }
    } catch {
      errors.url = "Use a complete URL beginning with http:// or https://.";
    }
  }

  if (
    !Number.isInteger(values.expected_status_code) ||
    values.expected_status_code < 100 ||
    values.expected_status_code > 599
  ) {
    errors.expected_status_code = "Enter a whole HTTP status code from 100 to 599.";
  }

  if (!Number.isInteger(values.slow_threshold_ms) || values.slow_threshold_ms < 1) {
    errors.slow_threshold_ms = "Enter a threshold of at least 1 millisecond.";
  }

  if (
    !Number.isInteger(values.check_interval_seconds) ||
    values.check_interval_seconds < MIN_CHECK_INTERVAL_SECONDS
  ) {
    errors.check_interval_seconds = `Enter an interval of at least ${MIN_CHECK_INTERVAL_SECONDS} seconds.`;
  }

  return errors;
}

function FieldFeedback({
  hintId,
  errorId,
  hint,
  error,
}: {
  hintId: string;
  errorId: string;
  hint: string;
  error?: string;
}) {
  return (
    <div className="mt-1.5 space-y-1 text-xs leading-5">
      <p id={hintId} className="text-slate-500">
        {hint}
      </p>
      {error && (
        <p id={errorId} className="font-medium text-rose-700">
          {error}
        </p>
      )}
    </div>
  );
}

export function ServiceForm({
  initialValues = defaults,
  submitting,
  submitLabel,
  submittingLabel = "Saving service…",
  error,
  onSubmit,
  onSubmitSuccess,
  onCancel,
  onValuesChange,
}: ServiceFormProps) {
  const [values, setValues] = useState<ServiceInput>(initialValues);
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const initialValuesRef = useRef(initialValues);
  const submitInFlightRef = useRef(false);
  const submittingRef = useRef(submitting);
  const navigationApprovedRef = useRef(false);
  const hasUnsavedChangesRef = useRef(false);
  const { registerNavigationGuard, confirmNavigation } = useUnsavedChanges();
  const hasUnsavedChanges = serviceFields.some(
    (field) => values[field] !== initialValuesRef.current[field],
  );
  hasUnsavedChangesRef.current = hasUnsavedChanges;
  submittingRef.current = submitting;
  const validationErrorCount = Object.values(fieldErrors).filter(Boolean).length;
  const blocker = useBlocker(
    ({ currentLocation, nextLocation }) =>
      (hasUnsavedChanges || submitting || submitInFlightRef.current) &&
      !navigationApprovedRef.current &&
      (currentLocation.pathname !== nextLocation.pathname ||
        currentLocation.search !== nextLocation.search ||
        currentLocation.hash !== nextLocation.hash),
  );

  useEffect(() => {
    if (!hasUnsavedChanges && !submitting && !submitInFlightRef.current) return;

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      event.preventDefault();
      event.returnValue = "";
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, [hasUnsavedChanges, submitting]);

  useEffect(
    () =>
      registerNavigationGuard(() => {
        if (submittingRef.current || submitInFlightRef.current) {
          window.alert(SAVING_IN_PROGRESS_MESSAGE);
          return false;
        }

        if (!hasUnsavedChangesRef.current) return true;
        if (!window.confirm(UNSAVED_CHANGES_MESSAGE)) return false;

        navigationApprovedRef.current = true;
        return true;
      }),
    [registerNavigationGuard],
  );

  useEffect(() => {
    if (blocker.state !== "blocked") return;

    if (submittingRef.current || submitInFlightRef.current) {
      window.alert(SAVING_IN_PROGRESS_MESSAGE);
      blocker.reset();
      return;
    }

    if (window.confirm(UNSAVED_CHANGES_MESSAGE)) {
      blocker.proceed();
    } else {
      blocker.reset();
    }
  }, [blocker]);

  const setField = <K extends keyof ServiceInput>(field: K, value: ServiceInput[K]) => {
    setValues((current) => ({ ...current, [field]: value }));
    setFieldErrors((current) => {
      if (!current[field]) return current;
      const nextErrors = { ...current };
      delete nextErrors[field];
      return nextErrors;
    });
    onValuesChange?.();
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (submitting || submitInFlightRef.current) return;

    const nextErrors = validateService(values);
    setFieldErrors(nextErrors);
    const firstInvalidField = serviceFields.find((field) => nextErrors[field]);

    if (firstInvalidField) {
      window.requestAnimationFrame(() => {
        document.getElementById(fieldIds[firstInvalidField])?.focus();
      });
      return;
    }

    submitInFlightRef.current = true;
    try {
      const succeeded = await onSubmit({
        ...values,
        name: values.name.trim(),
        url: values.url.trim(),
      });
      if (succeeded) {
        navigationApprovedRef.current = true;
        onSubmitSuccess();
      }
    } finally {
      submitInFlightRef.current = false;
    }
  };

  const handleCancel = () => {
    if (!confirmNavigation()) return;
    onCancel();
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6" noValidate aria-busy={submitting}>
      {error && (
        <div
          role="alert"
          className="notice-enter rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700"
        >
          <p className="font-semibold">We couldn&apos;t save this service.</p>
          <p className="mt-1 leading-6">{error}</p>
          <p className="mt-1 text-xs leading-5 text-rose-600">Your entries are still here so you can review and try again.</p>
        </div>
      )}

      {validationErrorCount > 0 && (
        <div
          role="alert"
          className="notice-enter rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700"
        >
          <p className="font-semibold">
            Check {validationErrorCount === 1 ? "the highlighted field" : `${validationErrorCount} highlighted fields`}.
          </p>
          <p className="mt-1 text-xs leading-5 text-rose-600">Each issue is explained directly below its field.</p>
        </div>
      )}

      <div className="premium-panel rounded-3xl p-5 sm:p-7">
        <div className="mb-6 flex items-start gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-cyan-50 text-cyan-700 ring-1 ring-cyan-100">
            <GlobeIcon className="h-5 w-5" />
          </span>
          <div>
            <p className="mb-1 text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-700">Step 01</p>
            <h2 className="font-semibold text-slate-950">Service details</h2>
            <p className="mt-1 text-sm leading-6 text-slate-500">
              Give the endpoint a clear name and provide the full public URL to monitor.
            </p>
          </div>
        </div>

        <div className="grid gap-5">
          <div className={fieldErrors.name ? "[&_input]:border-rose-400 [&_input]:bg-rose-50/30" : ""}>
            <FormField
              id="service-name"
              name="name"
              label="Service name"
              placeholder="e.g. Production API"
              autoComplete="organization"
              required
              maxLength={100}
              disabled={submitting}
              aria-invalid={Boolean(fieldErrors.name)}
              aria-describedby={`service-name-hint${fieldErrors.name ? " service-name-error" : ""}`}
              aria-errormessage={fieldErrors.name ? "service-name-error" : undefined}
              value={values.name}
              onChange={(event) => setField("name", event.target.value)}
            />
            <FieldFeedback
              hintId="service-name-hint"
              errorId="service-name-error"
              hint="Choose a short label you will recognise in dashboards and alerts."
              error={fieldErrors.name}
            />
          </div>
          <div className={fieldErrors.url ? "[&_input]:border-rose-400 [&_input]:bg-rose-50/30" : ""}>
            <FormField
              id="service-url"
              name="url"
              label="Service URL"
              type="url"
              inputMode="url"
              autoComplete="url"
              spellCheck={false}
              placeholder="https://api.example.com/health"
              required
              disabled={submitting}
              aria-invalid={Boolean(fieldErrors.url)}
              aria-describedby={`service-url-hint${fieldErrors.url ? " service-url-error" : ""}`}
              aria-errormessage={fieldErrors.url ? "service-url-error" : undefined}
              value={values.url}
              onChange={(event) => setField("url", event.target.value)}
            />
            <FieldFeedback
              hintId="service-url-hint"
              errorId="service-url-error"
              hint="Include http:// or https://. A lightweight health endpoint gives the clearest signal."
              error={fieldErrors.url}
            />
          </div>
        </div>
      </div>

      <div className="premium-panel rounded-3xl p-5 sm:p-7">
        <div className="mb-6">
          <p className="mb-1 text-[10px] font-bold uppercase tracking-[0.2em] text-violet-600">Step 02</p>
          <h2 className="font-semibold text-slate-950">Monitoring rules</h2>
          <p className="mt-1 text-sm leading-6 text-slate-500">
            These settings determine when the service is classified as healthy, slow, or down.
          </p>
        </div>

        <div className="grid gap-5 sm:grid-cols-3">
          <div className={fieldErrors.expected_status_code ? "[&_input]:border-rose-400 [&_input]:bg-rose-50/30" : ""}>
            <FormField
              id="expected-status"
              name="expected_status_code"
              label="Expected status"
              type="number"
              inputMode="numeric"
              min={100}
              max={599}
              step={1}
              required
              disabled={submitting}
              aria-invalid={Boolean(fieldErrors.expected_status_code)}
              aria-describedby={`expected-status-hint${fieldErrors.expected_status_code ? " expected-status-error" : ""}`}
              aria-errormessage={fieldErrors.expected_status_code ? "expected-status-error" : undefined}
              value={displayNumberInput(values.expected_status_code)}
              onChange={(event) => setField("expected_status_code", parseNumberInput(event.target.value))}
            />
            <FieldFeedback
              hintId="expected-status-hint"
              errorId="expected-status-error"
              hint="The HTTP response that counts as healthy, usually 200."
              error={fieldErrors.expected_status_code}
            />
          </div>
          <div className={fieldErrors.slow_threshold_ms ? "[&_input]:border-rose-400 [&_input]:bg-rose-50/30" : ""}>
            <FormField
              id="slow-threshold"
              name="slow_threshold_ms"
              label="Slow threshold (ms)"
              type="number"
              inputMode="numeric"
              min={1}
              step={1}
              required
              disabled={submitting}
              aria-invalid={Boolean(fieldErrors.slow_threshold_ms)}
              aria-describedby={`slow-threshold-hint${fieldErrors.slow_threshold_ms ? " slow-threshold-error" : ""}`}
              aria-errormessage={fieldErrors.slow_threshold_ms ? "slow-threshold-error" : undefined}
              value={displayNumberInput(values.slow_threshold_ms)}
              onChange={(event) => setField("slow_threshold_ms", parseNumberInput(event.target.value))}
            />
            <FieldFeedback
              hintId="slow-threshold-hint"
              errorId="slow-threshold-error"
              hint="Checks taking longer than this are marked slow."
              error={fieldErrors.slow_threshold_ms}
            />
          </div>
          <div className={fieldErrors.check_interval_seconds ? "[&_input]:border-rose-400 [&_input]:bg-rose-50/30" : ""}>
            <FormField
              id="check-interval"
              name="check_interval_seconds"
              label="Check interval (sec)"
              type="number"
              inputMode="numeric"
              min={MIN_CHECK_INTERVAL_SECONDS}
              step={1}
              required
              disabled={submitting}
              aria-invalid={Boolean(fieldErrors.check_interval_seconds)}
              aria-describedby={`check-interval-hint${fieldErrors.check_interval_seconds ? " check-interval-error" : ""}`}
              aria-errormessage={fieldErrors.check_interval_seconds ? "check-interval-error" : undefined}
              value={displayNumberInput(values.check_interval_seconds)}
              onChange={(event) => setField("check_interval_seconds", parseNumberInput(event.target.value))}
            />
            <FieldFeedback
              hintId="check-interval-hint"
              errorId="check-interval-error"
              hint={`The minimum check interval is ${MIN_CHECK_INTERVAL_SECONDS} seconds.`}
              error={fieldErrors.check_interval_seconds}
            />
          </div>
        </div>
      </div>

      <div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button
          type="button"
          onClick={handleCancel}
          disabled={submitting}
          className="rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 shadow-sm transition hover:border-slate-300 hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={submitting}
          className="primary-action group inline-flex items-center justify-center gap-2 rounded-xl bg-[#07111f] px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {submitting && (
            <span
              className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white motion-reduce:animate-none"
              aria-hidden="true"
            />
          )}
          <span aria-live="polite">{submitting ? submittingLabel : submitLabel}</span>
          {!submitting && <ArrowRightIcon className="h-4 w-4 transition group-hover:translate-x-0.5" />}
        </button>
      </div>
    </form>
  );
}
