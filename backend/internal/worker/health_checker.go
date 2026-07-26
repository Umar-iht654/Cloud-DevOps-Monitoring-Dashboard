package worker

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/metrics"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"gorm.io/gorm"
)

// HealthChecker stores everything needed to run background service checks.
type HealthChecker struct {
	// DB gives the health checker access to PostgreSQL through GORM.
	DB *gorm.DB

	// Client is used to send HTTP requests to monitored service URLs.
	Client *http.Client

	// ScanInterval controls how often the worker scans the database for services that need checking.
	ScanInterval time.Duration
}

// NewHealthChecker creates a new background health checker.
func NewHealthChecker(db *gorm.DB) *HealthChecker {
	// This returns a HealthChecker with database access, an HTTP client and a scan interval.
	return &HealthChecker{
		// This stores the database connection inside the health checker.
		DB: db,

		// This creates an HTTP client with a timeout so requests do not hang forever.
		Client: &http.Client{
			// This stops a health check if the service does not respond within 10 seconds.
			Timeout: 10 * time.Second,
		},

		// This makes the worker scan for due services every 10 seconds.
		ScanInterval: 10 * time.Second,
	}
}

// Start begins the background health checking loop.
func (h *HealthChecker) Start() {
	// This logs that the background worker has started.
	log.Println("Background health checker started")

	// This gives the API server a short moment to fully start before checks begin.
	time.Sleep(2 * time.Second)

	// This runs one health check scan immediately after startup.
	h.runOnce()

	// This creates a ticker that triggers repeatedly based on the scan interval.
	ticker := time.NewTicker(h.ScanInterval)

	// This stops the ticker if the Start function ever exits.
	defer ticker.Stop()

	// This keeps the worker running forever while the backend is running.
	for range ticker.C {
		// This runs another scan each time the ticker triggers.
		h.runOnce()
	}
}

// runOnce scans all services and checks the ones that are due.
func (h *HealthChecker) runOnce() {
	// This creates a slice to store services loaded from the database.
	var services []models.Service

	// This fetches all monitored services from the database.
	if err := h.DB.Find(&services).Error; err != nil {
		// This logs the database error instead of crashing the backend.
		log.Println("Failed to fetch services for health check:", err)

		// This stops the current scan because services could not be loaded.
		return
	}

	// This loops through each monitored service.
	for _, service := range services {
		// This checks whether this service is due for another health check.
		shouldCheck := h.shouldCheckService(service)

		// This skips the service if it is not due yet.
		if !shouldCheck {
			// This moves to the next service in the loop.
			continue
		}

		// This checks the service URL and records the result.
		h.checkService(service)
	}
}

// shouldCheckService decides whether a service should be checked now.
func (h *HealthChecker) shouldCheckService(service models.Service) bool {
	// This creates a variable to store the latest health check for this service.
	var latestCheck models.HealthCheck

	// This gets the most recent health check for the service.
	err := h.DB.Where("service_id = ?", service.ID).Order("checked_at DESC").First(&latestCheck).Error

	// This checks whether the service has never been checked before.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// This returns true because new services should be checked immediately.
		return true
	}

	// This checks whether another database error happened.
	if err != nil {
		// This logs the error instead of crashing the backend.
		log.Println("Failed to fetch latest health check:", err)

		// This returns false because it is safer to skip this service for now.
		return false
	}

	// This copies the service check interval into a local variable.
	intervalSeconds := service.CheckIntervalSeconds

	// This checks whether the interval is invalid.
	if intervalSeconds <= 0 {
		// This uses 60 seconds as a safe default interval.
		intervalSeconds = 60
	}

	// This calculates the next time the service is allowed to be checked.
	nextCheckTime := latestCheck.CheckedAt.Add(time.Duration(intervalSeconds) * time.Second)

	// This returns true if the current time is after the next allowed check time.
	return time.Now().After(nextCheckTime)
}

// checkService sends an HTTP request to one service and stores the result.
func (h *HealthChecker) checkService(service models.Service) {
	// This records the start time so response time can be calculated.
	startTime := time.Now()

	// This creates a GET request for the service URL.
	request, err := http.NewRequest(http.MethodGet, service.URL, nil)

	// This checks whether the request could not be created.
	if err != nil {
		// This records the service as down because the URL could not be requested.
		h.saveHealthCheck(service, "down", nil, nil, err.Error())

		// This stops the function because there is no valid request to send.
		return
	}

	// This sends the HTTP request to the monitored service.
	response, err := h.Client.Do(request)

	// This calculates how long the request took in milliseconds.
	responseTimeMs := int(time.Since(startTime).Milliseconds())

	// This checks whether the HTTP request failed completely.
	if err != nil {
		// This records the service as down because the request failed.
		h.saveHealthCheck(service, "down", nil, &responseTimeMs, err.Error())

		// This stops the function because there is no successful response to inspect.
		return
	}

	// This makes sure the response body is closed after the check is finished.
	defer response.Body.Close()

	// This stores the returned HTTP status code.
	httpStatusCode := response.StatusCode

	// This starts with an empty error message.
	errorMessage := ""

	// This starts by assuming the service is online.
	status := "online"

	// This checks whether the returned status code does not match the expected status code.
	if httpStatusCode != service.ExpectedStatusCode {
		// This marks the service as down because it returned the wrong HTTP status code.
		status = "down"

		// This stores a useful error message explaining the status code mismatch.
		errorMessage = fmt.Sprintf("expected status %d but got %d", service.ExpectedStatusCode, httpStatusCode)

		// This checks whether the response was correct but too slow.
	} else if responseTimeMs > service.SlowThresholdMs {
		// This marks the service as slow because it took longer than the configured threshold.
		status = "slow"
	}

	// This saves the health check result in the database.
	h.saveHealthCheck(service, status, &httpStatusCode, &responseTimeMs, errorMessage)
}

// saveHealthCheck stores the check result and updates the service current status.
func (h *HealthChecker) saveHealthCheck(service models.Service, status string, httpStatusCode *int, responseTimeMs *int, errorMessage string) {
	// This creates a new health check record ready to be saved.
	healthCheck := models.HealthCheck{
		// This links the health check to the monitored service.
		ServiceID: service.ID,

		// This stores whether the service was online, slow or down.
		Status: status,

		// This stores the HTTP status code if one was received.
		HTTPStatusCode: httpStatusCode,

		// This stores the response time if one was measured.
		ResponseTimeMs: responseTimeMs,

		// This stores an error message if the check failed.
		ErrorMessage: errorMessage,

		// This stores the time the check was completed.
		CheckedAt: time.Now(),
	}

	// This inserts the health check record into the health_checks table.
	if err := h.DB.Create(&healthCheck).Error; err != nil {
		// This logs the error if the health check could not be saved.
		log.Println("Failed to save health check:", err)

		// This stops the function because the health check save failed.
		return
	}

	// This creates an alert if the service has just moved into a down state.
	h.createDowntimeAlertIfNeeded(service, status, healthCheck, errorMessage)

	// This starts with a zero duration in case the response time is missing.
	healthCheckDuration := time.Duration(0)

	// This converts the response time from milliseconds into a Go duration.
	if responseTimeMs != nil {
		healthCheckDuration = time.Duration(*responseTimeMs) * time.Millisecond
	}

	// This records the health check result for Prometheus.
	metrics.RecordHealthCheck(status, healthCheckDuration)

	// This updates the current_status field on the monitored service.
	if err := h.DB.Model(&models.Service{}).Where("id = ?", service.ID).Update("current_status", status).Error; err != nil {
		// This logs the error if the service status could not be updated.
		log.Println("Failed to update service current status:", err)

		// This stops the function because the status update failed.
		return
	}

	// This logs the result of the health check.
	log.Printf("Checked service %s: %s", service.Name, status)
}

// createDowntimeAlertIfNeeded creates an alert when a service changes from not-down to down.
func (h *HealthChecker) createDowntimeAlertIfNeeded(service models.Service, status string, healthCheck models.HealthCheck, errorMessage string) {
	// This checks whether the latest check result is not down.
	if status != "down" {
		// This stops because we only create downtime alerts in this first version.
		return
	}

	// This checks whether the service was already down before this check.
	if service.CurrentStatus == "down" {
		// This stops because we do not want duplicate alerts every time the worker checks an already-down service.
		return
	}

	// This creates a default message if the health checker did not provide one.
	if errorMessage == "" {
		// This gives the alert a useful fallback explanation.
		errorMessage = "The service failed its health check."
	}

	// This creates a short title for the alert.
	title := fmt.Sprintf("%s is down", service.Name)

	// This creates a longer message with the service URL and the failure reason.
	message := fmt.Sprintf("%s is currently down. URL: %s. Reason: %s", service.Name, service.URL, errorMessage)

	// This creates the alert record that will be saved in PostgreSQL.
	alert := models.Alert{
		// This stores the user who owns the service.
		UserID: service.UserID,

		// This stores the service that triggered the alert.
		ServiceID: service.ID,

		// This links the alert to the health check that triggered it.
		HealthCheckID: &healthCheck.ID,

		// This stores the alert type.
		Type: "service_down",

		// This marks downtime as a critical alert.
		Severity: "critical",

		// This stores the short alert title.
		Title: title,

		// This stores the detailed alert message.
		Message: message,
	}

	// This inserts the alert into the alerts table.
	if err := h.DB.Create(&alert).Error; err != nil {
		// This logs the error but does not stop the health checker, because monitoring should continue.
		log.Println("Failed to create downtime alert:", err)

		// This stops only the alert creation function.
		return
	}

	// This logs that an alert was created.
	log.Printf("Created downtime alert for service %s", service.Name)
}
