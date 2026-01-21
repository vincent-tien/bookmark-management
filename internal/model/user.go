package model

import "time"

type User struct {
	ID          string `gorm:"type:uuid;primaryKey;column:id"`
	Username    string `gorm:"type:varchar(50);uniqueIndex;column:username"`
	Password    string `gorm:"column:password"`
	DisplayName string `gorm:"column:display_name;type:varchar(50);"`
	Email       string `gorm:"column:email;type:varchar(100);uniqueIndex" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
