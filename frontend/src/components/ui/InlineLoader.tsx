export function InlineLoader({ label = "Loading" }: { label?: string }) {
  return (
    <div role="status" aria-live="polite" className="flex min-h-52 items-center justify-center rounded-3xl border border-slate-200 bg-white">
      <div className="text-center">
        <div className="mx-auto mb-3 h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-cyan-500" />
        <p className="text-sm text-slate-500">{label}…</p>
      </div>
    </div>
  );
}
