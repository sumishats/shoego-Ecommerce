package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAdminOrders(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := usecase.GetAdminOrders(page, limit, search, status, sortBy)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to get orders", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "order fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)

}

func GetAdminOrderDetail(c *gin.Context) {
	orderID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid order id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	data, err := usecase.GetAdminOrderDetail(uint(orderID64))
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to get order detail", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order detail fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

func UpdateAdminOrderStatus(c *gin.Context) {
	orderID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid order id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.ChangeOrderStatus(uint(orderID64), req.OrderStatus)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to update order status", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order status updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

//inventory

func GetAdminInventory(c *gin.Context) {
	search := c.Query("search")
	stockFilter := c.Query("stock_filter")
	sortBy := c.Query("sort")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := usecase.GetAdminInventory(page, limit, search, stockFilter, sortBy)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch inventory", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "inventory fetched successfully", data, nil)
	c.JSON(http.StatusOK, successRes)
}

func UpdateAdminProductStock(c *gin.Context) {
	productID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid product id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.ChangeProductStock(uint(productID64), req.Stock)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to update stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "stock updated successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

//user order management

func GetUserOrders(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || page <= 0 {
		page = 1
	}
	search := c.Query("search")

	orders, err := usecase.GetUserOrderList(userID, search, page, limit)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch orders", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "orders fetched successfully", orders, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetOrderDetail(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")

	order, err := usecase.GetUserOrderDetail(userID, orderID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to fetch order detail", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order detail fetched successfully", order, nil)
	c.JSON(http.StatusOK, successRes)
}

func CancelOrder(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")

	var req models.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.CancelOrder(userID, orderID, req.Reason)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to cancel order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order cancelled successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
func CancelOrderItem(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")
	itemID64, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid item id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.CancelOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.CancelOrderItem(userID, orderID, uint(itemID64), req.Reason)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to cancel order item", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order item cancelled successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
func ReturnOrder(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")

	var req models.ReturnOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.ReturnOrder(userID, orderID, req.Reason)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to return order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order returned successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
func ReturnOrderItem(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")
	itemID64, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid item id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	var req models.ReturnOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err = usecase.ReturnOrderItem(userID, orderID, uint(itemID64), req.Reason)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to return order item", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order item returned successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// invoice
func DownloadInvoice(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orderID := c.Param("order_id")

	pdfBytes, err := usecase.GenerateInvoice(userID, orderID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to generate invoice", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=invoice-"+orderID+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
