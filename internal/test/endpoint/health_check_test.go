package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupTestHttp  func(api apipkg.Engine) *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "success case",
			setupTestHttp: func(api apipkg.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, routers.Endpoints.HealthCheck, nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
		},
	}

	cfg := &config.Config{
		ServiceName: "bookmark_service",
		InstanceId:  "",
	}

	mockRedis := redisPkg.InitMockRedis(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := apipkg.New(cfg, mockRedis)
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
