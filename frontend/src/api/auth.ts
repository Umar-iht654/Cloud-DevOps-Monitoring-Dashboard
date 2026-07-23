import { api } from "./client";
import type {
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

export function getCurrentUser() {
  return api.get<CurrentUserResponse>("/api/auth/me");
}
