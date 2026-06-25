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

func GetAllProductOffers(page, limit int) ([]domain.ProductOffer, int64, error) {

	var offers []domain.ProductOffer
	var totalCount int64

	offset := (page - 1) * limit

	database.DB.Model(&domain.ProductOffer{}).Count(&totalCount)

	err := database.DB.Preload("Product").Limit(limit).Offset(offset).Find(&offers).Error

	return offers, totalCount, err
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

func GetAllCategoryOffers(page, limit int) ([]domain.CategoryOffer, int64, error) {

	var offers []domain.CategoryOffer
	var totalCount int64

	offset := (page - 1) * limit

	// count total records
	err := database.DB.Model(&domain.CategoryOffer{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// paginated fetch
	err = database.DB.Preload("Category").Limit(limit).Offset(offset).Find(&offers).Error

	return offers, totalCount, err
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
func GetProductOfferByID(id uint) (*domain.ProductOffer, error) {

	var offer domain.ProductOffer

	err := database.DB.First(&offer, id).Error

	if err != nil {
		return nil, err
	}

	return &offer, nil
}
func UpdateProductOffer(offer *domain.ProductOffer) error {
	return database.DB.Save(offer).Error
}
func GetCategoryOfferByID(id uint) (*domain.CategoryOffer, error) {

	var offer domain.CategoryOffer

	err := database.DB.First(&offer, id).Error
	if err != nil {
		return nil, err
	}

	return &offer, nil
}
func UpdateCategoryOffer(offer *domain.CategoryOffer) error {

	return database.DB.Save(offer).Error
}
