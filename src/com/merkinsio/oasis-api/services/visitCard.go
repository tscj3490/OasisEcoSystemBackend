package services

import (
	"com/merkinsio/oasis-api/domain"
)

func StoreVisitCard(visitCard *domain.VisitCard) (*domain.VisitCard, error) {
	visitCardID, err := domain.CreateVisitCardByVisitCard(visitCard)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetVisitCardByID(*visitCardID)

	if err != nil {
		return nil, err
	}

	return res, nil
}
