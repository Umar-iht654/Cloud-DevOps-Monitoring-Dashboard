package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler stores the dependencies needed by dashboard routes.
type DashboardHandler struct {
	// DB gives the dashboard handler access to PostgreSQL through GORM.
	DB *gorm.DB
}

// DashboardSummaryResponse defines the response body for the main dashboard summary.
type DashboardSummaryResponse struct {
	// TotalServices stores the total number of services owned by the logged-in user.
	TotalServices int64 `json:"total_services"`

	// OnlineServices stores the number of services currently marked as online.
	OnlineServices int64 `json:"online_services"`

	// SlowServices stores the number of services currently marked as slow.
	SlowServices int64 `json:"slow_services"`

	// DownServices stores the number of services currently marked as down.
	DownServices int64 `json:"down_services"`

	// UnknownServices stores the number of services that have not been checked yet.
	UnknownServices int64 `json:"unknown_services"`

	// TotalChecks stores the total number of health checks across the user's services.
	TotalChecks int64 `json:"total_checks"`

	// SuccessfulChecks stores the number of health checks where services were reachable.
	SuccessfulChecks int64 `json:"successful_checks"`

	// FailedChecks stores the number of health checks where services were down.
	FailedChecks int64 `json:"failed_checks"`

	// AverageUptimePercentage stores overall uptime across the user's monitored services.
	AverageUptimePercentage float64 `json:"average_uptime_percentage"`

	// AverageResponseTimeMs stores the average response time across the user's health checks.
	AverageResponseTimeMs int `json:"average_response_time_ms"`

	// LastCheckedAt stores the latest health check timestamp across the user's services.
	LastCheckedAt *time.Time `json:"last_checked_at"`
}

// NewDashboardHandler creates a new DashboardHandler with database access.
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	// This returns a pointer to a DashboardHandler so routes can use its methods.
	return &DashboardHandler{
		DB: db,
	}
}

// GetSummary returns the main dashboard summary for the logged-in user.
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	// This gets the authenticated user's ID from the request context.
	userID, ok := getUserIDFromContext(c)

	// This checks whether the user ID was missing or invalid.
	if !ok {
		// This returns a 401 response because the user is not authenticated.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authenticated",
		})

		// This stops the handler because dashboard data cannot be loaded without a user.
		return
	}

	// This creates a variable to store the total number of services.
	var totalServices int64

	// This counts all services owned by the logged-in user.
	if err := h.DB.Model(&models.Service{}).Where("user_id = ?", userID).Count(&totalServices).Error; err != nil {
		// This returns a 500 response if the total service count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate total services",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the number of online services.
	var onlineServices int64

	// This counts services owned by the user that are currently online.
	if err := h.DB.Model(&models.Service{}).Where("user_id = ? AND current_status = ?", userID, "online").Count(&onlineServices).Error; err != nil {
		// This returns a 500 response if the online service count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate online services",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the number of slow services.
	var slowServices int64

	// This counts services owned by the user that are currently slow.
	if err := h.DB.Model(&models.Service{}).Where("user_id = ? AND current_status = ?", userID, "slow").Count(&slowServices).Error; err != nil {
		// This returns a 500 response if the slow service count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate slow services",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the number of down services.
	var downServices int64

	// This counts services owned by the user that are currently down.
	if err := h.DB.Model(&models.Service{}).Where("user_id = ? AND current_status = ?", userID, "down").Count(&downServices).Error; err != nil {
		// This returns a 500 response if the down service count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate down services",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the number of unknown services.
	var unknownServices int64

	// This counts services owned by the user that have not been checked yet.
	if err := h.DB.Model(&models.Service{}).Where("user_id = ? AND current_status = ?", userID, "unknown").Count(&unknownServices).Error; err != nil {
		// This returns a 500 response if the unknown service count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate unknown services",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the total number of health checks.
	var totalChecks int64

	// This counts health checks linked to services owned by the logged-in user.
	if err := h.DB.Model(&models.HealthCheck{}).
		Joins("JOIN services ON services.id = health_checks.service_id").
		Where("services.user_id = ?", userID).
		Count(&totalChecks).Error; err != nil {
		// This returns a 500 response if the total health check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate total health checks",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store successful health checks.
	var successfulChecks int64

	// This counts checks where the user's services were online or slow.
	if err := h.DB.Model(&models.HealthCheck{}).
		Joins("JOIN services ON services.id = health_checks.service_id").
		Where("services.user_id = ? AND health_checks.status IN ?", userID, []string{"online", "slow"}).
		Count(&successfulChecks).Error; err != nil {
		// This returns a 500 response if the successful check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate successful health checks",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store failed health checks.
	var failedChecks int64

	// This counts checks where the user's services were down.
	if err := h.DB.Model(&models.HealthCheck{}).
		Joins("JOIN services ON services.id = health_checks.service_id").
		Where("services.user_id = ? AND health_checks.status = ?", userID, "down").
		Count(&failedChecks).Error; err != nil {
		// This returns a 500 response if the failed check count failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate failed health checks",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable for the overall average uptime percentage.
	averageUptimePercentage := 0.0

	// This checks whether there are health checks before dividing.
	if totalChecks > 0 {
		// This calculates uptime as successful checks divided by total checks.
		averageUptimePercentage = (float64(successfulChecks) / float64(totalChecks)) * 100

		// This rounds the uptime percentage to two decimal places.
		averageUptimePercentage = math.Round(averageUptimePercentage*100) / 100
	}

	// This creates a variable to store the average response time.
	var averageResponseTime float64

	// This calculates the average response time across all user-owned service checks.
	if err := h.DB.Model(&models.HealthCheck{}).
		Joins("JOIN services ON services.id = health_checks.service_id").
		Where("services.user_id = ? AND health_checks.response_time_ms IS NOT NULL", userID).
		Select("COALESCE(AVG(health_checks.response_time_ms), 0)").
		Scan(&averageResponseTime).Error; err != nil {
		// This returns a 500 response if the average response time calculation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to calculate average response time",
		})

		// This stops the handler because the summary cannot be calculated safely.
		return
	}

	// This creates a variable to store the latest health check.
	var latestCheck models.HealthCheck

	// This tries to find the most recent health check across the user's services.
	latestCheckErr := h.DB.Model(&models.HealthCheck{}).
		Joins("JOIN services ON services.id = health_checks.service_id").
		Where("services.user_id = ?", userID).
		Order("health_checks.checked_at DESC").
		First(&latestCheck).Error

	// This creates a nullable timestamp for the latest health check time.
	var lastCheckedAt *time.Time

	// This checks whether a latest health check was found.
	if latestCheckErr == nil {
		// This stores the latest health check timestamp.
		lastCheckedAt = &latestCheck.CheckedAt
	}

	// This builds the dashboard summary response.
	summary := DashboardSummaryResponse{
		TotalServices:           totalServices,
		OnlineServices:          onlineServices,
		SlowServices:            slowServices,
		DownServices:            downServices,
		UnknownServices:         unknownServices,
		TotalChecks:             totalChecks,
		SuccessfulChecks:        successfulChecks,
		FailedChecks:            failedChecks,
		AverageUptimePercentage: averageUptimePercentage,
		AverageResponseTimeMs:   int(math.Round(averageResponseTime)),
		LastCheckedAt:           lastCheckedAt,
	}

	// This returns the dashboard summary to the frontend.
	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
