package repository

import (
	"github.com/alexey-dobry/auth-service/internal/domain/model"
)

type UserRepository interface {
	Add(model.User) error

	GetOneByMail(email string) (model.User, error)

	GetOneByID(ID uint) (model.User, error)

	Close() error
}
