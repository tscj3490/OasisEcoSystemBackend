package domain

import (
	"com/merkinsio/oasis-api/constants"

	"gopkg.in/mgo.v2/bson"
)

type Statistics struct {
	FinishedIssues int64   `bson:"finishedIssues" json:"finishedIssues"`
	TotalCost      float64 `bson:"totalCost" json:"totalCost"`
}

type PlantationStatistics struct {
	Statistics `bson:",inline"`
}

func GetPlantationStatisticsByPlantationID(plantationID bson.ObjectId) (PlantationStatistics, error) {
	result := PlantationStatistics{}
	if err := DB.C("issues").Pipe([]bson.M{
		bson.M{"$match": bson.M{"status": constants.StatusFinished, "active": true,
			"location.plantationId": plantationID}},
		bson.M{"$group": bson.M{"_id": nil, "finishedIssues": bson.M{"$sum": 1}, "totalCost": bson.M{"$sum": "$workOrderInfo.totalCost"}}},
	}).One(&result); err != nil {
		return result, err
	}
	return result, nil
}
