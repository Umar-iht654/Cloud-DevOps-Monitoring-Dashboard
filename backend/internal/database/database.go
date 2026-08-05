package database

import (
	"log"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/monitoring"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) *gorm.DB {
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL database")

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		// This creates or updates the users table.
		&models.User{},

		// This creates or updates the pending registrations table.
		&models.PendingRegistration{},

		// This creates or updates the services table.
		&models.Service{},

		// This creates or updates the health_checks table.
		&models.HealthCheck{},

		// This creates or updates the alerts table.
		&models.Alert{},

		// This creates or updates the hourly service summaries table.
		&models.HourlyServiceSummary{},

		// This creates or updates the daily service summaries table.
		&models.DailyServiceSummary{},
	)

	if err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// This updates any existing services that were created before the new minimum interval rule.
	if err := db.Model(&models.Service{}).
		Where("check_interval_seconds < ?", monitoring.MinCheckIntervalSeconds).
		Update("check_interval_seconds", monitoring.MinCheckIntervalSeconds).Error; err != nil {
		// This stops startup if old service intervals could not be normalised.
		log.Fatal("Failed to normalise service check intervals:", err)
	}

	log.Println("Database migrations completed")
}
