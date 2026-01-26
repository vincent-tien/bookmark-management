package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/internal/middleware"
)

// GetUserIDFromContext extracts the user ID from the JWT middleware context.
// Returns the user ID as a string and a boolean indicating success.
// If the user ID is not found or is not a string, it returns an empty string and false.
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userId, ok := c.Get(middleware.UserIDKey)
	if !ok {
		return "", false
	}

	userIdValue, ok := userId.(string)
	if !ok {
		return "", false
	}

	return userIdValue, true
}
