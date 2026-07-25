package middleware

import (
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records Prometheus metrics for every HTTP request.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This stores the time when the request first reached the backend.
		startTime := time.Now()

		// This lets Gin continue to the actual route handler.
		c.Next()

		// This gets the route pattern, such as /api/services/:id instead of /api/services/12.
		path := c.FullPath()

		// This uses the raw URL path if Gin cannot find a route pattern.
		if path == "" {
			path = c.Request.URL.Path
		}

		// This calculates how long the request took.
		duration := time.Since(startTime)

		// This records the request method, route path, status code and duration.
		metrics.RecordHTTPRequest(c.Request.Method, path, c.Writer.Status(), duration)
	}
}
