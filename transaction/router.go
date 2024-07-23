package transaction

import (
	"example/auth"
	"example/campaign"
	"example/payment"
	"example/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := auth.NewService()

	paymentService := payment.NewService()

	campaignRepository := campaign.NewRepository(db)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	transactionRepository := NewRepository(db)
	transactionService := NewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := NewTransactionHandler(transactionService)

	api.GET("/campaigns/:id/transactions", user.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", user.AuthMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", user.AuthMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
}
