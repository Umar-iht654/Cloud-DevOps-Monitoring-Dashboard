import { Link } from "react-router-dom";
import { ArrowLeftIcon, PulseIcon } from "../components/ui/Icons";

export function NotFoundPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[#07111f] px-6 text-center text-white">
      <div>
        <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300 ring-1 ring-cyan-300/20">
          <PulseIcon className="h-7 w-7" />
        </div>
        <p className="mt-6 text-sm font-semibold text-cyan-300">404</p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight">Page not found</h1>
        <p className="mt-3 text-sm text-slate-400">The page you requested does not exist.</p>
        <Link to="/dashboard" className="mt-7 inline-flex items-center gap-2 rounded-xl bg-white px-4 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-100">
          <ArrowLeftIcon className="h-4 w-4" />
          Return to dashboard
        </Link>
      </div>
    </main>
  );
}
