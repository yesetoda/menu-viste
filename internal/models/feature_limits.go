package models

// FeatureLimits represents the features available in a subscription tier
type FeatureLimits struct {
	// Core Limits
	MaxRestaurants   int `json:"max_restaurants"`
	MaxCategories    int `json:"max_categories"`
	MaxMenuItems     int `json:"max_menu_items"`
	MaxStaffAccounts int `json:"max_staff_accounts"`

	// Activity Log Features
	ActivityLogEnabled bool `json:"activity_log_enabled"`
	ActivityLogDays    int  `json:"activity_log_days"`

	// Analytics Features
	AnalyticsEnabled     bool `json:"analytics_enabled"`
	AnalyticsHistoryDays int  `json:"analytics_history_days"`

	// Ranking & Visibility
	SearchPriorityBoost int `json:"search_priority_boost"`
}
