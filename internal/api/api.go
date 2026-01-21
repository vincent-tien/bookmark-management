package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vincent-tien/bookmark-management/docs"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/handler"
	"github.com/vincent-tien/bookmark-management/internal/repository"
	"github.com/vincent-tien/bookmark-management/internal/routers"
	"github.com/vincent-tien/bookmark-management/internal/service"
	validationPkg "github.com/vincent-tien/bookmark-management/pkg/validation"
	"gorm.io/gorm"
)

const (
	Version = "v1"
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
	app         *gin.Engine
	cfg         *config.Config
	redisClient *redis.Client
	db          *gorm.DB
}

// Start starts the HTTP server on the configured port.
// It also registers the Swagger documentation endpoint.
// Returns an error if the server fails to start.
func (a *api) Start() error {
	docs.SwaggerInfo.Host = a.cfg.AppHostName
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
func New(cfg *config.Config, redisClient *redis.Client, db *gorm.DB) Engine {
	a := &api{
		app:         gin.New(),
		cfg:         cfg,
		redisClient: redisClient,
		db:          db,
	}
	a.registerValidators()
	a.registerEP()
	return a
}

// registerValidators registers custom validation functions
func (a *api) registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validationPkg.RegisterCustomValidators(v); err != nil {
			panic(fmt.Sprintf("Failed to register custom validators: %v", err))
		}
	}
}

// registerEP registers all API endpoints and sets up their dependencies.
func (a *api) registerEP() {
	a.registerHealthCheckEndpoint()
	a.registerLinkShortenEndpoint()
	a.registerUsersEndpoint()
}

// registerHealthCheckEndpoint registers the health check endpoint.
func (a *api) registerHealthCheckEndpoint() {
	uuidSvc := service.NewUuid()
	repo := repository.NewPingRedis(a.redisClient)
	healthCheckHandler := handler.NewHealthCheck(uuidSvc, a.cfg, repo)
	a.app.GET(routers.Endpoints.HealthCheck, healthCheckHandler.DoCheck)
}

// registerLinkShortenEndpoint registers the link shorten endpoint.
func (a *api) registerLinkShortenEndpoint() {
	urlStorage := repository.NewUrlStorage(a.redisClient)
	linkShortenSvc := service.NewUrlShorten(urlStorage)
	linkShortenHandler := handler.NewLinkShorten(linkShortenSvc)

	apiVersion := a.app.Group(fmt.Sprintf("/%s", Version))
	{
		apiVersion.POST(routers.Endpoints.LinkShorten, linkShortenHandler.Create)
		apiVersion.GET(routers.Endpoints.LinkRedirect, linkShortenHandler.Redirect)
	}
}

// registerUsersEndpoint registers the API endpoint for user-related operations at the path specified in Endpoints.Users.
func (a *api) registerUsersEndpoint() {
	userRepo := repository.NewUserRepository(a.db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	apiVersion := a.app.Group(fmt.Sprintf("/%s", Version))
	{
		apiVersion.POST(routers.Endpoints.UserRegister, userHandler.Register)
	}
}
