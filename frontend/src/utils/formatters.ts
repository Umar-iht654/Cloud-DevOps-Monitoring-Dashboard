export function formatDateTime(value: string | null | undefined) {
  if (!value) return "Not checked yet";

  return new Intl.DateTimeFormat("en-GB", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

export function formatRelativeTime(value: string | null | undefined) {
  if (!value) return "Not checked yet";

  const timestamp = new Date(value).getTime();
  const seconds = Math.round((timestamp - Date.now()) / 1000);
  const formatter = new Intl.RelativeTimeFormat("en", { numeric: "auto" });

  if (Math.abs(seconds) < 60) return formatter.format(seconds, "second");

  const minutes = Math.round(seconds / 60);
  if (Math.abs(minutes) < 60) return formatter.format(minutes, "minute");

  const hours = Math.round(minutes / 60);
  if (Math.abs(hours) < 24) return formatter.format(hours, "hour");

  return formatter.format(Math.round(hours / 24), "day");
}

export function formatPercentage(value: number) {
  return `${value.toFixed(value % 1 === 0 ? 0 : 2)}%`;
}

export function formatMilliseconds(value: number | null | undefined) {
  return value == null ? "—" : `${value.toLocaleString()} ms`;
}

export function hostnameFromUrl(value: string) {
  try {
    return new URL(value).hostname;
  } catch {
    return value;
  }
}
