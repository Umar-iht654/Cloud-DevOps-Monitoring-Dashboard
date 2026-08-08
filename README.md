# Cloud Service Monitoring Dashboard

A Dockerised cloud monitoring platform for tracking service uptime, response times, failures, alerts and reliability metrics using Go, PostgreSQL, Prometheus and Grafana.

## Contributors

This project is currently being developed collaboratively as a portfolio project by:

- Umar Ihtesham
- Ateeq Ur Rehman

## Overview

Modern applications often rely on multiple services such as frontend websites, backend APIs, authentication systems, payment services and admin panels. If one of these services becomes slow, unreliable or unavailable, it can affect users and business operations.

This monitoring dashboard provides a central place to track the health of services. Users can add services, automatically check whether they are online, slow or down, record response times, view failure and alert history, monitor reliability over time and analyse longer-term monitoring reports.

This project demonstrates practical backend development, frontend development, DevOps fundamentals, monitoring, containerisation, CI/CD, database design and production-focused software engineering.

## Key Features

- **Service Monitoring**
  - Add, edit and delete monitored services.
  - Automatically track service health.

- **Health Checks**
  - Record HTTP status codes, response times and timestamps.
  - Configurable service check intervals.
  - Store historical monitoring data with configurable retention.

- **Status Detection**
  - Classify services as online, slow or down.

- **Dashboard**
  - View total services, online services, slow services, down services and average uptime.
  - View live service health information.

- **Service Detail Page**
  - View recent health checks.
  - View response-time history.
  - View uptime percentage and failure count.
  - View service-specific alert history.

- **Authentication**
  - User registration and login.
  - JWT-based protected routes.
  - User-owned monitored services.
  - Email verification flow.

- **Alerts**
  - Automatically create downtime alerts.
  - View global and service-specific alert history.
  - Track repeated outages and recoveries.

- **Monitoring Reports**
  - 7-day and 30-day monitoring reports.
  - Overall uptime, total checks, failed checks and response-time summaries.
  - Best and worst service reliability.
  - Daily uptime, response-time and failure charts.
  - Service reliability comparison.

- **Observability**
  - Prometheus metrics.
  - Grafana dashboards.

- **Dockerised Setup**
  - Run PostgreSQL, backend, frontend, Prometheus and Grafana together using Docker Compose.
  - Container health checks and startup dependencies.

- **Continuous Integration**
  - Frontend lint and production build checks.
  - Backend test and build checks.
  - Backend and frontend Docker image build checks using GitHub Actions.

## Tech Stack

| Area | Technology |
|---|---|
| Backend | Go, Gin, GORM |
| Frontend | React, Vite, TypeScript, Tailwind CSS, Recharts |
| Database | PostgreSQL |
| Authentication | JWT |
| Monitoring | Prometheus |
| Visualisation | Grafana |
| Containerisation | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Cloud | Planned |

## Application Pages

- Login
- Register
- Email Verification
- Dashboard
- Add/Edit Service
- Service Detail Page
- Alerts
- Reports
- Grafana Metrics Dashboard

## Database Models

The application uses the following main models:

- **User** — stores account, authentication and verification details.
- **Service** — stores monitored service details such as name, URL, expected status code, check interval and slow response threshold.
- **HealthCheck** — stores recorded monitoring checks including status, HTTP status code, response time and timestamp.
- **Alert** — stores downtime alert history and service failure events.

## Current Progress

The core monitoring platform is now functional.

Implemented:

- Go backend and React frontend.
- PostgreSQL database.
- JWT authentication.
- User registration and login.
- Email verification flow.
- Protected service CRUD API.
- Automatic background health checks.
- Health-check history.
- Dashboard summary metrics.
- Service detail pages.
- Alert history.
- 7-day and 30-day monitoring reports.
- Prometheus metrics.
- Grafana dashboards.
- Full Docker Compose stack.
- Frontend and backend Docker images.
- GitHub Actions CI for frontend, backend and Docker builds.

Remaining work includes production email delivery, additional automated testing, final polish and cloud deployment.

## Development Plan

1. Set up the Go backend and React frontend. ✅
2. Design the PostgreSQL database schema. ✅
3. Implement authentication and protected routes. ✅
4. Implement service creation and management. ✅
5. Create the background health checker. ✅
6. Store health-check history. ✅
7. Build the dashboard and service detail pages. ✅
8. Add alerts and alert history. ✅
9. Add monitoring reports and analytics. ✅
10. Expose Prometheus metrics. ✅
11. Connect Prometheus and Grafana. ✅
12. Dockerise the frontend and backend. ✅
13. Run the full stack through Docker Compose. ✅
14. Add GitHub Actions CI. ✅
15. Configure production email delivery.
16. Expand automated testing.
17. Deploy the project to a cloud virtual machine.

## Minimum Viable Product

The MVP includes:

- Authentication
- Email verification
- Service management
- Automatic health checks
- Online / Slow / Down status detection
- Health-check history
- Uptime percentage calculation
- Dashboard
- Alerts
- Monitoring reports
- Service detail pages
- Docker Compose setup
- Prometheus metrics
- Grafana dashboards
- GitHub Actions CI

## Future Improvements

- Public status page
- Incident timeline
- Slack or Discord notifications
- CSV report export
- Cloud storage for reports
- Advanced analytics
- Additional report ranges
- Automated frontend test suite
- Extended backend test coverage
- Cloud deployment

## Goal

The goal of this project is to build a realistic, deployable Cloud Service Monitoring Dashboard that demonstrates real-world software engineering skills, including backend development, frontend development, service monitoring, Docker, PostgreSQL, Prometheus, Grafana, CI/CD and collaborative GitHub workflows.
