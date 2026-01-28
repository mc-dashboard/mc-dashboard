package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/rohanvsuri/minecraft-dashboard/internal/auth"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
	"github.com/rohanvsuri/minecraft-dashboard/internal/router"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := config.Load()

	auth.InitHandlers(cfg)
	auth.InitAllowedEmails(cfg)
	auth.InitSessionStore(cfg);
	auth.InitOAuth(cfg);

	r := router.NewRouter(cfg);

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}