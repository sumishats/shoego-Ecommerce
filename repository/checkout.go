package repository

import (
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
)

func GetCheckoutAddresses(userID uint) ([]domain.Address, error) {
	var addresses []domain.Address

	err := database.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error

	return addresses, err
}

func GetCartItemsForCheckout(userID uint) ([]domain.CartItem, error) {
	var cartItems []domain.CartItem

	err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").Preload("Product").Preload("Product.Images").Preload("Product.Category").
		Where("carts.user_id = ?", userID).Find(&cartItems).Error

	return cartItems, err
}

func GetAddressByIDAndUserID(addressID, userID uint) (domain.Address, error) {
	var address domain.Address

	err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error

	return address, err
}

func CreateOrderWithItems(order *domain.Order, orderItems []domain.OrderItem) error {
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

		return nil
	})
}

func ReduceProductStock(productID uint, quantity int) error {
	return database.DB.Model(&domain.Product{}).Where("id = ? AND stock >= ?", productID, quantity).Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

func ClearCartAfterOrder(userID uint) error {
	return database.DB.Where("cart_id IN (?)",database.DB.Model(&domain.Cart{}).Select("id").Where("user_id = ?", userID),).Delete(&domain.CartItem{}).Error
}


