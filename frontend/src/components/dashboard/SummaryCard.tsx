import type { ComponentType, CSSProperties, SVGProps } from "react";
import { positionSpotlight } from "../../utils/motion";

interface SummaryCardProps {
  label: string;
  value: string | number;
  helper: string;
  icon: ComponentType<SVGProps<SVGSVGElement>>;
  tone?: "cyan" | "emerald" | "amber" | "rose" | "slate";
  className?: string;
  featured?: boolean;
  delay?: number;
}

const tones = {
  cyan: "bg-cyan-50 text-cyan-700 ring-cyan-100",
  emerald: "bg-emerald-50 text-emerald-700 ring-emerald-100",
  amber: "bg-amber-50 text-amber-700 ring-amber-100",
  rose: "bg-rose-50 text-rose-700 ring-rose-100",
  slate: "bg-slate-100 text-slate-700 ring-slate-200",
};

const accents = {
  cyan: "#06b6d4",
  emerald: "#10b981",
  amber: "#f59e0b",
  rose: "#f43f5e",
  slate: "#7c3aed",
};

export function SummaryCard({
  label,
  value,
  helper,
  icon: Icon,
  tone = "cyan",
  className = "",
  featured = false,
  delay = 0,
}: SummaryCardProps) {
  const style = {
    "--metric-accent": accents[tone],
    animationDelay: `${delay}ms`,
  } as CSSProperties;

  return (
    <article
      className={`premium-surface metric-card rounded-2xl border border-slate-200/80 bg-white p-5 shadow-[0_1px_2px_rgba(15,23,42,0.03)] ${featured ? "metric-card-featured" : ""} ${className}`}
      style={style}
      onPointerMove={positionSpotlight}
    >
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-slate-500">{label}</p>
          <p className={`metric-tabular mt-2 font-semibold tracking-[-0.04em] text-slate-950 ${featured ? "text-4xl" : "text-3xl"}`}>{value}</p>
        </div>
        <span className={`metric-icon flex h-10 w-10 items-center justify-center rounded-xl ring-1 ring-inset transition-transform duration-300 ${tones[tone]}`}>
          <Icon className="h-5 w-5" />
        </span>
      </div>
      <p className="mt-4 text-xs leading-5 text-slate-500">{helper}</p>
    </article>
  );
}
