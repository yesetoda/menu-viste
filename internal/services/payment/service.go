package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	var latestSub persistence.GetLatestSubscriptionByOwnerRow
	var activeSub persistence.GetActiveSubscriptionByOwnerRow
	var err error

	// 1. Get latest subscription (any status) for reuse logic
	latestSub, err = s.queries.GetLatestSubscriptionByOwner(ctx, input.OwnerID)
	hasLatest := err == nil

	// 2. Get active subscription for amount calculation (upgrades/renewals)
	activeSub, err = s.queries.GetActiveSubscriptionByOwner(ctx, input.OwnerID)
	hasActive := err == nil

	newPlan, err := s.queries.GetSubscriptionPlanBySlug(ctx, input.Plan)
	if err != nil {
		return nil, fmt.Errorf("invalid plan: %s", input.Plan)
	}

	// 3. Reuse incomplete subscription if it matches the plan
	if hasLatest && latestSub.Status == persistence.SubscriptionStatusIncomplete && latestSub.PlanSlug == newPlan.Slug {
		input.SubscriptionID = latestSub.ID
		log.Printf("[PaymentService] Reusing incomplete subscription: %v", latestSub.ID)
	} else {
		// 4. Create a new subscription
		createSubParams := persistence.CreateSubscriptionParams{
			OwnerID:            input.OwnerID,
			PlanID:             newPlan.ID,
			Status:             persistence.SubscriptionStatusIncomplete,
			CurrentPeriodStart: pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		// If there's an active subscription, use its end date if it's close to expiring
		if hasActive {
			if (time.Until(activeSub.CurrentPeriodEnd.Time).Hours() / 24) <= 5 {
				createSubParams.CurrentPeriodEnd = pgtype.Timestamp{Time: time.Now().AddDate(0, 1, 0), Valid: true}
			} else {
				createSubParams.CurrentPeriodEnd = activeSub.CurrentPeriodEnd
			}
		} else {
			createSubParams.CurrentPeriodEnd = pgtype.Timestamp{Time: time.Now().AddDate(0, 1, 0), Valid: true}
		}

		newSub, err := s.queries.CreateSubscription(ctx, createSubParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
		input.SubscriptionID = newSub.ID
		log.Printf("[PaymentService] Created new incomplete subscription: %v", newSub.ID)
	}

	email, name, err := s.getUserDetails(ctx, input)
	if err != nil {
		return nil, err
	}

	amount, err := s.calculateAmount(ctx, input.Type, activeSub, newPlan)
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

func (s *Service) calculateAmount(ctx context.Context, newPlanType string, currentPlan persistence.GetActiveSubscriptionByOwnerRow, newPlan persistence.SubscriptionPlan) (float64, error) {
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
