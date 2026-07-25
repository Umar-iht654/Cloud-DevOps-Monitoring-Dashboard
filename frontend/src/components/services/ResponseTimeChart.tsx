import { useId } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { usePrefersReducedMotion } from "../../hooks/usePrefersReducedMotion";
import type { HealthCheck } from "../../types/api";

interface ResponseTimeChartProps {
  healthChecks: HealthCheck[];
  slowThresholdMs?: number;
}

export function ResponseTimeChart({ healthChecks, slowThresholdMs }: ResponseTimeChartProps) {
  const prefersReducedMotion = usePrefersReducedMotion();
  const summaryId = useId();
  const data = [...healthChecks]
    .reverse()
    .map((check) => ({
      timestamp: check.checked_at,
      responseTime: check.response_time_ms,
      status: check.status,
    }));
  const measuredSamples = data.filter((check) => check.responseTime != null);

  if (healthChecks.length === 0) {
    return (
      <div className="flex h-72 items-center justify-center rounded-2xl bg-slate-50 text-center">
        <div>
          <p className="font-medium text-slate-700">No checks recorded yet</p>
          <p className="mt-1 text-sm text-slate-500">The chart will appear after the first health check.</p>
        </div>
      </div>
    );
  }

  if (measuredSamples.length === 0) {
    return (
      <div className="flex h-72 items-center justify-center rounded-2xl bg-slate-50 px-6 text-center">
        <div>
          <p className="font-medium text-slate-700">No successful response-time samples yet</p>
          <p className="mt-1 text-sm text-slate-500">
            {healthChecks.length.toLocaleString()} {healthChecks.length === 1 ? "check is" : "checks are"} recorded,
            but none include a response-time measurement.
          </p>
        </div>
      </div>
    );
  }

  const firstDay = new Date(data[0].timestamp).toDateString();
  const spansMultipleDays = data.some((check) => new Date(check.timestamp).toDateString() !== firstDay);
  const xAxisFormatter = new Intl.DateTimeFormat(
    "en-GB",
    spansMultipleDays
      ? { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" }
      : { hour: "2-digit", minute: "2-digit" },
  );
  const summaryDateFormatter = new Intl.DateTimeFormat("en-GB", {
    dateStyle: "medium",
    timeStyle: "short",
  });
  const responseTimes = measuredSamples.map((check) => check.responseTime as number);
  const minimumResponseTime = Math.min(...responseTimes);
  const maximumResponseTime = Math.max(...responseTimes);
  const missingSampleCount = data.length - measuredSamples.length;
  const downCheckCount = data.filter((check) => check.status === "down").length;
  const summaryParts = [
    `Response time history for ${data.length.toLocaleString()} recorded ${data.length === 1 ? "check" : "checks"}, from ${summaryDateFormatter.format(new Date(data[0].timestamp))} to ${summaryDateFormatter.format(new Date(data[data.length - 1].timestamp))}.`,
    `${measuredSamples.length.toLocaleString()} ${measuredSamples.length === 1 ? "check has" : "checks have"} a response-time measurement.`,
    missingSampleCount > 0
      ? `${missingSampleCount.toLocaleString()} ${missingSampleCount === 1 ? "check has" : "checks have"} no measurement and ${missingSampleCount === 1 ? "appears" : "appear"} as a gap.`
      : "Every check has a response-time measurement.",
    minimumResponseTime === maximumResponseTime
      ? `The measured response time is ${minimumResponseTime.toLocaleString()} milliseconds.`
      : `Measured response times range from ${minimumResponseTime.toLocaleString()} to ${maximumResponseTime.toLocaleString()} milliseconds.`,
    downCheckCount > 0
      ? `${downCheckCount.toLocaleString()} ${downCheckCount === 1 ? "check has" : "checks have"} a down status.`
      : "No checks have a down status.",
    slowThresholdMs != null
      ? `The slow threshold is ${slowThresholdMs.toLocaleString()} milliseconds.`
      : "No slow threshold is shown.",
  ];

  return (
    <div>
      <ul
        className="mb-3 flex list-none flex-wrap items-center gap-x-5 gap-y-2 text-xs text-slate-600"
        aria-label="Chart legend"
      >
        <li className="inline-flex items-center gap-2">
          <span className="h-0.5 w-6 rounded-full bg-cyan-600" aria-hidden="true" />
          Response time
        </li>
        {slowThresholdMs != null && (
          <li className="inline-flex items-center gap-2">
            <span className="w-6 border-t-2 border-dashed border-amber-500" aria-hidden="true" />
            Slow threshold ({slowThresholdMs.toLocaleString()} ms)
          </li>
        )}
        {missingSampleCount > 0 && (
          <li className="inline-flex items-center gap-2">
            <span className="flex w-6 items-center gap-1" aria-hidden="true">
              <span className="h-0.5 flex-1 bg-slate-400" />
              <span className="h-0.5 flex-1 bg-slate-400" />
            </span>
            Gap: no measurement
          </li>
        )}
      </ul>
      <div
        className="h-60 w-full sm:h-72"
        role="img"
        aria-label="Response time chart"
        aria-describedby={summaryId}
      >
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart
            data={data}
            margin={{ top: 18, right: 12, bottom: 4, left: 0 }}
            accessibilityLayer={false}
          >
            <defs>
              <linearGradient id="responseTimeArea" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#06b6d4" stopOpacity={0.42} />
                <stop offset="72%" stopColor="#7c3aed" stopOpacity={0.07} />
                <stop offset="100%" stopColor="#7c3aed" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid stroke="#dbe4ef" strokeDasharray="4 5" vertical={false} />
            <XAxis
              dataKey="timestamp"
              tickFormatter={(value: string) => xAxisFormatter.format(new Date(value))}
              tick={{ fill: "#64748b", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              minTickGap={spansMultipleDays ? 44 : 28}
            />
            <YAxis
              tickFormatter={(value: number) => `${value}ms`}
              tick={{ fill: "#64748b", fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              width={64}
            />
            <Tooltip
              labelFormatter={(value) =>
                new Intl.DateTimeFormat("en-GB", {
                  dateStyle: "medium",
                  timeStyle: "medium",
                }).format(new Date(String(value)))
              }
              formatter={(value) =>
                value == null
                  ? ["No response-time measurement", "Response time"]
                  : [`${Number(value).toLocaleString()} ms`, "Response time"]
              }
              contentStyle={{
                color: "#e2e8f0",
                background: "rgba(7, 17, 31, 0.96)",
                border: "1px solid rgba(103, 232, 249, 0.16)",
                borderRadius: "14px",
                boxShadow: "0 18px 38px rgba(15, 23, 42, 0.22)",
                fontSize: "12px",
              }}
              labelStyle={{ color: "#94a3b8", marginBottom: "6px" }}
              itemStyle={{ color: "#67e8f9" }}
            />
            {slowThresholdMs != null && (
              <ReferenceLine
                y={slowThresholdMs}
                ifOverflow="extendDomain"
                stroke="#f59e0b"
                strokeDasharray="5 5"
                strokeOpacity={0.65}
                label={{
                  value: `Slow at ${slowThresholdMs} ms`,
                  position: "insideTopRight",
                  fill: "#b45309",
                  fontSize: 10,
                }}
              />
            )}
            <Area
              type="monotone"
              dataKey="responseTime"
              stroke="#0891b2"
              fill="url(#responseTimeArea)"
              fillOpacity={1}
              strokeWidth={2.75}
              dot={false}
              activeDot={{ r: 5, fill: "#0891b2", stroke: "white", strokeWidth: 2 }}
              connectNulls={false}
              isAnimationActive={!prefersReducedMotion}
              animationDuration={750}
              animationEasing="ease-out"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      <p id={summaryId} className="sr-only">
        {summaryParts.join(" ")}
      </p>
    </div>
  );
}
