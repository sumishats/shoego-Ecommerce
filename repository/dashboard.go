package repository

import (
	"shoego/database"
	"shoego/models"
)

func GetAdminBestSellingProducts() ([]models.BestSellingProductResponse, error) {

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

func GetAdminBestSellingCategories() ([]models.BestSellingCategoryResponse, error) {

	var categories []models.BestSellingCategoryResponse

	err := database.DB.Table("order_items").
		Select(`
			categories.id AS category_id,
			categories.name AS category_name,
			SUM(order_items.quantity) AS total_sold
		`).
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.order_status = ?", "delivered").
		Group("categories.id, categories.name").
		Order("total_sold DESC").
		Limit(10).
		Scan(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetAdminSalesChart(filter string) ([]models.SalesChartResponse, error) {

	var sales []models.SalesChartResponse

	var label string
	var groupBy string

	switch filter {

	case "daily":
		label = "TO_CHAR(DATE(created_at), 'DD Mon')"
		groupBy = "TO_CHAR(DATE(created_at), 'DD Mon')"

	case "weekly":
		label = "TO_CHAR(DATE_TRUNC('week', created_at), 'DD Mon YYYY')"
		groupBy = "DATE_TRUNC('week', created_at)"

	case "monthly":
		label = "TO_CHAR(DATE_TRUNC('month', created_at), 'Mon')"
		groupBy = "DATE_TRUNC('month', created_at)"

	case "yearly":
		label = "TO_CHAR(DATE_TRUNC('year', created_at), 'YYYY')"
		groupBy = "DATE_TRUNC('year', created_at)"

	default:
		label = "TO_CHAR(DATE_TRUNC('month', created_at), 'Mon')"
		groupBy = "DATE_TRUNC('month', created_at)"
	}

	err := database.DB.Table("orders").Select(label+" AS label, SUM(final_amount) AS total_sales").Where("order_status = ?", "delivered").
		Group(groupBy).
		Order(groupBy).
		Scan(&sales).Error

	if err != nil {
		return nil, err
	}

	return sales, nil
}
