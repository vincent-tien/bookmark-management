package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vincent-tien/bookmark-management/pkg/response"
)

func BindJson[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}
	return &req, validateStruct(c, &req)
}

func validateStruct(c *gin.Context, v any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(v); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return err
	}
	return nil
}
