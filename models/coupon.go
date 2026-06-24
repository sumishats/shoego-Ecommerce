package models

type CreateCouponRequest struct {
	Code           string  `json:"code" binding:"required"`
	DiscountType   string  `json:"discount_type" binding:"required"`
	DiscountAmount float64 `json:"discount_amount" binding:"required"`
	MinimumAmount  float64 `json:"minimum_amount" binding:"required"`
	ExpiryDate     string  `json:"expiry_date" binding:"required"`
	UsageLimit     int     `json:"usage_limit" binding:"required"`
}

type CouponResponse struct {
	ID             uint    `json:"id"`
	Code           string  `json:"code"`
	DiscountType   string  `json:"discount_type"`
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

type UpdateCouponRequest struct {
	Code           string  `json:"code" binding:"required"`
	DiscountType   string  `json:"discount_type" binding:"required"`
	DiscountAmount float64 `json:"discount_amount" binding:"required"`
	MinimumAmount  float64 `json:"minimum_amount" binding:"required"`
	ExpiryDate     string  `json:"expiry_date" binding:"required"`
	UsageLimit     int     `json:"usage_limit" binding:"required"`
	IsActive       bool    `json:"is_active"`
}

type UserCouponResponse struct {
	Code           string  `json:"code"`
	DiscountType   string  `json:"discount_type"`
	DiscountAmount float64 `json:"discount_amount"`
	MinimumAmount  float64 `json:"minimum_amount"`
	Eligible       bool    `json:"eligible"`
	Message        string  `json:"message,omitempty"`
}
type CouponListResponse struct {
	Coupons    []CouponResponse `json:"coupons"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalCount int64            `json:"total_count"`
	TotalPages int              `json:"total_pages"`
}
type UserCouponListResponse struct {
	Coupons    []UserCouponResponse `json:"coupons"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalCount int64                `json:"total_count"`
	TotalPages int                  `json:"total_pages"`
}