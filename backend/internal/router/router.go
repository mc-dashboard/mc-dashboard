package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/rohanvsuri/minecraft-dashboard/internal/auth"
	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
	"github.com/rohanvsuri/minecraft-dashboard/internal/graph"
	"github.com/rohanvsuri/minecraft-dashboard/internal/minecraft"
)

func NewRouter(cfg *config.Config, authService *auth.Service, resolver *graph.Resolver, minecraftHandler *minecraft.MinecraftHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/login", authService.HandleLogin)
	r.Get("/auth/google/callback", authService.HandleCallback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authService.RequireAuth)

		r.Get("/logout", authService.HandleLogout)
		r.Get("/api/user", authService.HandleGetUser)

		r.Post("/api/minecraft/start", minecraftHandler.StartServer)
		r.Post("/api/minecraft/stop", minecraftHandler.StopServer)
		r.Get("/api/minecraft/status", minecraftHandler.GetServerStatus)
		r.Get("/api/minecraft/players", minecraftHandler.GetOnlinePlayers)
		r.Post("/api/minecraft/command", minecraftHandler.ExecuteCommand)

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
			Resolvers: resolver,
		}))
		r.Handle("/graphql", srv)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	return r
}
