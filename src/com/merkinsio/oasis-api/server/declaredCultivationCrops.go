package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//GetAllDeclaredCultivation get all the cultivation crop
func GetAllDeclaredCultivationCrops(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var declaredCultivationCrops []domain.DeclaredCultivation
	var err error
	declaredCultivationCrops, err = domain.GetAllDeclaredCultivations()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, declaredCultivationCrops)
	}
}

//CreateDeclaredCultivationCrop create the cultivation crop
func CreateDeclaredCultivationCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.DeclaredCultivation
	parseBody(r, &body)

	declaredCultivation, err := services.StoreDeclaredCultivation(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, declaredCultivation)
}

//DeleteDeclaredCultivationCrop delete the cultivation crop
func DeleteDeclaredCultivationCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	declaredCultivationID := parseStringToObjectID(vars["declaredCultivationId"])

	if declaredCultivationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid param id %v", vars["declaredCultivationId"]))
		return
	}

	if err := domain.DeleteDeclaredCultivationByID(*declaredCultivationID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

//UpdateDeclaredCultivationCrop updates the cultivation crop
func UpdateDeclaredCultivationCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.DeclaredCultivation
	parseBody(r, &body)

	vars := mux.Vars(r)

	declaredCultivationIDStr := vars["declaredCultivationId"]

	if valid := bson.IsObjectIdHex(declaredCultivationIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid declaredCultivationId")
		return
	}

	declaredCultivationID := bson.ObjectIdHex(declaredCultivationIDStr)
	body.ID = declaredCultivationID

	err := domain.UpdateDeclaredCultivationByDeclaredCultivation(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		declaredCultivation, err := domain.GetDeclaredCultivationByID(declaredCultivationID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, declaredCultivation)
		}
	}
}
