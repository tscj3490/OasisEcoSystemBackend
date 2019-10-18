package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type Variaty struct {
	Entity              `bson:",inline"`
	CodeVariaty         int64  `bson:"codeVariaty" json:"codeVariaty"`
	Name                string `bson:"name" json:"name"`
	DeclaredCultivation string `bson:"declaredCultivation" json:"declaredCultivation"`
}

//Validate validates the Param struct
func (variaty *Variaty) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(variaty)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreateParamByParam creates a new param with the param values brought
func CreateVariatyByVariaty(variaty *Variaty) (*bson.ObjectId, error) {
	variaty.InitializeNewData()

	if _, err := variaty.Validate(); err != nil {
		log.Errorf("Error in domain.CreateVariatyByVariaty.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("variatyCrops").Insert(variaty); err != nil {
		log.Errorf("Error in domain.CreateVariatyByVariaty.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &variaty.ID, nil
}

func GetVariaty() ([]Variaty, error) {
	result := []Variaty{}

	if err := DB.C("variatyCrops").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetVariatyByID(variatyID bson.ObjectId) (*Variaty, error) {
	var variaty Variaty
	if err := DB.C("variatyCrops").Find(bson.M{"_id": variatyID}).One(&variaty); err != nil {
		log.Errorf("Error in domain.GetVariatyByID -> error: %v", err.Error())
		return nil, err
	}

	return &variaty, nil
}

func DeleteVariatyByID(variatyID bson.ObjectId) error {

	if err := DB.C("variatyCrops").UpdateId(variatyID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdateVariatyByVariaty(variaty Variaty) error {
	variaty.UpdatedAt = time.Now().UTC()
	if err := DB.C("variatyCrops").UpdateId(variaty.ID, variaty); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateVariatyByVariaty -> error: %v", err.Error())
		}
		return err
	}

	return nil
}

func GetConditionalVariaties(declaredCultivation string) ([]Variaty, error) {
	result := []Variaty{}

	if err := DB.C("variatyCrops").Find(bson.M{"active": true, "declaredCultivation": declaredCultivation}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}
