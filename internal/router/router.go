package router

import (
	"example/internal/module/auth"
	"example/internal/module/campaign"
	"example/internal/module/transaction"
	"example/internal/module/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(db *gorm.DB) {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(cors.Default())
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	// Initialize each module's routes
	auth.SetupAuthRoutes(api, db)
	user.SetupUserRoutes(api, db)
	campaign.SetupCampaignRoutes(api, db)
	transaction.SetupTransactionRoutes(api, db)

	router.Run(":8080")
}
