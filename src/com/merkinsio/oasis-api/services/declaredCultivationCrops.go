package services

import (
	"com/merkinsio/oasis-api/domain"
)

func StoreDeclaredCultivation(declaredCultivation *domain.DeclaredCultivation) (*domain.DeclaredCultivation, error) {
	declaredCultivationID, err := domain.CreateDeclaredCultivationByDeclaredCultivation(declaredCultivation)
	if err != nil {
		return nil, err
	}

	result, err := domain.GetDeclaredCultivationByID(*declaredCultivationID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
