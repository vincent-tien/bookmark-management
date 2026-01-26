package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils"
)

const bearerPrefix = "Bearer "

// UserIDKey is the Gin context key under which the authenticated user's ID (from JWT "sub" claim) is stored.
const UserIDKey = "userId"

type JwtAuth interface {
	JwtAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtValidator jwtUtils.JwtValidator
}

// NewJwtAuth returns a new jwtAuth middleware that uses the given JWT validator.
func NewJwtAuth(validator jwtUtils.JwtValidator) JwtAuth {
	return &jwtAuth{jwtValidator: validator}
}

// JwtAuth returns a Gin middleware function that validates JWT tokens and
// stores the user_id from the token content to the Gin context.
//
// It expects the JWT token to be passed in the Authorization header.
// If the header is empty, it aborts the request with a 401 Unauthorized status.
//
// After validating the token, it verifies the token content and stores the user_id
// to the Gin context.
//
// It can be used to protect routes that require authentication.
func (j *jwtAuth) JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		claims, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token content"})
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
