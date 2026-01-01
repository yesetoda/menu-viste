package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *persistence.Queries
}

func NewService(queries *persistence.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

type InitiatePaymentInput struct {
	OwnerID        uuid.UUID
	SubscriptionID uuid.UUID
	Plan           string // slug
	Email          string
	Name           string
	Type           string // update, upgrade
}

type InitiatePaymentResponse struct {
	CheckoutURL string `json:"checkout_url"`
}

func (s *Service) InitiatePayment(ctx context.Context, input InitiatePaymentInput) (*InitiatePaymentResponse, error) {
	log.Printf("[PaymentService] Initiating payment for owner: %v, plan: %s, type: %s", utils.UUIDToString(input.OwnerID), input.Plan, input.Type)

	var sub persistence.GetSubscriptionByOwnerRow
	if (input.SubscriptionID == uuid.UUID{}) {
		sub, err := s.queries.GetSubscriptionByOwner(ctx, input.OwnerID)
		if err == nil {
			input.SubscriptionID = sub.ID

		} else {
			// If no subscription exists, create one (pending/incomplete)
			// We need the plan details first
			plan, err := s.queries.GetSubscriptionPlanBySlug(ctx, input.Plan)
			if err != nil {
				return nil, fmt.Errorf("invalid plan: %s", input.Plan)
			}

			// Create subscription
			newSub, err := s.queries.CreateSubscription(ctx, persistence.CreateSubscriptionParams{
				OwnerID:            input.OwnerID,
				PlanID:             plan.ID,
				Status:             persistence.SubscriptionStatusIncomplete,
				CurrentPeriodStart: pgtype.Timestamp{Time: time.Now(), Valid: true},
				CurrentPeriodEnd:   pgtype.Timestamp{Time: time.Now().AddDate(0, 1, 0), Valid: true},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create subscription: %w", err)
			}
			input.SubscriptionID = newSub.ID
		}
	}

	newPlan, err := s.queries.GetSubscriptionPlanBySlug(ctx, input.Plan)

	if err != nil {
		return nil, fmt.Errorf("invalid plan: %s", input.Plan)
	}

	email, name, err := s.getUserDetails(ctx, input)
	if err != nil {
		return nil, err
	}
	fmt.Printf("this is the current plan %+v \n", sub)
	fmt.Printf("this is the new plan %+v \n", newPlan)

	amount, err := s.calculateAmount(ctx, input.Type, sub, newPlan)
	if err != nil {
		return nil, err
	}

	txRef := "tx_" + uuid.New().String()
	if err := s.createTransactionRecord(ctx, input, newPlan, amount, txRef); err != nil {
		return nil, err
	}

	checkoutURL, err := s.initializeChapaTransaction(newPlan, email, name, fmt.Sprintf("%.2f", amount), txRef)
	if err != nil {
		return nil, err
	}

	s.createInvoiceRecord(ctx, input, newPlan, amount, txRef)

	return &InitiatePaymentResponse{CheckoutURL: checkoutURL}, nil
}

func (s *Service) calculateAmount(ctx context.Context, newPlanType string, currentPlan persistence.GetSubscriptionByOwnerRow, newPlan persistence.SubscriptionPlan) (float64, error) {
	amount := float64(newPlan.PriceMonthly)
	planType := "update"
	if currentPlan.PlanSlug != newPlan.Slug {
		planType = "upgrade"
	}

	remainingDays := time.Until(currentPlan.CurrentPeriodEnd.Time).Hours() / 24
	if remainingDays <= 5 {
		log.Println(planType, " plan", currentPlan.PlanSlug, "tfb", newPlan.Slug, "==>", amount)
		return amount, nil
	}

	totalDays := currentPlan.CurrentPeriodEnd.Time.Sub(currentPlan.CurrentPeriodStart.Time).Hours() / 24
	if totalDays <= 0 {
		log.Println(planType, "plan", currentPlan.PlanSlug, "-->", newPlan.Slug, "==>", amount)
		return amount, nil
	}

	oldPlan, _ := s.queries.GetSubscriptionPlanBySlug(ctx, currentPlan.PlanSlug)
	priceDiff := float64(newPlan.PriceMonthly - oldPlan.PriceMonthly)
	if priceDiff > 0 {
		if newPlanType != "update" {
			amount = float64(newPlan.PriceMonthly)
			log.Println(planType, "exisitng plan", oldPlan.Slug, "-->", newPlan.Slug, "==>", amount)
		}
	}

	return amount, nil
}

func (s *Service) createTransactionRecord(ctx context.Context, input InitiatePaymentInput, plan persistence.SubscriptionPlan, amount float64, txRef string) error {
	_, err := s.queries.CreatePaymentTransaction(ctx, persistence.CreatePaymentTransactionParams{
		OwnerID:   input.OwnerID,
		Amount:    utils.ToNumeric(amount),
		Currency:  plan.Currency,
		Status:    "pending",
		TxRef:     txRef,
		Reference: pgtype.Text{String: fmt.Sprintf("%s:%s", input.Type, plan.Slug), Valid: true},
	})
	return err
}

func (s *Service) createInvoiceRecord(ctx context.Context, input InitiatePaymentInput, plan persistence.SubscriptionPlan, amount float64, txRef string) {
	_, err := s.queries.CreateInvoice(ctx, persistence.CreateInvoiceParams{
		OwnerID:            input.OwnerID,
		InvoiceNumber:      txRef,
		Amount:             utils.ToNumeric(amount),
		Currency:           plan.Currency,
		Status:             persistence.InvoiceStatusPending,
		SubscriptionID:     input.SubscriptionID,
		BillingPeriodStart: pgtype.Timestamp{Time: time.Now(), Valid: true},
		BillingPeriodEnd:   pgtype.Timestamp{Time: time.Now().AddDate(0, 1, 0), Valid: true},
	})
	if err != nil {
		log.Printf("[PaymentService] Failed to save invoice record: %v", err)
	}
}

func (s *Service) getUserDetails(ctx context.Context, input InitiatePaymentInput) (string, string, error) {
	if input.Email != "" && input.Name != "" {
		return input.Email, input.Name, nil
	}

	user, err := s.queries.GetUserByID(ctx, input.OwnerID)
	if err != nil {
		log.Printf("[PaymentService] User not found: %v", utils.UUIDToString(input.OwnerID))
		return "", "", fmt.Errorf("user not found")
	}

	email := input.Email
	if email == "" {
		email = user.Email
	}

	name := input.Name
	if name == "" {
		name = user.FullName
	}

	return email, name, nil
}

func (s *Service) initializeChapaTransaction(plan persistence.SubscriptionPlan, email, name, amountStr, txRef string) (string, error) {
	payload := map[string]interface{}{
		"amount":       amountStr,
		"currency":     plan.Currency,
		"email":        email,
		"first_name":   name,
		"tx_ref":       txRef,
		"callback_url": os.Getenv("CHAPA_CALLBACK_URL"),
		"return_url":   fmt.Sprintf("%s?tx_ref=%s", os.Getenv("CHAPA_RETURN_URL"), txRef),
		// "cancel_url":   os.Getenv("CHAPA_CANCEL_URL"), // Uncomment if Chapa supports it directly or handle via return_url logic
		"customization": map[string]string{
			"title":       plan.Name + " Plan",
			"description": "Monthly subscription",
		},
	}

	body, _ := json.Marshal(payload)

	reqHttp, _ := http.NewRequest(
		"POST",
		"https://api.chapa.co/v1/transaction/initialize",
		bytes.NewBuffer(body),
	)

	reqHttp.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))
	reqHttp.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(reqHttp)
	if err != nil {
		log.Printf("[PaymentService] Chapa request failed: %v", err)
		return "", fmt.Errorf("failed to initiate payment: %w", err)
	}
	defer resp.Body.Close()

	var chapaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&chapaResp); err != nil {
		log.Printf("[PaymentService] Failed to decode Chapa response: %v", err)
		return "", fmt.Errorf("failed to decode payment response")
	}

	if status, ok := chapaResp["status"].(string); !ok || status != "success" {
		log.Printf("[PaymentService] Chapa returned error: %v", chapaResp)
		return "", fmt.Errorf("payment initiation failed")
	}

	data, ok := chapaResp["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid chapa response data")
	}

	checkoutURL, ok := data["checkout_url"].(string)
	if !ok {
		return "", fmt.Errorf("checkout_url not found")
	}

	return checkoutURL, nil
}

func (s *Service) VerifyPayment(ctx context.Context, txRef string) (bool, error) {
	verifyURL := "https://api.chapa.co/v1/transaction/verify/" + txRef

	req, _ := http.NewRequest("GET", verifyURL, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var verifyResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return false, err
	}

	if status, ok := verifyResp["status"].(string); !ok || status != "success" {
		return false, nil
	}

	data, ok := verifyResp["data"].(map[string]interface{})
	if !ok {
		return false, nil
	}

	if data["status"] == "success" {
		return true, nil
	}

	return false, nil
}

func (s *Service) HandleWebhook(ctx context.Context, txRef string, status string) error {
	log.Printf("[PaymentService] Handling webhook for tx: %s, status: %s", txRef, status)

	verified, err := s.VerifyPayment(ctx, txRef)
	if err != nil {
		return err
	}
	if !verified {
		return fmt.Errorf("payment verification failed")
	}

	// Update invoice status
	invoice, err := s.queries.UpdateInvoiceStatus(ctx, persistence.UpdateInvoiceStatusParams{
		InvoiceNumber: txRef,
		Status:        persistence.InvoiceStatusPaid,
	})
	if err != nil {
		log.Printf("[PaymentService] Failed to update invoice status: %v", err)
		return err
	}

	// Activate Subscription
	// Fetch transaction to get plan details
	tx, err := s.queries.GetPaymentTransactionByTxRef(ctx, txRef)
	if err != nil {
		log.Printf("[PaymentService] Failed to fetch transaction: %v", err)
		return fmt.Errorf("failed to fetch transaction")
	}

	var planID uuid.UUID
	if tx.Reference.Valid {
		log.Printf("[PaymentService] Transaction reference: %s", tx.Reference.String)
		parts := strings.Split(tx.Reference.String, ":")
		if len(parts) == 2 {
			planSlug := parts[1]
			log.Printf("[PaymentService] Parsed plan slug: %s", planSlug)
			plan, err := s.queries.GetSubscriptionPlanBySlug(ctx, planSlug)
			if err != nil {
				log.Printf("[PaymentService] Failed to fetch plan: %v", err)
			} else {
				planID = plan.ID
				log.Printf("[PaymentService] Resolved plan ID: %v", planID)
			}
		} else {
			log.Printf("[PaymentService] Invalid transaction reference format: %s", tx.Reference.String)
		}
	} else {
		log.Printf("[PaymentService] Transaction reference is invalid/empty")
	}

	now := time.Now()
	expiresAt := now.AddDate(0, 1, 0) // 30 days

	updateParams := persistence.UpdateSubscriptionParams{
		ID:                 invoice.SubscriptionID, // Use ID from invoice
		Status:             persistence.NullSubscriptionStatus{SubscriptionStatus: persistence.SubscriptionStatusActive, Valid: true},
		CurrentPeriodStart: pgtype.Timestamp{Time: now, Valid: true},
		CurrentPeriodEnd:   pgtype.Timestamp{Time: expiresAt, Valid: true},
	}

	if planID != uuid.Nil {
		updateParams.PlanID = &planID
	} else {
		// Fallback: keep existing plan ID if we couldn't resolve new one
		// But UpdateSubscription requires PlanID to be set or it will use COALESCE if we pass something?
		// The query uses COALESCE($1, plan_id). If we pass zero UUID, it might update to zero UUID?
		// No, UUID zero value is 0000... which is a valid UUID but likely not a valid FK.
		// We should fetch the current subscription to get the current plan ID if we don't have a new one?
		// Or we can rely on COALESCE if we pass nil? But we can't pass nil for UUID value type.
		// Actually, if we look at the query: plan_id = COALESCE($1, plan_id)
		// If $1 is NULL, it keeps existing.
		// But in Go, uuid.UUID is a value type, not a pointer. It cannot be nil.
		// So we must pass a valid UUID.
		// If planID is Nil (0000...), we should probably fetch the current subscription to get the ID.
		// OR, we can change the query to accept NULL?
		// The query param is `PlanID uuid.UUID`. It's not nullable in the struct generated by sqlc unless we configure it.
		// Let's check UpdateSubscriptionParams again.
		// It is `PlanID uuid.UUID`.
		// So we MUST provide a PlanID.
		// If we don't have a new plan ID, we should use the existing one.
		// So we need to fetch the subscription first.
		sub, err := s.queries.GetSubscriptionByOwner(ctx, invoice.OwnerID)
		if err != nil {
			return fmt.Errorf("failed to fetch subscription: %w", err)
		}
		updateParams.PlanID = &sub.PlanID
	}

	_, err = s.queries.UpdateSubscription(ctx, updateParams)
	if err != nil {
		log.Printf("[PaymentService] Failed to activate subscription: %v", err)
		return err
	}

	return nil
}
