package main

import (
	"github.com/alexey-dobry/goauth/internal/config"
	"github.com/alexey-dobry/goauth/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.NewLogger(cfg.Env)

	logger.Info("Successfully initialized logger")

}
