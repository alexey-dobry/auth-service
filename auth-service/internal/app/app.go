package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alexey-dobry/auth-service/internal/config"
	"github.com/alexey-dobry/auth-service/internal/repository"
	"github.com/alexey-dobry/auth-service/pkg/logger"
	"google.golang.org/grpc"
)

type App interface {
	Run(context.Context) error
}

type app struct {
	authServer    *grpc.Server
	serverAddress string
	repository    repository.UserRepository
	logger        logger.Logger
}

func New(cfg config.Config, logger logger.Logger, repository repository.UserRepository, authServer *grpc.Server) App {
	var a app

	a.logger = logger.WithFields("layer", "app")

	a.serverAddress = fmt.Sprintf(":%d", cfg.GRPC.Port)

	a.repository = repository

	a.authServer = authServer
	a.logger.Info("app was built")
	return &a
}

func (a *app) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	lis, err := net.Listen("tcp", a.serverAddress)
	if err != nil {
		a.logger.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		a.logger.Infof("starting grpc server at address %s...", a.serverAddress)
		if err := a.authServer.Serve(lis); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errMsg := fmt.Errorf("serving grpc server error: %s", err)
				a.logger.Error(errMsg)
				cancel()
			}
		}
	}()
	a.logger.Info("app running...")

	select {
	case <-quit:
		a.logger.Info("shutdown signal received")
	case <-ctx.Done():
		a.logger.Info("context canceled")
	}

	a.logger.Info("stopping all services")

	cancel()
	if err := lis.Close(); err != nil {
		a.logger.Warnf("net listener closing ended with error: %s", err)
	}
	wg.Wait()

	if err := a.repository.Close(); err != nil {
		a.logger.Warnf("store closing ended with error: %s", err)
	}

	a.logger.Info("app was gracefully shutdowned")
	return nil
}
