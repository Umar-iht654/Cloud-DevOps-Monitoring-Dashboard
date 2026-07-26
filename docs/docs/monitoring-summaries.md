# Monitoring Summary Aggregation

The backend stores raw health checks for recent monitoring detail.

Raw health checks are useful for recent debugging, charts and health-check history, but they can grow quickly over time.

To preserve long-term monitoring statistics, the backend creates hourly and daily summary rows.

## Raw health checks

Raw health checks are individual records created by the background health checker.

They are used for recent service history and detailed debugging.

## Hourly summaries

Hourly summaries store one row per service per completed hour.

They include:

- total checks
- successful checks
- failed checks
- response-time sample count
- average response time
- minimum response time
- maximum response time
- uptime percentage

## Daily summaries

Daily summaries store one row per service per completed day.

They are created from hourly summaries.

## Why summaries are needed

Health-check retention can delete old raw health-check rows.

Summary rows allow the application to preserve long-term uptime and response-time statistics without keeping every raw check forever.

## Alerts

Alerts are separate from raw health checks and summaries.

Alerts represent important events, such as a service going down.

Alerts are kept even when old raw health checks are removed.

## Current limitation

The existing dashboard and service summary API still read from raw health checks.

A future update can use summary tables for longer-term reporting.