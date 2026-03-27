package handlers

import (
	"mime/multipart"
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

//admin product management

func AddProduct(c *gin.Context) {
	var req models.AddProductRequest

	//request is multipart form data 
	req.Name = c.PostForm("name")
	req.Description = c.PostForm("description")
	req.SKU = c.PostForm("sku")

	
	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid price", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.Price = price

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.Stock = stock

	brandID64, err := strconv.ParseUint(c.PostForm("brand_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid brand id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.BrandID = uint(brandID64)

	categoryID64, err := strconv.ParseUint(c.PostForm("category_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid category id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.CategoryID = uint(categoryID64)

	//get all form including images 
	form, err := c.MultipartForm()
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid multipart form", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	files := form.File["images"]

	err = usecase.AddProduct(req, files)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to add product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product added successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func EditProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID64, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.EditProductRequest

	req.Name = c.PostForm("name")
	req.Description = c.PostForm("description")
	req.SKU = c.PostForm("sku")

	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid price", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.Price = price

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.Stock = stock

	brandID64, err := strconv.ParseUint(c.PostForm("brand_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid brand id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.BrandID = uint(brandID64)

	categoryID64, err := strconv.ParseUint(c.PostForm("category_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid category id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	req.CategoryID = uint(categoryID64)

	var files []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files = form.File["images"]
	}

	err = usecase.EditProduct(uint(productID64), req, files)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to edit product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.DeleteProduct(uint(productID64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to delete product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetProducts(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	data, err := usecase.GetProducts(page, limit)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "products fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}


//admin category management

func AddCategory(c *gin.Context) {
	var data models.AddCategory

	if err := c.ShouldBindJSON(&data); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.AddCategory(data)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to add category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "category added successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func EditCategory(c *gin.Context) {
	idParam := c.Param("id")
	categoryID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid category id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var data models.UpdateCategory

	if err := c.ShouldBindJSON(&data); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.EditCategory(uint(categoryID64), data)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to update category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "category updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	categoryID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid category id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.DeleteCategory(uint(categoryID64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to delete category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "category deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetCategories(c *gin.Context) {
	search := c.Query("search")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		limit = 5
	}

	data, err := usecase.GetCategories(search, page, limit)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch categories", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "categories fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

//user product

func GetUserProducts(c *gin.Context) {
	search := c.Query("search")
	sort := c.Query("sort")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1")) 
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categoryID64, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)
	brandID64, _ := strconv.ParseUint(c.DefaultQuery("brand_id", "0"), 10, 64)

	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)

	//url query params to struct
	query := models.UserProductQuery{
		Search:     search,
		Sort:       sort,
		CategoryID: uint(categoryID64),
		BrandID:    uint(brandID64),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       page,
		Limit:      limit,
	}

	data, err := usecase.GetUserProducts(query)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "products fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetUserProductDetails(c *gin.Context) {
	idParam := c.Param("id")
	productID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	data, err := usecase.GetUserProductDetails(uint(productID64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to fetch product details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product details fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

//check product is currently availbale 
func ValidateUserProductAvailability(c *gin.Context) {
	idParam := c.Param("id")
	productID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.ValidateUserProductAvailability(uint(productID64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "product unavailable", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product is available", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

//user category

func GetUserCategories(c *gin.Context) {
	search := c.Query("search") 
	data, err := usecase.GetUserCategories(search)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch categories", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "categories fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}