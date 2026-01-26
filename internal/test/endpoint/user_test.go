package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/test/fixture"
	"gorm.io/gorm"
)

func TestUserRegisterEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(api apipkg.Engine) *httptest.ResponseRecorder
		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "test@example.com",
					Password:    fixture.ValidTestPassword(),
				}
				return executeJSONRequest(api, http.MethodPost, getUserRegisterEndpoint(), reqBody)
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data    dto.RegisterResponseDto `json:"data"`
					Message string                  `json:"message"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Register an user successfully!", resp.Message)
				assert.NotEmpty(t, resp.Data.ID)
				assert.Equal(t, "testuser", resp.Data.Username)
				assert.Equal(t, "test@example.com", resp.Data.Email)
				assert.Equal(t, "Test User", resp.Data.DisplayName)
			},
		},
		{
			name: "bad request - invalid JSON",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				return executeRequest(api, http.MethodPost, getUserRegisterEndpoint(), "invalid json")
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				// Invalid JSON returns InternalErrorResponse, not validation error
				assert.Equal(t, "Something went wrong", resp.Message)
				// Details can be nil or empty for invalid JSON
			},
		},
		{
			name: "bad request - missing required fields",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{}
				return executeJSONRequest(api, http.MethodPost, getUserRegisterEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Invalid request", resp.Message)
				assert.NotEmpty(t, resp.Details)
			},
		},
		{
			name: "bad request - invalid email",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "invalid-email",
					Password:    fixture.ValidTestPassword(),
				}
				return executeJSONRequest(api, http.MethodPost, getUserRegisterEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Invalid request", resp.Message)
				assert.NotEmpty(t, resp.Details)
				assert.Contains(t, rec.Body.String(), "Email is invalid email")
			},
		},
		{
			name: "bad request - weak password",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.RegisterRequestDto{
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "test@example.com",
					Password:    "weak",
				}
				return executeJSONRequest(api, http.MethodPost, getUserRegisterEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Invalid request", resp.Message)
				assert.NotEmpty(t, resp.Details)
			},
		},
	}

	cfg := defaultTestConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setup := setupTestInfrastructure(t, cfg, false)
			rec := tc.setupTestHttp(setup.app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}

func TestUserLoginEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder
		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				testUser := createTestUserWithDefaults(t, db)
				reqBody := dto.LoginRequestDto{
					Username:    testUser.Username,
					RawPassword: fixture.ValidTestPassword(),
				}
				return executeJSONRequest(api, http.MethodPost, getUserLoginEndpoint(), reqBody)
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data    string `json:"data"`
					Message string `json:"message"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Logged in successfully!", resp.Message)
				assert.NotEmpty(t, resp.Data)
			},
		},
		{
			name: "bad request - invalid JSON",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				return executeRequest(api, http.MethodPost, getUserLoginEndpoint(), "invalid json")
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				// Invalid JSON returns InternalErrorResponse, not validation error
				assert.Equal(t, "Something went wrong", resp.Message)
				// Details can be nil or empty for invalid JSON
			},
		},
		{
			name: "bad request - missing required fields",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{}
				return executeJSONRequest(api, http.MethodPost, getUserLoginEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Invalid request", resp.Message)
				assert.NotEmpty(t, resp.Details)
			},
		},
		{
			name: "bad request - invalid credentials - wrong password",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				testUser := createTestUserWithDefaults(t, db)
				reqBody := dto.LoginRequestDto{
					Username:    testUser.Username,
					RawPassword: fixture.WrongTestPassword(),
				}
				return executeJSONRequest(api, http.MethodPost, getUserLoginEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "invalid")
			},
		},
		{
			name: "bad request - invalid credentials - user not found",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				reqBody := dto.LoginRequestDto{
					Username:    "nonexistent",
					RawPassword: fixture.ValidTestPassword(),
				}
				return executeJSONRequest(api, http.MethodPost, getUserLoginEndpoint(), reqBody)
			},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// Should return internal server error when user not found
				assert.Contains(t, rec.Body.String(), "message")
			},
		},
		{
			name: "bad request - password too short",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB) *httptest.ResponseRecorder {
				reqBody := dto.LoginRequestDto{
					Username:    "testuser",
					RawPassword: "short",
				}
				return executeJSONRequest(api, http.MethodPost, getUserLoginEndpoint(), reqBody)
			},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string   `json:"message"`
					Details []string `json:"details"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Invalid request", resp.Message)
				assert.NotEmpty(t, resp.Details)
			},
		},
	}

	cfg := defaultTestConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setup := setupTestInfrastructure(t, cfg, false)
			rec := tc.setupTestHttp(t, setup.app, setup.mockDB)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}

func TestUserGetProfileEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder
		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				testUser := createTestUserWithDefaults(t, db)
				mockJwtValidator.SetUserID(testUser.ID)
				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				req.Header.Set("Authorization", "Bearer mock.token.from.fixture")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data    dto.UserProfileResponseDto `json:"data"`
					Message string                     `json:"message"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.NotEmpty(t, resp.Data.UserId)
				assert.Equal(t, "testuser", resp.Data.Username)
				assert.Equal(t, "test@example.com", resp.Data.Email)
				assert.Equal(t, "Test User", resp.Data.DisplayName)
				assert.NotEmpty(t, resp.Data.CreatedAt)
				assert.NotEmpty(t, resp.Data.UpdatedAt)
			},
		},
		{
			name: "unauthorized - missing authorization header",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "Authorization is required")
			},
		},
		{
			name: "unauthorized - invalid token format",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				req.Header.Set("Authorization", "InvalidFormat token")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "Bearer token")
			},
		},
		{
			name: "unauthorized - invalid token",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				// Set mock validator to return error
				mockJwtValidator.SetShouldReturnError(true)

				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "invalid token")
			},
		},
		{
			name: "unauthorized - empty token",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				req.Header.Set("Authorization", "Bearer ")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "required")
			},
		},
		{
			name: "unauthorized - user not found",
			setupTestHttp: func(t *testing.T, api apipkg.Engine, db *gorm.DB, mockJwtValidator *fixture.MockJwtValidator) *httptest.ResponseRecorder {
				// Set userID in mock validator for non-existent user
				nonExistentUserID := "00000000-0000-0000-0000-000000000000"
				mockJwtValidator.SetUserID(nonExistentUserID)

				req := httptest.NewRequest(http.MethodGet, getUserProfileEndpoint(), nil)
				req.Header.Set("Authorization", "Bearer mock.token.from.fixture")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "Invalid Token")
			},
		},
	}

	cfg := defaultTestConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setup := setupTestInfrastructure(t, cfg, true)
			rec := tc.setupTestHttp(t, setup.app, setup.mockDB, setup.mockJwtValidator)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}
