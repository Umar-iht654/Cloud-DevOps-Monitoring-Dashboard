// Converts a service name into a URL-friendly slug.
export function serviceSlug(name: string) {
  const slug = name
    .trim()
    .toLowerCase()
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");

  return slug || "service";
}

// Builds the canonical service details URL.
export function servicePath(id: number | string, name: string) {
  return `/services/${id}/${serviceSlug(name)}`;
}

// Builds the canonical edit URL.
export function serviceEditPath(id: number | string, name: string) {
  return `${servicePath(id, name)}/edit`;
}