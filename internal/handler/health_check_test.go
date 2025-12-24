package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service/mocks"
)

func TestUuidService_DoCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func() *mocks.Uuid
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success case",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, routers.Endpoints.HealthCheck, nil)
			},
			setupMockSvc: func() *mocks.Uuid {
				mockSvc := mocks.NewUuid(t)
				mockSvc.On("Generate").Return("12345678-1234-5678-9abc-def012345678", nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `{"message":"OK","service_name":"bookmark_service","instance_id":"12345678-1234-5678-9abc-def012345678"}`,
		},
		{
			name: "internal server err",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, routers.Endpoints.HealthCheck, nil)
			},
			setupMockSvc: func() *mocks.Uuid {
				mockSvc := mocks.NewUuid(t)
				mockSvc.On("Generate").Return("", errors.New("something wrong"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   `Failed to generate uuid`,
		},
	}

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceId:  "",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)

			mockSvc := tc.setupMockSvc()

			handler := NewHealthCheck(mockSvc, cfg)
			handler.DoCheck(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
