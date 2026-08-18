import { Link } from "react-router-dom";
import type { OverviewServiceReport } from "../../types/api";
import { formatMilliseconds, formatPercentage } from "../../utils/formatters";
import { servicePath } from "../../utils/serviceRoutes";

interface ServiceReliabilityTableProps {
  services: OverviewServiceReport[];
}

export function ServiceReliabilityTable({ services }: ServiceReliabilityTableProps) {
  if (services.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-200 bg-slate-50 px-5 py-10 text-center">
        <p className="font-medium text-slate-800">No service summaries yet</p>
        <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-slate-500">
          This table will appear after a monitored service has completed an hourly summary.
        </p>
      </div>
    );
  }

  return (
    <>
      <div className="space-y-3 md:hidden">
        {services.map((service) => (
          <article key={service.service_id} className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
            <div className="flex items-start justify-between gap-4">
              <Link
                to={servicePath(service.service_id, service.service_name)}
                className="min-w-0 break-words font-semibold text-slate-950 underline decoration-cyan-300/70 decoration-2 underline-offset-4 transition hover:text-cyan-700"
              >
                {service.service_name}
              </Link>
              <span className="metric-tabular shrink-0 rounded-full bg-cyan-50 px-2.5 py-1 text-xs font-semibold text-cyan-800">
                {formatPercentage(service.uptime_percentage)} uptime
              </span>
            </div>
            <dl className="mt-4 grid grid-cols-3 gap-3 border-t border-slate-100 pt-4 text-sm">
              <div>
                <dt className="text-xs text-slate-500">Checks</dt>
                <dd className="metric-tabular mt-1 font-semibold text-slate-800">{service.total_checks.toLocaleString()}</dd>
              </div>
              <div>
                <dt className="text-xs text-slate-500">Failed</dt>
                <dd className="metric-tabular mt-1 font-semibold text-rose-700">{service.failed_checks.toLocaleString()}</dd>
              </div>
              <div>
                <dt className="text-xs text-slate-500">Avg. time</dt>
                <dd className="metric-tabular mt-1 font-semibold text-slate-800">{formatMilliseconds(service.average_response_time_ms)}</dd>
              </div>
            </dl>
          </article>
        ))}
      </div>

      <div
        className="hidden overflow-x-auto rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-600 focus-visible:ring-offset-2 md:block"
        role="region"
        aria-label="Scrollable service reliability table"
        tabIndex={0}
      >
        <table className="w-full min-w-[720px] text-left">
          <caption className="sr-only">
            Service reliability report including uptime, total checks, failed checks and average response time.
          </caption>
          <thead>
            <tr className="border-b border-slate-200 text-[11px] font-semibold uppercase tracking-wide text-slate-500">
              <th scope="col" className="px-4 py-3">Service</th>
              <th scope="col" className="px-4 py-3 text-right">Uptime</th>
              <th scope="col" className="px-4 py-3 text-right">Checks</th>
              <th scope="col" className="px-4 py-3 text-right">Failed</th>
              <th scope="col" className="px-4 py-3 text-right">Avg. response</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 text-sm">
            {services.map((service) => (
              <tr key={service.service_id} className="transition hover:bg-slate-50/80">
                <th scope="row" className="max-w-[22rem] px-4 py-4 text-left font-medium">
                  <Link
                    to={servicePath(service.service_id, service.service_name)}
                    className="break-words text-slate-900 underline decoration-cyan-300/70 decoration-2 underline-offset-4 transition hover:text-cyan-700"
                  >
                    {service.service_name}
                  </Link>
                </th>
                <td className="metric-tabular px-4 py-4 text-right font-semibold text-cyan-800">{formatPercentage(service.uptime_percentage)}</td>
                <td className="metric-tabular px-4 py-4 text-right text-slate-700">{service.total_checks.toLocaleString()}</td>
                <td className="metric-tabular px-4 py-4 text-right text-rose-700">{service.failed_checks.toLocaleString()}</td>
                <td className="metric-tabular px-4 py-4 text-right text-slate-700">{formatMilliseconds(service.average_response_time_ms)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
