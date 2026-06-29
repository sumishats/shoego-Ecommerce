package repository

import (
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
)

func CreateCoupon(coupon *domain.Coupon) error {
	return database.DB.Create(coupon).Error
}
func GetCouponByCode(code string) (*domain.Coupon, error) {
	var coupon domain.Coupon

	err := database.DB.Where("code=?", code).First(&coupon).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &coupon, nil
}
func GetAllCoupons(page, limit int) ([]domain.Coupon, int64, error) {

	var coupons []domain.Coupon
	var totalCount int64

	db := database.DB.Model(&domain.Coupon{})

	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := db.Offset(offset).Limit(limit).Find(&coupons).Error; err != nil {

		return nil, 0, err
	}

	return coupons, totalCount, nil
}

func GetCouponByID(id uint) (*domain.Coupon, error) {

	var coupon domain.Coupon

	err := database.DB.First(&coupon, id).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func GetActiveCoupons() ([]domain.Coupon, error) {
	var coupons []domain.Coupon

	err := database.DB.Where("is_active = ?", true).Find(&coupons).Error

	return coupons, err
}

func UpdateCoupon(coupon *domain.Coupon) error {
	return database.DB.Save(coupon).Error
}

func DeleteCoupon(id uint) error {
	return database.DB.Delete(&domain.Coupon{}, id).Error
}


func UpdateCartCoupon(userID uint, code string) error {

	return database.DB.Model(&domain.Cart{}).
		Where("user_id = ?", userID).
		Update("coupon_code", code).Error
}
