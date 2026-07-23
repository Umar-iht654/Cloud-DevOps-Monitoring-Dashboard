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
  const data = [...healthChecks]
    .reverse()
    .filter((check) => check.response_time_ms != null)
    .map((check) => ({
      timestamp: check.checked_at,
      responseTime: check.response_time_ms,
      status: check.status,
    }));

  if (data.length === 0) {
    return (
      <div className="flex h-72 items-center justify-center rounded-2xl bg-slate-50 text-center">
        <div>
          <p className="font-medium text-slate-700">No response-time data yet</p>
          <p className="mt-1 text-sm text-slate-500">The chart will appear after the first completed check.</p>
        </div>
      </div>
    );
  }

  return (
    <div
      className="h-60 w-full sm:h-72"
      role="img"
      aria-label="Response time history in milliseconds, including the configured slow threshold"
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
            tickFormatter={(value: string) =>
              new Intl.DateTimeFormat("en-GB", {
                hour: "2-digit",
                minute: "2-digit",
              }).format(new Date(value))
            }
            tick={{ fill: "#64748b", fontSize: 11 }}
            axisLine={false}
            tickLine={false}
            minTickGap={28}
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
            formatter={(value) => [`${Number(value).toLocaleString()} ms`, "Response time"]}
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
            isAnimationActive={!prefersReducedMotion}
            animationDuration={750}
            animationEasing="ease-out"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
