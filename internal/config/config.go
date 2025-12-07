package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		URL string
	}
	JWT struct {
		Secret     string
		Expiration string
	}
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &Config{}

	// Load Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Load Database configuration
	cfg.Database.URL = getEnv("DATABASE_URL", "")
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Load JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "")
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	cfg.JWT.Expiration = getEnv("JWT_EXPIRATION", "24h")

	return cfg
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
