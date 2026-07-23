import { useState, type FormEvent } from "react";
import type { ServiceInput } from "../../types/api";
import { ArrowRightIcon, GlobeIcon } from "../ui/Icons";
import { FormField } from "../ui/FormField";

interface ServiceFormProps {
  initialValues?: ServiceInput;
  submitting: boolean;
  submitLabel: string;
  error?: string;
  onSubmit: (values: ServiceInput) => Promise<void>;
  onCancel: () => void;
}

const defaults: ServiceInput = {
  name: "",
  url: "",
  expected_status_code: 200,
  slow_threshold_ms: 750,
  check_interval_seconds: 60,
};

export function ServiceForm({
  initialValues = defaults,
  submitting,
  submitLabel,
  error,
  onSubmit,
  onCancel,
}: ServiceFormProps) {
  const [values, setValues] = useState<ServiceInput>(initialValues);
  const [validationError, setValidationError] = useState("");

  const setField = <K extends keyof ServiceInput>(field: K, value: ServiceInput[K]) => {
    setValues((current) => ({ ...current, [field]: value }));
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setValidationError("");

    if (!values.name.trim()) {
      setValidationError("Enter a service name.");
      return;
    }

    try {
      const parsedUrl = new URL(values.url);
      if (!["http:", "https:"].includes(parsedUrl.protocol)) throw new Error();
    } catch {
      setValidationError("Enter a complete URL beginning with http:// or https://.");
      return;
    }

    if (values.check_interval_seconds < 10) {
      setValidationError("Check interval must be at least 10 seconds.");
      return;
    }

    await onSubmit({
      ...values,
      name: values.name.trim(),
      url: values.url.trim(),
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {(error || validationError) && (
        <div role="alert" className="notice-enter rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
          {validationError || error}
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
          <FormField
            id="service-name"
            label="Service name"
            placeholder="e.g. Production API"
            required
            maxLength={100}
            value={values.name}
            onChange={(event) => setField("name", event.target.value)}
          />
          <FormField
            id="service-url"
            label="Service URL"
            type="url"
            placeholder="https://api.example.com/health"
            hint="Use a dedicated health endpoint when one is available."
            required
            value={values.url}
            onChange={(event) => setField("url", event.target.value)}
          />
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
          <FormField
            id="expected-status"
            label="Expected status"
            type="number"
            min={100}
            max={599}
            required
            value={values.expected_status_code}
            onChange={(event) => setField("expected_status_code", Number(event.target.value))}
            hint="Usually 200."
          />
          <FormField
            id="slow-threshold"
            label="Slow threshold (ms)"
            type="number"
            min={1}
            required
            value={values.slow_threshold_ms}
            onChange={(event) => setField("slow_threshold_ms", Number(event.target.value))}
            hint="Above this is slow."
          />
          <FormField
            id="check-interval"
            label="Check interval (sec)"
            type="number"
            min={10}
            required
            value={values.check_interval_seconds}
            onChange={(event) => setField("check_interval_seconds", Number(event.target.value))}
            hint="Minimum 10 seconds."
          />
        </div>
      </div>

      <div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button
          type="button"
          onClick={onCancel}
          disabled={submitting}
          className="rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 shadow-sm transition hover:border-slate-300 hover:bg-slate-50 disabled:opacity-60"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={submitting}
          className="primary-action group inline-flex items-center justify-center gap-2 rounded-xl bg-[#07111f] px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-slate-900/10 focus:outline-none focus:ring-4 focus:ring-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {submitting ? "Saving…" : submitLabel}
          {!submitting && <ArrowRightIcon className="h-4 w-4 transition group-hover:translate-x-0.5" />}
        </button>
      </div>
    </form>
  );
}
