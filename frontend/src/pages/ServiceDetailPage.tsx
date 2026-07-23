import { useCallback, useEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { getApiErrorMessage } from "../api/client";
import {
  deleteService,
  getHealthChecks,
  getService,
  getServiceSummary,
} from "../api/services";
import { HealthCheckTable } from "../components/services/HealthCheckTable";
import { ResponseTimeChart } from "../components/services/ResponseTimeChart";
import { ErrorState } from "../components/ui/ErrorState";
import {
  ArrowLeftIcon,
  EditIcon,
  ExternalLinkIcon,
  RefreshIcon,
  TrashIcon,
} from "../components/ui/Icons";
import { InlineLoader } from "../components/ui/InlineLoader";
import { StatusBadge } from "../components/ui/StatusBadge";
import type { HealthCheck, Service, ServiceSummary } from "../types/api";
import {
  formatDateTime,
  formatMilliseconds,
  formatPercentage,
} from "../utils/formatters";

export function ServiceDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [service, setService] = useState<Service | null>(null);
  const [summary, setSummary] = useState<ServiceSummary | null>(null);
  const [healthChecks, setHealthChecks] = useState<HealthCheck[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState("");
  const requestVersion = useRef(0);

  const loadService = useCallback(async (silent = false) => {
    if (!id) return;
    const currentRequest = ++requestVersion.current;
    if (silent) setRefreshing(true);
    else setLoading(true);
    setError("");

    try {
      const [serviceResponse, summaryResponse, checksResponse] = await Promise.all([
        getService(id),
        getServiceSummary(id),
        getHealthChecks(id, 100),
      ]);
      if (currentRequest !== requestVersion.current) return;
      setService(serviceResponse.data.service);
      setSummary(summaryResponse.data.summary);
      setHealthChecks(checksResponse.data.health_checks);
    } catch (requestError) {
      if (currentRequest !== requestVersion.current) return;
      setError(getApiErrorMessage(requestError, "Unable to load this service."));
    } finally {
      if (currentRequest === requestVersion.current) {
        setLoading(false);
        setRefreshing(false);
      }
    }
  }, [id]);

  useEffect(() => {
    void loadService();
    const interval = window.setInterval(() => void loadService(true), 30_000);
    return () => {
      window.clearInterval(interval);
      requestVersion.current += 1;
    };
  }, [loadService]);

  const handleDelete = async () => {
    if (!id || !service) return;
    const confirmed = window.confirm(
      `Delete ${service.name}? Its complete health-check history will also be removed.`,
    );
    if (!confirmed) return;

    setDeleting(true);
    try {
      await deleteService(id);
      navigate("/dashboard", { replace: true });
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, "Unable to delete the service."));
      setDeleting(false);
    }
  };

  const successMessage = (location.state as { created?: boolean; updated?: boolean } | null)?.created
    ? "Service added. Its first health check will run shortly."
    : (location.state as { updated?: boolean } | null)?.updated
      ? "Settings updated. Status will refresh after the next check."
      : "";

  if (loading) {
    return <div className="mx-auto max-w-[1500px] px-4 py-9 sm:px-6 lg:px-8"><InlineLoader label="Loading service history" /></div>;
  }

  if (error && !service) {
    return <div className="mx-auto max-w-4xl px-4 py-9 sm:px-6"><ErrorState message={error} onRetry={() => void loadService()} /></div>;
  }

  if (!service || !summary) return null;

  return (
    <div className="mx-auto w-full max-w-[1500px] px-4 py-7 sm:px-6 lg:px-8 lg:py-9">
      <Link to="/dashboard" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 transition hover:text-slate-900">
        <ArrowLeftIcon className="h-4 w-4" />
        Back to dashboard
      </Link>

      {successMessage && (
        <div role="status" className="notice-enter mt-5 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
          {successMessage}
        </div>
      )}
      {error && (
        <div role="alert" className="notice-enter mt-5 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</div>
      )}

      <header className="detail-hero mt-5 flex flex-col gap-6 rounded-[1.75rem] p-6 text-white lg:flex-row lg:items-center lg:justify-between lg:p-8">
        <div className="relative z-10 min-w-0">
          <div className="mb-3 flex flex-wrap items-center gap-3">
            <StatusBadge status={service.current_status} live />
            <span className="text-xs text-slate-400">Checking every {service.check_interval_seconds}s</span>
          </div>
          <h1 className="break-words text-3xl font-semibold tracking-[-0.04em] text-white sm:text-4xl">{service.name}</h1>
          <a href={service.url} target="_blank" rel="noreferrer" className="mt-2 inline-flex max-w-full items-center gap-2 truncate text-sm text-cyan-300 transition hover:text-cyan-200">
            <span className="truncate">{service.url}</span>
            <ExternalLinkIcon className="h-4 w-4 shrink-0" />
          </a>
        </div>
        <div className="relative z-10 flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => void loadService(true)}
            disabled={refreshing}
            className="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/[0.06] px-3.5 py-2.5 text-sm font-semibold text-slate-200 backdrop-blur-sm transition hover:bg-white/10 disabled:opacity-60"
          >
            <RefreshIcon className={`h-4 w-4 ${refreshing ? "animate-spin" : ""}`} />
            Refresh
          </button>
          <Link to={`/services/${id}/edit`} className="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/[0.06] px-3.5 py-2.5 text-sm font-semibold text-slate-200 backdrop-blur-sm transition hover:bg-white/10">
            <EditIcon className="h-4 w-4" />
            Edit
          </Link>
          <button
            type="button"
            onClick={() => void handleDelete()}
            disabled={deleting}
            className="inline-flex items-center gap-2 rounded-xl border border-rose-300/20 bg-rose-400/10 px-3.5 py-2.5 text-sm font-semibold text-rose-200 backdrop-blur-sm transition hover:bg-rose-400/15 disabled:opacity-60"
          >
            <TrashIcon className="h-4 w-4" />
            {deleting ? "Deleting…" : "Delete"}
          </button>
        </div>
      </header>

      <section className="mt-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        {[
          ["Uptime", summary.total_checks > 0 ? formatPercentage(summary.uptime_percentage) : "No data"],
          ["Average response", summary.total_checks > 0 ? formatMilliseconds(summary.average_response_time_ms) : "No data"],
          ["Total checks", summary.total_checks.toLocaleString()],
          ["Failed checks", summary.failed_checks.toLocaleString()],
          ["Last downtime", summary.last_down_at ? formatDateTime(summary.last_down_at) : "None recorded"],
        ].map(([label, value], index) => (
          <article key={label} className="detail-metric premium-surface rounded-2xl border border-slate-200/80 bg-white p-5" style={{ animationDelay: `${index * 55}ms` }}>
            <p className="text-xs font-medium uppercase tracking-wide text-slate-500">{label}</p>
            <p className="metric-tabular mt-2 text-lg font-semibold text-slate-900">{value}</p>
          </article>
        ))}
      </section>

      <section className="premium-panel mt-5 rounded-3xl p-5 sm:p-6">
        <div className="mb-5 flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 className="text-lg font-semibold text-slate-950">Response time</h2>
            <p className="mt-1 text-sm text-slate-500">Most recent {healthChecks.length} checks, shown chronologically.</p>
          </div>
          <p className="text-xs text-slate-500">Last checked {formatDateTime(summary.last_checked_at)}</p>
        </div>
        <ResponseTimeChart healthChecks={healthChecks} slowThresholdMs={service.slow_threshold_ms} />
      </section>

      <section className="premium-panel mt-5 rounded-3xl p-5 sm:p-6">
        <div className="mb-5">
          <h2 className="text-lg font-semibold text-slate-950">Recent health checks</h2>
          <p className="mt-1 text-sm text-slate-500">HTTP results, timing, and failure details from the monitoring worker.</p>
        </div>
        <HealthCheckTable healthChecks={healthChecks} />
      </section>
    </div>
  );
}
