package services

import (
	"com/merkinsio/oasis-api/domain"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

//StorePlot creates and returns the newly created plot
func StorePlot(plot *domain.Plot) (*domain.Plot, error) {
	plotNumber := strconv.Itoa(plot.PlotCode.ProvinceRec) + strconv.Itoa(plot.PlotCode.TownRec) + strconv.Itoa(plot.PlotCode.PolygonNum) + strconv.Itoa(plot.PlotCode.PlotNum) + strconv.Itoa(plot.PlotCode.Enclosure)
	plot.PlotNumber = plotNumber
	plotID, err := domain.CreatePlotByPlot(plot)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetPlotByID(*plotID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetPlotsWithGeoJSONByPlantationID(plantationID bson.ObjectId) ([]domain.Plot, error) {
	plots, err := domain.GetPlotsByPlantationID(plantationID)
	if err != nil {
		return nil, err
	}

	for i, plot := range plots {
		sigpacPlot, err := domain.GetSigpacPlotByPlotCode(plot.PlotCode)
		if err != nil {
			return nil, err
		} else {
			plots[i].GeoJSON = &sigpacPlot.GeoJSON
		}
	}

	return plots, nil
}
