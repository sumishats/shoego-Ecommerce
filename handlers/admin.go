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

// @Summary Admin Login
// @Description Login handler for admin
// @Tags Admin Authentication
// @Accept json
// @Produce json
// @Param  admin body models.AdminLogin true "Admin login details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/admin [post]

func AdminLogin(c *gin.Context) {
	var adminmodel models.AdminLogin

	if err := c.ShouldBindJSON(&adminmodel); err != nil {
		erres := response.ClientResponse(http.StatusBadGateway, "formate is not correct admin", nil, err)
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	if err := validator.New().Struct(adminmodel); err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "constrian are not satisfied ", nil, err.Error())
		c.JSON(http.StatusBadGateway, errRes)
		return

	}

	admin, err := usecase.AdminLogin(adminmodel)
	if errors.Is(err, models.PasswordIsNil) || errors.Is(err, models.PasswordIsNotCorrect) {

		erres := response.ClientResponse(http.StatusBadGateway, "", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return

	}

	if err != nil {

		erres := response.ClientResponse(http.StatusBadGateway, "server error from admin use case", nil, err.Error())
		c.JSON(http.StatusBadGateway, erres)
		return
	}

	succesRes := response.ClientResponse(http.StatusOK, "admin login succes ", admin, nil)
	c.JSON(http.StatusOK, succesRes)

}
