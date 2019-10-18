package services

import (
	"com/merkinsio/oasis-api/domain"

	"gopkg.in/mgo.v2/bson"
)

//StorePlantation creates and returns the newly created plantation
func StorePlantation(plantation *domain.Plantation) (*domain.Plantation, error) {
	plantationID, err := domain.CreatePlantationByPlantation(plantation)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetPlantationByID(*plantationID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetSigpacPlotsByPlantationID(plantationID bson.ObjectId) ([]domain.SigpacPlot, error) {
	plantationPlots, err := domain.GetPlotsByPlantationID(plantationID)
	if err != nil {
		return nil, err
	}
	var sigpacPlots []domain.SigpacPlot
	for _, plot := range plantationPlots {
		sigpacPlot, err := domain.GetSigpacPlotByPlotCode(plot.PlotCode)
		if err != nil {
			return nil, err
		}
		sigpacPlots = append(sigpacPlots, *sigpacPlot)
	}
	return sigpacPlots, nil
}
