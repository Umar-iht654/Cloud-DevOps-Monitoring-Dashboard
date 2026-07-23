import type { InputHTMLAttributes, ReactNode } from "react";

interface FormFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  hint?: string;
  icon?: ReactNode;
}

export function FormField({ label, hint, icon, id, ...props }: FormFieldProps) {
  return (
    <label htmlFor={id} className="form-field block">
      <span className="form-label mb-2 block text-sm font-semibold text-slate-700">{label}</span>
      <span className="relative block">
        {icon && (
          <span className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-slate-400">
            {icon}
          </span>
        )}
        <input
          id={id}
          {...props}
          className={`form-input w-full rounded-xl border border-slate-300 bg-white px-3.5 py-3 text-sm text-slate-900 outline-none transition placeholder:text-slate-500 disabled:cursor-not-allowed disabled:bg-slate-50 ${icon ? "pl-10" : ""}`}
        />
      </span>
      {hint && <span className="mt-1.5 block text-xs leading-5 text-slate-500">{hint}</span>}
    </label>
  );
}
