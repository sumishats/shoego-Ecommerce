package repository

import (
	"errors"
	"fmt"
	"shoego/database"
	"shoego/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	err := database.DB.Preload("User").Preload("Address").Preload("OrderItems.Product").Preload("OrderItems.Product.Images").Preload("OrderItems.Variant").First(&order, orderID).Error

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

		for _, item := range order.OrderItems {

			if item.VariantID == nil {

				var product domain.Product

				if err := database.DB.First(&product, item.ProductID).Error; err != nil {
					return err
				}

				product.Stock += item.Quantity

				if err := database.DB.Save(&product).Error; err != nil {
					return err
				}

			} else {

				var variant domain.ProductVariant

				if err := database.DB.First(&variant, *item.VariantID).Error; err != nil {
					return err
				}

				variant.Stock += item.Quantity

				if err := database.DB.Save(&variant).Error; err != nil {
					return err
				}
			}
		}
		if err := database.DB.Model(&domain.OrderItem{}).Where("order_id = ? AND item_status = ?", order.ID, "return_requested").Update("item_status", "returned").Error; err != nil {
			return err
		}
	}

	updates := map[string]interface{}{
		"order_status": status,
	}

	switch order.PaymentMethod {

	case "cod":

		switch status {

		case "pending", "shipped", "out_for_delivery":
			updates["payment_status"] = "pending"

		case "delivered":
			updates["payment_status"] = "paid"

		case "cancelled":
			updates["payment_status"] = "cancelled"

		case "returned":
			updates["payment_status"] = "refunded"
		}

	case "razorpay", "wallet":

		switch status {

		case "pending":
			updates["payment_status"] = "pending"

		case "delivered":
			updates["payment_status"] = "paid"

		case "cancelled":
			updates["payment_status"] = "refunded"

		case "returned":
			updates["payment_status"] = "refunded"
		}

	}
	return database.DB.Model(&domain.Order{}).Where("id = ?", orderID).Updates(updates).Error
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

func GetOrderByOrderID(userID uint, orderID string) (domain.Order, error) {
	var order domain.Order

	err := database.DB.
		Preload("Address").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Images").
		Preload("OrderItems.Variant").
		Where("user_id = ? AND order_id = ?", userID, orderID).
		First(&order).Error

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

func CancelOrderItemTransaction(userID uint,orderID string,itemID uint,reason string,) error {

    return database.DB.Transaction(func(tx *gorm.DB) error {

        var order domain.Order

        err := tx.Clauses(clause.Locking{Strength: "UPDATE"}). Preload("OrderItems").
            Where("user_id = ? AND order_id = ?", userID, orderID).First(&order).Error

        if err != nil {
            return err
        }

        if order.OrderStatus == "delivered" ||
            order.OrderStatus == "returned" {
            return errors.New("this order cannot be cancelled")
        }

        
        var item domain.OrderItem

        err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND order_id = ?", itemID, order.ID).First(&item).Error

        if err != nil {
            return err
        }

        if item.ItemStatus == "cancelled" {
            return errors.New("item already cancelled")
        }

        if item.ItemStatus == "returned" {
            return errors.New("returned item cannot be cancelled")
        }

        //restore product stock 
        result := tx.Model(&domain.Product{}).
            Where("id = ?", item.ProductID).UpdateColumn("stock",gorm.Expr("stock + ?", item.Quantity),)

        if result.Error != nil {
            return result.Error
        }

        if result.RowsAffected == 0 {
            return errors.New("product not found")
        }

        if item.VariantID != nil {

        result = tx.Model(&domain.ProductVariant{}).Where("id = ?", *item.VariantID).UpdateColumn("stock",gorm.Expr("stock + ?", item.Quantity),)

            if result.Error != nil {
                return result.Error
            }

            if result.RowsAffected == 0 {
                return errors.New("product variant not found")
            }
        }

        
        err = tx.Model(&domain.OrderItem{}).Where("id = ?", item.ID).
            Updates(map[string]interface{}{
                "item_status":        "cancelled",
                "cancellation_reason": reason,
            }).Error

        if err != nil {
            return err
        }

        allCancelled := true
        totalItemValue := 0.0

        for _, orderItem := range order.OrderItems {

            if orderItem.ItemStatus == "cancelled" {
                continue
            }

            totalItemValue += orderItem.TotalPrice

            if orderItem.ID != item.ID {
                allCancelled = false
            }
        }

        if allCancelled {
            order.OrderStatus = "cancelled"
        } else {
            order.OrderStatus = "partially_cancelled"
        }

        shouldRefund := order.PaymentStatus == "paid" && (order.PaymentMethod == "razorpay" ||order.PaymentMethod == "wallet")

        if shouldRefund {

            var wallet domain.Wallet

            err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error

            if errors.Is(err, gorm.ErrRecordNotFound) {

                wallet = domain.Wallet{
                    UserID:  userID,
                    Balance: 0,
                }

                if err = tx.Create(&wallet).Error; err != nil {
                    return err
                }

            } else if err != nil {
                return err
            }

            var alreadyRefunded float64

            err = tx.Table("wallet_transactions").Joins("JOIN wallets ON wallets.id = wallet_transactions.wallet_id",).
                Where("wallets.user_id = ? AND wallet_transactions.type = ? AND wallet_transactions.description LIKE ?",userID,"credit","Order Cancel refund: "+order.OrderID+"%",).
                Select("COALESCE(SUM(wallet_transactions.amount), 0)",).Scan(&alreadyRefunded).Error

            if err != nil {
                return err
            }

            remainingAmount := order.FinalAmount - alreadyRefunded

            if remainingAmount < 0 {
                remainingAmount = 0
            }

            refundAmount := 0.0

            if allCancelled {
                refundAmount = remainingAmount
            } else if totalItemValue > 0 {
                refundAmount = order.FinalAmount *
                    item.TotalPrice / totalItemValue

                if refundAmount > remainingAmount {
                    refundAmount = remainingAmount
                }
            }

            if refundAmount > 0 {

                err = tx.Model(&domain.Wallet{}).Where("id = ?", wallet.ID).UpdateColumn("balance",gorm.Expr("balance + ?", refundAmount),).Error

                if err != nil {
                    return err
                }

                transaction := domain.WalletTransaction{
                    WalletID: wallet.ID,
                    Amount:   refundAmount,
                    Type:     "credit",
                    Description: fmt.Sprintf(
                        "Order Cancel refund: %s item %d",
                        order.OrderID,
                        item.ID,
                    ),
                }

                if err = tx.Create(&transaction).Error; err != nil {
                    return err
                }
            }
        }

        updates := map[string]interface{}{
            "order_status": order.OrderStatus,
        }

        if allCancelled {
            if shouldRefund {
                updates["payment_status"] = "refunded"
            } else if order.PaymentStatus != "paid" {
                updates["payment_status"] = "cancelled"
            }
        }

        return tx.Model(&domain.Order{}).Where("id = ?", order.ID).Updates(updates).Error
    })
}