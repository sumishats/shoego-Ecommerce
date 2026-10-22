package repository

import (
	"errors"
	"math"
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePaymentTx(tx *gorm.DB, payment *domain.Payment) error {
	return tx.Create(payment).Error
}
func GetPaymentByOrderID(orderID uint) (*domain.Payment, error) {
	var payment domain.Payment

	err := database.DB.Where("order_id = ?", orderID).First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}
func UpdatePaymentStatusTx(tx *gorm.DB, orderID uint, status string) error {

	err := tx.Model(&domain.Payment{}).
		Where("order_id = ?", orderID).
		Update("status", status).Error
	if err != nil {
		return err
	}

	err = tx.Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateRazorpayDetailsTx(tx *gorm.DB, orderID uint, razorpayPaymentID string, razorpaySignature string) error {
	return tx.Model(&domain.Payment{}).Where("order_id = ?", orderID).Updates(map[string]interface{}{"razorpay_payment_id": razorpayPaymentID, "razorpay_signature": razorpaySignature}).Error
}

func CreatePendingOrderTransaction(order *domain.Order, orderItems []domain.OrderItem, payment *domain.Payment) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		payment.OrderID = order.ID

		if err := CreatePaymentTx(tx, payment); err != nil {
			return err
		}

		return nil
	})
}

func UpdateRazorpayOrderID(orderID uint, razorpayOrderID string) error {
	return database.DB.Model(&domain.Payment{}).Where("order_id = ?", orderID).Update("razorpay_order_id", razorpayOrderID).Error
}

func GetPaymentByRazorpayOrderID(razorpayOrderID string) (*domain.Payment, error) {

	var payment domain.Payment

	err := database.DB.Where("razorpay_order_id = ?", razorpayOrderID).First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func UpdateOrderAfterPayment(orderID uint) error {

	return database.DB.Model(&domain.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"payment_status": "paid",
			"order_status":   "placed",
		}).
		Error
}

// func ReduceStockAfterPayment(orderID uint) error {

// 	return database.DB.Transaction(func(tx *gorm.DB) error {

// 		var orderItems []domain.OrderItem

// 		err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error
// 		if err != nil {
// 			return err
// 		}

// 		for _, item := range orderItems {

// 			if item.VariantID != nil {

// 				result := tx.Model(&domain.ProductVariant{}).
// 					Where("id = ? AND stock >= ?", *item.VariantID, item.Quantity).
// 					Update("stock", gorm.Expr("stock - ?", item.Quantity))

// 				if result.Error != nil {
// 					return result.Error
// 				}

// 				if result.RowsAffected == 0 {
// 					return errors.New("variant stock unavailable")
// 				}

// 			} else {
// 				result := tx.Model(&domain.Product{}).
// 					Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
// 					Update("stock", gorm.Expr("stock - ?", item.Quantity))

// 				if result.Error != nil {
// 					return result.Error
// 				}

// 				if result.RowsAffected == 0 {
// 					return errors.New("product stock unavailable")
// 				}
// 			}
// 		}

// 		return nil
// 	})
// }

func ReduceStockAfterPaymentTx(tx *gorm.DB, orderID uint) error {
	var orderItems []domain.OrderItem

	err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return err
	}

	if len(orderItems) == 0 {
		return errors.New("no order items found")
	}

	for _, item := range orderItems {
		if item.VariantID != nil {
			result := tx.Model(&domain.ProductVariant{}).
				Where("id = ? AND stock >= ?", *item.VariantID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return errors.New("variant stock unavailable")
			}
		} else {
			result := tx.Model(&domain.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return errors.New("product stock unavailable")
			}
		}
	}

	return nil
}

func UpdateOrderPaymentStatus(orderID uint, status string) error {
	return database.DB.Model(&domain.Order{}).Where("id = ?", orderID).Update("payment_status", status).Error
}

// Process successful payment and stock deduction atomically.
func ProcessWebhookPaymentSuccessTx(
	tx *gorm.DB,
	razorpayOrderID string,
	razorpayPaymentID string,
	amountPaise int64,
) (uint, error) {
	var payment domain.Payment

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&payment).Error
	if err != nil {
		return 0, err
	}

	if payment.Status == "paid" {
		return payment.UserID, nil
	}

	expectedPaise := int64(math.Round(payment.Amount * 100))
	if amountPaise != expectedPaise {
		return 0, errors.New("payment amount mismatch")
	}

	// Update payment status and payment ID.
	err = tx.Model(&domain.Payment{}).
		Where("id = ?", payment.ID).
		Updates(map[string]interface{}{
			"status":              "paid",
			"razorpay_payment_id": razorpayPaymentID,
		}).Error
	if err != nil {
		return 0, err
	}

	// Update the order.
	err = tx.Model(&domain.Order{}).
		Where("id = ?", payment.OrderID).
		Updates(map[string]interface{}{
			"payment_status": "paid",
			"order_status":   "placed",
		}).Error
	if err != nil {
		return 0, err
	}

	// Deduct product/variant stock in this same transaction.
	err = ReduceStockAfterPaymentTx(tx, payment.OrderID)
	if err != nil {
		return 0, err
	}

	return payment.UserID, nil
}

// Record a failed payment without overwriting a paid payment.
func ProcessWebhookPaymentFailureTx(
	tx *gorm.DB,
	razorpayOrderID string,
	razorpayPaymentID string,
) error {
	var payment domain.Payment

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&payment).Error
	if err != nil {
		return err
	}

	if payment.Status == "paid" {
		return nil
	}

	err = tx.Model(&domain.Payment{}).
		Where("id = ? AND status <> ?", payment.ID, "paid").
		Updates(map[string]interface{}{
			"status":              "failed",
			"razorpay_payment_id": razorpayPaymentID,
		}).Error
	if err != nil {
		return err
	}

	return tx.Model(&domain.Order{}).
		Where("id = ? AND payment_status <> ?", payment.OrderID, "paid").
		Update("payment_status", "failed").Error
}
