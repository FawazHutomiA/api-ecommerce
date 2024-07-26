package transaction

import (
	"example/internal/middleware"
	"example/internal/repository"
	"example/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	authService := service.AuthNewService()

	paymentService := service.PaymentNewService()

	campaignRepository := repository.CampaignNewRepository(db)

	userRepository := repository.UserNewRepository(db)
	userService := service.UserNewService(userRepository)

	transactionRepository := repository.TransactionNewRepository(db)
	transactionService := service.TransactionNewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := NewTransactionHandler(transactionService)

	api.GET("/campaigns/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
}
