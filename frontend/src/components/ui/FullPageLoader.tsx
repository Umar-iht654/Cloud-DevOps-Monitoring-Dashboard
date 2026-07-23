import { PulseIcon } from "./Icons";

export function FullPageLoader({ label = "Loading" }: { label?: string }) {
  return (
    <div role="status" aria-live="polite" className="flex min-h-screen items-center justify-center bg-[#07111f] px-6 text-white">
      <div className="text-center">
        <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-400/10 ring-1 ring-cyan-300/20">
          <PulseIcon className="h-7 w-7 animate-pulse text-cyan-300" />
        </div>
        <p className="text-sm font-medium text-slate-300">{label}…</p>
      </div>
    </div>
  );
}
