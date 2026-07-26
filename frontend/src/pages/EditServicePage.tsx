import { useCallback, useEffect, useRef, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import { getService, updateService } from "../api/services";
import { ServiceForm } from "../components/services/ServiceForm";
import { ErrorState } from "../components/ui/ErrorState";
import { ArrowLeftIcon } from "../components/ui/Icons";
import { InlineLoader } from "../components/ui/InlineLoader";
import type { ServiceInput } from "../types/api";

export function EditServicePage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [initialValues, setInitialValues] = useState<ServiceInput | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [loadError, setLoadError] = useState("");
  const requestVersionRef = useRef(0);
  const submitInFlightRef = useRef(false);

  const loadService = useCallback(async () => {
    const currentRequest = ++requestVersionRef.current;
    setLoading(true);
    setLoadError("");

    if (!id) {
      setLoadError("A service ID is required to edit this service.");
      setLoading(false);
      return;
    }

    try {
      const { data } = await getService(id);
      if (currentRequest !== requestVersionRef.current) return;
      const service = data.service;
      setInitialValues({
        name: service.name,
        url: service.url,
        expected_status_code: service.expected_status_code,
        slow_threshold_ms: service.slow_threshold_ms,
        check_interval_seconds: service.check_interval_seconds,
      });
    } catch (requestError) {
      if (currentRequest !== requestVersionRef.current) return;
      setLoadError(getApiErrorMessage(requestError, "Unable to load the service."));
    } finally {
      if (currentRequest === requestVersionRef.current) setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    void loadService();
    return () => {
      requestVersionRef.current += 1;
    };
  }, [loadService]);

  const handleSubmit = async (values: ServiceInput) => {
    if (!id || submitInFlightRef.current) return false;
    submitInFlightRef.current = true;
    setSubmitting(true);
    setError("");

    try {
      await updateService(id, values);
      return true;
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, "Unable to update the service."));
      return false;
    } finally {
      submitInFlightRef.current = false;
      setSubmitting(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-4xl px-4 py-7 sm:px-6 lg:px-8 lg:py-9">
      <Link
        to={id ? `/services/${id}` : "/dashboard"}
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
        {id ? "Back to service" : "Back to dashboard"}
      </Link>

      <div className="mb-8 mt-5">
        <p className="text-sm font-semibold text-cyan-700">Service configuration</p>
        <h1 className="mt-1 text-3xl font-semibold tracking-[-0.03em] text-slate-950 sm:text-4xl">Edit service</h1>
        <p className="mt-2 text-sm leading-6 text-slate-500">Update the endpoint or adjust its monitoring thresholds.</p>
      </div>

      {loading ? (
        <InlineLoader label="Loading service" />
      ) : loadError || !initialValues ? (
        <ErrorState
          title="Service unavailable"
          message={loadError || "Service not found."}
          onRetry={id ? () => void loadService() : undefined}
        />
      ) : (
        <ServiceForm
          initialValues={initialValues}
          submitting={submitting}
          submitLabel="Save changes"
          submittingLabel="Saving changes…"
          error={error}
          onSubmit={handleSubmit}
          onSubmitSuccess={() =>
            navigate(`/services/${id}`, {
              replace: true,
              state: { updated: true },
            })
          }
          onCancel={() => navigate(`/services/${id}`)}
          onValuesChange={() => setError("")}
        />
      )}
    </div>
  );
}
