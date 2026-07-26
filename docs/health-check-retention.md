# Health Check Retention

Raw health checks are stored frequently as services are monitored.

To stop the database from growing forever, the backend includes a configurable retention cleaner.

## Environment variables

HEALTH_CHECK_RETENTION_DAYS=14  
HEALTH_CHECK_CLEANUP_INTERVAL_HOURS=24

## Behaviour

HEALTH_CHECK_RETENTION_DAYS controls how many days of raw health checks are kept.

For example:

HEALTH_CHECK_RETENTION_DAYS=14

This keeps the latest 14 days of raw health checks and deletes older health checks.

Set this to 0 to disable cleanup:

HEALTH_CHECK_RETENTION_DAYS=0

HEALTH_CHECK_CLEANUP_INTERVAL_HOURS controls how often the cleanup worker runs.

For example:

HEALTH_CHECK_CLEANUP_INTERVAL_HOURS=24

This runs cleanup once every 24 hours.

## Alerts

Alerts are kept when old health checks are deleted.

If an alert points to an old health check that is removed, the alert remains and its health_check_id is cleared.

## Current limitation

Until hourly and daily aggregation is added, service summary statistics are calculated from retained raw health checks only.

Hourly and daily aggregation will be added in a future feature so long-term uptime statistics can be preserved without keeping every raw check forever.