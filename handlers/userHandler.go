package handlers

import (
	"errors"
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

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
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format 🙌", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// CHEKING THE DATA ARE SENDED IN CORRECT FORMET OR NOT

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

func UserLoginWithPassword(c *gin.Context) {
	var LoginUser models.UserLogin

	if err := c.ShouldBindJSON(&LoginUser); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "Login field provided in wrong way ", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	if err := validator.New().Struct(LoginUser); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "Login field was wrong formate ahn", nil, err.Error())
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
