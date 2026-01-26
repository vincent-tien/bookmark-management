package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	logPkg "github.com/rs/zerolog/log"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	errorsPkg "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/service"
	"github.com/vincent-tien/bookmark-management/pkg/response"
	"github.com/vincent-tien/bookmark-management/pkg/utils"
)

/*
User defines the interface for user handlers.
It provides methods to handle user registration.
*/
type User interface {
	/*
		Register processes a user registration request and returns a response with created user details or an error status.
		@Summary User Registration
		@Description	Register a new user with display name, email, username, and password
		@param c *gin.Context	the Gin context object
		@return void
	*/
	Register(c *gin.Context)

	Login(c *gin.Context)

	GetProfile(c *gin.Context)

	UpdateProfile(c *gin.Context)
}

type user struct {
	userService service.User
}

// bindUpdateProfileRequest binds the request body to UpdateUserProfileRequestDto
// and sets the user ID from the JWT context for security.
// Returns the request DTO and a boolean indicating success.
// If binding fails or user ID cannot be extracted, it returns nil and false.
func bindUpdateProfileRequest(c *gin.Context) (*dto.UpdateUserProfileRequestDto, bool) {
	userId, ok := utils.GetUserIDFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return nil, false
	}

	req, err := utils.BindJson[dto.UpdateUserProfileRequestDto](c)
	if err != nil {
		return nil, false
	}

	// Set the user ID from JWT context (not from request body for security)
	req.UserId = userId
	return req, true
}

// Register processes a user registration request and returns a response with created user details or an error status.
//
//	@Summary User Registration
//	@Description	Register a new user with display name, email, username, and password
//	@Tags Users
//	@Accept json
//	@Produce json
//	@Param request body dto.RegisterRequestDto true "User registration request payload"
//	@Success 200 {object} dto.RegisterSuccessResponse "Successfully registered user"
//	@Failure 400 {object} response.Response "Invalid request body or validation error"
//	@Failure 500 {object} response.Response "Internal server error"
//	@Router /v1/users/register [post]
func (u *user) Register(c *gin.Context) {
	req, err := utils.BindJson[dto.RegisterRequestDto](c)
	if err != nil {
		return
	}

	responseDto, err := u.userService.Register(c, *req)
	if err != nil {
		logPkg.Error().Err(err).Msg("Failed to Register")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    responseDto,
		"message": "Register an user successfully!",
	})
}

// Login processes a user login request and returns a response with a JWT token or an error status.
//
//	@Summary User Login
//	@Description Login a user with username and password
//	@Tags Users
//	@Accept json
//	@Produce json
//	@Param request body dto.LoginRequestDto true "User login request payload"
//	@Success 200 {object} dto.LoginSuccessResponse "Successfully logged in user"
//	@Failure 400 {object} response.Response "Invalid request body or validation error"
//	@Failure 500 {object} response.Response "Internal server error"
//	@Router /v1/users/login [post]
func (u *user) Login(c *gin.Context) {
	req, err := utils.BindJson[dto.LoginRequestDto](c)
	if err != nil {
		return
	}

	token, err := u.userService.Login(c, *req)
	if err != nil {
		if errors.Is(err, errorsPkg.ErrInvalidAuth) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logPkg.Error().Err(err).Msg("Failed to Login")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
	}

	c.JSON(http.StatusOK, response.Success(token, "Logged in successfully!"))
}

// GetProfile returns the user profile information.
//
//	@Summary		Get user profile
//	@Description	Get user profile information
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200 {object} dto.UserProfileResponseDto "User profile"
//	@Failure		401 {object} response.Response "Unauthorized"
//	@Failure		500 {object} response.Response "Internal server error"
//	@Security		BearerAuth
//	@Router			/v1/self/info [get]
func (u *user) GetProfile(c *gin.Context) {
	userId, ok := utils.GetUserIDFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	userModel, err := u.userService.GetProfile(c, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Convert model.User to UserProfileResponseDto
	responseDto := dto.UserProfileResponseDto{
		UserId:      userModel.ID,
		DisplayName: userModel.DisplayName,
		Username:    userModel.Username,
		Email:       userModel.Email,
		CreatedAt:   userModel.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   userModel.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response.Success(responseDto))
}

// UpdateProfile updates the user profile information.
//
//	@Summary		Update user profile
//	@Description	Update user profile information (display name and/or email)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request body dto.UpdateUserProfileRequestDto true "User profile update request payload"
//	@Success		200 {object} response.Response "Profile updated successfully"
//	@Failure		400 {object} response.Response "Invalid request body or validation error"
//	@Failure		401 {object} response.Response "Unauthorized"
//	@Failure		500 {object} response.Response "Internal server error"
//	@Security		BearerAuth
//	@Router			/v1/self/info [put]
func (u *user) UpdateProfile(c *gin.Context) {
	req, ok := bindUpdateProfileRequest(c)
	if !ok {
		return
	}

	err := u.userService.UpdateProfile(c, *req)
	if err != nil {
		logPkg.Error().Err(err).Msg("Failed to UpdateProfile")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Edit current user successfully!",
	})
}

func NewUserHandler(us service.User) User {
	return &user{
		userService: us,
	}
}
