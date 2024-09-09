package auth

import (
	"example/internal/service"

	userRepository "example/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories and services
	userRepository := userRepository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)

	authService := service.AuthNewService(userRepository)
	authHandler := NewAuthHandler(authService, userService)

	api.POST("/register", authHandler.RegisterUser)                 // Register user
	api.POST("/sessions", authHandler.Login)                        // User login
	api.POST("/email_checkers", authHandler.CheckEmailAvailability) // Check email availability
	api.POST("/login", authHandler.HandleLogin)                     // OAuth2 login
	api.GET("/callback", authHandler.HandleCallback)
}
