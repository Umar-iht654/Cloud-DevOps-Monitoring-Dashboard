package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertHandler stores the dependencies needed by the alert routes.
type AlertHandler struct {
	// DB gives the alert handler access to PostgreSQL through GORM.
	DB *gorm.DB
}

// NewAlertHandler creates a new AlertHandler with database access.
func NewAlertHandler(db *gorm.DB) *AlertHandler {
	// This returns a pointer to an AlertHandler so routes can use its methods.
	return &AlertHandler{
		DB: db,
	}
}

// getAlertLimit reads and validates the optional limit query parameter.
func getAlertLimit(c *gin.Context) int {
	// This sets a safe default number of alerts to return.
	limit := 50

	// This reads the optional limit query parameter.
	limitQuery := c.Query("limit")

	// This checks whether the limit query parameter was provided.
	if limitQuery == "" {
		// This returns the default limit.
		return limit
	}

	// This tries to convert the limit query parameter into a number.
	parsedLimit, err := strconv.Atoi(limitQuery)

	// This checks whether the limit was invalid.
	if err != nil || parsedLimit <= 0 {
		// This returns the default limit when the provided value is not usable.
		return limit
	}

	// This prevents users from requesting too many alerts at once.
	if parsedLimit > 200 {
		// This caps the maximum result size.
		return 200
	}

	// This returns the validated limit.
	return parsedLimit
}

// GetAlerts returns recent alerts owned by the logged-in user.
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	// This gets the authenticated user's ID from the request context.
	userID, ok := getUserIDFromContext(c)

	// This checks whether the user ID was missing or invalid.
	if !ok {
		// This returns a 401 response because the user is not authenticated.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authenticated",
		})

		// This stops the handler because alerts require a logged-in user.
		return
	}

	// This gets the validated alert limit.
	limit := getAlertLimit(c)

	// This creates a slice to store the alerts loaded from the database.
	var alerts []models.Alert

	// This loads the user's alerts, newest first, including the related service.
	if err := h.DB.Preload("Service").Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&alerts).Error; err != nil {
		// This returns a 500 response if alerts could not be loaded.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get alerts",
		})

		// This stops the handler because the database query failed.
		return
	}

	// This returns the alerts to the frontend.
	c.JSON(http.StatusOK, gin.H{
		"alerts":         alerts,
		"returned_count": len(alerts),
	})
}

// GetServiceAlerts returns recent alerts for one service owned by the logged-in user.
func (h *AlertHandler) GetServiceAlerts(c *gin.Context) {
	// This gets the authenticated user's ID from the request context.
	userID, ok := getUserIDFromContext(c)

	// This checks whether the user ID was missing or invalid.
	if !ok {
		// This returns a 401 response because the user is not authenticated.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authenticated",
		})

		// This stops the handler because alerts require a logged-in user.
		return
	}

	// This converts the service ID from the URL into a number.
	serviceID, err := strconv.Atoi(c.Param("id"))

	// This checks whether the service ID was invalid.
	if err != nil {
		// This returns a 400 response because the ID was not a number.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid service ID",
		})

		// This stops the handler because the service ID cannot be used.
		return
	}

	// This creates a variable to store the service if it exists.
	var service models.Service

	// This checks that the requested service belongs to the logged-in user.
	if err := h.DB.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		// This checks whether the service genuinely does not exist or does not belong to this user.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// This returns a 404 response because the service was not found for this user.
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Service not found",
			})

			// This stops the handler because users cannot view alerts for services they do not own.
			return
		}

		// This returns a 500 response because an unexpected database error happened.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check service ownership",
		})

		// This stops the handler because the ownership check failed unexpectedly.
		return
	}

	// This gets the validated alert limit.
	limit := getAlertLimit(c)

	// This creates a slice to store the alerts loaded from the database.
	var alerts []models.Alert

	// This loads alerts for the requested service, newest first.
	if err := h.DB.Preload("Service").Where("user_id = ? AND service_id = ?", userID, serviceID).Order("created_at DESC").Limit(limit).Find(&alerts).Error; err != nil {
		// This returns a 500 response if alerts could not be loaded.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get service alerts",
		})

		// This stops the handler because the database query failed.
		return
	}

	// This returns the service alerts to the frontend.
	c.JSON(http.StatusOK, gin.H{
		"service_id":     serviceID,
		"alerts":         alerts,
		"returned_count": len(alerts),
	})
}
