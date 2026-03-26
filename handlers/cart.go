package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

//-->user loggin and take  user id 
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		errRes := response.ClientResponse(http.StatusUnauthorized, "user not authorized", nil, "user_id not found in context")
		c.JSON(http.StatusUnauthorized, errRes)
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		errRes := response.ClientResponse(http.StatusUnauthorized, "invalid user", nil, "invalid user_id type")
		c.JSON(http.StatusUnauthorized, errRes)
		return 0, false
	}

	return userID, true
}
func AddToCart(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.AddToCart(userID, req.ProductID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to add to cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product added to cart successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetCart(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	cart, err := usecase.GetCart(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to get cart", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "cart fetched successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)
}

func UpdateCartQuantity(c *gin.Context) {
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
	productID := uint(productID64)

	var req models.UpdateCartQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.UpdateCartQuantity(userID, productID, req.Action)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to update cart quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "cart quantity updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func RemoveCartItem(c *gin.Context) {
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
	productID := uint(productID64)

	err = usecase.RemoveCartItem(userID, productID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to remove cart item", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "cart item removed successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func ValidateCartCheckout(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	err := usecase.ValidateCartForCheckout(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "cart is not valid for checkout", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "cart is valid for checkout", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
