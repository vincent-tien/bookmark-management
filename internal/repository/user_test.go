package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/dto"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/test/fixture"
	"gorm.io/gorm"
)

// setupTestDB creates a test database with user fixture
func setupTestDB(t *testing.T) *gorm.DB {
	return fixture.NewFixture(t, &fixture.UserFixture{})
}

// normalizeTimeFields sets CreatedAt and UpdatedAt to zero time for comparison
func normalizeTimeFields(user *model.User) {
	if user != nil {
		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}
	}
}

// normalizeTimeFieldsForComparison normalizes time fields for both result and expected user
func normalizeTimeFieldsForComparison(result, expected *model.User) {
	normalizeTimeFields(result)
	normalizeTimeFields(expected)
}

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupDb         func(t *testing.T) *gorm.DB
		inputUser       *model.User
		expectedOut     *model.User
		expectErrString string
		verifyFunc      func(db *gorm.DB, user *model.User)
	}{
		{
			name:    "create success",
			setupDb: setupTestDB,
			inputUser: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd999",
				DisplayName: "John Doo",
				Username:    "John Doo",
				Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
				Email:       "john.doo@example.com",
			},
			expectErrString: "",
			expectedOut: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd999",
				DisplayName: "John Doo",
				Username:    "John Doo",
				Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
				Email:       "john.doo@example.com",
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}
				err := db.Where("username = ?", user.Username).First(checkUser).Error
				assert.Nil(t, err)
				normalizeTimeFieldsForComparison(checkUser, user)
				assert.Equal(t, checkUser, user)
			},
		},
		{
			name:    "error on duplicate username",
			setupDb: setupTestDB,
			inputUser: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd995",
				DisplayName: "John Doe",
				Username:    "John Doe",
				Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
				Email:       "john.doe.dup@example.com",
			},
			expectErrString: "UNIQUE constraint failed: users.username",
			expectedOut:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDb(t)
			testRepo := NewUserRepository(db)
			result, err := testRepo.CreateUser(ctx, tc.inputUser)

			if err != nil {
				assert.Nil(t, result)
				assert.ErrorContains(t, err, tc.expectErrString)
				return
			}

			assert.NotNil(t, result)
			normalizeTimeFieldsForComparison(result, tc.expectedOut)
			assert.Equal(t, tc.expectedOut, result)

			if tc.verifyFunc != nil {
				tc.verifyFunc(db, tc.expectedOut)
			}
		})
	}
}

func TestUser_GetUserById(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupDb         func(t *testing.T) *gorm.DB
		inputId         string
		expectedOut     *model.User
		expectErrString string
	}{
		{
			name:    "get user by id success",
			setupDb: setupTestDB,
			inputId: "deb745af-1a62-4efa-99a0-f06b274bd993",
			expectedOut: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "John Doe",
				Username:    "John Doe",
				Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
				Email:       "john.doe@example.com",
			},
			expectErrString: "",
		},
		{
			name:    "get user by id not found",
			setupDb: setupTestDB,
			inputId:         "deb745af-1a62-4efa-99a0-f06b274bd999",
			expectedOut:     nil,
			expectErrString: "record not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDb(t)
			testRepo := NewUserRepository(db)
			result, err := testRepo.GetUserById(ctx, tc.inputId)

			verifyGetUserResult(t, result, err, tc.expectedOut, tc.expectErrString)
		})
	}
}

func TestUser_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupDb         func(t *testing.T) *gorm.DB
		inputUsername   string
		expectedOut     *model.User
		expectErrString string
	}{
		{
			name:    "get user by username success",
			setupDb: setupTestDB,
			inputUsername: "John Doe",
			expectedOut: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "John Doe",
				Username:    "John Doe",
				Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
				Email:       "john.doe@example.com",
			},
			expectErrString: "",
		},
		{
			name:    "get user by username not found",
			setupDb: setupTestDB,
			inputUsername:   "NonExistentUser",
			expectedOut:     nil,
			expectErrString: "record not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDb(t)
			testRepo := NewUserRepository(db)
			result, err := testRepo.GetUserByUsername(ctx, tc.inputUsername)

			verifyGetUserResult(t, result, err, tc.expectedOut, tc.expectErrString)
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupDb         func(t *testing.T) *gorm.DB
		inputDto        dto.UpdateUserProfileRequestDto
		expectErrString string
		verifyFunc      func(db *gorm.DB, userId string, expectedUser *model.User)
	}{
		{
			name:    "update display name success",
			setupDb: setupTestDB,
			inputDto: dto.UpdateUserProfileRequestDto{
				UserId:      "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "John Updated",
				Email:       "",
			},
			expectErrString: "",
			verifyFunc: func(db *gorm.DB, userId string, expectedUser *model.User) {
				verifyUpdatedUserFields(t, db, userId, expectedUser, "John Doe", "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106")
			},
		},
		{
			name:    "update email success",
			setupDb: setupTestDB,
			inputDto: dto.UpdateUserProfileRequestDto{
				UserId:      "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "",
				Email:       "john.updated@example.com",
			},
			expectErrString: "",
			verifyFunc: func(db *gorm.DB, userId string, expectedUser *model.User) {
				checkUser := &model.User{}
				err := db.Where("id = ?", userId).First(checkUser).Error
				assert.Nil(t, err)
				assert.Equal(t, expectedUser.Email, checkUser.Email)
				assert.Equal(t, expectedUser.DisplayName, checkUser.DisplayName)
				// Username and Password should remain unchanged
				assert.Equal(t, "John Doe", checkUser.Username)
				assert.Equal(t, "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106", checkUser.Password)
			},
		},
		{
			name:    "update both display name and email success",
			setupDb: setupTestDB,
			inputDto: dto.UpdateUserProfileRequestDto{
				UserId:      "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			expectErrString: "",
			verifyFunc: func(db *gorm.DB, userId string, expectedUser *model.User) {
				verifyUpdatedUserFields(t, db, userId, expectedUser, "John Doe", "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106")
			},
		},
		{
			name:    "update with no fields to update success",
			setupDb: setupTestDB,
			inputDto: dto.UpdateUserProfileRequestDto{
				UserId:      "deb745af-1a62-4efa-99a0-f06b274bd993",
				DisplayName: "",
				Email:       "",
			},
			expectErrString: "",
			verifyFunc: func(db *gorm.DB, userId string, expectedUser *model.User) {
				verifyUpdatedUserFields(t, db, userId, &model.User{DisplayName: "John Doe", Email: "john.doe@example.com"}, "John Doe", "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106")
			},
		},
		{
			name:    "error on user not found",
			setupDb: setupTestDB,
			inputDto: dto.UpdateUserProfileRequestDto{
				UserId:      "deb745af-1a62-4efa-99a0-f06b274bd999",
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			expectErrString: "record not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDb(t)
			testRepo := NewUserRepository(db)
			err := testRepo.UpdateProfile(ctx, tc.inputDto)

			if tc.expectErrString != "" {
				assert.ErrorContains(t, err, tc.expectErrString)
				return
			}

			assert.Nil(t, err)

			if tc.verifyFunc != nil {
				// Build expected user for verification
				expectedUser := &model.User{
					DisplayName: tc.inputDto.DisplayName,
					Email:       tc.inputDto.Email,
				}
				// If DisplayName is empty, use the original value
				if tc.inputDto.DisplayName == "" {
					expectedUser.DisplayName = "John Doe"
				}
				// If Email is empty, use the original value
				if tc.inputDto.Email == "" {
					expectedUser.Email = "john.doe@example.com"
				}
				tc.verifyFunc(db, tc.inputDto.UserId, expectedUser)
			}
		})
	}
}

// verifyGetUserResult is a helper function that verifies the result of GetUserById or GetUserByUsername operations.
// It handles common error checking and assertion logic, including normalizing time fields for comparison.
func verifyGetUserResult(t *testing.T, result *model.User, err error, expectedOut *model.User, expectErrString string) {
	if err != nil {
		assert.Nil(t, result)
		assert.ErrorContains(t, err, expectErrString)
		return
	}

	assert.NotNil(t, result)
	normalizeTimeFieldsForComparison(result, expectedOut)
	assert.Equal(t, expectedOut, result)
}

// verifyUpdatedUserFields verifies that user fields were updated correctly and unchanged fields remain the same
func verifyUpdatedUserFields(t *testing.T, db *gorm.DB, userId string, expectedUser *model.User, expectedUsername, expectedPassword string) {
	t.Helper()
	checkUser := &model.User{}
	err := db.Where("id = ?", userId).First(checkUser).Error
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.DisplayName, checkUser.DisplayName)
	assert.Equal(t, expectedUser.Email, checkUser.Email)
	// Username and Password should remain unchanged
	assert.Equal(t, expectedUsername, checkUser.Username)
	assert.Equal(t, expectedPassword, checkUser.Password)
}
