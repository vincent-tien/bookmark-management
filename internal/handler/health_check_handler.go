package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/internal/config"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

type HealthCheckHandler interface {
	DoCheck(c *gin.Context)
}

type HealthCheckServiceHandler struct {
	svc service.Uuid
	cfg *config.Config
}

func NewUuidHandler(svc service.Uuid, cfg *config.Config) HealthCheckHandler {
	return &HealthCheckServiceHandler{
		svc: svc,
		cfg: cfg,
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
func (h *HealthCheckServiceHandler) DoCheck(c *gin.Context) {
	var err error
	uuid := h.cfg.InstanceId

	if uuid == "" {
		uuid, err = h.svc.Generate()
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, GenerateUuidResponse{
			Message:     "OK",
			ServiceName: h.cfg.ServiceName,
			InstanceId:  uuid,
		})
	}
}
