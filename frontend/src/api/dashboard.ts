import { api } from "./client";
import type { DashboardResponse } from "../types/api";

export function getDashboardSummary() {
  return api.get<DashboardResponse>("/api/dashboard/summary");
}
