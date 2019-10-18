package services

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/domain"
	"errors"
)

func GetParamsByTreatmentType(treatmentType string) ([]domain.Param, error) {
	params := []domain.Param{}
	var err error
	switch treatmentType {
	case constants.TreatmentCultivation:
		params, err = domain.GetParamsByType("treatmentCultivation")
	case constants.TreatmentSowing:
		params, err = domain.GetParamsByType("treatmentSowing")
	case constants.TreatmentPlant:
		params, err = domain.GetParamsByType("treatmentPlant")
	case constants.TreatmentFertilization:
		params, err = domain.GetParamsByType("treatmentFertilization")
	case constants.TreatmentPruning:
		params, err = domain.GetParamsByType("treatmentPruning")
	case constants.TreatmentPhytosanitary:
		params, err = domain.GetParamsByType("treatmentPhytosanitary")
	case constants.TreatmentIrrigation:
		params, err = domain.GetParamsByType("treatmentIrrigation")
	case constants.TreatmentTwigsRemoval:
		params, err = domain.GetParamsByType("treatmentTwigsRemoval")
	case constants.TreatmentHarvesting:
		params, err = domain.GetParamsByType("treatmentHarvesting")
	case constants.TreatmentAnalysis:
		params, err = domain.GetParamsByType("treatmentAnalysis")
	case constants.TreatmentOtherTask:
		params, err = domain.GetParamsByType("treatmentOtherTask")
	default:
		err = errors.New("invalid treatmentType")
	}
	if err != nil {
		return nil, err
	}

	return params, nil
}

func StoreParam(param *domain.Param) (*domain.Param, error) {
	paramID, err := domain.CreateParamByParam(param)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetParamByID(*paramID)

	if err != nil {
		return nil, err
	}

	return res, nil
}
