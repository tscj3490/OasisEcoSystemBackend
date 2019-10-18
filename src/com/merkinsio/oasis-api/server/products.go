package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var products []domain.Product
	var err error
	products, err = domain.GetProducts()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, products)
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Product
	if err := parseBody(r, &body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := services.StoreProduct(body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, product)
	}

}

func DeleteProduct(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	productID := parseStringToObjectID(vars["productId"])

	if productID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid product id %v", vars["productId"]))
		return
	}

	if err := domain.DeleteProductByID(*productID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}
