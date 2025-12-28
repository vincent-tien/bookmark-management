package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		setupMockSvc   func(t *testing.T) *mocks.UrlShorten
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
				ctx.Request = httptest.NewRequest(http.MethodPost, routers.Endpoints.LinkShorten, bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("Shorten", mock.Anything, dto.LinkShortenRequestDto{
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
				ctx.Request = httptest.NewRequest(http.MethodPost, routers.Endpoints.LinkShorten, strings.NewReader("invalid json"))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T) *mocks.UrlShorten {
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
				ctx.Request = httptest.NewRequest(http.MethodPost, routers.Endpoints.LinkShorten, bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T) *mocks.UrlShorten {
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
				ctx.Request = httptest.NewRequest(http.MethodPost, routers.Endpoints.LinkShorten, bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				mockSvc.On("Shorten", mock.Anything, dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}).Return("", errors.New("database error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "",
		},
		{
			name: "duplicate key retry success",
			setupRequest: func(ctx *gin.Context) {
				reqBody := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}
				jsonData, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, routers.Endpoints.LinkShorten, bytes.NewBuffer(jsonData))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T) *mocks.UrlShorten {
				mockSvc := mocks.NewUrlShorten(t)
				reqDto := dto.LinkShortenRequestDto{
					ExpInSeconds: 3600,
					Url:          "https://google.com",
				}
				// First call returns duplicate key error, second call succeeds
				mockSvc.On("Shorten", mock.Anything, reqDto).Return("", e.ErrKeyAlreadyExists).Once()
				mockSvc.On("Shorten", mock.Anything, reqDto).Return("xyz12345", nil).Once()
				return mockSvc
			},
			expectedStatus: http.StatusCreated,
			expectedResp:   `{"code":"xyz12345","message":"Shorten URL generated successfully!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockSvc(t)
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
