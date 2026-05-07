package models


type CheckoutItemResponse struct{
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ItemTotal float64 `json:"item_total"`
}


type CheckoutAddressResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
	IsDefault bool   `json:"is_default"`
}

type CheckoutPageResponse struct {
	Addresses      []CheckoutAddressResponse `json:"addresses"`
	Items          []CheckoutItemResponse    `json:"items"`
	Subtotal       float64                   `json:"subtotal"`
	TaxAmount      float64                   `json:"tax_amount"`
	DiscountAmount float64                   `json:"discount_amount"`
	ShippingCharge float64                   `json:"shipping_charge"`
	FinalAmount    float64                   `json:"final_amount"`
}

type PlaceOrderRequest struct {
	AddressID     uint   `json:"address_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type PlaceOrderResponse struct {
	OrderID     string  `json:"order_id"`
	Message     string  `json:"message"`
	FinalAmount float64 `json:"final_amount"`
}

type OrderSuccessResponse struct {
	OrderID  string `json:"order_id"`
	Message  string `json:"message"`
	Image    string `json:"image"`
	Buttons  []ButtonResponse `json:"buttons"`
}

type ButtonResponse struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}