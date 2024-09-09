package user

import (
	"example/internal/service"
	"example/pkg/middleware"

	userRepository "example/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories and services
	userRepository := userRepository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)
	userHandler := NewUserHandler(userService)

	// Authenticated routes (require JWT authentication)
	authUserRoutes := api.Group("/users")
	authUserRoutes.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		authUserRoutes.POST("/avatars", userHandler.UploadAvatar) // Upload user avatar
		authUserRoutes.GET("/me", userHandler.GetUserByJWT)       // Get user by JWT token
		authUserRoutes.GET("", userHandler.GetUsers)              // Get list of users
	}
}
