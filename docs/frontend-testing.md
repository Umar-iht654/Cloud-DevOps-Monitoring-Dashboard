# Frontend Testing Report

## Overview

This document records the frontend quality assurance completed for the Cloud Service Monitoring Dashboard.

- Test date: 25 July 2026
- Tester: Ateeq Ur Rehman
- Browser: Google Chrome on macOS

The tests covered authentication, routing, dashboard behaviour, service management, monitoring history, selected error and recovery states, responsive design, keyboard accessibility and production build checks.

## Test Environment

- Local React frontend connected to the local Go API
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
| AUTH-06 | Register a new account | Account is created and the login page is shown with the email pre-filled | Pass |
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

## Known Limitations and Follow-up Work

- There is currently no automated frontend unit, component or end-to-end test suite. This report covers manual browser testing plus lint and production build verification.
- Testing was completed in Chrome on macOS. Safari, Firefox and Windows browser testing remain future cross-browser checks.
- Logout removes the token from the browser, but the backend does not currently revoke a copied token before its 24-hour expiry.
- A backend outage during initial session restoration does not yet have a dedicated offline screen.
- The initial dashboard load requests its main datasets together, so one failed request prevents a partial initial render.
- Cross-account service isolation was not included in this test run.
- Independent detail-data failures and delete-request failures remain follow-up fault-injection cases.
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
10. Remove any test accounts and services created during testing.
