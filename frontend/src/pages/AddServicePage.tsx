import { useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import { createService } from "../api/services";
import { ServiceForm } from "../components/services/ServiceForm";
import { ArrowLeftIcon } from "../components/ui/Icons";
import type { ServiceInput } from "../types/api";

export function AddServicePage() {
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const submitInFlightRef = useRef(false);

  const handleSubmit = async (values: ServiceInput) => {
    if (submitInFlightRef.current) return;
    submitInFlightRef.current = true;
    setSubmitting(true);
    setError("");

    try {
      const { data } = await createService(values);
      navigate(`/services/${data.service.id}`, {
        replace: true,
        state: { created: true },
      });
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, "Unable to create the service."));
    } finally {
      submitInFlightRef.current = false;
      setSubmitting(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-4xl px-4 py-7 sm:px-6 lg:px-8 lg:py-9">
      <Link
        to="/dashboard"
        aria-disabled={submitting}
        tabIndex={submitting ? -1 : undefined}
        onClick={(event) => {
          if (submitting) event.preventDefault();
        }}
        className={`inline-flex items-center gap-2 text-sm font-semibold text-slate-500 transition hover:text-slate-900 ${
          submitting ? "pointer-events-none opacity-50" : ""
        }`}
      >
        <ArrowLeftIcon className="h-4 w-4" />
        Back to dashboard
      </Link>

      <div className="mb-8 mt-5">
        <p className="text-sm font-semibold text-cyan-700">Monitoring setup</p>
        <h1 className="mt-1 text-3xl font-semibold tracking-[-0.03em] text-slate-950 sm:text-4xl">Add a service</h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-500">
          Add a website, API, or health endpoint. Its first check will run automatically after creation.
        </p>
      </div>

      <ServiceForm
        submitting={submitting}
        submitLabel="Start monitoring"
        submittingLabel="Starting monitoring…"
        error={error}
        onSubmit={handleSubmit}
        onCancel={() => navigate("/dashboard")}
        onValuesChange={() => setError("")}
      />
    </div>
  );
}
