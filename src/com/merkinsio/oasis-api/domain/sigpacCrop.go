package domain

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type SigpacCrop struct {
	Entity `bson:",inline"`

	Code              string `json:"code" bson:"code"`
	SigpacDescription string `json:"sigpacDescription" bson:"sigpacDescription"`
}

func GetSigpacCropByID(sigpacCropID bson.ObjectId) (*SigpacCrop, error) {
	var sigpacCrop SigpacCrop
	if err := DB.C("sigpacCrops").Find(bson.M{"_id": sigpacCropID, "active": true}).One(&sigpacCrop); err != nil {
		log.Errorf("Error in domain.GetSigpacCropByID -> error: %v", err.Error())
		return nil, err
	}

	return &sigpacCrop, nil
}

func GetAllSigpacCrops() ([]SigpacCrop, error) {
	result := []SigpacCrop{}

	if err := DB.C("sigpacCrops").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}
