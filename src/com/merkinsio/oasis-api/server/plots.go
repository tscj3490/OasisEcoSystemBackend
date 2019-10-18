package server

import (
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//CreatePlot creates and returns to the client the created plot
func CreatePlot(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Plot
	parseBody(r, &body)

	plot, err := services.StorePlot(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plot)
}

//GetPlotByID creates and returns to the client the created plot
func GetPlotByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	plotID := parseStringToObjectID(vars["plotId"])

	if plotID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plotID %v", vars["plotId"]))
		return
	}

	plot, err := domain.GetPlotByID(*plotID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, plot)
	}
}

func GetPlotsByPlantationID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	plantationID := parseStringToObjectID(vars["plantationId"])

	if plantationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plantationID %v", vars["plantationId"]))
		return
	}

	plots, err := domain.GetPlotsByPlantationID(*plantationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, plots)
	}
}

func DeletePlot(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	plotIDStr := vars["plotId"]
	var plotID bson.ObjectId

	if isValid := bson.IsObjectIdHex(plotIDStr); isValid == false {
		respondWithError(w, http.StatusBadRequest, "invalid id")
	} else {
		plotID = bson.ObjectIdHex(plotIDStr)
	}

	if err := domain.DeletePlotByID(plotID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdatePlot(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Plot
	parseBody(r, &body)

	vars := mux.Vars(r)

	plotIDStr := vars["plotId"]

	if valid := bson.IsObjectIdHex(plotIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid plotId")
		return
	}

	plotID := bson.ObjectIdHex(plotIDStr)
	body.ID = plotID

	err := domain.UpdatePlot(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		plot, err := domain.GetPlotByID(plotID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, plot)
		}
	}
}

func GetPlotsWithGeoJSONByPlantationID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	plantationID := parseStringToObjectID(vars["plantationId"])

	if plantationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plantationID %v", vars["plantationId"]))
		return
	}

	plots, err := services.GetPlotsWithGeoJSONByPlantationID(*plantationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, plots)
	}
}
