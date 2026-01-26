package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/repository/mocks"
	jwtUtilsMocks "github.com/vincent-tien/bookmark-management/pkg/jwtUtils/mocks"
	"github.com/vincent-tien/bookmark-management/pkg/utils"
)

// validateTestResult is a helper function to validate test results and errors.
// It handles the common pattern of either calling a custom validateResult function
// or asserting on expectedError.
func validateTestResult[T any](t *testing.T, result T, err error, expectedError error, validateResult func(*testing.T, T, error)) {
	t.Helper()
	if validateResult != nil {
		validateResult(t, result, err)
	} else {
		if expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, expectedError, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUser_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMockRepo  func(t *testing.T) *mocks.User
		request        dto.RegisterRequestDto
		expectedError  error
		validateResult func(t *testing.T, resp dto.RegisterResponseDto, err error)
	}{
		{
			name: "success",
			setupMockRepo: func(t *testing.T) *mocks.User {
				mockRepo := mocks.NewUser(t)
				// Mock CreateUser to succeed - simulate GORM's BeforeCreate hook
				mockRepo.On("CreateUser", t.Context(), mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
					u := args.Get(1).(*model.User)
					// Simulate GORM's BeforeCreate hook - generate UUID if ID is empty
					if u.ID == "" {
						userID, err := uuid.NewV7()
						if err == nil {
							u.ID = userID.String()
						}
					}
					assert.Equal(t, "johndoe", u.Username)
					assert.Equal(t, "John Doe", u.DisplayName)
					assert.Equal(t, "john.doe@example.com", u.Email)
					assert.NotEmpty(t, u.ID)
					assert.True(t, utils.VerifyPassword("Password123!", u.Password))
				}).Return(func(ctx context.Context, u *model.User) *model.User {
					// Return the same user model that was passed in
					return u
				}, nil)

				return mockRepo
			},
			request: dto.RegisterRequestDto{
				DisplayName: "John Doe",
				Username:    "johndoe",
				Password:    "Password123!",
				Email:       "john.doe@example.com",
			},
			expectedError: nil,
			validateResult: func(t *testing.T, resp dto.RegisterResponseDto, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "johndoe", resp.Username)
				assert.Equal(t, "John Doe", resp.DisplayName)
				assert.Equal(t, "john.doe@example.com", resp.Email)
				assert.NotEmpty(t, resp.CreatedAt)
				assert.NotEmpty(t, resp.UpdatedAt)
				// Validate UUID format
				_, parseErr := uuid.Parse(resp.ID)
				assert.NoError(t, parseErr)
			},
		},
		{
			name: "repository error",
			setupMockRepo: func(t *testing.T) *mocks.User {
				mockRepo := mocks.NewUser(t)
				// Mock CreateUser to return an error
				mockRepo.On("CreateUser", t.Context(), mock.AnythingOfType("*model.User")).Return(nil, assert.AnError)

				return mockRepo
			},
			request: dto.RegisterRequestDto{
				DisplayName: "John Doe",
				Username:    "johndoe",
				Password:    "Password123!",
				Email:       "john.doe@example.com",
			},
			expectedError:  assert.AnError,
			validateResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := tc.setupMockRepo(t)
			ctx := t.Context()
			mockJwtGen := jwtUtilsMocks.NewJwtGenerator(t)
			service := NewUserService(mockRepo, mockJwtGen)

			resp, err := service.Register(ctx, tc.request)
			validateTestResult(t, resp, err, tc.expectedError, tc.validateResult)
		})
	}
}
