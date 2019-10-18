package services

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//StoreIssue creates and returns the newly created issue
func StoreIssue(issue *domain.Issue) (*domain.Issue, error) {
	issueID, err := domain.CreateIssueByIssue(issue)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetIssueByID(*issueID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateIssue(issue *domain.Issue, currentUser domain.User) (*domain.Issue, error) {
	now := time.Now().UTC()
	if issue.VisitInfo != nil {
		if issue.VisitInfo.TechnicianID != nil {
			technician, err := domain.GetUserByID(*issue.VisitInfo.TechnicianID)
			if err != nil {
				return nil, err
			}
			issue.VisitInfo.Technician = *technician
		} else {
			issue.VisitInfo.Technician = currentUser
			issue.VisitInfo.TechnicianID = &currentUser.ID
		}
		if !issue.VisitInfo.VisitDate.IsZero() {
			issue.VisitInfo.VisitDate = now
		}
		issue.Status = constants.StatusPendingVisit
		if issue.VisitInfo.VisitComment != "" || issue.VisitInfo.TemplateID != nil { // This changes the state depending on what comes
			if issue.VisitInfo.NeedsWorkOrder {
				issue.Status = constants.StatusPendingWorkOrder
			} else {
				issue.Status = constants.StatusFinished
			}
		}
	}

	if issue.WorkOrderInfo != nil {
		issue.VisitInfo.NeedsWorkOrder = true
		if issue.WorkOrderInfo.Accepted {
			issue.Status = constants.StatusAccepted
		} else {
			issue.Status = constants.StatusWorkOrder
			if issue.WorkOrderInfo.Products != nil && len(issue.WorkOrderInfo.Products) > 0 {
				for i, product := range issue.WorkOrderInfo.Products {
					product, err := domain.GetProductByID(*product.ProductID)
					if err != nil {
						return nil, err
					}
					issue.WorkOrderInfo.Products[i].ProductToUse = *product
				}
			}
		}
	}

	if issue.WorkOrderInfo != nil && !issue.WorkOrderInfo.WorkDate.IsZero() {
		if issue.WorkOrderInfo.WorkerID != nil {
			worker, err := domain.GetUserByID(*issue.WorkOrderInfo.WorkerID)
			if err != nil {
				return nil, err
			}
			issue.WorkOrderInfo.Worker = *worker
		} else {
			issue.WorkOrderInfo.Worker = currentUser
			issue.WorkOrderInfo.WorkerID = &currentUser.ID
		}
		issue.VisitInfo.NeedsWorkOrder = true
		issue.WorkOrderInfo.Accepted = true
		issue.Status = constants.StatusAssigned
	}
	err := domain.UpdateIssueByIssue(*issue)

	if err != nil {
		return nil, err
	}
	res, err := domain.GetIssueByID(issue.ID)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func StartWorkingIssue(issue *domain.Issue) (*domain.Issue, error) {
	issue.Status = constants.StatusWorking
	if issue.WorkOrderInfo.WorkOrderTimer != nil && issue.WorkOrderInfo.WorkOrderTimer.TimesPaused > 0 {
		issue.WorkOrderInfo.WorkOrderTimer.RestartTime = time.Now().UTC()
	} else {
		workOrderTimer := domain.WorkOrderTimer{
			StartTime: time.Now().UTC(),
		}
		issue.WorkOrderInfo.WorkOrderTimer = &workOrderTimer
	}
	err := domain.UpdateIssueByIssue(*issue)

	if err != nil {
		return nil, err
	}
	res, err := domain.GetIssueByID(issue.ID)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func PauseWorkingIssue(issue *domain.Issue) (*domain.Issue, error) {
	issue.Status = constants.StatusPaused
	pauseTime := time.Now().UTC()
	var currentTimeWorked time.Duration

	if issue.WorkOrderInfo.WorkOrderTimer.TimesPaused > 0 {
		issue.WorkOrderInfo.WorkOrderTimer.TimesPaused++
		timeWorked := issue.WorkOrderInfo.WorkOrderTimer.TimeWorking
		timeSinceRestart := pauseTime.Sub(issue.WorkOrderInfo.WorkOrderTimer.RestartTime)
		currentTimeWorked = timeWorked + timeSinceRestart
	} else {
		issue.WorkOrderInfo.WorkOrderTimer.TimesPaused = 1
		currentTimeWorked = pauseTime.Sub(issue.WorkOrderInfo.WorkOrderTimer.StartTime)
	}
	issue.WorkOrderInfo.WorkOrderTimer.TimeWorking = currentTimeWorked
	issue.WorkOrderInfo.WorkOrderTimer.LastPauseTime = time.Now().UTC()

	err := domain.UpdateIssueByIssue(*issue)

	if err != nil {
		return nil, err
	}
	res, err := domain.GetIssueByID(issue.ID)
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func GetOwnedPlantationIssuesByUserID(userID bson.ObjectId) ([]domain.Issue, error) {
	result := []domain.Issue{}
	plantations, err := domain.GetPlantationsByOwnerID(userID)
	if err != nil {
		return nil, err
	}
	for _, plantation := range plantations {
		issues, err := domain.GetIssuesByPlantationID(plantation.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, issues...)
	}
	return result, nil
}
