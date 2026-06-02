package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsRegion          string

	SessionSecret string //nolint:gosec // field name, not an actual secret value
	SessionMaxAge int

	FrontendURL        string
	EnableDebugLogging bool
	IsProduction       bool

	// RCON connection for Minecraft server monitoring
	MinecraftHost         string
	MinecraftRCONPort     string
	MinecraftRCONPassword string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:            getEnv("PORT", "8080"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		GoogleClientID:        os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		AwsAccessKeyId:        os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsRegion:             getEnv("AWS_REGION", "us-east-1"),
		SessionSecret:         os.Getenv("SESSION_SECRET"),
		SessionMaxAge:         86400,
		FrontendURL:           getEnv("FRONTEND_URL", "http://localhost:3000"),
		EnableDebugLogging:    os.Getenv("ENABLE_DEBUG_LOGGING") == "true",
		IsProduction:          os.Getenv("ENV") == "production",
		MinecraftHost:         os.Getenv("MINECRAFT_HOST"),
		MinecraftRCONPort:     getEnv("MINECRAFT_RCON_PORT", "25575"),
		MinecraftRCONPassword: os.Getenv("MINECRAFT_RCON_PASSWORD"),
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
	if cfg.AwsAccessKeyId == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID is required")
	}
	if cfg.AwsSecretAccessKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY is required")
	}
	if cfg.AwsRegion == "" {
		return nil, fmt.Errorf("AWS_REGION is required")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
	}
	if cfg.MinecraftHost == "" {
		return nil, fmt.Errorf("MINECRAFT_HOST is required")
	}
	if cfg.MinecraftRCONPassword == "" {
		return nil, fmt.Errorf("MINECRAFT_RCON_PASSWORD is required")
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
