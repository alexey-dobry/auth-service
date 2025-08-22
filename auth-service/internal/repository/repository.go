package repository

import (
	"github.com/alexey-dobry/auth-service/internal/domain/model"
)

type UserRepository interface {
	Add(model.User) error

	GetOne(email string) (model.User, error)

	Close() error
}
