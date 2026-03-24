package repository

import (
	"shoego/database"
	"shoego/domain"
	"shoego/models"
	"strings"
)

//admin product

// insert new product into product table
func CreateProduct(product *domain.Product) error {
	return database.DB.Create(product).Error
}

// create a new images into product_images table
func CreateProductImages(images []domain.ProductImage) error {
	return database.DB.Create(&images).Error
}

func GetProductByID(productID uint) (*domain.Product, error) {
	var product domain.Product
	err := database.DB.Preload("Images").First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// get products with images , desc order and pagination
func GetProducts(limit int, offset int) ([]domain.Product, error) {
	var products []domain.Product

	err := database.DB.Preload("Images").Order("created_at desc").Limit(limit).Offset(offset).Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

// get total num of products in product table
func CountProducts() (int64, error) {
	var count int64
	err := database.DB.Model(&domain.Product{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateProduct(productID uint, data map[string]interface{}) error {
	return database.DB.Model(&domain.Product{}).Where("id = ?", productID).Updates(data).Error
}

// dlt old images of the product and insert new images
func DeleteProductImages(productID uint) error {
	return database.DB.Where("product_id = ?", productID).Delete(&domain.ProductImage{}).Error
}

// find the product by id and delete the product
func SoftDeleteProduct(productID uint) error {
	return database.DB.Delete(&domain.Product{}, productID).Error
}

//admin category management

func CreateCategory(category *domain.Category) error {
	return database.DB.Create(category).Error
}

func GetCategoryByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := database.DB.First(&category, id).Error
	return &category, err
}

func UpdateCategory(category *domain.Category) error {
	return database.DB.Save(category).Error
}

func SoftDeleteCategory(id uint) error {
	return database.DB.Delete(&domain.Category{}, id).Error
}

func IsCategoryExists(name string) (bool, error) {
	var count int64
	err := database.DB.Model(&domain.Category{}).Where("LOWER(name) = LOWER(?)", name).Count(&count).Error

	return count > 0, err
}

func CheckCategoryExistsForUpdate(name string, id uint) (bool, error) {
	//check in any category have same name
	var count int64
	err := database.DB.Model(&domain.Category{}).Where("LOWER(name)=LOWER(?) AND id != ?", name, id).Count(&count).Error

	return count > 0, err
}

func ListCategories(search string, page, limit int) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	query := database.DB.Model(&domain.Category{})

	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit

	//fetch data with desc order and pagination
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&categories).Error

	return categories, total, err
}

//user product

func GetUserProducts(query models.UserProductQuery) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalCount int64

	db := database.DB.Model(&domain.Product{}).
		Preload("Images").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.is_listed = ? AND categories.is_listed = ?", true, true)

	if strings.TrimSpace(query.Search) != "" {
		searchValue := "%" + strings.TrimSpace(query.Search) + "%"
		db = db.Where("LOWER(products.name) LIKE LOWER(?)", searchValue)
	}

	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	if query.BrandID > 0 {
		db = db.Where("brand_id = ?", query.BrandID)
	}

	if query.MinPrice > 0 {
		db = db.Where("price >= ?", query.MinPrice)
	}

	if query.MaxPrice > 0 {
		db = db.Where("price <= ?", query.MaxPrice)
	}

	//sort product
	switch query.Sort {
	case "price_asc":
		db = db.Order("products.price ASC")
	case "price_desc":
		db = db.Order("products.price DESC")
	case "name_asc":
		db = db.Order("products.name ASC")
	case "name_desc":
		db = db.Order("products.name DESC")
	case "new_arrivals":
		db = db.Order("products.created_at DESC")
	default:
		db = db.Order("created_at DESC")
	}

	//count total page
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(offset).Limit(query.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func GetUserProductDetails(productID uint) (*domain.Product, error) {
	//fetch  product details by id
	var product domain.Product

	//err := database.DB.Preload("Images").Preload("Category").Where("id = ? AND is_listed = ?", productID, true).First(&product).Error

	err := database.DB.Model(&domain.Product{}).
		Preload("Images").
		Preload("Category").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.id = ? AND products.is_listed = ? AND categories.is_listed = ?", productID, true, true).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetRelatedUserProducts(categoryID uint, productID uint, limit int) ([]domain.Product, error) {
	//fetching related product same category
	var products []domain.Product

	//err := database.DB.Preload("Images").Where("category_id = ? AND id != ? AND is_listed = ?", categoryID, productID, true).Order("created_at DESC").Limit(limit).Find(&products).Error
	err := database.DB.Model(&domain.Product{}).
		Preload("Images").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.category_id = ? AND products.id != ? AND products.is_listed = ? AND categories.is_listed = ?", categoryID, productID, true, true).
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

//user category

func GetUserCategories(search string) ([]domain.Category, error) {
	var categories []domain.Category
	
	db := database.DB.Model(&domain.Category{}).Where("is_listed = ?", true)

	if strings.TrimSpace(search) != "" {
		searchValue := "%" + strings.TrimSpace(search) + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", searchValue)
	}
	err := db.Order("categories.name ASC").Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}
