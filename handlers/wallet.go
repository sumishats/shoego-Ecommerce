package handlers

import (
	"net/http"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func GetWallet(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	wallet, err := usecase.GetWallet(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to get wallet", nil, err.Error()))
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "wallet fetched successfully", wallet, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetWalletHistory(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	history, err := usecase.GetWalletHistory(userID)
	if err != nil {

		c.JSON(http.StatusBadRequest,response.ClientResponse(http.StatusBadRequest,"failed to get wallet history",nil,err.Error()))
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "wallet history fetched successfully", history, nil)
	c.JSON(http.StatusOK, successRes)
}
