package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// NewRouter creates and returns a new HTTP router with all defined routes.
// It connects HTTP endpoints to their respective handler functions.
func NewRouter(h HandlerInterface) (http.Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("invalid handler")
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Get("/health", h.HealthHandler)

	r.Route("/v1/packs", func(r chi.Router) {
		r.Get("/", h.GetPacks)
		r.Post("/", h.PostPacks)
		r.Delete("/{size}", h.DeletePacks)
	})

	r.Post("/v1/order", h.CalculateOrder)

	r.NotFound(h.NotFoundHandler)

	return r, nil
}
