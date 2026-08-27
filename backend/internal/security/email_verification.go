package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GenerateEmailVerificationToken creates a secure random token for email verification.
func GenerateEmailVerificationToken() (string, string, error) {
	// This creates 32 random bytes for the verification token.
	tokenBytes := make([]byte, 32)

	// This fills the byte slice using the operating system's secure random generator.
	if _, err := rand.Read(tokenBytes); err != nil {
		// This returns the error if secure token generation fails.
		return "", "", err
	}

	// This converts the random bytes into a URL-safe token for the email link.
	rawToken := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// This hashes the raw token before storing it in the database.
	tokenHash := HashEmailVerificationToken(rawToken)

	// This returns the raw token for the email and the hash for the database.
	return rawToken, tokenHash, nil
}

// GenerateVerificationSessionToken creates a secure random token for browser verification sessions.
func GenerateVerificationSessionToken() (string, string, error) {
	return GenerateEmailVerificationToken()
}

// HashEmailVerificationToken hashes a raw verification token for database lookup.
func HashEmailVerificationToken(rawToken string) string {
	// This calculates a SHA-256 hash of the raw token.
	hash := sha256.Sum256([]byte(rawToken))

	// This converts the hash bytes into a stable hex string.
	return hex.EncodeToString(hash[:])
}

// HashVerificationSessionToken hashes a raw temporary session token for database lookup.
func HashVerificationSessionToken(rawToken string) string {
	return HashEmailVerificationToken(rawToken)
}
