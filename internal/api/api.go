package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/handler"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

// Engine defines the interface for the API engine.
// It provides methods to start the server and serve HTTP requests.
type Engine interface {
	// Start starts the HTTP server on the configured port.
	// It also registers the Swagger documentation endpoint.
	// Returns an error if the server fails to start.
	Start() error
	// ServeHTTP serves HTTP requests using the underlying gin engine.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app *gin.Engine
	cfg *config.Config
}

// Start starts the HTTP server on the configured port.
// It also registers the Swagger documentation endpoint.
// Returns an error if the server fails to start.
func (a *api) Start() error {
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP serves HTTP requests using the underlying gin engine.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// New creates and initializes a new API engine instance.
// It sets up the gin router, registers all endpoints, and returns an Engine interface.
// The configuration is used to set up the application settings.
func New(cfg *config.Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}
	a.registerEP(cfg)
	return a
}

func (a *api) registerEP(cfg *config.Config) {
	uuidSvc := service.NewUuid()
	uuidHandler := handler.NewHealthCheck(uuidSvc, cfg)
	a.app.GET(routers.Endpoints.HealthCheck, uuidHandler.DoCheck)
}
