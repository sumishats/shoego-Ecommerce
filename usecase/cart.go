package usecase

import (
	"errors"
	"shoego/models"
	"shoego/repository"

	"gorm.io/gorm"
)

const MaxCartQuantityPerProduct = 5

func AddToCart(userID, productID uint) error {
	product, err := repository.GetProductForCart(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if !product.IsListed {
		return errors.New("product is unavailable")
	}

	if !product.Category.IsListed {
		return errors.New("category is unavailable")
	}

	if product.Stock <= 0 {
		return errors.New("product is out of stock")
	}

	cart, err := repository.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	//check if item already in cart
	item, err := repository.GetCartItem(cart.ID, productID)
	if err == nil {
		newQty := item.Quantity + 1

		if newQty > product.Stock {
			return errors.New("cannot add more than available stock")
		}

		if newQty > MaxCartQuantityPerProduct {
			return errors.New("maximum quantity limit reached")
		}

		return repository.UpdateCartItemQuantity(item.ID, newQty)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return repository.CreateCartItem(cart.ID, productID, 1)
}

func GetCart(userID uint) (*models.CartResponse, error) {
	items, err := repository.GetCartItemsByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := &models.CartResponse{
		Items:           []models.CartItemResponse{},
		TotalAmount:     0,
		CheckoutAllowed: true,
	}

	for _, item := range items {
		product := item.Product

		status := "available"
		isAvailable := true

		if !product.IsListed {
			status = "product unavailable"
			isAvailable = false
		} else if !product.Category.IsListed {
			status = "category unavailable"
			isAvailable = false
		} else if product.Stock <= 0 {
			status = "out of stock"
			isAvailable = false
		} else if item.Quantity > product.Stock {
			status = "quantity exceeds stock"
			isAvailable = false
		}

		if !isAvailable {
			resp.CheckoutAllowed = false
		}

		var images []string
		for _, img := range product.Images {
			images = append(images, img.ImageURL)
		}

		//caculate subtotal for one cart item
		subtotal := float64(item.Quantity) * product.Price
		if isAvailable {
			resp.TotalAmount += subtotal
		}

		resp.Items = append(resp.Items, models.CartItemResponse{
			ProductID:    product.ID,
			Name:         product.Name,
			Price:        product.Price,
			Quantity:     item.Quantity,
			Stock:        product.Stock,
			CategoryName: product.Category.Name,
			Images:       images,
			Subtotal:     subtotal,
			Status:       status,
			IsAvailable:  isAvailable,
		})
	}

	return resp, nil
}

func UpdateCartQuantity(userID, productID uint, action string) error {
	item, err := repository.GetCartItemByUserIDAndProductID(userID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cart item not found")
		}
		return err
	}

	product, err := repository.GetProductForCart(productID)
	if err != nil {
		return err
	}

	switch action {
	case "increment":
		newQty := item.Quantity + 1

		if !product.IsListed || !product.Category.IsListed {
			return errors.New("product is unavailable")
		}

		if product.Stock <= 0 {
			return errors.New("product is out of stock")
		}

		if newQty > product.Stock {
			return errors.New("cannot add more than available stock")
		}

		if newQty > MaxCartQuantityPerProduct {
			return errors.New("maximum quantity limit reached")
		}

		return repository.UpdateCartItemQuantity(item.ID, newQty)

	case "decrement":
		newQty := item.Quantity - 1

		if newQty <= 0 {
			return repository.DeleteCartItem(item.CartID, productID)
		}

		return repository.UpdateCartItemQuantity(item.ID, newQty)

	default:
		return errors.New("invalid action")
	}
}

func RemoveCartItem(userID, productID uint) error {
	item, err := repository.GetCartItemByUserIDAndProductID(userID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cart item not found")
		}
		return err
	}

	return repository.DeleteCartItem(item.CartID, productID)
}

func ValidateCartForCheckout(userID uint) error {
	cart, err := GetCart(userID)
	if err != nil {
		return err
	}

	if len(cart.Items) == 0 {
		return errors.New("cart is empty")
	}

	if !cart.CheckoutAllowed {
		return errors.New("some items in the cart are unavailable")
	}

	return nil
}
