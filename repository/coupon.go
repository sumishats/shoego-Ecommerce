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

func GetAllCoupons() ([]domain.Coupon, error) {
	var coupons []domain.Coupon

	err := database.DB.Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}
func GetCouponByID(id uint) (*domain.Coupon, error) {

	var coupon domain.Coupon

	err := database.DB.
		First(&coupon, id).
		Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}
func DeleteCoupon(id uint) error {
	return database.DB.Delete(&domain.Coupon{}, id).Error
}
func UpdateCartCoupon(userID uint, code string, discount float64) error {

	return database.DB.
		Model(&domain.Cart{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"coupon_code":     code,
			"discount_amount": discount,
		}).Error
}