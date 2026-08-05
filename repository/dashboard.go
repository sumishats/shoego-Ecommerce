package repository

import (
	"shoego/database"
	"shoego/models"
)

func GetAdminBestSellingProducts() ([]models.BestSellingProductResponse, error){

	var products []models.BestSellingProductResponse

	err := database.DB.Table("order_items").
		Select(`
			products.id AS product_id,
			products.name AS product_name,
			SUM(order_items.quantity) AS total_sold
		`).
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.order_status = ?", "delivered").
		Group("products.id, products.name").
		Order("total_sold DESC").
		Limit(10).Scan(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}