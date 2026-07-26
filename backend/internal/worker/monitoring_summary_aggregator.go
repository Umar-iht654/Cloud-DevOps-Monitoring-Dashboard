package worker

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// MonitoringSummaryAggregator creates hourly and daily summaries from raw health checks.
type MonitoringSummaryAggregator struct {
	// DB gives the aggregator access to PostgreSQL through GORM.
	DB *gorm.DB

	// AggregationInterval controls how often aggregation runs.
	AggregationInterval time.Duration
}

// NewMonitoringSummaryAggregator creates a new monitoring summary aggregator.
func NewMonitoringSummaryAggregator(db *gorm.DB) *MonitoringSummaryAggregator {
	// This returns an aggregator with database access and a default hourly interval.
	return &MonitoringSummaryAggregator{
		// This stores the database connection.
		DB: db,

		// This runs aggregation once every hour.
		AggregationInterval: time.Hour,
	}
}

// Start begins the summary aggregation loop.
func (a *MonitoringSummaryAggregator) Start() {
	// This logs that the aggregation worker has started.
	log.Println("Monitoring summary aggregator started")

	// This creates a ticker that runs aggregation repeatedly.
	ticker := time.NewTicker(a.AggregationInterval)

	// This stops the ticker if Start ever exits.
	defer ticker.Stop()

	// This keeps the aggregation worker running while the backend is running.
	for range ticker.C {
		// This runs another aggregation pass.
		a.RunOnce()
	}
}

// RunOnce creates or updates hourly and daily monitoring summaries.
func (a *MonitoringSummaryAggregator) RunOnce() {
	// This creates hourly summaries from raw health checks.
	if err := a.aggregateHourlySummaries(); err != nil {
		// This logs the error without crashing the backend.
		log.Println("Hourly summary aggregation failed:", err)

		// This stops because daily summaries depend on hourly summaries.
		return
	}

	// This creates daily summaries from hourly summaries.
	if err := a.aggregateDailySummaries(); err != nil {
		// This logs the error without crashing the backend.
		log.Println("Daily summary aggregation failed:", err)

		// This stops the current aggregation run.
		return
	}

	// This logs that aggregation completed.
	log.Println("Monitoring summary aggregation completed")
}

// aggregateHourlySummaries creates one summary row per service per completed hour.
func (a *MonitoringSummaryAggregator) aggregateHourlySummaries() error {
	// This calculates the start of the current hour.
	currentHourStart := time.Now().Truncate(time.Hour)

	// This query creates hourly summaries from raw health checks.
	query := `
		INSERT INTO hourly_service_summaries (
			user_id,
			service_id,
			period_start,
			period_end,
			total_checks,
			successful_checks,
			failed_checks,
			response_time_sample_count,
			average_response_time_ms,
			min_response_time_ms,
			max_response_time_ms,
			uptime_percentage,
			created_at,
			updated_at
		)
		SELECT
			services.user_id,
			health_checks.service_id,
			date_trunc('hour', health_checks.checked_at) AS period_start,
			date_trunc('hour', health_checks.checked_at) + INTERVAL '1 hour' AS period_end,
			COUNT(*) AS total_checks,
			SUM(CASE WHEN health_checks.status IN ('online', 'slow') THEN 1 ELSE 0 END) AS successful_checks,
			SUM(CASE WHEN health_checks.status = 'down' THEN 1 ELSE 0 END) AS failed_checks,
			COUNT(health_checks.response_time_ms) AS response_time_sample_count,
			ROUND(COALESCE(AVG(health_checks.response_time_ms), 0))::integer AS average_response_time_ms,
			MIN(health_checks.response_time_ms) AS min_response_time_ms,
			MAX(health_checks.response_time_ms) AS max_response_time_ms,
			ROUND(
				(
					SUM(CASE WHEN health_checks.status IN ('online', 'slow') THEN 1 ELSE 0 END)::numeric
					/
					COUNT(*)::numeric
				) * 100,
				2
			) AS uptime_percentage,
			NOW() AS created_at,
			NOW() AS updated_at
		FROM health_checks
		INNER JOIN services ON services.id = health_checks.service_id
		WHERE health_checks.checked_at < ?
		GROUP BY
			services.user_id,
			health_checks.service_id,
			date_trunc('hour', health_checks.checked_at)
		ON CONFLICT (service_id, period_start)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			period_end = EXCLUDED.period_end,
			total_checks = EXCLUDED.total_checks,
			successful_checks = EXCLUDED.successful_checks,
			failed_checks = EXCLUDED.failed_checks,
			response_time_sample_count = EXCLUDED.response_time_sample_count,
			average_response_time_ms = EXCLUDED.average_response_time_ms,
			min_response_time_ms = EXCLUDED.min_response_time_ms,
			max_response_time_ms = EXCLUDED.max_response_time_ms,
			uptime_percentage = EXCLUDED.uptime_percentage,
			updated_at = NOW();
	`

	// This runs the hourly aggregation query.
	if err := a.DB.Exec(query, currentHourStart).Error; err != nil {
		// This returns a useful wrapped error to the caller.
		return fmt.Errorf("failed to aggregate hourly summaries: %w", err)
	}

	// This returns nil because hourly aggregation succeeded.
	return nil
}

// aggregateDailySummaries creates one summary row per service per completed day.
func (a *MonitoringSummaryAggregator) aggregateDailySummaries() error {
	// This query creates daily summaries from hourly summaries.
	query := `
		INSERT INTO daily_service_summaries (
			user_id,
			service_id,
			period_start,
			period_end,
			total_checks,
			successful_checks,
			failed_checks,
			response_time_sample_count,
			average_response_time_ms,
			min_response_time_ms,
			max_response_time_ms,
			uptime_percentage,
			created_at,
			updated_at
		)
		SELECT
			hourly_service_summaries.user_id,
			hourly_service_summaries.service_id,
			date_trunc('day', hourly_service_summaries.period_start) AS period_start,
			date_trunc('day', hourly_service_summaries.period_start) + INTERVAL '1 day' AS period_end,
			SUM(hourly_service_summaries.total_checks) AS total_checks,
			SUM(hourly_service_summaries.successful_checks) AS successful_checks,
			SUM(hourly_service_summaries.failed_checks) AS failed_checks,
			SUM(hourly_service_summaries.response_time_sample_count) AS response_time_sample_count,
			CASE
				WHEN SUM(hourly_service_summaries.response_time_sample_count) > 0 THEN
					ROUND(
						SUM(
							hourly_service_summaries.average_response_time_ms
							*
							hourly_service_summaries.response_time_sample_count
						)::numeric
						/
						SUM(hourly_service_summaries.response_time_sample_count)::numeric
					)::integer
				ELSE 0
			END AS average_response_time_ms,
			MIN(hourly_service_summaries.min_response_time_ms) AS min_response_time_ms,
			MAX(hourly_service_summaries.max_response_time_ms) AS max_response_time_ms,
			ROUND(
				(
					SUM(hourly_service_summaries.successful_checks)::numeric
					/
					SUM(hourly_service_summaries.total_checks)::numeric
				) * 100,
				2
			) AS uptime_percentage,
			NOW() AS created_at,
			NOW() AS updated_at
		FROM hourly_service_summaries
		WHERE hourly_service_summaries.period_start < date_trunc('day', NOW())
		GROUP BY
			hourly_service_summaries.user_id,
			hourly_service_summaries.service_id,
			date_trunc('day', hourly_service_summaries.period_start)
		ON CONFLICT (service_id, period_start)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			period_end = EXCLUDED.period_end,
			total_checks = EXCLUDED.total_checks,
			successful_checks = EXCLUDED.successful_checks,
			failed_checks = EXCLUDED.failed_checks,
			response_time_sample_count = EXCLUDED.response_time_sample_count,
			average_response_time_ms = EXCLUDED.average_response_time_ms,
			min_response_time_ms = EXCLUDED.min_response_time_ms,
			max_response_time_ms = EXCLUDED.max_response_time_ms,
			uptime_percentage = EXCLUDED.uptime_percentage,
			updated_at = NOW();
	`

	// This runs the daily aggregation query.
	if err := a.DB.Exec(query).Error; err != nil {
		// This returns a useful wrapped error to the caller.
		return fmt.Errorf("failed to aggregate daily summaries: %w", err)
	}

	// This returns nil because daily aggregation succeeded.
	return nil
}
