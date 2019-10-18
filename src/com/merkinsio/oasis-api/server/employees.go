package server

import (
	"com/merkinsio/oasis-api/domain"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetAllEmployees(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	employees, err := domain.GetAllEmployees()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, employees)
	}

}

func GetEmployeesByCooperativeID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	cooperativeIDStr := vars["cooperativeId"]

	if valid := bson.IsObjectIdHex(cooperativeIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	cooperativeID := bson.ObjectIdHex(cooperativeIDStr)

	employees, err := domain.GetEmployeesByCooperativeID(&cooperativeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, employees)
	}

}
