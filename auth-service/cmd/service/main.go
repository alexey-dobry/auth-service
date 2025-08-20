package main

import (
	"context"
	"flag"

	"github.com/alexey-dobry/auth-service/internal/app"
	"github.com/alexey-dobry/auth-service/internal/config"
	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	pg "github.com/alexey-dobry/auth-service/internal/repository/postgresql"
	"github.com/alexey-dobry/auth-service/internal/server/grpc"
	"github.com/alexey-dobry/auth-service/pkg/logger/zap"
)

func main() {
	flag.Parse()

	cfg := config.MustLoad()

	logger := zap.NewLogger(cfg.Logger)

	logger.Info("Successfully initialized logger")

	repository, err := pg.New(cfg.Repository)
	if err != nil {
		logger.Error("Failed to create user repository: %s", err)
		panic("Failed to create user repository")
	}

	jwtHandler, err := jwt.NewHandler(cfg.JWT)
	if err != nil {
		logger.Error("Failed to create jwt handler: %s", err)
		panic("Failed to create jwt handler")
	}

	authServer := grpc.NewServer(logger, repository, jwtHandler)

	application := app.New(cfg, logger, repository, authServer)

	application.Run(context.Background())
}
