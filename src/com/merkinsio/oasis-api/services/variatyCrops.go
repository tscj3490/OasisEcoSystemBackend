package services

import (
	"com/merkinsio/oasis-api/domain"
)

func StoreVariaty(variaty *domain.Variaty) (*domain.Variaty, error) {
	variatyID, err := domain.CreateVariatyByVariaty(variaty)
	if err != nil {
		return nil, err
	}

	result, err := domain.GetVariatyByID(*variatyID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
