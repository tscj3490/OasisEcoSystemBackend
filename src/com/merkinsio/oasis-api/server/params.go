package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetAllParams(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var params []domain.Param
	var err error
	params, err = domain.GetParams()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, params)
	}
}

func GetConditionalParams(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var params []domain.Param
	var err error
	vars := mux.Vars(r)

	treatmentType := vars["treatment"]
	params, err = services.GetParamsByTreatmentType(treatmentType)

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, params)
	}
}

func CreateParam(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Param
	parseBody(r, &body)

	param, err := services.StoreParam(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, param)
}

func DeleteParam(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	paramID := parseStringToObjectID(vars["paramId"])

	if paramID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid param id %v", vars["paramId"]))
		return
	}

	if err := domain.DeleteParamByID(*paramID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdateParam(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Param
	parseBody(r, &body)

	vars := mux.Vars(r)

	paramIDStr := vars["paramId"]

	if valid := bson.IsObjectIdHex(paramIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	paramID := bson.ObjectIdHex(paramIDStr)
	body.ID = paramID

	err := domain.UpdateParamByParam(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		param, err := domain.GetParamByID(paramID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, param)
		}
	}
}
