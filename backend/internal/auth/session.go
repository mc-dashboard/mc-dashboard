package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
)

var Store *sessions.CookieStore

func InitSessionStore(cfg *config.Config) {
	Store = sessions.NewCookieStore([]byte(cfg.SessionSecret))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return Store.Get(r, "auth-session")
}
