package repository

import (
	"shoego/database"
	"shoego/domain"
	"time"
)

func CreateProductOffer(offer *domain.ProductOffer) error {
	return database.DB.Create(offer).Error
}

func GetProductOfferByProductID(productID uint) (*domain.ProductOffer, error) {

	var offer domain.ProductOffer

	err := database.DB.
		Where("product_id = ? AND is_active = ?", productID, true).
		First(&offer).Error

	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func GetAllProductOffers() ([]domain.ProductOffer, error) {

	var offers []domain.ProductOffer

	err := database.DB.
		Preload("Product").
		Find(&offers).Error

	return offers, err
}
func DeleteProductOffer(offerID uint) error {

	return database.DB.
		Delete(&domain.ProductOffer{}, offerID).
		Error
}

func GetActiveProductOffer(productID uint) (*domain.ProductOffer, error) {

	var offer domain.ProductOffer

	err := database.DB.
		Where(
			"product_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ?",
			productID,
			true,
			time.Now(),
			time.Now(),
		).
		First(&offer).Error

	if err != nil {
		return nil, err
	}

	return &offer, nil
}
func CreateCategoryOffer(offer *domain.CategoryOffer) error {
	return database.DB.Create(offer).Error
}
func DeleteCategoryOffer(id uint) error {
	return database.DB.Delete(&domain.CategoryOffer{}, id).Error
}

func GetAllCategoryOffers() ([]domain.CategoryOffer, error) {

	var offers []domain.CategoryOffer

	err := database.DB.Preload("Category").Find(&offers).Error

	return offers, err
}
func GetActiveCategoryOffer(categoryID uint) (*domain.CategoryOffer, error) {

	var offer domain.CategoryOffer

	err := database.DB.
		Where(
			"category_id = ? AND is_active = ? AND start_date <= NOW() AND end_date >= NOW()",
			categoryID,
			true,
		).
		First(&offer).Error

	if err != nil {
		return nil, err
	}

	return &offer, nil
}

// func GetCategoryByID(id uint) (*domain.Category, error) {
// 	var category domain.Category

// 	err := database.DB.First(&category, id).Error

// 	if err != nil {
// 		return nil, err
// 	}

//		return &category, nil/
//	}
func GetUserByReferralCode(code string) (*domain.User, error) {

	var user domain.User

	err := database.DB.
		Where("referral_code = ?", code).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}