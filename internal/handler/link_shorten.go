package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

// LinkShorten defines the interface for link shortening handlers.
// It provides methods to handle link shortening operations.
type LinkShorten interface {
	// Create handles the creation of a shortened link.
	// It validates the request, generates a short code, and stores the mapping.
	Create(c *gin.Context)
}

type linkShorten struct {
	svc service.UrlShorten
}

// NewLinkShorten creates and returns a new link shortening handler instance.
// It initializes the handler with a URL shortening service.
// Returns a LinkShorten interface implementation.
func NewLinkShorten(svc service.UrlShorten) LinkShorten {
	return &linkShorten{
		svc: svc,
	}
}

// Create CreateShortLink godoc
//
// @Summary      Create a shortened link
// @Description  Generate a short URL with expiration time
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        request body dto.LinkShortenRequestDto true "Shorten link request payload"
// @Success      200 {object} dto.LinkShortenResponseDto
// @Failure      400 {object} dto.ErrorResponse "Invalid request body or validation error"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /v1/links/shorten [post]
func (s *linkShorten) Create(c *gin.Context) {
	var req dto.LinkShortenRequestDto
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Prepare()

	code, err := s.svc.Shorten(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	res := dto.LinkShortenResponseDto{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	}
	c.JSON(http.StatusCreated, res)
}
