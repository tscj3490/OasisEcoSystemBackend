package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetAllCooperatives Returns all the cooperatives
func GetAllCooperatives(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var cooperatives []domain.Cooperative
	var err error
	cooperatives, err = domain.GetCooperatives()

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, cooperatives)
	}
}

//GetCooperativeByID returns the cooperative whose Id you put on the query
func GetCooperativeByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	cooperativeIDStr := vars["cooperativeId"]
	coopID := parseStringToObjectID(cooperativeIDStr)
	if coopID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID %v", cooperativeIDStr))
		return
	}

	coop, err := domain.GetCooperativeByID(*coopID)

	if err != nil {
		errorCode := http.StatusBadRequest
		errorMSG := err.Error()
		if err.Error() == "not found" {
			errorCode = http.StatusNotFound
			errorMSG = fmt.Sprintf("Cooperative not found for id %v", cooperativeIDStr)
		}

		respondWithError(w, errorCode, errorMSG)
	} else {
		respondWithJSON(w, http.StatusOK, coop)
	}
}

//CreateCooperative creates and returns to the client the created cooperative
func CreateCooperative(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Cooperative
	parseBody(r, &body)

	cooperative, err := services.StoreCooperative(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, cooperative)
}

func DeleteCooperative(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	cooperativeID := parseStringToObjectID(vars["cooperativeId"])

	if cooperativeID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid cooperative id %v", vars["cooperativeId"]))
		return
	}

	if err := domain.DeleteCooperativeByID(*cooperativeID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdateCooperative(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Cooperative
	parseBody(r, &body)

	vars := mux.Vars(r)

	cooperativeIDStr := vars["cooperativeId"]

	if valid := bson.IsObjectIdHex(cooperativeIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid cooperativeId")
		return
	}

	cooperativeID := bson.ObjectIdHex(cooperativeIDStr)
	body.ID = cooperativeID

	err := domain.UpdateCooperativeByCooperative(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		cooperative, err := domain.GetCooperativeByID(cooperativeID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, cooperative)
		}
	}
}

func GetCooperativeWorkers(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	cooperativeIDStr := vars["cooperativeId"]
	coopID := parseStringToObjectID(cooperativeIDStr)
	if coopID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ID %v", cooperativeIDStr))
		return
	}

	workers, err := domain.GetCooperativeWorkersByID(*coopID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, workers)
	}
}
