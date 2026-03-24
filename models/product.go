package models

// admin side product and category
type AddProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BrandID     uint    `json:"brand_id"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
}

type EditProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BrandID     uint    `json:"brand_id"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  uint    `json:"category_id"`
}

type ProductResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BrandID     uint     `json:"brand_id"`
	SKU         string   `json:"sku"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	CategoryID  uint     `json:"category_id"`
	IsListed    bool     `json:"is_listed"`
	Images      []string `json:"images"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalCount int64             `json:"total_count"`
	TotalPages int               `json:"total_pages"`
}

type AddCategory struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategory struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsListed    bool   `json:"is_listed"`
}

type CategoryListResponse struct {
	Categories  []CategoryResponse `json:"categories"`
	TotalCount  int64              `json:"total_count"`
	CurrentPage int                `json:"current_page"`
	Limit       int                `json:"limit"`
	TotalPages  int                `json:"total_pages"`
}

//user side product and category

type UserProductQuery struct {
	Search     string
	Sort       string
	CategoryID uint
	BrandID    uint
	MinPrice   float64
	MaxPrice   float64
	Page       int
	Limit      int
}

type UserProductResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	BrandID     uint     `json:"brand_id"`
	SKU         string   `json:"sku"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	CategoryID  uint     `json:"category_id"`
	IsListed    bool     `json:"is_listed"`
	Images      []string `json:"images"`
}

type UserProductListResponse struct {
	Products   []UserProductResponse `json:"products"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalCount int64                 `json:"total_count"`
	TotalPages int                   `json:"total_pages"`
}

type UserProductDetailResponse struct {
	ID              uint                  `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	BrandID         uint                  `json:"brand_id"`
	SKU             string                `json:"sku"`
	Price           float64               `json:"price"`
	Stock           int                   `json:"stock"`
	CategoryID      uint                  `json:"category_id"`
	CategoryName    string                `json:"category_name"`
	IsListed        bool                  `json:"is_listed"`
	Images          []string              `json:"images"`
	Status          string                `json:"status"`
	Breadcrumbs     []string              `json:"breadcrumbs"`
	Highlights      []string              `json:"highlights"`
	RelatedProducts []UserProductResponse `json:"related_products"`
}

type UserCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
