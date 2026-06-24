package repository

import (
	"errors"
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
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

// ADDED
func ReduceStockAfterPayment(orderID uint) error {

	var orderItems []domain.OrderItem

	err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error

	if err != nil {
		return err
	}

	for _, item := range orderItems {

		var product domain.Product

		err := database.DB.First(&product, item.ProductID).Error

		if err != nil {
			return err
		}
		if product.Stock < item.Quantity {
			return errors.New("insufficient stock")
		}
		product.Stock -= item.Quantity

		err = database.DB.
			Save(&product).
			Error

		if err != nil {
			return err
		}
	}

	return nil
}
