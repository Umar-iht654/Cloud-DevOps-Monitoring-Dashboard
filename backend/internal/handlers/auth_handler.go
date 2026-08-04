package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/email"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler stores the dependencies needed by the authentication routes.
type AuthHandler struct {
	// DB gives the auth handler access to the PostgreSQL database through GORM.
	DB *gorm.DB

	// JWTSecret is used to sign login tokens.
	JWTSecret string

	// EmailSender sends verification emails.
	EmailSender *email.Sender
}

// RegisterRequest defines the JSON body expected when a user registers.
type RegisterRequest struct {
	// Name stores the user's name from the request body.
	Name string `json:"name"`

	// Email stores the user's email from the request body.
	Email string `json:"email"`

	// Password stores the user's raw password from the request body.
	Password string `json:"password"`
}

// LoginRequest defines the JSON body expected when a user logs in.
type LoginRequest struct {
	// Email stores the user's email from the request body.
	Email string `json:"email"`

	// Password stores the user's raw password from the request body.
	Password string `json:"password"`
}

// ResendVerificationRequest defines the JSON body expected when resending verification emails.
type ResendVerificationRequest struct {
	// Email stores the user's email from the request body.
	Email string `json:"email"`
}

// NewAuthHandler creates a new AuthHandler with database access and a JWT secret.
func NewAuthHandler(db *gorm.DB, jwtSecret string, emailSender *email.Sender) *AuthHandler {
	// This returns a pointer to an AuthHandler so routes can use its methods.
	return &AuthHandler{
		DB:          db,
		JWTSecret:   jwtSecret,
		EmailSender: emailSender,
	}
}

// Register handles new user registration by creating a pending registration first.
func (h *AuthHandler) Register(c *gin.Context) {
	// This creates a variable to store the incoming JSON request body.
	var req RegisterRequest

	// This tries to convert the incoming JSON body into the RegisterRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// This removes extra spaces before or after the user's name.
	req.Name = strings.TrimSpace(req.Name)

	// This lowercases the email and removes extra spaces so duplicate checks are consistent.
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// This removes extra spaces before or after the password.
	req.Password = strings.TrimSpace(req.Password)

	// This checks that all required fields were provided.
	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name, email and password are required",
		})
		return
	}

	// This checks that the password has a minimum length.
	if len(req.Password) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password must be at least 7 characters long",
		})
		return
	}

	// This creates a variable to store any existing real user found with the same email.
	var existingUser models.User

	// This searches the users table for a real verified account with the submitted email.
	err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error

	// If there is no error, a real user with that email already exists.
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email is already registered",
		})
		return
	}

	// This checks whether the database error was something other than "user not found".
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check existing user",
		})
		return
	}

	// This creates a variable to store any existing pending registration for the same email.
	var existingPendingRegistration models.PendingRegistration

	// This searches the pending registrations table for the submitted email.
	pendingErr := h.DB.Where("email = ?", req.Email).First(&existingPendingRegistration).Error

	// If there is no error, a pending registration already exists.
	if pendingErr == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "A verification email has already been sent. Please check your email or request a new verification link.",
		})
		return
	}

	// This checks whether the database error was something other than "pending registration not found".
	if !errors.Is(pendingErr, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check pending registration",
		})
		return
	}

	// This hashes the raw password so the real password is never stored in the database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// This checks whether password hashing failed.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	// This creates a secure email verification token.
	verificationToken, verificationTokenHash, err := security.GenerateEmailVerificationToken()

	// This checks whether token generation failed.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create verification token",
		})
		return
	}

	// This stores the current time so all verification timestamps are consistent.
	now := time.Now()

	// This stores when the verification token should expire.
	verificationExpiresAt := now.Add(1 * time.Hour)

	// This creates a pending registration instead of a real user.
	pendingRegistration := models.PendingRegistration{
		Name:                    req.Name,
		Email:                   req.Email,
		PasswordHash:            string(hashedPassword),
		VerificationTokenHash:   verificationTokenHash,
		VerificationExpiresAt:   verificationExpiresAt,
		VerificationSentAt:      now,
		VerificationDailyCount:  1,
		VerificationWindowStart: now,
	}

	// This inserts the pending registration into the database.
	if err := h.DB.Create(&pendingRegistration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create pending registration",
		})
		return
	}

	// This sends or logs the verification email.
	if err := h.EmailSender.SendVerificationEmail(pendingRegistration.Email, verificationToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Registration started but verification email could not be sent. Please request a new verification email.",
		})
		return
	}

	// This returns a 201 response showing that verification is required before the account is created.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration started. Please verify your email to create your account.",
	})
}

// VerifyEmail verifies a pending registration and creates the real user account.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	// This reads the raw verification token from the URL query string.
	rawToken := strings.TrimSpace(c.Query("token"))

	// This checks that a token was provided.
	if rawToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Verification token is required",
		})
		return
	}

	// This hashes the submitted token so it can be compared with the stored hash.
	tokenHash := security.HashEmailVerificationToken(rawToken)

	// This stores the pending registration found by verification token hash.
	var pendingRegistration models.PendingRegistration

	// This looks up the pending registration with the matching token hash.
	if err := h.DB.Where("verification_token_hash = ?", tokenHash).First(&pendingRegistration).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired verification token",
		})
		return
	}

	// This checks whether the token has expired.
	if time.Now().After(pendingRegistration.VerificationExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired verification token",
		})
		return
	}

	// This creates the real user and deletes the pending registration in one database transaction.
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// This checks whether a real user was already created with this email.
		var existingUser models.User
		if err := tx.Where("email = ?", pendingRegistration.Email).First(&existingUser).Error; err == nil {
			return errors.New("email is already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// This creates the real user after email verification succeeds.
		user := models.User{
			Name:         pendingRegistration.Name,
			Email:        pendingRegistration.Email,
			PasswordHash: pendingRegistration.PasswordHash,
		}

		// This inserts the real user into the users table.
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// This deletes the pending registration so the token cannot be reused.
		if err := tx.Delete(&pendingRegistration).Error; err != nil {
			return err
		}

		// This commits the transaction.
		return nil
	})

	// This checks whether verification failed.
	if err != nil {
		if err.Error() == "email is already registered" {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email is already registered",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to verify email",
		})
		return
	}

	// This returns a success response.
	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully. You can now log in.",
	})
}

// ResendVerificationEmail sends a new verification email if the request is allowed.
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	// This creates a variable to store the incoming JSON request body.
	var req ResendVerificationRequest

	// This tries to convert the incoming JSON body into the request struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// This normalises the email before database lookup.
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// This checks that an email was provided.
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is required",
		})
		return
	}

	// This generic message avoids revealing whether an email exists or is pending.
	genericMessage := "If the account exists and is awaiting verification, a verification email will be sent if allowed."

	// This stores the pending registration found by email.
	var pendingRegistration models.PendingRegistration

	// This searches for a pending registration by email.
	if err := h.DB.Where("email = ?", req.Email).First(&pendingRegistration).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": genericMessage,
		})
		return
	}

	// This stores the current time.
	now := time.Now()

	// This enforces a 10-minute resend cooldown.
	if now.Sub(pendingRegistration.VerificationSentAt) < 10*time.Minute {
		c.JSON(http.StatusOK, gin.H{
			"message": genericMessage,
		})
		return
	}

	// This resets the daily window if more than 24 hours have passed.
	if now.Sub(pendingRegistration.VerificationWindowStart) >= 24*time.Hour {
		pendingRegistration.VerificationWindowStart = now
		pendingRegistration.VerificationDailyCount = 0
	}

	// This enforces a daily email-send limit.
	if pendingRegistration.VerificationDailyCount >= 5 {
		c.JSON(http.StatusOK, gin.H{
			"message": genericMessage,
		})
		return
	}

	// This creates a new secure verification token.
	verificationToken, verificationTokenHash, err := security.GenerateEmailVerificationToken()

	// This checks whether token generation failed.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create verification token",
		})
		return
	}

	// This updates the pending registration with the new token details.
	pendingRegistration.VerificationTokenHash = verificationTokenHash
	pendingRegistration.VerificationExpiresAt = now.Add(1 * time.Hour)
	pendingRegistration.VerificationSentAt = now
	pendingRegistration.VerificationDailyCount++

	// This saves the updated pending registration.
	if err := h.DB.Save(&pendingRegistration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update verification token",
		})
		return
	}

	// This sends or logs the verification email.
	if err := h.EmailSender.SendVerificationEmail(pendingRegistration.Email, verificationToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send verification email",
		})
		return
	}

	// This returns the generic success message.
	c.JSON(http.StatusOK, gin.H{
		"message": genericMessage,
	})
}

// Login handles user login.
func (h *AuthHandler) Login(c *gin.Context) {
	// This creates a variable to store the incoming login JSON body.
	var req LoginRequest

	// This tries to convert the incoming JSON body into the LoginRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		// This returns a 400 response if the request body is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})

		// This stops the function so no more code runs after the error.
		return
	}

	// This lowercases the email and removes extra spaces.
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// This removes extra spaces before or after the password.
	req.Password = strings.TrimSpace(req.Password)

	// This checks that both email and password were provided.
	if req.Email == "" || req.Password == "" {
		// This returns a 400 response if email or password is missing.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email and password are required",
		})

		// This stops the function because login cannot continue without both fields.
		return
	}

	// This creates a variable to store the user found in the database.
	var user models.User

	// This searches the database for a user with the submitted email.
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// This returns a 401 response if no user exists with that email.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})

		// This stops the function so login fails safely.
		return
	}

	// This compares the submitted password with the hashed password stored in the database.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// This returns a 401 response if the password is incorrect.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})

		// This stops the function so no token is created for an invalid login.
		return
	}

	// This creates a new JWT token with user information and an expiry time.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// This signs the JWT token using the secret from the environment variables.
	signedToken, err := token.SignedString([]byte(h.JWTSecret))

	// This checks whether token signing failed.
	if err != nil {
		// This returns a 500 response if the token could not be created.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create login token",
		})

		// This stops the function because login cannot finish without a token.
		return
	}

	// This returns the login token and basic user details to the frontend.
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   signedToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Me returns the currently logged-in user's details.
func (h *AuthHandler) Me(c *gin.Context) {
	// This gets the userID value that was added to the request by the auth middleware.
	userIDValue, exists := c.Get("userID")

	// This checks whether the middleware added a userID to the request.
	if !exists {
		// This returns a 401 response if the user is not authenticated.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User is not authenticated",
		})

		// This stops the function because there is no authenticated user.
		return
	}

	// This converts the userID value from the request context into a uint.
	userID, ok := userIDValue.(uint)

	// This checks whether the userID was stored in the expected format.
	if !ok {
		// This returns a 500 response if the userID in context is invalid.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID in request context",
		})

		// This stops the function because the user cannot be safely loaded.
		return
	}

	// This creates a variable to store the user loaded from the database.
	var user models.User

	// This searches the database for the authenticated user by ID.
	if err := h.DB.First(&user, userID).Error; err != nil {
		// This returns a 404 response if the user cannot be found.
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})

		// This stops the function because there is no user to return.
		return
	}

	// This returns the authenticated user's public details.
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
