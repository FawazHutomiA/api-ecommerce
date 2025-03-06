package cloudinary

import (
	"example/pkg/app"
	"example/pkg/middleware"

	"github.com/go-chi/chi"
)

func SetupCloudinaryRoutes(r chi.Router, app app.AppConfig) {
	// Initialize repositories and services
	cloudinaryService := NewCloudinaryService(app)
	cloudinaryHandler := NewCloudinaryHandler(app, cloudinaryService)

	r.With(middleware.AuthMiddleware).Route("/upload-file", func(r chi.Router) {
		r.Post("/", cloudinaryHandler.UploadImage)
	})
}
