package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Account        string `gorm:"uniqueIndex"`
	HashedPassword string
	Name           string
}

func (u *User) TableName() string {
	return "users"
}
