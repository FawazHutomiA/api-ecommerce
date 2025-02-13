package warehouse

import (
	"example/internal/repository/postgresql/token"
	"example/internal/repository/postgresql/warehouse"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupWarehouseRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	tokenRepository := token.NewTokenRepository(app)

	warehouseRepository := warehouse.NewWarehouseRepository(app)
	warehouseService := NewWarehouseService(app, warehouseRepository, tokenRepository)
	warehouseHandler := NewWarehouseHandler(app, warehouseService)

	r.With(middleware.AuthMiddleware).Route("/warehouses", func(r chi.Router) {
		r.Get("/", warehouseHandler.ListPaginate)
		r.Get("/{id:[a-fA-F0-9-]{36}}", warehouseHandler.Detail)
	})
}
