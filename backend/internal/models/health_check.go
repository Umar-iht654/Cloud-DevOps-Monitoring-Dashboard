package models

import "time"

type HealthCheck struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ServiceID      uint      `gorm:"not null;index" json:"service_id"`
	Status         string    `gorm:"size:20;not null" json:"status"`
	HTTPStatusCode *int      `json:"http_status_code"`
	ResponseTimeMs *int      `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message"`
	CheckedAt      time.Time `gorm:"not null;index" json:"checked_at"`
	CreatedAt      time.Time `json:"created_at"`
}
