package pg

import (
	"fmt"
	"time"

	"github.com/alexey-dobry/auth-service/internal/domain/model"
	"github.com/alexey-dobry/auth-service/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 10
const delay = 2 * time.Second

type UserRepository struct {
	db *gorm.DB
}

func New(cfg Config) (repository.UserRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DatabaseName, cfg.Port)

	var db *gorm.DB
	var err error
	for range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn))
		if err == nil {
			break
		}

		time.Sleep(delay)
	}
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(model.User{})
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		db: db,
	}, nil
}

func (ur *UserRepository) Close() error {
	sqlDB, _ := ur.db.DB()
	err := sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
