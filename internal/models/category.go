package models

import (
	"mime/multipart"
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
	Name         string                `form:"name" binding:"required"`
	Description  string                `form:"description,omitempty"`
	Icon         *multipart.FileHeader `form:"icon,omitempty"`
	DisplayOrder int32                 `form:"display_order"`
	IsActive     bool                  `form:"is_active"`
}

type UpdateCategoryRequest struct {
	Name         *string               `form:"name,omitempty"`
	Description  *string               `form:"description,omitempty"`
	Icon         *multipart.FileHeader `form:"icon,omitempty"`
	DisplayOrder *int32                `form:"display_order,omitempty"`
	IsActive     *bool                 `form:"is_active,omitempty"`
}
