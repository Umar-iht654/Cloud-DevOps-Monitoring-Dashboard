import type { ServiceStatus } from "../../types/api";

const statusStyles: Record<ServiceStatus, string> = {
  online: "border-emerald-200 bg-emerald-50 text-emerald-700",
  slow: "border-amber-200 bg-amber-50 text-amber-700",
  down: "border-rose-200 bg-rose-50 text-rose-700",
  unknown: "border-slate-200 bg-slate-100 text-slate-600",
};

const dotStyles: Record<ServiceStatus, string> = {
  online: "bg-emerald-500",
  slow: "bg-amber-500",
  down: "bg-rose-500",
  unknown: "bg-slate-400",
};

export function StatusBadge({ status, live = false }: { status: ServiceStatus; live?: boolean }) {
  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full border px-2.5 py-1 text-xs font-semibold capitalize ${statusStyles[status]}`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${live && status === "online" ? "live-dot" : ""} ${dotStyles[status]}`} />
      {status}
    </span>
  );
}
