package usecase

import (
	"errors"
	"math"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"strconv"
	"strings"
)

func GetAdminOrders(page, limit int, search, status, sortBy, date string) (*models.AdminOrderPaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	orders, totalCount, err := repository.GetAdminOrders(search, status, sortBy, date, limit, offset)
	if err != nil {
		return nil, err
	}

	var orderResponse []models.AdminOrderListResponse

	for _, order := range orders {
		orderResponse = append(orderResponse, models.AdminOrderListResponse{
			ID:            order.ID,
			OrderID:       "ORD" + strconv.FormatUint(uint64(order.ID), 10), //uint to string like orderid
			Date:          order.CreatedAt,
			UserName:      order.User.Name,
			UserEmail:     order.User.Email,
			UserPhone:     order.User.Phone,
			ItemCount:     len(order.OrderItems),
			FinalAmount:   order.FinalAmount,
			OrderStatus:   order.OrderStatus,
			PaymentMethod: order.PaymentMethod,
			PaymentStatus: order.PaymentStatus,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &models.AdminOrderPaginationResponse{
		Orders:     orderResponse,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func GetAdminOrderDetail(orderID uint) (*models.AdminOrderDetailResponse, error) {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	var items []models.AdminOrderItemResponse
	for _, item := range order.OrderItems {
		image := ""
		if len(item.Product.Images) > 0 {
			image = item.Product.Images[0].ImageURL
		}
		items = append(items, models.AdminOrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    float64(item.Quantity) * item.Price,
			Image:       image,
		})
	}

	return &models.AdminOrderDetailResponse{
		ID:            order.ID,
		OrderID:       "ORD" + strconv.FormatUint(uint64(order.ID), 10),
		Date:          order.CreatedAt,
		OrderStatus:   order.OrderStatus,
		FinalAmount:   order.FinalAmount,
		PaymentMethod: order.PaymentMethod,
		PaymentStatus: order.PaymentStatus,
		ItemCount:     len(order.OrderItems),
		User: models.AdminOrderUserResponse{
			ID:    order.User.ID,
			Name:  order.User.Name,
			Email: order.User.Email,
			Phone: order.User.Phone,
		},
		Address: models.AdminOrderAddressResponse{
			Name:      order.Address.Name,
			Phone:     order.Address.Phone,
			HouseName: order.Address.HouseName,
			Street:    order.Address.Street,
			City:      order.Address.City,
			State:     order.Address.State,
			Pincode:   order.Address.Pincode,
		},
		Items: items,
	}, nil
}

func ChangeOrderStatus(orderID uint, status string) error {
	status = strings.TrimSpace(strings.ToLower(status))

	validStatuses := map[string]bool{
		"pending":          true,
		"shipped":          true,
		"out_for_delivery": true,
		"delivered":        true,
		"cancelled":        true,
		"returned":         true,
		"return_requested": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid order status")
	}
	if status == "returned" {

		order, err := repository.GetOrderByID(orderID)
		if err != nil {
			return err
		}
	
		if order.OrderStatus != "returned" {
	
			err = CreditWallet(
				order.UserID,
				order.FinalAmount,
				"Return refund",
			)
	
			if err != nil {
				return err
			}
		}
	}
	return repository.UpdateOrderStatus(orderID, status)
}

//intentory manage admin

func GetAdminInventory(page, limit int, search, stockFilter, sortBy string) (*models.AdminInventoryPaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, totalCount, err := repository.GetAdminInventory(search, stockFilter, sortBy, limit, offset)
	if err != nil {
		return nil, err
	}

	var productResponses []models.AdminInventoryResponse

	for _, product := range products {
		stockStatus := "in stock"
		if product.Stock == 0 {
			stockStatus = "out of stock"
		} else if product.Stock <= 5 {
			stockStatus = "low stock"
		}

		productResponses = append(productResponses, models.AdminInventoryResponse{
			ProductID:    product.ID,
			ProductName:  product.Name,
			SKU:          product.SKU,
			CategoryName: product.Category.Name,
			Stock:        product.Stock,
			Price:        product.Price,
			IsListed:     product.IsListed,
			StockStatus:  stockStatus,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &models.AdminInventoryPaginationResponse{
		Products:   productResponses,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func ChangeProductStock(productID uint, stock int) error {
	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return repository.UpdateProductStock(productID, stock)
}

//user order managemnt

func GetUserOrderList(userID uint, search string, page, limit int) ([]models.OrderListItemResponse, error) {
	offset := (page - 1) * limit
	orders, err := repository.GetUserOrders(userID, search, limit, offset)
	if err != nil {
		return nil, err
	}
	var response []models.OrderListItemResponse
	for _, order := range orders {
		response = append(response, models.OrderListItemResponse{
			OrderID:       order.OrderID,
			OrderStatus:   order.OrderStatus,
			OrderDate:     order.CreatedAt,
			FinalAmount:   order.FinalAmount,
			PaymentMethod: order.PaymentMethod,
		})
	}
	return response, nil
}

func GetUserOrderDetail(userID uint, orderID string) (*models.OrderDetailResponse, error) {
	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return nil, err
	}

	var items []models.OrderItemResponse
	for _, item := range order.OrderItems {
		image := ""
		if len(item.Product.Images) > 0 {
			image = item.Product.Images[0].ImageURL
		}

		items = append(items, models.OrderItemResponse{
			ItemID:      item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Image:       image,
			Quantity:    item.Quantity,
			Price:       item.Price,
			TotalPrice:  item.TotalPrice,
			ItemStatus:  item.ItemStatus,
		})
	}

	return &models.OrderDetailResponse{
		OrderID:        order.OrderID,
		OrderStatus:    order.OrderStatus,
		OrderDate:      order.CreatedAt,
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		Subtotal:       order.Subtotal,
		TaxAmount:      order.TaxAmount,
		CouponDiscount: order.CouponDiscount,
		ShippingCharge: order.ShippingCharge,
		FinalAmount:    order.FinalAmount,
		Items:          items,
	}, nil
}

func CancelOrder(userID uint, orderID, reason string) error {
	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return err
	}

	if order.OrderStatus == "cancelled" {
		return errors.New("order already cancelled")
	}

	if order.OrderStatus == "delivered" {
		return errors.New("delivered order cannot be cancelled")
	}

	if order.OrderStatus == "returned" {
		return errors.New("returned order cannot be cancelled")
	}

	for i := range order.OrderItems {
		item := &order.OrderItems[i]

		if item.ItemStatus == "cancelled" || item.ItemStatus == "returned" {
			continue
		}

		item.ItemStatus = "cancelled"
		item.CancellationReason = reason

		if err := repository.UpdateOrderItem(item); err != nil {
			return err
		}

		if err := repository.IncrementProductStock(item.ProductID, item.Quantity); err != nil {
			return err
		}

	}

	payment, err := repository.GetPaymentByOrderID(order.ID)

	if err == nil && payment.Status == "paid" {

		err = CreditWallet(
			order.UserID,
			order.FinalAmount,
			"Order cancellation refund",
		)

		if err != nil {
			return err
		}
	}

	order.OrderStatus = "cancelled"
	order.CancellationReason = reason

	return repository.UpdateOrder(&order)
}

func CancelOrderItem(userID uint, orderID string, itemID uint, reason string) error {
	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return err
	}

	if order.OrderStatus == "delivered" {
		return errors.New("delivered order item cannot be cancelled")
	}

	if order.OrderStatus == "returned" {
		return errors.New("returned order item cannot be cancelled")
	}

	item, err := repository.GetOrderItemByID(order.ID, itemID)
	if err != nil {
		return err
	}

	if item.ItemStatus == "cancelled" {
		return errors.New("item already cancelled")
	}

	if item.ItemStatus == "returned" {
		return errors.New("returned item cannot be cancelled")
	}

	item.ItemStatus = "cancelled"
	item.CancellationReason = reason

	if err := repository.UpdateOrderItem(&item); err != nil {
		return err
	}

	if err := repository.IncrementProductStock(item.ProductID, item.Quantity); err != nil {
		return err
	}

	allCancelled := true
	for _, orderItem := range order.OrderItems {
		if orderItem.ID == item.ID {
			continue
		}
		if orderItem.ItemStatus != "cancelled" {
			allCancelled = false
			break
		}
	}

	if allCancelled {
		order.OrderStatus = "cancelled"
	} else {
		order.OrderStatus = "partially_cancelled"
	}

	return repository.UpdateOrder(&order)
}

func ReturnOrder(userID uint, orderID, reason string) error {
	if reason == "" {
		return errors.New("return reason is required")
	}

	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return err
	}

	if order.OrderStatus != "delivered" {
		return errors.New("only delivered orders can be returned")
	}

	for i := range order.OrderItems {
		item := &order.OrderItems[i]

		if item.ItemStatus == "cancelled" {
			continue
		}

		item.ItemStatus = "return_requested"
		item.ReturnReason = reason
		// err := repository.IncrementProductStock(
		// 	item.ProductID,
		// 	item.Quantity,
		// )
		// if err != nil {
		// 	return err
		// }

		// item.ItemStatus = "returned"
		// item.ReturnReason = reason

		if err := repository.UpdateOrderItem(item); err != nil {
			return err
		}
	}

	order.OrderStatus = "return_requested"
	order.ReturnReason = reason

	return repository.UpdateOrder(&order)
}

func ReturnOrderItem(userID uint, orderID string, itemID uint, reason string) error {
	if reason == "" {
		return errors.New("return reason is required")
	}

	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return err
	}

	if order.OrderStatus != "delivered" && order.OrderStatus != "partially_returned" {
		return errors.New("only delivered order items can be returned")
	}

	item, err := repository.GetOrderItemByID(order.ID, itemID)
	if err != nil {
		return err
	}

	if item.ItemStatus == "returned" {
		return errors.New("item already returned")
	}

	if item.ItemStatus == "cancelled" {
		return errors.New("cancelled item cannot be returned")
	}

	item.ItemStatus = "return_requested"
	item.ReturnReason = reason

	// err = repository.IncrementProductStock(
	// 	item.ProductID,
	// 	item.Quantity,
	// )
	// if err != nil {
	// 	return err
	// }

	// item.ItemStatus = "returned"
	// item.ReturnReason = reason

	if err := repository.UpdateOrderItem(&item); err != nil {
		return err
	}

	items, err := repository.GetOrderItemsByOrderID(order.ID)
	if err != nil {
		return err
	}

	allReturned := true
	for _, orderItem := range items {
		if orderItem.ItemStatus != "returned" {
			allReturned = false
			break
		}
	}

	if allReturned {
		order.OrderStatus = "return_requested"
		order.ReturnReason = reason
	} else {
		order.OrderStatus = "partially_returned"
	}

	return repository.UpdateOrder(&order)
}

func GenerateInvoice(userID uint, orderID string) ([]byte, error) {
	order, err := repository.GetOrderByOrderID(userID, orderID)
	if err != nil {
		return nil, err
	}

	if order.OrderStatus == "cancelled" {
		return nil, errors.New("invoice not available for cancelled order")
	}

	return helper.GenerateInvoicePDF(order)
}
