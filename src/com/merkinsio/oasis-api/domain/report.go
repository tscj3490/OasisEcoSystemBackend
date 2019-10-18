package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Report struct {
	Entity           `bson:",inline"`
	WorkerID         bson.ObjectId `json:"workerId" bson:"workerId"`
	ConsumedTime     time.Duration `json:"consumedTime" bson:"consumedTime"`
	ClientAssessment string        `json:"clientAssessment" bson:"clientAssessment"` //NOTE: We may need
	WorkerComment    string        `json:"workerComment" bson:"workerComment"`
	ClientSignature  File          `json:"clientSignature" bson:"clientSignature"`
	IssueID          bson.ObjectId `json:"issueId" bson:"issueId"`
}
