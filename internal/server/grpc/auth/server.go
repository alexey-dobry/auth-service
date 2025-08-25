package auth

import (
	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	"github.com/alexey-dobry/auth-service/internal/repository"
	pb "github.com/alexey-dobry/auth-service/pkg/gen/go"
	"github.com/alexey-dobry/auth-service/pkg/logger"
)

type ServerAPI struct {
	pb.UnimplementedAuthServer

	logger     logger.Logger
	repository repository.UserRepository
	jwtHandler jwt.JWTHandler
}

func New(logger logger.Logger, repository repository.UserRepository, jwtHandler jwt.JWTHandler) *ServerAPI {
	return &ServerAPI{
		repository: repository,
		logger:     logger.WithFields("layer", "grpc server api", "server", "manager"),
		jwtHandler: jwtHandler,
	}
}
