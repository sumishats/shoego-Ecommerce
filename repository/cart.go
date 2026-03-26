package repository

import (
	"errors"
	"shoego/database"
	"shoego/domain"

	"gorm.io/gorm"
)

func GetCartByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := database.DB.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func CreateCart(userID uint) (*domain.Cart, error) {
	cart := domain.Cart{
		UserID: userID,
	}
	if err := database.DB.Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

//create cart if not exist for user
func GetOrCreateCart(userID uint) (*domain.Cart, error) {
	cart, err := GetCartByUserID(userID)
	if err == nil {
		return cart, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return CreateCart(userID)
	}

	return nil, err
}

//fetch  product form db by id
func GetProductForCart(productID uint) (*domain.Product, error) {
	var product domain.Product
	err := database.DB.Preload("Category").Preload("Images").First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetCartItem(cartID, productID uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := database.DB.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateCartItem(cartID, productID uint, quantity int) error {
	item := domain.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return database.DB.Create(&item).Error
}

func UpdateCartItemQuantity(cartItemID uint, quantity int) error {
	return database.DB.Model(&domain.CartItem{}).Where("id = ?", cartItemID).Update("quantity", quantity).Error
}

func DeleteCartItem(cartID, productID uint) error {
	return database.DB.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&domain.CartItem{}).Error
}

//fetch all cart items for user with product details
func GetCartItemsByUserID(userID uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	err := database.DB.
		Model(&domain.CartItem{}).
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Images").
		Find(&items).Error

	return items, err
}

func GetCartItemByUserIDAndProductID(userID, productID uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := database.DB.
		Model(&domain.CartItem{}).
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ? AND cart_items.product_id = ?", userID, productID).
		First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}