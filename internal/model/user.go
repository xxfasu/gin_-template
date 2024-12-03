package model

import (
	"gorm.io/gorm"
)

type User struct {
	Id        uint           `gorm:"primarykey"`
	UserId    string         `gorm:"unique;not null"`
	Nickname  string         `gorm:"not null"`
	Password  string         `gorm:"not null"`
	Email     string         `gorm:"not null"`
	CreatedAt int64          `gorm:"create_at"`
	UpdatedAt int64          `gorm:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"delete_at"`
}

func (u *User) TableName() string {
	return "users"
}
