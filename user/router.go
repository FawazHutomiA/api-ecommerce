package user

import (
	"example/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := auth.NewService()

	userRepository := NewRepository(db)
	userService := NewService(userRepository)
	userHandler := NewUserHandler(userService, authService)

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/user", AuthMiddleware(authService, userService), userHandler.GetUserByJWT)
	api.GET("/users", AuthMiddleware(authService, userService), userHandler.GetUsers)

	api.POST("/login", userHandler.HandleLogin)
	api.GET("/callback", userHandler.HandleCallback)
}
