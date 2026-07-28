function SkeletonBlock({ className }: { className: string }) {
  return <div className={`skeleton rounded-xl ${className}`} aria-hidden="true" />;
}

export function AlertsPageSkeleton() {
  return (
    <div role="status" aria-label="Loading alerts" className="space-y-5">
      <span className="sr-only">Loading alerts</span>
      <div className="dashboard-hero rounded-[1.75rem] p-6 sm:p-8">
        <SkeletonBlock className="h-4 w-32 bg-slate-700" />
        <SkeletonBlock className="mt-5 h-10 w-64 max-w-full bg-slate-700" />
        <SkeletonBlock className="mt-4 h-5 w-[34rem] max-w-full bg-slate-700" />
      </div>
      <div className="premium-panel rounded-3xl p-5 sm:p-6">
        <div className="flex items-center justify-between gap-4">
          <div className="space-y-2">
            <SkeletonBlock className="h-6 w-40" />
            <SkeletonBlock className="h-4 w-64 max-w-full" />
          </div>
          <SkeletonBlock className="h-10 w-24" />
        </div>
        <div className="mt-6 space-y-3">
          {[0, 1, 2, 3].map((item) => (
            <SkeletonBlock key={item} className="h-20 w-full" />
          ))}
        </div>
      </div>
    </div>
  );
}
