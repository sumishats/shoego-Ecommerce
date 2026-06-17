package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"shoego/config"
	"shoego/database"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"strings"
	"time"
)

func CreateRazorpayOrder(userID uint, req models.CreateRazorpayOrderRequest) (*models.CreateRazorpayOrderResponse, error) {

	_, err := repository.GetAddressByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, err
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
		OrderStatus:    "pending",
		PaymentMethod:  "razorpay",
		PaymentStatus:  "pending",
		Subtotal:       checkout.Subtotal,
		TaxAmount:      checkout.TaxAmount,
		DiscountAmount: checkout.DiscountAmount,
		ShippingCharge: checkout.ShippingCharge,
		FinalAmount:    checkout.FinalAmount,
	}

	var orderItems []domain.OrderItem

	for _, item := range cartItems {

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:  item.Product.ID,
			Quantity:   item.Quantity,
			Price:      item.Product.Price,
			TotalPrice: item.Product.Price * float64(item.Quantity),
		},
		)
	}

	payment := &domain.Payment{
		UserID:        userID,
		Amount:        checkout.FinalAmount,
		PaymentMethod: "razorpay",
		Status:        "pending",
	}

	err = repository.CreatePendingOrderTransaction(order, orderItems, payment)
	if err != nil {
		return nil, err
	}

	amount := int64(checkout.FinalAmount * 100)

	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  order.OrderID,
	}

	body, err := config.RazorpayClient.Order.Create(data, nil)
	if err != nil {
		return nil, err
	}

	razorpayOrderID := body["id"].(string)

	err = repository.UpdateRazorpayOrderID(order.ID, razorpayOrderID)
	if err != nil {
		return nil, err
	}

	return &models.CreateRazorpayOrderResponse{
		OrderID:         order.OrderID,
		RazorpayOrderID: razorpayOrderID,
		// Amount:          amount,
		Currency:    "INR",
		Key:         config.AppConfig.RAZORPAY_KEY_ID,
		FinalAmount: checkout.FinalAmount,
	}, nil
}

func VerifyPayment(req models.VerifyPaymentRequest) error {

	payment, err := repository.GetPaymentByRazorpayOrderID(req.RazorpayOrderID)
	if err != nil {
		return err
	}

	//razorpay verification
	data := payment.RazorpayOrderID + "|" + req.RazorpayPaymentID

	h := hmac.New(sha256.New, []byte(config.AppConfig.RAZORPAY_KEY_SECRET))

	h.Write([]byte(data))

	expectedSignature := hex.EncodeToString(h.Sum(nil))

	req.RazorpaySignature = strings.TrimSpace(req.RazorpaySignature)

	if expectedSignature != req.RazorpaySignature {
		return errors.New("payment verification failed")
	}

	err = repository.UpdateRazorpayDetailsTx(database.DB, payment.OrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if err != nil {
		return err
	}

	err = repository.UpdatePaymentStatusTx(database.DB, payment.OrderID, "paid")
	if err != nil {
		return err
	}

	err = repository.UpdateOrderAfterPayment(payment.OrderID)
	if err != nil {
		return err
	}
	// ADDED
	err = repository.ReduceStockAfterPayment(payment.OrderID)
	if err != nil {
		return err
	}

	// Debug logs
	fmt.Println("Payment ID:", payment.ID)
	fmt.Println("Order ID:", payment.OrderID)
	fmt.Println("User ID:", payment.UserID)

	// Clear cart
	err = repository.ClearCartByUserID(payment.UserID)
	if err != nil {
		return err
	}

	fmt.Println("Cart cleared successfully for user:", payment.UserID)

	return nil

}
