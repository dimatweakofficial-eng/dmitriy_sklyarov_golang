package user

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email    string `gorm:"index"`
	Password string
	Name     string
}
