package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/rohanvsuri/minecraft-dashboard/internal/auth"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
	"github.com/rohanvsuri/minecraft-dashboard/internal/db"
	"github.com/rohanvsuri/minecraft-dashboard/internal/graph"
	"github.com/rohanvsuri/minecraft-dashboard/internal/router"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	authService := auth.NewService(cfg, pool)

	resolver := &graph.Resolver{DB: pool}

	r := router.NewRouter(cfg, authService, resolver)

	log.Printf("Server running on http://localhost:%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
