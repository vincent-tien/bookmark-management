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
				setupJSONRequest(ctx, http.MethodPost, getRegisterEndpoint(), validRegisterRequest())
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				now := time.Now().Format(time.RFC3339)
				expectedReq := validRegisterRequest()
				expectedResp := dto.RegisterResponseDto{
					ID:          "test-uuid-123",
					Username:    expectedReq.Username,
					DisplayName: expectedReq.DisplayName,
					Email:       expectedReq.Email,
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
				setupInvalidJSONRequest(ctx, http.MethodPost, getRegisterEndpoint())
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
				setupJSONRequest(ctx, http.MethodPost, getRegisterEndpoint(), map[string]interface{}{})
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
				reqBody := validRegisterRequest()
				reqBody.Email = "invalid-email"
				setupJSONRequest(ctx, http.MethodPost, getRegisterEndpoint(), reqBody)
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
				reqBody := validRegisterRequest()
				reqBody.Password = "short"
				setupJSONRequest(ctx, http.MethodPost, getRegisterEndpoint(), reqBody)
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
				setupJSONRequest(ctx, http.MethodPost, getRegisterEndpoint(), validRegisterRequest())
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				expectedReq := validRegisterRequest()
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
			rec, ctx := createTestContext()
			tc.setupRequest(ctx)
			mockSvc := tc.setupMockSvc(t, ctx)
			handler := NewUserHandler(mockSvc)
			handler.Register(ctx)

			assertResponse(t, rec, tc.expectedStatus, tc.expectedResp)
			mockSvc.AssertExpectations(t)
		})
	}
}

// setupJSONRequest creates a JSON request and sets it on the gin context
func setupJSONRequest(ctx *gin.Context, method, endpoint string, body interface{}) {
	var bodyReader *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonData)
	} else {
		bodyReader = bytes.NewBuffer([]byte{})
	}
	ctx.Request = httptest.NewRequest(method, endpoint, bodyReader)
	ctx.Request.Header.Set("Content-Type", "application/json")
}

// setupInvalidJSONRequest creates a request with invalid JSON
func setupInvalidJSONRequest(ctx *gin.Context, method, endpoint string) {
	ctx.Request = httptest.NewRequest(method, endpoint, strings.NewReader("invalid json"))
	ctx.Request.Header.Set("Content-Type", "application/json")
}

// validRegisterRequest returns a valid RegisterRequestDto for testing
func validRegisterRequest() dto.RegisterRequestDto {
	return dto.RegisterRequestDto{
		DisplayName: "John Doe",
		Username:    "johndoe",
		Password:    "Password123!",
		Email:       "john.doe@example.com",
	}
}

// createTestContext creates a gin test context with a recorder
func createTestContext() (*httptest.ResponseRecorder, *gin.Context) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	return rec, ctx
}

// assertResponse asserts the response based on the expected status code
func assertResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedResp string) {
	t.Helper()
	assert.Equal(t, expectedStatus, rec.Code)
	actualBody := rec.Body.String()

	switch {
	case expectedResp != "":
		// For success case, check that response contains the expected message
		assert.Contains(t, actualBody, expectedResp)
		if expectedStatus == http.StatusOK {
			assert.Contains(t, actualBody, "data")
		}
	case expectedStatus == http.StatusBadRequest:
		// For bad request, verify it contains message field
		assert.Contains(t, actualBody, "message")
		// For validation errors, it should contain details
		if strings.Contains(actualBody, "Invalid request") {
			assert.Contains(t, actualBody, "details")
		}
	case expectedStatus == http.StatusInternalServerError:
		// For internal server error, verify it contains error message
		assert.Contains(t, actualBody, "message")
		assert.Contains(t, actualBody, "Something went wrong")
	}
}

func getRegisterEndpoint() string {
	return fmt.Sprintf("/v1/%s", routers.Endpoints.UserRegister)
}
