package repository

import (
	"shoego/database"
	"shoego/domain"
	"strings"

	"gorm.io/gorm"
)

func GetAdminOrders(search, status, sortBy, date string, limit, offset int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var totalCount int64

	query := database.DB.Model(&domain.Order{}).Preload("User").Preload("OrderItems")

	if search != "" {
		search = strings.ToLower(search)
		query = query.Joins("JOIN users ON users.id = orders.user_id").Where("LOWER(users.name) LIKE ? OR LOWER(users.email) LIKE ? OR CAST(orders.id AS TEXT) LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		status = strings.ToLower(status)
		query = query.Where("orders.order_status = ?", status)
	}

	if date != "" {
		query = query.Where("DATE(orders.created_at)=?", date)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	switch sortBy {
	case "asc":
		query = query.Order("orders.created_at ASC")
	case "status_asc":
		query = query.Order("orders.order_status ASC")
	case "status_desc":
		query = query.Order("orders.order_status DESC")
	default:
		query = query.Order("orders.created_at DESC")
	}

	err = query.Limit(limit).Offset(offset).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func GetOrderByID(orderID uint) (*domain.Order, error) {
	var order domain.Order

	err := database.DB.Preload("User").Preload("Address").Preload("OrderItems").Preload("OrderItems.Product").Preload("OrderItems.Product.Images").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrderStatus(orderID uint, status string) error {
	var order domain.Order

	if err := database.DB.Preload("OrderItems").First(&order, orderID).Error; err != nil {
		return err
	}

	if status == "returned" {

		// Stock return
		for _, item := range order.OrderItems {

			var product domain.Product

			if err := database.DB.First(&product, item.ProductID).Error; err != nil {
				return err
			}

			product.Stock += item.Quantity

			if err := database.DB.Save(&product).Error; err != nil {
				return err
			}
		}

		// NEW CODE START
		err := database.DB.Model(&domain.OrderItem{}).
			Where("order_id = ? AND item_status = ?", order.ID, "return_requested").
			Update("item_status", "returned").Error

		if err != nil {
			return err
		}

	}

	updates := map[string]interface{}{
		"order_status": status,
	}

	if order.PaymentMethod == "cod" {
		switch status {
		case "delivered":
			updates["payment_status"] = "paid"

		case "pending", "shipped", "out_for_delivery":
			updates["payment_status"] = "pending"

		case "cancelled":
			updates["payment_status"] = "cancelled"

		case "returned":
			updates["payment_status"] = "refunded"
		}
	}

	return database.DB.Model(&domain.Order{}).
		Where("id = ?", orderID).
		Updates(updates).Error
}

// inventory admin

func GetAdminInventory(search, stockFilter, sortBy string, limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalCount int64

	query := database.DB.Model(&domain.Product{}).Preload("Category")

	if search != "" {
		search = strings.ToLower(search)
		query = query.
			Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(products.name) LIKE ? OR LOWER(products.sku) LIKE ? OR LOWER(categories.name) LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	switch stockFilter {
	case "out_of_stock":
		query = query.Where("products.stock = ?", 0)
	case "low_stock":
		query = query.Where("products.stock > ? AND products.stock <= ?", 0, 5)
	case "in_stock":
		query = query.Where("products.stock > ?", 0)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	switch sortBy {
	case "stock_asc":
		query = query.Order("products.stock ASC")
	case "stock_desc":
		query = query.Order("products.stock DESC")
	case "name_asc":
		query = query.Order("products.name ASC")
	case "name_desc":
		query = query.Order("products.name DESC")
	default:
		query = query.Order("products.created_at DESC")
	}

	err = query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func UpdateProductStock(productID uint, stock int) error {
	return database.DB.Model(&domain.Product{}).Where("id = ?", productID).Update("stock", stock).Error
}

// user order managemnt
func GetUserOrders(userID uint, search string, limit, offset int) ([]domain.Order, error) {
	var orders []domain.Order
	query := database.DB.Where("user_id=?", userID).Order("created_at DESC")

	if search != "" {
		query = query.Where("order_id ILIKE ? OR order_status ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

//	func GetOrderByOrderID(userID uint, orderID string) (domain.Order, error) {
//		var order domain.Order
//		err := database.DB.Preload("OrderItems.Product.Images").Preload("Address").Where(" user_id AND order_id = ?", order.UserID, orderID).First(&order).Error
//		return order, err
//	}
func GetOrderByOrderID(userID uint, orderID string) (domain.Order, error) {
	var order domain.Order

	err := database.DB.Preload("OrderItems.Product.Images").Preload("Address").Where("user_id = ? AND order_id = ?", userID, orderID).First(&order).Error
	return order, err
}

func GetOrderByOrderIDPayment(orderID string) (*domain.Order, error) {
	var order domain.Order

	err := database.DB.Where("order_id = ?", orderID).First(&order).Error
	return &order, err
}

func UpdateOrder(order *domain.Order) error {
	return database.DB.Save(order).Error
}

func UpdateOrderItem(item *domain.OrderItem) error {
	return database.DB.Save(item).Error
}

func GetOrderItemByID(orderID uint, itemID uint) (domain.OrderItem, error) {
	var item domain.OrderItem
	err := database.DB.Where("order_id = ? AND id = ?", orderID, itemID).First(&item).Error
	return item, err
}

func IncrementProductStock(productID uint, qty int) error {
	return database.DB.Model(&domain.Product{}).Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", qty)).Error
}

func GetOrderItemsByOrderID(orderID uint) ([]domain.OrderItem, error) {
	var items []domain.OrderItem

	err := database.DB.Where("order_id = ?", orderID).Find(&items).Error

	return items, err
}
