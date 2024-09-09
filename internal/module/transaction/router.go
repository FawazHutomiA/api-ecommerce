package transaction

import (
	"example/internal/middleware"
	"example/internal/repository"
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize services and repositories
	paymentService := service.PaymentNewService()
	campaignRepository := repository.CampaignNewRepository(db)
	transactionRepository := repository.TransactionNewRepository(db)
	transactionService := service.TransactionNewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := NewTransactionHandler(transactionService)

	// Public routes (accessible without authentication)
	publicTransactionRoutes := api.Group("/transactions")
	{
		// Notification endpoint does not require authentication
		publicTransactionRoutes.POST("/notification", transactionHandler.GetNotification)
	}

	// Authenticated routes (require JWT authentication)
	authTransactionRoutes := api.Group("/transactions")
	authTransactionRoutes.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		authTransactionRoutes.GET("", transactionHandler.GetUserTransactions) // Get user transactions
		authTransactionRoutes.POST("", transactionHandler.CreateTransaction)  // Create a new transaction
	}

	// Authenticated routes for campaign-specific transactions
	authCampaignTransactionRoutes := api.Group("/campaigns/:id/transactions")
	authCampaignTransactionRoutes.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to these routes
	{
		authCampaignTransactionRoutes.GET("", transactionHandler.GetCampaignTransactions) // Get transactions for a specific campaign
	}
}
