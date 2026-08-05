package models

type BestSellingProductResponse struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}