package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/routers"
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
				return executeRequest(api, http.MethodGet, routers.Endpoints.HealthCheck, "")
			},
			expectedStatus: http.StatusOK,
		},
	}

	cfg := defaultTestConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setup := setupTestInfrastructureSimple(t, cfg)
			rec := tc.setupTestHttp(setup.app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
