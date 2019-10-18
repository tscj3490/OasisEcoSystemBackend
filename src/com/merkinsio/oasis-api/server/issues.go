package server

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//CreateIssue creates and returns to the client the created issue
func CreateIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Issue
	parseBody(r, &body)

	currentUser := context.Get(r, "currentUser").(*domain.User)

	fmt.Println(currentUser.ID.Hex())

	if currentUser == nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	plantation, err := domain.GetPlantationByID(body.Location.PlantationID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	plot, err := domain.GetPlotByID(body.Location.PlotID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.VisitInfo != nil && body.WorkOrderInfo != nil {
		//TODO: definir cómo proceder en este caso!
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	body.Issuer = *currentUser
	body.Location.PlantationName = plantation.Name
	body.Location.PlotCode = plot.PlotCode
	body.CropsType = plot.CropsName
	body.CropVariety = plot.Variety
	body.Status = constants.StatusCreated
	body.IssueNumber = int64(time.Now().UTC().UnixNano() * int64(time.Nanosecond) / int64(time.Millisecond))

	if visitInfo := body.VisitInfo; visitInfo != nil {
		if *currentUser.CooperativeID != plantation.CooperativeID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		if currentUser.Role != constants.RoleTechnician && currentUser.Role != constants.RoleAdmin {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		visitInfo.Technician = *currentUser
		visitInfo.TechnicianID = &visitInfo.Technician.ID
		if visitInfo.NeedsWorkOrder {
			body.Status = constants.StatusPendingWorkOrder
		} else {
			body.Status = constants.StatusFinished
		}
		body.IssueComment = "Incidencia creada a partir de visita"
	}

	if workOrderInfo := body.WorkOrderInfo; workOrderInfo != nil {
		if *currentUser.CooperativeID != plantation.CooperativeID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		if currentUser.Role != constants.RoleWorker && currentUser.Role != constants.RoleAdmin {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		now := time.Now().UTC()
		workOrderInfo.WorkDate = &now
		workOrderInfo.WorkOrderTimer = &domain.WorkOrderTimer{
			StartTime: now,
		}
		workOrderInfo.Worker = *currentUser
		workOrderInfo.WorkerID = &workOrderInfo.Worker.ID
		workOrderInfo.Accepted = true
		if workOrderInfo.Completed {
			body.Status = constants.StatusFinished
		} else {
			body.Status = constants.StatusWorking
		}
		body.IssueComment = "Incidencia creada a partir de orden de trabajo"

		if body.VisitInfo == nil {
			body.VisitInfo = &domain.Visit{
				NeedsWorkOrder: true,
				Technician:     *currentUser,
				TechnicianID:   &currentUser.ID,
				VisitDate:      now,
			}
		} else {
			body.VisitInfo.NeedsWorkOrder = true
		}
	}

	issue, err := services.StoreIssue(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, issue)
}

//UpdateIssue updates the given issue in the database. If it's being updated into a visit and the visitor is not specified, the current user will be assigned. The same is applied if it's being updated into a task, if the worker is not specified, the current user will be assigned.
func UpdateIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := context.Get(r, "currentUser").(*domain.User)
	var body domain.Issue
	parseBody(r, &body)

	vars := mux.Vars(r)

	issueIDStr := vars["issueId"]

	if valid := bson.IsObjectIdHex(issueIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	issueID := bson.ObjectIdHex(issueIDStr)
	body.ID = issueID

	issue, err := services.UpdateIssue(&body, *currentUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issue)
	}
}

// StartWorkingIssue starts the the given issue's workOrder timer.
func StartWorkingIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)

	issueIDStr := vars["issueId"]

	if valid := bson.IsObjectIdHex(issueIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	issueID := bson.ObjectIdHex(issueIDStr)
	issue, err := domain.GetIssueByID(issueID)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Issue not found for id %v", issueIDStr))
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	issue.ID = issueID
	issue, err = services.StartWorkingIssue(issue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issue)
	}
}

//PauseWorkingIssue stops the given issue's worOrder timer.
func PauseWorkingIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)

	issueIDStr := vars["issueId"]

	if valid := bson.IsObjectIdHex(issueIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	issueID := bson.ObjectIdHex(issueIDStr)
	issue, err := domain.GetIssueByID(issueID)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Issue not found for id %v", issueIDStr))
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	issue.ID = issueID
	issue, err = services.PauseWorkingIssue(issue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issue)
	}
}

//FinishIssue marks the given issue as finished.
func FinishIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)

	issueIDStr := vars["issueId"]

	if valid := bson.IsObjectIdHex(issueIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	issueID := bson.ObjectIdHex(issueIDStr)
	issue, err := domain.GetIssueByID(issueID)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Issue not found for id %v", issueIDStr))
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	issue.ID = issueID
	issue.Status = constants.StatusFinished

	err = domain.UpdateIssueByIssue(*issue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		res, err := domain.GetIssueByID(issue.ID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, res)
		}
	}
}

//AcceptWorkOrderIssue accepts the workOrder of the given issue
func AcceptWorkOrderIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)

	issueIDStr := vars["issueId"]

	if valid := bson.IsObjectIdHex(issueIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	issueID := bson.ObjectIdHex(issueIDStr)
	issue, err := domain.GetIssueByID(issueID)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Issue not found for id %v", issueIDStr))
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	issue.ID = issueID
	issue.WorkOrderInfo.Accepted = true
	issue.Status = constants.StatusAccepted

	err = domain.UpdateIssueByIssue(*issue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		res, err := domain.GetIssueByID(issue.ID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, res)
		}
	}
}

//GetAllIssues returns all the unfinished issues that the current User has access to
func GetAllIssues(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var issues []domain.Issue
	var err error

	currentUser := context.Get(r, "currentUser").(*domain.User)

	switch currentUser.Role {
	case constants.RoleAdmin:
		issues, err = domain.GetAllIssues()
	case constants.RoleTechnician:
		issues, err = domain.GetAllIssues()
	case constants.RoleFarmer:
		issues, err = services.GetOwnedPlantationIssuesByUserID(currentUser.ID)
	default:
		issues, err = domain.GetIssuesByUserID(currentUser.ID)
	}

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issues)
	}
}

//GetMyIssues returns the currentUser's personal issues (not sent to any cooperative)
func GetMyIssues(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var issues []domain.Issue
	var err error
	vars := mux.Vars(r)
	var body domain.GetMyIssuesBody

	plantationIDStr := vars["plantationId"]

	if valid := bson.IsObjectIdHex(plantationIDStr); valid == false {
		respondWithError(w, http.StatusBadRequest, "Invalid issueId")
		return
	}

	plantationID := bson.ObjectIdHex(plantationIDStr)

	currentUser := context.Get(r, "currentUser").(*domain.User)

	body.PlantationID = plantationID
	body.UserID = currentUser.ID

	issues, err = domain.GetMyIssuesByUserID(body)

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issues)
	}
}

//GetAllFinishedIssues returns all the finished issues that the current user has access to
func GetAllFinishedIssues(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var issues []domain.Issue
	var err error

	currentUser := context.Get(r, "currentUser").(*domain.User)

	switch currentUser.Role {
	case constants.RoleAdmin:
		issues, err = domain.GetAllFinishedIssues()
	default:
		issues, err = domain.GetFinishedIssuesByUserID(currentUser.ID)
	}

	if err != nil && err.Error() != "not found" {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issues)
	}
}

//GetIssueByID returns the issue that matches the ID given
func GetIssueByID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	issueID := parseStringToObjectID(vars["issueId"])

	if issueID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid issueID %v", vars["issueId"]))
		return
	}

	issue, err := domain.GetIssueByID(*issueID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, issue)
	}
}

//DeleteIssue updates the given issue to inactivate it so it will no longer be fetched by any method
func DeleteIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	issueID := parseStringToObjectID(vars["issueId"])

	if issueID == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid issue id %v", vars["issueId"]))
		return
	}

	if err := domain.DeleteIssueByID(*issueID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

//CreateTaskIssue creates and returns to the client the created issue from task
func CreateTaskIssue(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var body domain.Issue
	parseBody(r, &body)

	currentUser := context.Get(r, "currentUser").(*domain.User)

	fmt.Println(currentUser.ID.Hex())

	if currentUser == nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	plantation, err := domain.GetPlantationByID(body.Location.PlantationID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	plot, err := domain.GetPlotByID(body.Location.PlotID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.VisitInfo != nil && body.WorkOrderInfo != nil {
		//TODO: definir cómo proceder en este caso!
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	body.Issuer = *currentUser
	body.Location.PlantationName = plantation.Name
	body.Location.PlotCode = plot.PlotCode
	body.CropsType = plot.CropsName
	body.CropVariety = plot.Variety
	body.Status = constants.StatusCreated
	body.IssueNumber = int64(time.Now().UTC().UnixNano() * int64(time.Nanosecond) / int64(time.Millisecond))

	if visitInfo := body.VisitInfo; visitInfo != nil {
		if *currentUser.CooperativeID != plantation.CooperativeID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		if currentUser.Role != constants.RoleTechnician && currentUser.Role != constants.RoleAdmin {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		visitInfo.Technician = *currentUser
		visitInfo.TechnicianID = &visitInfo.Technician.ID
		if visitInfo.NeedsWorkOrder {
			body.Status = constants.StatusPendingWorkOrder
		} else {
			body.Status = constants.StatusFinished
		}
		body.IssueComment = "Incidencia creada a partir de visita"
	}

	if workOrderInfo := body.WorkOrderInfo; workOrderInfo != nil {
		if *currentUser.CooperativeID != plantation.CooperativeID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		workOrderInfo.Accepted = true
		body.Status = constants.StatusAccepted
		body.IssueComment = "Incidencia creada a partir de orden de trabajo"

		if body.VisitInfo == nil {
			body.VisitInfo = &domain.Visit{
				NeedsWorkOrder: true,
				Technician:     *currentUser,
				TechnicianID:   &currentUser.ID,
				VisitDate:      time.Now(),
			}
		} else {
			body.VisitInfo.NeedsWorkOrder = true
		}
	}

	issue, err := services.StoreIssue(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, issue)
}

//GetIssueStatuses returns every status the issues go by
func GetIssueStatuses(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var statuses []string
	statuses = append(statuses, constants.StatusCreated)
	statuses = append(statuses, constants.StatusPendingVisit)
	statuses = append(statuses, constants.StatusPendingWorkOrder)
	statuses = append(statuses, constants.StatusWorkOrder)
	statuses = append(statuses, constants.StatusAccepted)
	statuses = append(statuses, constants.StatusAssigned)
	statuses = append(statuses, constants.StatusWorking)
	statuses = append(statuses, constants.StatusPaused)
	statuses = append(statuses, constants.StatusPendingAssessment)
	statuses = append(statuses, constants.StatusFinished)
	statuses = append(statuses, constants.StatusCancelled)
	statuses = append(statuses, constants.StatusDelayed)
	respondWithJSON(w, http.StatusOK, statuses)
}

//GetTreatmentTypes returns every treatment type
func GetTreatmentTypes(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var treatments []string
	treatments = append(treatments, constants.TreatmentCultivation)
	treatments = append(treatments, constants.TreatmentSowing)
	treatments = append(treatments, constants.TreatmentPlant)
	treatments = append(treatments, constants.TreatmentFertilization)
	treatments = append(treatments, constants.TreatmentPruning)
	treatments = append(treatments, constants.TreatmentPhytosanitary)
	treatments = append(treatments, constants.TreatmentIrrigation)
	treatments = append(treatments, constants.TreatmentTwigsRemoval)
	treatments = append(treatments, constants.TreatmentHarvesting)
	treatments = append(treatments, constants.TreatmentAnalysis)
	treatments = append(treatments, constants.TreatmentOtherTask)
	respondWithJSON(w, http.StatusOK, treatments)
}
