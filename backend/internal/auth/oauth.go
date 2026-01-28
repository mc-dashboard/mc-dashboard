package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
)

var GoogleOAuthConfig *oauth2.Config

func InitOAuth(cfg *config.Config) {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}