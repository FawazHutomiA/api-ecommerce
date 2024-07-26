package user

import (
	"example/internal/middleware"
	"example/internal/service"

	userRepository "example/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := service.AuthNewService()

	userRepository := userRepository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)
	userHandler := NewUserHandler(userService, authService)

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/user", middleware.AuthMiddleware(authService, userService), userHandler.GetUserByJWT)
	api.GET("/users", middleware.AuthMiddleware(authService, userService), userHandler.GetUsers)

	api.POST("/login", userHandler.HandleLogin)
	api.GET("/callback", userHandler.HandleCallback)
}
