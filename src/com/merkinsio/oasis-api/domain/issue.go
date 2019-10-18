package domain

import (
	"com/merkinsio/oasis-api/constants"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Issue holds the issue information
type Issue struct {
	Entity                 `bson:",inline"`
	Pictures               []File        `json:"pictures" bson:"pictures"`
	CropsType              string        `json:"cropsType" bson:"cropsType"`
	CropVariety            string        `json:"cropVariety" bson:"cropVariety"`
	ReportObservations     string        `json:"reportObservations" bson:"reportObservations"`
	Location               IssueLocation `json:"location" bson:"location"`
	VisitInfo              *Visit        `json:"visitInfo,omitempty" bson:"visitInfo,omitempty"`
	Issuer                 User          `json:"issuer" bson:"issuer"`
	WorkOrderInfo          *WorkOrder    `json:"workOrderInfo,omitempty" bson:"workOrderInfo,omitempty"`
	EstimatedExecutionTime time.Duration `json:"estimatedExecutionTime" bson:"estimatedExecutionTime"`
	Importance             string        `json:"importance" bson:"importance"` // TODO: Add constants to different severities
	Reports                []Report      `json:"reports" bson:"reports"`
	IssueNumber            int64         `json:"issueNumber" bson:"issueNumber"` //TODO: This number must be autogen.
	IssueComment           string        `json:"issueComment" bson:"issueComment"`
	Status                 string        `json:"status" bson:"status"` // TODO: LastStatutsChange Date
	SendToCooperative      bool          `json:"sendToCooperative" bson:"sendToCooperative"`
}

// IssueLocation holds the issue location information
type IssueLocation struct {
	PlantationID   bson.ObjectId    `json:"plantationId" bson:"plantationId"`
	PlotID         bson.ObjectId    `json:"plotId" bson:"plotId"`
	PlantationName string           `json:"plantationName" bson:"plantationName"`
	Coordinates    Coordinates      `json:"coordinates" bson:"coordinates"`
	PlotCode       PlotCodification `json:"plotCode" bson:"plotClode"`
}

// Visit holds the visit information
type Visit struct {
	VisitDate            time.Time                 `json:"visitDate" bson:"visitDate"`
	TechnicianID         *bson.ObjectId            `json:"technicianId" bson:"technicianId"`
	Technician           User                      `json:"technician" bson:"technician"`
	NeedsWorkOrder       bool                      `json:"needsWorkOrder" bson:"needsWorkOrder"`
	Images               []File                    `json:"images,omitempty" bson:"images,omitempty"`
	TemplateID           *bson.ObjectId            `json:"templateId,omitempty" bson:"templateId,omitempty"` // TODO: TemplateId and Data will be a new Object and Visit will rreceive an array of them
	Data                 map[string][]SamplingItem `json:"data,omitempty" bson:"data,omitempty"`
	VisitComment         string                    `json:"visitComment,omitempty" bson:"visitComment,omitempty"`
	CommentForTechnician string                    `json:"commentForTechnician,omitempty" bson:"commentForTechnician,omitempty"`
}

// WorkOrder holds the work order information
type WorkOrder struct {
	WorkDate             *time.Time         `json:"workDate,omitempty" bson:"workDate,omitempty"`
	WorkerID             *bson.ObjectId     `json:"workerId,omitempty" bson:"workerId,omitempty"`
	Worker               User               `json:"worker,omitempty" bson:"worker,omitempty"`
	Treatment            string             `json:"treatment" bson:"treatment"`
	TreatmentDescription map[string]string  `json:"treatmentDescription" bson:"treatmentDescription"`
	Comments             string             `json:"comments" bson:"comments"`
	WorkOrderTimer       *WorkOrderTimer    `json:"workOrderTimer,omitempty" bson:"workOrderTimer,omitempty"`
	TotalCost            float64            `json:"totalCost,omitempty" bson:"totalCost,omitempty"`
	Accepted             bool               `json:"accepted" bson:"accepted"`
	Products             []WorkOrderProduct `json:"products" bson:"products"`
	Completed            bool               `json:"completed, bson"-"`
}

// WorkOrderTimer holds the work order timer information
type WorkOrderTimer struct {
	StartTime     time.Time     `json:"startTime" bson:"startTime"`
	TimesPaused   int64         `json:"timesPaused" bson:"timesPaused"`
	LastPauseTime time.Time     `json:"lastPauseTime,omitempty" bson:"lastPauseTime,omitempty"`
	TimeWorking   time.Duration `json:"timeWorking,omitempty" bson:"timeWorking,omitempty"`
	RestartTime   time.Time     `json:"restartTime,omitempty" bson:"restartTime,omitempty"`
}

type WorkOrderProduct struct {
	ProductID            *bson.ObjectId    `json:"productId,omitempty" bson:"productId,omitempty"`
	ProductToUse         Product           `json:"productToUse,omitempty" bson:"productToUse,omitempty"`
	TreatmentDescription map[string]string `json:"treatmentDescription" bson:"treatmentDescription"`
}

// CreateIssueByIssue creates a new issue
func CreateIssueByIssue(issue *Issue) (*bson.ObjectId, error) {
	issue.InitializeNewData()

	if err := DB.C("issues").Insert(issue); err != nil {
		return nil, err
	}

	return &issue.ID, nil
}

//GetIssueByID returns a issue by its given ID
func GetIssueByID(issueID bson.ObjectId) (*Issue, error) {
	var issue Issue
	if err := DB.C("issues").Find(bson.M{"_id": issueID}).One(&issue); err != nil {
		logrus.Errorf("Error in domain.GetIssueByID -> error: %v", err.Error())
		return nil, err
	}

	return &issue, nil
}

func UpdateIssueByIssue(issue Issue) error {
	issue.UpdatedAt = time.Now().UTC()
	if err := DB.C("issues").UpdateId(issue.ID, issue); err != nil {
		if err.Error() != "not found" {
			logrus.Errorf("Error in domain.UpdateIssueByIssue -> error: %v", err.Error())
		}
		return err
	}

	return nil
}

func GetAllIssues() ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "sendToCooperative": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetAllFinishedIssues() ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "status": constants.StatusFinished, "sendToCooperative": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetIssuesByUserID(userID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{
		"$and": []bson.M{
			bson.M{"active": true},
			bson.M{"$or": []bson.M{
				bson.M{"issuer._id": userID, "sendToCooperative": true},
				bson.M{"technician._id": userID, "sendToCooperative": true},
				bson.M{"worker._id": userID, "sendToCooperative": true},
			}},
		},
	}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

type GetMyIssuesBody struct {
	UserID       bson.ObjectId `bson:"userId" json:"userId"`
	PlantationID bson.ObjectId `bson:"plantationId" json:"plantationId"`
}

func GetMyIssuesByUserID(body GetMyIssuesBody) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "issuer._id": body.UserID, "location.plantationId": body.PlantationID, "sendToCooperative": false}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetAllIssuedIssuesByUserID(userID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "issuer._id": userID}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetAllMyIssuesByUserID(userID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "issuer._id": userID, "sendToCooperative": false}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetFinishedIssuesByUserID(userID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{
		"$and": []bson.M{
			bson.M{"active": true, "status": constants.StatusFinished},
			bson.M{"$or": []bson.M{
				bson.M{"issuer._id": userID},
				bson.M{"technician._id": userID},
				bson.M{"worker._id": userID},
			}},
		},
	}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetIssuesByPlantationID(plantationID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "location.plantationId": plantationID, "sendToCooperative": true, "status": bson.M{"$ne": constants.StatusFinished}}).All(&result); err != nil {
		return result, nil
	}

	return result, nil
}

func GetFinishedIssuesByPlantationID(plantationID bson.ObjectId) ([]Issue, error) {
	result := []Issue{}

	if err := DB.C("issues").Find(bson.M{"active": true, "location.plantationId": plantationID, "status": constants.StatusFinished}).All(&result); err != nil {
		return result, nil
	}

	return result, nil
}

func DeleteIssueByID(issueID bson.ObjectId) error {

	if err := DB.C("issues").UpdateId(issueID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}
