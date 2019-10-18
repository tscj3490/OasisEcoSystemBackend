package server

import (
	"com/merkinsio/oasis-api/domain"
	"net/http"
)

func GetAllProvinces(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var provinces []domain.Province
	var err error
	provinces, err = domain.GetProvinces()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, provinces)
	}
}
