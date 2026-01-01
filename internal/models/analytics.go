package models

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsEvent struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	EventType    string     `json:"event_type"` // page_view, item_view, category_view
	VisitorID    string     `json:"visitor_id"`
	SessionID    uuid.UUID  `json:"session_id"`
	TargetID     *uuid.UUID `json:"target_id,omitempty"`
	IPAddress    string     `json:"ip_address,omitempty"`
	DeviceType   string     `json:"device_type,omitempty"`
	Browser      string     `json:"browser,omitempty"`
	Os           string     `json:"os,omitempty"`
	Country      string     `json:"country,omitempty"`
	City         string     `json:"city,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CreateAnalyticsEventRequest struct {
	RestaurantID uuid.UUID `json:"restaurant_id"`
	EventType    string    `json:"event_type"`
	VisitorID    string    `json:"visitor_id"`
	TargetID     uuid.UUID `json:"target_id"`
	IPAddress    string    `json:"ip_address"`
	DeviceType   string    `json:"device_type"`
	Browser      string    `json:"browser"`
	OS           string    `json:"os"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
}

type AnalyticsAggregate struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	Date         time.Time  `json:"date"`
	Hour         *int32     `json:"hour,omitempty"`
	MetricType   string     `json:"metric_type"`
	TargetID     *uuid.UUID `json:"target_id,omitempty"`
	Value        int32      `json:"value"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
