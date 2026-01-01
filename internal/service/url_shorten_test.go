package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	e "github.com/vincent-tien/bookmark-management/internal/errors"
	"github.com/vincent-tien/bookmark-management/internal/repository/mocks"
)

func TestUrlShorten_Shorten(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		name                    string
		setupMockUrlStorageRepo func() *mocks.UrlStorage
		request                 dto.LinkShortenRequestDto
		expectedError           error
		validateResult          func(t *testing.T, code string, err error)
	}{
		{
			name: "success",
			setupMockUrlStorageRepo: func() *mocks.UrlStorage {
				mockStorage := mocks.NewUrlStorage(t)
				//// Mock CheckKeyExists to return false (key doesn't exist)
				mockStorage.On("CheckKeyExists", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				})).Return(false, nil)
				// Mock Store to succeed
				mockStorage.On("Store", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				}), mock.Anything).Return(nil)

				return mockStorage
			},
			request: dto.LinkShortenRequestDto{
				Url:          "https://example.com",
				ExpInSeconds: 3600,
			},
			expectedError: nil,
			validateResult: func(t *testing.T, code string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
				assert.Len(t, code, 8)
			},
		},
		{
			name: "key already exists",
			setupMockUrlStorageRepo: func() *mocks.UrlStorage {
				mockStorage := mocks.NewUrlStorage(t)
				// Mock CheckKeyExists to return true (key exists) - can be called multiple times during retries
				mockStorage.On("CheckKeyExists", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				})).Return(true, nil)

				return mockStorage
			},
			request: dto.LinkShortenRequestDto{
				Url:          "https://example.com",
				ExpInSeconds: 3600,
			},
			expectedError:  e.ErrKeyAlreadyExists,
			validateResult: nil,
		},
		{
			name: "Store returns error",
			setupMockUrlStorageRepo: func() *mocks.UrlStorage {
				mockStorage := mocks.NewUrlStorage(t)
				// Mock CheckKeyExists to return false (key doesn't exist)
				mockStorage.On("CheckKeyExists", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				})).Return(false, nil)
				// Mock Store to return an error
				mockStorage.On("Store", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				}), mock.Anything).Return(assert.AnError)

				return mockStorage
			},
			request: dto.LinkShortenRequestDto{
				Url:          "https://example.com",
				ExpInSeconds: 3600,
			},
			expectedError:  assert.AnError,
			validateResult: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tc.setupMockUrlStorageRepo()
			service := NewUrlShorten(mockStorage)

			ctx := t.Context()
			code, err := service.Shorten(ctx, tc.request)

			if tc.validateResult != nil {
				tc.validateResult(t, code, err)
			} else {
				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestUrlShorten_GetUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                    string
		setupMockUrlStorageRepo func() *mocks.UrlStorage
		validateResult          func(t *testing.T, code string, err error)
	}{
		{
			name: "normal case",
			setupMockUrlStorageRepo: func() *mocks.UrlStorage {
				mockStorage := mocks.NewUrlStorage(t)
				//// Mock CheckKeyExists to return false (key doesn't exist)
				mockStorage.On("GetUrl", mock.Anything, mock.MatchedBy(func(code string) bool {
					return len(code) == 8
				})).Return("https://google.com", nil)

				return mockStorage
			},
			validateResult: func(t *testing.T, url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://google.com", url)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := tc.setupMockUrlStorageRepo()
			service := NewUrlShorten(mockStorage)

			ctx := t.Context()

			code := "12345678"
			url, err := service.GetUrl(ctx, code)

			tc.validateResult(t, url, err)
		})
	}
}
