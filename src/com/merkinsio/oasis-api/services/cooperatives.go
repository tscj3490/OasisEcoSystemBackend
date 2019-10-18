package services

import "com/merkinsio/oasis-api/domain"

func StoreCooperative(cooperative *domain.Cooperative) (*domain.Cooperative, error) {
	cooperativeID, err := domain.CreateCooperativeByCooperative(cooperative)

	if err != nil {
		return nil, err
	}

	res, err := domain.GetCooperativeByID(*cooperativeID)

	if err != nil {
		return nil, err
	}

	return res, nil
}
