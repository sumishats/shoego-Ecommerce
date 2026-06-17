package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func CreateRazorpayOrder(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.CreateRazorpayOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	result, err := usecase.CreateRazorpayOrder(userID, req)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to create razorpay order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "razorpay order created successfully", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func VerifyPayment(c *gin.Context) {

	var req models.VerifyPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid request", nil, err.Error()))
		return
	}

	err := usecase.VerifyPayment(req)
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "payment verification failed", nil, err.Error()))
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "payment verified successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
