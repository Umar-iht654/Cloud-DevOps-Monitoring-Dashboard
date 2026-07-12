package main

import (
	"net/http"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/config"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.Connect(cfg.DatabaseURL)
	database.Migrate(db)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "cloud-devops-monitoring-backend",
		})
	})

	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend API is running",
		})
	})

	router.GET("/api/db-status", func(c *gin.Context) {
		sqlDB, err := db.DB()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to access database connection",
			})
			return
		}

		err = sqlDB.Ping()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Database is not reachable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Database connection is healthy",
		})
	})

	router.Run(":" + cfg.Port)
}
