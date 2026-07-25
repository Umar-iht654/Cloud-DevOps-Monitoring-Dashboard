import { useId, useState } from "react";
import type { HealthCheck } from "../../types/api";
import { formatDateTime, formatMilliseconds } from "../../utils/formatters";
import { StatusBadge } from "../ui/StatusBadge";

const INITIAL_VISIBLE_CHECKS = 20;
const MAX_VISIBLE_CHECKS = 100;

const checkToneClasses: Record<HealthCheck["status"], { card: string; row: string }> = {
  online: {
    card: "border-slate-200/80 bg-slate-50/70",
    row: "border-slate-100",
  },
  slow: {
    card: "border-amber-200/80 bg-amber-50/60",
    row: "border-amber-100 bg-amber-50/40",
  },
  down: {
    card: "border-rose-200/80 bg-rose-50/60",
    row: "border-rose-100 bg-rose-50/40",
  },
  unknown: {
    card: "border-slate-200/80 bg-slate-50/70",
    row: "border-slate-100",
  },
};

export function HealthCheckTable({ healthChecks }: { healthChecks: HealthCheck[] }) {
  const [showAll, setShowAll] = useState(false);
  const historyId = useId();
  const cappedHealthChecks = healthChecks.slice(0, MAX_VISIBLE_CHECKS);
  const visibleHealthChecks = showAll
    ? cappedHealthChecks
    : cappedHealthChecks.slice(0, INITIAL_VISIBLE_CHECKS);
  const hasMoreChecks = cappedHealthChecks.length > INITIAL_VISIBLE_CHECKS;

  if (healthChecks.length === 0) {
    return (
      <div className="rounded-2xl bg-slate-50 px-6 py-12 text-center">
        <p className="font-medium text-slate-700">No checks recorded yet</p>
        <p className="mt-1 text-sm text-slate-500">The background worker will populate this history automatically.</p>
      </div>
    );
  }

  return (
    <>
      <div id={historyId}>
        <div className="space-y-3 sm:hidden">
          {visibleHealthChecks.map((check) => (
            <article
              key={check.id}
              className={`health-row rounded-2xl border p-4 ${checkToneClasses[check.status].card}`}
            >
              <div className="flex items-center justify-between gap-3">
                <StatusBadge status={check.status} />
                <time className="text-xs text-slate-500" dateTime={check.checked_at}>
                  {formatDateTime(check.checked_at)}
                </time>
              </div>
              <div className="mt-4 grid grid-cols-2 gap-3">
                <div>
                  <p className="text-[10px] font-semibold uppercase tracking-wide text-slate-500">HTTP code</p>
                  <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">
                    {check.http_status_code ?? "—"}
                  </p>
                </div>
                <div>
                  <p className="text-[10px] font-semibold uppercase tracking-wide text-slate-500">Response</p>
                  <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">
                    {formatMilliseconds(check.response_time_ms)}
                  </p>
                </div>
              </div>
              {check.error_message && (
                <p className="mt-3 break-words rounded-lg bg-rose-50 px-3 py-2 text-xs leading-5 text-rose-700 [overflow-wrap:anywhere]">
                  {check.error_message}
                </p>
              )}
            </article>
          ))}
        </div>

        <div className="hidden overflow-x-auto sm:block">
          <table className="w-full min-w-[680px] text-left">
            <caption className="sr-only">
              Showing {visibleHealthChecks.length.toLocaleString()} of {cappedHealthChecks.length.toLocaleString()} recent
              health checks, including status, HTTP code, response time, checked time, and error details.
            </caption>
            <thead>
              <tr className="border-b border-slate-200 text-[11px] font-semibold uppercase tracking-wide text-slate-500">
                <th className="pb-3 pr-4" scope="col">
                  Status
                </th>
                <th className="px-4 pb-3" scope="col">
                  HTTP code
                </th>
                <th className="px-4 pb-3" scope="col">
                  Response time
                </th>
                <th className="px-4 pb-3" scope="col">
                  Checked at
                </th>
                <th className="pb-3 pl-4" scope="col">
                  Details
                </th>
              </tr>
            </thead>
            <tbody>
              {visibleHealthChecks.map((check) => (
                <tr
                  key={check.id}
                  className={`health-row border-b last:border-0 ${checkToneClasses[check.status].row}`}
                >
                  <td className="py-4 pr-4">
                    <StatusBadge status={check.status} />
                  </td>
                  <td className="metric-tabular px-4 py-4 text-sm font-medium text-slate-700">
                    {check.http_status_code ?? "—"}
                  </td>
                  <td className="metric-tabular px-4 py-4 text-sm font-medium text-slate-700">
                    {formatMilliseconds(check.response_time_ms)}
                  </td>
                  <td className="px-4 py-4 text-sm text-slate-500">
                    <time dateTime={check.checked_at}>{formatDateTime(check.checked_at)}</time>
                  </td>
                  <td className="max-w-xs break-words py-4 pl-4 text-sm leading-5 text-slate-500">
                    {check.error_message || "No errors"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      {hasMoreChecks && (
        <div className="mt-5 flex justify-center">
          <button
            type="button"
            className="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-slate-300 hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500 focus-visible:ring-offset-2"
            aria-controls={historyId}
            aria-expanded={showAll}
            onClick={() => setShowAll((current) => !current)}
          >
            {showAll
              ? "Show fewer"
              : `Show more (${(cappedHealthChecks.length - INITIAL_VISIBLE_CHECKS).toLocaleString()})`}
          </button>
        </div>
      )}
    </>
  );
}
