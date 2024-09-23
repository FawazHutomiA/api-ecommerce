package router

import (
	"example/internal/module/user"
	"example/pkg/app"

	"github.com/go-chi/chi"
)

func SetupRoutes(r *chi.Mux, app app.AppConfig) {
	// API V1
	r.Route("/api/v1", func(r chi.Router) {
		user.SetupUserRoutes(r, app)
	})
}
