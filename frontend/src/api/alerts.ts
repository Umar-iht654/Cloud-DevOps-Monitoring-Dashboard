import type { AlertsResponse, ServiceAlertsResponse } from "../types/api";
import { api } from "./client";

export function getAlerts(limit = 50) {
  return api.get<AlertsResponse>("/api/alerts", {
    params: { limit },
  });
}

export function getServiceAlerts(id: string | number, limit = 8) {
  return api.get<ServiceAlertsResponse>(`/api/services/${id}/alerts`, {
    params: { limit },
  });
}
