package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	DisplayOrder int32     `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    uuid.UUID `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description,omitempty"`
	Icon         string `json:"icon,omitempty"`
	DisplayOrder int32  `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	DisplayOrder *int32  `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}
