package repository

import (
	"errors"
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

	err := database.DB.Joins("JOIN carts ON carts.id = cart_items.cart_id").Preload("Product").Preload("Variant").Preload("Product.Images").Preload("Product.Category").
		Where("carts.user_id = ?", userID).Find(&cartItems).Error

	return cartItems, err
}

func GetAddressByIDAndUserID(addressID, userID uint) (domain.Address, error) {
	var address domain.Address

	err := database.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error

	return address, err
}
func CreateOrderTransaction(
	order *domain.Order,
	orderItems []domain.OrderItem,
	userID uint,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Save Order Items
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
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
		// Clear Cart

		if err := tx.
			Where(
				"cart_id IN (?)",
				tx.Model(&domain.Cart{}).
					Select("id").
					Where("user_id = ?", userID),
			).
			Delete(&domain.CartItem{}).Error; err != nil {

			return err
		}

		return nil
	})
}
