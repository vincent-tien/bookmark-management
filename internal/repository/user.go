package repository

import (
	"context"
	"fmt"

	"github.com/vincent-tien/bookmark-management/internal/dto"
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

	UpdateProfile(ctx context.Context, dto dto.UpdateUserProfileRequestDto) error
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

func (u *user) UpdateProfile(ctx context.Context, dto dto.UpdateUserProfileRequestDto) error {
	// First, get the existing user
	existingUser, err := u.GetUserById(ctx, dto.UserId)
	if err != nil {
		return err
	}

	// Update only the fields that are provided (non-empty)
	updates := make(map[string]interface{})
	if dto.DisplayName != "" {
		updates["display_name"] = dto.DisplayName
	}
	if dto.Email != "" {
		updates["email"] = dto.Email
	}

	// If no fields to update, return the existing user
	if len(updates) == 0 {
		return nil
	}

	err = u.db.WithContext(ctx).Model(existingUser).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
