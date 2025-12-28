package dto

// ErrorResponse represents error response
//
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error code
	// example: VALIDATION_ERROR
	Code string `json:"code"`

	// Error message
	// example: invalid url format
	Message string `json:"message"`
}
