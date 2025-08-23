package auth

import (
	"context"
	"testing"
	"time"

	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	"github.com/alexey-dobry/auth-service/internal/domain/model"
	"github.com/alexey-dobry/auth-service/internal/domain/utils"
	"github.com/alexey-dobry/auth-service/internal/repository/mock"

	authv1 "github.com/alexey-dobry/auth-service/pkg/gen/go"
	"github.com/alexey-dobry/auth-service/pkg/logger/zap"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// hashPassword, err := utils.HashPassword("GOOOOlang123")
	// assert.Equal(nil, err)

	// user := model.User{
	// 	Username:     "Klaasje 51",
	// 	Email:        "KlausShwab@gmail.com",
	// 	HashPassword: hashPassword,
	// 	FirstName:    "Klaasje",
	// 	LastName:     "Amadeus",
	// 	IsAdmin:      false,
	// }

	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().Add(gomock.Any()).Return(nil)

	jh, err := jwt.NewHandler(jwt.Config{
		AccessSecret:  "privateKey",
		RefreshSecret: "veryPrivateKey",
		TTL: jwt.TTL{
			AccessTTL:  time.Minute * 120,
			RefreshTTL: time.Minute * 10000,
		},
	})
	assert.Equal(nil, err)

	logger := zap.NewLogger(zap.Config{})

	s := New(logger, mockUserRepository, jh)

	req := &authv1.RegisterRequest{
		Username:  "Klaasje 51",
		Email:     "KlausShwab@gmail.com",
		Password:  "GOOOOlang123",
		FirstName: "Klaasje",
		LastName:  "Amadeus",
	}

	_, err = s.Register(ctx, req)
	assert.Equal(nil, err)
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	password := "GOOOOlang123"

	hashPassword, err := utils.HashPassword(password)
	assert.Equal(nil, err)

	user := model.User{
		Username:     "Klaasje 51",
		Email:        "KlausShwab@gmail.com",
		HashPassword: hashPassword,
		FirstName:    "Klaasje",
		LastName:     "Amadeus",
		IsAdmin:      false,
	}

	mockUserRepository := mock.NewMockUserRepository(ctrl)
	mockUserRepository.EXPECT().GetOne("KlausShwab@gmail.com").Return(user, nil)

	jh, err := jwt.NewHandler(jwt.Config{
		AccessSecret:  "privateKey",
		RefreshSecret: "veryPrivateKey",
		TTL: jwt.TTL{
			AccessTTL:  time.Minute * 120,
			RefreshTTL: time.Minute * 10000,
		},
	})
	assert.Equal(nil, err)

	logger := zap.NewLogger(zap.Config{})

	s := New(logger, mockUserRepository, jh)

	req := &authv1.LoginRequest{
		Email:    "KlausShwab@gmail.com",
		Password: password,
	}

	_, err = s.Login(ctx, req)
	assert.Equal(nil, err)
}

func TestRefresh(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepository := mock.NewMockUserRepository(ctrl)

	jh, err := jwt.NewHandler(jwt.Config{
		AccessSecret:  "privateKey",
		RefreshSecret: "veryPrivateKey",
		TTL: jwt.TTL{
			AccessTTL:  time.Minute * 120,
			RefreshTTL: time.Minute * 10000,
		},
	})
	assert.Equal(nil, err)

	logger := zap.NewLogger(zap.Config{})

	s := New(logger, mockUserRepository, jh)

	refreshToken, _, err := jh.GenerateJWTPair(jwt.Claims{
		ID:        1,
		Username:  "Klaasje 51",
		FirstName: "Klaasje",
		LastName:  "Amadeus",
		IsAdmin:   false,
	})

	req := &authv1.RefreshRequest{
		RefreshToken: refreshToken,
	}

	_, err = s.Refresh(ctx, req)

	assert.Equal(nil, err)
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepository := mock.NewMockUserRepository(ctrl)

	jh, err := jwt.NewHandler(jwt.Config{
		AccessSecret:  "privateKey",
		RefreshSecret: "veryPrivateKey",
		TTL: jwt.TTL{
			AccessTTL:  time.Minute * 120,
			RefreshTTL: time.Minute * 10000,
		},
	})
	assert.Equal(nil, err)

	logger := zap.NewLogger(zap.Config{})

	s := New(logger, mockUserRepository, jh)

	_, accessToken, err := jh.GenerateJWTPair(jwt.Claims{
		ID:        1,
		Username:  "Klaasje 51",
		FirstName: "Klaasje",
		LastName:  "Amadeus",
		IsAdmin:   false,
	})

	req := &authv1.ValidateRequest{
		AccessToken: accessToken,
	}

	_, err = s.Validate(ctx, req)

	assert.Equal(nil, err)
}
