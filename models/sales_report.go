package models

import "time"

type SalesReportOrderResponse struct {
	OrderID        string    `json:"order_id"`
	OrderDate      time.Time `json:"order_date"`
	CustomerName   string    `json:"customer_name"`
	PaymentMethod  string    `json:"payment_method"`
	OrderStatus    string    `json:"order_status"`
	TotalProducts  int       `json:"total_products"`
	OrderAmount    float64   `json:"order_amount"`
	OfferDiscount  float64   `json:"offer_discount"`
	CouponDiscount float64   `json:"coupon_discount"`
	TotalDiscount  float64   `json:"total_discount"`
	NetAmount      float64   `json:"net_amount"`
}

type SalesReportResponse struct {
	FromDate       string                     `json:"from_date"`
	ToDate         string                     `json:"to_date"`
	TotalOrders    int64                       `json:"total_orders"`
	TotalProducts  int                        `json:"total_products"`
	GrossSales     float64                    `json:"gross_sales"`
	OfferDiscount  float64                    `json:"offer_discount"`
	CouponDiscount float64                    `json:"coupon_discount"`
	TotalDiscount  float64                    `json:"total_discount"`
	NetSales       float64                    `json:"net_sales"`
	Orders         []SalesReportOrderResponse `json:"orders"`
}
