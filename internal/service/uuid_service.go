package service

import "github.com/google/uuid"

//go:generate mockery --name=Uuid --filename=uuid_service.go

// Uuid defines the interface for UUID generation services.
// It provides a method to generate unique identifiers.
type Uuid interface {
	// Generate creates a new UUID v7 string.
	// Returns the UUID as a string and an error if generation fails.
	Generate() (string, error)
}

// UuidService implements the Uuid interface for generating UUIDs.
type UuidService struct {
}

// NewUuid creates and returns a new instance of UuidService.
// It implements the Uuid interface for UUID generation.
func NewUuid() Uuid {
	return &UuidService{}
}

// Generate creates a new UUID v7 string.
// Returns the UUID as a string and an error if generation fails.
func (u *UuidService) Generate() (string, error) {
	v7, err := uuid.NewV7()

	if err != nil {
		return "", err
	}

	return v7.String(), nil
}
