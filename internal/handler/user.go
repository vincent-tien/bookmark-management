package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logPkg "github.com/rs/zerolog/log"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/service"
	"github.com/vincent-tien/bookmark-management/pkg/response"
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
}

type user struct {
	userService service.User
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
func (u user) Register(c *gin.Context) {
	var req dto.RegisterRequestDto
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	responseDto, err := u.userService.Register(c, req)
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

func NewUserHandler(us service.User) User {
	return &user{
		userService: us,
	}
}
