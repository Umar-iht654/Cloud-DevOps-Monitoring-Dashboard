package config

import (
	"log"
	"os"
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

	// This returns all loaded config values so other parts of the backend can use them.
	return Config{
		Port:         port,
		DatabaseURL:  databaseURL,
		JWTSecret:    jwtSecret,
		FrontendURLs: cleanFrontendURLs,
	}
}
