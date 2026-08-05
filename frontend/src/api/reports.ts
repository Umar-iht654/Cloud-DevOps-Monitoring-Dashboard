import { api } from "./client";
import type { OverviewReportResponse, ReportRange } from "../types/api";

export function getOverviewReport(range: ReportRange) {
  return api.get<OverviewReportResponse>("/api/reports/overview", {
    params: { range },
  });
}
