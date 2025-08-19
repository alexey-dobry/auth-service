package app

import (
	"time"

	grpcapp "github.com/alexey-dobry/auth-service/internal/app/grpc"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *zap.SugaredLogger, grpcPort string, tockenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
