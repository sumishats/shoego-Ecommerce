package usecase

import (
	"shoego/models"
	"shoego/repository"
)

func GetAdminBestSellingProducts() ([]models.BestSellingProductResponse, error) {

	products, err := repository.GetAdminBestSellingProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}
