package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID       `json:"id"`
	RestaurantID uuid.UUID       `json:"restaurant_id"`
	CategoryID   uuid.UUID       `json:"category_id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Price        float64         `json:"price"`
	Currency     string          `json:"currency"`
	Images       json.RawMessage `json:"images"`
	Allergens    json.RawMessage `json:"allergens"`
	DietaryTags  json.RawMessage `json:"dietary_tags"`
	SpiceLevel   int32           `json:"spice_level"`
	Calories     int32           `json:"calories,omitempty"`
	IsAvailable  bool            `json:"is_available"`
	DisplayOrder int32           `json:"display_order"`
	ViewCount    int32           `json:"view_count"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CreateMenuItemRequest struct {
	CategoryID   uuid.UUID       `json:"category_id" binding:"required"`
	Name         string          `json:"name" binding:"required"`
	Description  string          `json:"description,omitempty"`
	Price        float64         `json:"price" binding:"required"`
	Currency     string          `json:"currency" binding:"required"`
	Images       json.RawMessage `json:"images,omitempty"`
	Allergens    json.RawMessage `json:"allergens,omitempty"`
	DietaryTags  json.RawMessage `json:"dietary_tags,omitempty"`
	SpiceLevel   int32           `json:"spice_level"`
	Calories     int32           `json:"calories,omitempty"`
	IsAvailable  bool            `json:"is_available"`
	DisplayOrder int32           `json:"display_order"`
}

type UpdateMenuItemRequest struct {
	CategoryID   *uuid.UUID      `json:"category_id,omitempty"`
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Price        *float64        `json:"price,omitempty"`
	Images       json.RawMessage `json:"images,omitempty"`
	Allergens    json.RawMessage `json:"allergens,omitempty"`
	DietaryTags  json.RawMessage `json:"dietary_tags,omitempty"`
	SpiceLevel   *int32          `json:"spice_level,omitempty"`
	Calories     *int32          `json:"calories,omitempty"`
	IsAvailable  *bool           `json:"is_available,omitempty"`
	DisplayOrder *int32          `json:"display_order,omitempty"`
}
