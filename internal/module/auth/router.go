package auth

import (
	"example/internal/repository/postgresql/user"
	"example/pkg/app"

	"github.com/go-chi/chi"
)

func SetupAuthRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	userRepository := user.NewUserRepository(app)

	authService := NewAuthService(app, userRepository)
	authHandler := NewAuthHandler(app, authService)

	r.Post("/register", authHandler.Register)
	r.Post("/session", authHandler.Login)
	r.Post("/login", authHandler.HandleLogin) // OAuth2 login
	r.Get("/callback", authHandler.HandleCallback)
	r.Get("/verify", authHandler.VerifyEmail)
}
