# Frontend Testing Report

## Overview

This document records the frontend quality assurance completed for the Cloud Service Monitoring Dashboard.

- Test dates: 25 and 28 July, and 5 and 6 August 2026
- Tester: Ateeq Ur Rehman
- Browser: Google Chrome and Chromium on macOS

The tests covered authentication and email verification, routing, dashboard behaviour, service management, monitoring history, alert history, reports, selected error and recovery states, responsive design, keyboard accessibility, the production container and production build checks.

## Test Environment

- Local React frontend connected to the local Go API
- Containerised production frontend running through Docker Compose
- Database: PostgreSQL running through Docker
- Health-check worker: running through the Go backend
- Responsive viewports: 320px, 768px and desktop

All executed manual test cases passed after the fixes documented below. Rows marked `Source verified` were checked through code review, linting and the production build.

## Authentication and Session Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| AUTH-01 | Open a protected route while signed out | Redirect to the login page | Pass |
| AUTH-02 | Submit an empty or invalid login form | Field-specific validation is shown and the first invalid field receives focus | Pass |
| AUTH-03 | Submit incorrect login credentials | A clear authentication error is shown without leaving the page | Pass |
| AUTH-04 | Submit valid login credentials | The user is signed in and taken to the dashboard | Pass |
| AUTH-05 | Submit invalid registration values | Required fields, email format, password length and password matching are validated | Pass |
| AUTH-06 | Start a new registration | Registration is held pending, then the user is taken to the verification screen with the email pre-filled | Pass |
| AUTH-07 | Register an existing email address | The form remains populated and shows `Email is already registered` | Pass |
| AUTH-08 | Use the password visibility controls | Each password field can be shown and hidden independently | Pass |
| AUTH-09 | Sign out | The token is cleared, the login page is shown and a success notice appears | Pass after fix |
| AUTH-10 | Enter a password with leading or trailing spaces | The frontend rejects the value before sending it to the backend | Pass after fix |
| AUTH-11 | Reload with a valid stored token | The authenticated session and user are restored | Pass |
| AUTH-12 | `/api/auth/me` returns `404` for a removed user | The stale token is removed and the session-expired flow is triggered | Source verified |
| AUTH-13 | Sign in after requesting a protected service URL | The user returns to the originally requested service page | Pass |
| AUTH-14 | Cancel Sign out while a service form has unsaved values | The form values and authenticated session both remain intact | Pass after fix |
| AUTH-15 | Confirm Sign out while a service form has unsaved values | The changes are discarded, the session is cleared and the sign-out notice is displayed | Pass |
| AUTH-16 | Attempt to sign out while a service save is active | Sign out is blocked until the request settles | Source verified |
| AUTH-17 | Open `/verify-email` without a token | A clear check-inbox screen and resend form are displayed | Pass |
| AUTH-18 | Open `/verify-email?token=INVALID_TOKEN` | A useful invalid-or-expired-link error and resend option are shown | Pass |
| AUTH-19 | Open a valid locally logged verification link | The pending registration is verified and the success screen links back to sign in | Pass |
| AUTH-20 | Request another link for a non-pending email | The generic backend response is shown without revealing account state | Pass |
| AUTH-21 | Request another verification link while a request is active | The resend button is disabled to prevent duplicate frontend submissions | Source verified |
| AUTH-22 | Sign in after verification | The verified account can sign in and reaches the dashboard | Pass |

## Routing Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| ROUTE-01 | Open `/` | Redirect to `/dashboard` | Pass |
| ROUTE-02 | Open an unknown path | A custom not-found page is displayed | Pass |
| ROUTE-03 | Open a service route with a non-numeric ID | The page displays `Invalid service ID` | Pass |
| ROUTE-04 | Open a numeric service ID that does not exist | The page displays `Service not found` | Pass |
| ROUTE-05 | Move between lazy-loaded pages | The requested page loads with an appropriate loading state | Pass |
| ROUTE-06 | Navigate between pages | The document title and main-content focus are updated | Pass |

## Dashboard Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| DASH-01 | Load an account with no services | Summary values show zero and the first-service empty state is displayed | Pass |
| DASH-02 | Load an account with monitored services | Service cards and summary totals match the API data | Pass |
| DASH-03 | Refresh the dashboard manually | Data is reloaded and the button shows progress | Pass |
| DASH-04 | Leave the dashboard open | Data refreshes automatically after the configured interval | Pass |
| DASH-05 | Stop the backend after data has loaded | Existing data remains visible and a clear retryable error notice appears | Pass |
| DASH-06 | Restart the backend | The dashboard recovers on a later refresh without a page reload | Pass |
| DASH-07 | Display long service content | Names and URLs wrap or truncate without breaking the layout | Pass |

## Service Form, Create and Edit Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| SERV-01 | Submit an empty service form | All required validation messages are shown | Pass |
| SERV-02 | Enter an invalid URL | Only complete `http://` or `https://` URLs are accepted | Pass |
| SERV-03 | Enter invalid status, threshold or interval values | Numeric range and whole-number validation is shown | Pass |
| SERV-04 | Create a valid service | The service is saved and its detail page is displayed | Pass |
| SERV-05 | View a newly created service | The initial no-data state is shown until the first check completes | Pass |
| SERV-06 | Wait for the health-check worker and refresh | Status, summary, chart and history populate from real API data | Pass |
| SERV-07 | Open the edit page | Existing values are loaded into the form | Pass |
| SERV-08 | Save valid edits | Updated values are shown on the service detail page | Pass |
| SERV-09 | Cancel a form containing unsaved values | A discard confirmation is displayed and cancelling keeps the values | Pass |
| SERV-10 | Use browser Back or Forward with unsaved values | A discard confirmation is displayed and cancelling keeps the current page and values | Pass after fix |
| SERV-11 | Confirm browser Back after the warning | Navigation proceeds and the unfinished changes are discarded | Pass |
| SERV-12 | Save a form after adding the navigation blocker | Successful submission is not blocked and redirects normally | Pass |
| SERV-13 | Submit repeatedly while a request is active | Controls remain disabled and duplicate requests are prevented | Source verified |
| SERV-14 | Reload a page containing unsaved form values | Chrome displays its native warning and cancelling retains the page and values | Pass |
| SERV-15 | Open the edit page while the backend is unavailable | A clear loading error is shown and Try again recovers after restart | Pass |
| SERV-16 | Save an edit while the backend is unavailable | The entered values remain and a clear retryable error is shown | Pass |
| SERV-17 | Restart the backend and retry the failed edit | The preserved values save successfully without re-entry | Pass |
| SERV-18 | Use in-app navigation or browser history while a service save is active | In-app navigation is blocked until the request settles; the successful redirect still proceeds | Source verified |
| SERV-19 | Reload or close while a service save is active | Chrome displays its native warning so the user can choose whether to remain | Source verified |

## Service Detail and Monitoring Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| DETAIL-01 | Load a monitored service | Current status, URL, interval and summary statistics are displayed | Pass |
| DETAIL-02 | Load recorded health checks | The response-time chart uses chronological API data | Source verified |
| DETAIL-03 | Review the chart without relying on graphics | An accessible text summary of the measurements is available | Pass |
| DETAIL-04 | Load recent health checks | Status, time, HTTP code, response time and errors render correctly | Pass |
| DETAIL-05 | Use Show more and Show fewer | Additional history is displayed and can be collapsed | Pass |
| DETAIL-06 | View a service with no history | Helpful empty states are shown for the chart and history | Pass |
| DETAIL-07 | Refresh the detail page manually | Service, summary and history data are refreshed | Pass |
| DETAIL-08 | Open the service URL | The URL is presented as a safe external link | Source verified |

## Delete Dialog Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| DELETE-01 | Open the delete action | A clearly labelled confirmation dialog is displayed | Pass |
| DELETE-02 | Open the dialog with a keyboard | Focus begins on the safe `Keep service` action | Pass |
| DELETE-03 | Move focus with Tab and Shift+Tab | Focus remains trapped between the two dialog actions | Pass |
| DELETE-04 | Press Escape | The dialog closes and focus returns to the original Delete button | Pass |
| DELETE-05 | Confirm deletion of the temporary service | The service and its history are removed and the dashboard shows a success notice | Pass |

## Alert History Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| ALERT-01 | Open `/alerts` or `/alerts/` directly while signed in | The alert page loads with the correct title, active navigation item and main-content focus | Pass after fix |
| ALERT-02 | Load an account with no alerts | Summary values show zero and a helpful empty state is displayed | Pass |
| ALERT-03 | Load recorded alerts | Counts, severity, service links, messages and timestamps match the API data, with newest alerts first | Pass |
| ALERT-04 | Change a monitored service from available to down | One critical `Service Down` alert is created and appears globally and on the service page | Pass |
| ALERT-05 | Allow repeated failed checks while the service remains down | The original alert remains and duplicate alerts are not added | Pass |
| ALERT-06 | Recover the service and cause a later outage | The later outage creates a second alert while preserving the first | Pass |
| ALERT-07 | Refresh alert history manually | Alert data is reloaded and the refresh control reports progress | Pass |
| ALERT-08 | Leave the alert page open | Alert data refreshes automatically after 30 seconds | Pass |
| ALERT-09 | Stop the backend after alerts have loaded | Existing alerts remain visible, a retryable warning appears and the page recovers after restart | Pass |
| ALERT-10 | View alert history on a service with no alerts and one with alerts | The empty state and service-specific history are both displayed correctly | Pass |
| ALERT-11 | Inspect the service-specific alert table | The Service column is omitted and the accessible caption describes only the columns present | Pass after fix |
| ALERT-12 | Review alerts at 320px and 768px | Cards and tables adapt without page-level horizontal overflow or clipped content | Pass after fix |
| ALERT-13 | Navigate to an overflowing alert table by keyboard | The scrollable region receives focus and has a visible focus indicator | Pass after fix |
| ALERT-14 | Fail only the service-alert request | The rest of the service detail page can still load with a separate alert-history error state | Source verified |
| ALERT-15 | Open the delete confirmation for a service with alerts | The warning clearly states that the service, health checks and alert history will be removed | Pass after fix |

## Reports Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| REPORT-01 | Open `/reports` while signed in | The report page loads with the correct title, active navigation item and accessible main content | Pass |
| REPORT-02 | Load an account with no completed summaries | A helpful warming-up empty state is displayed | Pass |
| REPORT-03 | Load temporary completed summary data | Summary cards, reliability highlights and service comparison show the expected values | Pass |
| REPORT-04 | Switch between 7-day and 30-day ranges | The selected period and report data update correctly | Pass |
| REPORT-05 | Review the daily report charts | Uptime, response-time and failed-check charts render with accessible text summaries | Pass |
| REPORT-06 | Refresh a populated report | Report data reloads and the updated time changes | Pass |
| REPORT-07 | Review reports on an iPhone SE viewport | Navigation collapses correctly and the summary, charts and service comparison remain usable | Pass |
| REPORT-08 | Remove the temporary report data and refresh | The page returns to the warming-up empty state | Pass |

## Responsive and Accessibility Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| UX-01 | Review dashboard, auth, service form and detail pages at 320px | Content remains readable with no clipped actions or horizontal page overflow | Pass |
| UX-02 | Review the service detail page at 768px | Header actions, summary cards, chart and history adapt to the tablet width | Pass |
| UX-03 | Review the application at desktop width | Sidebar, content grid, chart and table use the available space correctly | Pass |
| UX-04 | Open the mobile navigation drawer | Navigation, account details and sign-out are available | Pass |
| UX-05 | Press Escape in the mobile drawer | The drawer closes and focus returns to the menu button | Pass |
| A11Y-01 | Inspect form controls | Inputs have programmatic labels, descriptions and error relationships | Pass |
| A11Y-02 | Trigger form errors | Errors use live regions and focus moves to the first invalid field | Pass |
| A11Y-03 | Navigate by keyboard | Interactive controls have visible focus and sensible keyboard order | Pass |
| A11Y-04 | Inspect page landmarks | Main navigation, main content and a skip link are present | Pass |
| A11Y-05 | Inspect status changes | Success, warning and error notices use appropriate status or alert semantics | Pass |

## Containerised Frontend Tests

| ID | Test | Expected result | Status |
|---|---|---|---|
| CONT-01 | Build the production frontend image | TypeScript, Vite and the Docker build complete successfully | Pass |
| CONT-02 | Start the full Docker Compose stack | Frontend, backend, PostgreSQL, Prometheus and Grafana start successfully | Pass |
| CONT-03 | Request `/healthz` from the frontend container | Nginx returns `200 OK` with a plain-text `ok` response | Pass |
| CONT-04 | Open or reload `/alerts`, `/services/6` and an unknown client route | Nginx serves the React application so client-side routing can resolve the page | Pass |
| CONT-05 | Request built and missing asset paths | Hashed assets use long-lived immutable caching and missing assets return `404` | Pass |
| CONT-06 | Inspect document and server headers | Application documents use `no-cache` and the Nginx version is not exposed | Pass |
| CONT-07 | Sign in and use dashboard, service and alert pages through the container | The production frontend communicates with the Dockerised API successfully | Pass |
| CONT-08 | Inspect the running frontend container | Nginx configuration is valid and the process uses the unprivileged `nginx` user | Pass |
| CONT-09 | Stop and restart the backend while the frontend remains available | Static pages remain available, stale alert data is preserved and polling recovers after restart | Pass |

## Static Verification

| Check | Command | Status |
|---|---|---|
| Oxlint | `npm run lint` | Pass |
| TypeScript and production bundle | `npm run build` | Pass |
| Whitespace and patch integrity | `git diff --check` | Pass |

## Defects Found and Resolved

### FE-QA-01 — Browser history could discard unsaved service changes

- Severity: High
- Previous behaviour: Cancel and link navigation were protected, but browser Back and Forward discarded the form without warning.
- Resolution: A shared navigation guard now covers links and browser history.
- Regression result: Cancelling the warning retains the page and values; confirming it allows navigation; successful form submissions remain unaffected.

### FE-QA-02 — Sign-out success message was not displayed

- Severity: Medium
- Previous behaviour: Authentication state changed before the navigation state reached the login page.
- Resolution: The sign-out notice is stored in session storage before the token and user state are cleared.
- Regression result: The login page reliably displays `You have signed out successfully.`

### FE-QA-03 — A deleted user could retain a stale token

- Severity: Medium
- Previous behaviour: The client cleared sessions for `401` responses, while `/api/auth/me` returns `404` when the token belongs to a removed user.
- Resolution: A `404` response from `/api/auth/me` now follows the same invalid-session cleanup flow.
- Verification: Lint, type checking and production build pass. A live deleted-user test was not performed because it would require deleting an account.

### FE-QA-04 — Password whitespace validation differed from the backend

- Severity: Low
- Previous behaviour: A password could pass client length validation and then be trimmed or rejected by the backend.
- Resolution: Registration now rejects leading and trailing password spaces with a specific field error.
- Regression result: The invalid value is stopped before any API request is made.

### FE-QA-05 — Cancelling Sign out could clear the session before navigation stopped

- Severity: High
- Previous behaviour: On a dirty service form, the logout side effects ran before the route blocker asked whether to discard the form. Cancelling could leave the user on the form but already signed out.
- Resolution: The active form now registers a shared navigation guard. Sign out must pass that guard before changing the token, user state or sign-out notice.
- Regression result: Cancelling retains the edited value and authenticated user; confirming signs out once, navigates to login and shows the success notice.

### FE-QA-06 — Navigation could race an active service save

- Severity: Medium
- Previous behaviour: Sign out, browser history, sidebar navigation or reload could leave the page while a create or update request was in progress, allowing an aborted request or a late success redirect.
- Resolution: The form now blocks in-app route and history navigation while the request is active, shows the browser's native warning for reload or close, and explicitly allows its own successful post-save redirect.
- Verification: Source review, lint, type checking and production build pass.

### FE-QA-07 — A trailing slash broke alert-page metadata and navigation state

- Severity: Low
- Previous behaviour: `/alerts/` displayed the alert page but used the not-found document title and did not mark Alerts as the active navigation item.
- Resolution: Route metadata and navigation matching now use a normalised pathname.
- Regression result: `/alerts` and `/alerts/` both show `Alert history | Cloud Monitor` and the active Alerts navigation state.

### FE-QA-08 — The alert loading screen used an invalid skeleton class

- Severity: Low
- Previous behaviour: The alert-page loading component used a class that was not defined by the shared skeleton styles and did not represent the summary cards.
- Resolution: The loading component now uses the shared skeleton class and includes summary-card placeholders.
- Verification: Lint, type checking and the production build pass.

### FE-QA-09 — Service alert history described a column that was not present

- Severity: Low
- Previous behaviour: The service-specific table caption referred to a Service column even though that column is only shown on the global alert page.
- Resolution: The accessible caption now changes with the table variant.
- Regression result: The service-specific caption accurately describes severity, type, message and created time.

### FE-QA-10 — The alert table overflow was not keyboard accessible

- Severity: Medium
- Previous behaviour: The tablet-width global table could scroll horizontally but its scroll container could not receive keyboard focus. The service-specific table also kept an unnecessarily wide minimum size.
- Resolution: Alert table wrappers are labelled focusable regions with a visible focus indicator, and the service-specific table uses a smaller minimum width.
- Regression result: The global scroll region receives keyboard focus, while the service-specific table fits without horizontal scrolling at 768px.

### FE-QA-11 — The delete warning omitted alert history

- Severity: Medium
- Previous behaviour: Deleting a service also removed its alerts, but the confirmation only warned about the service and health checks.
- Resolution: The dialog now explicitly includes alert history in the permanent deletion warning.
- Regression result: The updated warning is shown before any deletion can be confirmed.

## Known Limitations and Follow-up Work

- There is currently no automated frontend unit, component or end-to-end test suite. This report covers manual browser testing plus lint and production build verification.
- Testing was completed in Chromium-based browsers on macOS. Safari, Firefox and Windows browser testing remain future cross-browser checks.
- Logout removes the token from the browser, but the backend does not currently revoke a copied token before its 24-hour expiry.
- A backend outage during initial session restoration does not yet have a dedicated offline screen.
- The initial dashboard load requests its main datasets together, so one failed request prevents a partial initial render.
- Cross-account service isolation was not included in this test run.
- Independent detail-data failures and delete-request failures remain follow-up fault-injection cases.
- Alert history is currently limited to 50 global records and eight records on a service page, with no pagination, filters or acknowledgement controls.
- Only downtime transitions create alert records. Recovery events are not stored in alert history, and a service already down when alerting is introduced is not backfilled until it recovers and fails again.
- Successful background alert refreshes update the page silently and are not announced to screen-reader users.
- The frontend API address is set when the production image is built, so custom ports or origins must be kept in sync with backend CORS settings.
- SMTP is not configured locally. The backend logs verification links for local testing, but delivery to a real inbox still needs a provider and a separate retest.
- Prometheus and Grafana are separate operational interfaces and were not included in this React frontend test report.

## Retest Checklist

Run these checks after future frontend changes:

1. Run `npm run lint`.
2. Run `npm run build`.
3. Verify registration, login, session restoration and sign-out.
4. Create and edit a temporary service.
5. Confirm that Back, Forward, Cancel and reload protect unsaved form values.
6. Confirm that a health check populates the dashboard, chart and history.
7. Verify dashboard and detail-page recovery after an API failure.
8. Recheck the 320px, 768px and desktop layouts.
9. Recheck mobile drawer and delete dialog focus handling.
10. Verify the alert empty state, first outage, duplicate suppression, recovery and later re-outage.
11. Verify manual and automatic alert refresh, including recovery from a backend outage.
12. Build and start the production frontend container, then check `/healthz`, direct client routes and asset caching.
13. Verify reports with both empty and populated summary data, then remove the temporary fixtures.
14. Verify the verification screen, invalid link, resend response and post-verification sign-in flow.
15. Remove any test accounts and services created during testing.
