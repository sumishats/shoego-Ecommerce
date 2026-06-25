package models

type CreateProductOfferRequest struct {
	ProductID          uint    `json:"product_id"`
	OfferName          string  `json:"offer_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
}
type ProductOfferResponse struct {
	ID                 uint    `json:"id"`
	ProductID          uint    `json:"product_id"`
	ProductName        string  `json:"product_name"`
	OfferName          string  `json:"offer_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	IsActive           bool    `json:"is_active"`
}

type CreateCategoryOfferRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	OfferName  string `json:"offer_name" binding:"required"`

	DiscountPercentage float64 `json:"discount_percentage"`

	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
type CategoryOfferResponse struct {
	ID                 uint    `json:"id"`
	CategoryID         uint    `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	OfferName          string  `json:"offer_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
	IsActive           bool    `json:"is_active"`
}

type ReferralResponse struct {
	ReferralCode string `json:"referral_code"`
	ReferralLink string `json:"referral_link"`
}
type UpdateProductOfferRequest struct {
	ProductID          uint    `json:"product_id"`
	OfferName          string  `json:"offer_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
}

type UpdateCategoryOfferRequest struct {
	CategoryID         uint    `json:"category_id"`
	OfferName          string  `json:"offer_name"`
	DiscountPercentage float64 `json:"discount_percentage"`
	StartDate          string  `json:"start_date"`
	EndDate            string  `json:"end_date"`
}
type PaginatedResponse struct {
	Offers     []ProductOfferResponse `json:"offers"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalCount int64                  `json:"total_count"`
	TotalPages int                    `json:"total_pages"`
}
type PaginatedCategoryOfferResponse struct {
	Offers     []CategoryOfferResponse `json:"offers"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalCount int64                   `json:"total_count"`
	TotalPages int                     `json:"total_pages"`
}