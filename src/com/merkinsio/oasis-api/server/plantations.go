package server

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//CreatePlantation creates and returns to the client the created plantation
func CreatePlantation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Plantation
	parseBody(r, &body)

	plantation, err := services.StorePlantation(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plantation)
}

func GetPlantationsByCurrentUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)

	// cooperativeID := parseStringToObjectID(currentUser.CooperativeID)

	if currentUser.Role == constants.RoleFarmer {
		plantations, err := domain.GetPlantationsByOwnerID(currentUser.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, plantations)
		}
	} else if currentUser.CooperativeID == nil {
		// respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid cooperativeID %v", currentUser.CooperativeID))

		// Return an empty array to avoid errors (the user may not habve been assigned yet)
		respondWithJSON(w, http.StatusOK, []domain.Plantation{})
		return
	} else {

		plantations, err := domain.GetPlantationsByUserCooperativeID(*currentUser.CooperativeID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, plantations)
		}
	}
}

func GetPlantationsAndPlotsByCurrentUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)

	if currentUser.Role == constants.RoleFarmer {
		plantations, err := domain.GetPlantationsAndPlotsByOwnerId(currentUser.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			for i, plantation := range plantations {
				if plantation.Plots == nil {
					plots, err := domain.GetPlotsByPlantationID(plantation.ID)
					if err != nil {
						respondWithError(w, http.StatusInternalServerError, err.Error())
					} else {
						plantations[i].Plots = &plots
					}
				}
			}
			respondWithJSON(w, http.StatusOK, plantations)
		}
	} else if currentUser.CooperativeID == nil {
		// respondWithError(w, http.StatusNotFound, fmt.Sprintf("Invalid cooperativeID %v", currentUser.CooperativeID))

		// Return an empty array to avoid errors (the user may not habve been assigned yet)
		respondWithJSON(w, http.StatusOK, []domain.Plantation{})
		return
	} else {
		plantations, err := domain.GetPlantationsAndPlotsByCooperativeID(*currentUser.CooperativeID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, plantations)
		}
	}
}

func GetPlantationIssues(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO: Handle current user role
	// currentUser := context.Get(r, "currentUser").(*domain.User)
	vars := mux.Vars(r)

	plantationID := parseStringToObjectID(vars["plantationId"])

	if plantationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plantation ID %v", vars["plantationId"]))
		return
	}

	issues, err := domain.GetIssuesByPlantationID(*plantationID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issues)
	}
}

func GetPlantationFinishedIssues(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO: Handle current user role
	// currentUser := context.Get(r, "currentUser").(*domain.User)
	vars := mux.Vars(r)

	plantationID := parseStringToObjectID(vars["plantationId"])

	if plantationID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid plantation ID %v", vars["plantationId"]))
		return
	}

	issues, err := domain.GetFinishedIssuesByPlantationID(*plantationID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issues)
	}
}

func DeletePlantation(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	plantationIDStr := vars["plantationId"]
	var plantationID bson.ObjectId

	if isValid := bson.IsObjectIdHex(plantationIDStr); isValid == false {
		respondWithError(w, http.StatusBadRequest, "invalid id")
	} else {
		plantationID = bson.ObjectIdHex(plantationIDStr)
	}

	if err := domain.DeletePlantationByID(plantationID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdatePlantationByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body *domain.Plantation
	parseBody(r, &body)

	vars := mux.Vars(r)

	plantationIDStr := vars["plantationId"]

	if valid := bson.IsObjectIdHex(plantationIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid plantationId")
		return
	}

	plantationID := bson.ObjectIdHex(plantationIDStr)
	body.ID = plantationID

	err := domain.UpdatePlantationByID(*body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		plantation, err := domain.GetPlantationByID(plantationID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, plantation)
		}
	}
}

func GetPlantationSigpacPlots(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	plantationIDStr := vars["plantationId"]

	if valid := bson.IsObjectIdHex(plantationIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid plantationId")
		return
	}

	plantationID := bson.ObjectIdHex(plantationIDStr)

	sigpacPlots, err := services.GetSigpacPlotsByPlantationID(plantationID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacPlots)
	}
}

func GetPlantationSigpacPlotsWithNoUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	plantationIDStr := vars["plantationId"]

	if valid := bson.IsObjectIdHex(plantationIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid plantationId")
		return
	}

	plantationID := bson.ObjectIdHex(plantationIDStr)

	sigpacPlots, err := services.GetPlotsWithGeoJSONByPlantationID(plantationID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, sigpacPlots)
	}
}
