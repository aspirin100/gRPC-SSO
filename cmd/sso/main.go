package main

import (
	"log/slog"
	"os"

	"github.com/aspirin100/gRPC-SSO/internal/app"
	"github.com/aspirin100/gRPC-SSO/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logg := setupLogger(cfg.Env)
	logg.Info("logger setuped", slog.String("env", cfg.Env))

	app := app.New(logg, cfg.GRPC.Port,
		cfg.StoragePath, cfg.RefreshTTL, cfg.AccessTTL)

	app.GRPCServer.Run()
	defer app.GRPCServer.GracefulStop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
