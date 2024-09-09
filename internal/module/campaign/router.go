package campaign

import (
	"example/internal/middleware"
	"example/internal/repository"
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCampaignRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize services and repositories
	campaignRepository := repository.CampaignNewRepository(db)
	campaignService := service.CampaignNewService(campaignRepository)
	campaignHandler := NewCampaignHandler(campaignService)

	// Public routes (accessible without authentication)
	publicCampaignRoutes := api.Group("/campaigns")
	{
		publicCampaignRoutes.GET("", campaignHandler.GetCampaigns)    // List all campaigns
		publicCampaignRoutes.GET("/:id", campaignHandler.GetCampaign) // Get a specific campaign by ID
	}

	// Authenticated routes (require JWT authentication)
	authCampaignRoutes := api.Group("/campaigns")
	authCampaignRoutes.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		authCampaignRoutes.POST("", campaignHandler.CreateCampaign)     // Create a new campaign
		authCampaignRoutes.PUT("/:id", campaignHandler.UpdateCampaign)  // Update an existing campaign
		authCampaignRoutes.POST("/images", campaignHandler.UploadImage) // Upload campaign image
	}
}
