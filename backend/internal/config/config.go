package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort         string
	// Environment string // "development" | "production"
	DatabaseURL        string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	
	SessionSecret      string
	SessionMaxAge      int
	
	FrontendURL        string
	EnableDebugLogging bool
}

func Load() (*Config, error) {

	cfg := &Config{
		ServerPort:         getEnv("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		SessionMaxAge:      86400,
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		EnableDebugLogging: os.Getenv("ENABLE_DEBUG_LOGGING") == "true",
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if cfg.GoogleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
	}

	return cfg, nil

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
