package usecase

import (
	"errors"
	"fmt"
	"math"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"time"
)

func GetCheckoutPage(userID uint) (*models.CheckoutPageResponse, error) {

	// 1. Addresses
	addresses, err := repository.GetCheckoutAddresses(userID)
	if err != nil {
		return nil, err
	}

	var addressRes []models.CheckoutAddressResponse
	for _, addr := range addresses {
		addressRes = append(addressRes, models.CheckoutAddressResponse{
			ID:        addr.ID,
			Name:      addr.Name,
			Phone:     addr.Phone,
			HouseName: addr.HouseName,
			Street:    addr.Street,
			City:      addr.City,
			State:     addr.State,
			Pincode:   addr.Pincode,
			IsDefault: addr.IsDefault,
		})
	}

	cartItems, err := repository.GetCartItemsForCheckout(userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var items []models.CheckoutItemResponse

	for _, item := range cartItems {

		product := item.Product

		if product.ID == 0 {
			continue
		}

		if !product.IsListed {
			return nil, errors.New(product.Name + " unavailable")
		}

		if !product.Category.IsListed {
			return nil, errors.New(product.Name + " category unavailable")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New(product.Name + " out of stock")
		}

		image := ""
		if len(product.Images) > 0 {
			image = product.Images[0].ImageURL
		}

		discountedPrice := product.Price

		offerPercentage := 0.0


		productOffer, err := repository.GetActiveProductOffer(product.ID)

		if err == nil {

			offerPercentage = productOffer.DiscountPercentage
		}

		categoryOffer, err := repository.GetActiveCategoryOffer(product.CategoryID)

		if err == nil {

			if categoryOffer.DiscountPercentage > offerPercentage {

				offerPercentage = categoryOffer.DiscountPercentage
			}
		}


		if offerPercentage > 0 {

			discountedPrice =product.Price -(product.Price*offerPercentage)/100
		}

		items = append(items, models.CheckoutItemResponse{
			ProductID: product.ID,
			Name:      product.Name,
			Image:     image,
			Quantity:  item.Quantity,
			Price:     discountedPrice,
			ItemTotal: discountedPrice * float64(item.Quantity),
		})

	}

	
	cart, err := repository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	subtotal, err := repository.GetCartSubtotal(userID)
	if err != nil {
		return nil, err
	}

	tax := math.Round((subtotal*0.05)*100) / 100

	
	shipping := 50.0
	if subtotal >= 1000 {
		shipping = 0
	}

	
	couponCode := ""
	couponDiscount := 0.0

	if cart != nil {
		couponCode = cart.CouponCode
		couponDiscount = cart.DiscountAmount
	}

	finalAmount := math.Round((subtotal+tax+shipping-couponDiscount)*100) / 100

	return &models.CheckoutPageResponse{
		Addresses:      addressRes,
		Items:          items,
		Subtotal:       subtotal,
		TaxAmount:      tax,
		CouponCode:     couponCode,
		CouponDiscount: couponDiscount,
		ShippingCharge: shipping,
		FinalAmount:    finalAmount,
	}, nil
}


func PlaceCODOrder(userID uint, req models.PlaceOrderRequest) (*models.PlaceOrderResponse, error) {
	if req.PaymentMethod != "cod" {
		return nil, errors.New("only Cash on Delivery is available now")
	}

	_, err := repository.GetAddressByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	checkout, err := GetCheckoutPage(userID)
	if err != nil {
		return nil, err
	}

	cartItems, err := repository.GetCartItemsForCheckout(userID)
	if err != nil {
		return nil, err
	}

	orderID := fmt.Sprintf("ORD%d", time.Now().UnixNano())

	order := &domain.Order{
		OrderID:        orderID,
		UserID:         userID,
		AddressID:      req.AddressID,
		OrderStatus:    "placed",
		PaymentMethod:  "cod",
		Subtotal:       checkout.Subtotal,
		TaxAmount:      checkout.TaxAmount,
		CouponDiscount: checkout.CouponDiscount,
		ShippingCharge: checkout.ShippingCharge,
		FinalAmount:    checkout.FinalAmount,
	}

	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		product := item.Product

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:  product.ID,
			Quantity:   item.Quantity,
			Price:      product.Price,
			TotalPrice: product.Price * float64(item.Quantity),
		})
	}

	if err := repository.CreateOrderTransaction(order, orderItems, userID); err != nil {
		return nil, err
	}

	return &models.PlaceOrderResponse{
		OrderID:     orderID,
		Message:     "order placed successfully",
		FinalAmount: checkout.FinalAmount,
	}, nil
}
func PlaceWalletOrder(userID uint, req models.PlaceOrderRequest) (*models.PlaceOrderResponse, error) {
	if req.PaymentMethod != "wallet" {
		return nil, errors.New("invalid payment method")
	}
	_, err := repository.GetAddressByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, errors.New("address not found")
	}
	checkout, err := GetCheckoutPage(userID)
	if err != nil {
		return nil, err
	}
	wallet, err := repository.GetWalletByUserID(userID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	if wallet.Balance < checkout.FinalAmount {
		return nil, errors.New("insufficient wallet balance")
	}
	cartItems, err := repository.GetCartItemsForCheckout(userID)
	if err != nil {
		return nil, err
	}
	orderID := fmt.Sprintf("ORD%d", time.Now().UnixNano())
	order := &domain.Order{
		OrderID:       orderID,
		UserID:        userID,
		AddressID:     req.AddressID,
		OrderStatus:   "placed",
		PaymentMethod: "wallet",

		PaymentStatus:  "paid",
		Subtotal:       checkout.Subtotal,
		TaxAmount:      checkout.TaxAmount,
		CouponDiscount: checkout.CouponDiscount,
		ShippingCharge: checkout.ShippingCharge,
		FinalAmount:    checkout.FinalAmount,
	}
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		product := item.Product
		orderItems = append(orderItems, domain.OrderItem{
			ProductID:  product.ID,
			Quantity:   item.Quantity,
			Price:      product.Price,
			TotalPrice: product.Price * float64(item.Quantity),
		})
	}
	err = repository.CreateOrderTransaction(order, orderItems, userID)
	if err != nil {
		return nil, err
	}
	newBalance := wallet.Balance - checkout.FinalAmount

	err = repository.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return nil, err
	}
	transaction := &domain.WalletTransaction{
		WalletID:    wallet.ID,
		Amount:      checkout.FinalAmount,
		Type:        "debit",
		Description: "order payment",
	}
	err = repository.CreateWalletTransaction(transaction)
	if err != nil {
		return nil, err
	}
	return &models.PlaceOrderResponse{
		OrderID:     orderID,
		Message:     "wallet payment successful",
		FinalAmount: checkout.FinalAmount,
	}, nil
}

func GetOrderSuccess(orderID string) (*models.OrderSuccessResponse, error) {
	_, err := repository.GetOrderByOrderIDPayment(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	return &models.OrderSuccessResponse{
		OrderID: orderID,
		Message: "Thank you! Your order has been placed successfully.",
		Image:   "/static/images/order-success.png",
		Buttons: []models.ButtonResponse{
			{
				Label: "View Order Details",
				URL:   "/orders/" + orderID,
			},
			{
				Label: "Continue Shopping",
				URL:   "/products",
			},
		},
	}, nil
}
