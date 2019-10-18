package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type DeclaredCultivation struct {
	Entity                  `bson:",inline"`
	CodeDeclaredCultivation int64  `bson:"codeDeclaredCultivation" json:"codeDeclaredCultivation"`
	Name                    string `bson:"name" json:"name"`
	Season                  string `json:"season" bson:"season"`
	ColorCode               string `json:"colorCode" bson:"colorCode"`
}

//Validate validates the Param struct
func (declaredCultivation *DeclaredCultivation) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(declaredCultivation)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreateParamByParam creates a new param with the param values brought
func CreateDeclaredCultivationByDeclaredCultivation(declaredCultivation *DeclaredCultivation) (*bson.ObjectId, error) {
	declaredCultivation.InitializeNewData()

	if _, err := declaredCultivation.Validate(); err != nil {
		log.Errorf("Error in domain.CreateDeclaredCultivationByDeclaredCultivation.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("declaredCultivationCrops").Insert(declaredCultivation); err != nil {
		log.Errorf("Error in domain.CreateDeclaredCultivationByDeclaredCultivation.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &declaredCultivation.ID, nil
}

func GetAllDeclaredCultivations() ([]DeclaredCultivation, error) {
	result := []DeclaredCultivation{}

	if err := DB.C("declaredCultivationCrops").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetDeclaredCultivationByID(declaredCultivationID bson.ObjectId) (*DeclaredCultivation, error) {
	var declaredCultivation DeclaredCultivation
	if err := DB.C("declaredCultivationCrops").Find(bson.M{"_id": declaredCultivationID}).One(&declaredCultivation); err != nil {
		log.Errorf("Error in domain.GetDeclaredCultivationByID -> error: %v", err.Error())
		return nil, err
	}

	return &declaredCultivation, nil
}

func DeleteDeclaredCultivationByID(declaredCultivationID bson.ObjectId) error {

	if err := DB.C("declaredCultivationCrops").UpdateId(declaredCultivationID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdateDeclaredCultivationByDeclaredCultivation(declaredCultivation DeclaredCultivation) error {
	declaredCultivation.UpdatedAt = time.Now().UTC()
	if err := DB.C("declaredCultivationCrops").UpdateId(declaredCultivation.ID, declaredCultivation); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateDeclaredCultivationByDeclaredCultivation -> error: %v", err.Error())
		}
		return err
	}

	return nil
}
