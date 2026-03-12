package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func ResendOTP(c *gin.Context) {

	var req models.ResendOTP

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	err := usecase.ResendOTP(req.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "OTP sent again")
}

func VerifyOTP(c *gin.Context) {

	var otpReq models.VerifyOTP

	if err := c.ShouldBindJSON(&otpReq); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	user, err := usecase.VerifyOTPAndCreateUser(otpReq)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// ✅ Build response struct
	userRes := response.UserResponse{
		ID:    user.User.ID,
		Name:  user.User.Name,
		Email: user.User.Email,
		Phone: user.User.Phone,
	}
	// ✅ Send response
	successRes := response.ClientResponse(
		http.StatusOK,
		"OTP verified successfully",
		userRes,
		nil,
	)

	c.JSON(http.StatusOK, successRes)
}

func ForgotPassword(c *gin.Context) {
	var req models.ForgotPassword

	if err := c.ShouldBindJSON(&req); err != nil {
		errReq := response.ClientResponse(http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errReq)
		return
	}

	err := usecase.ForgotPassword(req.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "OTP sent to your email ")
}

func ResetPassword(c *gin.Context) {
	var req models.ResetPassword

	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	userToken, err := usecase.ResetPassword(req.Email, req.OTP, req.NewPassword)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, err.Error(), nil, nil)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Password reset successfully", userToken, nil)
	c.JSON(http.StatusOK, successRes)
}
