package models

import "time"

// Alert stores an important monitoring event, such as a service going down.
type Alert struct {
	// ID is the unique primary key for the alert.
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID links the alert to the owner of the monitored service.
	UserID uint `gorm:"not null;index" json:"user_id"`

	// ServiceID links the alert to the monitored service.
	ServiceID uint `gorm:"not null;index" json:"service_id"`

	// HealthCheckID optionally links the alert to the health check that triggered it.
	HealthCheckID *uint `gorm:"index" json:"health_check_id,omitempty"`

	// Type stores the kind of alert, such as service_down.
	Type string `gorm:"size:50;not null;index" json:"type"`

	// Severity stores how serious the alert is.
	Severity string `gorm:"size:20;not null;default:critical" json:"severity"`

	// Title stores a short human-readable alert title.
	Title string `gorm:"size:150;not null" json:"title"`

	// Message stores a longer explanation of what happened.
	Message string `gorm:"not null" json:"message"`

	// CreatedAt stores when the alert was created.
	CreatedAt time.Time `json:"created_at"`

	// Service allows the API to include basic service details with the alert.
	Service Service `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"service,omitempty"`

	// User allows GORM to understand the user relationship, but we do not return it in JSON.
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// HealthCheck allows GORM to understand the health check relationship, but we do not return it in JSON.
	HealthCheck HealthCheck `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
