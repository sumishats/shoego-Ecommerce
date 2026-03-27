package usecase

import (
	"errors"
	"math"
	"mime/multipart"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"strconv"
)

// admin product
func AddProduct(req models.AddProductRequest, files []*multipart.FileHeader) error {

	if len(files) < 3 {
		return errors.New("minimum 3 images required")
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		BrandID:     req.BrandID,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		IsListed:    true,
	}

	err := repository.CreateProduct(product)
	if err != nil {
		return err
	}

	var images []domain.ProductImage

	for _, file := range files {
		path, err := helper.SaveProductImage(file, "./images")
		if err != nil {
			return err
		}

		images = append(images, domain.ProductImage{
			ProductID: product.ID,
			ImageURL:  path,
		})
	}

	return repository.CreateProductImages(images)
}

func EditProduct(productID uint, req models.EditProductRequest, files []*multipart.FileHeader) error {

	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}

	err = repository.UpdateProduct(productID, map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"brand_id":    req.BrandID,
		"sku":         req.SKU,
		"price":       req.Price,
		"stock":       req.Stock,
		"category_id": req.CategoryID,
	})
	if err != nil {
		return err
	}

	// if new images are uploaded, replace old images
	if len(files) > 0 {
		for _, img := range product.Images {
			helper.DeleteFileIfExists(img.ImageURL)
		}

		err = repository.DeleteProductImages(productID)
		if err != nil {
			return err
		}

		var images []domain.ProductImage

		for _, file := range files {
			path, err := helper.SaveProductImage(file, "./images")
			if err != nil {
				return err
			}

			images = append(images, domain.ProductImage{
				ProductID: productID,
				ImageURL:  path,
			})
		}

		err = repository.CreateProductImages(images)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteProduct(id uint) error {
	return repository.SoftDeleteProduct(id)
}

func GetProducts(pageStr, limitStr string) (*models.ProductListResponse, error) {

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	products, err := repository.GetProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := repository.CountProducts()
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	var resp []models.ProductResponse

	for _, p := range products {

		var images []string
		for _, img := range p.Images {
			images = append(images, img.ImageURL)
		}

		//response for each product
		resp = append(resp, models.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			BrandID:     p.BrandID,
			SKU:         p.SKU,
			Price:       p.Price,
			Stock:       p.Stock,
			CategoryID:  p.CategoryID,
			IsListed:    p.IsListed,
			Images:      images,
		})
	}

	return &models.ProductListResponse{
		Products:   resp,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

//admin category managent

func AddCategory(data models.AddCategory) error {
	exists, _ := repository.IsCategoryExists(data.Name)
	if exists {
		return errors.New("category already exists")
	}

	category := domain.Category{
		Name:        data.Name,
		Description: data.Description,
	}

	return repository.CreateCategory(&category)
}

func EditCategory(id uint, data models.UpdateCategory) error {
	category, err := repository.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	exists, _ := repository.CheckCategoryExistsForUpdate(data.Name, id)
	if exists {
		return errors.New("category already exists")
	}

	category.Name = data.Name
	category.Description = data.Description

	return repository.UpdateCategory(category)
}

func DeleteCategory(id uint) error {
	_, err := repository.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return repository.SoftDeleteCategory(id)
}

func GetCategories(search string, page, limit int) (map[string]interface{}, error) {
	categories, total, err := repository.ListCategories(search, page, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":  categories,
		"total": total,
		"page":  page,
		"limit": limit,
	}, nil
}

//user product

func GetUserProducts(query models.UserProductQuery) (*models.UserProductListResponse, error) {
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	products, totalCount, err := repository.GetUserProducts(query)
	if err != nil {
		return nil, err
	}

	var resp []models.UserProductResponse

	for _, p := range products {
		var images []string
		for _, img := range p.Images {
			images = append(images, img.ImageURL)
		}

		resp = append(resp, models.UserProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			BrandID:     p.BrandID,
			SKU:         p.SKU,
			Price:       p.Price,
			Stock:       p.Stock,
			CategoryID:  p.CategoryID,
			IsListed:    p.IsListed,
			Images:      images,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(query.Limit)))

	return &models.UserProductListResponse{
		Products:   resp,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func GetUserProductDetails(productID uint) (*models.UserProductDetailResponse, error) {
	product, err := repository.GetUserProductDetails(productID)
	if err != nil {
		return nil, errors.New("product unavailable or blocked")
	}

	var images []string
	for _, img := range product.Images {
		images = append(images, img.ImageURL)
	}

	//set product status is available or out of stock
	status := "available"
	if product.Stock <= 0 {
		status = "out_of_stock"
	}

	related, err := repository.GetRelatedUserProducts(product.CategoryID, product.ID, 4)
	if err != nil {
		return nil, err
	}

	//related product images
	var relatedProducts []models.UserProductResponse
	for _, p := range related {
		var relImages []string
		for _, img := range p.Images {
			relImages = append(relImages, img.ImageURL)
		}

		relatedProducts = append(relatedProducts, models.UserProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			BrandID:     p.BrandID,
			SKU:         p.SKU,
			Price:       p.Price,
			Stock:       p.Stock,
			CategoryID:  p.CategoryID,
			IsListed:    p.IsListed,
			Images:      relImages,
		})
	}

	return &models.UserProductDetailResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		BrandID:      product.BrandID,
		SKU:          product.SKU,
		Price:        product.Price,
		Stock:        product.Stock,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.Name,
		IsListed:     product.IsListed,
		Images:       images,
		Status:       status,
		Breadcrumbs:  []string{"Home", "Products", product.Category.Name, product.Name},
		Highlights: []string{
			"Comfortable for daily wear",
			"Durable quality",
			"Stylish design",
		},
		RelatedProducts: relatedProducts,
	}, nil
}

func ValidateUserProductAvailability(productID uint) error {
	product, err := repository.GetUserProductDetails(productID)
	if err != nil {
		return errors.New("product unavailable or blocked")
	}

	if product.Stock <= 0 {
		return errors.New("product is out of stock")
	}

	return nil
}

// user category
func GetUserCategories(search string) ([]models.UserCategoryResponse, error) {
	categories, err := repository.GetUserCategories(search)
	if err != nil {
		return nil, err
	}

	var resp []models.UserCategoryResponse
	for _, c := range categories {
		resp = append(resp, models.UserCategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return resp, nil
}
