package user

import (
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupUserRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	userRepository := user.NewUserRepository(app)
	userService := NewUserService(app, userRepository)
	userHandler := NewUserHandler(app, userService)

	r.With(middleware.AuthMiddleware).Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.ListPaginate)
		r.Get("/me", userHandler.DetailByJwt)
	})
}
