package models

import "time"

const (
	VerificationSessionStatusPending  = "pending"
	VerificationSessionStatusVerified = "verified"
	VerificationSessionStatusExpired  = "expired"
	VerificationSessionStatusConsumed = "consumed"
)

// VerificationSession stores a short-lived browser session waiting for email verification.
type VerificationSession struct {
	// ID is the unique primary key for the verification session.
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID links the session to the verified user once email verification succeeds.
	UserID *uint `gorm:"index" json:"user_id,omitempty"`

	// User stores the optional verified user relationship.
	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	// Email stores the email address used to start registration.
	Email string `gorm:"size:255;not null;index" json:"email"`

	// TokenHash stores the hashed temporary session token, not the raw token.
	TokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// Status stores the current lifecycle state for the temporary session.
	Status string `gorm:"size:20;not null;index" json:"status"`

	// ExpiresAt stores when the temporary session stops being usable.
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`

	// ConsumedAt stores when the temporary session was exchanged for normal login.
	ConsumedAt *time.Time `gorm:"index" json:"consumed_at,omitempty"`

	// CreatedAt stores when the verification session was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores when the verification session was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
