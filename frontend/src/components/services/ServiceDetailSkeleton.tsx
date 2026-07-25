function SkeletonBlock({ className }: { className: string }) {
  return <div className={`skeleton-block rounded-xl ${className}`} />;
}

export function ServiceDetailSkeleton() {
  return (
    <div role="status" aria-live="polite" aria-label="Loading service history">
      <span className="sr-only">Loading service history…</span>

      <SkeletonBlock className="h-5 w-36" />

      <section
        className="mt-5 rounded-[1.75rem] bg-slate-900 p-6 lg:p-8"
        aria-hidden="true"
      >
        <SkeletonBlock className="h-7 w-28 bg-slate-700" />
        <SkeletonBlock className="mt-4 h-10 w-80 max-w-full bg-slate-700" />
        <SkeletonBlock className="mt-3 h-4 w-64 max-w-full bg-slate-700" />
      </section>

      <section className="mt-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-5" aria-hidden="true">
        {[0, 1, 2, 3, 4].map((item) => (
          <article
            key={item}
            className={`rounded-2xl border border-slate-200/80 bg-white p-5 ${
              item === 4 ? "sm:col-span-2 xl:col-span-1" : ""
            }`}
          >
            <SkeletonBlock className="h-3 w-24" />
            <SkeletonBlock className="mt-3 h-6 w-28 max-w-full" />
          </article>
        ))}
      </section>

      {[0, 1].map((item) => (
        <section
          key={item}
          className="mt-5 rounded-3xl border border-slate-200/80 bg-white p-5 sm:p-6"
          aria-hidden="true"
        >
          <SkeletonBlock className="h-5 w-40" />
          <SkeletonBlock className="mt-2 h-3 w-72 max-w-full" />
          <SkeletonBlock className={`mt-6 w-full ${item === 0 ? "h-64" : "h-52"}`} />
        </section>
      ))}
    </div>
  );
}
