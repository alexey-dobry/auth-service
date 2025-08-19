package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/alexey-dobry/auth-service/internal/server/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.SugaredLogger
	gRPCServer *grpc.Server
	port       string
}

func New(log *zap.SugaredLogger, port string) *App {
	grpcServer := grpc.NewServer()

	authgrpc.Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	log := a.log

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Infof("grpc server is running; addr: %s", l.Addr().String())

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
