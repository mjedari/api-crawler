package router

import (
	"github.com/gorilla/mux"
	"github.com/mjedari/vgang-project/app/handler"
	"github.com/mjedari/vgang-project/app/wiring"
	"net/http"
)

func NewRouter() *mux.Router {
	productHandler := handler.NewProductHandler(wiring.Wiring.GetStorage())

	logger := handler.LoggerMiddleware
	limiter := handler.RateLimiterMiddleware

	router := mux.NewRouter()

	router.HandleFunc("/{key:[a-zA-Z0-9]{5}}", logger(limiter(productHandler.GetProduct))).Methods(http.MethodGet)
	router.HandleFunc("/product/all", logger(limiter(productHandler.GetAll))).Methods(http.MethodGet)
	router.HandleFunc("/product/short-links", logger(limiter(productHandler.GetShortLinks))).Methods(http.MethodGet)

	return router
}
