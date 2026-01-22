package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	// Combined regex pattern for password validation (RE2 compatible):
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
	upperRegex    = regexp.MustCompile(`[A-Z]`)
	lowerRegex    = regexp.MustCompile(`[a-z]`)
	numberRegex   = regexp.MustCompile(`[0-9]`)
	specialRegex  = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// RegisterCustomValidators registers custom validation functions
func RegisterCustomValidators(v *validator.Validate) error {
	// Register password strength validator
	if err := v.RegisterValidation("strong_password", validateStrongPassword); err != nil {
		return err
	}
	return nil
}

// validateStrongPassword validates password using regex patterns:
// - At least 8 characters long
// - At least 1 uppercase letter
// - At least 1 lowercase letter
// - At least 1 number
// - At least 1 special character
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Use regex patterns to validate all conditions
	return passwordRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		lowerRegex.MatchString(password) &&
		numberRegex.MatchString(password) &&
		specialRegex.MatchString(password)
}
