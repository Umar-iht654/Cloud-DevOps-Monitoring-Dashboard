# StatusWatch

A production-deployed cloud service monitoring platform for tracking uptime, response times, failures, alerts and reliability metrics using Go, React, PostgreSQL, Prometheus and Grafana.

Live Demo: https://statuswatch.duckdns.org

## Contributors

Built collaboratively by:

- Umar Ihtesham
- Ateeq Ur Rehman

## Overview

StatusWatch is a full-stack monitoring platform for tracking the availability and performance of web services and APIs.

Users can add services to monitor, automatically run health checks, analyse uptime and latency, receive downtime alerts, inspect historical incidents and generate reliability reports from a central dashboard.

The project demonstrates practical software engineering across backend development, frontend development, observability, security, containerisation, CI/CD and cloud deployment.

## Features

- Service Monitoring - automatically monitor websites and APIs at configurable intervals.
- Health Checks - record HTTP status codes, response times, failures and timestamps.
- Status Detection - classify services as Online, Slow or Down.
- Dashboard - view uptime, latency, service health and system-wide monitoring statistics.
- Service Analytics - inspect response-time history, failure counts and uptime percentages.
- Downtime Alerts - automatically create alerts and send email notifications when services fail.
- Monitoring Reports - analyse 7-day and 30-day uptime, failures, latency and reliability trends.
- Authentication - JWT authentication, email verification and protected user-owned resources.
- Cross-device Verification - verify an account from a phone while automatically completing sign-in on the original browser.
- Observability - expose Prometheus metrics and visualise application behaviour through Grafana.
- Data Retention and Aggregation - retain recent raw checks while generating hourly and daily monitoring summaries.
- Containerisation - run the full application stack using Docker Compose.
- CI/CD - automatically test and build frontend, backend and Docker images with GitHub Actions.
- Production Deployment - deployed to an Oracle Cloud VM behind Caddy with HTTPS.

## Tech Stack

| Area | Technology |
|---|---|
| Backend | Go, Gin, GORM |
| Frontend | React, TypeScript, Vite, Tailwind CSS, Recharts |
| Database | PostgreSQL |
| Authentication | JWT |
| Monitoring | Prometheus |
| Visualisation | Grafana |
| Containerisation | Docker, Docker Compose |
| CI/CD | GitHub Actions |
| Reverse Proxy | Caddy |
| Cloud | Oracle Cloud Infrastructure |
| HTTPS | Automatic TLS via Caddy |

## Architecture

```text
                         Internet
                            |
                            v
                     Caddy / HTTPS
                            |
                +-----------+-----------+
                |                       |
                v                       v
          React Frontend            Go API
                                       |
                         +-------------+-------------+
                         |             |             |
                         v             v             v
                    PostgreSQL    Health Checker  Prometheus
                                                     |
                                                     v
                                                   Grafana
```

The background monitoring worker periodically checks registered services and stores health-check results in PostgreSQL. Aggregated data is used by the dashboard and reports, while application metrics are exposed to Prometheus and visualised through Grafana.

## Production Deployment

StatusWatch is deployed on an Oracle Cloud Infrastructure VM using Docker Compose.

Production traffic flows through Caddy, which provides:

- HTTPS/TLS
- HTTP to HTTPS redirects
- reverse proxy routing
- frontend and API routing

Internal services such as PostgreSQL and Prometheus are not exposed publicly. Grafana administration is accessed through an SSH tunnel.

## Engineering Highlights

### Secure Authentication

- JWT-based authentication
- password hashing
- protected routes
- user-level data isolation
- email verification
- verification token expiry
- rate-limited verification resends
- secure temporary verification sessions for cross-device auto-login

### Monitoring Pipeline

```text
Service
  |
  v
Background health checker
  |
  v
HTTP request + latency measurement
  |
  v
Online / Slow / Down classification
  |
  v
PostgreSQL
  |
  v
Dashboard / Alerts / Reports
```

### Reliability Data

Raw monitoring checks are retained for recent analysis while hourly and daily summaries preserve longer-term reliability information without indefinitely growing the raw health-check table.

### Observability

The backend exposes Prometheus metrics which are collected by Prometheus and displayed through provisioned Grafana dashboards.

## CI

GitHub Actions validates changes by running:

- frontend linting
- frontend production builds
- backend tests
- backend builds
- frontend Docker builds
- backend Docker builds

## Future Improvements

The core platform is complete. Potential future improvements include:

- bounded worker pool for concurrent health checks
- improved scheduling for large numbers of monitored services
- additional database indexes and query optimisation
- distributed monitoring workers
- queue-based monitoring jobs
- load and performance testing
- public service status pages
- Slack or Discord notifications
- CSV or PDF report exports
- scheduled email reports
- expanded automated frontend testing

## Purpose

StatusWatch was built as a portfolio project to demonstrate the design, implementation and deployment of a production-style monitoring system using modern full-stack and DevOps technologies.
