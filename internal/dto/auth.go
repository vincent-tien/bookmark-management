package dto

// LoginRequestDto represents request payload for user login
//
// swagger:model LoginRequestDto
type LoginRequestDto struct {
	// Username
	// required: true
	// example: johndoe
	Username string `json:"username"`

	// Password
	// required: true
	// example: SecurePass123!
	RawPassword string `json:"password" binding:"gte=8"`
}

// LoginSuccessResponse represents the success response for user login
//
// swagger:model LoginSuccessResponse
type LoginSuccessResponse struct {
	// JWT token
	// example: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
	Data string `json:"data"`

	// Success message
	// example: Logged in successfully!
	Message string `json:"message"`
}
