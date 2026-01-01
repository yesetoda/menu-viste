package rest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"menuvista/internal/services/payment"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookService *payment.WebhookService
}

func NewWebhookHandler(webhookService *payment.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/chapa", h.ChapaWebhook)
	r.POST("/chapa/test", h.ChapaWebhookTest)
}

func (h *WebhookHandler) ChapaWebhook(c *gin.Context) {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📩 RECEIVED CHAPA WEBHOOK")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Failed to read webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
		return
	}

	signature := c.GetHeader("Chapa-Signature")
	if signature == "" {
		log.Printf("❌ Missing Chapa-Signature header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Chapa-Signature header"})
		return
	}

	if err := h.webhookService.ProcessWebhook(c.Request.Context(), body, signature); err != nil {
		log.Printf("⚠️  Webhook processing error: %v", err)
		// Return 200 to Chapa to stop retries if it's a permanent error or we've logged it
		c.JSON(http.StatusOK, gin.H{
			"status":  "received",
			"message": fmt.Sprintf("Webhook received but processing error: %v", err),
		})
		return
	}

	log.Printf("✅ Webhook processed successfully")
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Webhook processed successfully",
		"timestamp": time.Now(),
	})
}

func (h *WebhookHandler) ChapaWebhookTest(c *gin.Context) {
	// This is for local testing as described in the guide
	// It simulates a webhook call
	webhookType := c.DefaultQuery("type", "success")
	txRef := c.DefaultQuery("tx_ref", fmt.Sprintf("tx_test_%d", time.Now().Unix()))

	payload := map[string]interface{}{
		"event": "payment." + webhookType,
		"data": map[string]string{
			"status":     "completed",
			"first_name": "Test",
			"last_name":  "User",
			"email":      "test@example.com",
			"amount":     "19.99",
			"currency":   "ETB",
			"tx_ref":     txRef,
			"reference":  fmt.Sprintf("ref_%d", time.Now().Unix()),
		},
	}
	if webhookType == "failed" {
		payload["data"].(map[string]string)["status"] = "failed"
	}

	body, _ := json.Marshal(payload)

	// In a real test, we'd calculate the signature, but for the test endpoint
	// we might want to skip it or provide a way to test it.
	// For now, let's just call ProcessWebhook with a dummy signature if we want to bypass verification
	// or implement a test-only bypass in WebhookService.
	// Actually, let's just use the real secret if it's set.

	secret := os.Getenv("CHAPA_WEBHOOK_SECRET")
	hmacObj := hmac.New(sha256.New, []byte(secret))
	hmacObj.Write(body)
	signature := hex.EncodeToString(hmacObj.Sum(nil))

	c.Request.Header.Set("Chapa-Signature", signature)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	h.ChapaWebhook(c)
}
