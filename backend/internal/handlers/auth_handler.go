package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Umar-iht654/Cloud-DevOps-Monitoring-Dashboard/backend/internal/models"
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

// NewAuthHandler creates a new AuthHandler with database access and a JWT secret.
func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
	// This returns a pointer to an AuthHandler so routes can use its methods.
	return &AuthHandler{
		DB:        db,
		JWTSecret: jwtSecret,
	}
}

// Register handles new user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	// This creates a variable to store the incoming JSON request body.
	var req RegisterRequest

	// This tries to convert the incoming JSON body into the RegisterRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		// This returns a 400 response if the request body is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})

		// This stops the function so no more code runs after the error.
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
		// This returns a 400 response if any required field is missing.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name, email and password are required",
		})

		// This stops the function so an invalid user is not created.
		return
	}

	// This checks that the password has a minimum length.
	if len(req.Password) < 7 {
		// This returns a 400 response if the password is too short.
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password must be at least 7 characters long",
		})

		// This stops the function so a weak password is not saved.
		return
	}

	// This creates a variable to store any existing user found with the same email.
	var existingUser models.User

	// This searches the database for a user with the submitted email.
	err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error

	// If there is no error, a user with that email already exists.
	if err == nil {
		// This returns a 409 response because the email is already taken.
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email is already registered",
		})

		// This stops the function so duplicate users are not created.
		return
	}

	// This checks whether the database error was something other than "user not found".
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// This returns a 500 response because the database check failed unexpectedly.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check existing user",
		})

		// This stops the function because it is unsafe to continue after a database error.
		return
	}

	// This hashes the raw password so the real password is never stored in the database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// This checks whether password hashing failed.
	if err != nil {
		// This returns a 500 response if the password could not be hashed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
		})

		// This stops the function because the user should not be created without a hashed password.
		return
	}

	// This creates a new User model ready to be saved in the database.
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// This inserts the new user into the users table.
	if err := h.DB.Create(&user).Error; err != nil {
		// This returns a 500 response if the user could not be saved.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})

		// This stops the function because registration failed.
		return
	}

	// This returns a 201 response showing that the user was created successfully.
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
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
