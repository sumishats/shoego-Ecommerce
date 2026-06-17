package models


type CreateRazorpayOrderRequest struct {
	AddressID uint `json:"address_id" binding:"required"`
}

type CreateRazorpayOrderResponse struct {
	OrderID         string  `json:"order_id"`
	RazorpayOrderID string  `json:"razorpay_order_id"`
	Currency        string  `json:"currency"`
	Key             string  `json:"key"`
	FinalAmount     float64 `json:"final_amount"`
}

type VerifyPaymentRequest struct {
	// OrderID            string `json:"order_id" binding:"required"`
	RazorpayOrderID    string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID  string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature  string `json:"razorpay_signature" binding:"required"`
}