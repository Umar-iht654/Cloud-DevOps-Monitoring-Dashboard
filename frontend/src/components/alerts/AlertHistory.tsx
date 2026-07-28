import { Link } from "react-router-dom";
import type { Alert } from "../../types/api";
import { formatDateTime } from "../../utils/formatters";
import { AlertIcon, ArrowRightIcon } from "../ui/Icons";

interface AlertHistoryProps {
  alerts: Alert[];
  emptyTitle?: string;
  emptyMessage?: string;
  showService?: boolean;
}

const severityStyles: Record<string, string> = {
  critical: "border-rose-200 bg-rose-50 text-rose-700",
  warning: "border-amber-200 bg-amber-50 text-amber-700",
  info: "border-cyan-200 bg-cyan-50 text-cyan-700",
};

function formatAlertType(value: string) {
  return value
    .split("_")
    .filter(Boolean)
    .map((part) => `${part[0]?.toUpperCase() ?? ""}${part.slice(1)}`)
    .join(" ");
}

function serviceName(alert: Alert) {
  return alert.service?.name || `Service #${alert.service_id}`;
}

function SeverityBadge({ severity }: { severity: string }) {
  const normalizedSeverity = severity.toLowerCase();

  return (
    <span
      className={`inline-flex rounded-full border px-2.5 py-1 text-[11px] font-semibold capitalize ${
        severityStyles[normalizedSeverity] ??
        "border-slate-200 bg-slate-100 text-slate-700"
      }`}
    >
      {severity || "Unknown"}
    </span>
  );
}

export function AlertHistory({
  alerts,
  emptyTitle = "No alerts recorded",
  emptyMessage = "Downtime transitions will appear here when a monitored service needs attention.",
  showService = true,
}: AlertHistoryProps) {
  if (alerts.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-200 bg-slate-50/70 px-6 py-12 text-center">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-50 text-emerald-700 ring-1 ring-emerald-100">
          <AlertIcon className="h-6 w-6" />
        </div>
        <p className="mt-4 font-semibold text-slate-800">{emptyTitle}</p>
        <p className="mx-auto mt-1 max-w-lg text-sm leading-6 text-slate-500">
          {emptyMessage}
        </p>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-3 md:hidden">
        {alerts.map((alert) => (
          <article
            key={alert.id}
            className="rounded-2xl border border-slate-200/80 bg-white p-4 shadow-sm"
          >
            <div className="flex flex-wrap items-center justify-between gap-3">
              <SeverityBadge severity={alert.severity} />
              <time className="text-xs text-slate-500" dateTime={alert.created_at}>
                {formatDateTime(alert.created_at)}
              </time>
            </div>
            <h3 className="mt-4 break-words font-semibold text-slate-950">
              {alert.title}
            </h3>
            <p className="mt-2 break-words text-sm leading-6 text-slate-600 [overflow-wrap:anywhere]">
              {alert.message}
            </p>
            <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-slate-100 pt-3">
              <span className="text-xs font-medium text-slate-500">
                {formatAlertType(alert.type)}
              </span>
              {showService && (
                <Link
                  to={`/services/${alert.service_id}`}
                  className="inline-flex min-w-0 items-center gap-1.5 text-sm font-semibold text-cyan-700 transition hover:text-cyan-900"
                >
                  <span className="max-w-52 truncate">{serviceName(alert)}</span>
                  <ArrowRightIcon className="h-4 w-4 shrink-0" />
                </Link>
              )}
            </div>
          </article>
        ))}
      </div>

      <div className="hidden overflow-x-auto md:block">
        <table className="w-full min-w-[760px] text-left">
          <caption className="sr-only">
            Recent alert history including severity, service, type, message and created time.
          </caption>
          <thead>
            <tr className="border-b border-slate-200 text-[11px] font-semibold uppercase tracking-wide text-slate-500">
              <th className="pb-3 pr-4" scope="col">
                Severity
              </th>
              {showService && (
                <th className="px-4 pb-3" scope="col">
                  Service
                </th>
              )}
              <th className="px-4 pb-3" scope="col">
                Alert
              </th>
              <th className="px-4 pb-3" scope="col">
                Message
              </th>
              <th className="pb-3 pl-4" scope="col">
                Created
              </th>
            </tr>
          </thead>
          <tbody>
            {alerts.map((alert) => (
              <tr
                key={alert.id}
                className="border-b border-slate-100 align-top last:border-0 hover:bg-slate-50/70"
              >
                <td className="py-4 pr-4">
                  <SeverityBadge severity={alert.severity} />
                </td>
                {showService && (
                  <td className="px-4 py-4">
                    <Link
                      to={`/services/${alert.service_id}`}
                      className="inline-flex max-w-48 items-center gap-1.5 font-semibold text-cyan-700 transition hover:text-cyan-900"
                    >
                      <span className="truncate">{serviceName(alert)}</span>
                      <ArrowRightIcon className="h-4 w-4 shrink-0" />
                    </Link>
                  </td>
                )}
                <td className="px-4 py-4">
                  <p className="max-w-48 break-words text-sm font-semibold text-slate-900">
                    {alert.title}
                  </p>
                  <p className="mt-1 text-xs font-medium text-slate-500">
                    {formatAlertType(alert.type)}
                  </p>
                </td>
                <td className="max-w-md break-words px-4 py-4 text-sm leading-6 text-slate-600 [overflow-wrap:anywhere]">
                  {alert.message}
                </td>
                <td className="whitespace-nowrap py-4 pl-4 text-sm text-slate-500">
                  <time dateTime={alert.created_at}>
                    {formatDateTime(alert.created_at)}
                  </time>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
