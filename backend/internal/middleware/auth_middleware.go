package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes by checking for a valid JWT token.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	// This returns a Gin middleware function.
	return func(c *gin.Context) {
		// This reads the Authorization header from the incoming request.
		authHeader := c.GetHeader("Authorization")

		// This checks whether the Authorization header is missing.
		if authHeader == "" {
			// This returns a 401 response because protected routes require a token.
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header is required",
			})

			// This stops Gin from running the next handler.
			c.Abort()

			// This stops the middleware function.
			return
		}

		// This checks that the Authorization header starts with "Bearer ".
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// This returns a 401 response because the token format is invalid.
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header must start with Bearer",
			})

			// This stops Gin from running the next handler.
			c.Abort()

			// This stops the middleware function.
			return
		}

		// This removes "Bearer " from the header so only the raw token remains.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// This creates an empty map where the JWT claims will be stored after parsing.
		claims := jwt.MapClaims{}

		// This parses and validates the token using the expected signing method and secret.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// This checks that the token was signed using an HMAC signing method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// This returns an error if the token uses an unexpected signing method.
				return nil, fmt.Errorf("unexpected signing method")
			}

			// This returns the JWT secret so the library can verify the token signature.
			return []byte(jwtSecret), nil
		})

		// This checks whether the token failed to parse or is invalid.
		if err != nil || !token.Valid {
			// This returns a 401 response if the token is invalid or expired.
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
			})

			// This stops Gin from running the next handler.
			c.Abort()

			// This stops the middleware function.
			return
		}

		// This extracts the user_id claim from the token as a float64.
		userIDFloat, ok := claims["user_id"].(float64)

		// This checks whether the user_id claim exists and is in the expected format.
		if !ok {
			// This returns a 401 response if the token does not contain a valid user_id.
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token payload",
			})

			// This stops Gin from running the next handler.
			c.Abort()

			// This stops the middleware function.
			return
		}

		// This stores the authenticated user's ID in the request context for handlers to use.
		c.Set("userID", uint(userIDFloat))

		// This allows the request to continue to the protected route handler.
		c.Next()
	}
}
