package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system.
//
// It has the following fields:
// - ID: the unique identifier of the user (type: uuid).
// - Username: the username of the user (type: varchar(50); unique index).
// - Password: the hashed password of the user (type: varchar(100); non-null).
// - DisplayName: the display name of the user (type: varchar(50); non-null).
// - Email: the email address of the user (type: varchar(100); unique index; non-null).
// - CreatedAt: the timestamp when the user is created (type: timestamp with time zone; non-null).
// - UpdatedAt: the timestamp when the user is updated (type: timestamp with time zone; non-null).
type User struct {
	ID          string `gorm:"type:uuid;primaryKey;column:id"`
	Username    string `gorm:"type:varchar(50);uniqueIndex;column:username"`
	Password    string `gorm:"column:password"`
	DisplayName string `gorm:"column:display_name;type:varchar(50);"`
	Email       string `gorm:"column:email;type:varchar(100);uniqueIndex" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		userID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		u.ID = userID.String()
	}

	return nil
}
