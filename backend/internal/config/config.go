package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Stores the port number the backend server will run on.
	Port string

	// Stores the PostgreSQL connection string.
	DatabaseURL string

	// Stores the secret key used to sign and verify JWT tokens.
	JWTSecret string
}

func LoadConfig() Config {
	// Tries to load environment variables from the .env file.
	err := godotenv.Load()

	if err != nil {
		// Prints this message if .env is missing, but still allows system environment variables to work.
		log.Println("No .env file found, using system environment variables")
	}

	// Reads the PORT value from the environment.
	port := os.Getenv("PORT")

	if port == "" {
		// Uses 8080 as a fallback if PORT is not set.
		port = "8080"
	}

	// Reads the database connection string from the environment.
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		// Stops the app because the backend cannot work without the database connection string.
		log.Fatal("DATABASE_URL is not set")
	}

	// Reads the JWT secret from the environment.
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		// Stops the app because login tokens cannot be created safely without a secret.
		log.Fatal("JWT_SECRET is not set")
	}

	return Config{
		// Adds the server port to the returned config.
		Port: port,

		// Adds the database URL to the returned config.
		DatabaseURL: databaseURL,

		// Adds the JWT secret to the returned config.
		JWTSecret: jwtSecret,
	}
}
