package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

//Plantation Holds the plantation data
type Plantation struct {
	Entity      `bson:",inline"`
	SigpacID    string      `json:"sigpacId" bson:"sigpacId"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
	// OverseerID  bson.ObjectId `json:"overseerId" bson:"overseerId"`
	OwnerID       []bson.ObjectId `json:"ownerId" bson:"ownerId"` //Actually this is the ManagerId
	Name          string          `json:"name" bson:"name"`
	CooperativeID bson.ObjectId   `json:"cooperativeId" bson:"cooperativeId"`
	Plots         *[]Plot         `json:"plots,omitempty" bson:"plots,omitempty"`
	OwnerPlot     string          `json:"ownerPlot,omitempty" bson:"ownerPlot,omitempty"`
}

//Validate validates the Plantation struct
func (plantation *Plantation) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(plantation)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreatePlantationByPlantation creates a new plantation
func CreatePlantationByPlantation(plantation *Plantation) (*bson.ObjectId, error) {
	plantation.InitializeNewData()

	if _, err := plantation.Validate(); err != nil {
		log.Errorf("Error in domain.CreatePlantationByPlantation.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("plantations").Insert(plantation); err != nil {
		log.Errorf("Error in domain.CreatePlantationByPlantation.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &plantation.ID, nil
}

//GetPlantationByID returns a plantation by its given ID
func GetPlantationByID(plantationID bson.ObjectId) (*Plantation, error) {
	var plantation Plantation
	if err := DB.C("plantations").Find(bson.M{"_id": plantationID}).One(&plantation); err != nil {
		log.Errorf("Error in domain.GetPlantationByID -> error: %v", err.Error())
		return nil, err
	}

	return &plantation, nil
}

//GetPlantationsByUserCooperativeID returns an Array pf plantations that have the same CooperativeID as the user's
func GetPlantationsByUserCooperativeID(cooperativeID bson.ObjectId) ([]Plantation, error) {
	plantations := []Plantation{}
	// TODO: ROLES diffs
	if err := DB.C("plantations").Find(bson.M{"cooperativeId": cooperativeID, "active": true}).All(&plantations); err != nil {
		return plantations, err
	}

	return plantations, nil
}

func GetPlantationsByOwnerID(ownerID bson.ObjectId) ([]Plantation, error) {
	plantations := []Plantation{}
	// TODO: ROLES diffs
	if err := DB.C("plantations").Find(bson.M{"ownerId": ownerID, "active": true}).All(&plantations); err != nil {
		return plantations, err
	}

	return plantations, nil
}

func GetPlantationsAndPlotsByCooperativeID(cooperativeID bson.ObjectId) ([]Plantation, error) {
	plantations := []Plantation{}
	// TODO: ROLES diffs
	if err := DB.C("plantations").Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"cooperativeId": cooperativeID,
				"active":        true,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "plots",
				"localField":   "_id",
				"foreignField": "plantationId",
				"as":           "_plots",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"plots": bson.M{
					"$filter": bson.M{
						"input": "$_plots",
						"as":    "plot",
						"cond":  bson.M{"$eq": []interface{}{"$$plot.active", true}},
					},
				},
			},
		},
		bson.M{"$project": bson.M{
			"plots.shape": 0,
			"_plots":      0,
		}},
	}).All(&plantations); err != nil {
		return plantations, err
	}

	return plantations, nil
}

func GetPlantationsAndPlotsByOwnerId(ownerID bson.ObjectId) ([]Plantation, error) {
	plantations := []Plantation{}
	// TODO: ROLES diffs
	if err := DB.C("plantations").Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"ownerId": ownerID,
				"active":  true,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "plots",
				"localField":   "_id",
				"foreignField": "plantationId",
				"as":           "plots",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"plots": bson.M{
					"$filter": bson.M{
						"input": "$_plots",
						"as":    "plot",
						"cond":  bson.M{"$eq": []interface{}{"$$plot.active", true}},
					},
				},
			},
		},
		bson.M{"$project": bson.M{
			"plots.shape": 0,
			"_plots":      0,
		}},
	}).All(&plantations); err != nil {
		return plantations, err
	}

	return plantations, nil
}

func DeletePlantationByID(plantationID bson.ObjectId) error {

	if err := DB.C("plantations").UpdateId(plantationID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdatePlantationByID(plantation Plantation) error {
	plantation.UpdatedAt = time.Now().UTC()
	if err := DB.C("plantations").UpdateId(plantation.ID, plantation); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdatePlantationByID -> error: %v", err.Error())
		}
		return err
	}

	return nil
}
