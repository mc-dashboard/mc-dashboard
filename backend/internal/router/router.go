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
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware for frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/login", auth.HandleLogin)
	r.Get("/auth/google/callback", auth.HandleCallback)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/logout", auth.HandleLogout)
		r.Get("/api/user", auth.HandleGetUser)

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{},
		}))

		r.Handle("/graphql", srv)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Helloo World!"))
	})

	return r
}
