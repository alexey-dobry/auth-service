package auth

import (
	"context"
	"errors"
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
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Errorf("Failed to hash password: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		HashPassword: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsAdmin:      false,
	}

	if err = user.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user field value")
	}

	err = s.repository.Add(user)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil, status.Error(codes.AlreadyExists, "Account with specified email already exists")
	} else if err != nil {
		s.logger.Errorf("Failed to add new user to data: %s", err)
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
		s.logger.Errorf("Failed to generate token pair: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	response := pb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userCredentials := &model.UserCredentials{Email: req.Email, Password: req.Password}

	if err := userCredentials.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid login arguments")
	}

	user, err := s.repository.GetOneByMail(req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if !utils.CheckPasswordHash(req.Password, user.HashPassword) {
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
		s.logger.Errorf("Failed to generate token pair: %s", err)
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
		s.logger.Errorf("Failed validate refresh token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user, err := s.repository.GetOneByID(claims.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
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
		s.logger.Errorf("Failed validate access token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}
