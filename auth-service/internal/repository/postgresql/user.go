package pg

import (
	"github.com/alexey-dobry/auth-service/internal/domain/model"
	"github.com/google/uuid"
)

func (ur *UserRepository) Add(userData model.User) error {
	return ur.db.Create(&userData).Error
}

func (ur *UserRepository) GetOneByMail(email string) (model.User, error) {
	user := model.User{}

	result := ur.db.Select("username", "hash_password", "first_name", "last_name", "is_admin").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (ur *UserRepository) GetOneByID(ID uint) (model.User, error) {
	user := model.User{}

	result := ur.db.Select("username", "first_name", "last_name", "is_admin").Where("id = ?", ID).First(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (ur *UserRepository) Delete(userId uuid.UUID) error {
	return ur.db.Where("id = ?", userId.String()).Delete(&model.User{}).Error
}
