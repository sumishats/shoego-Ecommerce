package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

// admin coupon
func CreateCoupon(c *gin.Context) {
	var req models.CreateCouponRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid request", nil, err.Error()))
		return
	}

	err := usecase.CreateCoupon(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to create coupon", nil, err.Error()))
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "coupon created successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}
func GetAllCoupons(c *gin.Context) {
	coupons, err := usecase.GetAllCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to get coupons", nil, err.Error()))
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "coupons fetched successfully", coupons, nil)
	c.JSON(http.StatusOK, successRes)
}
func DeleteCoupon(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid coupon id", nil, err.Error()))
		return
	}

	err = usecase.DeleteCoupon(uint(id))
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to delete coupon", nil, err.Error()))
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "coupon deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// user coupon
func ApplyCoupon(c *gin.Context) {

	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.ClientResponse(400, "invalid request", nil, err.Error()))
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	result, err := usecase.ApplyCoupon(userID, req.Code)
	if err != nil {
		c.JSON(400, response.ClientResponse(400, err.Error(), nil, err.Error()))
		return
	}

	c.JSON(200, response.ClientResponse(200, "coupon applied", result, nil))
}
func RemoveCoupon(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	err := usecase.RemoveCoupon(userID)
	if err != nil {
		c.JSON(400, response.ClientResponse(400, "failed to remove coupon", nil, err.Error()))
		return
	}

	c.JSON(200, response.ClientResponse(200, "coupon removed", nil, nil))
}
