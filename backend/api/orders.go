package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateOrder(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}
	fmt.Println(req.Amount)
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	// amount := 50000 // Amount in paise (i.e., ₹500.00)
	receipt := "receipt#123"

	data := map[string]interface{}{
		"amount":          req.Amount,
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": 1,
	}

	body, err := server.razorPayClient.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, body)
}

func (server *Server) VerifySignature(c *gin.Context) {
	var req struct {
		OrderID   string `json:"order_id"`
		PaymentID string `json:"payment_id"`
		Signature string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	secret := server.config.RazorpaySecret
	payload := req.OrderID + "|" + req.PaymentID

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != req.Signature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Payment verified"})
}
