package domain

import (
	"time"
)

type Sampling struct {
	Name string `json:"name" bson:"name"`

	SampleDescription string `json:"sampleDescription" bson:"sampleDescription"`

	ThresholdDescription string `json:"thresholdDescription" bson:"thresholdDescription"`

	SeasonDescription string `json:"seasonDescription" bson:"seasonDescription"`

	Items []string `json:"items" bson:"items"`

	// Data []SamplingData `json:"data" bson:"data"`

	// PlotID bson.ObjectId `json:"plotId" bson:"plotId"`

	// WorkerID bson.ObjectId `json:"workerId" bson:"workerId"`

	// TODO: Maybe we'll have to add the crops reference here

}

type SamplingData struct {
	SamplingDate time.Time `json:"samplingDate" bson:"samplingDate"`

	SamplingItems []SamplingItem `json:"samplingItems" bson:"samplingItems"`
}

type SamplingItem struct {
	ItemName string `json:"itemName" bson:"itemName"`

	Threshold string `json:"threshold" bson:"threshold"`

	Treatment string `json:"treatment" bson:"treatment"`
}
