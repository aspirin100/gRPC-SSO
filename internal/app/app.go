package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/aspirin100/gRPC-SSO/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(
	logg *slog.Logger,
	port int,
	storagePath string,
	refreshTTL,
	accessTTL time.Duration,
) *App {
	grpcApp := grpcApp.New(logg, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
