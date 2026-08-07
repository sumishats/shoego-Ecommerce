package models

type BestSellingProductResponse struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}

type BestSellingCategoryResponse struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	TotalSold    int    `json:"total_sold"`
}
type BestSellingBrandResponse struct {
	BrandID   uint  `json:"brand_id"`
	TotalSold int64 `json:"total_sold"`
}
type SalesChartResponse struct {
	Label      string  `json:"label"`
	TotalSales float64 `json:"total_sales"`
}