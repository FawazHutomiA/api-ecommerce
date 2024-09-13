package transaction

import (
	"example/internal/module/payment"
	"example/internal/repository/postgresql/product"
	"example/internal/repository/postgresql/transaction"
	"example/internal/repository/postgresql/user"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupTransactionRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	productRepository := product.NewProductRepository(app)

	userRepository := user.NewUserRepository(app)

	paymentService := payment.PaymentNewService()

	transactionRepository := transaction.NewTransactionRepository(app)
	transactionService := NewTransactionService(app, transactionRepository, productRepository, userRepository, paymentService)
	transactionHandler := NewTransactionHandler(app, transactionService)

	r.Post("/transactions/notification", transactionHandler.GetNotification)

	r.With(middleware.AuthMiddleware).Route("/transactions", func(r chi.Router) {
		r.Get("/", transactionHandler.ListPaginate)
		r.Get("/{id:[a-fA-F0-9-]{36}}", transactionHandler.Detail)
		r.Get("/product/{id:[a-fA-F0-9-]{36}}", transactionHandler.DetailByProductID)
		r.Get("/user/{id:[a-fA-F0-9-]{36}}", transactionHandler.DetailByUserID)
		r.Post("/", transactionHandler.Create)
		r.Put("/{id:[a-fA-F0-9-]{36}}", transactionHandler.Update)
		r.Delete("/{id:[a-fA-F0-9-]{36}}", transactionHandler.Delete)
	})
}
