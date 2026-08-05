import { useId, useMemo } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { OverviewDailyReportPoint } from "../../types/api";
import { usePrefersReducedMotion } from "../../hooks/usePrefersReducedMotion";

type TrendMetric = "uptime" | "responseTime" | "failedChecks";

interface ReportTrendChartProps {
  title: string;
  description: string;
  data: OverviewDailyReportPoint[];
  metric: TrendMetric;
}

const metricConfig = {
  uptime: {
    dataKey: "uptime",
    label: "Uptime",
    stroke: "#0891b2",
    fillId: "uptime",
    unit: "%",
  },
  responseTime: {
    dataKey: "responseTime",
    label: "Average response time",
    stroke: "#7c3aed",
    fillId: "response-time",
    unit: "ms",
  },
  failedChecks: {
    dataKey: "failedChecks",
    label: "Failed checks",
    stroke: "#e11d48",
    fillId: "failed-checks",
    unit: "",
  },
} as const;

const axisDateFormatter = new Intl.DateTimeFormat("en-GB", {
  day: "numeric",
  month: "short",
});

const tooltipDateFormatter = new Intl.DateTimeFormat("en-GB", {
  dateStyle: "medium",
});

export function ReportTrendChart({
  title,
  description,
  data,
  metric,
}: ReportTrendChartProps) {
  const summaryId = useId();
  const prefersReducedMotion = usePrefersReducedMotion();
  const config = metricConfig[metric];
  const chartData = useMemo(
    () =>
      data.map((point) => ({
        timestamp: point.period_start,
        uptime: point.uptime_percentage,
        responseTime:
          point.response_time_sample_count > 0
            ? point.average_response_time_ms
            : null,
        failedChecks: point.failed_checks,
      })),
    [data],
  );

  const values = chartData
    .map((point) => point[config.dataKey])
    .filter((value): value is number => typeof value === "number");
  const latestValue = values.at(-1);
  const highestValue = values.length > 0 ? Math.max(...values) : null;
  const lowestValue = values.length > 0 ? Math.min(...values) : null;
  const formatValue = (value: number) =>
    metric === "uptime"
      ? `${value.toFixed(value % 1 === 0 ? 0 : 2)}%`
      : metric === "responseTime"
        ? `${value.toLocaleString()} ms`
        : value.toLocaleString();

  const summary =
    values.length === 0
      ? `${title} has no recorded values in this report range.`
      : `${title} has ${values.length} recorded ${values.length === 1 ? "value" : "values"}. The latest value is ${formatValue(latestValue ?? 0)}, with a range from ${formatValue(lowestValue ?? 0)} to ${formatValue(highestValue ?? 0)}.`;

  return (
    <article className="premium-surface rounded-2xl border border-slate-200/80 bg-white p-5 shadow-[0_1px_2px_rgba(15,23,42,0.03)] sm:p-6">
      <div>
        <h2 className="text-base font-semibold tracking-[-0.02em] text-slate-950">{title}</h2>
        <p className="mt-1 text-sm leading-6 text-slate-500">{description}</p>
      </div>

      {chartData.length === 0 ? (
        <div className="mt-6 flex h-60 items-center justify-center rounded-xl border border-dashed border-slate-200 bg-slate-50 px-5 text-center text-sm leading-6 text-slate-500">
          A trend will appear after monitoring summaries have been generated.
        </div>
      ) : (
        <div
          className="mt-5 h-60 w-full sm:h-64"
          role="img"
          aria-label={`${title} chart`}
          aria-describedby={summaryId}
        >
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData} margin={{ top: 12, right: 4, bottom: 0, left: -18 }}>
              <defs>
                <linearGradient id={`${config.fillId}-${summaryId.replace(/:/g, "")}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor={config.stroke} stopOpacity={0.32} />
                  <stop offset="100%" stopColor={config.stroke} stopOpacity={0.02} />
                </linearGradient>
              </defs>
              <CartesianGrid stroke="#dbe4ef" strokeDasharray="4 5" vertical={false} />
              <XAxis
                dataKey="timestamp"
                tickFormatter={(value: string) => axisDateFormatter.format(new Date(value))}
                tick={{ fill: "#64748b", fontSize: 11 }}
                axisLine={false}
                tickLine={false}
                minTickGap={28}
              />
              <YAxis
                domain={metric === "uptime" ? [0, 100] : ["auto", "auto"]}
                tickFormatter={(value: number) => `${value}${config.unit}`}
                tick={{ fill: "#64748b", fontSize: 11 }}
                axisLine={false}
                tickLine={false}
                width={58}
              />
              <Tooltip
                labelFormatter={(value) => tooltipDateFormatter.format(new Date(String(value)))}
                formatter={(value) =>
                  value == null
                    ? ["No measurement", config.label]
                    : [formatValue(Number(value)), config.label]
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
                itemStyle={{ color: "#e2e8f0" }}
              />
              <Area
                type="monotone"
                dataKey={config.dataKey}
                stroke={config.stroke}
                fill={`url(#${config.fillId}-${summaryId.replace(/:/g, "")})`}
                strokeWidth={2.5}
                dot={false}
                activeDot={{ r: 4, fill: config.stroke, stroke: "white", strokeWidth: 2 }}
                connectNulls={false}
                isAnimationActive={!prefersReducedMotion}
                animationDuration={650}
                animationEasing="ease-out"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      )}
      <p id={summaryId} className="sr-only">
        {summary}
      </p>
    </article>
  );
}
