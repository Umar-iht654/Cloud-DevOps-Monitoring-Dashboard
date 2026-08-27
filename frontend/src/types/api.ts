export type ServiceStatus = "online" | "slow" | "down" | "unknown";

export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Service {
  id: number;
  user_id: number;
  name: string;
  url: string;
  expected_status_code: number;
  slow_threshold_ms: number;
  check_interval_seconds: number;
  current_status: ServiceStatus;
  created_at: string;
  updated_at: string;
}

export interface ServiceInput {
  name: string;
  url: string;
  expected_status_code: number;
  slow_threshold_ms: number;
  check_interval_seconds: number;
}

export interface HealthCheck {
  id: number;
  service_id: number;
  status: ServiceStatus;
  http_status_code: number | null;
  response_time_ms: number | null;
  error_message: string;
  checked_at: string;
  created_at?: string;
}

export interface Alert {
  id: number;
  user_id: number;
  service_id: number;
  health_check_id?: number | null;
  type: string;
  severity: string;
  title: string;
  message: string;
  created_at: string;
  service?: Service;
}

export interface DashboardSummary {
  total_services: number;
  online_services: number;
  slow_services: number;
  down_services: number;
  unknown_services: number;
  total_checks: number;
  successful_checks: number;
  failed_checks: number;
  average_uptime_percentage: number;
  average_response_time_ms: number;
  last_checked_at: string | null;
}

export interface ServiceSummary {
  service_id: number;
  service_name: string;
  current_status: ServiceStatus;
  total_checks: number;
  successful_checks: number;
  failed_checks: number;
  uptime_percentage: number;
  average_response_time_ms: number;
  last_checked_at: string | null;
  last_down_at: string | null;
}

export interface ApiMessage {
  message: string;
}

export interface LoginResponse extends ApiMessage {
  token: string;
  user: User;
}

export interface RegisterResponse extends ApiMessage {
  verificationSessionToken: string;
  verificationSessionExpiresAt: string;
  verificationSessionExpiresInSeconds: number;
}

export interface ResendVerificationResponse extends ApiMessage {
  code?: string;
  retryAfterSeconds?: number;
  expiresAt?: string;
  verificationSessionToken?: string;
  verificationSessionExpiresAt?: string;
  verificationSessionExpiresInSeconds?: number;
}

export type VerificationSessionStatus = "pending" | "verified" | "expired" | "consumed";

export interface VerificationSessionStatusResponse extends Partial<ApiMessage> {
  status: VerificationSessionStatus;
  code?: string;
  expiresAt?: string;
  token?: string;
  user?: User;
}

export interface CurrentUserResponse {
  user: User;
}

export interface ServicesResponse {
  services: Service[];
}

export interface ServiceResponse {
  message?: string;
  service: Service;
}

export interface DashboardResponse {
  summary: DashboardSummary;
}

export interface ServiceSummaryResponse {
  summary: ServiceSummary;
}

export interface HealthChecksResponse {
  service_id: number;
  returned_count: number;
  health_checks: HealthCheck[];
}

export interface AlertsResponse {
  alerts: Alert[];
  returned_count: number;
}

export interface ServiceAlertsResponse extends AlertsResponse {
  service_id: number;
}

export type ReportRange = "7d" | "30d" | "90d";

export interface OverviewServiceReport {
  service_id: number;
  service_name: string;
  total_checks: number;
  successful_checks: number;
  failed_checks: number;
  response_time_sample_count: number;
  average_response_time_ms: number;
  min_response_time_ms: number | null;
  max_response_time_ms: number | null;
  uptime_percentage: number;
}

export interface OverviewDailyReportPoint {
  period_start: string;
  period_end: string;
  total_checks: number;
  successful_checks: number;
  failed_checks: number;
  response_time_sample_count: number;
  average_response_time_ms: number;
  min_response_time_ms: number | null;
  max_response_time_ms: number | null;
  uptime_percentage: number;
}

export interface OverviewReport {
  range: ReportRange;
  period_start: string;
  period_end: string;
  total_checks: number;
  successful_checks: number;
  failed_checks: number;
  response_time_sample_count: number;
  average_response_time_ms: number;
  min_response_time_ms: number | null;
  max_response_time_ms: number | null;
  uptime_percentage: number;
  services: OverviewServiceReport[];
  daily: OverviewDailyReportPoint[];
}

export interface OverviewReportResponse extends OverviewReport {}
