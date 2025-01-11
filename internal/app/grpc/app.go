package grpc

import (
	"log/slog"

	grpcAuth "github.com/aspirin100/gRPC-SSO/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	logg       *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logg *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	grpcAuth.Register(gRPCServer)

	return &App{
		logg:       logg,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
