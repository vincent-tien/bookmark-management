package service

import (
	"context"
	"time"

	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/repository"
	"github.com/vincent-tien/bookmark-management/pkg/utils"
)

//go:generate mockery --name=User --filename=user.go

// User defines the interface for user services.
// It provides methods to handle user registration.
type User interface {

	/*
		Register registers a new user.

		The function takes a context and a dto.RegisterRequestDto as parameters.
		It returns a dto.RegisterResponseDto and an error.
	*/
	Register(ctx context.Context, r dto.RegisterRequestDto) (dto.RegisterResponseDto, error)
}

type user struct {
	userRepository repository.User
}

func NewUserService(repo repository.User) User {
	return &user{
		userRepository: repo,
	}
}

func (u *user) Register(ctx context.Context, r dto.RegisterRequestDto) (dto.RegisterResponseDto, error) {
	// Hash the password
	hashedPassword := utils.HashPassword(r.Password)

	// Create user model
	userModel := &model.User{
		Username:    r.Username,
		Password:    hashedPassword,
		DisplayName: r.DisplayName,
		Email:       r.Email,
	}

	// Create user in repository
	createdUser, err := u.userRepository.CreateUser(ctx, userModel)
	if err != nil {
		return dto.RegisterResponseDto{}, err
	}

	// Convert to response DTO
	now := time.Now().Format(time.RFC3339)
	response := dto.RegisterResponseDto{
		ID:          createdUser.ID,
		Username:    createdUser.Username,
		DisplayName: createdUser.DisplayName,
		Email:       createdUser.Email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return response, nil
}
