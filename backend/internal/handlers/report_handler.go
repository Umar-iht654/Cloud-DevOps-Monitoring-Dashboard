package handlers

import (
	"database/sql"
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

func maxTime(left time.Time, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}

func calculateOverviewPeriodStart(requestedStart time.Time, earliestServiceCreatedAt sql.NullTime) time.Time {
	if earliestServiceCreatedAt.Valid && earliestServiceCreatedAt.Time.After(requestedStart) {
		return earliestServiceCreatedAt.Time
	}

	return requestedStart
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

	// This gets the current time.
	now := time.Now()

	// This gets the start of today in the backend's local timezone.
	currentDayStart := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	// This calculates the oldest completed daily summary period that should be returned.
	startTime := currentDayStart.AddDate(0, 0, -days)

	// This stores the report rows returned from the database.
	var summaries []models.DailyServiceSummary

	// This reads daily summaries for the selected service and date range.
	if err := h.DB.
		Where("service_id = ? AND period_start >= ? AND period_start < ?", service.ID, startTime, currentDayStart).
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

	// This gets the start of the current hour.
	currentHourStart := time.Now().Truncate(time.Hour)

	// This calculates the oldest completed hourly summary period that should be returned.
	startTime := currentHourStart.Add(-time.Duration(hours) * time.Hour)

	// This stores the report rows returned from the database.
	var summaries []models.HourlyServiceSummary

	// This reads hourly summaries for the selected service and time range.
	if err := h.DB.
		Where("service_id = ? AND period_start >= ? AND period_start < ?", service.ID, startTime, currentHourStart).
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

	// This gets the current time.
	now := time.Now()

	// This gets the start of the current hour.
	currentHourStart := now.Truncate(time.Hour)

	// This gets the start of today in the backend's local timezone.
	currentDayStart := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	// This calculates the requested maximum lookback window.
	// For a 7-day report, this means the previous 6 completed days plus today so far.
	requestedPeriodStart := currentDayStart.AddDate(0, 0, -(days - 1))

	var earliestServiceCreatedAt sql.NullTime
	if err := h.DB.
		Model(&models.Service{}).
		Where("user_id = ?", userID).
		Select("MIN(created_at)").
		Scan(&earliestServiceCreatedAt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load report period"})
		return
	}

	periodStart := calculateOverviewPeriodStart(requestedPeriodStart, earliestServiceCreatedAt)
	periodEnd := now
	completedHourSummaryStart := maxTime(periodStart, currentDayStart)
	rawCurrentHourStart := maxTime(periodStart, currentHourStart)

	// This stores the overall totals returned from the database.
	var totals overviewReportTotals

	// This calculates totals from completed daily summaries, completed hours today and raw current-hour checks.
	totalsQuery := `
	WITH report_rows AS (
		SELECT
			total_checks,
			successful_checks,
			failed_checks,
			response_time_sample_count,
			average_response_time_ms,
			min_response_time_ms,
			max_response_time_ms
		FROM daily_service_summaries
		WHERE user_id = ?
		AND period_end > ?
		AND period_start < ?

		UNION ALL

		SELECT
			total_checks,
			successful_checks,
			failed_checks,
			response_time_sample_count,
			average_response_time_ms,
			min_response_time_ms,
			max_response_time_ms
		FROM hourly_service_summaries
		WHERE user_id = ?
		AND period_end > ?
		AND period_start < ?

		UNION ALL

		SELECT
			COUNT(*)::bigint AS total_checks,
			COALESCE(SUM(CASE WHEN health_checks.status IN ('online', 'slow') THEN 1 ELSE 0 END), 0)::bigint AS successful_checks,
			COALESCE(SUM(CASE WHEN health_checks.status = 'down' THEN 1 ELSE 0 END), 0)::bigint AS failed_checks,
			COUNT(health_checks.response_time_ms)::bigint AS response_time_sample_count,
			ROUND(COALESCE(AVG(health_checks.response_time_ms), 0))::integer AS average_response_time_ms,
			MIN(health_checks.response_time_ms) AS min_response_time_ms,
			MAX(health_checks.response_time_ms) AS max_response_time_ms
		FROM health_checks
		INNER JOIN services ON services.id = health_checks.service_id
		WHERE services.user_id = ?
		AND health_checks.checked_at >= ?
		AND health_checks.checked_at < ?
	)
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
	FROM report_rows;
`

	// This runs the totals query.
	if err := h.DB.Raw(
		totalsQuery,
		userID,
		periodStart,
		currentDayStart,
		userID,
		completedHourSummaryStart,
		currentHourStart,
		userID,
		rawCurrentHourStart,
		periodEnd,
	).Scan(&totals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load overview totals"})
		return
	}

	// This stores per-service report rows.
	var services []OverviewServiceReport

	// This calculates one report row for each owned service with available monitoring data.
	servicesQuery := `
		WITH report_rows AS (
			SELECT
				service_id,
				total_checks,
				successful_checks,
				failed_checks,
				response_time_sample_count,
				average_response_time_ms,
				min_response_time_ms,
				max_response_time_ms
			FROM daily_service_summaries
			WHERE user_id = ?
			AND period_end > ?
			AND period_start < ?

			UNION ALL

			SELECT
				service_id,
				total_checks,
				successful_checks,
				failed_checks,
				response_time_sample_count,
				average_response_time_ms,
				min_response_time_ms,
				max_response_time_ms
			FROM hourly_service_summaries
			WHERE user_id = ?
			AND period_end > ?
			AND period_start < ?

			UNION ALL

			SELECT
				health_checks.service_id,
				COUNT(*)::bigint AS total_checks,
				COALESCE(SUM(CASE WHEN health_checks.status IN ('online', 'slow') THEN 1 ELSE 0 END), 0)::bigint AS successful_checks,
				COALESCE(SUM(CASE WHEN health_checks.status = 'down' THEN 1 ELSE 0 END), 0)::bigint AS failed_checks,
				COUNT(health_checks.response_time_ms)::bigint AS response_time_sample_count,
				ROUND(COALESCE(AVG(health_checks.response_time_ms), 0))::integer AS average_response_time_ms,
				MIN(health_checks.response_time_ms) AS min_response_time_ms,
				MAX(health_checks.response_time_ms) AS max_response_time_ms
			FROM health_checks
			INNER JOIN services ON services.id = health_checks.service_id
			WHERE services.user_id = ?
			AND health_checks.checked_at >= ?
			AND health_checks.checked_at < ?
			GROUP BY health_checks.service_id
		)
		SELECT
			report_rows.service_id AS service_id,
			services.name AS service_name,
			COALESCE(SUM(report_rows.total_checks), 0)::bigint AS total_checks,
			COALESCE(SUM(report_rows.successful_checks), 0)::bigint AS successful_checks,
			COALESCE(SUM(report_rows.failed_checks), 0)::bigint AS failed_checks,
			COALESCE(SUM(report_rows.response_time_sample_count), 0)::bigint AS response_time_sample_count,
			CASE
				WHEN COALESCE(SUM(report_rows.response_time_sample_count), 0) > 0 THEN
					ROUND(
						SUM(
							report_rows.average_response_time_ms
							*
							report_rows.response_time_sample_count
						)::numeric
						/
						SUM(report_rows.response_time_sample_count)::numeric
					)::integer
				ELSE 0
			END AS average_response_time_ms,
			MIN(report_rows.min_response_time_ms) AS min_response_time_ms,
			MAX(report_rows.max_response_time_ms) AS max_response_time_ms,
			CASE
				WHEN COALESCE(SUM(report_rows.total_checks), 0) > 0 THEN
					ROUND(
						(
							SUM(report_rows.successful_checks)::numeric
							/
							SUM(report_rows.total_checks)::numeric
						) * 100,
						2
					)
				ELSE 0
			END AS uptime_percentage
		FROM report_rows
		INNER JOIN services ON services.id = report_rows.service_id
		WHERE services.user_id = ?
		GROUP BY report_rows.service_id, services.name
		ORDER BY uptime_percentage ASC, failed_checks DESC, services.name ASC;
	`

	// This runs the per-service report query.
	if err := h.DB.Raw(
		servicesQuery,
		userID,
		periodStart,
		currentDayStart,
		userID,
		completedHourSummaryStart,
		currentHourStart,
		userID,
		rawCurrentHourStart,
		periodEnd,
		userID,
	).Scan(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load service report data"})
		return
	}

	// This stores daily trend rows across all owned services.
	var daily []OverviewDailyReportPoint

	// This calculates daily trend rows across all owned services.
	// Previous completed days come from daily summaries.
	// Today so far is built from completed hourly summaries plus raw current-hour checks.
	dailyQuery := `
		WITH report_rows AS (
			SELECT
				period_start,
				period_end,
				total_checks,
				successful_checks,
				failed_checks,
				response_time_sample_count,
				average_response_time_ms,
				min_response_time_ms,
				max_response_time_ms
			FROM daily_service_summaries
			WHERE user_id = ?
			AND period_end > ?
			AND period_start < ?

			UNION ALL

			SELECT
				date_trunc('day', period_start) AS period_start,
				? AS period_end,
				total_checks,
				successful_checks,
				failed_checks,
				response_time_sample_count,
				average_response_time_ms,
				min_response_time_ms,
				max_response_time_ms
			FROM hourly_service_summaries
			WHERE user_id = ?
			AND period_end > ?
			AND period_start < ?

			UNION ALL

			SELECT
				date_trunc('day', health_checks.checked_at) AS period_start,
				? AS period_end,
				COUNT(*)::bigint AS total_checks,
				COALESCE(SUM(CASE WHEN health_checks.status IN ('online', 'slow') THEN 1 ELSE 0 END), 0)::bigint AS successful_checks,
				COALESCE(SUM(CASE WHEN health_checks.status = 'down' THEN 1 ELSE 0 END), 0)::bigint AS failed_checks,
				COUNT(health_checks.response_time_ms)::bigint AS response_time_sample_count,
				ROUND(COALESCE(AVG(health_checks.response_time_ms), 0))::integer AS average_response_time_ms,
				MIN(health_checks.response_time_ms) AS min_response_time_ms,
				MAX(health_checks.response_time_ms) AS max_response_time_ms
			FROM health_checks
			INNER JOIN services ON services.id = health_checks.service_id
			WHERE services.user_id = ?
			AND health_checks.checked_at >= ?
			AND health_checks.checked_at < ?
			GROUP BY date_trunc('day', health_checks.checked_at)
		)
		SELECT
			period_start AS period_start,
			MAX(period_end) AS period_end,
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
		FROM report_rows
		GROUP BY period_start
		ORDER BY period_start ASC;
	`

	// This runs the daily trend query.
	if err := h.DB.Raw(
		dailyQuery,
		userID,
		periodStart,
		currentDayStart,
		periodEnd,
		userID,
		completedHourSummaryStart,
		currentHourStart,
		periodEnd,
		userID,
		rawCurrentHourStart,
		periodEnd,
	).Scan(&daily).Error; err != nil {
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
