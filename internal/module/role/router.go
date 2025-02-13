package role

import (
	"example/internal/repository/postgresql/role"
	"example/internal/repository/postgresql/token"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupRoleRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	tokenRepository := token.NewTokenRepository(app)

	roleRepository := role.NewRoleRepository(app)
	roleService := NewRoleService(app, roleRepository, tokenRepository)
	roleHandler := NewRoleHandler(app, roleService)

	r.With(middleware.AuthMiddleware).Route("/roles", func(r chi.Router) {
		r.Get("/", roleHandler.ListPaginate)
		r.Get("/{id:[a-fA-F0-9-]{36}}", roleHandler.Detail)
	})
}
