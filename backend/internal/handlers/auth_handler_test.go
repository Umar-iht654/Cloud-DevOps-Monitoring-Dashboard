package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/email"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/security"
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

	if err := db.AutoMigrate(&models.User{}, &models.PendingRegistration{}, &models.VerificationSession{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	handler := NewAuthHandler(db, "test-secret", email.NewSender("", "", "", "", "", "http://localhost:5173"))
	router := gin.New()
	router.POST("/api/auth/register", handler.Register)
	router.GET("/api/auth/verify-email", handler.VerifyEmail)
	router.POST("/api/auth/resend-verification", handler.ResendVerificationEmail)
	router.POST("/api/auth/verification-session/status", handler.CheckVerificationSessionStatus)

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

func performRegisterRequest(t *testing.T, router *gin.Engine, email string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(gin.H{
		"name":     "New User",
		"email":    email,
		"password": "Password1",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func performVerifyEmailRequest(t *testing.T, router *gin.Engine, rawToken string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/auth/verify-email?token="+rawToken, nil)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func performVerificationSessionStatusRequest(t *testing.T, router *gin.Engine, rawToken string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(gin.H{"token": rawToken})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/auth/verification-session/status", bytes.NewReader(body))
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

func decodeAnyJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var payload map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	return payload
}

func TestRegisterCreatesVerificationSession(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)

	beforeRegister := time.Now()
	response := performRegisterRequest(t, router, "  NewUser@example.com ")
	afterRegister := time.Now()
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	rawToken, ok := payload["verificationSessionToken"].(string)
	if !ok || rawToken == "" {
		t.Fatal("expected raw verification session token in registration response")
	}

	if payload["verificationSessionExpiresInSeconds"] != float64(180) {
		t.Fatalf("expected 180 second expiry, got %v", payload["verificationSessionExpiresInSeconds"])
	}

	var session models.VerificationSession
	if err := db.Where("email = ?", "newuser@example.com").First(&session).Error; err != nil {
		t.Fatalf("failed to load verification session: %v", err)
	}

	if session.Status != models.VerificationSessionStatusPending {
		t.Fatalf("expected pending status, got %q", session.Status)
	}

	if session.TokenHash == rawToken {
		t.Fatal("verification session stored the raw token")
	}

	if session.TokenHash != security.HashVerificationSessionToken(rawToken) {
		t.Fatal("verification session did not store the returned token hash")
	}

	if session.ExpiresAt.Before(beforeRegister.Add(verificationLifetime)) {
		t.Fatalf("expected expiry at least %s after request start, got %s", verificationLifetime, session.ExpiresAt.Sub(beforeRegister))
	}

	if session.ExpiresAt.After(afterRegister.Add(verificationLifetime)) {
		t.Fatalf("expected expiry no later than %s after response, got %s", verificationLifetime, session.ExpiresAt.Sub(afterRegister))
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "newuser@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	if !pendingRegistration.VerificationExpiresAt.Equal(session.ExpiresAt) {
		t.Fatal("expected pending registration and verification session to share the same expiry")
	}
}

func TestRegisterSetsEmailVerificationTokenExpiryToThreeMinutes(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)

	beforeRegister := time.Now()
	response := performRegisterRequest(t, router, "tokenexpiry@example.com")
	afterRegister := time.Now()

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "tokenexpiry@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	if pendingRegistration.VerificationExpiresAt.Before(beforeRegister.Add(verificationLifetime)) {
		t.Fatalf("expected email token expiry at least %s after request start, got %s", verificationLifetime, pendingRegistration.VerificationExpiresAt.Sub(beforeRegister))
	}

	if pendingRegistration.VerificationExpiresAt.After(afterRegister.Add(verificationLifetime)) {
		t.Fatalf("expected email token expiry no later than %s after response, got %s", verificationLifetime, pendingRegistration.VerificationExpiresAt.Sub(afterRegister))
	}
}

func TestVerifyEmailSucceedsBeforeExpiryAndMarksActiveVerificationSessionVerified(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	response := performRegisterRequest(t, router, "verifiedlater@example.com")
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "verifiedlater@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	emailToken := "email-token"
	pendingRegistration.VerificationTokenHash = security.HashEmailVerificationToken(emailToken)
	pendingRegistration.VerificationExpiresAt = time.Now().Add(verificationLifetime)
	if err := db.Save(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to save pending registration: %v", err)
	}

	verifyResponse := performVerifyEmailRequest(t, router, emailToken)
	if verifyResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, verifyResponse.Code, verifyResponse.Body.String())
	}

	rawSessionToken := payload["verificationSessionToken"].(string)
	var session models.VerificationSession
	if err := db.Where("token_hash = ?", security.HashVerificationSessionToken(rawSessionToken)).First(&session).Error; err != nil {
		t.Fatalf("failed to load verification session: %v", err)
	}

	if session.Status != models.VerificationSessionStatusVerified {
		t.Fatalf("expected verified status, got %q", session.Status)
	}

	if session.UserID == nil {
		t.Fatal("expected verified session to be linked to a user")
	}
}

func TestVerifyEmailFailsAfterExpiryAndDoesNotCreateUser(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	response := performRegisterRequest(t, router, "expiredlink@example.com")

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "expiredlink@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	emailToken := "expired-email-token"
	pendingRegistration.VerificationTokenHash = security.HashEmailVerificationToken(emailToken)
	pendingRegistration.VerificationExpiresAt = time.Now().Add(-time.Second)
	if err := db.Save(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to save expired pending registration: %v", err)
	}

	verifyResponse := performVerifyEmailRequest(t, router, emailToken)
	payload := decodeJSONResponse(t, verifyResponse)

	if verifyResponse.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, verifyResponse.Code, verifyResponse.Body.String())
	}

	if payload["code"] != "VERIFICATION_LINK_EXPIRED" {
		t.Fatalf("expected expired code, got %q", payload["code"])
	}

	if payload["message"] != "This verification link has expired. Please request a new verification email." {
		t.Fatalf("unexpected message %q", payload["message"])
	}

	var user models.User
	if err := db.Where("email = ?", "expiredlink@example.com").First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected no user to be created, got error %v", err)
	}

	var session models.VerificationSession
	if err := db.Where("email = ?", "expiredlink@example.com").First(&session).Error; err != nil {
		t.Fatalf("failed to load verification session: %v", err)
	}

	if session.Status == models.VerificationSessionStatusVerified || session.UserID != nil {
		t.Fatal("expired email verification marked a verification session valid")
	}
}

func TestVerifyEmailDoesNotMarkExpiredVerificationSessionVerified(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	response := performRegisterRequest(t, router, "expiredsession@example.com")

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "expiredsession@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	emailToken := "email-token"
	pendingRegistration.VerificationTokenHash = security.HashEmailVerificationToken(emailToken)
	pendingRegistration.VerificationExpiresAt = time.Now().Add(verificationLifetime)
	if err := db.Save(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to save pending registration: %v", err)
	}

	var session models.VerificationSession
	if err := db.Where("email = ?", "expiredsession@example.com").First(&session).Error; err != nil {
		t.Fatalf("failed to load verification session: %v", err)
	}

	session.ExpiresAt = time.Now().Add(-time.Second)
	if err := db.Save(&session).Error; err != nil {
		t.Fatalf("failed to expire verification session: %v", err)
	}

	verifyResponse := performVerifyEmailRequest(t, router, emailToken)
	if verifyResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, verifyResponse.Code, verifyResponse.Body.String())
	}

	var updatedSession models.VerificationSession
	if err := db.First(&updatedSession, session.ID).Error; err != nil {
		t.Fatalf("failed to reload verification session: %v", err)
	}

	if updatedSession.Status == models.VerificationSessionStatusVerified || updatedSession.UserID != nil {
		t.Fatal("expired verification session was marked valid")
	}
}

func TestVerificationSessionStatusReturnsPendingForActiveSession(t *testing.T) {
	_, router, _ := setupAuthHandlerTest(t)
	response := performRegisterRequest(t, router, "pendingstatus@example.com")
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	statusResponse := performVerificationSessionStatusRequest(t, router, payload["verificationSessionToken"].(string))
	statusPayload := decodeAnyJSONResponse(t, statusResponse)

	if statusResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, statusResponse.Code, statusResponse.Body.String())
	}

	if statusPayload["status"] != models.VerificationSessionStatusPending {
		t.Fatalf("expected pending status, got %v", statusPayload["status"])
	}

	if statusPayload["expiresAt"] == "" {
		t.Fatal("expected pending response to include expiresAt")
	}
}

func TestVerificationSessionStatusExchangesVerifiedSessionForJWTAndAllowsSameTokenRetry(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	registerResponse := performRegisterRequest(t, router, "exchange@example.com")
	registerPayload := decodeAnyJSONResponse(t, registerResponse)

	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, registerResponse.Code)
	}

	var pendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "exchange@example.com").First(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to load pending registration: %v", err)
	}

	emailToken := "exchange-email-token"
	pendingRegistration.VerificationTokenHash = security.HashEmailVerificationToken(emailToken)
	pendingRegistration.VerificationExpiresAt = time.Now().Add(verificationLifetime)
	if err := db.Save(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to save pending registration: %v", err)
	}

	verifyResponse := performVerifyEmailRequest(t, router, emailToken)
	if verifyResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, verifyResponse.Code, verifyResponse.Body.String())
	}

	rawSessionToken := registerPayload["verificationSessionToken"].(string)
	exchangeResponse := performVerificationSessionStatusRequest(t, router, rawSessionToken)
	exchangePayload := decodeAnyJSONResponse(t, exchangeResponse)

	if exchangeResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, exchangeResponse.Code, exchangeResponse.Body.String())
	}

	if exchangePayload["status"] != models.VerificationSessionStatusVerified {
		t.Fatalf("expected verified status, got %v", exchangePayload["status"])
	}

	if exchangePayload["token"] == "" {
		t.Fatal("expected exchange response to include a JWT")
	}

	userPayload, ok := exchangePayload["user"].(map[string]interface{})
	if !ok {
		t.Fatal("expected exchange response to include user")
	}

	if userPayload["email"] != "exchange@example.com" {
		t.Fatalf("expected exchanged user email, got %v", userPayload["email"])
	}

	var session models.VerificationSession
	if err := db.Where("token_hash = ?", security.HashVerificationSessionToken(rawSessionToken)).First(&session).Error; err != nil {
		t.Fatalf("failed to reload verification session: %v", err)
	}

	if session.Status != models.VerificationSessionStatusVerified || session.ConsumedAt != nil {
		t.Fatal("expected verification session to remain verified for retry until expiry")
	}

	secondExchangeResponse := performVerificationSessionStatusRequest(t, router, rawSessionToken)
	secondPayload := decodeAnyJSONResponse(t, secondExchangeResponse)

	if secondExchangeResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, secondExchangeResponse.Code)
	}

	if secondPayload["status"] != models.VerificationSessionStatusVerified {
		t.Fatalf("expected verified status on retry, got %v", secondPayload["status"])
	}

	if secondPayload["token"] == "" {
		t.Fatal("expected retry with same verification session token to issue a JWT")
	}
}

func TestVerificationSessionStatusDoesNotIssueJWTForExpiredSession(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	rawToken := "expired-session-token"
	session := models.VerificationSession{
		Email:     "expiredpoll@example.com",
		TokenHash: security.HashVerificationSessionToken(rawToken),
		Status:    models.VerificationSessionStatusVerified,
		ExpiresAt: time.Now().Add(-time.Second),
	}

	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("failed to create expired verification session: %v", err)
	}

	response := performVerificationSessionStatusRequest(t, router, rawToken)
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	if payload["status"] != models.VerificationSessionStatusExpired {
		t.Fatalf("expected expired status, got %v", payload["status"])
	}

	if payload["code"] != "VERIFICATION_SESSION_EXPIRED" {
		t.Fatalf("expected expired code, got %v", payload["code"])
	}

	if _, hasToken := payload["token"]; hasToken {
		t.Fatal("expired verification session issued a JWT")
	}
}

func TestVerificationSessionStatusHandlesInvalidTokenSafely(t *testing.T) {
	_, router, _ := setupAuthHandlerTest(t)

	response := performVerificationSessionStatusRequest(t, router, "not-a-real-token")
	payload := decodeJSONResponse(t, response)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}

	if payload["message"] != "Invalid verification session" {
		t.Fatalf("unexpected message %q", payload["message"])
	}
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

	var sessionCount int64
	if err := db.Model(&models.VerificationSession{}).Where("email = ?", "verified@example.com").Count(&sessionCount).Error; err != nil {
		t.Fatalf("failed to count verification sessions: %v", err)
	}

	if sessionCount != 0 {
		t.Fatalf("expected no verification session for already verified email, got %d", sessionCount)
	}
}

func TestResendVerificationEmailReturnsActiveLinkStateForPendingEmail(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	now := time.Now()
	pendingRegistration := models.PendingRegistration{
		Name:                    "Pending User",
		Email:                   "active@example.com",
		PasswordHash:            "hashed-password",
		VerificationTokenHash:   security.HashEmailVerificationToken("active-email-token"),
		VerificationExpiresAt:   now.Add(verificationLifetime),
		VerificationSentAt:      now.Add(-30 * time.Second),
		VerificationDailyCount:  1,
		VerificationWindowStart: now.Add(-1 * time.Hour),
	}

	if err := db.Create(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to create pending registration: %v", err)
	}

	response := performResendVerificationRequest(t, router, "active@example.com")
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if payload["code"] != "VERIFICATION_LINK_ALREADY_SENT" {
		t.Fatalf("expected active link code, got %v", payload["code"])
	}

	if payload["message"] != "A verification link has already been sent to this email address." {
		t.Fatalf("unexpected message %v", payload["message"])
	}

	if payload["retryAfterSeconds"] == nil {
		t.Fatal("expected retryAfterSeconds in active link response")
	}

	if _, hasToken := payload["verificationSessionToken"]; hasToken {
		t.Fatal("active link response should not create a new verification session token")
	}

	var sessionCount int64
	if err := db.Model(&models.VerificationSession{}).Where("email = ?", "active@example.com").Count(&sessionCount).Error; err != nil {
		t.Fatalf("failed to count verification sessions: %v", err)
	}

	if sessionCount != 0 {
		t.Fatalf("expected no new verification session while link is active, got %d", sessionCount)
	}
}

func TestResendVerificationEmailCreatesFreshVerificationSessionForPendingEmail(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	now := time.Now()
	oldRawToken := "old-email-token"
	oldSessionRawToken := "active-session-token"
	pendingRegistration := models.PendingRegistration{
		Name:                    "Pending User",
		Email:                   "pending@example.com",
		PasswordHash:            "hashed-password",
		VerificationTokenHash:   security.HashEmailVerificationToken(oldRawToken),
		VerificationExpiresAt:   now.Add(-time.Second),
		VerificationSentAt:      now.Add(-3 * time.Minute),
		VerificationDailyCount:  1,
		VerificationWindowStart: now.Add(-1 * time.Hour),
	}

	if err := db.Create(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to create pending registration: %v", err)
	}

	activeSession := models.VerificationSession{
		Email:     "pending@example.com",
		TokenHash: security.HashVerificationSessionToken(oldSessionRawToken),
		Status:    models.VerificationSessionStatusPending,
		ExpiresAt: now.Add(verificationLifetime),
	}
	if err := db.Create(&activeSession).Error; err != nil {
		t.Fatalf("failed to create active verification session: %v", err)
	}

	beforeResend := time.Now()
	response := performResendVerificationRequest(t, router, "pending@example.com")
	afterResend := time.Now()
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if payload["message"] != resendVerificationGenericMessage {
		t.Fatalf("unexpected message %q", payload["message"])
	}

	rawSessionToken, ok := payload["verificationSessionToken"].(string)
	if !ok || rawSessionToken == "" {
		t.Fatal("expected resend response to include a raw verification session token")
	}

	if payload["verificationSessionExpiresInSeconds"] != float64(180) {
		t.Fatalf("expected 180 second expiry, got %v", payload["verificationSessionExpiresInSeconds"])
	}

	var updatedPendingRegistration models.PendingRegistration
	if err := db.Where("email = ?", "pending@example.com").First(&updatedPendingRegistration).Error; err != nil {
		t.Fatalf("failed to reload pending registration: %v", err)
	}

	if updatedPendingRegistration.VerificationDailyCount != 2 {
		t.Fatalf("expected verification daily count 2, got %d", updatedPendingRegistration.VerificationDailyCount)
	}

	if updatedPendingRegistration.VerificationTokenHash == security.HashEmailVerificationToken(oldRawToken) {
		t.Fatal("expected verification token hash to be updated")
	}

	if updatedPendingRegistration.VerificationExpiresAt.Before(beforeResend.Add(verificationLifetime)) {
		t.Fatalf("expected fresh resend expiry at least %s after resend start, got %s", verificationLifetime, updatedPendingRegistration.VerificationExpiresAt.Sub(beforeResend))
	}

	if updatedPendingRegistration.VerificationExpiresAt.After(afterResend.Add(verificationLifetime)) {
		t.Fatalf("expected fresh resend expiry no later than %s after resend response, got %s", verificationLifetime, updatedPendingRegistration.VerificationExpiresAt.Sub(afterResend))
	}

	var updatedSession models.VerificationSession
	if err := db.First(&updatedSession, activeSession.ID).Error; err != nil {
		t.Fatalf("failed to reload active verification session: %v", err)
	}

	if updatedSession.Status != models.VerificationSessionStatusExpired {
		t.Fatalf("expected previous verification session to be expired, got %q", updatedSession.Status)
	}

	var freshSession models.VerificationSession
	if err := db.Where("token_hash = ?", security.HashVerificationSessionToken(rawSessionToken)).First(&freshSession).Error; err != nil {
		t.Fatalf("failed to load fresh verification session: %v", err)
	}

	if freshSession.ID == activeSession.ID {
		t.Fatal("expected resend to create a new verification session")
	}

	if freshSession.TokenHash == rawSessionToken {
		t.Fatal("fresh verification session stored the raw token")
	}

	if freshSession.Status != models.VerificationSessionStatusPending {
		t.Fatalf("expected fresh verification session to be pending, got %q", freshSession.Status)
	}

	if !freshSession.ExpiresAt.Equal(updatedPendingRegistration.VerificationExpiresAt) {
		t.Fatal("expected fresh verification session expiry to align with resent email token expiry")
	}

	var pendingSessionCount int64
	if err := db.Model(&models.VerificationSession{}).
		Where("email = ?", "pending@example.com").
		Where("status = ?", models.VerificationSessionStatusPending).
		Count(&pendingSessionCount).Error; err != nil {
		t.Fatalf("failed to count pending sessions: %v", err)
	}

	if pendingSessionCount != 1 {
		t.Fatalf("expected exactly one usable pending session after resend, got %d", pendingSessionCount)
	}

	verifyOldTokenResponse := performVerifyEmailRequest(t, router, oldRawToken)
	if verifyOldTokenResponse.Code != http.StatusBadRequest {
		t.Fatalf("expected old token to fail with status %d, got %d", http.StatusBadRequest, verifyOldTokenResponse.Code)
	}

	oldSessionStatusResponse := performVerificationSessionStatusRequest(t, router, oldSessionRawToken)
	oldSessionPayload := decodeAnyJSONResponse(t, oldSessionStatusResponse)

	if oldSessionStatusResponse.Code != http.StatusOK {
		t.Fatalf("expected old session status %d, got %d", http.StatusOK, oldSessionStatusResponse.Code)
	}

	if oldSessionPayload["status"] != models.VerificationSessionStatusExpired {
		t.Fatalf("expected old session to be expired, got %v", oldSessionPayload["status"])
	}

	if _, hasToken := oldSessionPayload["token"]; hasToken {
		t.Fatal("old verification session issued a JWT")
	}
}

func TestResendVerificationEmailFreshSessionCanAutoLoginAfterPhoneVerification(t *testing.T) {
	_, router, db := setupAuthHandlerTest(t)
	now := time.Now()
	pendingRegistration := models.PendingRegistration{
		Name:                    "Pending User",
		Email:                   "resendlogin@example.com",
		PasswordHash:            "hashed-password",
		VerificationTokenHash:   security.HashEmailVerificationToken("old-email-token"),
		VerificationExpiresAt:   now.Add(-time.Second),
		VerificationSentAt:      now.Add(-3 * time.Minute),
		VerificationDailyCount:  1,
		VerificationWindowStart: now.Add(-1 * time.Hour),
	}

	if err := db.Create(&pendingRegistration).Error; err != nil {
		t.Fatalf("failed to create pending registration: %v", err)
	}

	response := performResendVerificationRequest(t, router, "resendlogin@example.com")
	payload := decodeAnyJSONResponse(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	rawSessionToken := payload["verificationSessionToken"].(string)
	emailToken := "new-email-token"

	if err := db.Model(&models.PendingRegistration{}).
		Where("email = ?", "resendlogin@example.com").
		Updates(map[string]interface{}{
			"verification_token_hash": security.HashEmailVerificationToken(emailToken),
			"verification_expires_at": time.Now().Add(verificationLifetime),
		}).Error; err != nil {
		t.Fatalf("failed to set known email verification token: %v", err)
	}

	verifyResponse := performVerifyEmailRequest(t, router, emailToken)
	if verifyResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, verifyResponse.Code, verifyResponse.Body.String())
	}

	exchangeResponse := performVerificationSessionStatusRequest(t, router, rawSessionToken)
	exchangePayload := decodeAnyJSONResponse(t, exchangeResponse)

	if exchangeResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, exchangeResponse.Code, exchangeResponse.Body.String())
	}

	if exchangePayload["status"] != models.VerificationSessionStatusVerified {
		t.Fatalf("expected verified status, got %v", exchangePayload["status"])
	}

	if exchangePayload["token"] == "" {
		t.Fatal("expected fresh resend verification session to issue a JWT")
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
