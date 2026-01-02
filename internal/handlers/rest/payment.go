package rest

import (
	"fmt"
	"log"
	"net/http"

	"menuvista/internal/services/payment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	service *payment.Service
}

func NewPaymentHandler(service *payment.Service) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	log.Printf("[PaymentHandler] InitiatePayment request received")

	userIDVal, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusInternalServerError, "user_id not found in context", "INTERNAL_ERROR")
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		RespondError(c, http.StatusInternalServerError, "user_id in context is not of type uuid.UUID", "INTERNAL_ERROR")
		return
	}

	var req struct {
		Plan string `json:"plan"`
		Type string `json:"type"` // update, upgrade
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	// if req.Type == "" {
	// 	req.Type = "registration"
	// }

	input := payment.InitiatePaymentInput{
		OwnerID: userID,
		Plan:    req.Plan,
	}
	fmt.Println("this is the payment input", input)

	resp, err := h.service.InitiatePayment(c.Request.Context(), input)
	if err != nil {
		log.Printf("[PaymentHandler] InitiatePayment service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, resp, nil)
}

func (h *PaymentHandler) RenewSubscription(c *gin.Context) {
	// Wrapper for InitiatePayment with type=renewal
	c.Set("type", "renewal")
	h.InitiatePayment(c)
}

func (h *PaymentHandler) UpgradeSubscription(c *gin.Context) {
	// Wrapper for InitiatePayment with type=upgrade
	c.Set("type", "upgrade")
	h.InitiatePayment(c)
}

func (h *PaymentHandler) ChapaWebhook(c *gin.Context) {
	log.Printf("[PaymentHandler] ChapaWebhook request received")

	var payload struct {
		TxRef  string `json:"tx_ref"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("[PaymentHandler] Webhook bind error: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.service.HandleWebhook(c.Request.Context(), payload.TxRef, payload.Status)
	if err != nil {
		log.Printf("[PaymentHandler] Webhook service error: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *PaymentHandler) PaymentSuccess(c *gin.Context) {
	log.Printf("[PaymentHandler] 📥 Payment success callback received")
	log.Printf("[PaymentHandler]    Full URL: %s", c.Request.URL.String())
	log.Printf("[PaymentHandler]    Query Params: %v", c.Request.URL.Query())

	// Chapa sends different parameter names, try all possibilities
	txRef := c.Query("trx_ref") // Try trx_ref first
	if txRef == "" {
		txRef = c.Query("tx_ref") // Try tx_ref
	}
	if txRef == "" {
		txRef = c.Query("reference") // Try reference
	}

	status := c.Query("status")

	log.Printf("[PaymentHandler]    Extracted - txRef: %s, status: %s", txRef, status)

	if txRef == "" {
		log.Printf("[PaymentHandler] ❌ No transaction reference found in query params")
		c.HTML(http.StatusBadRequest, "payment_error.html", gin.H{
			"error": "Invalid payment reference. Please contact support if payment was deducted.",
		})
		return
	}

	// Verify payment
	log.Printf("[PaymentHandler] 🔍 Verifying payment with Chapa: %s", txRef)
	verified, err := h.service.VerifyPayment(c.Request.Context(), txRef)
	if err != nil {
		log.Printf("[PaymentHandler] ❌ Payment verification failed: %v", err)
		c.HTML(http.StatusInternalServerError, "payment_error.html", gin.H{
			"error": "Failed to verify payment: " + err.Error(),
		})
		return
	}

	if !verified {
		log.Printf("[PaymentHandler] ⚠️ Payment verification returned false")
		c.HTML(http.StatusOK, "payment_failed.html", gin.H{
			"status": status,
			"tx_ref": txRef,
		})
		return
	}

	// Payment successful - show success page
	log.Printf("[PaymentHandler] ✅ Payment verified successfully: %s", txRef)
	c.HTML(http.StatusOK, "payment_success.html", gin.H{
		"tx_ref": txRef,
		"status": "completed",
		// "amount": amount, // Need to fetch from DB or Chapa verify response
		// "currency": currency,
	})
}

func (h *PaymentHandler) PaymentCancel(c *gin.Context) {
	txRef := c.Query("trx_ref")

	c.HTML(http.StatusOK, "payment_cancelled.html", gin.H{
		"tx_ref":  txRef,
		"message": "Payment was cancelled. You can try again.",
	})
}
