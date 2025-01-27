package grpcApp //nolint:stylecheck

import (
	"fmt"
	"log/slog"
	"net"

	grpcAuth "github.com/aspirin100/gRPC-SSO/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	logg       *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       int
}

func New(logg *slog.Logger, authService grpcAuth.Auth, host string, port int) *App {
	gRPCServer := grpc.NewServer()

	grpcAuth.RegisterAuthServer(gRPCServer, authService)

	return &App{
		logg:       logg,
		gRPCServer: gRPCServer,
		host:       host,
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

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	logg.Info("grpc server is running",
		slog.String("addr", listener.Addr().String()))

	err = a.gRPCServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to run grpc server: %w", err)
	}

	return nil
}

func (a *App) GracefulStop() {
	a.gRPCServer.GracefulStop()
}
