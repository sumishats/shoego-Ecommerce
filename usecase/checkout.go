package usecase

import (
	"errors"
	"fmt"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"time"
)

func GetCheckoutPage(userID uint) (*models.CheckoutPageResponse, error) {
	addresses, err := repository.GetCheckoutAddresses(userID)
	if err != nil {
		return nil, err
	}

	cartItems, err := repository.GetCartItemsForCheckout(userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
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

	var itemRes []models.CheckoutItemResponse
	var subtotal float64

	for _, item := range cartItems {
		product := item.Product

		if !product.IsListed {
			return nil, errors.New(product.Name + " is unavailable")
		}

		if !product.Category.IsListed {
			return nil, errors.New(product.Name + " category is unavailable")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New(product.Name + " does not have enough stock")
		}

		image := ""
		if len(product.Images) > 0 {
			image = product.Images[0].ImageURL
		}

		itemTotal := product.Price * float64(item.Quantity)
		subtotal += itemTotal

		itemRes = append(itemRes, models.CheckoutItemResponse{
			ProductID: product.ID,
			Name:      product.Name,
			Image:     image,
			Quantity:  item.Quantity,
			Price:     product.Price,
			ItemTotal: itemTotal,
		})
	}

	taxAmount := subtotal * 0.05
	discountAmount := 0.0
	shippingCharge := 50.0

	if subtotal >= 1000 {
		shippingCharge = 0
	}

	finalAmount := subtotal + taxAmount + shippingCharge - discountAmount

	return &models.CheckoutPageResponse{
		Addresses:      addressRes,
		Items:          itemRes,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		DiscountAmount: discountAmount,
		ShippingCharge: shippingCharge,
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
		DiscountAmount: checkout.DiscountAmount,
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

	if err := repository.CreateOrderTransaction(order,orderItems,userID,);err != nil {
		return nil, err
	}

	return &models.PlaceOrderResponse{
		OrderID:     orderID,
		Message:     "order placed successfully",
		FinalAmount: checkout.FinalAmount,
	}, nil
}

func GetOrderSuccess(userID uint, orderID string) (*models.OrderSuccessResponse, error) {
	_, err := repository.GetOrderByOrderID(userID, orderID)
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
