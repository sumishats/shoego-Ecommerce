package usecase

import (
	"errors"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
)

func AddToWishlist(userID uint, req models.AddToWishlistRequest) error {
	product, err := repository.IsProductExistByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if !product.IsListed || !product.Category.IsListed {
		return errors.New("product is unavailable")
	}

	exists, err := repository.IsWishlistItemExists(userID, req.ProductID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("product already in wishlist")
	}

	item := &domain.Wishlist{
		UserID:    userID,
		ProductID: req.ProductID,
	}

	return repository.AddWishlistItem(item)
}

func RemoveFromWishlist(userID uint, productID uint) error {
	return repository.RemoveWishlistItem(userID, productID)
}


func GetWishlist(userID uint) (*models.WishlistResponse, error) {
	items, err := repository.GetWishlistByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responseItems []models.WishlistProductResponse

	for _, item := range items {
		product := item.Product

		var images []string
		for _, img := range product.Images {
			images = append(images, img.ImageURL)
		}

		status := "available"
		if !product.IsListed || !product.Category.IsListed {
			status = "product unavailable"
		} else if product.Stock <= 0 {
			status = "out of stock"
		}

		responseItems = append(responseItems, models.WishlistProductResponse{
			ProductID:    product.ID,
			Name:         product.Name,
			Description:  product.Description,
			Price:        product.Price,
			Stock:        product.Stock,
			CategoryID:   product.CategoryID,
			CategoryName: product.Category.Name,
			IsListed:     product.IsListed,
			Status:       status,
			Images:       images,
			IsWishlisted: true,
		})
	}

	return &models.WishlistResponse{
		Items: responseItems,
	}, nil
}