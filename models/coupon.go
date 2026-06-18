package models

type CreateCouponRequest struct {
	Code           string  `json:"code" binding:"required"`
	DiscountAmount float64 `json:"discount_amount" binding:"required"`
	MinimumAmount  float64 `json:"minimum_amount" binding:"required"`
	ExpiryDate     string  `json:"expiry_date" binding:"required"`
	UsageLimit     int     `json:"usage_limit" binding:"required"`
}

type CouponResponse struct {
	ID             uint    `json:"id"`
	Code           string  `json:"code"`
	DiscountAmount float64 `json:"discount_amount"`
	MinimumAmount  float64 `json:"minimum_amount"`
	ExpiryDate     string  `json:"expiry_date"`
	UsageLimit     int     `json:"usage_limit"`
	UsedCount      int     `json:"used_count"`
	IsActive       bool    `json:"is_active"`
}
type CouponApplyResponse struct {
	Code     string  `json:"code"`
	Subtotal float64 `json:"subtotal"`
	Discount float64 `json:"discount"`
	Final    float64 `json:"final"`
}