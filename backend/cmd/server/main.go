package main

import (
	"net/http"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

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

	router.Run(":" + cfg.Port)
}
