package models

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type UpdateCartQuantityRequest struct {
	Action    string `json:"action" binding:"required"` // increment or decrement
}

type CartItemResponse struct {
	ProductID    uint     `json:"product_id"`
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	Quantity     int      `json:"quantity"`
	Stock        int      `json:"stock"`
	CategoryName string   `json:"category_name"`
	Images       []string `json:"images"`
	Subtotal     float64  `json:"subtotal"`
	Status       string   `json:"status"`
	IsAvailable  bool     `json:"is_available"`
}

type CartResponse struct {
	Items           []CartItemResponse `json:"items"`
	TotalAmount     float64            `json:"total_amount"`
	CheckoutAllowed bool               `json:"checkout_allowed"`
}
