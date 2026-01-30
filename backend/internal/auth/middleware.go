package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rohanvsuri/minecraft-dashboard/internal/db"
)

// RequireAuth is a middleware that checks for a valid, authorized session.
func (s *Service) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getSession(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":    "unauthorized, no session found",
				"redirect": "/login",
			})
			return
		}

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

		allowed, err := db.UserExistsByEmail(s.db, email)
		if err != nil || !allowed {
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
