package handlers

import (
	"net/http"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	search := c.Query("search")

	users, err := usecase.GetAdminUsers(page, limit, search) 
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "users fetched successfully", users, nil)
	c.JSON(http.StatusOK, successRes)
}

func BlockUser(c *gin.Context) {
	idParam := c.Param("id") //get the user id from the url parameter

	//convert user id string to number 
	id64, err := strconv.ParseUint(idParam, 10, 64) 
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid user id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.BlockUser(uint(id64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to block user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "user blocked successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func UnblockUser(c *gin.Context) {
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid user id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.UnblockUser(uint(id64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to unblock user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "user unblocked successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}