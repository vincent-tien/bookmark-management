package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/test/fixture"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
	sqldbPkg "github.com/vincent-tien/bookmark-management/pkg/sqldb"
	"github.com/vincent-tien/bookmark-management/pkg/utils"
	"gorm.io/gorm"
)

// getProjectRoot finds the project root by looking for go.mod file
func getProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root (go.mod not found)")
		}
		dir = parent
	}
}

// testSetup contains common test infrastructure
type testSetup struct {
	mockRedis        *redis.Client
	mockDB           *gorm.DB
	jwtGenerator     jwtUtils.JwtGenerator
	jwtValidator     jwtUtils.JwtValidator
	mockJwtValidator *fixture.MockJwtValidator
	app              apipkg.Engine
}

// setupTestInfrastructure sets up common test infrastructure (Redis, DB, JWT, etc.)
func setupTestInfrastructure(t *testing.T, cfg *config.Config, useMockValidator bool) *testSetup {
	t.Helper()

	mockRedis := redisPkg.InitMockRedis(t)
	mockDB := sqldbPkg.InitMockDb(t)

	// Migrate user table
	require.NoError(t, mockDB.AutoMigrate(&model.User{}))

	projectRoot := getProjectRoot(t)
	privateKeyPath := filepath.Join(projectRoot, "pkg", "jwtUtils", "private.test.pem")

	jwtGenerator, err := jwtUtils.NewJwtGenerator(privateKeyPath)
	if err != nil {
		t.Fatalf("Failed to create JWT generator: %v", err)
	}

	var jwtValidator jwtUtils.JwtValidator
	var mockJwtValidator *fixture.MockJwtValidator

	if useMockValidator {
		mockJwtValidator = fixture.NewMockJwtValidator("")
		jwtValidator = mockJwtValidator
	} else {
		publicKeyPath := filepath.Join(projectRoot, "pkg", "jwtUtils", "public.test.pem")
		realValidator, err := jwtUtils.NewJwtValidator(publicKeyPath)
		if err != nil {
			t.Fatalf("Failed to create JWT validator: %v", err)
		}
		jwtValidator = realValidator
	}

	app := apipkg.New(cfg, mockRedis, mockDB, jwtGenerator, jwtValidator)

	return &testSetup{
		mockRedis:        mockRedis,
		mockDB:           mockDB,
		jwtGenerator:     jwtGenerator,
		jwtValidator:     jwtValidator,
		mockJwtValidator: mockJwtValidator,
		app:              app,
	}
}

// setupTestInfrastructureSimple sets up test infrastructure without user migration (for non-user tests)
func setupTestInfrastructureSimple(t *testing.T, cfg *config.Config) *testSetup {
	t.Helper()

	mockRedis := redisPkg.InitMockRedis(t)
	mockDB := sqldbPkg.InitMockDb(t)

	projectRoot := getProjectRoot(t)
	privateKeyPath := filepath.Join(projectRoot, "pkg", "jwtUtils", "private.test.pem")
	publicKeyPath := filepath.Join(projectRoot, "pkg", "jwtUtils", "public.test.pem")

	jwtGenerator, err := jwtUtils.NewJwtGenerator(privateKeyPath)
	if err != nil {
		t.Fatalf("Failed to create JWT generator: %v", err)
	}

	jwtValidator, err := jwtUtils.NewJwtValidator(publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to create JWT validator: %v", err)
	}

	app := apipkg.New(cfg, mockRedis, mockDB, jwtGenerator, jwtValidator)

	return &testSetup{
		mockRedis:    mockRedis,
		mockDB:       mockDB,
		jwtGenerator: jwtGenerator,
		jwtValidator: jwtValidator,
		app:          app,
	}
}

// executeJSONRequest executes an HTTP request with JSON body
func executeJSONRequest(api apipkg.Engine, method, endpoint string, body interface{}) *httptest.ResponseRecorder {
	var jsonData []byte
	if body != nil {
		jsonData, _ = json.Marshal(body)
	}

	var req *http.Request
	if jsonData != nil {
		req = httptest.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	} else {
		req = httptest.NewRequest(method, endpoint, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	return rec
}

// executeRequest executes an HTTP request with a raw body (string or nil)
func executeRequest(api apipkg.Engine, method, endpoint string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, endpoint, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, endpoint, nil)
	}

	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	return rec
}

// createTestUser creates a test user in the database with the given credentials
func createTestUser(t *testing.T, db *gorm.DB, username, email, displayName, password string) *model.User {
	t.Helper()
	hashedPassword := utils.HashPassword(password)
	testUser := &model.User{
		Username:    username,
		Password:    hashedPassword,
		DisplayName: displayName,
		Email:       email,
	}
	require.NoError(t, db.Create(testUser).Error)
	return testUser
}

// createTestUserWithDefaults creates a test user with default test values
func createTestUserWithDefaults(t *testing.T, db *gorm.DB) *model.User {
	return createTestUser(t, db, "testuser", "test@example.com", "Test User", fixture.ValidTestPassword())
}

// defaultTestConfig returns a default test configuration
func defaultTestConfig() *config.Config {
	return &config.Config{
		ServiceName: "bookmark_service",
		InstanceId:  "",
	}
}

// Helper functions for endpoint paths
func getUserRegisterEndpoint() string {
	return "/v1" + routers.Endpoints.UserRegister
}

func getUserLoginEndpoint() string {
	return "/v1" + routers.Endpoints.AuthLogin
}

func getUserProfileEndpoint() string {
	return "/v1" + routers.Endpoints.GetProfile
}

// Response validation helpers

// validateBadRequestResponse validates a bad request response with Message and Details
func validateBadRequestResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()
	var resp struct {
		Message string   `json:"message"`
		Details []string `json:"details"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, expectedMessage, resp.Message)
	assert.NotEmpty(t, resp.Details)
}

// validateInvalidJSONResponse validates an invalid JSON response
func validateInvalidJSONResponse(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	var resp struct {
		Message string   `json:"message"`
		Details []string `json:"details"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	// Invalid JSON returns InternalErrorResponse, not validation error
	assert.Equal(t, "Something went wrong", resp.Message)
	// Details can be nil or empty for invalid JSON
}

// validateUnauthorizedResponse validates an unauthorized response with Error field
func validateUnauthorizedResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedErrorSubstring string) {
	t.Helper()
	var resp struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Contains(t, resp.Error, expectedErrorSubstring)
}

// executeGetRequestWithAuth executes a GET request with Authorization header
// If token doesn't start with "Bearer ", it will be prefixed automatically
func executeGetRequestWithAuth(api apipkg.Engine, endpoint, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	if token != "" {
		authHeader := token
		if !strings.HasPrefix(token, "Bearer ") {
			authHeader = "Bearer " + token
		}
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	return rec
}
