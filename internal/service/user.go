package service

import (
	"context"
	"time"

	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/repository"
	"github.com/vincent-tien/bookmark-management/pkg/jwtUtils"
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
	/*
		Login logs in a user.

			The function takes a context and a dto.LoginRequestDto as parameters.
			It returns a JWT token and an error.
	*/
	Login(ctx context.Context, r dto.LoginRequestDto) (string, error)

	/*
		GetProfile retrieves a user by id.

		The function takes a context and a user id as parameters.
		It returns a user model and an error.
	*/
	GetProfile(ctx context.Context, userId string) (*model.User, error)

	UpdateProfile(ctx context.Context, requestDto dto.UpdateUserProfileRequestDto) error
}

type user struct {
	userRepository repository.User
	jwtGen         jwtUtils.JwtGenerator
}

func NewUserService(repo repository.User, jwtGen jwtUtils.JwtGenerator) User {
	return &user{
		userRepository: repo,
		jwtGen:         jwtGen,
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

func (u *user) Login(ctx context.Context, r dto.LoginRequestDto) (string, error) {
	// check user exist
	username, err := u.userRepository.GetUserByUsername(ctx, r.Username)
	if err != nil {
		return "", err
	}

	// check pass is valid
	isTokenValid := utils.VerifyPassword(r.RawPassword, username.Password)
	if !isTokenValid {
		return "", errors.ErrInvalidAuth
	}

	//create token
	jwtContent := u.jwtGen.GenerateContent(username.ID)
	jwtToken, err := u.jwtGen.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

/*
GetProfile retrieves a user by id.

The function takes a context and a user id as parameters.
It returns a user model and an error.
*/
func (u *user) GetProfile(ctx context.Context, userId string) (*model.User, error) {
	return u.userRepository.GetUserById(ctx, userId)
}

func (u *user) UpdateProfile(ctx context.Context, requestDto dto.UpdateUserProfileRequestDto) error {
	// First, check if the user exists
	_, err := u.userRepository.GetUserById(ctx, requestDto.UserId)
	if err != nil {
		return err
	}

	// Build updates map with only non-empty fields
	updates := make(map[string]interface{})
	if requestDto.DisplayName != "" {
		updates["display_name"] = requestDto.DisplayName
	}
	if requestDto.Email != "" {
		updates["email"] = requestDto.Email
	}

	// If no fields to update, return early
	if len(updates) == 0 {
		return nil
	}

	return u.userRepository.UpdateProfile(ctx, requestDto.UserId, updates)
}
