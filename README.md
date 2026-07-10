# Cloud DevOps Monitoring Dashboard

A Dockerised Cloud DevOps Monitoring Dashboard for tracking service uptime, response times, failures, alerts and reliability metrics using Go, PostgreSQL, Prometheus and Grafana.

## Contributors
This project is currently being developed collaboratively as a portfolio project by:

* Umar Ihtesham
* Ateeq Ur Rehman 

## Overview

Modern applications often rely on multiple services such as frontend websites, backend APIs, authentication systems, payment services and admin panels. If one of these services becomes slow, unreliable or unavailable, it can affect users and business operations.

This monitoring dashboard provides a central place to track the health of services. Users can add services, automatically check whether they are online, slow or down, record response times, view failure history and monitor reliability over time.

This project demonstrates practical backend development, DevOps fundamentals, cloud deployment, monitoring, containerisation, CI/CD, database design and production focused software engineering.

## Key Features

Service monitoring — add services and track their health automatically.

Health checks — record HTTP status codes, response times and check timestamps.

Status detection — classify services as online, slow or down.

Dashboard — view total services, online services, slow services, down services and average uptime.

Service detail page — view recent checks, response time history, uptime percentage and failure count.

Alerts — detect downtime and notify users when a service goes down.

Metrics — expose backend metrics for Prometheus.

Grafana dashboards — visualise reliability and system metrics.

Authentication — users can register, log in and access their own monitored services securely.

Dockerised setup — run the full application locally using Docker Compose.

## Tech Stack

| Area | Technology |
|---|---|
| Backend | Go |
| Frontend | React or Angular |
| Database | PostgreSQL |
| Monitoring | Prometheus |
| Visualisation | Grafana |
| Containerisation | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Cloud | AWS EC2, GCP Compute Engine or Azure VM |

## Planned Pages

Login and Register

Dashboard

Add/Edit Service

Service Detail Page

Alerts

Grafana Metrics Dashboard

Public Status Page

## Database Models

The application will use the following main models:

User — stores account and login details.

Service — stores monitored service details such as name, URL, expected status code and slow response threshold.

HealthCheck — stores each recorded check, including status, HTTP code, response time and timestamp.

Alert — stores downtime alerts and notification history.

Incident — stores service failure events and recovery details.

## Development Plan

Set up the Go backend and frontend application.

Design the PostgreSQL database schema.

Add Dockerfiles and Docker Compose for local development.

Build user registration, login and protected routes.

Implement service creation and management.

Create the background health checker.

Store health check history in the database.

Build the dashboard with service status cards.

Add service detail pages with response time history.

Add alerts for service downtime.

Expose Prometheus metrics from the backend.

Connect Prometheus and Grafana.

Add GitHub Actions for testing and build checks.

Deploy the project to a cloud virtual machine.

## Minimum Viable Product

This project is currently being developed.

The MVP will include authentication, service management, automatic health checks, online/slow/down status detection, health check history, uptime percentage calculation, a dashboard, service detail pages, Docker Compose setup and basic Prometheus metrics.

## Future Improvements

* Public status page
* Incident timeline
* Email alerts
* Slack or Discord alerts
* CSV report export
* Cloud storage for reports
* Advanced analytics
* Weekly and monthly uptime reports

## Goal

The goal of this project is to build a realistic, deployable Cloud DevOps Monitoring Dashboard that demonstrates real-world software engineering skills, including backend development, service monitoring, Docker, PostgreSQL, Prometheus, Grafana, CI/CD, cloud deployment and collaborative GitHub workflows.
