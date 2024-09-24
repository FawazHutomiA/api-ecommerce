package auth

import (
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"

	"github.com/go-chi/chi"
)

func SetupAuthRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	userRepository := user.NewUserRepository(app)
	tokenRepository := token.NewTokenRepository(app)

	authService := NewAuthService(app, userRepository, tokenRepository)
	authHandler := NewAuthHandler(app, authService)

	r.Post("/register", authHandler.Register)
	r.Post("/session", authHandler.Login)
}
