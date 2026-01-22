package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Response represents a generic API response
//
// swagger:model Response
type Response struct {
	// Response message
	// example: Invalid request
	Message string `json:"message"`

	// Additional response details (optional)
	// example: ["email is invalid email", "password is invalid min"]
	Details any `json:"details,omitempty"`
}

// common response messages
var (
	InternalErrorResponse = Response{Message: "Something went wrong", Details: nil}
	InvalidRequestError   = Response{Message: "Invalid request", Details: nil}
)

// InputFieldError Package response contains common response messages and helpers
func InputFieldError(e error) Response {
	if ok := errors.As(e, &validator.ValidationErrors{}); !ok {
		return InternalErrorResponse
	}

	var errs []string
	for _, err := range e.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid "+err.Tag())
	}

	return Response{Message: "Invalid request", Details: errs}
}
