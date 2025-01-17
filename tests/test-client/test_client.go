package test_client

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/aspirin100/gRPC-SSO/internal/config"
	ssov1 "github.com/aspirin100/gRPC-SSO/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type TestClient struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *TestClient) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	// for control if test is too long
	ctx, cancelFunc := context.WithTimeout(
		context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelFunc()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithIdleTimeout(cfg.GRPC.Timeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &TestClient{
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
