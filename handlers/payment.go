package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"shoego/config"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func CreateRazorpayOrder(c *gin.Context) {

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.CreateRazorpayOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	result, err := usecase.CreateRazorpayOrder(userID, req)

	if err != nil {
		fmt.Println("RAZORPAY ERROR:", err)
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to create razorpay order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "razorpay order created successfully", result, nil)
	c.JSON(http.StatusOK, successRes)
}
func VerifyPayment(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.VerifyPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid request", nil, err.Error()))
		return
	}

	err := usecase.VerifyPayment(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "payment verification failed", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "payment verified successfully", nil, nil))
}

func HandlePaymentFailure(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.PaymentFailureRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid payment failure request", nil, err.Error()))
		return
	}

	err := usecase.HandlePaymentFailure(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "could not confirm payment failure", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "payment failure recorded", nil, nil))
}

func RazorpayWebhook(c *gin.Context) {
	webhookSecret := config.AppConfig.RAZORPAY_WEBHOOK_SECRET

	if webhookSecret == "" {
		log.Println("Razorpay webhook secret is missing")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "webhook configuration error",
		})
		return
	}

	// Read the original body; Razorpay's signature depends on these exact bytes.
	rawBody, err := io.ReadAll(
		io.LimitReader(c.Request.Body, (1<<20)+1),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not read webhook body",
		})
		return
	}

	if len(rawBody) > 1<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"message": "webhook body is too large",
		})
		return
	}

	receivedSignature := c.GetHeader("X-Razorpay-Signature")
	if receivedSignature == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing webhook signature",
		})
		return
	}

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(rawBody)
	expectedSignature := mac.Sum(nil)

	receivedBytes, err := hex.DecodeString(receivedSignature)
	if err != nil || !hmac.Equal(receivedBytes, expectedSignature) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid webhook signature",
		})
		return
	}

	// Signature is valid. Process the event in the usecase layer.
	if err := usecase.HandleRazorpayWebhook(rawBody); err != nil {
		log.Printf("Razorpay webhook processing failed: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "webhook processing failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "webhook processed successfully",
	})
}
