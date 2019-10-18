package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type VisitCard struct {
	Entity `bson:",inline"`

	Name          string        `bson:"name" json:"name"`
	Samplings     []Sampling    `json:"samplings" bson:"samplings"`
	CooperativeID bson.ObjectId `json:"cooperativeId" bson:"cooperativeId"`
}

func (visitCard *VisitCard) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(visitCard)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

func GetVisitCardsByCooperativeID(cooperativeID bson.ObjectId) ([]VisitCard, error) {
	result := []VisitCard{}

	if err := DB.C("visitCards").Find(bson.M{"active": true, "cooperativeId": cooperativeID}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetVisitCardByID(visitCardID bson.ObjectId) (*VisitCard, error) {
	var result VisitCard

	if err := DB.C("visitCards").Find(bson.M{"active": true, "_id": visitCardID}).One(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateVisitCardByVisitCard(visitCard *VisitCard) (*bson.ObjectId, error) {
	visitCard.InitializeNewData()

	if _, err := visitCard.Validate(); err != nil {
		log.Errorf("Error in domain.CreateVisitCardByVisitCard.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("visitCards").Insert(visitCard); err != nil {
		log.Errorf("Error in domain.CreateVisitCardByVisitCard.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &visitCard.ID, nil
}

func DeleteVisitCardByID(visitCardID bson.ObjectId) error {

	if err := DB.C("visitCards").UpdateId(visitCardID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdateVisitCardByVisitCard(visitCard VisitCard) error {
	visitCard.UpdatedAt = time.Now().UTC()
	if err := DB.C("visitCards").UpdateId(visitCard.ID, visitCard); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateVisitCardByID -> error: %v", err.Error())
		}
		return err
	}

	return nil
}
