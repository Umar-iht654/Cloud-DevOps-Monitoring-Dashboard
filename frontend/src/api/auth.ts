import { api } from "./client";
import type {
  ApiMessage,
  CurrentUserResponse,
  LoginResponse,
  RegisterResponse,
} from "../types/api";

export function login(email: string, password: string) {
  return api.post<LoginResponse>("/api/auth/login", { email, password });
}

export function register(name: string, email: string, password: string) {
  return api.post<RegisterResponse>("/api/auth/register", {
    name,
    email,
    password,
  });
}

export function verifyEmail(token: string) {
  return api.get<ApiMessage>("/api/auth/verify-email", {
    params: { token },
  });
}

export function resendVerificationEmail(email: string) {
  return api.post<ApiMessage>("/api/auth/resend-verification", { email });
}

export function getCurrentUser() {
  return api.get<CurrentUserResponse>("/api/auth/me");
}
