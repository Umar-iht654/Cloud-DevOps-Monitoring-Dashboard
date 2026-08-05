function SkeletonBlock({ className }: { className: string }) {
  return <div className={`skeleton-block rounded-xl ${className}`} aria-hidden="true" />;
}

export function ReportsPageSkeleton() {
  return (
    <div role="status" aria-live="polite" aria-label="Loading monitoring reports">
      <span className="sr-only">Loading monitoring reports…</span>

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4" aria-hidden="true">
        {[0, 1, 2, 3].map((item) => (
          <article
            key={item}
            className="rounded-2xl border border-slate-200/80 bg-white p-5"
          >
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <SkeletonBlock className="h-4 w-28" />
                <SkeletonBlock className="mt-3 h-9 w-20" />
              </div>
              <SkeletonBlock className="h-10 w-10" />
            </div>
            <SkeletonBlock className="mt-5 h-3 w-36" />
          </article>
        ))}
      </section>

      <section className="mt-6 grid gap-4 lg:grid-cols-2" aria-hidden="true">
        {[0, 1].map((item) => (
          <article
            key={item}
            className="rounded-2xl border border-slate-200/80 bg-white p-5"
          >
            <SkeletonBlock className="h-4 w-28" />
            <SkeletonBlock className="mt-3 h-6 w-48 max-w-full" />
            <SkeletonBlock className="mt-5 h-3 w-32" />
          </article>
        ))}
      </section>

      <section className="mt-8 grid gap-5 xl:grid-cols-3" aria-hidden="true">
        {[0, 1, 2].map((item) => (
          <article
            key={item}
            className="rounded-2xl border border-slate-200/80 bg-white p-5"
          >
            <SkeletonBlock className="h-5 w-36" />
            <SkeletonBlock className="mt-3 h-3 w-52 max-w-full" />
            <SkeletonBlock className="mt-6 h-56 w-full" />
          </article>
        ))}
      </section>

      <section className="mt-8 rounded-2xl border border-slate-200/80 bg-white p-5" aria-hidden="true">
        <SkeletonBlock className="h-5 w-48" />
        <SkeletonBlock className="mt-3 h-3 w-72 max-w-full" />
        <SkeletonBlock className="mt-6 h-72 w-full" />
      </section>
    </div>
  );
}
