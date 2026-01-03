package utils

import (
	"fmt"
	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// AllowedKeys defines the allowed filter keys for each entity
var AllowedKeys = map[string][]string{
	"restaurants":      {"id", "name", "status", "owner_id", "slug", "search", "cuisine_type", "city", "country", "is_published"},
	"users":            {"id", "email", "status", "role", "search", "owner_id", "restaurant_id", "is_active"},
	"categories":       {"id", "restaurant_id", "name", "search", "is_active"},
	"items":            {"id", "restaurant_id", "category_id", "name", "is_available", "search"},
	"activity_logs":    {"id", "restaurant_id", "user_id", "action_type", "action_category", "target_type", "target_id", "success", "search"},
	"analytics_events": {"id", "restaurant_id", "event_type", "visitor_id", "session_id", "target_id", "country", "city", "search"},
	"subscriptions":    {"id", "owner_id", "plan_id", "status"},
	"invoices":         {"id", "subscription_id", "owner_id", "status"},
}

// FilterBuilder handles validation and construction of filters
type FilterBuilder struct {
	entity      string
	allowedKeys []string
}

// NewFilterBuilder creates a new FilterBuilder for an entity
func NewFilterBuilder(entity string) *FilterBuilder {
	return &FilterBuilder{
		entity:      entity,
		allowedKeys: AllowedKeys[entity],
	}
}

// ValidateAndParse parses query parameters and validates them against allowed keys
func (fb *FilterBuilder) ValidateAndParse(c *gin.Context) (models.FilterParams, error) {
	filters := models.FilterParams{
		SortBy:  c.DefaultQuery("sort_by", "created_at"),
		SortDir: c.DefaultQuery("sort_dir", "desc"),
	}

	queryParams := c.Request.URL.Query()
	for key := range queryParams {
		// Skip pagination and sorting keys
		if key == "page" || key == "page_size" || key == "sort_by" || key == "sort_dir" {
			continue
		}

		if !fb.isAllowed(key) {
			return filters, fmt.Errorf("filter key '%s' is not allowed for %s", key, fb.entity)
		}

		val := c.Query(key)
		if val == "" {
			continue
		}

		switch key {
		case "id":
			filters.ID = &val
		case "name":
			filters.Name = &val
		case "status":
			filters.Status = &val
		case "owner_id":
			filters.OwnerID = &val
		case "restaurant_id":
			filters.RestaurantID = &val
		case "plan_id":
			filters.PlanID = &val
		case "subscription_id":
			filters.SubscriptionID = &val
		case "role":
			filters.Role = &val
		case "search":
			filters.Search = &val
		case "email":
			filters.Email = &val
		case "slug":
			filters.Slug = &val
		case "category_id":
			filters.CategoryID = &val
		case "is_available":
			isAvailable := val == "true"
			filters.IsAvailable = &isAvailable
		case "is_active":
			isActive := val == "true"
			filters.IsActive = &isActive
		case "is_published":
			isPublished := val == "true"
			filters.IsPublished = &isPublished
		case "cuisine_type":
			filters.CuisineType = &val
		case "city":
			filters.City = &val
		case "country":
			filters.Country = &val
		case "action_type":
			filters.ActionType = &val
		case "action_category":
			filters.ActionCategory = &val
		case "event_type":
			filters.EventType = &val
		case "user_id":
			filters.UserID = &val
		case "target_type":
			filters.TargetType = &val
		case "target_id":
			filters.TargetID = &val
		case "success":
			success := val == "true"
			filters.Success = &success
		case "visitor_id":
			filters.VisitorID = &val
		case "session_id":
			filters.SessionID = &val
		}
	}

	return filters, nil
}

func (fb *FilterBuilder) isAllowed(key string) bool {
	for _, k := range fb.allowedKeys {
		if k == key {
			return true
		}
	}
	return false
}

// ToUUID converts a string pointer to uuid.UUID
func ToUUID(s *string) uuid.UUID {
	if s == nil || *s == "" {
		return uuid.Nil
	}
	u, err := uuid.Parse(*s)
	if err != nil {
		return uuid.Nil
	}
	return u
}

// FromUUID converts uuid.UUID to string
func FromUUID(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

// ToText converts a string pointer to pgtype.Text
func ToText(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// ToRestaurantStatus converts a string pointer to persistence.RestaurantStatus with NULL handling
// func ToRestaurantStatus(s *string) (persistence.RestaurantStatus, bool) {
// 	if s == nil || *s == "" {
// 		return "", false
// 	}
// 	return persistence.RestaurantStatus(*s), true
// }

// ToUserRole converts a string pointer to persistence.UserRole with NULL handling
func ToUserRole(s *string) (persistence.UserRole, bool) {
	if s == nil || *s == "" {
		return "", false
	}
	return persistence.UserRole(*s), true
}

// ToInvoiceStatus converts a string pointer to persistence.InvoiceStatus with NULL handling
func ToInvoiceStatus(s *string) (persistence.InvoiceStatus, bool) {
	if s == nil || *s == "" {
		return "", false
	}
	return persistence.InvoiceStatus(*s), true
}

// ToSubscriptionStatus converts a string pointer to persistence.SubscriptionStatus with NULL handling
func ToSubscriptionStatus(s *string) (persistence.SubscriptionStatus, bool) {
	if s == nil || *s == "" {
		return "", false
	}
	return persistence.SubscriptionStatus(*s), true
}

// ToBool converts a bool pointer to pgtype.Bool
func ToBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}
