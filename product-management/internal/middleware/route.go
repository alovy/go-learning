package middleware

import (
	"database/sql"
	"product-api/internal/controller"
	"product-api/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(db *sql.DB) *chi.Mux {
	// Set up the router
	r := chi.NewRouter()

	r.Get("/generate-token", controller.GenerateJWTToken)

	// Create the service
	productService := service.NewProductService(db)

	r.Route("/products", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Mount("/", controller.NewProduct(productService).Router())
	})

	return r
}
