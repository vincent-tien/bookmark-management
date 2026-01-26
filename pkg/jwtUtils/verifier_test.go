package jwtUtils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// TestJwtValidator_ValidateToken tests the ValidateToken method of the jwtValidator
// struct. It uses a list of test cases to validate the functionality of
// the method, including the validation of valid tokens, invalid tokens,
// and the expected output and error.
func TestJwtValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	// First, generate a valid token for testing
	gen, err := NewJwtGenerator(filepath.FromSlash("./private.test.pem"))
	if err != nil {
		t.Fatalf("Failed to create JWT generator: %v", err)
	}

	validClaims := jwt.MapClaims{
		"id":   "1234",
		"name": "John",
	}
	validToken, err := gen.GenerateToken(validClaims)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	testCases := []struct {
		name           string
		publicKeyPath  string
		tokenString    string
		expectedClaims jwt.MapClaims
		expectedError  error
	}{
		{
			name:           "valid token",
			publicKeyPath:  filepath.FromSlash("./public.test.pem"),
			tokenString:    validToken,
			expectedClaims: validClaims,
			expectedError:  nil,
		},
		{
			name:           "invalid token - malformed",
			publicKeyPath:  filepath.FromSlash("./public.test.pem"),
			tokenString:    "invalid.token.string",
			expectedClaims: nil,
			expectedError:  errInvalidToken,
		},
		{
			name:           "invalid token - empty string",
			publicKeyPath:  filepath.FromSlash("./public.test.pem"),
			tokenString:    "",
			expectedClaims: nil,
			expectedError:  errInvalidToken,
		},
		{
			name:           "invalid token - wrong signature",
			publicKeyPath:  filepath.FromSlash("./public.test.pem"),
			tokenString:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNCwibmFtZSI6IkpvaG4ifQ.invalid_signature",
			expectedClaims: nil,
			expectedError:  errInvalidToken,
		},
		{
			name:           "invalid public key path",
			publicKeyPath:  filepath.FromSlash("./nonexistent.pem"),
			tokenString:    validToken,
			expectedClaims: nil,
			expectedError:  os.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := NewJwtValidator(tc.publicKeyPath)
			if tc.expectedError != nil && errors.Is(tc.expectedError, os.ErrNotExist) {
				// Handle initialization error for invalid key path
				assert.Error(t, err)
				assert.True(t, errors.Is(err, os.ErrNotExist))
				assert.Nil(t, validator)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, validator)

			// Type assert to access ValidateToken method
			jwtValidator, ok := validator.(*jwtValidator)
			assert.True(t, ok, "validator should be of type *jwtValidator")

			claims, err := jwtValidator.ValidateToken(tc.tokenString)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tc.expectedClaims != nil {
					assert.Equal(t, tc.expectedClaims["id"], claims["id"])
					assert.Equal(t, tc.expectedClaims["name"], claims["name"])
				}
			}
		})
	}
}
