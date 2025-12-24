package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

type HealthCheck interface {
	DoCheck(c *gin.Context)
}

type healthCheckService struct {
	svc  service.Uuid
	cfg  *config.Config
	uuid string
}

func NewHealthCheck(svc service.Uuid, cfg *config.Config) HealthCheck {
	var err error

	uuid := cfg.InstanceId

	if uuid == "" {
		uuid, err = svc.Generate()
	}

	if err != nil {
		log.Printf("Failed to generate uuid: %v", err)
		uuid = ""
	}

	return &healthCheckService{
		svc:  svc,
		cfg:  cfg,
		uuid: uuid,
	}
}

type GenerateUuidResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceId  string `json:"instance_id"`
}

// DoCheck  	godoc
// @Summary		health check
// @Tags		utils
// @Description	check health
// @Accept			json
// @Router			/health-check [get]
func (h *healthCheckService) DoCheck(c *gin.Context) {
	if h.uuid == "" {
		c.String(http.StatusInternalServerError, "Failed to generate uuid")
	} else {
		c.JSON(http.StatusOK, GenerateUuidResponse{
			Message:     "OK",
			ServiceName: h.cfg.ServiceName,
			InstanceId:  h.uuid,
		})
	}
}
