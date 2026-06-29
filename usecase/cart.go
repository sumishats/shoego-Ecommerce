package usecase

import (
	"errors"
	"shoego/models"
	"shoego/repository"

	"gorm.io/gorm"
)

const MaxCartQuantityPerProduct = 5

func AddToCart(userID uint, productID uint, variantID uint) error {

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

	variant, err := repository.GetVariantByID(variantID)
	if err != nil {
		return errors.New("variant not found")
	}

	if variant.ProductID != productID {
		return errors.New("invalid variant")
	}

	if variant.Stock <= 0 {
		return errors.New("variant out of stock")
	}

	cart, err := repository.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	item, err := repository.GetCartItem(
		cart.ID,
		productID,
		variantID,
	)

	if err == nil {

		newQty := item.Quantity + 1

		if newQty > variant.Stock {
			return errors.New("cannot add more than available stock")
		}

		if newQty > MaxCartQuantityPerProduct {
			return errors.New("maximum quantity limit reached")
		}

		err = repository.UpdateCartItemQuantity(
			item.ID,
			newQty,
		)

		if err != nil {
			return err
		}

		_ = repository.RemoveProductFromWishlist(
			userID,
			productID,
		)

		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = repository.CreateCartItem(
		cart.ID,
		productID,
		variantID,
		1,
	)

	if err != nil {
		return err
	}

	_ = repository.RemoveProductFromWishlist(
		userID,
		productID,
	)

	return nil
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

		stock := product.Stock

		// If this cart item has a variant, use variant stock instead
		if item.VariantID != nil {
			stock = item.Variant.Stock
		}

		if !product.IsListed {

			status = "product unavailable"
			isAvailable = false

		} else if !product.Category.IsListed {

			status = "category unavailable"
			isAvailable = false

		} else if stock <= 0 {

			status = "out of stock"
			isAvailable = false

		} else if item.Quantity > stock {

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

		discountedPrice := product.Price

		offerPercentage := 0.0
		offerName := ""

		productOffer, err := repository.GetActiveProductOffer(product.ID)

		if err == nil {
			offerPercentage = productOffer.DiscountPercentage
			offerName = productOffer.OfferName
		}

		categoryOffer, err := repository.GetActiveCategoryOffer(product.CategoryID)

		if err == nil {

			if categoryOffer.DiscountPercentage > offerPercentage {

				offerPercentage = categoryOffer.DiscountPercentage
				offerName = categoryOffer.OfferName
			}
		}

		if offerPercentage > 0 {
			discountedPrice = product.Price - (product.Price*offerPercentage)/100
		}

		subtotal := discountedPrice * float64(item.Quantity)

		if isAvailable {
			resp.TotalAmount += subtotal
		}
		size := ""

		if item.VariantID != nil {
			size = item.Variant.Size
		}

		resp.Items = append(resp.Items, models.CartItemResponse{
			ProductID:       product.ID,
			VariantID:       item.VariantID,
			Size:            size,
			Name:            product.Name,
			Price:           product.Price,
			DiscountedPrice: discountedPrice,
			OfferPercentage: offerPercentage,
			OfferName:       offerName,
			Quantity:        item.Quantity,
			Stock:           item.Variant.Stock,
			CategoryName:    product.Category.Name,
			Images:          images,
			Subtotal:        subtotal,
			Status:          status,
			IsAvailable:     isAvailable,
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
