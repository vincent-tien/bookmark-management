package service

import "github.com/google/uuid"

//go:generate mockery --name=Uuid --filename=uuid_service.go
type Uuid interface {
	Generate() (string, error)
}

type UuidService struct {
}

func NewUuid() Uuid {
	return &UuidService{}
}

func (u *UuidService) Generate() (string, error) {
	v7, err := uuid.NewV7()

	if err != nil {
		return "", err
	}

	return v7.String(), nil
}
