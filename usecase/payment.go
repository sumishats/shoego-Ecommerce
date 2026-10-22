package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"shoego/config"
	"shoego/database"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	
	orderID := "ORD-" + uuid.NewString()

	order := &domain.Order{
		OrderID:        orderID,
		UserID:         userID,
		AddressID:      req.AddressID,
		OrderStatus:    "pending",
		PaymentMethod:  "razorpay",
		PaymentStatus:  "pending",
		Subtotal:       checkout.Subtotal,
		TaxAmount:      checkout.TaxAmount,
		CouponDiscount: checkout.CouponDiscount,
		ShippingCharge: checkout.ShippingCharge,
		FinalAmount:    checkout.FinalAmount,
	}

	var orderItems []domain.OrderItem

	for _, item := range cartItems {

		price := item.Product.Price

		bestOffer := 0.0

		productOffer, err := repository.GetActiveProductOffer(item.ProductID)
		if err == nil {
			bestOffer = productOffer.DiscountPercentage
		}

		categoryOffer, err := repository.GetActiveCategoryOffer(item.Product.CategoryID)
		if err == nil && categoryOffer.DiscountPercentage > bestOffer {
			bestOffer = categoryOffer.DiscountPercentage
		}

		if bestOffer > 0 {
			price = price - (price*bestOffer)/100
		}

		price = math.Round(price*100) / 100

		var variantID *uint

		if item.VariantID != nil {
			variantID = item.VariantID
		}

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:  item.Product.ID,
			VariantID:  variantID,
			Quantity:   item.Quantity,
			Price:      price,
			TotalPrice: math.Round(price*float64(item.Quantity)*100) / 100,
		})
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

	amount := int64(math.Round(checkout.FinalAmount * 100))

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
		Currency:        "INR",
		Key:             config.AppConfig.RAZORPAY_KEY_ID,
		FinalAmount:     checkout.FinalAmount,
	}, nil
}

func VerifyPayment(userID uint, req models.VerifyPaymentRequest) error {
	
	payment, err := repository.GetPaymentByRazorpayOrderID(
		req.RazorpayOrderID,
	)
	if err != nil {
		return err
	}

	
	if payment.UserID != userID {
		return errors.New("payment does not belong to this user")
	}


	data := payment.RazorpayOrderID + "|" + req.RazorpayPaymentID

	h := hmac.New(
		sha256.New,
		[]byte(config.AppConfig.RAZORPAY_KEY_SECRET),
	)
	h.Write([]byte(data))

	expectedSignature := hex.EncodeToString(h.Sum(nil))
	receivedSignature := strings.TrimSpace(req.RazorpaySignature)

	if !hmac.Equal(
		[]byte(expectedSignature),
		[]byte(receivedSignature),
	) {
		return errors.New("payment verification failed")
	}

	
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var currentPayment domain.Payment

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", payment.ID).
			First(&currentPayment).Error
		if err != nil {
			return err
		}

		if currentPayment.UserID != userID {
			return errors.New("payment does not belong to this user")
		}

		
		if currentPayment.Status == "paid" {
			return nil
		}

		
		err = repository.UpdateRazorpayDetailsTx(
			tx,
			currentPayment.OrderID,
			req.RazorpayPaymentID,
			receivedSignature,
		)
		if err != nil {
			return err
		}

		
		err = repository.UpdatePaymentStatusTx(
			tx,
			currentPayment.OrderID,
			"paid",
		)
		if err != nil {
			return err
		}

		
		err = tx.Model(&domain.Order{}).
			Where("id = ?", currentPayment.OrderID).
			Update("order_status", "placed").Error
		if err != nil {
			return err
		}

		
		err = repository.ReduceStockAfterPaymentTx(
			tx,
			currentPayment.OrderID,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}


	if err := repository.ClearCartByUserID(userID); err != nil {
		fmt.Println("Cart clearing failed:", err)
	}

	return nil
}

func HandlePaymentFailure(userID uint, req models.PaymentFailureRequest) error {

	payment, err := repository.GetPaymentByRazorpayOrderID(req.RazorpayOrderID)
	if err != nil {
		return err
	}

	
	if payment.UserID != userID {
		return errors.New("payment does not belong to this user")
	}

	
	result, err := config.RazorpayClient.Payment.Fetch(
		req.RazorpayPaymentID,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	razorpayOrderID, _ := result["order_id"].(string)
	status, _ := result["status"].(string)

	if razorpayOrderID != req.RazorpayOrderID {
		return errors.New("payment does not belong to this razorpay order")
	}

	if status != "failed" {
		return errors.New("razorpay has not confirmed payment failure")
	}

	
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var currentPayment domain.Payment

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", payment.ID).
			First(&currentPayment).Error
		if err != nil {
			return err
		}

		
		if currentPayment.UserID != userID {
			return errors.New("payment does not belong to this user")
		}

		if currentPayment.Status == "paid" {
			return errors.New("payment is already completed")
		}

		
		err = tx.Model(&domain.Payment{}).
			Where("id = ? AND status <> ?", currentPayment.ID, "paid").
			Updates(map[string]interface{}{
				"status":              "failed",
				"razorpay_payment_id": req.RazorpayPaymentID,
			}).Error
		if err != nil {
			return err
		}

		
		err = tx.Model(&domain.Order{}).
			Where("id = ? AND payment_status <> ?", currentPayment.OrderID, "paid").
			Update("payment_status", "failed").Error
		if err != nil {
			return err
		}

		return nil
	})
}

type RazorpayWebhookPayload struct {
	Event string `json:"event"`

	Payload struct {
		Payment struct {
			Entity struct {
				ID       string `json:"id"`
				OrderID  string `json:"order_id"`
				Status   string `json:"status"`
				Amount   int64  `json:"amount"`
				Currency string `json:"currency"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

func HandleRazorpayWebhook(rawBody []byte) error {
	var event RazorpayWebhookPayload

	if err := json.Unmarshal(rawBody, &event); err != nil {
		return errors.New("invalid webhook payload")
	}

	payment := event.Payload.Payment.Entity

	switch event.Event {
	case "payment.captured":
		if payment.ID == "" || payment.OrderID == "" {
			return errors.New("missing Razorpay payment details")
		}

		if payment.Status != "captured" {
			return errors.New("payment is not captured")
		}

		if payment.Currency != "INR" {
			return errors.New("unexpected payment currency")
		}

		return handleWebhookPaymentSuccess(
			payment.OrderID,
			payment.ID,
			payment.Amount,
		)

	case "payment.failed":
		if payment.ID == "" || payment.OrderID == "" {
			return errors.New("missing Razorpay payment details")
		}

		if payment.Status != "failed" {
			return errors.New("payment is not failed")
		}

		return handleWebhookPaymentFailure(
			payment.OrderID,
			payment.ID,
		)

	case "order.paid":

		return nil

	default:
		
		return nil
	}
}

func handleWebhookPaymentSuccess(razorpayOrderID string, razorpayPaymentID string, amountPaise int64) error {
	var userID uint

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		id, err := repository.ProcessWebhookPaymentSuccessTx(
			tx,
			razorpayOrderID,
			razorpayPaymentID,
			amountPaise,
		)
		if err != nil {
			return err
		}

		userID = id
		return nil
	})
	if err != nil {
		return err
	}

	// Cart clearing is outside the payment transaction.
	if userID != 0 {
		if err := repository.ClearCartByUserID(userID); err != nil {
			fmt.Println("Cart clearing failed:", err)
		}
	}

	return nil
}

func handleWebhookPaymentFailure(razorpayOrderID string, razorpayPaymentID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		return repository.ProcessWebhookPaymentFailureTx(
			tx,
			razorpayOrderID,
			razorpayPaymentID,
		)
	})
}
