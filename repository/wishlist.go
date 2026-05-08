package repository

import (
	"errors"
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
)

func IsProductExistByID(productID uint) (*domain.Product, error) {
	var product domain.Product
	err := database.DB.Preload("Images").Preload("Category").First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func IsWishlistItemExists(userID uint, productID uint) (bool, error) {
	var wishlist domain.Wishlist
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&wishlist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AddWishlistItem(item *domain.Wishlist) error {
	return database.DB.Create(item).Error
}
func RemoveWishlistItem(userID uint, productID uint) error {
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.Wishlist{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func GetWishlistByUserID(userID uint) ([]domain.Wishlist, error) {
	var items []domain.Wishlist

	err := database.DB.Preload("Product").Preload("Product.Images").Preload("Product.Category").Where("user_id = ?", userID).Order("created_at DESC").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func RemoveProductFromWishlist(userID, productID uint) error {
	result := database.DB.
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&domain.Wishlist{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
