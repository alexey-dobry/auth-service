package pg

import (
	"github.com/alexey-dobry/auth-service/internal/domain/model"
	"github.com/alexey-dobry/auth-service/internal/repository"
	"github.com/google/uuid"
)

func (ur *UserRepository) Add(userData model.User) error {
	return ur.db.Create(userData).Error
}

func (ur *UserRepository) GetOne(email string) (model.User, error) {
	user := model.User{}

	result := ur.db.Select("password", "role", "id", "first_name", "last_name").Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (ur *UserRepository) UpdateUser(userId uuid.UUID, newData repository.UpdateUserParams) error {
	return ur.db.Model(model.User{}).Where("id = ?", userId.String()).Update("password", newData.NewPassword).Error
}

func (ur *UserRepository) Delete(userId uuid.UUID) error {
	return ur.db.Where("id = ?", userId.String()).Delete(&model.User{}).Error
}
