import type { CSSProperties } from "react";
import { Link } from "react-router-dom";
import type { Service, ServiceSummary } from "../../types/api";
import {
  formatMilliseconds,
  formatPercentage,
  formatRelativeTime,
  hostnameFromUrl,
} from "../../utils/formatters";
import { servicePath } from "../../utils/serviceRoutes";
import { ArrowRightIcon, ClockIcon, ExternalLinkIcon, GlobeIcon } from "../ui/Icons";
import { StatusBadge } from "../ui/StatusBadge";
import { positionSpotlight } from "../../utils/motion";

interface ServiceCardProps {
  service: Service;
  summary?: ServiceSummary;
  summaryError?: boolean;
  summaryLoading?: boolean;
  index?: number;
}

const statusAccents = {
  online: "#10b981",
  slow: "#f59e0b",
  down: "#f43f5e",
  unknown: "#94a3b8",
};

export function ServiceCard({
  service,
  summary,
  summaryError = false,
  summaryLoading = false,
  index = 0,
}: ServiceCardProps) {
  const hasData = Boolean(summary && summary.total_checks > 0 && !summaryError);
  const uptime = hasData ? Math.max(0, Math.min(100, summary?.uptime_percentage ?? 0)) : 0;
  const accent = statusAccents[service.current_status];
  const uptimeAccent = uptime >= 99 ? "#10b981" : uptime >= 95 ? "#f59e0b" : "#f43f5e";
  const cardStyle = {
    "--service-accent": accent,
    animationDelay: `${Math.min(index, 8) * 55}ms`,
  } as CSSProperties;
  const uptimeRing = {
    background: hasData
      ? `conic-gradient(${uptimeAccent} ${uptime * 3.6}deg, #e2e8f0 ${uptime * 3.6}deg)`
      : "conic-gradient(#cbd5e1 28deg, #e2e8f0 28deg)",
  };

  return (
    <article
      className="premium-surface service-card group flex h-full flex-col rounded-2xl border border-slate-200/80 bg-white p-5 shadow-[0_1px_2px_rgba(15,23,42,0.03)]"
      style={cardStyle}
      onPointerMove={positionSpotlight}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex min-w-0 items-start gap-3">
          <span className="service-icon flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-slate-100 text-slate-600 ring-1 ring-inset ring-slate-200">
            <GlobeIcon className="h-5 w-5" />
          </span>
          <div className="min-w-0">
            <h3 className="line-clamp-2 font-semibold text-slate-950" title={service.name}>
              {service.name}
            </h3>
            <a
              href={service.url}
              target="_blank"
              rel="noreferrer"
              aria-label={`Open ${service.name} in a new tab`}
              className="mt-1 flex max-w-full items-center gap-1.5 truncate text-xs text-slate-500 hover:text-cyan-700"
              onClick={(event) => event.stopPropagation()}
            >
              <span className="truncate">{hostnameFromUrl(service.url)}</span>
              <ExternalLinkIcon className="h-3.5 w-3.5 shrink-0" />
            </a>
          </div>
        </div>
        <StatusBadge status={service.current_status} live />
      </div>

      <div className="mt-5 flex items-center gap-5 border-y border-slate-100 py-4">
        <div className="relative h-20 w-20 shrink-0 rounded-full p-[7px]" style={uptimeRing}>
          <div className="flex h-full w-full flex-col items-center justify-center rounded-full bg-white shadow-inner">
            <p className="metric-tabular text-sm font-bold tracking-tight text-slate-900">
              {summaryLoading && !summary
                ? "…"
                : hasData
                  ? formatPercentage(summary?.uptime_percentage ?? 0)
                  : "—"}
            </p>
            <p className="text-[9px] font-semibold uppercase tracking-[0.14em] text-slate-500">Uptime</p>
          </div>
        </div>
        <div className="grid min-w-0 flex-1 grid-cols-2 gap-4">
          <div>
            <p className="text-[11px] font-medium uppercase tracking-wide text-slate-500">Average</p>
            <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">
              {summaryError
                ? "Unavailable"
                : summaryLoading && !summary
                  ? "Loading…"
                : summary && summary.total_checks > 0
                  ? formatMilliseconds(summary.average_response_time_ms)
                  : "No data"}
            </p>
          </div>
          <div>
            <p className="text-[11px] font-medium uppercase tracking-wide text-slate-500">Checks</p>
            <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">
              {summaryError
                ? "Unavailable"
                : summaryLoading && !summary
                  ? "Loading…"
                  : (summary?.total_checks.toLocaleString() ?? "No data")}
            </p>
          </div>
        </div>
      </div>

      <div className="mt-auto flex items-center justify-between gap-3 pt-4">
        <span className="flex min-w-0 items-center gap-1.5 truncate text-xs text-slate-500">
          <ClockIcon className="h-3.5 w-3.5 shrink-0" />
          {service.current_status === "unknown"
            ? "Pending first check"
            : summaryLoading && !summary
              ? "Loading metrics"
            : summaryError
              ? "Metrics unavailable"
              : formatRelativeTime(summary?.last_checked_at)}
        </span>
        <Link
          to={servicePath(service.id, service.name)}
          className="inline-flex shrink-0 items-center gap-1.5 rounded-lg px-1 py-1 text-xs font-semibold text-cyan-700 transition hover:text-cyan-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-700/30"
        >
          View details
          <ArrowRightIcon className="service-arrow h-3.5 w-3.5" />
        </Link>
      </div>
    </article>
  );
}
