package router

import (
	"example/pkg/app"
	"example/pkg/log"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func Router(app app.AppConfig) {
	logger := log.New()
	router := chi.NewRouter()

	// Use CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	SetupRoutes(router, app)

	logger.Infof("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Failed to start server:", err)
	}

}
