package handlers

import (
	"errors"
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddToWishlist(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.AddToWishlist(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, err.Error(), nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "product added to wishlist", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}

func RemoveFromWishlist(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	productIDStr := c.Param("product_id")
	productID64, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.RemoveFromWishlist(userID, uint(productID64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errRes := response.ClientResponse(http.StatusNotFound, "wishlist item not found", nil, err.Error())
			c.JSON(http.StatusNotFound, errRes)
			return
		}

		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to remove wishlist item", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product removed from wishlist", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetWishlist(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	wishlist, err := usecase.GetWishlist(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch wishlist", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "wishlist fetched successfully", wishlist, nil)
	c.JSON(http.StatusOK, successRes)
}