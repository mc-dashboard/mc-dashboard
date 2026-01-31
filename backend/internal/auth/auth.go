package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
	"github.com/rohanvsuri/minecraft-dashboard/internal/db"
)

// Service holds all auth-related dependencies.
type Service struct {
	oauthConfig  *oauth2.Config
	sessionStore *sessions.CookieStore
	db           *pgxpool.Pool
	frontendURL  string
}

// NewService creates the auth service with all dependencies.
func NewService(cfg *config.Config, pool *pgxpool.Pool) *Service {
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	// In production with HTTPS, cookies must have Secure flag
	// SameSite=None is required for cross-site cookies with Secure flag
	sameSite := http.SameSiteLaxMode
	if cfg.IsProduction {
		sameSite = http.SameSiteNoneMode
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   cfg.IsProduction,
		SameSite: sameSite,
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Service{
		oauthConfig:  oauthConfig,
		sessionStore: store,
		db:           pool,
		frontendURL:  cfg.FrontendURL,
	}
}

func (s *Service) getSession(r *http.Request) (*sessions.Session, error) {
	return s.sessionStore.Get(r, "auth-session")
}

func (s *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()

	session, _ := s.getSession(r)
	session.Values["state"] = state
	session.Save(r, w)

	url := s.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Service) HandleCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := s.getSession(r)

	savedState, _ := session.Values["state"].(string)
	if r.URL.Query().Get("state") != savedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	allowed, err := db.UserExistsByEmail(s.db, userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to check authorization", http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "You are not authorized to access this application", http.StatusForbidden)
		return
	}

	session.Values["email"] = userInfo.Email
	session.Values["name"] = userInfo.Name
	session.Save(r, w)

	http.Redirect(w, r, s.frontendURL, http.StatusSeeOther)
}

func (s *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.getSession(r)
	session.Values["email"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out",
	})
}

func (s *Service) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	session, _ := s.getSession(r)

	email, _ := session.Values["email"].(string)
	name, _ := session.Values["name"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email": email,
		"name":  name,
	})
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
