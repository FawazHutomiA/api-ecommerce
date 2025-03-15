package user

import (
	"example/internal/repository/postgresql/role"
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupUserRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	tokenRepository := token.NewTokenRepository(app)
	roleRepo := role.NewRoleRepository(app)

	userRepository := user.NewUserRepository(app)
	userService := NewUserService(app, userRepository, tokenRepository, roleRepo)
	userHandler := NewUserHandler(app, userService)

	// Middleware Auth + SuperAdminMiddleware hanya untuk route tertentu
	r.With(middleware.AuthMiddleware).Route("/users", func(r chi.Router) {
		r.Get("/me", userHandler.UserMe) // Ini bisa diakses oleh semua user yang login

		// Route khusus Super Admin
		r.With(middleware.SuperAdminMiddleware).Group(func(r chi.Router) {
			r.Get("/", userHandler.ListPaginate)
			r.Get("/{id:[a-fA-F0-9-]{36}}", userHandler.Detail)
			r.Post("/", userHandler.Create)
			r.Put("/{id:[a-fA-F0-9-]{36}}", userHandler.Update)
			r.Delete("/{id:[a-fA-F0-9-]{36}}", userHandler.Delete)
		})
	})
}
