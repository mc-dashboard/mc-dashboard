package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/rohanvsuri/minecraft-dashboard/internal/auth"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
	"github.com/rohanvsuri/minecraft-dashboard/internal/db"
	"github.com/rohanvsuri/minecraft-dashboard/internal/graph"
	"github.com/rohanvsuri/minecraft-dashboard/internal/lambda"
	"github.com/rohanvsuri/minecraft-dashboard/internal/minecraft"
	"github.com/rohanvsuri/minecraft-dashboard/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	authService := auth.NewService(cfg, pool)

	lambdaService := lambda.NewLambdaService(cfg)

	rconClient := minecraft.NewRCONClient(
		cfg.MinecraftHost,
		cfg.MinecraftRCONPort,
		cfg.MinecraftRCONPassword,
	)

	if err := rconClient.Connect(); err != nil {
		log.Printf("Failed to connect to Minecraft RCON: %v. RCON features will be unavailable.", err)
	} else {
		defer func() {
			if err := rconClient.Disconnect(); err != nil {
				log.Printf("Failed to disconnect RCON: %v", err)
			}
		}()
		log.Println("Minecraft RCON connected successfully")
	}

	minecraftHandler := minecraft.NewMinecraftHandler(lambdaService, rconClient)

	resolver := &graph.Resolver{DB: pool}

	r := router.NewRouter(cfg, authService, resolver, minecraftHandler)

	log.Printf("Server running on http://localhost:%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}
