package fixture

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils/mocks"
)

// NewMockJwtValidator creates a new mock JWT validator using mockery.
// The mock is registered with the test for automatic assertion cleanup.
func NewMockJwtValidator(t *testing.T) *mocks.JwtValidator {
	return mocks.NewJwtValidator(t)
}

// SetupMockJwtValidatorWithUserID configures the mock to return claims with the given userID.
// This should be called before making the HTTP request that triggers token validation.
func SetupMockJwtValidatorWithUserID(m *mocks.JwtValidator, userID string) {
	m.On("ValidateToken", mock.Anything).Return(jwt.MapClaims{
		"sub": userID,
		"iat": float64(1234567890),
		"exp": float64(1234567890 + 86400), // 24 hours later
	}, nil).Once()
}

// SetupMockJwtValidatorWithError configures the mock to return an error.
// This should be called before making the HTTP request that triggers token validation.
func SetupMockJwtValidatorWithError(m *mocks.JwtValidator) {
	m.On("ValidateToken", mock.Anything).Return(nil, jwt.ErrSignatureInvalid).Once()
}
