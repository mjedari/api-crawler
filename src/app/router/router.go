package router

import (
	"github.com/gorilla/mux"
	"github.com/mjedari/vgang-project/src/app/handler"
	"github.com/mjedari/vgang-project/src/app/wiring"
	"net/http"
)

func NewRouter() *mux.Router {
	// Initialize handlers
	productHandler := handler.NewProductHandler(wiring.Wiring.GetStorage())

	// Initialize middleware
	logger := handler.LoggerMiddleware
	limiter := handler.RateLimiterMiddleware

	// Create new mux router
	router := mux.NewRouter()

	// Define routes with middleware
	router.HandleFunc("/{key:[a-zA-Z0-9]{5}}", logger(limiter(productHandler.GetProduct))).Methods(http.MethodGet)
	router.HandleFunc("/product/all", logger(limiter(productHandler.GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/product/short-links", logger(limiter(productHandler.GetShortLinks))).Methods(http.MethodGet)

	return router
}
