package campaign

import (
	"example/auth"
	"example/module/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCampaignRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := auth.NewService()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	campaignRepository := NewRepository(db)
	campaignService := NewService(campaignRepository)
	campaignHandler := NewCampaignHandler(campaignService)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", user.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", user.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", user.AuthMiddleware(authService, userService), campaignHandler.UploadImage)
}
