package grpcApp

import (
	"fmt"
	"log/slog"
	"net"

	grpcAuth "github.com/aspirin100/gRPC-SSO/internal/grpc/auth"
	"github.com/aspirin100/gRPC-SSO/pkg/logger/sl"
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

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"
	logg := a.logg.With(slog.String("op", op))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		logg.Error("net.Listen error: %w", sl.Err(err))
		return fmt.Errorf("%s:%w", op, err)
	}

	logg.Info("grpc server is running",
		slog.String("addr", listener.Addr().String()))

	err = a.gRPCServer.Serve(listener)
	if err != nil {
		logg.Error("grpc serving error: %w", sl.Err(err))
		return fmt.Errorf("failed to run grpc server: %w", err)
	}

	return nil
}

func (a *App) GracefulStop() {
	const op = "grpcApp.Stop"
	logg := a.logg.With(slog.String("op", op))

	a.gRPCServer.GracefulStop()

	logg.Info("grpc server stopped", slog.Int("port", a.port))
}
