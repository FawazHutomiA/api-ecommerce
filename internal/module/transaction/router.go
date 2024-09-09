package transaction

import (
	"example/internal/repository"
	"example/internal/service"
	"example/pkg/middleware"

	userRepository "example/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize services and repositories
	paymentService := service.PaymentNewService()

	campaignRepository := repository.CampaignNewRepository(db)

	userRepository := userRepository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)

	transactionRepository := repository.TransactionNewRepository(db)
	transactionService := service.TransactionNewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := NewTransactionHandler(transactionService, userService)

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
