import { Link, useNavigate } from "react-router-dom";
import { ArrowLeftIcon, PulseIcon } from "../components/ui/Icons";

export function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <main className="relative isolate flex min-h-screen items-center justify-center overflow-hidden bg-[#07111f] px-6 py-12 text-center text-white">
      <div
        aria-hidden="true"
        className="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_50%_20%,rgba(34,211,238,0.14),transparent_24rem),radial-gradient(circle_at_80%_90%,rgba(124,58,237,0.16),transparent_26rem)]"
      />
      <div
        aria-hidden="true"
        className="absolute inset-0 -z-10 opacity-20 [background-image:linear-gradient(rgba(255,255,255,0.06)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.06)_1px,transparent_1px)] [background-size:38px_38px]"
      />

      <section className="w-full max-w-xl rounded-3xl border border-white/10 bg-white/[0.045] px-6 py-10 shadow-2xl shadow-black/20 backdrop-blur-xl sm:px-10 sm:py-12">
        <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300 ring-1 ring-cyan-300/20">
          <PulseIcon className="h-7 w-7" />
        </div>
        <p className="mt-6 text-sm font-semibold uppercase tracking-[0.2em] text-cyan-300">
          Error 404
        </p>
        <h1 className="mt-3 text-3xl font-semibold tracking-tight sm:text-4xl">
          This page is off the radar
        </h1>
        <p className="mx-auto mt-4 max-w-md text-sm leading-6 text-slate-400 sm:text-base">
          The address may be incorrect, or the page may have moved. Return to your monitoring
          workspace or retrace your last step.
        </p>

        <div className="mt-8 flex flex-col-reverse justify-center gap-3 sm:flex-row">
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="inline-flex min-h-12 items-center justify-center rounded-xl border border-white/10 bg-white/5 px-5 py-3 text-sm font-semibold text-slate-200 transition hover:border-white/20 hover:bg-white/10 hover:text-white focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-cyan-300/30"
          >
            Go back
          </button>
          <Link
            to="/dashboard"
            className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-white px-5 py-3 text-sm font-semibold text-slate-900 shadow-lg shadow-black/10 transition hover:bg-cyan-50 focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-cyan-300/40"
          >
            <ArrowLeftIcon className="h-4 w-4" />
            Return to dashboard
          </Link>
        </div>
      </section>
    </main>
  );
}
