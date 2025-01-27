package app

import (
	"fmt"
	"log/slog"
	"time"

	grpcApp "github.com/aspirin100/gRPC-SSO/internal/app/grpc"
	"github.com/aspirin100/gRPC-SSO/internal/config"
	"github.com/aspirin100/gRPC-SSO/internal/service/auth"
	"github.com/aspirin100/gRPC-SSO/internal/storage"
)

type App struct {
	GRPCServer *grpcApp.App
}

type AppConfig struct {
	host        string
	port        int
	storagePath string
	refreshTTL  time.Duration
	accessTTL   time.Duration
	secretKey   string
	reflection  bool
}

func New(
	logg *slog.Logger,
	cfg *AppConfig,
) (*App, error) {
	storage, err := storage.New(logg, cfg.storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to construct storage: %w", err)
	}

	// service layer constructor
	authService := auth.New(
		logg,
		storage,
		cfg.accessTTL, cfg.refreshTTL,
		cfg.secretKey)

	// business logic layer constructor
	grpcApplication := grpcApp.New(logg,
		authService, cfg.host, cfg.port, cfg.reflection)


	return &App{
		GRPCServer: grpcApplication,
	}, nil
}

func NewAppConfig(cfg *config.Config, reflection bool) *AppConfig {
	appCfg := &AppConfig{
		port:        cfg.GRPC.Port,
		storagePath: cfg.StoragePath,
		refreshTTL:  cfg.RefreshTTL,
		accessTTL:   cfg.AccessTTL,
		secretKey:   cfg.SecretKey,
		reflection:  reflection,
	}

	return appCfg
}
