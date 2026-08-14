# Final QA Testing Report

This document records the final end-to-end quality assurance testing for the Cloud Service Monitoring Dashboard.

The purpose of this test run is to confirm that the main backend, frontend, database, monitoring, alerting, reporting, email, Docker and DevOps features work together as a complete system.

## Test Environment

| Item | Details |
|---|---|
| Tester | Umar Ihtesham |
| Test date | 14 August 2026 |
| Browser | Google Chrome |
| Operating system | Windows |
| Backend | Go, Gin, GORM |
| Frontend | React, TypeScript, Vite, Tailwind |
| Database | PostgreSQL |
| Containers | Docker Compose |
| Monitoring | Prometheus and Grafana |
| Email provider | Brevo SMTP |
| Environment | Local full-stack Docker Compose |

## Overall Result

| Area | Result | Notes |
|---|---|---|
| Authentication | Not tested yet | Registration, verification, login and logout must be tested |
| Service management | Not tested yet | Create, edit, delete and validation must be tested |
| Health monitoring | Not tested yet | Online, slow and down checks must be tested |
| Alerts | Not tested yet | Alert history and duplicate prevention must be tested |
| Email notifications | Not tested yet | Verification and downtime emails must be tested |
| Reports | Not tested yet | Reports page and API-backed summary data must be tested |
| Docker Compose | Not tested yet | Full stack should start with Docker Compose |
| Prometheus | Not tested yet | Backend metrics should be scraped successfully |
| Grafana | Not tested yet | Provisioned dashboard should load |
| CI | Not tested yet | GitHub Actions should pass |
| Security checks | Not tested yet | User isolation, env secrets and protected routes must be tested |

## Final QA Checklist

### 1. Authentication and email verification

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Register with valid details | Registration starts and user is sent to the email verification flow |  | Not tested |
| Invalid registration validation | Invalid/missing fields show clear validation errors |  | Not tested |
| Duplicate email registration | Existing or pending email is rejected safely |  | Not tested |
| Real email verification | Verification email arrives in the user's inbox |  | Not tested |
| Valid verification link | Account is verified and user can continue to login |  | Not tested |
| Invalid verification token | Invalid or expired token shows a useful error state |  | Not tested |
| Resend verification flow | User can request another verification email |  | Not tested |
| Rate limiting email sends | Repeated resend attempts are limited by cooldown and daily limit |  | Not tested |
| Login with verified account | User can log in and access the dashboard |  | Not tested |
| Login with wrong password | Login fails with a clear error |  | Not tested |
| Logout | Token is cleared and user is returned to login |  | Not tested |
| Session persistence after refresh | Logged-in user remains authenticated after page refresh |  | Not tested |
| Protected route redirect | Logged-out user is redirected away from protected pages |  | Not tested |

### 2. Service management

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Create service | Authenticated user can create a monitored service |  | Not tested |
| Create service form validation | Invalid URL, missing name or invalid interval is rejected |  | Not tested |
| Edit service | User can update service name, URL and settings |  | Not tested |
| Delete service | User can delete a service and related data is removed |  | Not tested |
| Service detail page | Service detail page loads correct service data |  | Not tested |
| Service ownership checks | User cannot access another user's services |  | Not tested |

### 3. Health monitoring

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Online service check | Working service is marked online |  | Not tested |
| Slow service check | Slow service is marked slow when above threshold |  | Not tested |
| Down service check | Broken service is marked down |  | Not tested |
| Health-check history | Recent checks are stored and displayed |  | Not tested |
| Service health summary stats | Uptime, failed checks and response-time data are shown |  | Not tested |
| Response-time chart | Chart renders response-time history correctly |  | Not tested |
| Health-check retention cleanup | Old raw health checks are cleaned up safely |  | Not tested |
| Hourly summary aggregation | Completed hourly summary rows are generated |  | Not tested |
| Daily summary aggregation | Completed daily summary rows are generated |  | Not tested |

### 4. Alerts and email notifications

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Downtime alert creation | Service moving from non-down to down creates a `service_down` alert |  | Not tested |
| Downtime alert email | Service owner receives a real downtime email |  | Not tested |
| Duplicate alert prevention | Repeated down checks do not create repeated alerts |  | Not tested |
| Duplicate email prevention | Repeated down checks do not send repeated emails |  | Not tested |
| Alerts page | Global alerts page displays alert history |  | Not tested |
| Service-specific alerts | Service detail page shows relevant alerts for that service |  | Not tested |
| Alert history refresh | Alert page refreshes data correctly |  | Not tested |
| SMTP missing fallback behaviour | Backend logs/skips email safely when SMTP credentials are missing |  | Not tested |
| Email send failure handling | Email failure is logged and does not crash the backend |  | Not tested |

### 5. Reports

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Reports page | User can open the reports page from navigation |  | Not tested |
| Reports 7-day range | 7-day report loads summary cards, charts and table |  | Not tested |
| Reports 30-day range | 30-day report loads summary cards, charts and table |  | Not tested |
| Reports empty state | Clear empty state appears when no summary data exists |  | Not tested |
| Reports refresh button | Manual refresh reloads report data |  | Not tested |
| Report data isolation | User only sees report data for their own services |  | Not tested |

### 6. Frontend states and responsiveness

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Frontend loading states | Loading states appear while data is being fetched |  | Not tested |
| Frontend error states | Useful error messages appear when requests fail |  | Not tested |
| Frontend empty states | Empty states appear when there is no data |  | Not tested |
| Mobile responsiveness | Main pages remain usable on small screens |  | Not tested |
| Navigation | Sidebar/top-level navigation works correctly |  | Not tested |

### 7. Docker, Prometheus, Grafana and CI

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Docker Compose clean rebuild | Full stack starts successfully after clean rebuild |  | Not tested |
| Docker Compose full stack | PostgreSQL, backend, frontend, Prometheus and Grafana run together |  | Not tested |
| Backend health endpoint | `/health` returns successful response |  | Not tested |
| Frontend health endpoint | Frontend health check returns successful response |  | Not tested |
| Prometheus metrics endpoint | `/metrics` returns Prometheus-formatted metrics |  | Not tested |
| Prometheus target status | Backend target is UP in Prometheus |  | Not tested |
| Grafana dashboard | Provisioned Grafana dashboard loads correctly |  | Not tested |
| GitHub Actions frontend checks | Frontend install, lint and build pass in CI |  | Not tested |
| GitHub Actions backend checks | Backend test and build pass in CI |  | Not tested |
| GitHub Actions Docker builds | Backend and frontend Docker images build in CI |  | Not tested |

### 8. Security and configuration

| Test | Expected result | Actual result | Status |
|---|---|---|---|
| Basic user data isolation | One user cannot view another user's services, alerts or reports |  | Not tested |
| `.env` not committed | Real secrets are not tracked by Git |  | Not tested |
| JWT required for protected APIs | Protected API routes reject missing/invalid tokens |  | Not tested |
| Email secrets loaded from environment | SMTP credentials are read from `.env` or deployment env vars |  | Not tested |

## Known Limitations

- PDF export is not included in the final build and is planned as a future improvement.
- Scheduled email reports are not included in the final build and are planned as a future improvement.
- Domain-authenticated email sending is planned for deployment to improve deliverability and provide a branded sender address.
- Backend currently relies mostly on compile-level `go test ./...` checks because dedicated backend unit tests have not yet been added.
- Automated frontend end-to-end tests are not yet implemented.
- Alert acknowledgement is not currently implemented.
- Recovery alert emails are not currently implemented.
- Alert pagination and filtering are not currently implemented.

## Final QA Notes

Use this section to record anything found during testing.

### Issues found

- None recorded yet.

### Fixes made during QA

- None recorded yet.

### Final decision

Not complete yet. Final QA must be completed before the project is marked as finished.