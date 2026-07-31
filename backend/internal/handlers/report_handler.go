package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReportHandler stores the database connection for report endpoints.
type ReportHandler struct {
	// DB gives the handler access to PostgreSQL through GORM.
	DB *gorm.DB
}

// NewReportHandler creates a new report handler.
func NewReportHandler(db *gorm.DB) *ReportHandler {
	// This returns a handler with database access.
	return &ReportHandler{DB: db}
}

// DailyServiceReportResponse is the JSON response returned by the daily service report endpoint.
type DailyServiceReportResponse struct {
	// ServiceID identifies the monitored service.
	ServiceID uint `json:"service_id"`

	// ServiceName stores the human-readable service name.
	ServiceName string `json:"service_name"`

	// Days stores the number of days requested.
	Days int `json:"days"`

	// Data stores the daily summary rows for the selected period.
	Data []models.DailyServiceSummary `json:"data"`
}

// GetServiceDailyReport returns daily summary data for one service owned by the logged-in user.
func (h *ReportHandler) GetServiceDailyReport(c *gin.Context) {
	// This gets the logged-in user's ID from the JWT middleware.
	userIDValue, exists := c.Get("userID")

	// This checks that the JWT middleware actually provided a user ID.
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	// This converts the user ID into the expected uint type.
	userID, ok := userIDValue.(uint)

	// This checks that the stored user ID has the expected type.
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user context"})
		return
	}

	// This reads the service ID from the URL path.
	serviceIDParam := c.Param("id")

	// This converts the service ID from text into an unsigned integer.
	serviceIDUint64, err := strconv.ParseUint(serviceIDParam, 10, 64)

	// This handles invalid service IDs such as letters or negative values.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid service ID"})
		return
	}

	// This converts the parsed service ID into a uint for GORM.
	serviceID := uint(serviceIDUint64)

	// This reads the optional days query parameter.
	daysParam := c.DefaultQuery("days", "30")

	// This converts the days query parameter from text into a number.
	days, err := strconv.Atoi(daysParam)

	// This validates the requested report window.
	if err != nil || days < 1 || days > 90 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "days must be between 1 and 90"})
		return
	}

	// This stores the service if it exists and belongs to the logged-in user.
	var service models.Service

	// This checks service ownership before returning any report data.
	if err := h.DB.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		// This returns 404 only when the service does not exist or does not belong to the user.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
			return
		}

		// This returns 500 for unexpected database errors.
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check service ownership"})
		return
	}

	// This calculates the oldest daily summary period that should be returned.
	startTime := time.Now().AddDate(0, 0, -days)

	// This stores the report rows returned from the database.
	var summaries []models.DailyServiceSummary

	// This reads daily summaries for the selected service and date range.
	if err := h.DB.
		Where("service_id = ? AND period_start >= ?", service.ID, startTime).
		Order("period_start ASC").
		Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load daily report"})
		return
	}

	// This returns clean JSON for the frontend reports page.
	c.JSON(http.StatusOK, DailyServiceReportResponse{
		ServiceID:   service.ID,
		ServiceName: service.Name,
		Days:        days,
		Data:        summaries,
	})
}
