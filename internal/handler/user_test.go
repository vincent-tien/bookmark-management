package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service/mocks"
	validationPkg "github.com/vincent-tien/bookmark-management/pkg/validation"
)

func init() {
	// Register custom validators for tests
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validationPkg.RegisterCustomValidators(v); err != nil {
			panic(fmt.Sprintf("Failed to register custom validators: %v", err))
		}
	}
}

func TestUser_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.User
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success case",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "Password123!",
					Email:       "john.doe@example.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				now := time.Now().Format(time.RFC3339)
				expectedReq := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "Password123!",
					Email:       "john.doe@example.com",
				}
				expectedResp := dto.RegisterResponseDto{
					ID:          "test-uuid-123",
					Username:    "johndoe",
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				mockSvc.On("Register", ctx, expectedReq).Return(expectedResp, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `"message":"Register an user successfully!"`,
		},
		{
			name: "bad request - invalid JSON",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), strings.NewReader("invalid json"))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "bad request - missing required fields",
			setupRequest: func(ctx *gin.Context) {
				reqBody := map[string]interface{}{}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "bad request - invalid email format",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "Password123!",
					Email:       "invalid-email",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "bad request - password too short",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "short",
					Email:       "john.doe@example.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "Password123!",
					Email:       "john.doe@example.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getRegisterEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				expectedReq := dto.RegisterRequestDto{
					DisplayName: "John Doe",
					Username:    "johndoe",
					Password:    "Password123!",
					Email:       "john.doe@example.com",
				}
				mockSvc.On("Register", ctx, expectedReq).Return(dto.RegisterResponseDto{}, errors.New("database error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)
			mockSvc := tc.setupMockSvc(t, ctx)
			handler := NewUserHandler(mockSvc)
			handler.Register(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedResp != "" {
				// For success case, check that response contains the expected message
				actualBody := rec.Body.String()
				assert.Contains(t, actualBody, tc.expectedResp)
				// Also verify that data field exists in success response
				if tc.expectedStatus == http.StatusOK {
					assert.Contains(t, actualBody, "data")
				}
			} else if tc.expectedStatus == http.StatusBadRequest {
				// For bad request, verify it contains message field
				actualBody := rec.Body.String()
				assert.Contains(t, actualBody, "message")
				// For validation errors, it should contain details
				if strings.Contains(actualBody, "Invalid request") {
					assert.Contains(t, actualBody, "details")
				}
			} else if tc.expectedStatus == http.StatusInternalServerError {
				// For internal server error, verify it contains error message
				actualBody := rec.Body.String()
				assert.Contains(t, actualBody, "message")
				assert.Contains(t, actualBody, "Something went wrong")
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func getRegisterEndpoint() string {
	return fmt.Sprintf("/v1/%s", routers.Endpoints.UserRegister)
}
