package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID             uuid.UUID       `json:"id"`
	RestaurantID   uuid.UUID       `json:"restaurant_id"`
	UserID         uuid.UUID       `json:"user_id"`
	ActionType     string          `json:"action_type"`
	ActionCategory string          `json:"action_category"`
	Description    string          `json:"description,omitempty"`
	TargetType     string          `json:"target_type,omitempty"`
	TargetID       *uuid.UUID      `json:"target_id,omitempty"`
	TargetName     string          `json:"target_name,omitempty"`
	BeforeValue    json.RawMessage `json:"before_value,omitempty"`
	AfterValue     json.RawMessage `json:"after_value,omitempty"`
	IPAddress      string          `json:"ip_address,omitempty"`
	UserAgent      string          `json:"user_agent,omitempty"`
	DeviceType     string          `json:"device_type,omitempty"`
	Browser        string          `json:"browser,omitempty"`
	Os             string          `json:"os,omitempty"`
	Success        bool            `json:"success"`
	CreatedAt      time.Time       `json:"created_at"`

	// Join fields
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

type CreateActivityLogRequest struct {
	RestaurantID   uuid.UUID              `json:"restaurant_id"`
	UserID         uuid.UUID              `json:"user_id"`
	ActionType     string                 `json:"action_type"`
	ActionCategory string                 `json:"action_category"`
	Description    string                 `json:"description"`
	TargetType     string                 `json:"target_type"`
	TargetID       uuid.UUID              `json:"target_id"`
	TargetName     string                 `json:"target_name"`
	BeforeValue    map[string]interface{} `json:"before_value"`
	AfterValue     map[string]interface{} `json:"after_value"`
	IPAddress      string                 `json:"ip_address"`
	UserAgent      string                 `json:"user_agent"`
	DeviceType     string                 `json:"device_type"`
	Browser        string                 `json:"browser"`
	OS             string                 `json:"os"`
	Success        bool                   `json:"success"`
}
