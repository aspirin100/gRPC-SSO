package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/aspirin100/gRPC-SSO/internal/app/grpc"
	"github.com/aspirin100/gRPC-SSO/internal/service/auth"
	"github.com/aspirin100/gRPC-SSO/internal/storage/sqlite"
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
	storage, err := sqlite.New(logg, storagePath)
	if err != nil {
		panic(err)
	}

	// service layer constructor
	authService := auth.New(logg, storage, accessTTL, refreshTTL)

	// business logic layer constructor
	grpcApp := grpcApp.New(logg, authService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
