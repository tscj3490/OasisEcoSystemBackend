package server

import (
	"com/merkinsio/oasis-api/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func GetSigPacPlotByPlotID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	plotID := parseStringToObjectID(vars["plotId"])

	if plotID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plotID %v", vars["plotId"]))
		return
	}

	plot, err := domain.GetPlotByID(*plotID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	sigpacPlot, err := domain.GetSigpacPlotByPlotCode(plot.PlotCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacPlot)
	}
}

func GetSigPacPlotByPlotIDWithNoUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	plotIDStr := vars["plotId"]

	if valid := bson.IsObjectIdHex(plotIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid plotId")
		return
	}

	plotID := bson.ObjectIdHex(plotIDStr)

	plot, err := domain.GetPlotByID(plotID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	sigpacPlot, err := domain.GetSigpacPlotByPlotCode(plot.PlotCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacPlot)
	}
}

func GetSigpacPlot(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	provinceRec := r.FormValue("provinceRec")
	townRec := r.FormValue("townRec")
	polygonNum := r.FormValue("polygonNum")
	plotNum := r.FormValue("plotNum")
	enclosure := r.FormValue("enclosure")
	plotCode := domain.PlotCodification{}

	provinceRecInt, err := strconv.Atoi(provinceRec)
	if err != nil {
		log.Errorf("Error in provinceRec parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	townRecInt, err := strconv.Atoi(townRec)
	if err != nil {
		log.Errorf("Error in townRec parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	polygonNumInt, err := strconv.Atoi(polygonNum)
	if err != nil {
		log.Errorf("Error in polygonNum parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	plotNumInt, err := strconv.Atoi(plotNum)
	if err != nil {
		log.Errorf("Error in plotNum parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	enclosureInt, err := strconv.Atoi(enclosure)
	if err != nil {
		log.Errorf("Error in enclosure parser -> error: %v", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	plotCode.ProvinceRec = provinceRecInt
	plotCode.TownRec = townRecInt
	plotCode.PolygonNum = polygonNumInt
	plotCode.PlotNum = plotNumInt
	plotCode.Enclosure = enclosureInt

	sigpacPlot, err := domain.GetSigpacPlotByPlotCode(plotCode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacPlot)
	}
}
