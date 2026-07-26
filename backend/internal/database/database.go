package database

import (
	"log"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
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

		// This creates or updates the services table.
		&models.Service{},

		// This creates or updates the health_checks table.
		&models.HealthCheck{},

		// This creates or updates the alerts table.
		&models.Alert{},
	)

	if err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	log.Println("Database migrations completed")
}
