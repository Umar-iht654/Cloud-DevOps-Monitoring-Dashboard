package models

import "time"

// HourlyServiceSummary stores compressed monitoring statistics for one service during one hour.
type HourlyServiceSummary struct {
	// ID is the unique primary key for the hourly summary.
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID links the summary to the owner of the monitored service.
	UserID uint `gorm:"not null;index" json:"user_id"`

	// ServiceID links the summary to the monitored service.
	ServiceID uint `gorm:"not null;index;uniqueIndex:idx_hourly_service_period" json:"service_id"`

	// PeriodStart stores the start of the hourly summary period.
	PeriodStart time.Time `gorm:"not null;uniqueIndex:idx_hourly_service_period" json:"period_start"`

	// PeriodEnd stores the end of the hourly summary period.
	PeriodEnd time.Time `gorm:"not null" json:"period_end"`

	// TotalChecks stores the total number of checks during the period.
	TotalChecks int64 `gorm:"not null;default:0" json:"total_checks"`

	// SuccessfulChecks stores checks where the service was online or slow.
	SuccessfulChecks int64 `gorm:"not null;default:0" json:"successful_checks"`

	// FailedChecks stores checks where the service was down.
	FailedChecks int64 `gorm:"not null;default:0" json:"failed_checks"`

	// ResponseTimeSampleCount stores how many checks had a response-time measurement.
	ResponseTimeSampleCount int64 `gorm:"not null;default:0" json:"response_time_sample_count"`

	// AverageResponseTimeMs stores the average response time for measured checks.
	AverageResponseTimeMs int `gorm:"not null;default:0" json:"average_response_time_ms"`

	// MinResponseTimeMs stores the fastest measured response time.
	MinResponseTimeMs *int `json:"min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest measured response time.
	MaxResponseTimeMs *int `json:"max_response_time_ms"`

	// UptimePercentage stores the percentage of checks where the service was reachable.
	UptimePercentage float64 `gorm:"type:numeric(5,2);not null;default:0" json:"uptime_percentage"`

	// CreatedAt stores when the summary row was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores when the summary row was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// Service allows GORM to understand the service relationship.
	Service Service `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"service,omitempty"`

	// User allows GORM to understand the user relationship, but we do not return it in JSON.
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// DailyServiceSummary stores compressed monitoring statistics for one service during one day.
type DailyServiceSummary struct {
	// ID is the unique primary key for the daily summary.
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID links the summary to the owner of the monitored service.
	UserID uint `gorm:"not null;index" json:"user_id"`

	// ServiceID links the summary to the monitored service.
	ServiceID uint `gorm:"not null;index;uniqueIndex:idx_daily_service_period" json:"service_id"`

	// PeriodStart stores the start of the daily summary period.
	PeriodStart time.Time `gorm:"not null;uniqueIndex:idx_daily_service_period" json:"period_start"`

	// PeriodEnd stores the end of the daily summary period.
	PeriodEnd time.Time `gorm:"not null" json:"period_end"`

	// TotalChecks stores the total number of checks during the period.
	TotalChecks int64 `gorm:"not null;default:0" json:"total_checks"`

	// SuccessfulChecks stores checks where the service was online or slow.
	SuccessfulChecks int64 `gorm:"not null;default:0" json:"successful_checks"`

	// FailedChecks stores checks where the service was down.
	FailedChecks int64 `gorm:"not null;default:0" json:"failed_checks"`

	// ResponseTimeSampleCount stores how many checks had a response-time measurement.
	ResponseTimeSampleCount int64 `gorm:"not null;default:0" json:"response_time_sample_count"`

	// AverageResponseTimeMs stores the weighted average response time for measured checks.
	AverageResponseTimeMs int `gorm:"not null;default:0" json:"average_response_time_ms"`

	// MinResponseTimeMs stores the fastest measured response time.
	MinResponseTimeMs *int `json:"min_response_time_ms"`

	// MaxResponseTimeMs stores the slowest measured response time.
	MaxResponseTimeMs *int `json:"max_response_time_ms"`

	// UptimePercentage stores the percentage of checks where the service was reachable.
	UptimePercentage float64 `gorm:"type:numeric(5,2);not null;default:0" json:"uptime_percentage"`

	// CreatedAt stores when the summary row was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores when the summary row was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// Service allows GORM to understand the service relationship.
	Service Service `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"service,omitempty"`

	// User allows GORM to understand the user relationship, but we do not return it in JSON.
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
