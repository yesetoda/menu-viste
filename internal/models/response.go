package models

import "fmt"

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Meta       *Meta       `json:"meta,omitempty"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
	Code       string `json:"code,omitempty"`
}

// Meta contains pagination and filtering metadata
type Meta struct {
	Page         int  `json:"page"`
	PageSize     int  `json:"page_size"`
	TotalPages   int  `json:"total_pages"`
	TotalRecords int  `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

// PaginationParams contains query parameters for pagination
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

// FilterParams contains common filter parameters
type FilterParams struct {
	ID             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Status         *string `json:"status,omitempty"`
	OwnerID        *string `json:"owner_id,omitempty"`
	RestaurantID   *string `json:"restaurant_id,omitempty"`
	PlanID         *string `json:"plan_id,omitempty"`
	SubscriptionID *string `json:"subscription_id,omitempty"`
	Role           *string `json:"role,omitempty"`
	Search         *string `json:"search,omitempty"`
	Email          *string `json:"email,omitempty"`
	Slug           *string `json:"slug,omitempty"`
	CategoryID     *string `json:"category_id,omitempty"`
	IsAvailable    *bool   `json:"is_available,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
	IsPublished    *bool   `json:"is_published,omitempty"`
	CuisineType    *string `json:"cuisine_type,omitempty"`
	City           *string `json:"city,omitempty"`
	Country        *string `json:"country,omitempty"`
	ActionType     *string `json:"action_type,omitempty"`
	ActionCategory *string `json:"action_category,omitempty"`
	EventType      *string `json:"event_type,omitempty"`
	UserID         *string `json:"user_id,omitempty"`
	TargetType     *string `json:"target_type,omitempty"`
	TargetID       *string `json:"target_id,omitempty"`
	Success        *bool   `json:"success,omitempty"`
	VisitorID      *string `json:"visitor_id,omitempty"`
	SessionID      *string `json:"session_id,omitempty"`
	SortBy         string  `json:"sort_by"`
	SortDir        string  `json:"sort_dir"`
}

// NewPaginationParams creates pagination params with defaults
func NewPaginationParams(page, pageSize int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// CalculateMeta calculates pagination metadata
func CalculateMeta(page, pageSize, totalRecords int) *Meta {
	fmt.Println("[CalculateMeta] page:", page)
	fmt.Println("[CalculateMeta] pageSize:", pageSize)
	fmt.Println("[CalculateMeta] totalRecords:", totalRecords)
	totalPages := (totalRecords + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return &Meta{
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}
}
