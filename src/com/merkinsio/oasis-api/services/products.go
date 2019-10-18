package services

import (
	"com/merkinsio/oasis-api/domain"
)

func StoreProduct(product *domain.Product) (*domain.Product, error) {
	productID, err := domain.CreateProductByProduct(product)
	if err != nil {
		return nil, err
	}

	result, err := domain.GetProductByID(*productID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
