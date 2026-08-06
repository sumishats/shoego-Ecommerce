package handlers

import (
	"net/http"

	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func GetAdminBestSellingProducts(c *gin.Context) {

	products, err := usecase.GetAdminBestSellingProducts()
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch best selling products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "best selling products fetched successfully", products, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetAdminBestSellingCategories(c *gin.Context) {

	categories, err := usecase.GetAdminBestSellingCategories()
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch best selling categories", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "best selling categories fetched successfully", categories, nil)
	c.JSON(http.StatusOK, successRes)
}
func GetAdminSalesChart(c *gin.Context) {

	filter := c.Query("filter")

	sales, err := usecase.GetAdminSalesChart(filter)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch sales chart", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "sales chart fetched successfully", sales, nil)
	c.JSON(http.StatusOK, successRes)
}
