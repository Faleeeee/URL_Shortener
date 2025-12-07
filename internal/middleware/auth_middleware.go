package middleware

import (
	"net/http"
	"strings"

	"github.com/Faleeeee/URL_Shortener/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, "Authorization header required", "AUTH_REQUIRED", "Missing Authorization header")
			c.Abort()
			return
		}

		var token string

		// Check if it starts with "Bearer "
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Extract token after "Bearer "
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Assume the entire header is the token (auto-add Bearer support)
			token = authHeader
		}

		// Validate token is not empty after extraction
		if token == "" {
			utils.SendError(c, http.StatusUnauthorized, "Token is empty", "INVALID_TOKEN", "Authorization token is empty")
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			if err == utils.ErrExpiredToken {
				utils.SendError(c, http.StatusUnauthorized, "Token has expired", "TOKEN_EXPIRED", "The provided token has expired")
			} else {
				utils.SendError(c, http.StatusUnauthorized, "Invalid token", "INVALID_TOKEN", "The provided token is invalid")
			}
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
