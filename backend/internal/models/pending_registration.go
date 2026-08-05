package models

import "time"

// PendingRegistration stores a registration that has not been email-verified yet.
type PendingRegistration struct {
	// ID is the unique primary key for the pending registration.
	ID uint `gorm:"primaryKey" json:"id"`

	// Name stores the user's submitted name.
	Name string `gorm:"size:100;not null" json:"name"`

	// Email stores the user's submitted email address.
	Email string `gorm:"size:255;uniqueIndex;not null" json:"email"`

	// PasswordHash stores the hashed password while the registration is pending.
	PasswordHash string `gorm:"not null" json:"-"`

	// VerificationTokenHash stores a hashed verification token, not the raw token.
	VerificationTokenHash string `gorm:"size:64;uniqueIndex;not null" json:"-"`

	// VerificationExpiresAt stores when the verification token expires.
	VerificationExpiresAt time.Time `gorm:"not null" json:"-"`

	// VerificationSentAt stores when the most recent verification email was sent.
	VerificationSentAt time.Time `gorm:"not null" json:"-"`

	// VerificationDailyCount stores how many verification emails were sent in the current 24-hour window.
	VerificationDailyCount int `gorm:"not null;default:0" json:"-"`

	// VerificationWindowStart stores when the current resend-count window started.
	VerificationWindowStart time.Time `gorm:"not null" json:"-"`

	// CreatedAt stores when the pending registration was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores when the pending registration was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
