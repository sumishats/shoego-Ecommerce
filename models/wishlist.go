package models

type AddToWishlistRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type WishlistProductResponse struct {
	ProductID     uint     `json:"product_id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Stock         int      `json:"stock"`
	CategoryID    uint     `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	IsListed      bool     `json:"is_listed"`
	Status        string   `json:"status"`
	Images        []string `json:"images"`
	IsWishlisted  bool     `json:"is_wishlisted"`
}

type WishlistResponse struct {
	Items []WishlistProductResponse `json:"items"`
}