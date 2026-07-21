package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthCheckHandler stores the dependencies needed by health check history routes.
type HealthCheckHandler struct {
	// DB gives the health check handler access to PostgreSQL through GORM.
	DB *gorm.DB
}

// ServiceSummaryResponse defines the response body for a service summary.
type ServiceSummaryResponse struct {
	// ServiceID stores the ID of the monitored service.
	ServiceID uint `json:"service_id"`

	// ServiceName stores the name of the monitored service.
	ServiceName string `json:"service_name"`

	// CurrentStatus stores the latest known service status.
	CurrentStatus string `json:"current_status"`

	// TotalChecks stores the total number of health checks recorded for this service.
	TotalChecks int64 `json:"total_checks"`

	// SuccessfulChecks stores checks where the service was online or slow.
	SuccessfulChecks int64 `json:"successful_checks"`

	// FailedChecks stores checks where the service was down.
	FailedChecks int64 `json:"failed_checks"`

	// UptimePercentage stores the percentage of checks where the service was reachable.
	UptimePercentage float64 `json:"uptime_percentage"`

	// AverageResponseTimeMs stores the average response time across recorded checks.
	AverageResponseTimeMs int `json:"average_response_time_ms"`

	// LastCheckedAt stores the timestamp of the most recent health check.
	LastCheckedAt *time.Time `json:"last_checked_at"`

	// LastDownAt stores the timestamp of the most recent failed check.
	LastDownAt *time.Time `json:"last_down_at"`
}

// NewHealthCheckHandler creates a new HealthCheckHandler with database access.
func NewHealthCheckHandler(db *gorm.DB) *HealthCheckHandler {
	// This returns a pointer to a HealthCheckHandler so routes can use its methods.
	return &HealthCheckHandler{
		DB: db,
	}
}

// getOwnedService loads a service only if it belongs to the authenticated user.
func (h *HealthCheckHandler) getOwnedService(c *gin.Context) (*models.Service, bool) {
	// This gets the authenticated user's ID from the request context.
	userID, ok := getUserIDFromContext(c)

	// This checks whether the user ID was missing or invalid.
	if !ok {
		// This returns a 401 response because the user is not authenticated.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authenticated",
		})

		// This returns nil and false because the service cannot be loaded.
		return nil, false
	}

	// This reads the service ID from the URL parameter.
	serviceIDParam := c.Param("id")

	// This converts the service ID from a string into an unsigned integer.
	serviceID, err := strconv.ParseUint(serviceIDParam, 10, 64)

	// This checks whether the service ID was invalid.
	if err != nil {
		// This returns a 400 response because the service ID is not valid.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid service ID",
		})

		// This returns nil and false because the service ID cannot be used.
		return nil, false
	}

	// This creates a variable to store the service found in the database.
	var service models.Service

	// This searches for a service with the matching ID and authenticated user ID.
	err = h.DB.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error

	// This checks whether the service was not found.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// This returns a 404 response if the service does not exist or does not belong to the user.
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Service not found",
		})

		// This returns nil and false because there is no service to use.
		return nil, false
	}

	// This checks whether another database error happened.
	if err != nil {
		// This returns a 500 response if the database query failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch service",
		})

		// This returns nil and false because the service could not be loaded.
		return nil, false
	}

	// This returns the owned service and true because it was loaded successfully.
	return &service, true
}

// GetHealthChecks returns recent health check records for one owned service.
func (h *HealthCheckHandler) GetHealthChecks(c *gin.Context) {
	// This loads the service only if it belongs to the logged-in user.
	service, ok := h.getOwnedService(c)

	// This stops the handler if the service could not be loaded.
	if !ok {
		// This stops the function because getOwnedService already returned the response.
		return
	}

	// This sets the default number of health checks to return.
	limit := 50

	// This reads the optional limit query parameter from the URL.
	limitParam := c.Query("limit")

	// This checks whether a limit query parameter was provided.
	if limitParam != "" {
		// This tries to convert the limit query parameter into an integer.
		parsedLimit, err := strconv.Atoi(limitParam)

		// This checks whether the limit is valid and greater than zero.
		if err == nil && parsedLimit > 0 {
			// This updates the limit using the valid query parameter.
			limit = parsedLimit
		}
	}

	// This prevents very large responses by capping the limit at 200.
	if limit > 200 {
		// This sets the maximum allowed limit to 200.
		limit = 200
	}

	// This creates a slice to store the health checks found for the service.
	var healthChecks []models.HealthCheck

	// This fetches recent health checks for the service, newest first.
	if err := h.DB.Where("service_id = ?", service.ID).Order("checked_at DESC").Limit(limit).Find(&healthChecks).Error; err != nil {
		// This returns a 500 response if the health checks could not be loaded.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch health checks",
		})

		// This stops the handler because the database query failed.
		return
	}

	// This returns the service ID and its recent health checks.
	c.JSON(http.StatusOK, gin.H{
		"service_id":     service.ID,
		"health_checks":  healthChecks,
		"returned_count": len(healthChecks),
	})
}

// GetServiceSummary returns calculated monitoring statistics for one owned service.
func (h *HealthCheckHandler) GetServiceSummary(c *gin.Context) {
	// This loads the service only if it belongs to the logged-in user.
	service, ok := h.getOwnedService(c)

	// This stops the handler if the service could not be loaded.
	if !ok {
		// This stops the function because getOwnedService already returned the response.
		return
	}

	// This creates a variable to store the total number of checks.
	var totalChecks int64

	// This counts all health checks for the service.
	if err := h.DB.Model(&models.HealthCheck{}).Where("service_id = ?", service.ID).Count(&totalChecks).Error; err != nil {
		// This returns a 500 response if the total check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate total checks",
		})

		// This stops the handler because the summary cannot be calculated.
		return
	}

	// This creates a variable to store checks where the service was reachable.
	var successfulChecks int64

	// This counts checks where the service was online or slow.
	if err := h.DB.Model(&models.HealthCheck{}).Where("service_id = ? AND status IN ?", service.ID, []string{"online", "slow"}).Count(&successfulChecks).Error; err != nil {
		// This returns a 500 response if the successful check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate successful checks",
		})

		// This stops the handler because the summary cannot be calculated.
		return
	}

	// This creates a variable to store checks where the service was down.
	var failedChecks int64

	// This counts checks where the service was down.
	if err := h.DB.Model(&models.HealthCheck{}).Where("service_id = ? AND status = ?", service.ID, "down").Count(&failedChecks).Error; err != nil {
		// This returns a 500 response if the failed check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate failed checks",
		})

		// This stops the handler because the summary cannot be calculated.
		return
	}

	// This creates a variable for the uptime percentage.
	uptimePercentage := 0.0

	// This checks whether there are any health checks before dividing.
	if totalChecks > 0 {
		// This calculates uptime as online-or-slow checks divided by total checks.
		uptimePercentage = (float64(successfulChecks) / float64(totalChecks)) * 100

		// This rounds the uptime percentage to two decimal places.
		uptimePercentage = math.Round(uptimePercentage*100) / 100
	}

	// This creates a variable to store the average response time.
	var averageResponseTime float64

	// This calculates the average response time, treating missing values as ignored by AVG.
	if err := h.DB.Model(&models.HealthCheck{}).Where("service_id = ? AND response_time_ms IS NOT NULL", service.ID).Select("COALESCE(AVG(response_time_ms), 0)").Scan(&averageResponseTime).Error; err != nil {
		// This returns a 500 response if the average response time calculation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate average response time",
		})

		// This stops the handler because the summary cannot be calculated.
		return
	}

	// This creates a variable to store the most recent health check.
	var latestCheck models.HealthCheck

	// This gets the latest health check for the service.
	latestCheckErr := h.DB.Where("service_id = ?", service.ID).Order("checked_at DESC").First(&latestCheck).Error

	// This creates a nullable last checked timestamp for the response.
	var lastCheckedAt *time.Time

	// This checks whether a latest health check was found.
	if latestCheckErr == nil {
		// This stores the latest checked timestamp.
		lastCheckedAt = &latestCheck.CheckedAt

		// This checks whether the error was something other than record not found.
	} else if !errors.Is(latestCheckErr, gorm.ErrRecordNotFound) {
		// This returns a 500 response if the latest check query failed unexpectedly.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch latest health check",
		})

		// This stops the handler because the summary cannot be safely completed.
		return
	}

	// This creates a variable to store the most recent down health check.
	var lastDownCheck models.HealthCheck

	// This gets the latest failed health check for the service.
	lastDownErr := h.DB.Where("service_id = ? AND status = ?", service.ID, "down").Order("checked_at DESC").First(&lastDownCheck).Error

	// This creates a nullable last down timestamp for the response.
	var lastDownAt *time.Time

	// This checks whether a latest down check was found.
	if lastDownErr == nil {
		// This stores the latest downtime timestamp.
		lastDownAt = &lastDownCheck.CheckedAt

		// This checks whether the error was something other than record not found.
	} else if !errors.Is(lastDownErr, gorm.ErrRecordNotFound) {
		// This returns a 500 response if the last down query failed unexpectedly.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch last downtime",
		})

		// This stops the handler because the summary cannot be safely completed.
		return
	}

	// This builds the final service summary response.
	summary := ServiceSummaryResponse{
		ServiceID:             service.ID,
		ServiceName:           service.Name,
		CurrentStatus:         service.CurrentStatus,
		TotalChecks:           totalChecks,
		SuccessfulChecks:      successfulChecks,
		FailedChecks:          failedChecks,
		UptimePercentage:      uptimePercentage,
		AverageResponseTimeMs: int(math.Round(averageResponseTime)),
		LastCheckedAt:         lastCheckedAt,
		LastDownAt:            lastDownAt,
	}

	// This returns the calculated service summary.
	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
