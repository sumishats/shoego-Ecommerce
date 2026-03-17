package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)
// @Summary  OTP login
// @Description Send OTP to Authenticate user
// @Tags User OTP Login
// @Accept json
// @Produce json
// @Param phone body models.OTPData true "phone number details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /send-otp [post]

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

// @Summary Verify OTP
// @Description Verify OTP by passing the OTP in order to authenticate user
// @Tags User OTP Login
// @Accept json
// @Produce json
// @Param phone body models.VerifyData true "Verify OTP Details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /verify-otp [post]

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

	//return user response struct 
	userRes := response.UserResponse{
		ID:    user.User.ID,
		Name:  user.User.Name,
		Email: user.User.Email,
		Phone: user.User.Phone,
	}
	// Send success response
	successRes := response.ClientResponse(
		http.StatusOK,
		"OTP verified successfully",
		userRes,
		nil,
	)

	c.JSON(http.StatusOK, successRes)
}

// @Summary Forgot Password
// @Description Send OTP to user email to reset password
// @Tags User Password
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Forgot Password Details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /forgot-password [post]

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

// @Summary Reset Password
// @Description Reset user password using OTP verification
// @Tags User Password
// @Accept json
// @Produce json
// @Param data body models.ResetPasswordRequest true "Reset Password Details"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /reset-password [post]

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
