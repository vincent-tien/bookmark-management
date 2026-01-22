package repository

import (
	"context"

	"github.com/vincent-tien/bookmark-management/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=User --filename=user.go

// User defines the interface for user repository.
// It provides methods to handle user creation.
type User interface {
	// CreateUser creates a new user.
	// It takes a context and a User model as input and returns the created user and an error if any.
	CreateUser(ctx context.Context, uModel *model.User) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) User {
	return &user{db: db}
}

func (u *user) CreateUser(ctx context.Context, uModel *model.User) (*model.User, error) {
	err := u.db.WithContext(ctx).Create(uModel).Error
	if err != nil {
		return nil, err
	}
	return uModel, nil
}
