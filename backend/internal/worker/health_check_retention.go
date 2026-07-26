package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"gorm.io/gorm"
)

// HealthCheckRetentionCleaner deletes old raw health-check records.
type HealthCheckRetentionCleaner struct {
	// DB gives the cleaner access to PostgreSQL through GORM.
	DB *gorm.DB

	// RetentionDays controls how many days of raw health checks are kept.
	RetentionDays int

	// CleanupInterval controls how often the cleaner runs.
	CleanupInterval time.Duration
}

// NewHealthCheckRetentionCleaner creates a new health-check retention cleaner.
func NewHealthCheckRetentionCleaner(db *gorm.DB, retentionDays int, cleanupIntervalHours int) *HealthCheckRetentionCleaner {
	// This returns a cleaner with database access and retention settings.
	return &HealthCheckRetentionCleaner{
		// This stores the database connection.
		DB: db,

		// This stores how many days of raw health checks should be kept.
		RetentionDays: retentionDays,

		// This converts the cleanup interval from hours into a Go duration.
		CleanupInterval: time.Duration(cleanupIntervalHours) * time.Hour,
	}
}

// Start begins the retention cleanup loop.
func (c *HealthCheckRetentionCleaner) Start() {
	// This disables cleanup when retention days is set to 0.
	if c.RetentionDays <= 0 {
		// This logs that retention cleanup is disabled.
		log.Println("Health check retention cleanup disabled")

		// This stops the worker from running.
		return
	}

	// This logs that the retention cleaner has started.
	log.Printf("Health check retention cleaner started: keeping %d days of raw checks", c.RetentionDays)

	// This runs cleanup once at startup.
	c.RunOnce()

	// This creates a ticker that runs cleanup repeatedly.
	ticker := time.NewTicker(c.CleanupInterval)

	// This stops the ticker if Start ever exits.
	defer ticker.Stop()

	// This keeps the cleanup worker running while the backend is running.
	for range ticker.C {
		// This runs another cleanup pass.
		c.RunOnce()
	}
}

// RunOnce deletes health checks older than the configured retention period.
func (c *HealthCheckRetentionCleaner) RunOnce() {
	// This disables cleanup when retention days is set to 0.
	if c.RetentionDays <= 0 {
		// This stops because cleanup is disabled.
		return
	}

	// This calculates the oldest timestamp that should be kept.
	cutoffTime := time.Now().AddDate(0, 0, -c.RetentionDays)

	// This stores how many health-check rows were deleted.
	var deletedRows int64

	// This runs alert unlinking and health-check deletion inside one database transaction.
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		// This selects old health-check IDs that are about to be deleted.
		oldHealthChecks := tx.Model(&models.HealthCheck{}).
			Select("id").
			Where("checked_at < ?", cutoffTime)

		// This removes old health-check references from alerts before deleting the old checks.
		if err := tx.Model(&models.Alert{}).
			Where("health_check_id IN (?)", oldHealthChecks).
			Updates(map[string]interface{}{
				"health_check_id": nil,
			}).Error; err != nil {
			// This returns the error so the transaction rolls back.
			return fmt.Errorf("failed to unlink old health checks from alerts: %w", err)
		}

		// This deletes old raw health-check rows.
		result := tx.Where("checked_at < ?", cutoffTime).Delete(&models.HealthCheck{})

		// This checks whether the delete failed.
		if result.Error != nil {
			// This returns the error so the transaction rolls back.
			return fmt.Errorf("failed to delete old health checks: %w", result.Error)
		}

		// This stores the number of deleted rows for logging after the transaction.
		deletedRows = result.RowsAffected

		// This commits the transaction because cleanup succeeded.
		return nil
	})

	// This checks whether the cleanup transaction failed.
	if err != nil {
		// This logs the error without crashing the backend.
		log.Println("Health check retention cleanup failed:", err)

		// This stops the current cleanup run.
		return
	}

	// This logs the cleanup result.
	log.Printf("Health check retention cleanup completed: deleted %d checks older than %d days", deletedRows, c.RetentionDays)
}
