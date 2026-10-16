package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type checkoutHandoff struct {
	AccessToken string
	ExpiresAt   time.Time
}

var (
	handoffMu     sync.Mutex
	handoffTokens = make(map[string]checkoutHandoff)
)

func CreateCheckoutHandoff(c *gin.Context) {
	value, exists := c.Get("auth_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ClientResponse(http.StatusUnauthorized,"authentication token missing",nil,nil,))
		return
	}

	accessToken, ok := value.(string)
	if !ok || accessToken == "" {
		c.JSON(http.StatusUnauthorized, response.ClientResponse(http.StatusUnauthorized,"invalid authentication token",nil,nil))
		return
	}

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		c.JSON(http.StatusInternalServerError, response.ClientResponse(http.StatusInternalServerError,"could not create checkout link",nil,
			err.Error(),
		))
		return
	}

	token := hex.EncodeToString(randomBytes)

	handoffMu.Lock()

	for key, item := range handoffTokens {
		if time.Now().After(item.ExpiresAt) {
			delete(handoffTokens, key)
		}
	}

	handoffTokens[token] = checkoutHandoff{
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(60 * time.Second),
	}

	handoffMu.Unlock()

	link := "http://localhost:8080/checkout/handoff?token=" + token

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK,"checkout handoff link created",
		gin.H{
			"url":                link,
			"expires_in_seconds": 60,
		},
		nil,
	))
}

func OpenCheckoutHandoff(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest,"handoff token is required",nil,nil,))
		return
	}

	handoffMu.Lock()
	item, exists := handoffTokens[token]
	delete(handoffTokens, token)
	handoffMu.Unlock()

	if !exists || time.Now().After(item.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, response.ClientResponse(http.StatusUnauthorized,"checkout link is invalid or expired; create a new one",nil,nil,))
		return
	}

	c.Header("Cache-Control", "no-store")
	c.Header("Referrer-Policy", "no-referrer")

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		item.AccessToken,
		15*60,
		"/",
		"",
		false, 
		true,  
	)

	c.Redirect(http.StatusSeeOther, "/static/index.html")
}

func GetCheckoutPage(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	checkout, err := usecase.GetCheckoutPage(userID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to load checkout page", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "checkout page loaded successfully", checkout, nil)
	c.JSON(http.StatusOK, successRes)
}

func PlaceCODOrder(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	order, err := usecase.PlaceCODOrder(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to place order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order placed successfully", order, nil)
	c.JSON(http.StatusOK, successRes)
}

func PlaceWalletOrder(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	var req models.PlaceOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	order, err := usecase.PlaceWalletOrder(userID, req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to place wallet order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	sucessRes := response.ClientResponse(http.StatusOK, "wallet order placed successfully", order, nil)
	c.JSON(http.StatusOK, sucessRes)
}
func GetOrderSuccess(c *gin.Context) {
	// userID, ok := getUserIDFromContext(c)
	// if !ok {
	// 	return
	// }

	orderID := c.Param("order_id")

	success, err := usecase.GetOrderSuccess(orderID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusNotFound, "order not found", nil, err.Error())
		c.JSON(http.StatusNotFound, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "order success page loaded", success, nil)
	c.JSON(http.StatusOK, successRes)
}
