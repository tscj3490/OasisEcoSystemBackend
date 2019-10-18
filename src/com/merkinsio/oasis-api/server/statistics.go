package server

import (
	"com/merkinsio/oasis-api/domain"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPlantationStatistics(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO: Handle current user role
	// currentUser := context.Get(r, "currentUser").(*domain.User)
	vars := mux.Vars(r)

	plantationID := parseStringToObjectID(vars["plantationId"])

	if plantationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plantation ID %v", vars["plantationId"]))
		return
	}

	plantationStatistics, err := domain.GetPlantationStatisticsByPlantationID(*plantationID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, plantationStatistics)
	}
}
