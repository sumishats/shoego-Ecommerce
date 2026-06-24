package usecase

import (
	"errors"
	"math"
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
	if req.DiscountType != "fixed" && req.DiscountType != "percentage" {
		return errors.New("invalid discount type")
	}

	coupon := &domain.Coupon{
		Code:           req.Code,
		DiscountType:   req.DiscountType,
		DiscountAmount: req.DiscountAmount,
		MinimumAmount:  req.MinimumAmount,
		ExpiryDate:     expiryDate,
		UsageLimit:     req.UsageLimit,
		IsActive:       true,
	}

	return repository.CreateCoupon(coupon)
}

func GetAllCoupons(page, limit int) (*models.CouponListResponse, error) {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	coupons, totalCount, err := repository.GetAllCoupons(page, limit)
	if err != nil {
		return nil, err
	}

	var response []models.CouponResponse

	for _, coupon := range coupons {

		response = append(response, models.CouponResponse{
			ID:             coupon.ID,
			Code:           coupon.Code,
			DiscountType:   coupon.DiscountType,
			DiscountAmount: coupon.DiscountAmount,
			MinimumAmount:  coupon.MinimumAmount,
			ExpiryDate:     coupon.ExpiryDate.Format("2006-01-02"),
			UsageLimit:     coupon.UsageLimit,
			UsedCount:      coupon.UsedCount,
			IsActive:       coupon.IsActive,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &models.CouponListResponse{
		Coupons:    response,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func GetAvailableCoupons(userID uint, page, limit int) (*models.UserCouponListResponse, error) {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	subtotal, err := repository.GetCartSubtotal(userID)
	if err != nil {
		return nil, err
	}

	coupons, err := repository.GetActiveCoupons()
	if err != nil {
		return nil, err
	}

	var resp []models.UserCouponResponse

	for _, coupon := range coupons {

		if time.Now().After(coupon.ExpiryDate) {
			continue
		}

		if coupon.UsedCount >= coupon.UsageLimit {
			continue
		}

		item := models.UserCouponResponse{
			Code:           coupon.Code,
			DiscountType:   coupon.DiscountType,
			DiscountAmount: coupon.DiscountAmount,
			MinimumAmount:  coupon.MinimumAmount,
		}

		if subtotal >= coupon.MinimumAmount {
			item.Eligible = true
		} else {
			item.Message = "minimum purchase amount not met"
		}

		resp = append(resp, item)
	}

	totalCount := int64(len(resp))
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	start := (page - 1) * limit

	if start >= len(resp) {
		resp = []models.UserCouponResponse{}
	} else {

		end := start + limit

		if end > len(resp) {
			end = len(resp)
		}

		resp = resp[start:end]
	}

	return &models.UserCouponListResponse{
		Coupons:    resp,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func UpdateCoupon(id uint, req models.UpdateCouponRequest) error {

	coupon, err := repository.GetCouponByID(id)
	if err != nil {
		return errors.New("coupon not found")
	}

	if req.DiscountType != "fixed" && req.DiscountType != "percentage" {
		return errors.New("invalid discount type")
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return errors.New("invalid expiry date")
	}

	coupon.Code = req.Code
	coupon.DiscountType = req.DiscountType
	coupon.DiscountAmount = req.DiscountAmount
	coupon.MinimumAmount = req.MinimumAmount
	coupon.ExpiryDate = expiryDate
	coupon.UsageLimit = req.UsageLimit
	coupon.IsActive = req.IsActive

	return repository.UpdateCoupon(coupon)
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

	var discount float64

	if coupon.DiscountType == "percentage" {
		discount = subtotal * coupon.DiscountAmount / 100
	} else {
		discount = coupon.DiscountAmount
	}

	if discount > subtotal {
		discount = subtotal
	}

	final := subtotal - discount
	err = repository.UpdateCartCoupon(userID,coupon.Code,discount,)
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
