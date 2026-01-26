package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vincent-tien/bookmark-management/pkg/response"
)

// BindJson binds the request body to a struct of type T and validates it.
//
// If binding fails, it returns nil and an error. If validation fails, it
// returns nil and an error.
//
// It is a convenience wrapper around gin.Context.ShouldBindJSON and
// validateStruct.
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
