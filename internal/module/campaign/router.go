package campaign

import (
	"example/internal/middleware"
	"example/internal/repository"
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCampaignRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := service.AuthNewService()

	userRepository := repository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)

	campaignRepository := repository.CampaignNewRepository(db)
	campaignService := service.CampaignNewService(campaignRepository)
	campaignHandler := NewCampaignHandler(campaignService)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadImage)
}
