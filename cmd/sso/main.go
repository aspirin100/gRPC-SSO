package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	logg.Info("current secret key", slog.String("sKey", cfg.SecretKey))

	application := app.New(logg, cfg.GRPC.Port,
		cfg.StoragePath, cfg.RefreshTTL, cfg.AccessTTL, cfg.SecretKey)

	go application.GRPCServer.MustRun()

	// graceful stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.GracefulStop()

	logg.Info("sso server stopped")
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
