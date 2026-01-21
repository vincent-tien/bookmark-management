package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-tien/bookmark-management/internal/model"
	"github.com/vincent-tien/bookmark-management/internal/test/fixture"
	"gorm.io/gorm"
)

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
			name: "create success",
			setupDb: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserFixture{})
			},
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
				// Skip comparing CreatedAt and UpdatedAt fields by setting them to zero
				timeZero := time.Time{}
				checkUser.CreatedAt = timeZero
				checkUser.UpdatedAt = timeZero
				user.CreatedAt = timeZero
				user.UpdatedAt = timeZero
				assert.Equal(t, checkUser, user)
			},
		},
		{
			name: "error on duplicate username",
			setupDb: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserFixture{})
			},
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

			// Skip comparing CreatedAt and UpdatedAt fields by setting them to zero
			timeZero := time.Time{}
			result.CreatedAt = timeZero
			result.UpdatedAt = timeZero
			tc.expectedOut.CreatedAt = timeZero
			tc.expectedOut.UpdatedAt = timeZero
			assert.Equal(t, tc.expectedOut, result)

			if tc.verifyFunc != nil {
				tc.verifyFunc(db, tc.expectedOut)
			}
		})
	}
}
