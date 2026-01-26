package fixture

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils"
)

// MockJwtValidator is a mock implementation of JwtValidator that returns predefined claims.
type MockJwtValidator struct {
	// userID is the user ID that will be returned in the token claims
	userID string
	// shouldReturnError determines if ValidateToken should return an error
	shouldReturnError bool
}

// NewMockJwtValidator creates a new mock JWT validator with the specified userID.
func NewMockJwtValidator(userID string) *MockJwtValidator {
	return &MockJwtValidator{
		userID:            userID,
		shouldReturnError: false,
	}
}

// NewMockJwtValidatorWithError creates a new mock JWT validator that will return an error.
func NewMockJwtValidatorWithError() *MockJwtValidator {
	return &MockJwtValidator{
		shouldReturnError: true,
	}
}

// SetUserID sets the user ID that will be returned in token claims.
func (m *MockJwtValidator) SetUserID(userID string) {
	m.userID = userID
	m.shouldReturnError = false
}

// SetShouldReturnError sets whether the validator should return an error.
func (m *MockJwtValidator) SetShouldReturnError(shouldReturnError bool) {
	m.shouldReturnError = shouldReturnError
}

// ValidateToken implements the JwtValidator interface.
// It returns token claims containing the userID in the "sub" field.
func (m *MockJwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if m.shouldReturnError {
		return nil, jwt.ErrSignatureInvalid
	}

	// Return claims with userID in "sub" field (as expected by the middleware)
	return jwt.MapClaims{
		"sub": m.userID,
		"iat": 1234567890,
		"exp": 1234567890 + 86400, // 24 hours later
	}, nil
}

// Ensure MockJwtValidator implements jwtUtils.JwtValidator interface
var _ jwtUtils.JwtValidator = (*MockJwtValidator)(nil)
