package user

import (
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupUserRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	tokenRepository := token.NewTokenRepository(app)

	userRepository := user.NewUserRepository(app)
	userService := NewUserService(app, userRepository, tokenRepository)
	userHandler := NewUserHandler(app, userService)

	r.With(middleware.AuthMiddleware).Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.ListPaginate)
		r.Get("/{id:[a-fA-F0-9-]{36}}", userHandler.Detail)
		r.Post("/", userHandler.Create)
		r.Put("/{id:[a-fA-F0-9-]{36}}", userHandler.Update)
		r.Delete("/{id:[a-fA-F0-9-]{36}}", userHandler.Delete)
		r.Get("/me", userHandler.UserMe)
	})
}
