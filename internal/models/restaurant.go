package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID            uuid.UUID       `json:"id"`
	OwnerID       uuid.UUID       `json:"owner_id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description,omitempty"`
	CuisineType   string          `json:"cuisine_type,omitempty"`
	Phone         string          `json:"phone,omitempty"`
	Email         string          `json:"email,omitempty"`
	Website       string          `json:"website,omitempty"`
	Address       string          `json:"address,omitempty"`
	City          string          `json:"city,omitempty"`
	Country       string          `json:"country,omitempty"`
	LogoURL       string          `json:"logo_url,omitempty"`
	CoverImageURL string          `json:"cover_image_url,omitempty"`
	ThemeSettings json.RawMessage `json:"theme_settings"`
	IsPublished   bool            `json:"is_published"`
	Status        string          `json:"status"`
	ViewCount     int32           `json:"view_count"`
	RankScore     float64         `json:"rank_score"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CreateRestaurantRequest struct {
	Name          string          `json:"name" binding:"required"`
	Slug          string          `json:"slug" binding:"required"`
	Description   string          `json:"description,omitempty"`
	CuisineType   string          `json:"cuisine_type,omitempty"`
	Phone         string          `json:"phone,omitempty"`
	Email         string          `json:"email,omitempty"`
	Website       string          `json:"website,omitempty"`
	Address       string          `json:"address,omitempty"`
	City          string          `json:"city,omitempty"`
	Country       string          `json:"country,omitempty"`
	ThemeSettings json.RawMessage `json:"theme_settings,omitempty"`
}

type UpdateRestaurantRequest struct {
	Name          *string         `json:"name,omitempty"`
	Description   *string         `json:"description,omitempty"`
	CuisineType   *string         `json:"cuisine_type,omitempty"`
	Phone         *string         `json:"phone,omitempty"`
	Email         *string         `json:"email,omitempty"`
	Website       *string         `json:"website,omitempty"`
	Address       *string         `json:"address,omitempty"`
	City          *string         `json:"city,omitempty"`
	Country       *string         `json:"country,omitempty"`
	LogoURL       *string         `json:"logo_url,omitempty"`
	CoverImageURL *string         `json:"cover_image_url,omitempty"`
	ThemeSettings json.RawMessage `json:"theme_settings,omitempty"`
	IsPublished   *bool           `json:"is_published,omitempty"`
}

type RestaurantFilter struct {
	OwnerID     *string `json:"owner_id,omitempty"`
	Status      *string `json:"status,omitempty"`
	Search      *string `json:"search,omitempty"`
	CuisineType *string `json:"cuisine_type,omitempty"`
	City        *string `json:"city,omitempty"`
	Country     *string `json:"country,omitempty"`
	IsPublished *bool   `json:"is_published,omitempty"`
}
