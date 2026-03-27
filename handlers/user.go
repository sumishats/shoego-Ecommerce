package handlers

import (
	"errors"
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// @Summary SignUp functionality for user
// @Description SignUp functionality at the user side
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param user body models.SignupDetail true "User Details Input"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /signup [post]

func Signup(c *gin.Context) {

	var usersign models.SignupDetail

	if err := c.ShouldBindJSON(&usersign); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// check the data sending in correct or not format
	if err := validator.New().Struct(usersign); err != nil {

		errres := response.ClientResponse(404, "They are not in format", nil, err.Error())
		c.JSON(http.StatusBadGateway, errres)
		return
	}

	err := usecase.UsersignUp(usersign)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "user signup format error ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "OTP sent to email successfully", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

//userlogin

// @Summary LogIn functionality for user
// @Description LogIn functionality at the user side
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User Details Input"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /login [post]

func UserLoginWithPassword(c *gin.Context) {
	var LoginUser models.UserLogin

	if err := c.ShouldBindJSON(&LoginUser); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "Login field provided in wrong way ", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	if err := validator.New().Struct(LoginUser); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "Login field was wrong formate ", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	LogedUser, err := usecase.UserLogged(LoginUser)

	if errors.Is(err, models.ErrEmailNotFound) {

		erres := response.ClientResponse(http.StatusBadRequest, "invalid email", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}
	if err != nil {

		erres := response.ClientResponse(500, "server error from usecase", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	successres := response.ClientResponse(http.StatusCreated, "succesed login user", LogedUser, nil)

	c.JSON(http.StatusOK, successres)
}

//user profile
func GetProfile(c *gin.Context) {
	
	userIDVal, exists := c.Get("user_id")
	if !exists {
		errRes := response.ClientResponse(http.StatusUnauthorized, "user not found in token", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "invalid user id type", nil, nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	data, err := usecase.GetUserProfile(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch profile", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "profile fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

func EditProfile(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req models.EditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.EditUserProfile(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to update profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "profile updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}


func ChangePassword(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.ChangePassword(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to change password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "password changed successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func RequestEmailChange(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req models.EmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.RequestEmailChange(userID, req.NewEmail)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to send otp", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "otp sent to new email", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func VerifyEmailChange(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req models.VerifyEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.VerifyEmailChange(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to verify email change", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "email updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func AddAddress(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	var req models.AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.AddUserAddress(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to add address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "address added successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func EditAddress(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid address id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.EditAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid input", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.EditUserAddress(userID, uint(id64), req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to edit address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "address updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func DeleteAddress(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid address id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.DeleteUserAddress(userID, uint(id64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to delete address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "address deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errRes := response.ClientResponse(http.StatusBadRequest, "authorization header is missing", nil, "no token found")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.Logout(authHeader)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to logout", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes:=response.ClientResponse(http.StatusOK,"logout successful",nil,nil)
	c.JSON(http.StatusOK,successRes)
}
