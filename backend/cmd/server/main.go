package main

import (
	"net/http"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/config"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/database"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/handlers"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/middleware"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/worker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// This loads environment variables such as PORT, DATABASE_URL and JWT_SECRET.
	cfg := config.LoadConfig()

	// This connects the backend to the PostgreSQL database.
	db := database.Connect(cfg.DatabaseURL)

	// This creates or updates the database tables using the GORM models.
	database.Migrate(db)

	// This creates the authentication handler and gives it database access plus the JWT secret.
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	// This creates the service handler and gives it database access.
	serviceHandler := handlers.NewServiceHandler(db)

	// This creates the health check handler and gives it database access.
	healthCheckHandler := handlers.NewHealthCheckHandler(db)

	// This creates the dashboard handler and gives it database access.
	dashboardHandler := handlers.NewDashboardHandler(db)

	// This creates the alert handler and gives it database access.
	alertHandler := handlers.NewAlertHandler(db)

	// This creates the background health checker and gives it database access.
	healthChecker := worker.NewHealthChecker(db)

	// This starts the health checker in a goroutine so it runs in the background while the API server runs.
	go healthChecker.Start()

	// This creates the health-check retention cleaner and gives it database access plus cleanup settings.
	retentionCleaner := worker.NewHealthCheckRetentionCleaner(
		db,
		cfg.HealthCheckRetentionDays,
		cfg.HealthCheckCleanupIntervalHours,
	)

	// This starts the retention cleaner in a goroutine so old raw health checks are cleaned up in the background.
	go retentionCleaner.Start()

	// This creates a new Gin router with default logging and recovery middleware.
	router := gin.Default()

	// This adds CORS middleware so approved frontend origins can call the backend from the browser.
	router.Use(cors.New(cors.Config{
		// This allows requests from the frontend URLs defined in the environment config.
		AllowOrigins: cfg.FrontendURLs,

		// This allows the frontend to use the HTTP methods needed by the API.
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// This allows the frontend to send JSON and JWT Authorization headers.
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},

		// This allows the browser to read the Content-Length response header if needed.
		ExposeHeaders: []string{"Content-Length"},

		// This disables cookie-based cross-origin credentials because we are using Bearer tokens.
		AllowCredentials: false,

		// This tells the browser how long it can cache the CORS preflight response.
		MaxAge: 12 * time.Hour,
	}))

	// This records Prometheus metrics for every backend request.
	router.Use(middleware.MetricsMiddleware())

	// This creates a public health check route for checking if the backend is running.
	router.GET("/health", func(c *gin.Context) {
		// This returns a successful JSON response for the health check.
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "cloud-devops-monitoring-backend",
		})
	})

	// This creates a public API status route.
	router.GET("/api/status", func(c *gin.Context) {
		// This returns a simple JSON message confirming the API is running.
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend API is running",
		})
	})

	// This exposes backend metrics in a format Prometheus can scrape.
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// This creates a public route to test whether the database connection is healthy.
	router.GET("/api/db-status", func(c *gin.Context) {
		// This gets the underlying SQL database connection from GORM.
		sqlDB, err := db.DB()

		// This checks whether GORM failed to provide the SQL database connection.
		if err != nil {
			// This returns a 500 response if the database connection cannot be accessed.
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to access database connection",
			})

			// This stops the route handler after the error response.
			return
		}

		// This sends a ping to PostgreSQL to check if the database is reachable.
		err = sqlDB.Ping()

		// This checks whether the database ping failed.
		if err != nil {
			// This returns a 500 response if PostgreSQL is not reachable.
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Database is not reachable",
			})

			// This stops the route handler after the error response.
			return
		}

		// This returns a successful response if the database connection is healthy.
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Database connection is healthy",
		})
	})

	// This creates a route group so all API routes start with /api.
	api := router.Group("/api")

	// This creates an authentication route group under /api/auth.
	auth := api.Group("/auth")

	// This registers the public user registration route.
	auth.POST("/register", authHandler.Register)

	// This registers the public user login route.
	auth.POST("/login", authHandler.Login)

	// This creates a separate auth group for routes that require a valid JWT token.
	protectedAuth := auth.Group("/")

	// This applies the JWT authentication middleware to the protected auth routes.
	protectedAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// This registers the protected route that returns the currently logged-in user.
	protectedAuth.GET("/me", authHandler.Me)

	// This creates a services route group under /api/services.
	services := api.Group("/services")

	// This applies JWT authentication middleware to all service routes.
	services.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// This registers the protected route for creating a monitored service.
	services.POST("", serviceHandler.CreateService)

	// This registers the protected route for getting all monitored services owned by the logged-in user.
	services.GET("", serviceHandler.GetServices)

	// This registers the protected route for getting one monitored service by ID.
	services.GET("/:id", serviceHandler.GetService)

	// This registers the protected route for updating one monitored service by ID.
	services.PUT("/:id", serviceHandler.UpdateService)

	// This registers the protected route for deleting one monitored service by ID.
	services.DELETE("/:id", serviceHandler.DeleteService)

	// This registers the protected route for getting recent health checks for one service.
	services.GET("/:id/health-checks", healthCheckHandler.GetHealthChecks)

	// This registers the protected route for getting calculated summary stats for one service.
	services.GET("/:id/summary", healthCheckHandler.GetServiceSummary)

	// This registers the protected route for getting alerts for one service.
	services.GET("/:id/alerts", alertHandler.GetServiceAlerts)

	// This creates a dashboard route group under /api/dashboard.
	dashboard := api.Group("/dashboard")

	// This applies JWT authentication middleware to all dashboard routes.
	dashboard.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// This registers the protected route for getting dashboard summary stats.
	dashboard.GET("/summary", dashboardHandler.GetSummary)

	// This creates an alerts route group under /api/alerts.
	alerts := api.Group("/alerts")

	// This applies JWT authentication middleware to all alert routes.
	alerts.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// This registers the protected route for getting recent alerts owned by the logged-in user.
	alerts.GET("", alertHandler.GetAlerts)

	// This starts the backend server using the configured port.
	router.Run(":" + cfg.Port)
}
