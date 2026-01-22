package dto

// RegisterRequestDto represents request payload for user registration
//
// swagger:model RegisterRequestDto
type RegisterRequestDto struct {
	// User's display name
	// required: true
	// example: John Doe
	DisplayName string `json:"display_name" binding:"required"`

	// User's email address
	// required: true
	// format: email
	// example: john@example.com
	Email string `json:"email" binding:"required,email"`

	// User's password (minimum 8 characters, must contain uppercase, lowercase, number, and special character)
	// required: true
	// minLength: 8
	// example: SecurePass123!
	Password string `json:"password" binding:"required,min=8,strong_password"`

	// User's unique username
	// required: true
	// example: johndoe
	Username string `json:"username" binding:"required"`
}

// RegisterResponseDto represents response payload for user registration
//
// swagger:model RegisterResponseDto
type RegisterResponseDto struct {
	// User ID
	// example: 123
	ID string `json:"id"`

	// User's username
	// example: johndoe
	Username string `json:"username"`

	// User's email address
	// example: john@example.com
	Email string `json:"email"`

	// User's display name
	// example: John Doe
	DisplayName string `json:"display_name"`

	// Account creation timestamp
	// example: 2024-01-01T00:00:00Z
	CreatedAt string `json:"created_at"`

	// Last update timestamp
	// example: 2024-01-01T00:00:00Z
	UpdatedAt string `json:"updated_at"`
}

// RegisterSuccessResponse represents the success response wrapper for user registration
//
// swagger:model RegisterSuccessResponse
type RegisterSuccessResponse struct {
	// User registration data
	Data RegisterResponseDto `json:"data"`

	// Success message
	// example: Register an user successfully!
	Message string `json:"message"`
}
