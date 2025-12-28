package repository

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	redisPkg "github.com/vincent-tien/bookmark-management/pkg/redis"
)

func TestUrlStorage_Store(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setupMock  func() *redis.Client
		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "store url",
			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				url, err := r.Get(ctx, "12345678").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			err := testRepo.Store(ctx, "12345678", dto.LinkShortenRequestDto{
				ExpInSeconds: 1,
				Url:          "https://google.com",
			})

			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

func TestUrlStorage_GetUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		code        string
		expectedUrl string
		expectErr   error
	}{
		{
			name: "successfully get existing url",
			setupMock: func(ctx context.Context) *redis.Client {
				redisMock := redisPkg.InitMockRedis(t)
				redisMock.Set(ctx, "12345678", "https://google.com", 0)
				return redisMock
			},
			code:        "12345678",
			expectedUrl: "https://google.com",
			expectErr:   nil,
		},
		{
			name: "get non-existent key returns redis.Nil error",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			code:        "nonexistent",
			expectedUrl: "",
			expectErr:   redis.Nil,
		},
		{
			name: "get url with different code",
			setupMock: func(ctx context.Context) *redis.Client {
				redisMock := redisPkg.InitMockRedis(t)
				redisMock.Set(ctx, "abcdefgh", "https://example.com", 0)
				return redisMock
			},
			code:        "abcdefgh",
			expectedUrl: "https://example.com",
			expectErr:   nil,
		},
		{
			name: "redis connection error",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisMock := tc.setupMock(ctx)
			testRepo := NewUrlStorage(redisMock)

			url, err := testRepo.GetUrl(ctx, tc.code)

			if tc.expectErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUrl, url)
			}
		})
	}
}

func TestUrlStorage_CheckKeyExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMock      func(ctx context.Context) *redis.Client
		code           string
		expectedExists bool
		expectErr      error
	}{
		{
			name: "key exists returns true",
			setupMock: func(ctx context.Context) *redis.Client {
				redisMock := redisPkg.InitMockRedis(t)
				redisMock.Set(ctx, "12345678", "https://google.com", 0)
				return redisMock
			},
			code:           "12345678",
			expectedExists: true,
			expectErr:      nil,
		},
		{
			name: "key does not exist returns false",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			code:           "nonexistent",
			expectedExists: false,
			expectErr:      nil,
		},
		{
			name: "check different existing key",
			setupMock: func(ctx context.Context) *redis.Client {
				redisMock := redisPkg.InitMockRedis(t)
				redisMock.Set(ctx, "abcdefgh", "https://example.com", 0)
				return redisMock
			},
			code:           "abcdefgh",
			expectedExists: true,
			expectErr:      nil,
		},
		{
			name: "check non-existent key with other keys present",
			setupMock: func(ctx context.Context) *redis.Client {
				redisMock := redisPkg.InitMockRedis(t)
				// Set a different key to ensure we're checking the right one
				redisMock.Set(ctx, "otherkey", "https://other.com", 0)
				return redisMock
			},
			code:           "notset",
			expectedExists: false,
			expectErr:      nil,
		},
		{
			name: "redis connection error",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisMock := tc.setupMock(ctx)
			testRepo := NewUrlStorage(redisMock)

			exists, err := testRepo.CheckKeyExists(ctx, tc.code)

			if tc.expectErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedExists, exists)
			}
		})
	}
}
