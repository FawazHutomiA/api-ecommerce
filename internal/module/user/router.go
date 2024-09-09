package user

import (
	"example/internal/middleware"
	"example/internal/service"

	userRepository "example/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories and services
	userRepository := userRepository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)
	userHandler := NewUserHandler(userService)

	// Public routes (accessible without authentication)
	publicUserRoutes := api.Group("/users")
	{
		publicUserRoutes.POST("", userHandler.RegisterUser)                          // Register user
		publicUserRoutes.POST("/sessions", userHandler.Login)                        // User login
		publicUserRoutes.POST("/email_checkers", userHandler.CheckEmailAvailability) // Check email availability
		publicUserRoutes.POST("/login", userHandler.HandleLogin)                     // OAuth2 login
		publicUserRoutes.GET("/callback", userHandler.HandleCallback)                // OAuth2 callback
	}

	// Authenticated routes (require JWT authentication)
	authUserRoutes := api.Group("/users")
	authUserRoutes.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		authUserRoutes.POST("/avatars", userHandler.UploadAvatar) // Upload user avatar
		authUserRoutes.GET("/user", userHandler.GetUserByJWT)     // Get user by JWT token
		authUserRoutes.GET("", userHandler.GetUsers)              // Get list of users
	}
}
