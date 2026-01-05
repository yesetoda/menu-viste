package models

import (
	"encoding/json"
	"mime/multipart"
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
	CategoryID   uuid.UUID             `form:"category_id" binding:"required"`
	Name         string                `form:"name" binding:"required"`
	Description  string                `form:"description,omitempty"`
	Price        float64               `form:"price" binding:"required"`
	Currency     string                `form:"currency" binding:"required"`
	Images       json.RawMessage       `form:"images,omitempty"`
	Allergens    json.RawMessage       `form:"allergens,omitempty"`
	DietaryTags  json.RawMessage       `form:"dietary_tags,omitempty"`
	SpiceLevel   int32                 `form:"spice_level"`
	Calories     int32                 `form:"calories,omitempty"`
	IsAvailable  bool                  `form:"is_available"`
	DisplayOrder int32                 `form:"display_order"`
	Image        *multipart.FileHeader `form:"image,omitempty"`
}

type UpdateMenuItemRequest struct {
	CategoryID   *uuid.UUID            `form:"category_id,omitempty"`
	Name         *string               `form:"name,omitempty"`
	Description  *string               `form:"description,omitempty"`
	Price        *float64              `form:"price,omitempty"`
	Currency     *string               `form:"currency,omitempty"`
	Images       json.RawMessage       `form:"images,omitempty"`
	Allergens    json.RawMessage       `form:"allergens,omitempty"`
	DietaryTags  json.RawMessage       `form:"dietary_tags,omitempty"`
	SpiceLevel   *int32                `form:"spice_level,omitempty"`
	Calories     *int32                `form:"calories,omitempty"`
	IsAvailable  *bool                 `form:"is_available,omitempty"`
	DisplayOrder *int32                `form:"display_order,omitempty"`
	Image        *multipart.FileHeader `form:"image,omitempty"`
}
