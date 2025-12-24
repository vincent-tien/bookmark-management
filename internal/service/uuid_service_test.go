package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-7][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func TestUuidService_Generate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "normal case",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := NewUuid()

			pass, err := svc.Generate()

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, true, uuidRegex.MatchString(pass))
		})
	}
}
