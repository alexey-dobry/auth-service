package grpc

import (
	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	"github.com/alexey-dobry/auth-service/internal/repository"
	"github.com/alexey-dobry/auth-service/internal/server/grpc/auth"
	authrpc "github.com/alexey-dobry/auth-service/pkg/gen/go"
	"github.com/alexey-dobry/auth-service/pkg/logger"

	"google.golang.org/grpc"
)

func NewServer(logger logger.Logger, repository repository.UserRepository, jwtHandler jwt.JWTHandler) *grpc.Server {
	s := grpc.NewServer()

	authrpc.RegisterAuthServer(s, auth.New(logger, repository, jwtHandler))

	return s
}
