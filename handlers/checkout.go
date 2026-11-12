package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"shoego/models"
	"shoego/response"
	"shoego/usecase"

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
		c.JSON(http.StatusUnauthorized,
			response.ClientResponse(
				http.StatusUnauthorized,
				"authentication token missing",
				nil, nil,
			))
		return
	}

	accessToken, ok := value.(string)
	if !ok || accessToken == "" {
		c.JSON(http.StatusUnauthorized,
			response.ClientResponse(
				http.StatusUnauthorized,
				"invalid authentication token",
				nil, nil,
			))
		return
	}

	// Get application base URL from environment
	baseURL := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")

	if baseURL == "" {
		c.JSON(http.StatusInternalServerError,
			response.ClientResponse(
				http.StatusInternalServerError,
				"application base URL is not configured",
				nil, nil,
			))
		return
	}

	// Generate secure random token
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		log.Println("Checkout token generation failed:", err)

		c.JSON(http.StatusInternalServerError,
			response.ClientResponse(
				http.StatusInternalServerError,
				"could not create checkout link",
				nil, nil,
			))
		return
	}

	token := hex.EncodeToString(randomBytes)

	now := time.Now()

	handoffMu.Lock()

	for key, item := range handoffTokens {
		if now.After(item.ExpiresAt) {
			delete(handoffTokens, key)
		}
	}

	handoffTokens[token] = checkoutHandoff{
		AccessToken: accessToken,
		ExpiresAt:   now.Add(60 * time.Second),
	}

	handoffMu.Unlock()

	link := baseURL + "/checkout/handoff?token=" + token

	c.JSON(http.StatusOK,
		response.ClientResponse(
			http.StatusOK,
			"checkout handoff link created",
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
		c.JSON(http.StatusBadRequest,
			response.ClientResponse(
				http.StatusBadRequest,
				"handoff token is required",
				nil, nil,
			))
		return
	}

	
	handoffMu.Lock()

	item, exists := handoffTokens[token]

	if exists {
		delete(handoffTokens, token)
	}

	handoffMu.Unlock()

	if !exists || time.Now().After(item.ExpiresAt) {
		c.JSON(http.StatusUnauthorized,
			response.ClientResponse(
				http.StatusUnauthorized,
				"checkout link is invalid or expired; create a new one",
				nil, nil,
			))
		return
	}


	c.Header("Cache-Control", "no-store")
	c.Header("Referrer-Policy", "no-referrer")

	
	secureCookie := os.Getenv("APP_ENV") == "production"

	c.SetSameSite(http.SameSiteLaxMode)
	
	c.SetCookie(
		"jwt_token",
		item.AccessToken,
		15*60,
		"/",
		"",
		secureCookie,
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
