package models

import (
	"encoding/json"
	"mime/multipart"
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
	Name          string                `form:"name" binding:"required"`
	Slug          string                `form:"slug" binding:"required"`
	Description   string                `form:"description,omitempty"`
	CuisineType   string                `form:"cuisine_type,omitempty"`
	Phone         string                `form:"phone,omitempty"`
	Email         string                `form:"email,omitempty"`
	Website       string                `form:"website,omitempty"`
	Address       string                `form:"address,omitempty"`
	City          string                `form:"city,omitempty"`
	Country       string                `form:"country,omitempty"`
	ThemeSettings json.RawMessage       `form:"theme_settings,omitempty"`
	Logo          *multipart.FileHeader `form:"logo,omitempty"`
	CoverImage    *multipart.FileHeader `form:"cover,omitempty"`
}

type UpdateRestaurantRequest struct {
	Name          *string               `form:"name,omitempty"`
	Description   *string               `form:"description,omitempty"`
	CuisineType   *string               `form:"cuisine_type,omitempty"`
	Phone         *string               `form:"phone,omitempty"`
	Email         *string               `form:"email,omitempty"`
	Website       *string               `form:"website,omitempty"`
	Address       *string               `form:"address,omitempty"`
	City          *string               `form:"city,omitempty"`
	Country       *string               `form:"country,omitempty"`
	ThemeSettings json.RawMessage       `form:"theme_settings,omitempty"`
	IsPublished   *bool                 `form:"is_published,omitempty"`
	Logo          *multipart.FileHeader `form:"logo,omitempty"`
	CoverImage    *multipart.FileHeader `form:"cover,omitempty"`
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
