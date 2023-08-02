package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mjedari/vgang-project/src/domain/contracts"
	"github.com/mjedari/vgang-project/src/domain/products"
	"net/http"
)

type ProductHandler struct {
	storage contracts.IStorage
}

func NewProductHandler(storage contracts.IStorage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Handle /product/all
	var productList []products.Product
	data := h.storage.FetchAll(r.Context(), "products:*")

	for _, item := range data {
		var product products.Product
		err := json.Unmarshal(item, &product)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		product.AddLink()

		productList = append(productList, product)

	}

	response, _ := json.Marshal(productList)
	writeResponse(w, response)
}

func (h *ProductHandler) GetShortLinks(w http.ResponseWriter, r *http.Request) {
	// Handle /product/short-links
	w.Write([]byte("there"))
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// todo: validation
	vars := mux.Vars(r)
	key := fmt.Sprintf("products:%v", vars["key"])

	data := h.storage.Fetch(r.Context(), key)
	var product products.Product

	err := json.Unmarshal(data, &product)
	if err != nil {
		// handle
	}

	product.AddLink()

	response, _ := json.Marshal(product)
	writeResponse(w, response)
}

func writeResponse(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
