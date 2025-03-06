package router

import (
	"example/internal/module/auth"
	"example/internal/module/role"
	"example/internal/module/user"
	"example/pkg/app"

	"github.com/go-chi/chi"
)

func SetupRoutes(r *chi.Mux, app app.AppConfig) {
	// API V1
	r.Route("/api/v1", func(r chi.Router) {
		auth.SetupAuthRoutes(r, app)
		user.SetupUserRoutes(r, app)
		role.SetupRoleRoutes(r, app)
	})
}
