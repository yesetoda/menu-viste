package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/services/email"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type WebhookService struct {
	queries      *persistence.Queries
	emailService *email.Service
}

func NewWebhookService(queries *persistence.Queries, emailService *email.Service) *WebhookService {
	return &WebhookService{
		queries:      queries,
		emailService: emailService,
	}
}

func (s *WebhookService) ProcessWebhook(ctx context.Context, body []byte, signature string) error {
	// 1. Verify signature
	if !s.verifySignature(body, signature) {
		return fmt.Errorf("invalid webhook signature")
	}

	// 2. Parse payload
	var payload models.ChapaWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// 3. Idempotency check: Save webhook
	webhook, err := s.queries.CreatePaymentWebhook(ctx, persistence.CreatePaymentWebhookParams{
		ProviderEventID: pgtype.Text{String: payload.Data.Reference, Valid: true},
		EventType:       payload.Event,
		Payload:         body,
	})
	if err != nil {
		log.Printf("[WebhookService] Failed to save webhook (might be duplicate): %v", err)
		return nil // Return nil to acknowledge receipt even if duplicate
	}
	if webhook.Processed.Bool {
		log.Printf("[WebhookService] Webhook already processed: %s", payload.Data.Reference)
		return nil
	}

	// 4. Process based on event type
	switch payload.Event {
	case "payment.success":
		err = s.handlePaymentSuccess(ctx, payload.Data)
	case "payment.failed":
		err = s.handlePaymentFailed(ctx, payload.Data)
	case "payment.pending":
		err = s.handlePaymentPending(ctx, payload.Data)
	default:
		log.Printf("[WebhookService] Unhandled event type: %s", payload.Event)
	}

	if err == nil {
		// Mark as processed
		_ = s.queries.MarkWebhookAsProcessed(ctx, pgtype.Text{String: payload.Data.Reference, Valid: true})
	}

	return err
}

func (s *WebhookService) verifySignature(body []byte, signature string) bool {
	secret := os.Getenv("CHAPA_WEBHOOK_SECRET")
	if secret == "" {
		log.Printf("[WebhookService] CHAPA_WEBHOOK_SECRET not set")
		return false
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func (s *WebhookService) handlePaymentSuccess(ctx context.Context, data models.ChapaWebhookData) error {
	log.Printf("[WebhookService] Processing payment success for tx: %s", data.TxRef)
	return s.CompletePayment(ctx, data.TxRef, data.Reference)
}

func (s *WebhookService) CompletePayment(ctx context.Context, txRef string, providerRef string) error {
	log.Printf("[WebhookService] Completing payment for tx: %s, providerRef: %s", txRef, providerRef)

	// 0. Idempotency check
	tx, err := s.queries.GetPaymentTransactionByTxRef(ctx, txRef)
	if err == nil && tx.Status == "completed" {
		log.Printf("[WebhookService] Payment already completed for tx: %s", txRef)
		return nil
	}

	// 1. Update transaction status
	_, err = s.queries.UpdatePaymentTransactionStatus(ctx, persistence.UpdatePaymentTransactionStatusParams{
		TxRef:                  txRef,
		Status:                 "completed",
		ProviderTransactionRef: pgtype.Text{String: providerRef, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// 2. Update invoice status
	invoice, err := s.queries.UpdateInvoiceStatus(ctx, persistence.UpdateInvoiceStatusParams{
		InvoiceNumber: txRef,
		Status:        persistence.InvoiceStatusPaid,
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to update invoice status: %v", err)
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	// 3. Activate Subscription
	now := time.Now()
	expiresAt := now.AddDate(0, 1, 0)

	_, err = s.queries.UpdateSubscription(ctx, persistence.UpdateSubscriptionParams{
		ID:                 invoice.SubscriptionID,
		Status:             persistence.NullSubscriptionStatus{SubscriptionStatus: persistence.SubscriptionStatusActive, Valid: true},
		CurrentPeriodStart: pgtype.Timestamp{Time: now, Valid: true},
		CurrentPeriodEnd:   pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	// 4. Mark old subscriptions as updated
	err = s.queries.UpdateOldSubscriptionsStatus(ctx, persistence.UpdateOldSubscriptionsStatusParams{
		OwnerID: invoice.OwnerID,
		ID:      invoice.SubscriptionID,
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to update old subscriptions status: %v", err)
	}

	// 5. Ensure user is active
	_, err = s.queries.UpdateUser(ctx, persistence.UpdateUserParams{
		ID:       invoice.OwnerID,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to activate user: %v", err)
	}

	user, err := s.queries.GetUserByID(ctx, invoice.OwnerID)
	if err != nil {
		log.Printf("[WebhookService] Failed to get user for email: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 6. Send confirmation email
	go func() {
		amount, _ := utils.NumericToFloat(invoice.Amount)
		if err := s.emailService.SendPaymentSuccessEmail(
			context.Background(),
			user.Email,
			user.FullName,
			invoice.InvoiceNumber,
			amount,
			invoice.Currency,
		); err != nil {
			log.Printf("[WebhookService] Failed to send payment success email: %v", err)
		}
	}()

	log.Printf("[WebhookService] Successfully completed payment for tx: %s", txRef)
	return nil
}

func (s *WebhookService) VerifyPayment(ctx context.Context, txRef string) (string, bool, error) {
	verifyURL := "https://api.chapa.co/v1/transaction/verify/" + txRef

	req, _ := http.NewRequest("GET", verifyURL, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	var verifyResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return "", false, err
	}

	if status, ok := verifyResp["status"].(string); !ok || status != "success" {
		return "", false, nil
	}

	data, ok := verifyResp["data"].(map[string]interface{})
	if !ok {
		return "", false, nil
	}

	providerRef, _ := data["reference"].(string)

	if data["status"] == "success" {
		return providerRef, true, nil
	}

	return providerRef, false, nil
}

func (s *WebhookService) handlePaymentFailed(ctx context.Context, data models.ChapaWebhookData) error {
	log.Printf("[WebhookService] Processing payment failure for tx: %s", data.TxRef)

	// Update transaction status
	_, err := s.queries.UpdatePaymentTransactionStatus(ctx, persistence.UpdatePaymentTransactionStatusParams{
		TxRef:                  data.TxRef,
		Status:                 "failed",
		ProviderTransactionRef: pgtype.Text{String: data.Reference, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Update invoice status
	invoice, err := s.queries.UpdateInvoiceStatus(ctx, persistence.UpdateInvoiceStatusParams{
		InvoiceNumber: data.TxRef,
		Status:        persistence.InvoiceStatusFailed,
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to update invoice status: %v", err)
	}

	// Mark subscription as past_due
	_, err = s.queries.UpdateSubscription(ctx, persistence.UpdateSubscriptionParams{
		ID:     invoice.SubscriptionID,
		Status: persistence.NullSubscriptionStatus{SubscriptionStatus: persistence.SubscriptionStatusPastDue, Valid: true},
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to mark subscription as past_due: %v", err)
	}

	// Schedule retry job
	_, err = s.queries.CreatePaymentRetryJob(ctx, persistence.CreatePaymentRetryJobParams{
		SubscriptionID: invoice.SubscriptionID,
		ScheduledFor:   pgtype.Timestamp{Time: time.Now().AddDate(0, 0, 1), Valid: true}, // Retry in 1 day
	})
	if err != nil {
		log.Printf("[WebhookService] Warning: Failed to schedule retry job: %v", err)
	}

	// Send failure email
	go func() {
		// Get user details
		user, err := s.queries.GetUserByID(context.Background(), invoice.OwnerID)
		if err != nil {
			log.Printf("[WebhookService] Failed to get user for email: %v", err)
			return
		}

		updatePaymentURL := os.Getenv("APP_BASE_URL") + "/billing"
		if err := s.emailService.SendPaymentFailedEmail(
			context.Background(),
			user.Email,
			user.FullName,
			updatePaymentURL,
		); err != nil {
			log.Printf("[WebhookService] Failed to send payment failed email: %v", err)
		}
	}()

	return nil
}

func (s *WebhookService) handlePaymentPending(ctx context.Context, data models.ChapaWebhookData) error {
	log.Printf("[WebhookService] Processing payment pending for tx: %s", data.TxRef)

	// Send pending email
	go func() {
		// Get transaction to find user
		tx, err := s.queries.GetPaymentTransactionByTxRef(context.Background(), data.TxRef)
		if err != nil {
			log.Printf("[WebhookService] Failed to get transaction for email: %v", err)
			return
		}

		user, err := s.queries.GetUserByID(context.Background(), tx.OwnerID)
		if err != nil {
			log.Printf("[WebhookService] Failed to get user for email: %v", err)
			return
		}

		if err := s.emailService.SendPaymentPendingEmail(
			context.Background(),
			user.Email,
			user.FullName,
		); err != nil {
			log.Printf("[WebhookService] Failed to send payment pending email: %v", err)
		}
	}()

	return nil
}
