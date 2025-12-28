package service

import (
	"context"

	"github.com/vincent-tien/bookmark-management/internal/dto"
	e "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/repository"
	"github.com/vincent-tien/bookmark-management/pkg/utils"
)

const (
	urlCodeLength = 8
)

//go:generate mockery --name=UrlShorten --filename=url_shorten.go

// UrlShorten defines the interface for URL shortening services.
// It provides methods to generate short codes and store URL mappings.
type UrlShorten interface {
	// Shorten generates a short code for the given URL and stores the mapping.
	// It returns the generated short code and an error if the operation fails.
	Shorten(ctx context.Context, r dto.LinkShortenRequestDto) (string, error)
}

type urlShorten struct {
	repo repository.UrlStorage
}

// NewUrlShorten creates and returns a new URL shortening service instance.
// It initializes the service with a URL storage repository.
// Returns a UrlShorten interface implementation.
func NewUrlShorten(repo repository.UrlStorage) UrlShorten {
	return &urlShorten{
		repo: repo,
	}
}

// Shorten generates a short code for the given URL and stores the mapping.
// It creates a random code, checks for duplicates, and stores the URL with expiration.
// Returns the generated short code and an error if the operation fails.
func (s *urlShorten) Shorten(ctx context.Context, r dto.LinkShortenRequestDto) (string, error) {
	code, err := utils.GenerateRandomString(urlCodeLength)
	if err != nil {
		return "", err
	}

	isExist, err := s.checkKeyExists(ctx, code)
	if isExist {
		return "", e.ErrKeyAlreadyExists
	}

	err = s.repo.Store(ctx, code, r)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *urlShorten) checkKeyExists(ctx context.Context, code string) (bool, error) {
	return s.repo.CheckKeyExists(ctx, code)
}
