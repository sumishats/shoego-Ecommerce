package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func GetCheckoutPage(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	checkout, err := usecase.GetCheckoutPage(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to load checkout page", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "checkout page loaded successfully", checkout, nil)
	c.JSON(http.StatusOK, successRes)
}
func PlaceCODOrder(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	order, err := usecase.PlaceCODOrder(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to place order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order placed successfully", order, nil)
	c.JSON(http.StatusOK, successRes)
}
func PlaceWalletOrder(c *gin.Context){
	userID,ok:=getUserIDFromContext(c)
	if !ok{
		return
	}
	var req models.PlaceOrderRequest

	if err:=c.ShouldBindJSON(&req);err!=nil{
		errRes:=response.ClientResponse(http.StatusBadRequest,"invalid request body",nil,err.Error())
		c.JSON(http.StatusBadRequest,errRes)
		return
	}
	order,err:=usecase.PlaceWalletOrder(userID,req)
	if err!=nil{
		errRes:=response.ClientResponse(http.StatusBadRequest,"failed to place wallet order",nil,err.Error())
		c.JSON(http.StatusBadRequest,errRes)
		return
	}
	sucessRes:=response.ClientResponse(http.StatusOK,"wallet order placed successfully",order,nil)
	c.JSON(http.StatusOK,sucessRes)
}
func GetOrderSuccess(c *gin.Context) {
	// userID, ok := getUserIDFromContext(c)
	// if !ok {
	// 	return
	// }

	orderID := c.Param("order_id")

	success, err := usecase.GetOrderSuccess(orderID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "order not found", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order success page loaded", success, nil)
	c.JSON(http.StatusOK, successRes)
}

