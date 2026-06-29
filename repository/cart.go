package repository

import (
	"errors"
	"fmt"
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

// create cart if not exist for user
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

// fetch  product form db by id
func GetProductForCart(productID uint) (*domain.Product, error) {
	var product domain.Product
	err := database.DB.Preload("Category").Preload("Images").First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetCartItem(cartID, productID, variantID uint) (*domain.CartItem, error) {

	var item domain.CartItem

	err := database.DB.
		Where(
			"cart_id = ? AND product_id = ? AND variant_id = ?",
			cartID,
			productID,
			&variantID,
		).
		First(&item).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func CreateCartItem(cartID, productID, variantID uint, quantity int) error {

	item := domain.CartItem{
		CartID:    cartID,
		ProductID: productID,
		VariantID: &variantID,
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

// fetch all cart items for user with product details
func GetCartItemsByUserID(userID uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	// err := database.DB.Model(&domain.CartItem{}).Joins("JOIN carts ON carts.id = cart_items.cart_id").Where("carts.user_id = ?", userID).Preload("Product").Preload("Variant").Preload("Product.Category").Preload("Product.Images").Find(&items).Error
	err := database.DB.
		Model(&domain.CartItem{}).
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Preload("Product").
		Preload("Variant").
		Preload("Product.Category").
		Preload("Product.Images").
		Find(&items).Error
	return items, err
}

func GetCartItemByUserIDAndProductID(userID, productID uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := database.DB.
		Model(&domain.CartItem{}).Joins("JOIN carts ON carts.id = cart_items.cart_id").Where("carts.user_id = ? AND cart_items.product_id = ?", userID, productID).First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func ClearCartByUserID(userID uint) error {

	tx := database.DB.Begin()

	var cart domain.Cart

	if err := tx.
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.
		Where("cart_id = ?", cart.ID).
		Delete(&domain.CartItem{}).Error; err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.
		Delete(&cart).Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func GetCartSubtotal(userID uint) (float64, error) {

	cartItems, err := GetCartItemsByUserID(userID)
	if err != nil {
		return 0, err
	}

	var subtotal float64

	for _, item := range cartItems {

		price := item.Product.Price

		bestOffer := 0.0

		// Product offer
		productOffer, err := GetActiveProductOffer(item.ProductID)
		if err == nil {
			bestOffer = productOffer.DiscountPercentage
		}

		// Category offer
		categoryOffer, err := GetActiveCategoryOffer(item.Product.CategoryID)
		if err == nil && categoryOffer.DiscountPercentage > bestOffer {
			bestOffer = categoryOffer.DiscountPercentage
		}

		if bestOffer > 0 {
			price = price - (price * bestOffer / 100)
		}

		subtotal += price * float64(item.Quantity)
	}

	return subtotal, nil
}

func UpdateCart(cart *domain.Cart) error {
	fmt.Println("upadte db")
	return database.DB.Save(cart).Error
}
