package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//GetAllVariaty get all the variaty crop
func GetAllVariaty(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var varieties []domain.Variaty
	var err error
	varieties, err = domain.GetVariaty()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, varieties)
	}
}

//CreateVariatyCrop create the variaty crop
func CreateVariatyCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Variaty
	parseBody(r, &body)

	variaty, err := services.StoreVariaty(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, variaty)
}

//DeleteVariatyCrop delete the variaty crop
func DeleteVariatyCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	variatyID := parseStringToObjectID(vars["variatyId"])

	if variatyID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid param id %v", vars["variatyId"]))
		return
	}

	if err := domain.DeleteVariatyByID(*variatyID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

//UpdateVariatyCrop updates the variaty crop
func UpdateVariatyCrop(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Variaty
	parseBody(r, &body)

	vars := mux.Vars(r)

	variatyIDStr := vars["variatyId"]

	if valid := bson.IsObjectIdHex(variatyIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid variatyId")
		return
	}

	variatyID := bson.ObjectIdHex(variatyIDStr)
	body.ID = variatyID

	err := domain.UpdateVariatyByVariaty(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		variaty, err := domain.GetVariatyByID(variatyID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, variaty)
		}
	}
}

func GetConditionalVariaties(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	declaredCultivation := r.FormValue("declaredCultivation")
	var varieties []domain.Variaty
	var err error
	varieties, err = domain.GetConditionalVariaties(declaredCultivation)

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, varieties)
	}
}
