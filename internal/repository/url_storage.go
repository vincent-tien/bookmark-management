package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vincent-tien/bookmark-management/internal/dto"
)

//go:generate mockery --name=UrlStorage --filename=url_storage.go

// UrlStorage defines the interface for URL storage operations.
// It provides methods to store, retrieve, and check the existence of URL mappings.
type UrlStorage interface {
	// Store stores a URL mapping with the given code and expiration time.
	// Returns an error if the storage operation fails.
	Store(ctx context.Context, code string, r dto.LinkShortenRequestDto) error
	// GetUrl retrieves the original URL associated with the given code.
	// Returns the URL string and an error if the code is not found or retrieval fails.
	GetUrl(ctx context.Context, code string) (string, error)
	// CheckKeyExists checks if a code already exists in storage.
	// Returns true if the code exists, false otherwise, and an error if the check fails.
	CheckKeyExists(ctx context.Context, code string) (bool, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new UrlStorage with the provided redis client.
// This allows for easy mocking in tests by passing a mock redis client.
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// Store stores a URL mapping with the given code and expiration time.
// Returns an error if the storage operation fails.
func (s *urlStorage) Store(ctx context.Context, code string, r dto.LinkShortenRequestDto) error {
	return s.c.Set(ctx, code, r.Url, time.Second*time.Duration(r.ExpInSeconds)).Err()
}

// GetUrl retrieves the original URL associated with the given code.
// Returns the URL string and an error if the code is not found or retrieval fails.
func (s *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}

// CheckKeyExists checks if a code already exists in storage.
// Returns true if the code exists, false otherwise, and an error if the check fails.
func (s *urlStorage) CheckKeyExists(ctx context.Context, code string) (bool, error) {
	count, err := s.c.Exists(ctx, code).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
