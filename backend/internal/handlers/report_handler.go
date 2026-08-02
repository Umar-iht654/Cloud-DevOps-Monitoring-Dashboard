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

// HourlyServiceReportResponse is the JSON response returned by the hourly service report endpoint.
type HourlyServiceReportResponse struct {
	// ServiceID identifies the monitored service.
	ServiceID uint `json:"service_id"`

	// ServiceName stores the human-readable service name.
	ServiceName string `json:"service_name"`

	// Hours stores the number of hours requested.
	Hours int `json:"hours"`

	// Data stores the hourly summary rows for the selected period.
	Data []models.HourlyServiceSummary `json:"data"`
}

// OverviewReportResponse is the JSON response returned by the overview report endpoint.
type OverviewReportResponse struct {
	// Range stores the requested report range, such as 7d or 30d.
	Range string `json:"range"`

	// PeriodStart stores the start of the report window.
	PeriodStart time.Time `json:"period_start"`

	// PeriodEnd stores the end of the report window.
	PeriodEnd time.Time `json:"period_end"`

	// TotalChecks stores the total number of checks across all owned services.
	TotalChecks int64 `json:"total_checks"`

	// SuccessfulChecks stores the total number of online or slow checks.
	SuccessfulChecks int64 `json:"successful_checks"`

	// FailedChecks stores the total number of down checks.
	FailedChecks int64 `json:"failed_checks"`

	// ResponseTimeSampleCount stores how many checks had response-time measurements.
	ResponseTimeSampleCount int64 `json:"response_time_sample_count"`

	// AverageResponseTimeMs stores the weighted average response time.
	AverageResponseTimeMs int `json:"average_response_time_ms"`

	// MinResponseTimeMs stores the fastest response time in the range.
	MinResponseTimeMs *int `json:"min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest response time in the range.
	MaxResponseTimeMs *int `json:"max_response_time_ms"`

	// UptimePercentage stores the overall uptime percentage for the report range.
	UptimePercentage float64 `json:"uptime_percentage"`

	// Services stores per-service report summaries.
	Services []OverviewServiceReport `json:"services"`

	// Daily stores daily trend points across all owned services.
	Daily []OverviewDailyReportPoint `json:"daily"`
}

// OverviewServiceReport stores report data for one service in the overview response.
type OverviewServiceReport struct {
	// ServiceID identifies the monitored service.
	ServiceID uint `json:"service_id" gorm:"column:service_id"`

	// ServiceName stores the human-readable service name.
	ServiceName string `json:"service_name" gorm:"column:service_name"`

	// TotalChecks stores the total checks for this service.
	TotalChecks int64 `json:"total_checks" gorm:"column:total_checks"`

	// SuccessfulChecks stores successful checks for this service.
	SuccessfulChecks int64 `json:"successful_checks" gorm:"column:successful_checks"`

	// FailedChecks stores failed checks for this service.
	FailedChecks int64 `json:"failed_checks" gorm:"column:failed_checks"`

	// ResponseTimeSampleCount stores how many checks had response-time measurements.
	ResponseTimeSampleCount int64 `json:"response_time_sample_count" gorm:"column:response_time_sample_count"`

	// AverageResponseTimeMs stores the weighted average response time for this service.
	AverageResponseTimeMs int `json:"average_response_time_ms" gorm:"column:average_response_time_ms"`

	// MinResponseTimeMs stores the fastest response time for this service.
	MinResponseTimeMs *int `json:"min_response_time_ms" gorm:"column:min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest response time for this service.
	MaxResponseTimeMs *int `json:"max_response_time_ms" gorm:"column:max_response_time_ms"`

	// UptimePercentage stores the uptime percentage for this service.
	UptimePercentage float64 `json:"uptime_percentage" gorm:"column:uptime_percentage"`
}

// OverviewDailyReportPoint stores one daily trend point in the overview response.
type OverviewDailyReportPoint struct {
	// PeriodStart stores the start of the daily period.
	PeriodStart time.Time `json:"period_start" gorm:"column:period_start"`

	// PeriodEnd stores the end of the daily period.
	PeriodEnd time.Time `json:"period_end" gorm:"column:period_end"`

	// TotalChecks stores the total checks for the day.
	TotalChecks int64 `json:"total_checks" gorm:"column:total_checks"`

	// SuccessfulChecks stores successful checks for the day.
	SuccessfulChecks int64 `json:"successful_checks" gorm:"column:successful_checks"`

	// FailedChecks stores failed checks for the day.
	FailedChecks int64 `json:"failed_checks" gorm:"column:failed_checks"`

	// ResponseTimeSampleCount stores how many checks had response-time measurements.
	ResponseTimeSampleCount int64 `json:"response_time_sample_count" gorm:"column:response_time_sample_count"`

	// AverageResponseTimeMs stores the weighted average response time for the day.
	AverageResponseTimeMs int `json:"average_response_time_ms" gorm:"column:average_response_time_ms"`

	// MinResponseTimeMs stores the fastest response time for the day.
	MinResponseTimeMs *int `json:"min_response_time_ms" gorm:"column:min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest response time for the day.
	MaxResponseTimeMs *int `json:"max_response_time_ms" gorm:"column:max_response_time_ms"`

	// UptimePercentage stores the uptime percentage for the day.
	UptimePercentage float64 `json:"uptime_percentage" gorm:"column:uptime_percentage"`
}

// overviewReportTotals stores the overall totals returned from the database.
type overviewReportTotals struct {
	// TotalChecks stores all checks in the report range.
	TotalChecks int64 `gorm:"column:total_checks"`

	// SuccessfulChecks stores all successful checks in the report range.
	SuccessfulChecks int64 `gorm:"column:successful_checks"`

	// FailedChecks stores all failed checks in the report range.
	FailedChecks int64 `gorm:"column:failed_checks"`

	// ResponseTimeSampleCount stores all response-time measurements in the report range.
	ResponseTimeSampleCount int64 `gorm:"column:response_time_sample_count"`

	// AverageResponseTimeMs stores the weighted average response time.
	AverageResponseTimeMs int `gorm:"column:average_response_time_ms"`

	// MinResponseTimeMs stores the fastest response time in the range.
	MinResponseTimeMs *int `gorm:"column:min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest response time in the range.
	MaxResponseTimeMs *int `gorm:"column:max_response_time_ms"`

	// UptimePercentage stores the overall uptime percentage.
	UptimePercentage float64 `gorm:"column:uptime_percentage"`
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

// GetServiceHourlyReport returns hourly summary data for one service owned by the logged-in user.
func (h *ReportHandler) GetServiceHourlyReport(c *gin.Context) {
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

	// This reads the optional hours query parameter.
	hoursParam := c.DefaultQuery("hours", "24")

	// This converts the hours query parameter from text into a number.
	hours, err := strconv.Atoi(hoursParam)

	// This validates the requested hourly report window.
	if err != nil || hours < 1 || hours > 168 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "hours must be between 1 and 168"})
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

	// This calculates the oldest hourly summary period that should be returned.
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	// This stores the report rows returned from the database.
	var summaries []models.HourlyServiceSummary

	// This reads hourly summaries for the selected service and time range.
	if err := h.DB.
		Where("service_id = ? AND period_start >= ?", service.ID, startTime).
		Order("period_start ASC").
		Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load hourly report"})
		return
	}

	// This returns clean JSON for the frontend reports page.
	c.JSON(http.StatusOK, HourlyServiceReportResponse{
		ServiceID:   service.ID,
		ServiceName: service.Name,
		Hours:       hours,
		Data:        summaries,
	})
}

// GetOverviewReport returns report data across all services owned by the logged-in user.
func (h *ReportHandler) GetOverviewReport(c *gin.Context) {
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

	// This reads the optional range query parameter.
	rangeParam := c.DefaultQuery("range", "7d")

	// This maps allowed range values to a number of days.
	allowedRanges := map[string]int{
		"7d":  7,
		"30d": 30,
		"90d": 90,
	}

	// This gets the number of days for the requested range.
	days, rangeIsAllowed := allowedRanges[rangeParam]

	// This rejects unsupported ranges.
	if !rangeIsAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"message": "range must be one of 7d, 30d or 90d"})
		return
	}

	// This calculates the start of the report window.
	periodStart := time.Now().AddDate(0, 0, -days)

	// This stores the end of the report window.
	periodEnd := time.Now()

	// This stores the overall totals returned from the database.
	var totals overviewReportTotals

	// This calculates totals across all daily summaries owned by the logged-in user.
	totalsQuery := `
		SELECT
			COALESCE(SUM(total_checks), 0)::bigint AS total_checks,
			COALESCE(SUM(successful_checks), 0)::bigint AS successful_checks,
			COALESCE(SUM(failed_checks), 0)::bigint AS failed_checks,
			COALESCE(SUM(response_time_sample_count), 0)::bigint AS response_time_sample_count,
			CASE
				WHEN COALESCE(SUM(response_time_sample_count), 0) > 0 THEN
					ROUND(
						SUM(average_response_time_ms * response_time_sample_count)::numeric
						/
						SUM(response_time_sample_count)::numeric
					)::integer
				ELSE 0
			END AS average_response_time_ms,
			MIN(min_response_time_ms) AS min_response_time_ms,
			MAX(max_response_time_ms) AS max_response_time_ms,
			CASE
				WHEN COALESCE(SUM(total_checks), 0) > 0 THEN
					ROUND(
						(
							SUM(successful_checks)::numeric
							/
							SUM(total_checks)::numeric
						) * 100,
						2
					)
				ELSE 0
			END AS uptime_percentage
		FROM daily_service_summaries
		WHERE user_id = ?
		AND period_start >= ?;
	`

	// This runs the totals query.
	if err := h.DB.Raw(totalsQuery, userID, periodStart).Scan(&totals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load overview totals"})
		return
	}

	// This stores per-service report rows.
	var services []OverviewServiceReport

	// This calculates one report row for each owned service with summary data.
	servicesQuery := `
		SELECT
			daily_service_summaries.service_id AS service_id,
			services.name AS service_name,
			COALESCE(SUM(daily_service_summaries.total_checks), 0)::bigint AS total_checks,
			COALESCE(SUM(daily_service_summaries.successful_checks), 0)::bigint AS successful_checks,
			COALESCE(SUM(daily_service_summaries.failed_checks), 0)::bigint AS failed_checks,
			COALESCE(SUM(daily_service_summaries.response_time_sample_count), 0)::bigint AS response_time_sample_count,
			CASE
				WHEN COALESCE(SUM(daily_service_summaries.response_time_sample_count), 0) > 0 THEN
					ROUND(
						SUM(
							daily_service_summaries.average_response_time_ms
							*
							daily_service_summaries.response_time_sample_count
						)::numeric
						/
						SUM(daily_service_summaries.response_time_sample_count)::numeric
					)::integer
				ELSE 0
			END AS average_response_time_ms,
			MIN(daily_service_summaries.min_response_time_ms) AS min_response_time_ms,
			MAX(daily_service_summaries.max_response_time_ms) AS max_response_time_ms,
			CASE
				WHEN COALESCE(SUM(daily_service_summaries.total_checks), 0) > 0 THEN
					ROUND(
						(
							SUM(daily_service_summaries.successful_checks)::numeric
							/
							SUM(daily_service_summaries.total_checks)::numeric
						) * 100,
						2
					)
				ELSE 0
			END AS uptime_percentage
		FROM daily_service_summaries
		INNER JOIN services ON services.id = daily_service_summaries.service_id
		WHERE daily_service_summaries.user_id = ?
		AND daily_service_summaries.period_start >= ?
		GROUP BY daily_service_summaries.service_id, services.name
		ORDER BY uptime_percentage ASC, failed_checks DESC, services.name ASC;
	`

	// This runs the per-service report query.
	if err := h.DB.Raw(servicesQuery, userID, periodStart).Scan(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load service report data"})
		return
	}

	// This stores daily trend rows across all owned services.
	var daily []OverviewDailyReportPoint

	// This calculates one trend row per day across all owned services.
	dailyQuery := `
		SELECT
			period_start AS period_start,
			period_end AS period_end,
			COALESCE(SUM(total_checks), 0)::bigint AS total_checks,
			COALESCE(SUM(successful_checks), 0)::bigint AS successful_checks,
			COALESCE(SUM(failed_checks), 0)::bigint AS failed_checks,
			COALESCE(SUM(response_time_sample_count), 0)::bigint AS response_time_sample_count,
			CASE
				WHEN COALESCE(SUM(response_time_sample_count), 0) > 0 THEN
					ROUND(
						SUM(average_response_time_ms * response_time_sample_count)::numeric
						/
						SUM(response_time_sample_count)::numeric
					)::integer
				ELSE 0
			END AS average_response_time_ms,
			MIN(min_response_time_ms) AS min_response_time_ms,
			MAX(max_response_time_ms) AS max_response_time_ms,
			CASE
				WHEN COALESCE(SUM(total_checks), 0) > 0 THEN
					ROUND(
						(
							SUM(successful_checks)::numeric
							/
							SUM(total_checks)::numeric
						) * 100,
						2
					)
				ELSE 0
			END AS uptime_percentage
		FROM daily_service_summaries
		WHERE user_id = ?
		AND period_start >= ?
		GROUP BY period_start, period_end
		ORDER BY period_start ASC;
	`

	// This runs the daily trend query.
	if err := h.DB.Raw(dailyQuery, userID, periodStart).Scan(&daily).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load daily report trend"})
		return
	}

	// This returns the full overview report response.
	c.JSON(http.StatusOK, OverviewReportResponse{
		Range:                   rangeParam,
		PeriodStart:             periodStart,
		PeriodEnd:               periodEnd,
		TotalChecks:             totals.TotalChecks,
		SuccessfulChecks:        totals.SuccessfulChecks,
		FailedChecks:            totals.FailedChecks,
		ResponseTimeSampleCount: totals.ResponseTimeSampleCount,
		AverageResponseTimeMs:   totals.AverageResponseTimeMs,
		MinResponseTimeMs:       totals.MinResponseTimeMs,
		MaxResponseTimeMs:       totals.MaxResponseTimeMs,
		UptimePercentage:        totals.UptimePercentage,
		Services:                services,
		Daily:                   daily,
	})
}
