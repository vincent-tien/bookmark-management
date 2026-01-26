package endpoint

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	apipkg "github.com/vincent-tien/bookmark-management/internal/api"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
	sqldbPkg "github.com/vincent-tien/bookmark-management/pkg/sqldb"
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
