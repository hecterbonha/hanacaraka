package routes

import (
	"com.hanacaraka/interfaces/http/handlers"
	"github.com/gorilla/mux"
)

// SetupProductRoutes configures all product-related routes
func SetupProductRoutes(router *mux.Router, productHandler *handlers.ProductHandler) {
	// Product CRUD operations
	router.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	router.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Additional product routes
	router.HandleFunc("/products/{id}/reviews", productHandler.GetProductReviews).Methods("GET")
	router.HandleFunc("/products/{id}/reviews", productHandler.CreateProductReview).Methods("POST")
	router.HandleFunc("/products/search", productHandler.SearchProducts).Methods("GET")
	router.HandleFunc("/products/categories", productHandler.GetProductCategories).Methods("GET")
	router.HandleFunc("/products/category/{category}", productHandler.GetProductsByCategory).Methods("GET")
}

// SetupProductAPIRoutes configures versioned API product routes
func SetupProductAPIRoutes(apiRouter *mux.Router, productHandler *handlers.ProductHandler) {
	// API v1 product routes
	apiRouter.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	apiRouter.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	apiRouter.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	apiRouter.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	apiRouter.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// API-specific product endpoints
	apiRouter.HandleFunc("/products/bulk", productHandler.BulkCreateProducts).Methods("POST")
	apiRouter.HandleFunc("/products/export", productHandler.ExportProducts).Methods("GET")
	apiRouter.HandleFunc("/products/import", productHandler.ImportProducts).Methods("POST")
	apiRouter.HandleFunc("/products/featured", productHandler.GetFeaturedProducts).Methods("GET")
	apiRouter.HandleFunc("/products/{id}/inventory", productHandler.UpdateProductInventory).Methods("PUT")
}
