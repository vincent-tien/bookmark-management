package repository

import (
	"context"
	"fmt"

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
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	GetUserById(ctx context.Context, userId string) (*model.User, error)

	UpdateProfile(ctx context.Context, userId string, updates map[string]interface{}) error
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

func (u *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.getUserByIField(ctx, "username", username)
}

func (u *user) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	return u.getUserByIField(ctx, "id", userId)
}

func (u *user) getUserByIField(ctx context.Context, fieldName, fieldValue string) (*model.User, error) {
	chosenUser := &model.User{}
	err := u.db.WithContext(ctx).Where(fmt.Sprintf("%s=?", fieldName), fieldValue).First(chosenUser).Error
	if err != nil {
		return nil, err
	}
	return chosenUser, nil
}

func (u *user) UpdateProfile(ctx context.Context, userId string, updates map[string]interface{}) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Updates(updates).Error
}
