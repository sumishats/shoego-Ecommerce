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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := usecase.GetAllCoupons(page, limit)
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to get coupons", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "coupons fetched successfully", data, nil))
}


func UpdateCoupon(c *gin.Context) {
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid coupon id", nil, err.Error()))
		return
	}
	
	var req models.UpdateCouponRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid request", nil, err.Error()))
		return
	}
	
	err = usecase.UpdateCoupon(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to update coupon", nil, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "coupon updated successfully", nil, nil))
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
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, err.Error(), nil, err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "coupon applied", result, nil))
}
func GetAvailableCoupons(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := usecase.GetAvailableCoupons(userID, page, limit)
	if err != nil {

		c.JSON(http.StatusBadRequest,response.ClientResponse(http.StatusBadRequest,"failed to fetch coupons",nil,err.Error(),),)
		return
	}

	c.JSON(http.StatusOK,response.ClientResponse(http.StatusOK,"available coupons fetched successfully",data,nil,),)
}
func RemoveCoupon(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	err := usecase.RemoveCoupon(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to remove coupon", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "coupon removed", nil, nil))
}
