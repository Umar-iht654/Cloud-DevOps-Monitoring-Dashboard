function SkeletonBlock({ className }: { className: string }) {
  return <div className={`skeleton-block rounded-xl ${className}`} />;
}

export function DashboardSkeleton() {
  return (
    <div role="status" aria-live="polite" aria-label="Loading dashboard">
      <span className="sr-only">Loading dashboard…</span>

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-12" aria-hidden="true">
        {[3, 2, 2, 2, 3].map((columnSpan, index) => (
          <article
            key={`${columnSpan}-${index}`}
            className={`rounded-2xl border border-slate-200/80 bg-white p-5 ${
              columnSpan === 3 ? "xl:col-span-3" : "xl:col-span-2"
            } ${index === 4 ? "sm:col-span-2 xl:col-span-3" : ""}`}
          >
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <SkeletonBlock className="h-4 w-24" />
                <SkeletonBlock className="mt-3 h-9 w-16" />
              </div>
              <SkeletonBlock className="h-10 w-10" />
            </div>
            <SkeletonBlock className="mt-5 h-3 w-32" />
          </article>
        ))}
      </section>

      <section
        className="mt-5 grid overflow-hidden rounded-2xl border border-slate-200/80 bg-white md:grid-cols-2"
        aria-hidden="true"
      >
        {[0, 1].map((item) => (
          <div
            key={item}
            className="px-5 py-4 first:border-b first:border-slate-200/80 md:first:border-b-0 md:first:border-r"
          >
            <SkeletonBlock className="h-3 w-36" />
            <SkeletonBlock className="mt-2 h-6 w-28" />
          </div>
        ))}
      </section>

      <section className="mt-9" aria-hidden="true">
        <div className="mb-4 flex items-end justify-between gap-4">
          <div className="flex-1">
            <SkeletonBlock className="h-5 w-40" />
            <SkeletonBlock className="mt-2 h-3 w-56 max-w-full" />
          </div>
          <SkeletonBlock className="h-6 w-20 rounded-full" />
        </div>
        <div className="grid gap-4 md:grid-cols-2 2xl:grid-cols-3">
          {[0, 1, 2].map((item) => (
            <article
              key={item}
              className="rounded-2xl border border-slate-200/80 bg-white p-5"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex flex-1 items-center gap-3">
                  <SkeletonBlock className="h-11 w-11 shrink-0" />
                  <div className="flex-1">
                    <SkeletonBlock className="h-4 w-36 max-w-full" />
                    <SkeletonBlock className="mt-2 h-3 w-24" />
                  </div>
                </div>
                <SkeletonBlock className="h-7 w-16 rounded-full" />
              </div>
              <div className="mt-5 flex items-center gap-5 border-y border-slate-100 py-4">
                <SkeletonBlock className="h-20 w-20 shrink-0 rounded-full" />
                <div className="grid flex-1 grid-cols-2 gap-4">
                  <SkeletonBlock className="h-10 w-full" />
                  <SkeletonBlock className="h-10 w-full" />
                </div>
              </div>
              <div className="mt-4 flex justify-between gap-4">
                <SkeletonBlock className="h-3 w-28" />
                <SkeletonBlock className="h-3 w-20" />
              </div>
            </article>
          ))}
        </div>
      </section>
    </div>
  );
}
