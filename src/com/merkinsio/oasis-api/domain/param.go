package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type Param struct {
	Entity                 `bson:",inline"`
	Name                   string `bson:"name" json:"name"`
	TreatmentCultivation   bool   `bson:"treatmentCultivation" json:"treatmentCultivation"`
	TreatmentSowing        bool   `bson:"treatmentSowing" json:"treatmentSowing"`
	TreatmentPlant         bool   `bson:"treatmentPlant" json:"treatmentPlant"`
	TreatmentFertilization bool   `bson:"treatmentFertilization" json:"treatmentFertilization"`
	TreatmentPruning       bool   `bson:"treatmentPruning" json:"treatmentPruning"`
	TreatmentPhytosanitary bool   `bson:"treatmentPhytosanitary" json:"treatmentPhytosanitary"`
	TreatmentIrrigation    bool   `bson:"treatmentIrrigation" json:"treatmentIrrigation"`
	TreatmentTwigsRemoval  bool   `bson:"treatmentTwigsRemoval" json:"treatmentTwigsRemoval"`
	TreatmentHarvesting    bool   `bson:"treatmentHarvesting" json:"treatmentHarvesting"`
	TreatmentAnalysis      bool   `bson:"treatmentAnalysis" json:"treatmentAnalysis"`
	TreatmentOtherTask     bool   `bson:"treatmentOtherTask" json:"treatmentOtherTask"`
}

//Validate validates the Param struct
func (param *Param) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(param)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreateParamByParam creates a new param with the param values brought
func CreateParamByParam(param *Param) (*bson.ObjectId, error) {
	param.InitializeNewData()

	if _, err := param.Validate(); err != nil {
		log.Errorf("Error in domain.CreateParamByParam.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("params").Insert(param); err != nil {
		log.Errorf("Error in domain.CreateParamByParam.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &param.ID, nil
}

func GetParams() ([]Param, error) {
	result := []Param{}

	if err := DB.C("params").Find(bson.M{"active": true}).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetParamsByType(treatmentType string) ([]Param, error) {
	result := []Param{}
	query := bson.M{}
	query[treatmentType] = true
	query["active"] = true

	if err := DB.C("params").Find(query).All(&result); err != nil {
		return result, err
	}

	return result, nil
}

func GetParamByID(paramID bson.ObjectId) (*Param, error) {
	var param Param
	if err := DB.C("params").Find(bson.M{"_id": paramID}).One(&param); err != nil {
		log.Errorf("Error in domain.GetParamByID -> error: %v", err.Error())
		return nil, err
	}

	return &param, nil
}

func DeleteParamByID(paramID bson.ObjectId) error {

	if err := DB.C("params").UpdateId(paramID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdateParamByParam(param Param) error {
	param.UpdatedAt = time.Now().UTC()
	if err := DB.C("params").UpdateId(param.ID, param); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdateParamByID -> error: %v", err.Error())
		}
		return err
	}

	return nil
}
