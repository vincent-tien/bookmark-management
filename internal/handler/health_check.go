package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/repository"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

// HealthCheck defines the interface for health check handlers.
// It provides a method to perform health check operations.
type HealthCheck interface {
	// DoCheck performs a health check and returns the service status.
	// It responds with the service name and instance ID.
	DoCheck(c *gin.Context)
}

type healthCheckHandler struct {
	svc           service.Uuid
	cfg           *config.Config
	uuid          string
	pingRedisRepo repository.PingRedis
}

// NewHealthCheck creates and returns a new health check handler instance.
// It initializes the handler with a UUID service, configuration, and Redis client.
// If no instance ID is provided in the config, it generates a new UUID.
// Returns a HealthCheck interface implementation.
func NewHealthCheck(svc service.Uuid, cfg *config.Config, repo repository.PingRedis) HealthCheck {
	var err error

	uuid := cfg.InstanceId

	if uuid == "" {
		uuid, err = svc.Generate()
	}

	if err != nil {
		log.Printf("Failed to generate uuid: %v", err)
		uuid = ""
	}

	return &healthCheckHandler{
		svc:           svc,
		cfg:           cfg,
		uuid:          uuid,
		pingRedisRepo: repo,
	}
}

// HealthCheckResponse represents the response structure for health check endpoints.
type HealthCheckResponse struct {
	Message     string `json:"message"`      // Status message
	ServiceName string `json:"service_name"` // Name of the service
	InstanceId  string `json:"instance_id"`  // Unique instance identifier
}

// DoCheck performs a health check and returns the service status.
// It checks the Redis connection and responds with the service name and instance ID.
// Returns HTTP 200 OK if the service is healthy, or HTTP 500 if UUID generation failed or Redis ping fails.
//
//	@Summary		health check
//	@Tags		utils
//	@Description	check health
//	@Accept			json
//	@Router			/health-check [get]
func (h *healthCheckHandler) DoCheck(c *gin.Context) {
	if h.uuid == "" {
		c.String(http.StatusInternalServerError, "Failed to generate uuid")
		return
	}

	err := h.pingRedisRepo.Ping(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.cfg.ServiceName,
		InstanceId:  h.uuid,
	})
}
