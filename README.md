# Cloud DevOps Monitoring Dashboard

A Dockerised cloud monitoring platform for tracking service uptime, response times, failures and reliability metrics using Go, PostgreSQL, Prometheus and Grafana.

## Contributors

This project is currently being developed collaboratively as a portfolio project by:

- Umar Ihtesham
- Ateeq Ur Rehman

## Overview

Modern applications often rely on multiple services such as frontend websites, backend APIs, authentication systems, payment services and admin panels. If one of these services becomes slow, unreliable or unavailable, it can affect users and business operations.

This monitoring dashboard provides a central place to track the health of services. Users can add services, automatically check whether they are online, slow or down, record response times, view failure history and monitor reliability over time.

This project demonstrates practical backend development, DevOps fundamentals, cloud deployment, monitoring, containerisation, CI/CD, database design and production-focused software engineering.

## Key Features

- **Service Monitoring**
  - Add services and track their health automatically.

- **Health Checks**
  - Record HTTP status codes, response times and check timestamps.

- **Status Detection**
  - Classify services as online, slow or down.

- **Dashboard**
  - View total services, online services, slow services, down services and average uptime.

- **Service Detail Page**
  - View recent checks, response-time history, uptime percentage and failure count.

- **Authentication**
  - Users can register, log in and securely access their own monitored services.

- **Observability**
  - Expose backend metrics for Prometheus.
  - Visualise reliability and system metrics with Grafana.

- **Dockerised Setup**
  - PostgreSQL currently runs locally using Docker Compose, with frontend and backend containerisation planned.

## Tech Stack

| Area | Technology |
|---|---|
| Backend | Go, Gin, GORM |
| Frontend | React, Vite, TypeScript, Tailwind CSS, Recharts |
| Database | PostgreSQL |
| Monitoring | Prometheus |
| Visualisation | Grafana |
| Containerisation | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Cloud | AWS EC2, GCP Compute Engine or Azure VM |

## Planned Pages

- Login and Register
- Dashboard
- Add/Edit Service
- Service Detail Page
- Grafana Metrics Dashboard

## Database Models

The application will use the following main models:

- **User** — stores account and login details.
- **Service** — stores monitored service details such as name, URL, expected status code and slow response threshold.
- **HealthCheck** — stores each recorded check, including status, HTTP status code, response time and timestamp.
- **Alert** — stores downtime alerts and notification history.
- **Incident** — stores service failure events and recovery details.

## Current Progress

The authentication system, service management API, background health checker, dashboard summary, health-check history and React frontend have been implemented. PostgreSQL currently runs through Docker Compose. Prometheus, Grafana, full application containerisation, CI/CD and cloud deployment remain in development.

## Development Plan

1. Set up the Go backend and frontend application.
2. Design the PostgreSQL database schema.
3. Add Dockerfiles and Docker Compose for local development.
4. Build user registration, login and protected routes.
5. Implement service creation and management.
6. Create the background health checker.
7. Store health check history in the database.
8. Build the dashboard with service status cards.
9. Add service detail pages with response-time history.
10. Expose Prometheus metrics from the backend.
11. Connect Prometheus and Grafana.
12. Add GitHub Actions for testing and build checks.
13. Deploy the project to a cloud virtual machine.

## Minimum Viable Product

This project is currently being developed.

The MVP includes:

- Authentication
- Service management
- Automatic health checks
- Online / Slow / Down status detection
- Health check history
- Uptime percentage calculation
- Dashboard
- Service detail pages
- Docker Compose setup
- Basic Prometheus metrics
- Grafana dashboards

## Future Improvements

- Public status page
- Incident timeline
- Email alerts
- Slack or Discord alerts
- CSV report export
- Cloud storage for reports
- Advanced analytics
- Weekly and monthly uptime reports

## Goal

The goal of this project is to build a realistic, deployable Cloud DevOps Monitoring Dashboard that demonstrates real-world software engineering skills, including backend development, service monitoring, Docker, PostgreSQL, Prometheus, Grafana, CI/CD, cloud deployment and collaborative GitHub workflows.
