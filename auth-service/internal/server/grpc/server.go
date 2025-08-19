package grpc

import "context"

type ServerAPI interface {
	Run(context.Context) error

	Stop() error
}
