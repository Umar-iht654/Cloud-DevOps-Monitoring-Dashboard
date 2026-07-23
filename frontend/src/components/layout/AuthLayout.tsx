import type { ReactNode } from "react";
import heroArtwork from "../../assets/hero.png";
import { LightRays } from "../effects/LightRays";
import { ActivityIcon, CheckIcon, PulseIcon } from "../ui/Icons";

const benefits = [
  "Track every endpoint from one dashboard",
  "Spot outages and degraded performance",
  "Turn health checks into useful history",
];

export function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="auth-shell grid min-h-screen lg:grid-cols-[1.08fr_0.92fr]">
      <section className="auth-visual relative hidden overflow-hidden border-r border-white/5 px-12 py-12 lg:flex lg:flex-col xl:px-20">
        <LightRays
          className="auth-light-rays"
          origin="top-center"
          speed={0.45}
          spread={0.58}
          length={2.4}
          mouseInfluence={0.07}
        />
        <div className="auth-orb absolute -left-32 top-1/4 h-96 w-96 rounded-full bg-cyan-400/10 blur-[100px]" />
        <div className="auth-orb auth-orb-delayed absolute -right-40 bottom-0 h-[28rem] w-[28rem] rounded-full bg-violet-600/12 blur-[120px]" />
        <div className="relative flex items-center gap-3">
          <div className="brand-pulse flex h-11 w-11 items-center justify-center rounded-xl bg-cyan-400 text-[#07111f] shadow-lg shadow-cyan-400/20">
            <PulseIcon className="h-6 w-6" />
          </div>
          <div>
            <p className="font-semibold tracking-tight text-white">Cloud Monitor</p>
            <p className="text-xs text-slate-500">Operations console</p>
          </div>
        </div>

        <div className="auth-copy relative z-10 my-auto max-w-2xl py-16">
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-cyan-300/15 bg-cyan-300/8 px-3 py-1.5 text-xs font-semibold text-cyan-300">
            <span className="live-dot h-1.5 w-1.5 rounded-full bg-emerald-400" />
            Reliability at a glance
          </div>
          <p className="max-w-xl text-5xl font-semibold leading-[1.05] tracking-[-0.045em] text-white xl:text-6xl">
            Know what is working before your users tell you.
          </p>
          <p className="mt-6 max-w-lg text-lg leading-8 text-slate-400">
            Monitor websites and APIs, understand response times, and catch failures from one focused workspace.
          </p>

          <div className="mt-9 space-y-3.5">
            {benefits.map((benefit) => (
              <div key={benefit} className="flex items-center gap-3 text-sm text-slate-300">
                <span className="flex h-6 w-6 items-center justify-center rounded-full bg-emerald-400/10 text-emerald-300 ring-1 ring-emerald-300/20">
                  <CheckIcon className="h-3.5 w-3.5" />
                </span>
                {benefit}
              </div>
            ))}
          </div>

          <div className="auth-showcase relative mt-10 h-44 max-w-xl xl:h-52" aria-hidden="true">
            <img
              src={heroArtwork}
              alt=""
              className="auth-art absolute -bottom-16 right-0 w-72 select-none object-contain xl:w-80"
            />
            <div className="telemetry-chip absolute bottom-4 left-0 rounded-2xl px-4 py-3 text-white">
              <div className="flex items-center gap-2 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                <ActivityIcon className="h-3.5 w-3.5 text-cyan-300" />
                Product preview
              </div>
              <p className="mt-1 text-sm font-semibold">One view. Every endpoint.</p>
            </div>
            <div className="telemetry-chip absolute right-4 top-0 rounded-2xl px-4 py-3 text-white xl:right-10">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-emerald-300">Live telemetry</p>
              <p className="mt-1 text-sm font-semibold">Availability &amp; latency</p>
            </div>
          </div>
        </div>

        <p className="relative text-xs text-slate-600">
          Built for clear, practical service monitoring.
        </p>
      </section>

      <main className="auth-form-canvas relative flex min-h-screen items-center justify-center px-5 py-10 sm:px-8">
        <div className="absolute left-5 top-6 flex items-center gap-2 lg:hidden">
          <div className="brand-pulse flex h-9 w-9 items-center justify-center rounded-xl bg-[#07111f] text-cyan-300">
            <PulseIcon className="h-5 w-5" />
          </div>
          <span className="font-semibold text-slate-900">Cloud Monitor</span>
        </div>
        <div className="auth-form-panel w-full max-w-md rounded-3xl p-6 sm:p-9">{children}</div>
      </main>
    </div>
  );
}
