package model

import (
	"github.com/alexey-dobry/auth-service/pkg/validator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"not null" validate:"required,min=1,max=200"`
	Email        string `gorm:"uniqueIndex" validate:"required,email"`
	HashPassword string `gorm:"not null" validate:"required"`
	FirstName    string `gorm:"not null" validate:"required,min=1,max=200"`
	LastName     string `gorm:"not null" validate:"required,min=1,max=200"`
	IsAdmin      bool   `gorm:"not null" validate:"required"`
}

func (u *User) Validate() error {
	return validator.V.Struct(u)
}

func (u *User) ValidateForLogin() error {
	var err error
	err = validator.V.Var(u.Email, "required,email")
	if err != nil {
		return err
	}
	err = validator.V.Var(u.HashPassword, "required")
	if err != nil {
		return err
	}
	return err
}
