import type { HealthCheck } from "../../types/api";
import { formatDateTime, formatMilliseconds } from "../../utils/formatters";
import { StatusBadge } from "../ui/StatusBadge";

export function HealthCheckTable({ healthChecks }: { healthChecks: HealthCheck[] }) {
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
      <div className="space-y-3 sm:hidden">
        {healthChecks.map((check) => (
          <article key={check.id} className="health-row rounded-2xl border border-slate-200/80 bg-slate-50/70 p-4">
            <div className="flex items-center justify-between gap-3">
              <StatusBadge status={check.status} />
              <span className="text-xs text-slate-500">{formatDateTime(check.checked_at)}</span>
            </div>
            <div className="mt-4 grid grid-cols-2 gap-3">
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-wide text-slate-500">HTTP code</p>
                <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">{check.http_status_code ?? "—"}</p>
              </div>
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-wide text-slate-500">Response</p>
                <p className="metric-tabular mt-1 text-sm font-semibold text-slate-800">{formatMilliseconds(check.response_time_ms)}</p>
              </div>
            </div>
            {check.error_message && <p className="mt-3 rounded-lg bg-rose-50 px-3 py-2 text-xs leading-5 text-rose-700">{check.error_message}</p>}
          </article>
        ))}
      </div>

      <div className="hidden overflow-x-auto sm:block">
      <table className="w-full min-w-[680px] text-left">
        <thead>
          <tr className="border-b border-slate-200 text-[11px] font-semibold uppercase tracking-wide text-slate-500">
            <th className="pb-3 pr-4">Status</th>
            <th className="pb-3 px-4">HTTP code</th>
            <th className="pb-3 px-4">Response time</th>
            <th className="pb-3 px-4">Checked at</th>
            <th className="pb-3 pl-4">Details</th>
          </tr>
        </thead>
        <tbody>
          {healthChecks.map((check) => (
            <tr key={check.id} className="health-row border-b border-slate-100 last:border-0">
              <td className="py-4 pr-4"><StatusBadge status={check.status} /></td>
              <td className="metric-tabular px-4 py-4 text-sm font-medium text-slate-700">{check.http_status_code ?? "—"}</td>
              <td className="metric-tabular px-4 py-4 text-sm font-medium text-slate-700">{formatMilliseconds(check.response_time_ms)}</td>
              <td className="px-4 py-4 text-sm text-slate-500">{formatDateTime(check.checked_at)}</td>
              <td className="max-w-xs break-words py-4 pl-4 text-sm leading-5 text-slate-500">
                {check.error_message || "No errors"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      </div>
    </>
  );
}
