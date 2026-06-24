package usecase

import (
	"errors"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"time"
)

func CreateCoupon(req models.CreateCouponRequest) error {

	existing, _ := repository.GetCouponByCode(req.Code)

	if existing != nil {
		return errors.New("coupon already exists")
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return errors.New("invalid expiry date")
	}

	if expiryDate.Before(time.Now()) {
		return errors.New("expiry date must be in future")
	}

	coupon := &domain.Coupon{
		Code:           req.Code,
		DiscountAmount: req.DiscountAmount,
		MinimumAmount:  req.MinimumAmount,
		ExpiryDate:     expiryDate,
		UsageLimit:     req.UsageLimit,
		IsActive:       true,
	}

	return repository.CreateCoupon(coupon)
}
func GetAllCoupons() ([]models.CouponResponse, error) {

	coupons, err := repository.GetAllCoupons()
	if err != nil {
		return nil, err
	}

	var response []models.CouponResponse

	for _, coupon := range coupons {

		response = append(response, models.CouponResponse{
			ID:             coupon.ID,
			Code:           coupon.Code,
			DiscountAmount: coupon.DiscountAmount,
			MinimumAmount:  coupon.MinimumAmount,
			ExpiryDate:     coupon.ExpiryDate.Format("2006-01-02"),
			UsageLimit:     coupon.UsageLimit,
			UsedCount:      coupon.UsedCount,
			IsActive:       coupon.IsActive,
		})
	}

	return response, nil
}
func DeleteCoupon(id uint) error {

	_, err := repository.GetCouponByID(id)
	if err != nil {
		return errors.New("coupon not found")
	}

	return repository.DeleteCoupon(id)
}

//user coupon

func ApplyCoupon(userID uint, code string) (*models.CouponApplyResponse, error) {

	cart, err := repository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	if cart.CouponCode != "" {
		return nil, errors.New("coupon already applied")
	}

	subtotal, err := repository.GetCartSubtotal(userID)
	if err != nil {
		return nil, err
	}

	coupon, err := repository.GetCouponByCode(code)
	if err != nil || coupon == nil {
		return nil, errors.New("invalid coupon")
	}

	if !coupon.IsActive {
		return nil, errors.New("coupon inactive")
	}

	if time.Now().After(coupon.ExpiryDate) {
		return nil, errors.New("coupon expired")
	}

	if coupon.UsedCount >= coupon.UsageLimit {
		return nil, errors.New("coupon limit reached")
	}

	if subtotal < coupon.MinimumAmount {
		return nil, errors.New("minimum amount not met")
	}

	discount := coupon.DiscountAmount

	if discount > subtotal {
		discount = subtotal
	}

	final := subtotal - discount
	err = repository.UpdateCartCoupon(
		userID,
		coupon.Code,
		discount,
	)
	if err != nil {
		return nil, err
	}

	return &models.CouponApplyResponse{
		Code:     coupon.Code,
		Subtotal: subtotal,
		Discount: discount,
		Final:    final,
	}, nil
}

func RemoveCoupon(userID uint) error {

	cart, err := repository.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	cart.CouponCode = ""
	cart.DiscountAmount = 0

	return repository.UpdateCart(cart)
}
