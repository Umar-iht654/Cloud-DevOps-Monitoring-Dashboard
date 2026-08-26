package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/email"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const resendVerificationGenericMessage = "If the account exists and is awaiting verification, a verification email will be sent if allowed."

func setupAuthHandlerTest(t *testing.T) (*AuthHandler, *gin.Engine, *gorm.DB) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.PendingRegistration{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	handler := NewAuthHandler(db, "test-secret", email.NewSender("", "", "", "", "", "http://localhost:5173"))
	router := gin.New()
	router.POST("/api/auth/resend-verification", handler.ResendVerificationEmail)

	return handler, router, db
}

func performResendVerificationRequest(t *testing.T, router *gin.Engine, email string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(gin.H{"email": email})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/auth/resend-verification", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func decodeJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]string {
	t.Helper()

	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	return payload
}

func TestResendVerificationEmailReturnsConflictForVerifiedEmail(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)

	if err := db.Create(&models.User{
		Name:         "Verified User",
		Email:        "verified@example.com",
		PasswordHash: "hashed-password",
	}).Error; err != nil {
		t.Fatalf("failed to create verified user: %v", err)
	}

	response := performResendVerificationRequest(t, router, "  VERIFIED@example.com ")
	payload := decodeJSONResponse(t, response)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, response.Code)
	}

	if payload["code"] != "EMAIL_ALREADY_VERIFIED" {
		t.Fatalf("expected EMAIL_ALREADY_VERIFIED code, got %q", payload["code"])
	}

	if payload["message"] != "This email has already been verified. Please sign in." {
		t.Fatalf("unexpected message %q", payload["message"])
	}
}

func TestResendVerificationEmailResendsForPendingUnverifiedEmail(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	now := time.Now()
	pendingRegistration := models.PendingRegistration{
		Name:                    "Pending User",
		Email:                   "pending@example.com",
		PasswordHash:            "hashed-password",
		VerificationTokenHash:   "old-token-hash",
		VerificationExpiresAt:   now.Add(30 * time.Minute),
		VerificationSentAt:      now.Add(-3 * time.Minute),
		VerificationDailyCount:  1,
		VerificationWindowStart: now.Add(-1 * time.Hour),
	}

	if err := db.Create(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to create pending registration: %v", err)
	}

	response := performResendVerificationRequest(t, router, "pending@example.com")
	payload := decodeJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if payload["message"] != resendVerificationGenericMessage {
		t.Fatalf("unexpected message %q", payload["message"])
	}

	var updatedPendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "pending@example.com").First(&updatedPendingRegistration).Error; err != nil {
		t.Fatalf("failed to reload pending registration: %v", err)
	}

	if updatedPendingRegistration.VerificationDailyCount != 2 {
		t.Fatalf("expected verification daily count 2, got %d", updatedPendingRegistration.VerificationDailyCount)
	}

	if updatedPendingRegistration.VerificationTokenHash == "old-token-hash" {
		t.Fatal("expected verification token hash to be updated")
	}
}

func TestResendVerificationEmailPreservesGenericResponseForNonexistentEmail(t *testing.T) {
	_, router, _ := setupAuthHandlerTest(t)

	response := performResendVerificationRequest(t, router, "missing@example.com")
	payload := decodeJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if payload["message"] != resendVerificationGenericMessage {
		t.Fatalf("unexpected message %q", payload["message"])
	}
}
