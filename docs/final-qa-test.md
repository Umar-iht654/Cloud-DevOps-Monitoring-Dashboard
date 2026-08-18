# Final QA Testing Report

This document records the final end-to-end quality assurance testing for the Cloud Service Monitoring Dashboard.

The purpose of this test run is to confirm that the main backend, frontend, database, monitoring, alerting, reporting, email, Docker and DevOps features work together as a complete system.

Testing also identified several usability and reporting improvements which were implemented and re-tested during the final QA process.

## Test Environment

| Item             | Details                           |
| ---------------- | --------------------------------- |
| Tester           | Umar Ihtesham                     |
| Test date        | 18 August 2026                    |
| Browser          | Google Chrome                     |
| Operating system | Windows                           |
| Backend          | Go, Gin, GORM                     |
| Frontend         | React, TypeScript, Vite, Tailwind |
| Database         | PostgreSQL                        |
| Containers       | Docker Compose                    |
| Monitoring       | Prometheus and Grafana            |
| Email provider   | Brevo SMTP                        |
| Environment      | Local full-stack Docker Compose   |

## Overall Result

| Area                | Result | Notes                                                                                                 |
| ------------------- | ------ | ----------------------------------------------------------------------------------------------------- |
| Authentication      | PASS   | Registration, verification, resend protection, login, logout and session handling tested successfully |
| Service management  | PASS   | Create, edit, delete, validation, ownership checks and readable service URLs tested successfully      |
| Health monitoring   | PASS   | Online, slow, down, recovery, history, statistics and response-time behaviour verified                |
| Alerts              | PASS   | Downtime alert creation, duplicate prevention, alert history and ownership isolation verified         |
| Email notifications | PASS   | Real verification and downtime emails delivered successfully through Brevo SMTP                       |
| Reports             | PASS   | 7-day and 30-day reports, live current-hour data, newly created services and report isolation tested  |
| Docker Compose      | PASS   | Full application stack starts and operates correctly                                                  |
| Prometheus          | PASS   | Backend metrics are exposed and scraped successfully                                                  |
| Grafana             | PASS   | Provisioned Grafana dashboard loads and displays monitoring data                                      |
| CI                  | PASS   | Frontend, backend and Docker GitHub Actions checks pass                                               |
| Security checks     | PASS   | Protected routes, user isolation, environment secrets and JWT behaviour verified                      |

**Overall QA Result: PASS**

No blocking defects remain in the tested MVP functionality.

## Final QA Checklist

### 1. Authentication and email verification

| Test                                   | Expected result                                                     | Actual result                                                                                             | Status |
| -------------------------------------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | ------ |
| Register with valid details            | Registration starts and user is sent to the email verification flow | Registration completed successfully and the account remained pending until verification                   | PASS   |
| Pending registration storage           | Permanent user should not be created before email verification      | Registration was stored as a pending registration and no permanent user was created before verification   | PASS   |
| Invalid registration validation        | Invalid/missing fields show clear validation errors                 | Invalid form submissions were rejected with appropriate frontend validation                               | PASS   |
| Duplicate email registration           | Existing or pending email is rejected safely                        | Duplicate registration attempts were handled safely                                                       | PASS   |
| Real email verification                | Verification email arrives in the user's inbox                      | Real verification email successfully delivered through Brevo SMTP                                         | PASS   |
| Verification email formatting          | Verification email should be readable and professionally formatted  | HTML verification email rendered correctly with branded layout and verification button                    | PASS   |
| Valid verification link                | Account is verified and user can continue to login                  | Verification succeeded, permanent user was created and pending registration was removed                   | PASS   |
| Invalid verification token             | Invalid token shows a useful error state                            | Invalid verification token was rejected correctly without creating an account                             | PASS   |
| Expired/invalid verification behaviour | Invalid or unusable links must not verify an account                | Verification flow failed safely for unusable tokens                                                       | PASS   |
| Resend verification flow               | User can request another verification email                         | Verification email resend worked successfully                                                             | PASS   |
| Verification resend countdown          | Frontend should show the two-minute resend cooldown clearly         | Added a visible countdown such as `Resend available in 1:42`; button remains disabled until cooldown ends | PASS   |
| Rate limiting email sends              | Repeated resend attempts are limited by cooldown and daily limit    | Two-minute cooldown and daily verification email limit operated correctly                                 | PASS   |
| Login with verified account            | User can log in and access the dashboard                            | Verified user logged in successfully                                                                      | PASS   |
| Login with wrong password              | Login fails with a clear error                                      | Incorrect password was rejected correctly                                                                 | PASS   |
| Logout                                 | Token is cleared and user is returned to login                      | Logout removed the stored authentication state and protected pages were no longer accessible              | PASS   |
| Session persistence after refresh      | Logged-in user remains authenticated after page refresh             | Session persisted correctly while the JWT remained valid                                                  | PASS   |
| Protected route redirect               | Logged-out user is redirected away from protected pages             | Dashboard, reports, alerts and service routes correctly required authentication                           | PASS   |

### 2. Service management

| Test                             | Expected result                                                               | Actual result                                                                          | Status |
| -------------------------------- | ----------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | ------ |
| Create service                   | Authenticated user can create a monitored service                             | Service was created successfully and appeared in the application                       | PASS   |
| Create service form validation   | Invalid URL, missing name or invalid interval is rejected                     | Frontend validation correctly rejected invalid values                                  | PASS   |
| Invalid service name             | Empty service name is rejected                                                | Empty name could not be submitted                                                      | PASS   |
| Invalid service URL              | Invalid or incomplete URLs are rejected                                       | Invalid URLs and unsupported schemes were rejected correctly                           | PASS   |
| Invalid expected HTTP status     | HTTP status must be between 100 and 599                                       | Invalid values such as 99 and 600 were rejected                                        | PASS   |
| Invalid slow threshold           | Threshold must be greater than zero                                           | Zero and negative thresholds were rejected                                             | PASS   |
| Minimum check interval           | Values below the 45-second minimum are rejected                               | 44 seconds was rejected and 45 seconds was accepted                                    | PASS   |
| Edit service                     | User can update service name, URL and settings                                | Service settings were updated and remained correct after refresh                       | PASS   |
| Delete service                   | User can delete a service and related data is removed                         | Disposable service was deleted successfully and did not reappear after refresh         | PASS   |
| Service detail page              | Service detail page loads correct service data                                | Name, URL, status, thresholds, interval, health history and alerts displayed correctly | PASS   |
| Unsaved changes protection       | User is warned before leaving a modified unsaved service form                 | Navigation warning appeared and correctly supported staying or discarding changes      | PASS   |
| Service ownership checks         | User cannot access another user's services                                    | Second test account could not access another user's service data                       | PASS   |
| Readable service URLs            | Service routes should be readable while keeping the ID as the real identifier | Routes changed from `/services/11` to `/services/11/backend-qa-service`                | PASS   |
| Legacy service URL compatibility | Existing ID-only URLs should still work                                       | ID-only service route redirects to the correct canonical ID/slug route                 | PASS   |
| Service rename updates URL       | Changing a service name should update the readable slug                       | Renaming a service resulted in the new canonical service-name URL                      | PASS   |
| Reports service links            | Service links from reports should use readable URLs                           | Report service links were updated to use the ID/service-name URL format                | PASS   |

### 3. Health monitoring

| Test                               | Expected result                                                        | Actual result                                                                                                           | Status |
| ---------------------------------- | ---------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | ------ |
| Online service check               | Working service is marked online                                       | `http://localhost:8080/health` returned the expected status and was classified as online                                | PASS   |
| Slow service check                 | Slow service is marked slow when above threshold                       | A deterministic delayed HTTP endpoint was used and correctly produced the slow state                                    | PASS   |
| Local response-time edge case      | Extremely fast local requests may measure 0 ms                         | Local Docker endpoint measured 0 ms, confirming why a dedicated delayed endpoint was required for deterministic testing | PASS   |
| Down service check                 | Broken service is marked down                                          | Unreachable endpoint such as `http://localhost:9999/health` was correctly marked down                                   | PASS   |
| Recovery check                     | Restored service changes from down back to online                      | Service successfully recovered after restoring the valid endpoint                                                       | PASS   |
| Health-check history               | Recent checks are stored and displayed                                 | Online, slow, down and recovery states were stored and visible in history                                               | PASS   |
| Health history database comparison | Frontend history should match PostgreSQL                               | Status, HTTP status, response time, error and timestamp data corresponded with stored rows                              | PASS   |
| Service health summary stats       | Uptime, failed checks and response-time data are shown                 | Summary values matched the underlying health-check records                                                              | PASS   |
| Uptime calculation                 | Online and slow checks count as successful while down counts as failed | Calculation matched `(online + slow) / total checks × 100`                                                              | PASS   |
| Average response time              | Displayed average should match recorded response samples               | Displayed average corresponded with stored response-time values                                                         | PASS   |
| Last downtime                      | Most recent down event is shown correctly                              | Last downtime matched the most recent stored `down` health check                                                        | PASS   |
| Response-time chart                | Chart renders response-time history correctly                          | Chart rendered correctly and showed a clear approximately 2-second spike from the delayed test endpoint                 | PASS   |
| Health-check retention cleanup     | Old raw health checks are cleaned up safely                            | Configured retention cleanup behaviour was verified                                                                     | PASS   |
| Hourly summary aggregation         | Completed hourly summary rows are generated                            | Hourly aggregation behaviour was verified                                                                               | PASS   |
| Daily summary aggregation          | Completed daily summary rows are generated                             | Daily aggregation behaviour was verified                                                                                | PASS   |

### 4. Alerts and email notifications

| Test                                 | Expected result                                                      | Actual result                                                                                                  | Status |
| ------------------------------------ | -------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | ------ |
| Downtime alert creation              | Service moving from non-down to down creates a `service_down` alert  | New online-to-down transition created the expected critical downtime alert                                     | PASS   |
| Downtime alert email                 | Service owner receives a real downtime email                         | Real downtime email successfully delivered through Brevo SMTP                                                  | PASS   |
| Downtime email formatting            | Downtime email should contain useful service and failure information | HTML email displayed service name, URL, reason and downtime timestamp correctly                                | PASS   |
| Duplicate alert prevention           | Repeated down checks do not create repeated alerts                   | Multiple consecutive down checks produced no duplicate alerts                                                  | PASS   |
| Duplicate email prevention           | Repeated down checks do not send repeated emails                     | Only the initial transition into down produced an email                                                        | PASS   |
| Recovery then new outage             | A later new online-to-down transition should create a new alert      | After recovery, a second outage generated a new alert and email correctly                                      | PASS   |
| Alerts page                          | Global alerts page displays alert history                            | Alert history loaded and showed recent outage events                                                           | PASS   |
| Alert ordering                       | Newest alerts should appear first                                    | Alerts were displayed newest first                                                                             | PASS   |
| Service-specific alerts              | Service detail page shows relevant alerts for that service           | Service-specific alert history was scoped correctly                                                            | PASS   |
| Alert history refresh                | Alert page refreshes data correctly                                  | Alert data updated correctly after refresh                                                                     | PASS   |
| Alert user isolation                 | One user cannot view another user's alerts                           | Second account could not access the first user's alerts                                                        | PASS   |
| SMTP missing fallback behaviour      | Backend logs/skips email safely when SMTP credentials are missing    | Development fallback behaviour safely avoids crashing the health worker                                        | PASS   |
| Email send failure handling          | Email failure is logged and does not crash the backend               | Email sending runs independently of the health-check transaction and failures do not roll back monitoring data | PASS   |
| Asynchronous downtime email delivery | SMTP network delays should not block the monitoring worker           | Downtime email is sent asynchronously and health checking continues normally                                   | PASS   |

### 5. Reports

| Test                              | Expected result                                                                            | Actual result                                                                                              | Status |
| --------------------------------- | ------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------- | ------ |
| Reports page                      | User can open the reports page from navigation                                             | Reports page loaded successfully                                                                           | PASS   |
| Reports 7-day range               | 7-day report loads summary cards, charts and table                                         | 7-day analytics loaded correctly                                                                           | PASS   |
| Reports 30-day range              | 30-day report loads summary cards, charts and table                                        | 30-day analytics loaded correctly                                                                          | PASS   |
| Maximum lookback behaviour        | 7-day/30-day range should not require a service to be that old                             | Reports now use all available data when the service is younger than the selected range                     | PASS   |
| Newly created service reporting   | New service should begin contributing to reports immediately                               | Newly created services now appear using available monitoring data without waiting for seven or thirty days | PASS   |
| Current incomplete hour           | Reports should include monitoring up to the current time                                   | Current-hour raw health checks are included instead of waiting for the hour to finish                      | PASS   |
| Historical summary use            | Older data should still use efficient hourly/daily summaries                               | Completed days use daily summaries and completed hours use hourly summaries                                | PASS   |
| Report double-count prevention    | Raw checks and summaries must not represent the same monitoring period twice               | Non-overlapping data boundaries prevent duplicate counting                                                 | PASS   |
| Report period metadata            | Displayed period should represent the actual available monitoring period                   | Young services show the available period from monitoring start up to the current time                      | PASS   |
| Current-day trend point           | Today's report point should include completed hours plus current raw checks                | Current day correctly combines hourly summaries and current-hour raw data                                  | PASS   |
| Reports summary metrics           | Uptime, successful checks, failures and latency calculations should be correct             | Report calculations displayed correctly                                                                    | PASS   |
| Service reliability ranking       | Best and worst performing services are identified correctly                                | Reliability ranking displayed successfully                                                                 | PASS   |
| Reports empty state               | Clear empty state appears when no monitoring data exists                                   | Empty state now only appears when genuinely no health-check data exists                                    | PASS   |
| Reports refresh button            | Manual refresh reloads report data                                                         | Refresh updated report information successfully                                                            | PASS   |
| Live report refresh               | Additional current-hour health checks should change totals without waiting for aggregation | Report totals updated using live data                                                                      | PASS   |
| Report data isolation             | User only sees report data for their own services                                          | Cross-user report access was prevented                                                                     | PASS   |
| Readable service links in reports | Service links should use canonical ID/service-name URLs                                    | Report table and highlights use readable canonical service URLs                                            | PASS   |

### 6. Frontend states and responsiveness

| Test                        | Expected result                                          | Actual result                                                                              | Status |
| --------------------------- | -------------------------------------------------------- | ------------------------------------------------------------------------------------------ | ------ |
| Frontend loading states     | Loading states appear while data is being fetched        | Loading/skeleton states displayed correctly across the main pages                          | PASS   |
| Frontend error states       | Useful error messages appear when requests fail          | Errors produced appropriate user-facing states                                             | PASS   |
| Frontend empty states       | Empty states appear when there is no data                | Empty service, alert and report states displayed correctly                                 | PASS   |
| Reports empty-state wording | Empty report state should match live reporting behaviour | Updated so it no longer incorrectly claims completed hourly/daily summaries are required   | PASS   |
| Form error accessibility    | Invalid fields should clearly communicate errors         | Form validation and invalid-state behaviour operated correctly                             | PASS   |
| Duplicate submit prevention | Forms should not submit multiple requests while saving   | Registration, resend and service forms prevented duplicate in-flight submissions           | PASS   |
| Mobile responsiveness       | Main pages remain usable on small screens                | Dashboard, authentication, services, alerts and reports remained usable on smaller layouts | PASS   |
| Desktop layout              | Main pages remain clear and usable on desktop            | Desktop layout worked correctly                                                            | PASS   |
| Navigation                  | Sidebar/top-level navigation works correctly             | Navigation between main application pages worked successfully                              | PASS   |
| Canonical URL navigation    | Service navigation uses readable canonical links         | Dashboard, reports, edit and service detail navigation use the improved routes             | PASS   |
| Browser refresh behaviour   | Main application pages remain functional after refresh   | Protected pages and service URLs continued to load correctly after refresh                 | PASS   |

### 7. Docker, Prometheus, Grafana and CI

| Test                            | Expected result                                                    | Actual result                                                             | Status |
| ------------------------------- | ------------------------------------------------------------------ | ------------------------------------------------------------------------- | ------ |
| Docker Compose clean rebuild    | Full stack starts successfully after clean rebuild                 | Stack rebuilt and started successfully using Docker Compose               | PASS   |
| Docker Compose full stack       | PostgreSQL, backend, frontend, Prometheus and Grafana run together | All core containers ran successfully together                             | PASS   |
| Docker container health         | Services should reach their configured healthy state               | Container health behaviour operated correctly                             | PASS   |
| Backend health endpoint         | `/health` returns successful response                              | Backend health endpoint responded successfully                            | PASS   |
| Frontend health endpoint        | Frontend health check returns successful response                  | Frontend health endpoint responded successfully                           | PASS   |
| PostgreSQL connectivity         | Backend can persist and query application data                     | Database connectivity remained stable throughout QA                       | PASS   |
| Docker networking               | Backend can monitor container-accessible endpoints correctly       | Internal monitoring routes and host Docker networking behaved as expected | PASS   |
| Prometheus metrics endpoint     | `/metrics` returns Prometheus-formatted metrics                    | Metrics endpoint returned successfully                                    | PASS   |
| Prometheus target status        | Backend target is UP in Prometheus                                 | Backend appeared as an active Prometheus target                           | PASS   |
| Prometheus health-check metrics | Health-check activity should produce monitoring metrics            | Monitoring actions were reflected in Prometheus metrics                   | PASS   |
| Grafana dashboard               | Provisioned Grafana dashboard loads correctly                      | Grafana dashboard loaded successfully                                     | PASS   |
| Grafana provisioning            | Dashboard configuration should be available automatically          | Provisioned monitoring dashboard was available without manual recreation  | PASS   |
| GitHub Actions frontend checks  | Frontend install, lint and build pass in CI                        | Frontend CI checks passed                                                 | PASS   |
| GitHub Actions backend checks   | Backend test and build pass in CI                                  | Backend CI checks passed                                                  | PASS   |
| GitHub Actions Docker builds    | Backend and frontend Docker images build in CI                     | Docker build checks passed                                                | PASS   |

### 8. Security and configuration

| Test                                         | Expected result                                                   | Actual result                                                             | Status |
| -------------------------------------------- | ----------------------------------------------------------------- | ------------------------------------------------------------------------- | ------ |
| Basic user data isolation                    | One user cannot view another user's services, alerts or reports   | Cross-user access was prevented across tested application areas           | PASS   |
| Service ownership isolation                  | Service detail access is scoped to the authenticated owner        | Second account could not access another user's service                    | PASS   |
| Alert ownership isolation                    | Alert APIs return only authenticated user's alerts                | Cross-user alerts were not exposed                                        | PASS   |
| Report ownership isolation                   | Reports contain only authenticated user's monitoring data         | Report data isolation worked correctly                                    | PASS   |
| `.env` not committed                         | Real secrets are not tracked by Git                               | Local secret configuration remained outside tracked repository files      | PASS   |
| JWT required for protected APIs              | Protected API routes reject missing/invalid tokens                | Protected backend resources required authenticated JWT access             | PASS   |
| Protected frontend routes                    | Logged-out users cannot remain on authenticated application pages | Protected frontend routes redirected correctly                            | PASS   |
| JWT session persistence                      | Valid stored JWT should keep the user logged in after refresh     | Authentication remained active after refresh                              | PASS   |
| Logout token clearing                        | Logout should remove local authentication state                   | JWT was removed on logout                                                 | PASS   |
| Email secrets loaded from environment        | SMTP credentials are read from `.env` or deployment env vars      | Brevo credentials were supplied through environment variables             | PASS   |
| Database credentials loaded from environment | Database configuration should not be hardcoded                    | PostgreSQL configuration is supplied through environment variables        | PASS   |
| JWT secret loaded from environment           | JWT signing secret should not be hardcoded                        | JWT secret is configured through environment variables                    | PASS   |
| SMTP password not exposed in source          | SMTP key should not appear in application code                    | SMTP secret remained outside committed source                             | PASS   |
| Verification token storage                   | Verification tokens should not be stored in plaintext             | Verification flow stores a verification token hash                        | PASS   |
| Verification token expiry                    | Verification links should have an expiry time                     | Expiry behaviour is included in the pending registration flow             | PASS   |
| Verification resend protection               | Resend endpoint should resist abuse                               | Cooldown, daily send limit and generic response behaviour are implemented | PASS   |
| Password storage                             | Passwords should not be stored in plaintext                       | Password hashes are used for stored accounts                              | PASS   |
| Service input validation                     | Backend and frontend validate monitored service input             | Service names, URLs, status codes, thresholds and intervals are validated | PASS   |
| URL protocol validation                      | Monitored URLs should be limited to appropriate HTTP protocols    | HTTP and HTTPS URLs are accepted while unsupported schemes are rejected   | PASS   |
| CORS configuration                           | Frontend origins should be configured rather than unrestricted    | Allowed frontend origins are controlled through configuration             | PASS   |

## Known Limitations

* CSV report export is not included in the current MVP and remains a future improvement.
* PDF export is not included in the final build and is planned as a future improvement.
* Scheduled email reports are not included in the final build and are planned as a future improvement.
* Domain-authenticated email sending is planned for cloud deployment to improve deliverability and provide a branded sender address.
* The current Brevo SMTP setup uses a verified sender and successfully supports local QA email delivery.
* Backend currently relies mostly on compile-level `go test ./...` validation because a larger dedicated backend unit/integration test suite has not yet been added.
* Automated frontend end-to-end browser tests are not yet implemented.
* Alert acknowledgement is not currently implemented.
* Recovery alert emails are not currently implemented.
* Alert pagination and advanced filtering are not currently implemented.
* Cloud deployment is not part of this local QA report and is tracked as the next major deployment phase.

## Final QA Notes

Use this section to record anything found during testing.

### Issues found

* Local Docker health checks against the backend itself were so fast that response time was often recorded as `0 ms`, making a 1 ms slow threshold unsuitable for deterministic slow-state testing.
* The original verification resend UX returned a generic success message but did not visually communicate the backend two-minute resend cooldown.
* Original service detail URLs only used the database ID, for example `/services/11`, which was functional but less readable.
* The original monitoring report behaviour relied primarily on completed hourly/daily summaries, meaning newly created services could initially have no report analytics during the current incomplete hour.
* Original report period presentation could imply a full 7-day or 30-day data range even when the monitored service had existed for a shorter period.
* Original report empty-state wording referred only to completed summary data and no longer matched the intended live-report behaviour after reporting improvements.

### Fixes made during QA

* Added a deterministic delayed local HTTP endpoint for slow-service QA. This confirmed:

  * Correct HTTP response within threshold → `online`
  * Correct HTTP response above threshold → `slow`
  * Failed/unreachable request → `down`
* Added a visible two-minute email verification resend countdown to the frontend.
* Disabled resend while the countdown is active and restart the timer after a successful resend.
* Improved service URLs from:

  * `/services/11`
  * to `/services/11/backend-qa-service`
* Kept the numeric database ID as the source of truth while adding a readable service-name slug.
* Added canonical route handling so legacy or incorrect slugs redirect to the correct service URL.
* Updated add-service, edit-service, dashboard and reports navigation to use the readable service URL.
* Updated service URLs automatically after renaming a service.
* Configured real Brevo SMTP email delivery.
* Verified real email verification messages arrive successfully.
* Verified real downtime alert emails arrive successfully.
* Improved verification and downtime emails with formatted HTML templates.
* Added safe escaping/sanitisation for dynamic email content and headers.
* Confirmed downtime email delivery is asynchronous so SMTP delays do not block health checks.
* Verified duplicate downtime alerts and emails are not created while a service remains continuously down.
* Verified a new alert and email are created after recovery followed by a new outage.
* Improved monitoring reports so 7-day and 30-day selections are maximum lookback windows rather than requiring the service to have existed for the full range.
* Newly created services now contribute report analytics immediately after health-check data exists.
* Added current incomplete-hour raw health checks to report analytics.
* Kept daily summaries for completed historical days and hourly summaries for completed hours.
* Added non-overlapping report source boundaries to prevent double counting.
* Updated report period metadata so young services display their actual available monitoring period.
* Updated current-day trend analytics to combine completed hourly data with current-hour raw health checks.
* Updated report empty-state wording so it appears only when no monitoring data exists.
* Re-tested report summary cards, reliability rankings, charts, tables and user isolation after the reporting changes.

### Final decision

**PASS — Final QA complete.**

The main MVP workflow has been tested successfully from account creation through monitoring and reporting:

```text
Register
↓
Receive verification email
↓
Verify account
↓
Login
↓
Create monitored service
↓
Automatic background health checks
↓
Online / Slow / Down status classification
↓
Health-check history and response-time analytics
↓
Downtime detection
↓
Alert creation
↓
Downtime email notification
↓
Recovery and future outage detection
↓
7-day / 30-day monitoring reports
↓
Prometheus metrics
↓
Grafana visualisation
```

The full local application stack operates successfully through Docker Compose, CI checks pass, user-owned data remains isolated, real SMTP email delivery works and no blocking MVP defects remain.

The project is ready to proceed to cloud deployment and final production configuration.
