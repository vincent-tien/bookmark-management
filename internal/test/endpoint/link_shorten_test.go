package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/routers"
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
		ApiVersion:  "v1",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := apipkg.New(cfg)
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedStatus == http.StatusBadRequest {
				assert.Contains(t, rec.Body.String(), "error")
			} else if tc.expectedStatus == http.StatusCreated {
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

func getApiEndpoint() string {
	return "/v1" + routers.Endpoints.LinkShorten
}
