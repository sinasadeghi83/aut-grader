package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	ServerPort string
	DbURL      string
	SecretKey  string
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() *AppConfig {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set. Please provide it.")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set. Please provide it.")
	}

	return &AppConfig{
		ServerPort: ":" + port,
		DbURL:      dbURL,
		SecretKey:  secretKey,
	}
}
