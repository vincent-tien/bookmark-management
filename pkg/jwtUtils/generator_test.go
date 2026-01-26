package jwtUtils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// TestJwtGenerator_GenerateToken tests the GenerateToken method of the jwtGenerator
// struct. It uses a list of test cases to validate the functionality of
// the method, including the validation of the key path, the input content,
// and the expected output and error.
func TestJwtGenerator_GenerateToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		keyPath       string
		inputContent  jwt.MapClaims
		expectOutput  string
		expectedError error
	}{
		{
			name:    "validate key path",
			keyPath: filepath.FromSlash("./private.test.pem"),
			inputContent: jwt.MapClaims{
				"id":   1234,
				"name": "John",
			},
			expectedError: nil,
			expectOutput:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNCwibmFtZSI6IkpvaG4ifQ.B0lYzj5pZWEnBn2aETGTtQdSSQpODGB1NtxJH2TLe9R3vnHT8RV0ZhV-GKBC3A1eGGsgvRGCmNk1Kds6f5rIUk3dVVcaabI38p6tEmxEpwWXmJ8Rid_UPlXx-0XdL9gKXTaDQ1Hjn3MzbzWfzb-t8brxauh5SoJxqnHoYkj5BMP3Crflu51wlRHddIkRooXKVxubinkrmeuZxdCf6oX09HXasuXrR2AVp0GZi6wL0ACQC-_NrCZRMRNkZV7ap70lmETTnlS5HpCShqkAHmAy49LQko7LRpWcFPft0VX-dTJZFOivlhTtXUfCvn99GzNwKE5fND1zcTNz6yXUEopV_g",
		},
		{
			name:          "invalid key path",
			keyPath:       filepath.FromSlash("./nonexistent.pem"),
			inputContent:  nil,
			expectOutput:  "",
			expectedError: os.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testGen, err := NewJwtGenerator(tc.keyPath)
			if tc.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tc.expectedError, os.ErrNotExist) {
					assert.True(t, errors.Is(err, os.ErrNotExist))
				} else {
					assert.Equal(t, tc.expectedError, err)
				}
				assert.Nil(t, testGen)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, testGen)

			res, err := testGen.GenerateToken(tc.inputContent)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectOutput, res)
		})
	}
}
