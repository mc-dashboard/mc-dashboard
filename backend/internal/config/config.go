package config

import (
	"os"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	SessionSecret      string
	AllowedEmails      map[string]bool
	ServerPort         string
	FrontendURL        string
}

func Load() *Config {
	allowedEmailsMap := make(map[string]bool)
	allowedEmailsMap["rsuri@irusmail.com"] = true; // eventually get this from the DB 

	return &Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		AllowedEmails:      allowedEmailsMap,
		ServerPort:         getEnv("PORT", "8080"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}