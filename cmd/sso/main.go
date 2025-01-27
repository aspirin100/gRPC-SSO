package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aspirin100/gRPC-SSO/internal/app"
	"github.com/aspirin100/gRPC-SSO/internal/config"
	"github.com/aspirin100/gRPC-SSO/pkg/logger/sl"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logg := setupLogger(cfg.Env)
	logg.Info("logger setuped", slog.String("env", cfg.Env))

	logg.Info("current secret key", slog.Int("length", len(cfg.SecretKey)))

	appConfig := app.NewAppConfig(cfg)

	application, err := app.New(logg, appConfig)
	if err != nil {
		logg.Debug("failed to create app instance", sl.Err(err))
		os.Exit(1)
	}

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
	default:
		log = slog.Default()
	}

	return log
}
