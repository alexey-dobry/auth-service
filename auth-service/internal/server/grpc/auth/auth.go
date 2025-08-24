package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	"github.com/alexey-dobry/auth-service/internal/domain/model"
	"github.com/alexey-dobry/auth-service/internal/domain/utils"
	pb "github.com/alexey-dobry/auth-service/pkg/gen/go"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, _ := utils.HashPassword(req.Password)

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		HashPassword: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsAdmin:      false,
	}

	err := user.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user field value")
	}

	err = s.repository.Add(user)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil, status.Error(codes.AlreadyExists, "Account with specified email already exists")
	} else if err != nil {
		errMsg := fmt.Sprintf("Error adding new user to data: %s", err)
		s.logger.Errorf(errMsg)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	refreshToken, accessToken, err := s.jwtHandler.GenerateJWTPair(jwt.Claims{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	})

	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate token pair: %s", err)
		s.logger.Errorf(errMsg)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	response := pb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	hashedPassword, _ := utils.HashPassword(req.Password)
	userMock := &model.User{Email: req.Email, HashPassword: hashedPassword}
	err := userMock.ValidateForLogin()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid login arguments")
	}

	user, err := s.repository.GetOneByMail(req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to get user data from database: %s", err)
		s.logger.Errorf(errMsg)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if !utils.CheckPasswordHash(req.Password, user.HashPassword) {
		s.logger.Info(user.ID)
		return nil, status.Error(codes.PermissionDenied, "Wrong password")
	}

	refreshToken, accessToken, err := s.jwtHandler.GenerateJWTPair(jwt.Claims{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	})

	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate token pair: %s", err)
		s.logger.Errorf(errMsg)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServerAPI) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	claims, err := s.jwtHandler.ValidateJWT(req.RefreshToken, jwt.RefreshToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		log.Print(err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user, err := s.repository.GetOneByID(claims.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to get user data from database: %s", err)
		s.logger.Errorf(errMsg)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	refreshToken, accessToken, err := s.jwtHandler.GenerateJWTPair(jwt.Claims{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	})

	return &pb.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *ServerAPI) Validate(ctx context.Context, req *pb.ValidateRequest) (*emptypb.Empty, error) {
	_, err := s.jwtHandler.ValidateJWT(req.AccessToken, jwt.AccessToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		log.Print(err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}
