package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	e "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/service"
)

// LinkShorten defines the interface for link shortening handlers.
// It provides methods to handle link shortening operations.
type LinkShorten interface {
	// Create handles the creation of a shortened link.
	// It validates the request, generates a short code, and stores the mapping.
	Create(c *gin.Context)
	// Redirect handles the redirection to the original URL based on the code.
	// It retrieves the original URL and redirects the user to it.
	Redirect(c *gin.Context)
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

	code, err := s.svc.Shorten(c, req)

	if err != nil {
		log.Error().Err(err).Msg("Failed to shorten URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	res := dto.LinkShortenResponseDto{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	}
	c.JSON(http.StatusCreated, res)
}

// Redirect RedirectLink godoc
//
// @Summary      Redirect to original URL
// @Description  Redirects user to the original URL based on the short code
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        code path string true "Short code"
// @Success      302 "Redirect to original URL"
// @Failure      404 {object} dto.ErrorResponse "URL not found"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /v1/links/redirect/{code} [get]
func (s *linkShorten) Redirect(c *gin.Context) {
	rawCode := c.Param("code")
	code := strings.TrimPrefix(rawCode, "/")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is required"})
		return
	}

	url, err := s.svc.GetUrl(c, code)
	if err != nil {
		// Check if it's a not found error
		if errors.Is(err, e.ErrUrlNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		log.Error().Err(err).Msg("Failed to get URL")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Redirect to the original URL
	c.Redirect(http.StatusFound, url)
}
