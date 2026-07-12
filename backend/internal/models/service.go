package models

import "time"

type Service struct {
	ID                   uint          `gorm:"primaryKey" json:"id"`
	UserID               uint          `gorm:"not null;index" json:"user_id"`
	Name                 string        `gorm:"size:100;not null" json:"name"`
	URL                  string        `gorm:"not null" json:"url"`
	ExpectedStatusCode   int           `gorm:"not null;default:200" json:"expected_status_code"`
	SlowThresholdMs      int           `gorm:"not null;default:750" json:"slow_threshold_ms"`
	CheckIntervalSeconds int           `gorm:"not null;default:60" json:"check_interval_seconds"`
	CurrentStatus        string        `gorm:"size:20;not null;default:unknown" json:"current_status"`
	HealthChecks         []HealthCheck `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"health_checks,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
}
