package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
)

type Subscription struct {
	ID                            uuid.UUID          `json:"id"`
	OwnerID                       uuid.UUID          `json:"owner_id"`
	PlanID                        uuid.UUID          `json:"plan_id"`
	Status                        SubscriptionStatus `json:"status"`
	CurrentPeriodStart            time.Time          `json:"current_period_start"`
	CurrentPeriodEnd              time.Time          `json:"current_period_end"`
	TrialEnd                      *time.Time         `json:"trial_end,omitempty"`
	CancelledAt                   *time.Time         `json:"cancelled_at,omitempty"`
	PaymentProviderSubscriptionID string             `json:"payment_provider_subscription_id,omitempty"`
	CreatedAt                     time.Time          `json:"created_at"`
	UpdatedAt                     time.Time          `json:"updated_at"`

	// Join fields
	PlanName string        `json:"plan_name,omitempty"`
	PlanSlug string        `json:"plan_slug,omitempty"`
	Features FeatureLimits `json:"features,omitempty"`
}

type SubscriptionPlan struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Slug         string        `json:"slug"`
	Description  string        `json:"description,omitempty"`
	PriceMonthly int32         `json:"price_monthly"`
	PriceAnnual  int32         `json:"price_annual,omitempty"`
	Currency     string        `json:"currency"`
	Features     FeatureLimits `json:"features"`
	DisplayOrder int32         `json:"display_order"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type SubscriptionDetailsResponse struct {
	PlanName      string             `json:"plan_name"`
	PlanSlug      string             `json:"plan_slug"`
	Price         int32              `json:"price"`
	Currency      string             `json:"currency"`
	Status        SubscriptionStatus `json:"status"`
	StartDate     time.Time          `json:"start_date"`
	EndDate       time.Time          `json:"end_date"`
	TrialEnd      *time.Time         `json:"trial_end,omitempty"`
	DaysRemaining int                `json:"days_remaining"`
	Features      FeatureLimits      `json:"features"`
}
