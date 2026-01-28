package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
)

var allowedEmails map[string]bool

func InitAllowedEmails(cfg *config.Config) {
	allowedEmails = cfg.AllowedEmails
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// check if session exists
		session, err := GetSession(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":    "unauthorized, no session found",
				"redirect": "/login",
			})
			return
		}

		// check if email exists
		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":    "unauthorized, no email found",
				"redirect": "/login",
			})
			return
		}

		// check if email is in whitelist
		if !allowedEmails[email] {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "forbidden, email not authorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
