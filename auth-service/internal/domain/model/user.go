package model

import (
	"github.com/alexey-dobry/auth-service/pkg/validator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"not null" validate:"required,min=1,max=200"`
	Email        string `gorm:"uniqueIndex" validate:"required,email"`
	HashPassword string `gorm:"not null" validate:"required,sha512"`
	FirstName    string `gorm:"not null" validate:"required,min=1,max=200"`
	LastName     string `gorm:"not null" validate:"required,min=1,max=200"`
	IsAdmin      bool   `gorm:"not null" validate:"required"`
}

func (u *User) Validate() error {
	return validator.V.Struct(u)
}
