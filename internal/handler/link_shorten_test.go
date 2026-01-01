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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	e "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service/mocks"
)

func TestLinkShorten_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success case",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("Shorten", ctx, dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}).Return("foobar", nil)
				return mockSvc
			},
			expectedStatus: http.StatusCreated,
			expectedResp:   `{"code":"foobar","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "bad request - invalid JSON",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, getEndpoint(), strings.NewReader("invalid json"))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				return mocks.NewUrlShorten(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "bad request - missing required fields",
			setupRequest: func(ctx *gin.Context) {
				reqBody := map[string]interface{}{}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				return mocks.NewUrlShorten(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, getEndpoint(), bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("Shorten", ctx, dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}).Return("", errors.New("redis error"))
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
			handler := NewLinkShorten(mockSvc)
			handler.Create(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedResp != "" {
				// Remove whitespace for comparison
				actualBody := strings.TrimSpace(rec.Body.String())
				expectedBody := strings.TrimSpace(tc.expectedResp)
				assert.Equal(t, expectedBody, actualBody)
			} else if tc.expectedStatus == http.StatusBadRequest {
				// For bad request, just verify it contains error
				assert.Contains(t, rec.Body.String(), "error")
			} else if tc.expectedStatus == http.StatusInternalServerError {
				// For internal server error, verify it contains error
				assert.Contains(t, rec.Body.String(), "error")
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestLinkShorten_GetUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten
		expectedStatus int
		expectedResp   string
		expectedLoc    string
	}{
		{
			name: "success case",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/foobar", nil)
				ctx.Params = gin.Params{gin.Param{Key: "code", Value: "foobar"}}
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("GetUrl", ctx, "foobar").Return("https://google.com", nil)
				return mockSvc
			},
			expectedStatus: http.StatusFound,
			expectedResp:   "",
			expectedLoc:    "https://google.com",
		},
		{
			name: "bad request - empty code",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				ctx.Params = gin.Params{gin.Param{Key: "code", Value: ""}}
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				return mocks.NewUrlShorten(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   `{"error":"code parameter is required"}`,
			expectedLoc:    "",
		},
		{
			name: "not found - URL not found",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/nonexistent", nil)
				ctx.Params = gin.Params{gin.Param{Key: "code", Value: "nonexistent"}}
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("GetUrl", ctx, "nonexistent").Return("", e.ErrUrlNotFound)
				return mockSvc
			},
			expectedStatus: http.StatusNotFound,
			expectedResp:   `{"error":"URL not found"}`,
			expectedLoc:    "",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/testcode", nil)
				ctx.Params = gin.Params{gin.Param{Key: "code", Value: "testcode"}}
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("GetUrl", ctx, "testcode").Return("", errors.New("redis connection error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   `{"error":"Internal Server Error"}`,
			expectedLoc:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockSvc(t, ctx)
			handler := NewLinkShorten(mockSvc)
			handler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedResp != "" {
				// Remove whitespace for comparison
				actualBody := strings.TrimSpace(rec.Body.String())
				expectedBody := strings.TrimSpace(tc.expectedResp)
				assert.Equal(t, expectedBody, actualBody)
			}
			if tc.expectedLoc != "" {
				// Check redirect location header
				assert.Equal(t, tc.expectedLoc, rec.Header().Get("Location"))
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func getEndpoint() string {
	return fmt.Sprintf("/v1/%s", routers.Endpoints.LinkShorten)
}
