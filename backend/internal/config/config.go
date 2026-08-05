package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config stores all environment variables the backend needs.
type Config struct {
	// Port stores the backend server port.
	Port string

	// DatabaseURL stores the PostgreSQL connection string.
	DatabaseURL string

	// JWTSecret stores the secret key used to sign JWT tokens.
	JWTSecret string

	// FrontendURLs stores the frontend origins that are allowed to call the backend.
	FrontendURLs []string

	// HealthCheckRetentionDays stores how many days of raw health checks should be kept.
	HealthCheckRetentionDays int

	// HealthCheckCleanupIntervalHours stores how often the cleanup worker should run.
	HealthCheckCleanupIntervalHours int

	// AppBaseURL stores the frontend URL used for verification links.
	AppBaseURL string

	// SMTPHost stores the SMTP server hostname.
	SMTPHost string

	// SMTPPort stores the SMTP server port.
	SMTPPort string

	// SMTPUsername stores the SMTP username.
	SMTPUsername string

	// SMTPPassword stores the SMTP password or API key.
	SMTPPassword string

	// SMTPFrom stores the email address used as the sender.
	SMTPFrom string
}

// LoadConfig loads environment variables from .env or the system environment.
func LoadConfig() Config {
	// This tries to load environment variables from a local .env file.
	err := godotenv.Load()

	// This allows the app to continue if no .env file exists, such as in some deployment environments.
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// This reads the PORT value from the environment.
	port := os.Getenv("PORT")

	// This uses 8080 as the default backend port if PORT is missing.
	if port == "" {
		port = "8080"
	}

	// This reads the DATABASE_URL value from the environment.
	databaseURL := os.Getenv("DATABASE_URL")

	// This stops the app if DATABASE_URL is missing because the backend needs PostgreSQL.
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// This reads the JWT_SECRET value from the environment.
	jwtSecret := os.Getenv("JWT_SECRET")

	// This stops the app if JWT_SECRET is missing because authentication would not be secure.
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// This reads the FRONTEND_URLS value from the environment.
	frontendURLsRaw := os.Getenv("FRONTEND_URLS")

	// This provides common local frontend URLs if FRONTEND_URLS is missing.
	if frontendURLsRaw == "" {
		frontendURLsRaw = "http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000,http://127.0.0.1:3000,http://localhost:4200,http://127.0.0.1:4200"
	}

	// This splits the comma-separated FRONTEND_URLS string into a slice.
	frontendURLs := strings.Split(frontendURLsRaw, ",")

	// This creates a clean slice for storing trimmed frontend URLs.
	cleanFrontendURLs := []string{}

	// This loops through each frontend URL from the environment variable.
	for _, frontendURL := range frontendURLs {
		// This removes extra spaces around each frontend URL.
		trimmedURL := strings.TrimSpace(frontendURL)

		// This checks that the frontend URL is not empty.
		if trimmedURL != "" {
			// This adds the cleaned frontend URL to the allowed list.
			cleanFrontendURLs = append(cleanFrontendURLs, trimmedURL)
		}
	}

	// This sets the default number of days to keep raw health checks.
	healthCheckRetentionDays := 14

	// This reads the optional HEALTH_CHECK_RETENTION_DAYS value from the environment.
	healthCheckRetentionDaysRaw := os.Getenv("HEALTH_CHECK_RETENTION_DAYS")

	// This checks whether the retention value was provided.
	if healthCheckRetentionDaysRaw != "" {
		// This converts the retention value from text into a number.
		parsedRetentionDays, err := strconv.Atoi(healthCheckRetentionDaysRaw)

		// This stops the app if the retention value is invalid.
		if err != nil || parsedRetentionDays < 0 {
			log.Fatal("HEALTH_CHECK_RETENTION_DAYS must be 0 or greater")
		}

		// This stores the configured retention value.
		healthCheckRetentionDays = parsedRetentionDays
	}

	// This sets the default cleanup interval to once per day.
	healthCheckCleanupIntervalHours := 24

	// This reads the optional HEALTH_CHECK_CLEANUP_INTERVAL_HOURS value from the environment.
	healthCheckCleanupIntervalHoursRaw := os.Getenv("HEALTH_CHECK_CLEANUP_INTERVAL_HOURS")

	// This checks whether the cleanup interval value was provided.
	if healthCheckCleanupIntervalHoursRaw != "" {
		// This converts the cleanup interval from text into a number.
		parsedCleanupIntervalHours, err := strconv.Atoi(healthCheckCleanupIntervalHoursRaw)

		// This stops the app if the cleanup interval is invalid.
		if err != nil || parsedCleanupIntervalHours <= 0 {
			log.Fatal("HEALTH_CHECK_CLEANUP_INTERVAL_HOURS must be greater than 0")
		}

		// This stores the configured cleanup interval.
		healthCheckCleanupIntervalHours = parsedCleanupIntervalHours
	}

	// This reads the frontend base URL used in email verification links.
	appBaseURL := os.Getenv("APP_BASE_URL")

	// This defaults to the local Docker frontend URL if APP_BASE_URL is missing.
	if appBaseURL == "" {
		appBaseURL = "http://localhost:3000"
	}

	// This reads SMTP email settings from the environment.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	// This returns all loaded config values so other parts of the backend can use them.
	return Config{
		Port:                            port,
		DatabaseURL:                     databaseURL,
		JWTSecret:                       jwtSecret,
		FrontendURLs:                    cleanFrontendURLs,
		HealthCheckRetentionDays:        healthCheckRetentionDays,
		HealthCheckCleanupIntervalHours: healthCheckCleanupIntervalHours,
		AppBaseURL:                      appBaseURL,
		SMTPHost:                        smtpHost,
		SMTPPort:                        smtpPort,
		SMTPUsername:                    smtpUsername,
		SMTPPassword:                    smtpPassword,
		SMTPFrom:                        smtpFrom,
	}
}
