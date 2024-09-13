package product

import (
	"example/internal/repository/postgresql/product"
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupProductRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	productRepository := product.NewProductRepository(app)
	productService := NewProductService(app, productRepository)
	productHandler := NewProductHandler(app, productService)

	r.With(middleware.AuthMiddleware).Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.ListPaginate)
		r.Get("/{id:[a-fA-F0-9-]{36}}", productHandler.Detail)
		r.Get("/{slug:[a-zA-Z0-9_-]+}", productHandler.DetailBySlug)
		r.Post("/", productHandler.Create)
		r.Put("/{id:[a-fA-F0-9-]{36}}", productHandler.Update)
		r.Delete("/{id:[a-fA-F0-9-]{36}}", productHandler.Delete)
	})
}
