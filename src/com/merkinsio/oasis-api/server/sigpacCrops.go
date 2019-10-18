package server

import (
	"com/merkinsio/oasis-api/domain"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetSigpacCropById(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	sigpacCropID := parseStringToObjectID(vars["sigpacCropId"])

	if sigpacCropID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid sigpacCropID %v", vars["sigpacCropId"]))
		return
	}

	sigpacCrop, err := domain.GetSigpacCropByID(*sigpacCropID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacCrop)
	}
}

func GetAllSigpacCrops(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var sigpacCrops []domain.SigpacCrop
	var err error
	sigpacCrops, err = domain.GetAllSigpacCrops()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacCrops)
	}
}
