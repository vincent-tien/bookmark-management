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
	errorsPkg "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/middleware"
	"github.com/vincent-tien/bookmark-management/internal/model"
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
				//nolint:gosec // NOSONAR - This is test data, not a real credential
				// This intentionally uses a short password to test validation
				testValue := "short"
				reqBody.Password = testValue // NOSONAR - test data for validation testing
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

func TestUser_Login(t *testing.T) {
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
				setupJSONRequest(ctx, http.MethodPost, getLoginEndpoint(), validLoginRequest())
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				expectedReq := validLoginRequest()
				expectedToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test.token"
				mockSvc.On("Login", ctx, expectedReq).Return(expectedToken, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `"message":"Logged in successfully!"`,
		},
		{
			name: "bad request - invalid JSON",
			setupRequest: func(ctx *gin.Context) {
				setupInvalidJSONRequest(ctx, http.MethodPost, getLoginEndpoint())
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
				setupJSONRequest(ctx, http.MethodPost, getLoginEndpoint(), map[string]interface{}{})
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
				reqBody := validLoginRequest()
				//nolint:gosec // NOSONAR - This is test data, not a real credential
				// This intentionally uses a short password to test validation
				testValue := "short"
				reqBody.RawPassword = testValue // NOSONAR - test data for validation testing
				setupJSONRequest(ctx, http.MethodPost, getLoginEndpoint(), reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "bad request - invalid auth",
			setupRequest: func(ctx *gin.Context) {
				setupJSONRequest(ctx, http.MethodPost, getLoginEndpoint(), validLoginRequest())
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				expectedReq := validLoginRequest()
				mockSvc.On("Login", ctx, expectedReq).Return("", errorsPkg.ErrInvalidAuth)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   `"error":"invalid username or password"`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				setupJSONRequest(ctx, http.MethodPost, getLoginEndpoint(), validLoginRequest())
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				expectedReq := validLoginRequest()
				mockSvc.On("Login", ctx, expectedReq).Return("", errors.New("database error"))
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
			handler.Login(ctx)

			assertResponse(t, rec, tc.expectedStatus, tc.expectedResp)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUser_GetProfile(t *testing.T) {
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
				setupGetRequest(ctx, http.MethodGet, getProfileEndpoint())
				setupUserIDInContext(ctx, "test-user-id-123")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				userId := "test-user-id-123"
				now := time.Now()
				expectedUser := &model.User{
					ID:          userId,
					Username:    "johndoe",
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				mockSvc.On("GetProfile", ctx, userId).Return(expectedUser, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `"id":"test-user-id-123"`,
		},
		{
			name: "unauthorized - missing user id in context",
			setupRequest: func(ctx *gin.Context) {
				setupGetRequest(ctx, http.MethodGet, getProfileEndpoint())
				// Don't set userId in context
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp:   `"error":"Invalid Token"`,
		},
		{
			name: "unauthorized - invalid user id type in context",
			setupRequest: func(ctx *gin.Context) {
				setupGetRequest(ctx, http.MethodGet, getProfileEndpoint())
				// Set userId with wrong type (int instead of string)
				ctx.Set(middleware.UserIDKey, 123)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp:   `"error":"Invalid Token"`,
		},
		{
			name: "unauthorized - user not found",
			setupRequest: func(ctx *gin.Context) {
				setupGetRequest(ctx, http.MethodGet, getProfileEndpoint())
				setupUserIDInContext(ctx, "non-existent-user-id")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				userId := "non-existent-user-id"
				mockSvc.On("GetProfile", ctx, userId).Return(nil, errors.New("user not found"))
				return mockSvc
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp:   `"error":"Invalid Token"`,
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				setupGetRequest(ctx, http.MethodGet, getProfileEndpoint())
				setupUserIDInContext(ctx, "test-user-id-123")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.User {
				mockSvc := mocks.NewUser(t)
				userId := "test-user-id-123"
				mockSvc.On("GetProfile", ctx, userId).Return(nil, errors.New("database error"))
				return mockSvc
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp:   `"error":"Invalid Token"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec, ctx := createTestContext()
			tc.setupRequest(ctx)
			mockSvc := tc.setupMockSvc(t, ctx)
			handler := NewUserHandler(mockSvc)
			handler.GetProfile(ctx)

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
		Password:    "Password123!", //nolint:gosec // NOSONAR - test data, not a real credential
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

func getLoginEndpoint() string {
	return fmt.Sprintf("/v1/%s", routers.Endpoints.AuthLogin)
}

func getProfileEndpoint() string {
	return fmt.Sprintf("/v1/%s", routers.Endpoints.GetProfile)
}

// setupGetRequest creates a GET request and sets it on the gin context
func setupGetRequest(ctx *gin.Context, method, endpoint string) {
	ctx.Request = httptest.NewRequest(method, endpoint, nil)
}

// setupUserIDInContext sets the user ID in the gin context
func setupUserIDInContext(ctx *gin.Context, userID string) {
	ctx.Set(middleware.UserIDKey, userID)
}

// validLoginRequest returns a valid LoginRequestDto for testing
func validLoginRequest() dto.LoginRequestDto {
	return dto.LoginRequestDto{
		Username:    "johndoe",
		RawPassword: "Password123!", //nolint:gosec // NOSONAR - test data, not a real credential
	}
}
