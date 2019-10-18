package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Cooperative holds the cooperative data
type Cooperative struct {
	Entity    `bson:",inline"`
	Name      string         `json:"name" bson:"name"`
	Address   string         `json:"address" bson:"address"`
	TaxIDCode string         `json:"taxIdCode" bson:"taxIdCode"` // Spanish C.I.F.
	Location  Coordinates    `json:"location" bson:"location"`
	AdminID   *bson.ObjectId `json:"adminId" bson:"adminId"`
}

//Validate validates the Cooperative struct
func (cooperative *Cooperative) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(cooperative)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

//GetCooperativeByID returns a cooperative by its given ID
func GetCooperativeByID(cooperativeID bson.ObjectId) (*Cooperative, error) {
	var cooperative Cooperative
	if err := DB.C("cooperatives").Find(bson.M{"_id": cooperativeID}).One(&cooperative); err != nil {
		log.Errorf("Error in domain.GetCooperativeByID -> error: %v", err.Error())
		return nil, err
	}

	return &cooperative, nil
}

// CreateCooperativeByCooperative creates a new cooperative with the cooperative values brought
func CreateCooperativeByCooperative(cooperative *Cooperative) (*bson.ObjectId, error) {
	cooperative.InitializeNewData()

	if _, err := cooperative.Validate(); err != nil {
		log.Errorf("Error in domain.CreateCooperativeByCooperative.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("cooperatives").Insert(cooperative); err != nil {
		log.Errorf("Error in domain.CreateCooperativeByCooperative.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &cooperative.ID, nil
}

func GetCooperatives() ([]Cooperative, error) {
	result := []Cooperative{}

	if err := DB.C("cooperatives").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func DeleteCooperativeByID(cooperativeID bson.ObjectId) error {

	if err := DB.C("cooperatives").UpdateId(cooperativeID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdateCooperativeByCooperative(cooperative Cooperative) error {
	cooperative.UpdatedAt = time.Now().UTC()
	if err := DB.C("cooperatives").UpdateId(cooperative.ID, cooperative); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateCooperativeByID -> error: %v", err.Error())
		}
		return err
	}

	return nil
}

// SetCooperativeOwner Set a cooperative owner (by changing user info)
func SetCooperativeOwner(userOwner UserAndOwner) error {
	var cooperative *Cooperative
	if err := DB.C("cooperatives").Find(bson.M{"adminId": userOwner.User.ID}).One(&cooperative); err != nil {
		log.Errorf("Error in domain.GetCooperativeByID.Find (SetCooperativeOwner) -> error: %v", err.Error())
	}
	if cooperative == nil {
	} else {
		if err := DB.C("cooperatives").UpdateId(cooperative.ID, bson.M{"$set": bson.M{"adminId": nil}}); err != nil {
			return err
		}
	}

	if userOwner.Owner {
		if err := DB.C("cooperatives").UpdateId(userOwner.User.CooperativeID, bson.M{"$set": bson.M{"adminId": userOwner.User.ID}}); err != nil {
			log.Errorf("Error in domain.SetCooperativeOwner.UpdateId -> error: %s", err.Error())
			return err
		}
	}
	return nil
}

func GetCooperativeWorkersByID(cooperativeID bson.ObjectId) ([]User, error) {
	result := []User{}
	if err := DB.C("users").Find(bson.M{
		"$and": []bson.M{
			bson.M{"cooperativeId": cooperativeID},
			bson.M{"role": "ROLE_WORKER"},
		}}).All(&result); err != nil {
		log.Errorf("Error in domain.SetCooperativeOwner.UpdateId -> error: %s", err.Error())
		return nil, err
	}
	return result, nil
}
