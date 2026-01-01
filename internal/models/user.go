package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleOwner UserRole = "owner"
	RoleStaff UserRole = "staff"
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	FullName      string     `json:"full_name"`
	Role          UserRole   `json:"role"`
	OwnerID       *uuid.UUID `json:"owner_id,omitempty"`
	RestaurantID  *uuid.UUID `json:"restaurant_id,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	FullName string   `json:"full_name" binding:"required"`
	Role     UserRole `json:"role" binding:"required,oneof=admin owner staff"`
	Phone    string   `json:"phone"`
	// PlanSlug string   `json:"plan_slug"` // Optional, defaults to free
}

type CreateStaffRequest struct {
	Email        string    `json:"email" binding:"required,email"`
	FullName     string    `json:"full_name" binding:"required"`
	Phone        string    `json:"phone,omitempty"`
	RestaurantID uuid.UUID `json:"restaurant_id" binding:"required"`
	OwnerID      uuid.UUID `json:"-"` // Set from context
}

type UpdateUserRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CheckoutURL  string `json:"checkout_url,omitempty"`
}
