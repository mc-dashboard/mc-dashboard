package auth

import (
	"fmt"
	"context"
	"encoding/json"
	"net/http"
	"crypto/rand"
	"encoding/base64"

	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
)

var cfg *config.Config

func InitHandlers(c *config.Config) {
	cfg = c
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate random state for CSRF protection
	state := GenerateRandomState()

	session, _ := GetSession(r)
	fmt.Println(session, state)
	session.Values["state"] = state
	session.Save(r, w)

	url := GoogleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := GetSession(r)

	// Verify state to prevent CSRF
	savedState, _ := session.Values["state"].(string)
	if r.URL.Query().Get("state") != savedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user info
	client := GoogleOAuthConfig.Client(context.Background(), token)
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

	// Check if user is allowed
	if !cfg.AllowedEmails[userInfo.Email] {
		http.Error(w, "You are not authorized to access this application", http.StatusForbidden)
		return
	}

	// Save user info in session
	session.Values["email"] = userInfo.Email
	session.Values["name"] = userInfo.Name
	// delete(session.Values, "state") // Clean up state
	session.Save(r, w)
	fmt.Println(session, "at the end of callback");

	// Redirect to frontend
	http.Redirect(w, r, cfg.FrontendURL, http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := GetSession(r)
	session.Values["email"] = ""
	session.Options.MaxAge = -1 // Delete cookie
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out",
	})
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	session, _ := GetSession(r)

	email, _ := session.Values["email"].(string)
	name, _ := session.Values["name"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email": email,
		"name":  name,
	})
}

func GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}