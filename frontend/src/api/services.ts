import { api } from "./client";
import type {
  ApiMessage,
  HealthChecksResponse,
  ServiceInput,
  ServiceResponse,
  ServicesResponse,
  ServiceSummaryResponse,
} from "../types/api";

export function getServices() {
  return api.get<ServicesResponse>("/api/services");
}

export function getService(id: string | number) {
  return api.get<ServiceResponse>(`/api/services/${id}`);
}

export function createService(input: ServiceInput) {
  return api.post<ServiceResponse>("/api/services", input);
}

export function updateService(id: string | number, input: Partial<ServiceInput>) {
  return api.put<ServiceResponse>(`/api/services/${id}`, input);
}

export function deleteService(id: string | number) {
  return api.delete<ApiMessage>(`/api/services/${id}`);
}

export function getServiceSummary(id: string | number) {
  return api.get<ServiceSummaryResponse>(`/api/services/${id}/summary`);
}

// This matches the backend default so the frontend does not request too much history.
export function getHealthChecks(id: string | number, limit = 25) {
  return api.get<HealthChecksResponse>(`/api/services/${id}/health-checks`, {
    params: { limit },
  });
}
