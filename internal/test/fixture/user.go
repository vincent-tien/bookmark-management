package fixture

import (
	"github.com/vincent-tien/bookmark-management/internal/model"
	"gorm.io/gorm"
)

type UserFixture struct {
	db *gorm.DB
}

func (u *UserFixture) SetupDB(db *gorm.DB) {
	u.db = db
}

func (u *UserFixture) DB() *gorm.DB {
	return u.db
}

func (u *UserFixture) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserFixture) GenerateData() error {
	db := u.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			ID:          "deb745af-1a62-4efa-99a0-f06b274bd993",
			DisplayName: "John Doe",
			Username:    "John Doe",
			Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
			Email:       "john.doe@example.com",
		},
		{
			ID:          "deb745af-1a62-4efa-99a0-f06b274bd994",
			DisplayName: "Jane Doe",
			Username:    "Jane Doe",
			Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
			Email:       "jane.doe@example.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}
