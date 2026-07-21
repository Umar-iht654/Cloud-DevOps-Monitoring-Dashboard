# Backend API Guide

This document explains how the frontend should connect to the backend.

The backend is written in Go and runs on:

```text
http://localhost:8080
```

The frontend should call this backend URL when making API requests.

---

## CORS Setup

The frontend and backend run on different ports during development.

Example:

```text
Backend:  http://localhost:8080
Frontend: http://localhost:5173
```

Because they are on different ports, the browser may block frontend requests unless the backend allows them.

This is handled using CORS.

In the backend `.env` file, there is a setting called:

```env
FRONTEND_URLS=http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000,http://127.0.0.1:3000,http://localhost:4200,http://127.0.0.1:4200
```

This tells the backend which frontend URLs are allowed to send requests to it.

Common frontend ports:

```text
React + Vite: http://localhost:5173
React CRA:    http://localhost:3000
Angular:      http://localhost:4200
```

If the frontend runs on a different port later, add that frontend URL to `FRONTEND_URLS` in `backend/.env`, then restart the backend.

The `.env` file is not pushed to GitHub because it can contain private local settings and secrets.

---

## Authentication

Some routes are public, but most routes require the user to be logged in.

When the user logs in, the backend returns a JWT token.

The frontend must store that token and send it in the request headers for protected routes.

Header format:

```text
Authorization: Bearer YOUR_TOKEN_HERE
```

Example:

```js
const response = await fetch("http://localhost:8080/api/services", {
  headers: {
    Authorization: `Bearer ${token}`,
  },
});
```

---

# Auth Routes

## Register

Creates a new user account.

```http
POST /api/auth/register
```

Request body:

```json
{
  "name": "Umar",
  "email": "umar@example.com",
  "password": "password123"
}
```

Success response:

```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "name": "Umar",
    "email": "umar@example.com"
  }
}
```

---

## Login

Logs in an existing user and returns a JWT token.

```http
POST /api/auth/login
```

Request body:

```json
{
  "email": "umar@example.com",
  "password": "password123"
}
```

Success response:

```json
{
  "message": "Login successful",
  "token": "JWT_TOKEN_HERE",
  "user": {
    "id": 1,
    "name": "Umar",
    "email": "umar@example.com"
  }
}
```

The frontend should save the `token`.

---

## Get Current User

Gets the logged-in user's details.

```http
GET /api/auth/me
```

Requires token:

```text
Authorization: Bearer JWT_TOKEN_HERE
```

Success response:

```json
{
  "user": {
    "id": 1,
    "name": "Umar",
    "email": "umar@example.com"
  }
}
```

---

# Service Routes

These routes are for managing monitored services.

All service routes require a JWT token.

---

## Create Service

Adds a service for the backend to monitor.

```http
POST /api/services
```

Request body:

```json
{
  "name": "Portfolio Website",
  "url": "https://example.com",
  "expected_status_code": 200,
  "slow_threshold_ms": 750,
  "check_interval_seconds": 60
}
```

What each field means:

```text
name = display name for the service
url = the URL that the backend will check
expected_status_code = the HTTP status code that means healthy, usually 200
slow_threshold_ms = response time limit before the service is classed as slow
check_interval_seconds = how often the backend should check the service
```

Success response:

```json
{
  "message": "Service created successfully",
  "service": {
    "id": 1,
    "user_id": 1,
    "name": "Portfolio Website",
    "url": "https://example.com",
    "expected_status_code": 200,
    "slow_threshold_ms": 750,
    "check_interval_seconds": 60,
    "current_status": "unknown"
  }
}
```

`current_status` starts as `unknown` because the background checker has not checked the service yet.

After the checker runs, it will become:

```text
online
slow
down
```

---

## Get All Services

Gets all services owned by the logged-in user.

```http
GET /api/services
```

Success response:

```json
{
  "services": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Portfolio Website",
      "url": "https://example.com",
      "expected_status_code": 200,
      "slow_threshold_ms": 750,
      "check_interval_seconds": 60,
      "current_status": "online"
    }
  ]
}
```

Frontend use:

```text
Dashboard service cards
Service list
Status badges
```

---

## Get One Service

Gets one service by ID.

```http
GET /api/services/:id
```

Example:

```http
GET /api/services/1
```

Success response:

```json
{
  "service": {
    "id": 1,
    "user_id": 1,
    "name": "Portfolio Website",
    "url": "https://example.com",
    "expected_status_code": 200,
    "slow_threshold_ms": 750,
    "check_interval_seconds": 60,
    "current_status": "online"
  }
}
```

---

## Update Service

Updates one service.

```http
PUT /api/services/:id
```

Example:

```http
PUT /api/services/1
```

Request body can include one or more fields:

```json
{
  "name": "Updated Portfolio Website",
  "slow_threshold_ms": 900
}
```

Success response:

```json
{
  "message": "Service updated successfully",
  "service": {
    "id": 1,
    "user_id": 1,
    "name": "Updated Portfolio Website",
    "url": "https://example.com",
    "expected_status_code": 200,
    "slow_threshold_ms": 900,
    "check_interval_seconds": 60,
    "current_status": "online"
  }
}
```

---

## Delete Service

Deletes one service.

```http
DELETE /api/services/:id
```

Example:

```http
DELETE /api/services/1
```

Success response:

```json
{
  "message": "Service deleted successfully"
}
```

---

# Health Check Routes

These routes return monitoring data that has been created by the background health checker.

All health check routes require a JWT token.

---

## Get Health Check History

Gets recent health checks for one service.

```http
GET /api/services/:id/health-checks
```

Example:

```http
GET /api/services/1/health-checks
```

You can limit the number of results:

```http
GET /api/services/1/health-checks?limit=5
```

Success response:

```json
{
  "service_id": 1,
  "returned_count": 2,
  "health_checks": [
    {
      "id": 10,
      "service_id": 1,
      "status": "online",
      "http_status_code": 200,
      "response_time_ms": 3,
      "error_message": "",
      "checked_at": "2026-07-21T22:04:34Z"
    },
    {
      "id": 9,
      "service_id": 1,
      "status": "online",
      "http_status_code": 200,
      "response_time_ms": 4,
      "error_message": "",
      "checked_at": "2026-07-21T22:04:24Z"
    }
  ]
}
```

Frontend use:

```text
Response time chart
Recent checks table
Status history
Failure history
```

---

## Get Service Summary

Gets calculated stats for one service.

```http
GET /api/services/:id/summary
```

Example:

```http
GET /api/services/1/summary
```

Success response:

```json
{
  "summary": {
    "service_id": 1,
    "service_name": "Portfolio Website",
    "current_status": "online",
    "total_checks": 20,
    "successful_checks": 20,
    "failed_checks": 0,
    "uptime_percentage": 100,
    "average_response_time_ms": 3,
    "last_checked_at": "2026-07-21T22:04:34Z",
    "last_down_at": null
  }
}
```

Frontend use:

```text
Service detail page
Uptime card
Average response time card
Failure count card
Last checked display
Last downtime display
```

---

# Dashboard Routes

## Get Dashboard Summary

Gets the main dashboard statistics for the logged-in user.

```http
GET /api/dashboard/summary
```

Requires token.

Success response:

```json
{
  "summary": {
    "total_services": 3,
    "online_services": 2,
    "slow_services": 1,
    "down_services": 0,
    "unknown_services": 0,
    "total_checks": 150,
    "successful_checks": 150,
    "failed_checks": 0,
    "average_uptime_percentage": 100,
    "average_response_time_ms": 4,
    "last_checked_at": "2026-07-21T22:04:34Z"
  }
}
```

Frontend use:

```text
Main dashboard cards
Total services card
Online services card
Slow services card
Down services card
Average uptime card
Average response time card
```

---

# Status Meanings

The backend uses these service statuses:

```text
unknown = service has not been checked yet
online = service returned the expected status code and was fast enough
slow = service returned the expected status code but response time was too high
down = service failed, timed out or returned the wrong status code
```

---

# Recommended Frontend Environment Variable

If using React + Vite, create a frontend `.env` file with:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Then frontend requests can use:

```js
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
```

Example login request:

```js
const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    email: "umar@example.com",
    password: "password123",
  }),
});
```

Example protected request:

```js
const response = await fetch(`${API_BASE_URL}/api/services`, {
  headers: {
    Authorization: `Bearer ${token}`,
  },
});
```

---

# Notes for Frontend Development

The frontend should store the JWT token after login.

For protected routes, always include:

```text
Authorization: Bearer TOKEN
```

A service may show as `unknown` immediately after being created. This is normal. The background checker updates it after it runs.

The backend already protects user data, so users only receive their own services and health checks.