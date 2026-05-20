package models

import "time"

type AdminOrderListResponse struct {
	ID            uint      `json:"id"`
	OrderID       string    `json:"order_id"`
	Date          time.Time `json:"date"`
	UserName      string    `json:"user_name"`
	UserEmail     string    `json:"user_email"`
	UserPhone     string    `json:"user_phone"`
	ItemCount     int       `json:"item_count"`
	FinalAmount   float64   `json:"final_amount"`
	OrderStatus   string    `json:"order_status"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
}

type AdminOrderDetailResponse struct {
	ID            uint                      `json:"id"`
	OrderID       string                    `json:"order_id"`
	Date          time.Time                 `json:"date"`
	OrderStatus   string                    `json:"order_status"`
	PaymentMethod string                    `json:"payment_method"`
	PaymentStatus string                    `json:"payment_status"`
	ItemCount     int                       `json:"item_count"`
	Subtotal      float64                   `json:"subtotal"`
	Discount      float64                   `json:"discount"`
	ShippingFee   float64                   `json:"shipping_fee"`
	FinalAmount   float64                   `json:"final_amount"`
	User          AdminOrderUserResponse    `json:"user"`
	Address       AdminOrderAddressResponse `json:"address"`
	Items         []AdminOrderItemResponse  `json:"items"`
}

type AdminOrderUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AdminOrderAddressResponse struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
}

type AdminOrderItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	Image       string  `json:"image"`
}

type UpdateOrderStatusRequest struct {
	OrderStatus string `json:"order_status" binding:"required"`
}

type AdminOrderPaginationResponse struct {
	Orders     []AdminOrderListResponse `json:"orders"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalCount int64                    `json:"total_count"`
	TotalPages int                      `json:"total_pages"`
	Image      string                   `json:"image"`
}

//inventory

type AdminInventoryResponse struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SKU          string  `json:"sku"`
	CategoryName string  `json:"category_name"`
	Stock        int     `json:"stock"`
	Price        float64 `json:"price"`
	IsListed     bool    `json:"is_listed"`
	StockStatus  string  `json:"stock_status"`
}

type AdminInventoryPaginationResponse struct {
	Products   []AdminInventoryResponse `json:"products"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalCount int64                    `json:"total_count"`
	TotalPages int                      `json:"total_pages"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" binding:"required"`
}

//user order details

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type CancelOrderItemRequest struct {
	Reason string `json:"reason"`
}

type ReturnOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type ReturnOrderItemRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type OrderListItemResponse struct {
	OrderID       string    `json:"order_id"`
	OrderStatus   string    `json:"order_status"`
	OrderDate     time.Time `json:"order_date"`
	FinalAmount   float64   `json:"final_amount"`
	PaymentMethod string    `json:"payment_method"`
}

type OrderItemResponse struct {
	ItemID      uint    `json:"item_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Image       string  `json:"image"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
	ItemStatus  string  `json:"item_status"`
}

type OrderDetailResponse struct {
	OrderID        string              `json:"order_id"`
	OrderStatus    string              `json:"order_status"`
	OrderDate      time.Time           `json:"order_date"`
	PaymentMethod  string              `json:"payment_method"`
	PaymentStatus  string              `json:"payment_status"`
	Subtotal       float64             `json:"subtotal"`
	TaxAmount      float64             `json:"tax_amount"`
	DiscountAmount float64             `json:"discount_amount"`
	ShippingCharge float64             `json:"shipping_charge"`
	FinalAmount    float64             `json:"final_amount"`
	Items          []OrderItemResponse `json:"items"`
}
