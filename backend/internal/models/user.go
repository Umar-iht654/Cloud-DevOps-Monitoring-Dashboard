package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Services     []Service `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"services,omitempty"`

	// Alerts stores alert records owned by this user.
	Alerts []Alert `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"alerts,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
