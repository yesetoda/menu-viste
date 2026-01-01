package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type PaymentTransaction struct {
	ID                     uuid.UUID     `json:"id"`
	OwnerID                uuid.UUID     `json:"owner_id"`
	Amount                 float64       `json:"amount"`
	Currency               string        `json:"currency"`
	Status                 PaymentStatus `json:"status"`
	TxRef                  string        `json:"tx_ref"`
	Reference              string        `json:"reference,omitempty"`
	ProviderTransactionRef string        `json:"provider_transaction_ref,omitempty"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
}

// ChapaWebhookPayload represents the payload sent by Chapa webhooks
type ChapaWebhookPayload struct {
	Event string           `json:"event"`
	Data  ChapaWebhookData `json:"data"`
}

type ChapaWebhookData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	TxRef     string `json:"tx_ref"`
	Reference string `json:"reference"`
	CreatedAt string `json:"created_at"`
}
