package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetVisitCards(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)

	visitCards, err := domain.GetVisitCardsByCooperativeID(*currentUser.CooperativeID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, visitCards)
	}
}

func GetVisitCardsByCooperativeID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	cooperativeID := parseStringToObjectID(vars["cooperativeId"])

	if cooperativeID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid cooperative id %v", vars["cooperativeId"]))
		return
	}

	if visitCards, err := domain.GetVisitCardsByCooperativeID(*cooperativeID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, visitCards)
	}
}

func GetVisitCardByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	visitCardID := parseStringToObjectID(vars["visitCardId"])

	if visitCardID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid visitCard id %v", vars["visitCardId"]))
		return
	}

	if visitCard, err := domain.GetVisitCardByID(*visitCardID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, visitCard)
	}
}

func CreateVisitCard(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.VisitCard
	parseBody(r, &body)
	currentUser := context.Get(r, "currentUser").(*domain.User)
	body.CooperativeID = *currentUser.CooperativeID

	visitCard, err := services.StoreVisitCard(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, visitCard)
}

func DeleteVisitCard(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	visitCardID := parseStringToObjectID(vars["visitCardId"])

	if visitCardID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid visitCard id %v", vars["visitCardId"]))
		return
	}

	if err := domain.DeleteVisitCardByID(*visitCardID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdateVisitCard(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.VisitCard
	parseBody(r, &body)

	vars := mux.Vars(r)

	visitCardIDStr := vars["visitCardId"]

	if valid := bson.IsObjectIdHex(visitCardIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	visitCardID := bson.ObjectIdHex(visitCardIDStr)
	body.ID = visitCardID

	err := domain.UpdateVisitCardByVisitCard(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		visitCard, err := domain.GetVisitCardByID(visitCardID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, visitCard)
		}
	}
}
