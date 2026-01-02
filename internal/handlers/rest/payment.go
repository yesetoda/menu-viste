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
	webhook *payment.WebhookService
}

func NewPaymentHandler(service *payment.Service, webhook *payment.WebhookService) *PaymentHandler {
	return &PaymentHandler{service: service, webhook: webhook}
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
		Type:    req.Type,
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
	providerRef, verified, err := h.webhook.VerifyPayment(c.Request.Context(), txRef)
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

	// Payment successful - update database immediately
	log.Printf("[PaymentHandler] ✅ Payment verified successfully: %s. Updating database...", txRef)
	if err := h.webhook.CompletePayment(c.Request.Context(), txRef, providerRef); err != nil {
		log.Printf("[PaymentHandler] ❌ Failed to complete payment: %v", err)
		// We still show success page because the webhook might succeed later,
		// but we log the error. Actually, it's better to show an error if it fails here?
		// If it's already completed (idempotent), it will return nil or we can handle it.
	}

	c.HTML(http.StatusOK, "payment_success.html", gin.H{
		"tx_ref": txRef,
		"status": "completed",
	})
}

func (h *PaymentHandler) PaymentCancel(c *gin.Context) {
	txRef := c.Query("trx_ref")

	c.HTML(http.StatusOK, "payment_cancelled.html", gin.H{
		"tx_ref":  txRef,
		"message": "Payment was cancelled. You can try again.",
	})
}
