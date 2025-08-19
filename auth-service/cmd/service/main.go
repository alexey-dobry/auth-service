package main

import (
	"github.com/alexey-dobry/auth-service/internal/app"
	"github.com/alexey-dobry/auth-service/internal/config"
	"github.com/alexey-dobry/auth-service/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.NewLogger(cfg.Env)

	logger.Info("Successfully initialized logger")

	application := app.New(logger, cfg.GRPC.Port, cfg.TockenTTL)

	application.GRPCSrv.Run()

}
