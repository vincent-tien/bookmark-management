package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service/mocks"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
)

func TestUuidService_DoCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func() *mocks.Uuid
		setupRedis     func() *redis.Client
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
			setupRedis: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `{"message":"OK","service_name":"bookmark_service","instance_id":"12345678-1234-5678-9abc-def012345678"}`,
		},
		{
			name: "internal server err - uuid generation failed",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, routers.Endpoints.HealthCheck, nil)
			},
			setupMockSvc: func() *mocks.Uuid {
				mockSvc := mocks.NewUuid(t)
				mockSvc.On("Generate").Return("", errors.New("something wrong"))
				return mockSvc
			},
			setupRedis: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   `Failed to generate uuid`,
		},
		{
			name: "internal server err - redis ping failed",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, routers.Endpoints.HealthCheck, nil)
			},
			setupMockSvc: func() *mocks.Uuid {
				mockSvc := mocks.NewUuid(t)
				mockSvc.On("Generate").Return("12345678-1234-5678-9abc-def012345678", nil)
				return mockSvc
			},
			setupRedis: func() *redis.Client {
				// Create a Redis client with invalid address to simulate connection failure
				return redis.NewClient(&redis.Options{
					Addr: "invalid:6379",
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   `Redis connection failed:`,
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
			redisClient := tc.setupRedis()

			handler := NewHealthCheck(mockSvc, cfg, redisClient)
			handler.DoCheck(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.name == "internal server err - redis ping failed" {
				// For Redis ping failure, just check that the response contains the error message
				assert.Contains(t, rec.Body.String(), tc.expectedResp)
			} else {
				assert.Equal(t, tc.expectedResp, rec.Body.String())
			}
		})
	}
}
