package server

import (
	"com/merkinsio/oasis-api/domain"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetTownsByProvince(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var towns []domain.Town
	province := r.FormValue("province")
	var err error
	provinceInt, err := strconv.Atoi(province)
	if err != nil {
		log.Errorf("Error in province parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	towns, err = domain.GetTownsByProvince(provinceInt)

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, towns)
	}
}
