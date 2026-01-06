package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
)

func TestLinkShortenEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(api apipkg.Engine) *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "success case",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, getApiEndpoint(), bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "bad request - invalid JSON",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, getApiEndpoint(), strings.NewReader("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad request - missing required fields",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := map[string]interface{}{}
				jsonData, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, getApiEndpoint(), bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad request - invalid URL",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "not-a-valid-url",
				}
				jsonData, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, getApiEndpoint(), bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success case - default expiration",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				reqBody := dto.LinkShortenRequestDto{
					Url: "https://example.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, getApiEndpoint(), bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusCreated,
		},
	}

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceId:  "",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRedis := redisPkg.InitMockRedis(t)

			app := apipkg.New(cfg, mockRedis)
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			switch tc.expectedStatus {
			case http.StatusBadRequest:
				assert.Contains(t, rec.Body.String(), "error")
			case http.StatusCreated:
				var resp struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

				assert.NotEmpty(t, resp.Code)
				assert.Equal(t, "Shorten URL generated successfully!", resp.Message)
			}
		})
	}
}

func TestRedirectLinkEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(api apipkg.Engine, mockRedis *redis.Client) *httptest.ResponseRecorder
		expectedStatus int
		expectedLoc    string
	}{
		{
			name: "success case",
			setupTestHttp: func(api apipkg.Engine, mockRedis *redis.Client) *httptest.ResponseRecorder {
				// Pre-populate Redis with a code->URL mapping
				ctx := context.Background()
				mockRedis.Set(ctx, "testcode123", "https://google.com", 0)

				req := httptest.NewRequest(http.MethodGet, getRedirectEndpoint("testcode123"), nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusFound,
			expectedLoc:    "https://google.com",
		},
		{
			name: "not found - URL not found",
			setupTestHttp: func(api apipkg.Engine, mockRedis *redis.Client) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, getRedirectEndpoint("nonexistent"), nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusNotFound,
			expectedLoc:    "",
		},
		{
			name: "bad request - empty code",
			setupTestHttp: func(api apipkg.Engine, mockRedis *redis.Client) *httptest.ResponseRecorder {
				// Test with empty code parameter (trailing slash makes code empty)
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
			expectedLoc:    "",
		},
		{
			name: "success case - different URL",
			setupTestHttp: func(api apipkg.Engine, mockRedis *redis.Client) *httptest.ResponseRecorder {
				// Pre-populate Redis with a different code->URL mapping
				ctx := context.Background()
				mockRedis.Set(ctx, "abc12345", "https://example.com", 0)

				req := httptest.NewRequest(http.MethodGet, getRedirectEndpoint("abc12345"), nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusFound,
			expectedLoc:    "https://example.com",
		},
	}

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceId:  "",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRedis := redisPkg.InitMockRedis(t)

			app := apipkg.New(cfg, mockRedis)
			rec := tc.setupTestHttp(app, mockRedis)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			switch tc.expectedStatus {
			case http.StatusFound:
				// Check redirect location header
				assert.Equal(t, tc.expectedLoc, rec.Header().Get("Location"))
			case http.StatusNotFound:
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "not found")
			case http.StatusBadRequest:
				var resp struct {
					Error string `json:"error"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp.Error, "code parameter is required")
			}
		})
	}
}

func getApiEndpoint() string {
	return "/v1" + routers.Endpoints.LinkShorten
}

func getRedirectEndpoint(code string) string {
	// The route uses wildcard *code, so we append the code directly
	basePath := "/v1" + strings.TrimSuffix(routers.Endpoints.LinkRedirect, "*code")
	return basePath + code
}
