package repository

import (
	"context"

	"github.com/vincent-tien/bookmark-management/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=User --filename=user.go

type User interface {
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
